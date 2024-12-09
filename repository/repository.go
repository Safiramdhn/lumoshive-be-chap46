package repository

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Repository struct {
	User    UserRepository
	Product ProductRepository
}

func NewRepository(db *gorm.DB, log *zap.Logger) *Repository {
	return &Repository{
		User:    NewUserRepository(db, log),
		Product: NewProductRepository(db, log),
	}
}
