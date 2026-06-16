package models

// Permission represents the permissions table.
type Permission struct {
	BaseModel
	Name        string `json:"name" gai:"column:name;size:50"`
	Code        string `json:"code" gai:"column:code;size:50;unique"`
	Description string `json:"description" gai:"column:description;size:255;nullable"`
	Type        string `json:"type" gai:"column:type"`
	ParentId    *int   `json:"parent_id" gai:"column:parent_id;nullable"`
	SortOrder   int    `json:"sort_order" gai:"column:sort_order"`
	Status      string `json:"status" gai:"column:status"`
}

// TableName returns the database table name.
func (p *Permission) TableName() string {
	return "permissions"
}