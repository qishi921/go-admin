package models

// Menu represents the menus table.
type Menu struct {
	BaseModel
	Name      string `json:"name" gai:"column:name;size:50"`
	Path      string `json:"path" gai:"column:path;size:100"`
	Icon      string `json:"icon" gai:"column:icon;size:50;nullable"`
	Component string `json:"component" gai:"column:component;size:100;nullable"`
	ParentId  *int   `json:"parent_id" gai:"column:parent_id;nullable"`
	SortOrder int    `json:"sort_order" gai:"column:sort_order"`
	Status    string `json:"status" gai:"column:status"`
}

// TableName returns the database table name.
func (m *Menu) TableName() string {
	return "menus"
}