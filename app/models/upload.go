package models

// Upload represents uploaded files.
type Upload struct {
	BaseModel
	FileName     string `json:"file_name" gai:"column:file_name;size:255"`
	OriginalName string `json:"original_name" gai:"column:original_name;size:255"`
	FilePath     string `json:"file_path" gai:"column:file_path;size:500"`
	FileSize     int64  `json:"file_size" gai:"column:file_size"`
	MimeType     string `json:"mime_type" gai:"column:mime_type;size:100"`
	Extension    string `json:"extension" gai:"column:extension;size:20"`
	UserId       *int   `json:"user_id" gai:"column:user_id;nullable"`
	Module       string `json:"module" gai:"column:module;size:50"` // avatar, attachment, etc.
	Status       string `json:"status" gai:"column:status"`         // active, deleted
}

// TableName returns the database table name.
func (u *Upload) TableName() string {
	return "uploads"
}