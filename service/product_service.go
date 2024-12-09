package service

import (
	"golang-chap46/models"
	"golang-chap46/repository"
)

type ProductService interface {
	GetAllProducts() ([]models.Product, error)
}

type productService struct {
	repo repository.Repository
}

func NewProductService(repo repository.Repository) ProductService {
	return &productService{repo}
}

func (s *productService) GetAllProducts() ([]models.Product, error) {
	return s.repo.Product.GetAll()
}
