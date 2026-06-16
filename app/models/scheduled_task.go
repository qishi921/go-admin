package models

// ScheduledTask 定时任务模型
type ScheduledTask struct {
	BaseModel
	Name        string  `json:"name" gai:"column:name;size:200"`
	Code        string  `json:"code" gai:"column:code;size:100;unique"`
	Description string  `json:"description" gai:"column:description;type:text"`
	CronExpr    string  `json:"cron_expr" gai:"column:cron_expr;size:100"` // cron 表达式
	Command     string  `json:"command" gai:"column:command;type:text"`    // 执行的命令或任务标识
	Params      string  `json:"params" gai:"column:params;type:text"`      // JSON 参数
	Status      string  `json:"status" gai:"column:status;size:20;default:enabled"` // enabled, disabled
	Priority    int     `json:"priority" gai:"column:priority;default:0"`
	MaxRetries  int     `json:"max_retries" gai:"column:max_retries;default:3"`
	Timeout     int     `json:"timeout" gai:"column:timeout;default:3600"` // 超时秒数
	LastRunAt   *string `json:"last_run_at" gai:"column:last_run_at"`
	NextRunAt   *string `json:"next_run_at" gai:"column:next_run_at"`
	RunCount    int     `json:"run_count" gai:"column:run_count;default:0"`
	FailCount   int     `json:"fail_count" gai:"column:fail_count;default:0"`
	LastError   string  `json:"last_error" gai:"column:last_error;type:text"`
}

// TableName 指定表名
func (ScheduledTask) TableName() string {
	return "scheduled_tasks"
}

// TaskExecution 任务执行记录模型
type TaskExecution struct {
	BaseModel
	TaskID      uint64  `json:"task_id" gai:"column:task_id;index"`
	TaskName    string  `json:"task_name" gai:"column:task_name;size:200"`
	Status      string  `json:"status" gai:"column:status;size:20"` // running, success, failed, timeout
	StartedAt   string  `json:"started_at" gai:"column:started_at"`
	FinishedAt  *string `json:"finished_at" gai:"column:finished_at"`
	Duration    int     `json:"duration" gai:"column:duration"` // 毫秒
	Output      string  `json:"output" gai:"column:output;type:text"`
	ErrorMsg    string  `json:"error_msg" gai:"column:error_msg;type:text"`
	RetryCount  int     `json:"retry_count" gai:"column:retry_count;default:0"`
	TriggeredBy string  `json:"triggered_by" gai:"column:triggered_by;size:50"` // scheduler, manual, api
}

// TableName 指定表名
func (TaskExecution) TableName() string {
	return "task_executions"
}