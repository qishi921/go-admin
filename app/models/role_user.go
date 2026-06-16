package models

// RoleUser represents the role_user pivot table.
type RoleUser struct {
	BaseModel
	RoleID int `json:"role_id" gai:"column:role_id"`
	UserID int `json:"user_id" gai:"column:user_id"`
}

// TableName returns the database table name.
func (ru *RoleUser) TableName() string {
	return "role_user"
}