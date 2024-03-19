package models

import "time"

type Directory struct {
	DirectoryID       string
	Name              string
	FullPath          string
	UserID            string
	Level             int
	ParentDirectoryID string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
