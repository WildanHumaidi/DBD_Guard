package repository

import (
	"context"

	"github.com/dbd-guard/backend/internal/domain"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) error
	GetUsersByRole(ctx context.Context, roleName string) ([]domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, userID string) error
	GetRoleByName(ctx context.Context, roleName string) (*domain.Role, error)
	GetUserByID(ctx context.Context, userID string) (*domain.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Preload("Role").Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) GetUsersByRole(ctx context.Context, roleName string) ([]domain.User, error) {
	var users []domain.User
	err := r.db.WithContext(ctx).Joins("JOIN roles ON roles.role_id = users.role_id").
		Where("roles.role_name = ?", roleName).
		Find(&users).Error
	return users, err
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) Delete(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&domain.User{}).Error
}

func (r *userRepository) GetRoleByName(ctx context.Context, roleName string) (*domain.Role, error) {
	var role domain.Role
	if err := r.db.WithContext(ctx).Where("role_name = ?", roleName).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, userID string) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
