package service

import (
	"errors"
	"golang-chap46/helper"
	"golang-chap46/models"
	"golang-chap46/repository"
)

type UserService interface {
	Register(userInput models.User) error
	Login(username, password string) (models.User, error)
}

type userService struct {
	repo repository.Repository
}

func NewUserService(repo repository.Repository) UserService {
	return &userService{repo}
}

func (s *userService) Register(userInput models.User) error {
	userInput.Password = helper.HashPassword(userInput.Password)
	return s.repo.User.CreateUser(&userInput)
}

func (s *userService) Login(email, password string) (models.User, error) {
	user, err := s.repo.User.FindByEmail(email)
	if err != nil {
		return user, errors.New("user not found")
	}

	if helper.CheckPassword(password, user.Password) {
		return user, errors.New("invalid password")
	}

	return user, nil
}
