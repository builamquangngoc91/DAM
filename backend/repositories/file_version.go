package repositories

import (
	"context"
	"dam/models"

	"gorm.io/gorm"
)

type FileVersionRepo struct {
	db *gorm.DB
}

type FileVersionRepoInterface interface {
	CreateFileVersion(ctx context.Context, fileVersion *models.FileVersion) error
	ListFileVersions(ctx context.Context, fileID string) ([]models.FileVersion, error)
}

func NewFileVersionRepo(db *gorm.DB) FileVersionRepoInterface {
	return &FileVersionRepo{
		db: db,
	}
}

func (r *FileVersionRepo) CreateFileVersion(ctx context.Context, fileVersion *models.FileVersion) error {
	return r.db.Create(fileVersion).WithContext(ctx).Error
}

func (r *FileVersionRepo) ListFileVersions(ctx context.Context, fileID string) ([]models.FileVersion, error) {
	fileVersions := []models.FileVersion{}
	err := r.db.Where("file_id = ?", fileID).Find(&fileVersions).WithContext(ctx).Error
	return fileVersions, err
}
