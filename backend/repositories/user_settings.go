package repositories

import (
	"context"

	"dam/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userSettingRepo struct {
	db *gorm.DB
}

type UserSettingRepoInterface interface {
	CreateUserSetting(ctx context.Context, userSetting *models.UserSetting) error
	UpdateUserSetting(ctx context.Context, userSetting *models.UserSetting) error
	GetUserSettingsByUserID(ctx context.Context, userID string, isForUpdate bool) (*models.UserSetting, error)
}

func NewUserSettingRepo(db *gorm.DB) UserSettingRepoInterface {
	return &userSettingRepo{
		db: db,
	}
}

func CreateUserSetting(ctx context.Context, db *gorm.DB, userSetting *models.UserSetting) error {
	return db.Create(userSetting).WithContext(ctx).Error
}

func (r *userSettingRepo) CreateUserSetting(ctx context.Context, userSetting *models.UserSetting) error {
	return CreateUserSetting(ctx, r.db, userSetting)
}

func UpdateUserSetting(ctx context.Context, db *gorm.DB, userSetting *models.UserSetting) error {
	return db.Where("user_setting_id = ?", userSetting.UserSettingID).Save(userSetting).WithContext(ctx).Error
}

func (r *userSettingRepo) UpdateUserSetting(ctx context.Context, userSetting *models.UserSetting) error {
	return UpdateUserSetting(ctx, r.db, userSetting)
}

func GetUserSettingsByUserID(ctx context.Context, db *gorm.DB, userID string, isForUpdate bool) (*models.UserSetting, error) {
	userSetting := &models.UserSetting{}
	db = db.Where("user_id = ?", userID).WithContext(ctx).First(userSetting)
	if isForUpdate {
		db = db.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	return userSetting, db.Error
}

func (r *userSettingRepo) GetUserSettingsByUserID(ctx context.Context, userID string, isForUpdate bool) (*models.UserSetting, error) {
	return GetUserSettingsByUserID(ctx, r.db, userID, isForUpdate)
}
