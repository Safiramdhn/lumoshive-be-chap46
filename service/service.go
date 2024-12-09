package service

import "golang-chap46/repository"

type Service struct {
	User    UserService
	Product ProductService
}

func NewService(repo repository.Repository) *Service {
	return &Service{
		User:    NewUserService(repo),
		Product: NewProductService(repo),
	}
}
