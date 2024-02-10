package service

import (
	"cloud-render/internal/dto"
	"cloud-render/internal/models"
)

type UserProvider interface {
	CreateUser(user models.User) error
}

type UserService struct {
	userProvider UserProvider
}

func NewUserService(userProvider UserProvider) *UserService {
	return &UserService{
		userProvider: userProvider,
	}
}

func (s *UserService) CreateUser(userDTO dto.CreateUserDTO) error {
	return s.userProvider.CreateUser(models.User{
		Login:    userDTO.Login,
		Email:    userDTO.Email,
		Password: userDTO.Password,
	})

}
