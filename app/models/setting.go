package models

// Setting represents system settings stored in database.
type Setting struct {
	BaseModel
	Key       string  `json:"key"       gai:"column:key;size:100;unique"`
	Value     string  `json:"value"     gai:"column:value;size:500"`
	Type      string  `json:"type"      gai:"column:type;size:20"`          // string, number, boolean, json
	GroupName string  `json:"group"     gai:"column:group_name;size:50"`    // system, email, etc. (use GroupName to avoid SQL reserved word)
	Label     string  `json:"label"     gai:"column:label;size:100"`        // Display name
	Options   *string `json:"options"   gai:"column:options;size:500;nullable"` // JSON options for select
	IsPublic  bool    `json:"is_public" gai:"column:is_public"`              // Can be accessed without auth
}

// TableName returns the database table name.
func (s *Setting) TableName() string {
	return "settings"
}