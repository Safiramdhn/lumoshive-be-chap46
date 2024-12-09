package repository

import (
	"golang-chap46/models"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	FindByEmail(email string) (models.User, error)
}

type userRepository struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewUserRepository(db *gorm.DB, log *zap.Logger) UserRepository {
	return &userRepository{db, log}
}

func (r *userRepository) CreateUser(user *models.User) error {
	r.log.Info("Creating user", zap.Any("user", user))
	return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return user, err
	}
	r.log.Info("Found user by email", zap.String("email", email), zap.Any("user", user))
	return user, nil
}
