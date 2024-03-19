package repositories

import (
	"context"
	"dam/models"

	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

type UserRepoInterface interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, userID string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
}

func NewUserRepo(db *gorm.DB) UserRepoInterface {
	return &UserRepo{db: db}
}

func (r *UserRepo) CreateUser(ctx context.Context, user *models.User) error {
	return r.db.Create(user).WithContext(ctx).Error
}

func (r *UserRepo) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	user := &models.User{}
	err := r.db.Where("user_id = ?", userID).First(user).WithContext(ctx).Error
	return user, err
}

func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	err := r.db.Where("email = ?", email).First(user).WithContext(ctx).Error
	return user, err
}

func (r *UserRepo) UpdateUser(ctx context.Context, user *models.User) error {
	return r.db.Save(user).WithContext(ctx).Error
}
