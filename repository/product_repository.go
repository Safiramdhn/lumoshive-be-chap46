package repository

import (
	"golang-chap46/models"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ProductRepository interface {
	GetAll() ([]models.Product, error)
}

type productRepository struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewProductRepository(db *gorm.DB, log *zap.Logger) ProductRepository {
	return &productRepository{db: db, log: log}
}

func (pr *productRepository) GetAll() ([]models.Product, error) {
	var products []models.Product
	err := pr.db.Find(&products).Error
	if err != nil {
		pr.log.Error("Error fetching products", zap.Error(err))
		return nil, err
	}

	return products, nil
}
