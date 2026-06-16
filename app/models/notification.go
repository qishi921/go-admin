package models

// Notification 通知消息模型
type Notification struct {
	BaseModel
	UserID     uint    `json:"user_id" gai:"column:user_id;index"`
	Title      string  `json:"title" gai:"column:title;size:200"`
	Content    string  `json:"content" gai:"column:content;type:text"`
	Type       string  `json:"type" gai:"column:type;size:50;default:system"` // system, alert, task, message
	Priority   int     `json:"priority" gai:"column:priority;default:0"`      // 0=普通, 1=重要, 2=紧急
	IsRead     bool    `json:"is_read" gai:"column:is_read;default:false"`
	ReadAt     *string `json:"read_at" gai:"column:read_at"`
	Channel    string  `json:"channel" gai:"column:channel;size:50;default:in_app"` // in_app, email, sms
	SentAt     *string `json:"sent_at" gai:"column:sent_at"`
	SendStatus string  `json:"send_status" gai:"column:send_status;size:20;default:pending"` // pending, sent, failed
	ErrorMsg   string  `json:"error_msg" gai:"column:error_msg;type:text"`
	Metadata   string  `json:"metadata" gai:"column:metadata;type:text"` // JSON 存储额外数据
}

// TableName 指定表名
func (Notification) TableName() string {
	return "notifications"
}

// NotificationTemplate 通知模板模型
type NotificationTemplate struct {
	BaseModel
	Code      string `json:"code" gai:"column:code;size:100;unique"`
	Name      string `json:"name" gai:"column:name;size:200"`
	Title     string `json:"title" gai:"column:title;size:200"`
	Content   string `json:"content" gai:"column:content;type:text"`
	Type      string `json:"type" gai:"column:type;size:50;default:system"`
	Channels  string `json:"channels" gai:"column:channels;size:200"` // JSON 数组: ["in_app","email"]
	Variables string `json:"variables" gai:"column:variables;type:text"` // JSON 对象说明变量
	IsActive  bool   `json:"is_active" gai:"column:is_active;default:true"`
}

// TableName 指定表名
func (NotificationTemplate) TableName() string {
	return "notification_templates"
}