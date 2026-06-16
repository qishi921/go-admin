package models

// OperationLog represents the operation_logs table.
type OperationLog struct {
	ID         uint64  `json:"id" gai:"column:id;primaryKey"`
	UserId     *int    `json:"user_id" gai:"column:user_id;nullable"`
	Username   string  `json:"username" gai:"column:username;size:50;nullable"`
	Action     string  `json:"action" gai:"column:action;size:50"`
	Method     string  `json:"method" gai:"column:method;size:10;nullable"`
	Path       string  `json:"path" gai:"column:path;size:255;nullable"`
	Ip         string  `json:"ip" gai:"column:ip;size:45;nullable"`
	UserAgent  string  `json:"user_agent" gai:"column:user_agent;size:500;nullable"`
	Params     string  `json:"params" gai:"column:params;nullable"`
	Result     string  `json:"result" gai:"column:result;nullable"`
	Duration   int     `json:"duration" gai:"column:duration"`
	Status     string  `json:"status" gai:"column:status"`
	// 数据变更快照字段
	ResourceTable string  `json:"resource_table" gai:"column:resource_table;size:100;nullable"`
	RecordId      *uint64 `json:"record_id" gai:"column:record_id;nullable"`
	OldData       string  `json:"old_data" gai:"column:old_data;type:text;nullable"`
	NewData       string  `json:"new_data" gai:"column:new_data;type:text;nullable"`
	ChangeType    string  `json:"change_type" gai:"column:change_type;size:20;nullable"` // create, update, delete
	CreatedAt     string  `json:"created_at" gai:"column:created_at"`
	UpdatedAt     string  `json:"updated_at" gai:"column:updated_at;nullable"`
	DeletedAt     *string `json:"deleted_at,omitempty" gai:"column:deleted_at;softDelete;nullable"`
}

// TableName returns the database table name.
func (o *OperationLog) TableName() string {
	return "operation_logs"
}
