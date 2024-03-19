package apis

import "time"

type UploadFileResponse struct {
	FileID string `json:"file_id"`
}

type File struct {
	FileID      string    `json:"file_id"`
	Name        string    `json:"name"`
	Size        int64     `json:"size"`
	Extension   string    `json:"extension"`
	UserID      string    `json:"user_id"`
	DirectoryID string    `json:"directory_id"`
	FullPath    string    `json:"full_path"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
