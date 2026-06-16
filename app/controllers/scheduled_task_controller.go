package controllers

import (
	"database/sql"
	"strconv"
	"time"

	ghttp "github.com/Hlgxz/gai/http"
	"github.com/Hlgxz/gai/database/orm"
)

// ScheduledTaskController 定时任务控制器
type ScheduledTaskController struct {
	DB        *orm.DB
	Scheduler interface {
		RunNow(taskID uint) error
	}
}

// TaskResponse 任务响应结构
type TaskResponse struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	CronExpr    string `json:"cron_expr"`
	Command     string `json:"command"`
	Params      string `json:"params"`
	Status      string `json:"status"`
	Priority    int    `json:"priority"`
	MaxRetries  int    `json:"max_retries"`
	Timeout     int    `json:"timeout"`
	LastRunAt   string `json:"last_run_at,omitempty"`
	NextRunAt   string `json:"next_run_at,omitempty"`
	RunCount    int    `json:"run_count"`
	FailCount   int    `json:"fail_count"`
	LastError   string `json:"last_error,omitempty"`
	CreatedAt   string `json:"created_at"`
}

// List 任务列表
func (tc *ScheduledTaskController) List(c *ghttp.Context) {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	keyword := c.Query("keyword", "")
	status := c.Query("status", "")

	query := `SELECT id, name, code, COALESCE(description, ''), cron_expr,
	          COALESCE(command, ''), COALESCE(params, ''), status, priority,
	          max_retries, timeout, COALESCE(last_run_at, ''), COALESCE(next_run_at, ''),
	          run_count, fail_count, COALESCE(last_error, ''), created_at
	          FROM scheduled_tasks WHERE deleted_at IS NULL`
	countQuery := `SELECT COUNT(*) FROM scheduled_tasks WHERE deleted_at IS NULL`
	args := []any{}
	countArgs := []any{}

	if keyword != "" {
		query += " AND name LIKE ?"
		countQuery += " AND name LIKE ?"
		args = append(args, "%"+keyword+"%")
		countArgs = append(countArgs, "%"+keyword+"%")
	}

	if status != "" {
		query += " AND status = ?"
		countQuery += " AND status = ?"
		args = append(args, status)
		countArgs = append(countArgs, status)
	}

	query += " ORDER BY id DESC LIMIT ? OFFSET ?"
	args = append(args, pageSize, offset)

	rows, err := tc.DB.SQL.Query(query, args...)
	if err != nil {
		c.Error(500, "查询失败: "+err.Error())
		return
	}
	defer rows.Close()

	var items []TaskResponse
	for rows.Next() {
		var item TaskResponse
		err := rows.Scan(
			&item.ID, &item.Name, &item.Code, &item.Description, &item.CronExpr,
			&item.Command, &item.Params, &item.Status, &item.Priority,
			&item.MaxRetries, &item.Timeout, &item.LastRunAt, &item.NextRunAt,
			&item.RunCount, &item.FailCount, &item.LastError, &item.CreatedAt,
		)
		if err != nil {
			continue
		}
		items = append(items, item)
	}

	var total int
	tc.DB.SQL.QueryRow(countQuery, countArgs...).Scan(&total)

	c.Success(map[string]any{
		"items":     items,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// Get 获取任务详情
func (tc *ScheduledTaskController) Get(c *ghttp.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Error(400, "无效的任务ID")
		return
	}

	var item TaskResponse
	err = tc.DB.SQL.QueryRow(`
		SELECT id, name, code, COALESCE(description, ''), cron_expr,
		       COALESCE(command, ''), COALESCE(params, ''), status, priority,
		       max_retries, timeout, COALESCE(last_run_at, ''), COALESCE(next_run_at, ''),
		       run_count, fail_count, COALESCE(last_error, ''), created_at
		FROM scheduled_tasks WHERE id = ? AND deleted_at IS NULL
	`, id).Scan(
		&item.ID, &item.Name, &item.Code, &item.Description, &item.CronExpr,
		&item.Command, &item.Params, &item.Status, &item.Priority,
		&item.MaxRetries, &item.Timeout, &item.LastRunAt, &item.NextRunAt,
		&item.RunCount, &item.FailCount, &item.LastError, &item.CreatedAt,
	)

	if err == sql.ErrNoRows {
		c.Error(404, "任务不存在")
		return
	}
	if err != nil {
		c.Error(500, "查询失败")
		return
	}

	c.Success(item)
}

// Create 创建任务
func (tc *ScheduledTaskController) Create(c *ghttp.Context) {
	var input struct {
		Name        string `json:"name"`
		Code        string `json:"code"`
		Description string `json:"description"`
		CronExpr    string `json:"cron_expr"`
		Command     string `json:"command"`
		Params      string `json:"params"`
		Status      string `json:"status"`
		Priority    int    `json:"priority"`
		MaxRetries  int    `json:"max_retries"`
		Timeout     int    `json:"timeout"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.Error(400, "无效的请求数据")
		return
	}

	if input.Name == "" || input.Code == "" || input.CronExpr == "" {
		c.Error(400, "名称、代码和Cron表达式不能为空")
		return
	}

	// 检查代码是否重复
	var exists int
	tc.DB.SQL.QueryRow(`SELECT 1 FROM scheduled_tasks WHERE code = ? AND deleted_at IS NULL`, input.Code).Scan(&exists)
	if exists == 1 {
		c.Error(400, "任务代码已存在")
		return
	}

	if input.Status == "" {
		input.Status = "enabled"
	}
	if input.MaxRetries == 0 {
		input.MaxRetries = 3
	}
	if input.Timeout == 0 {
		input.Timeout = 3600
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	result, err := tc.DB.SQL.Exec(`
		INSERT INTO scheduled_tasks (name, code, description, cron_expr, command, params,
			status, priority, max_retries, timeout, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, input.Name, input.Code, input.Description, input.CronExpr, input.Command, input.Params,
		input.Status, input.Priority, input.MaxRetries, input.Timeout, now, now)
	if err != nil {
		c.Error(500, "创建失败: "+err.Error())
		return
	}

	id, _ := result.LastInsertId()

	c.Success(map[string]any{
		"id":      id,
		"message": "创建成功",
	})
}

// Update 更新任务
func (tc *ScheduledTaskController) Update(c *ghttp.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Error(400, "无效的任务ID")
		return
	}

	var input map[string]any
	if err := c.BindJSON(&input); err != nil {
		c.Error(400, "无效的请求数据")
		return
	}

	var exists int
	tc.DB.SQL.QueryRow(`SELECT 1 FROM scheduled_tasks WHERE id = ? AND deleted_at IS NULL`, id).Scan(&exists)
	if exists == 0 {
		c.Error(404, "任务不存在")
		return
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	if name, ok := input["name"].(string); ok {
		tc.DB.SQL.Exec(`UPDATE scheduled_tasks SET name = ?, updated_at = ? WHERE id = ?`, name, now, id)
	}
	if description, ok := input["description"].(string); ok {
		tc.DB.SQL.Exec(`UPDATE scheduled_tasks SET description = ?, updated_at = ? WHERE id = ?`, description, now, id)
	}
	if cronExpr, ok := input["cron_expr"].(string); ok {
		tc.DB.SQL.Exec(`UPDATE scheduled_tasks SET cron_expr = ?, updated_at = ? WHERE id = ?`, cronExpr, now, id)
	}
	if status, ok := input["status"].(string); ok {
		tc.DB.SQL.Exec(`UPDATE scheduled_tasks SET status = ?, updated_at = ? WHERE id = ?`, status, now, id)
	}

	c.Success(map[string]string{"message": "更新成功"})
}

// Delete 删除任务
func (tc *ScheduledTaskController) Delete(c *ghttp.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Error(400, "无效的任务ID")
		return
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	result, err := tc.DB.SQL.Exec(`
		UPDATE scheduled_tasks SET deleted_at = ? WHERE id = ? AND deleted_at IS NULL
	`, now, id)
	if err != nil {
		c.Error(500, "删除失败")
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.Error(404, "任务不存在")
		return
	}

	c.Success(map[string]string{"message": "删除成功"})
}

// Toggle 启用/禁用任务
func (tc *ScheduledTaskController) Toggle(c *ghttp.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Error(400, "无效的任务ID")
		return
	}

	var currentStatus string
	err = tc.DB.SQL.QueryRow(`SELECT status FROM scheduled_tasks WHERE id = ? AND deleted_at IS NULL`, id).Scan(&currentStatus)
	if err == sql.ErrNoRows {
		c.Error(404, "任务不存在")
		return
	}

	newStatus := "enabled"
	if currentStatus == "enabled" {
		newStatus = "disabled"
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	tc.DB.SQL.Exec(`UPDATE scheduled_tasks SET status = ?, updated_at = ? WHERE id = ?`, newStatus, now, id)

	c.Success(map[string]string{
		"message": "状态已更新",
		"status":  newStatus,
	})
}

// RunNow 立即执行任务
func (tc *ScheduledTaskController) RunNow(c *ghttp.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.Error(400, "无效的任务ID")
		return
	}

	var taskName string
	err = tc.DB.SQL.QueryRow(`SELECT name FROM scheduled_tasks WHERE id = ? AND deleted_at IS NULL`, id).Scan(&taskName)
	if err == sql.ErrNoRows {
		c.Error(404, "任务不存在")
		return
	}

	// 创建执行记录
	now := time.Now().Format("2006-01-02 15:04:05")
	result, _ := tc.DB.SQL.Exec(`
		INSERT INTO task_executions (task_id, task_name, status, started_at, triggered_by, created_at, updated_at)
		VALUES (?, ?, 'running', ?, 'manual', ?, ?)
	`, id, taskName, now, now, now)

	execID, _ := result.LastInsertId()

	// 模拟执行
	go func() {
		time.Sleep(1 * time.Second)
		finishedAt := time.Now().Format("2006-01-02 15:04:05")
		tc.DB.SQL.Exec(`
			UPDATE task_executions SET status = 'success', finished_at = ?, duration = 1000, updated_at = ?
			WHERE id = ?
		`, finishedAt, finishedAt, execID)

		tc.DB.SQL.Exec(`
			UPDATE scheduled_tasks SET last_run_at = ?, run_count = run_count + 1, updated_at = ?
			WHERE id = ?
		`, now, now, id)
	}()

	c.Success(map[string]any{
		"message":      "任务已触发执行",
		"execution_id": execID,
	})
}

// Executions 任务执行记录列表
func (tc *ScheduledTaskController) Executions(c *ghttp.Context) {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	taskID := c.Query("task_id", "")

	query := `SELECT id, task_id, task_name, status, started_at, COALESCE(finished_at, ''),
	          duration, COALESCE(output, ''), COALESCE(error_msg, ''), retry_count, triggered_by, created_at
	          FROM task_executions WHERE deleted_at IS NULL`
	args := []any{}

	if taskID != "" {
		query += " AND task_id = ?"
		taskIDInt, _ := strconv.ParseUint(taskID, 10, 64)
		args = append(args, taskIDInt)
	}

	query += " ORDER BY id DESC LIMIT ? OFFSET ?"
	args = append(args, pageSize, offset)

	rows, err := tc.DB.SQL.Query(query, args...)
	if err != nil {
		c.Error(500, "查询失败")
		return
	}
	defer rows.Close()

	var items []map[string]any
	for rows.Next() {
		var id, taskID, duration, retryCount int
		var taskName, status, startedAt, finishedAt, output, errorMsg, triggeredBy, createdAt string
		rows.Scan(&id, &taskID, &taskName, &status, &startedAt, &finishedAt,
			&duration, &output, &errorMsg, &retryCount, &triggeredBy, &createdAt)
		items = append(items, map[string]any{
			"id":           id,
			"task_id":      taskID,
			"task_name":    taskName,
			"status":       status,
			"started_at":   startedAt,
			"finished_at":  finishedAt,
			"duration":     duration,
			"output":       output,
			"error_msg":    errorMsg,
			"retry_count":  retryCount,
			"triggered_by": triggeredBy,
			"created_at":   createdAt,
		})
	}

	c.Success(map[string]any{
		"items":     items,
		"page":      page,
		"page_size": pageSize,
	})
}

// Stats 任务统计
func (tc *ScheduledTaskController) Stats(c *ghttp.Context) {
	var totalTasks, enabledTasks int
	tc.DB.SQL.QueryRow(`SELECT COUNT(*) FROM scheduled_tasks WHERE deleted_at IS NULL`).Scan(&totalTasks)
	tc.DB.SQL.QueryRow(`SELECT COUNT(*) FROM scheduled_tasks WHERE status = 'enabled' AND deleted_at IS NULL`).Scan(&enabledTasks)

	today := time.Now().Format("2006-01-02")
	var todayExecutions int
	tc.DB.SQL.QueryRow(`SELECT COUNT(*) FROM task_executions WHERE started_at LIKE ? AND deleted_at IS NULL`, today+"%").Scan(&todayExecutions)

	var todaySuccess int
	tc.DB.SQL.QueryRow(`SELECT COUNT(*) FROM task_executions WHERE started_at LIKE ? AND status = 'success' AND deleted_at IS NULL`, today+"%").Scan(&todaySuccess)

	successRate := float64(0)
	if todayExecutions > 0 {
		successRate = float64(todaySuccess) / float64(todayExecutions) * 100
	}

	c.Success(map[string]any{
		"total_tasks":      totalTasks,
		"enabled_tasks":    enabledTasks,
		"today_executions": todayExecutions,
		"today_success":    todaySuccess,
		"success_rate":     successRate,
	})
}