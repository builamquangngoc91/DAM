package repositories

import (
	"context"
	"dam/models"

	"gorm.io/gorm"
)

type DirectoryRepoInterface interface {
	CreateDirectory(ctx context.Context, directory *models.Directory) error
	UpdateDirectory(ctx context.Context, directory *models.Directory) error
	GetDirectoryByID(ctx context.Context, directoryID string) (*models.Directory, error)
	GetDirectoryByFullPath(ctx context.Context, fullPath string) (*models.Directory, error)
	ListFilesOrFoldersByDirectoryID(ctx context.Context, directoryID string, orderBy string, limit, offset int) ([]models.FileOrFolder, error)
}

type DirectoryRepo struct {
	db *gorm.DB
}

func NewDirectoryRepo(db *gorm.DB) DirectoryRepoInterface {
	return &DirectoryRepo{db: db}
}

func (r *DirectoryRepo) CreateDirectory(ctx context.Context, directory *models.Directory) error {
	return r.db.Create(directory).WithContext(ctx).Error
}

func (r *DirectoryRepo) UpdateDirectory(ctx context.Context, directory *models.Directory) error {
	return r.db.Where("directory_id = ?", directory.DirectoryID).Save(directory).WithContext(ctx).Error
}

func (r *DirectoryRepo) GetDirectoryByID(ctx context.Context, directoryID string) (*models.Directory, error) {
	directory := &models.Directory{}
	err := r.db.Where("directory_id = ?", directoryID).First(directory).WithContext(ctx).Error
	return directory, err
}

func (r *DirectoryRepo) GetDirectoryByFullPath(ctx context.Context, fullPath string) (*models.Directory, error) {
	directory := &models.Directory{}
	err := r.db.Where("full_path = ?", fullPath).First(directory).WithContext(ctx).Error
	return directory, err
}

func (r *DirectoryRepo) ListFilesOrFoldersByDirectoryID(ctx context.Context, directoryID string, orderBy string, limit, offset int) ([]models.FileOrFolder, error) {
	filesOrFolders := []models.FileOrFolder{}
	err := r.db.
		WithContext(ctx).
		Raw(`
			SELECT id, parent_directory_id, name, full_path, created_at, updated_at, is_directory
			FROM (
				SELECT directory_id AS id, parent_directory_id, name, full_path, created_at, updated_at, true AS is_directory 
				FROM directories
				UNION
				SELECT file_id AS id, directory_id, name, full_path, created_at, updated_at, false AS is_directory
				FROM files
			) AS files_or_folders
			WHERE parent_directory_id = ?
			ORDER BY ?
			LIMIT ?
			OFFSET ?
		`, directoryID, orderBy, limit, offset).
		Scan(&filesOrFolders).
		Error

	return filesOrFolders, err
}
