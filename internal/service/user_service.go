package service

import (
	"shop/internal/model"
	"shop/internal/repository"
)

type UserService interface {
	Register(username, email, password string) error
	GetUser(id uint) (*model.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Register(username, email, password string) error {
	// In a real app, hash the password here
	user := &model.User{
		Username: username,
		Email:    email,
		Password: password,
	}
	return s.repo.Create(user)
}

func (s *userService) GetUser(id uint) (*model.User, error) {
	return s.repo.FindByID(id)
}
