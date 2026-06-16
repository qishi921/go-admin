package models

// RolePermission represents the role_permission pivot table.
type RolePermission struct {
	BaseModel
	RoleID       int `json:"role_id"       gai:"column:role_id"`
	PermissionID int `json:"permission_id" gai:"column:permission_id"`
}

// TableName returns the database table name.
func (rp *RolePermission) TableName() string {
	return "role_permission"
}