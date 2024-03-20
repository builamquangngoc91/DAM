package repositories

import (
	"context"
	"dam/models"

	"gorm.io/gorm"
)

type FileRepo struct {
	db *gorm.DB
}

type FileRepoInterface interface {
	CreateFile(ctx context.Context, file *models.File) error
	UpdateFile(ctx context.Context, file *models.File) error
	GetFileByID(ctx context.Context, fileID string) (*models.File, error)
	MoveDirectory(ctx context.Context, sourceDirectory, destinationDirectory *models.Directory) error
}

func NewFileRepo(db *gorm.DB) FileRepoInterface {
	return &FileRepo{db: db}
}

func (r *FileRepo) CreateFile(ctx context.Context, file *models.File) error {
	return r.db.Create(file).WithContext(ctx).Error
}

func (r *FileRepo) UpdateFile(ctx context.Context, file *models.File) error {
	return r.db.Where("file_id = ?", file.FileID).Save(file).WithContext(ctx).Error
}

func (r *FileRepo) GetFileByID(ctx context.Context, fileID string) (*models.File, error) {
	file := &models.File{}
	err := r.db.Where("file_id = ?", fileID).WithContext(ctx).First(file).Error
	return file, err
}

func (r *FileRepo) MoveDirectory(ctx context.Context, sourceDirectory, destinationDirectory *models.Directory) error {
	return r.db.
		WithContext(ctx).
		Exec(`
			UPDATE directories
			SET full_path = REPLACE(full_path, ?, ?)
			WHERE full_path LIKE ?
		`, sourceDirectory.FullPath, destinationDirectory.FullPath+"/"+sourceDirectory.DirectoryID, sourceDirectory.FullPath+"%").
		Error
}
