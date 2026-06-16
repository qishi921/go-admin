package scheduler

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Hlgxz/gai/database/orm"
	"github.com/user/admin-system/app/models"
)

// TaskFunc 任务执行函数类型
type TaskFunc func(params map[string]any) error

// Scheduler 定时任务调度器
type Scheduler struct {
	db         *orm.DB
	tasks      map[string]TaskFunc
	running    bool
	mu         sync.RWMutex
	stopChan   chan struct{}
}

// NewScheduler 创建调度器
func NewScheduler(db *orm.DB) *Scheduler {
	return &Scheduler{
		db:       db,
		tasks:    make(map[string]TaskFunc),
		stopChan: make(chan struct{}),
	}
}

// RegisterTask 注册任务处理函数
func (s *Scheduler) RegisterTask(code string, fn TaskFunc) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks[code] = fn
}

// Start 启动调度器
func (s *Scheduler) Start() {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return
	}
	s.running = true
	s.mu.Unlock()

	go s.runLoop()
	log.Println("[Scheduler] Started")
}

// Stop 停止调度器
func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.running {
		s.stopChan <- struct{}{}
		s.running = false
		log.Println("[Scheduler] Stopped")
	}
}

// runLoop 调度循环
func (s *Scheduler) runLoop() {
	// 每分钟检查一次
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.checkAndRun()
		case <-s.stopChan:
			return
		}
	}
}

// checkAndRun 检查并执行到期任务
func (s *Scheduler) checkAndRun() {
	// 获取所有启用的任务
	query := orm.Query[models.ScheduledTask](s.db).
		Where("status", "=", "enabled")
	tasks, err := orm.Get[models.ScheduledTask](query)
	if err != nil {
		log.Printf("[Scheduler] Failed to get tasks: %v", err)
		return
	}

	now := time.Now()
	for _, task := range tasks {
		if s.shouldRun(task, now) {
			go s.executeTask(task, "scheduler")
		}
	}
}

// shouldRun 判断任务是否应该执行
func (s *Scheduler) shouldRun(task models.ScheduledTask, now time.Time) bool {
	if task.NextRunAt == nil {
		return false
	}

	// 解析下次执行时间
	nextRun, err := time.Parse("2006-01-02 15:04:05", *task.NextRunAt)
	if err != nil {
		return false
	}

	// 检查是否到达执行时间（允许1分钟误差）
	return now.After(nextRun) || now.Sub(nextRun) < time.Minute
}

// executeTask 执行任务
func (s *Scheduler) executeTask(task models.ScheduledTask, triggeredBy string) {
	startTime := time.Now()
	nowStr := startTime.Format("2006-01-02 15:04:05")

	// 创建执行记录
	execution := &models.TaskExecution{
		TaskID:      task.ID,
		TaskName:    task.Name,
		Status:      "running",
		StartedAt:   nowStr,
		TriggeredBy: triggeredBy,
	}
	created, err := orm.Create[models.TaskExecution](s.db, execution)
	if err != nil {
		log.Printf("[Scheduler] Failed to create execution record: %v", err)
		return
	}

	// 获取任务处理函数
	s.mu.RLock()
	taskFunc, exists := s.tasks[task.Code]
	s.mu.RUnlock()

	var execErr error
	var output string

	if exists {
		// 执行任务
		var params map[string]any
		// 这里简化处理，实际需要解析 JSON
		execErr = taskFunc(params)
		if execErr != nil {
			output = execErr.Error()
		}
	} else {
		execErr = fmt.Errorf("task handler not registered: %s", task.Code)
	}

	// 更新执行记录
	endTime := time.Now()
	duration := int(endTime.Sub(startTime).Milliseconds())
	finishedAt := endTime.Format("2006-01-02 15:04:05")

	created.Status = "success"
	if execErr != nil {
		created.Status = "failed"
		created.ErrorMsg = execErr.Error()
	}
	created.FinishedAt = &finishedAt
	created.Duration = duration
	created.Output = output
	orm.Update[models.TaskExecution](s.db, created)

	// 更新任务状态
	task.LastRunAt = &nowStr
	task.RunCount++
	if execErr != nil {
		task.FailCount++
		task.LastError = execErr.Error()
	}
	orm.Update[models.ScheduledTask](s.db, &task)

	log.Printf("[Scheduler] Task '%s' executed, status: %s, duration: %dms", task.Name, created.Status, duration)
}

// RunNow 立即执行任务
func (s *Scheduler) RunNow(taskID uint) error {
	query := orm.Query[models.ScheduledTask](s.db).Where("id", "=", taskID)
	task, err := orm.First[models.ScheduledTask](query)
	if err != nil || task == nil {
		return fmt.Errorf("task not found")
	}

	go s.executeTask(*task, "manual")
	return nil
}

// ParseCron 解析 cron 表达式并计算下次执行时间
// 支持标准5字段格式: minute hour day month weekday
func ParseCron(cronExpr string, from time.Time) (time.Time, error) {
	// 简化实现：支持基本的 cron 表达式
	// 格式: minute hour day month weekday
	// 例如: "0 9 * * *" 表示每天早上9点

	fields := splitCronFields(cronExpr)
	if len(fields) != 5 {
		return time.Time{}, fmt.Errorf("invalid cron expression: expected 5 fields")
	}

	// 从下一分钟开始计算
	next := from.Add(time.Minute).Truncate(time.Minute)

	// 最多尝试1000次查找下次执行时间
	for i := 0; i < 1000; i++ {
		if matchesCron(next, fields) {
			return next, nil
		}
		next = next.Add(time.Minute)
	}

	return time.Time{}, fmt.Errorf("cannot find next execution time")
}

// splitCronFields 分割 cron 字段
func splitCronFields(expr string) []string {
	var fields []string
	start := 0
	for i := 0; i <= len(expr); i++ {
		if i == len(expr) || expr[i] == ' ' {
			if i > start {
				fields = append(fields, expr[start:i])
			}
			start = i + 1
		}
	}
	return fields
}

// matchesCron 检查时间是否匹配 cron 表达式
func matchesCron(t time.Time, fields []string) bool {
	return matchField(fields[0], t.Minute()) &&
		matchField(fields[1], t.Hour()) &&
		matchField(fields[2], t.Day()) &&
		matchField(fields[3], int(t.Month())) &&
		matchField(fields[4], int(t.Weekday()))
}

// matchField 检查单个字段是否匹配
func matchField(field string, value int) bool {
	if field == "*" {
		return true
	}

	// 处理逗号分隔的多个值
	hasComma := false
	for _, c := range field {
		if c == ',' {
			hasComma = true
			break
		}
	}

	if hasComma {
		start := 0
		for i := 0; i <= len(field); i++ {
			if i == len(field) || field[i] == ',' {
				if i > start {
					if matchSingleValue(field[start:i], value) {
						return true
					}
				}
				start = i + 1
			}
		}
		return false
	}

	return matchSingleValue(field, value)
}

// matchSingleValue 匹配单个值
func matchSingleValue(field string, value int) bool {
	// 处理范围 a-b
	for i := 0; i < len(field); i++ {
		if field[i] == '-' {
			start := parseInt(field[:i])
			end := parseInt(field[i+1:])
			return value >= start && value <= end
		}
	}

	// 处理步长 */n
	if len(field) >= 2 && field[0] == '/' {
		step := parseInt(field[1:])
		return value%step == 0
	}

	// 处理 */*
	if field == "*" {
		return true
	}

	// 精确匹配
	return parseInt(field) == value
}

// parseInt 解析整数
func parseInt(s string) int {
	var result int
	for _, c := range s {
		if c >= '0' && c <= '9' {
			result = result*10 + int(c-'0')
		}
	}
	return result
}

// CalculateNextRun 计算并更新任务的下次执行时间
func (s *Scheduler) CalculateNextRun(task *models.ScheduledTask) error {
	nextRun, err := ParseCron(task.CronExpr, time.Now())
	if err != nil {
		return err
	}

	nextRunStr := nextRun.Format("2006-01-02 15:04:05")
	task.NextRunAt = &nextRunStr
	return nil
}
