package apis

import "time"

type UploadFileResponse struct {
	FileVersionID string `json:"file_version_id"`
}

type MoveFilesRequest struct {
	SourceFileIDs          []string `json:"source_file_ids"`
	DestinationDirectoryID string   `json:"destination_directory_id"`
}

type UpdateFileRequest struct {
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

type ListFileVersions struct {
	File         File          `json:"file"`
	FileVersions []FileVersion `json:"file_versions"`
}

type File struct {
	FileID      string    `json:"file_id"`
	Name        string    `json:"name"`
	Size        int64     `json:"size"`
	Extension   string    `json:"extension"`
	UserID      string    `json:"user_id"`
	DirectoryID string    `json:"directory_id"`
	FullPath    string    `json:"full_path"`
	Description string    `json:"description"`
	Tags        []string  `json:"tags"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type FileVersion struct {
	FileVersionID string    `json:"file_version_id"`
	FileID        string    `json:"file_id"`
	Size          int64     `json:"size"`
	Extension     string    `json:"extension"`
	UserID        string    `json:"user_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
