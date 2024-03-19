package models

import "time"

type File struct {
	FileID      string
	Name        string
	Size        int64
	Extension   string
	UserID      string
	DirectoryID string
	FullPath    string
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
