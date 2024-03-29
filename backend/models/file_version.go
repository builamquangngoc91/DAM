package models

import "time"

type FileVersion struct {
	FileVersionID string
	FileID        string
	Size          int64
	Extension     string
	UserID        string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
