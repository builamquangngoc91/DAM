package apis

import "time"

type CreateDirectoryRequest struct {
	Name     string `json:"name"`
	ParentID string `json:"parent_id"`
}

type CreateDirectoryResponse struct {
	DirectoryID string `json:"directory_id"`
}

type UpdateDirectoryRequest struct {
	Name string `json:"name"`
}

type UpdateDirectoryResponse struct {
	DirectoryID string `json:"directory_id"`
}

type MoveDirectoriesRequest struct {
	SourceDirectoryIDs     []string `json:"source_directory_ids"`
	DestinationDirectoryID string   `json:"destination_directory_id"`
}

type Directory struct {
	DirectoryID       string    `json:"directory_id"`
	Name              string    `json:"name"`
	FullPath          string    `json:"full_path"`
	UserID            string    `json:"user_id"`
	Level             int       `json:"level"`
	ParentDirectoryID string    `json:"parent_directory_id"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
