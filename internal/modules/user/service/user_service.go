package service

import (
	"template/internal/modules/user/model"
	"template/internal/modules/user/repository"
)

type UserService interface {
	Register(email, password, name string) error
	GetAllUsers() ([]model.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo}
}

func (s *userService) Register(email, password, name string) error {
	user := &model.User{
		Email:    email,
		Password: password,
		Name:     name,
	}
	return s.repo.Create(user)
}

func (s *userService) GetAllUsers() ([]model.User, error) {
	return s.repo.FindAll()
}
