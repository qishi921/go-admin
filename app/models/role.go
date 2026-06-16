package models

// Role represents the roles table.
type Role struct {
	BaseModel
	Name        string `json:"name" gai:"column:name;size:50"`
	Code        string `json:"code" gai:"column:code;size:50;unique"`
	Description string `json:"description" gai:"column:description;size:255;nullable"`
	Status      string `json:"status" gai:"column:status"`
}

// TableName returns the database table name.
func (r *Role) TableName() string {
	return "roles"
}