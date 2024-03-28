package models

import (
	"time"

	"github.com/lib/pq"
)

type File struct {
	FileID      string
	Name        string
	Size        int64
	Extension   string
	UserID      string
	DirectoryID string
	FullPath    string
	Description string
	Tags        pq.StringArray `gorm:"type:_text"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type FileOrFolder struct {
	ID                string
	Name              string
	ParentDirectoryID string
	IsDirectory       bool
	FullPath          string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
