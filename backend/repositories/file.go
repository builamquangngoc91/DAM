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
	GetFileByID(ctx context.Context, fileID string) (*models.File, error)
}

func NewFileRepo(db *gorm.DB) FileRepoInterface {
	return &FileRepo{db: db}
}

func (r *FileRepo) CreateFile(ctx context.Context, file *models.File) error {
	return r.db.Create(file).WithContext(ctx).Error
}

func (r *FileRepo) GetFileByID(ctx context.Context, fileID string) (*models.File, error) {
	file := &models.File{}
	err := r.db.Where("file_id = ?", fileID).WithContext(ctx).First(file).Error
	return file, err
}
