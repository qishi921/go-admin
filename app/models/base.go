package models

import "time"

// BaseModel provides common fields with string timestamps for SQLite compatibility.
// Gai ORM's default Model uses time.Time which doesn't work with SQLite's string timestamps.
type BaseModel struct {
	ID        uint64 `json:"id" gai:"column:id;primaryKey"`
	CreatedAt string `json:"created_at" gai:"column:created_at"`
	UpdatedAt string `json:"updated_at" gai:"column:updated_at"`
	DeletedAt *string `json:"deleted_at,omitempty" gai:"column:deleted_at;softDelete"`
}

// GetID returns the model ID.
func (m *BaseModel) GetID() uint64 {
	return m.ID
}

// SetID sets the model ID.
func (m *BaseModel) SetID(id uint64) {
	m.ID = id
}

// NowStr returns current time as string in SQLite format.
func NowStr() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
