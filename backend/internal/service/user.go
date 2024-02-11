package service

import (
	"cloud-render/internal/dto"
	"cloud-render/internal/lib/tokenManager"
	"cloud-render/internal/models"
	"cloud-render/internal/repository"
	"errors"
)

type UserProvider interface {
	CreateUser(user models.User) error
	CheckCredentials(loginOrEmail, password string) (int64, error)
	UpdateRefreshToken(uid int64, refreshToken string) error
}

type UserService struct {
	userProvider UserProvider
	tokenManager *tokenManager.Manager
}

func NewUserService(userProvider UserProvider) *UserService {
	return &UserService{
		userProvider: userProvider,
	}
}

func (s *UserService) CreateUser(userDTO dto.CreateUserDTO) error {
	err := s.userProvider.CreateUser(models.User{
		Login:    userDTO.Login,
		Email:    userDTO.Email,
		Password: userDTO.Password,
	})
	if err != nil {
		if errors.Is(err, repository.ErrUserExists) {
			return ErrUserAlreadyExists
		}
		return err
	}
	return nil
}

func (s *UserService) AuthUser(userDTO dto.AuthUserDTO) (*dto.AuthUserDTO, error) {
	id, err := s.userProvider.CheckCredentials(userDTO.LoginOrEmail, userDTO.Password)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	accessToken, err := s.tokenManager.NewJWT(id)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.tokenManager.NewRT(id)
	if err != nil {
		return nil, err
	}

	err = s.userProvider.UpdateRefreshToken(id, refreshToken)
	if err != nil {
		return nil, err
	}

	return &dto.AuthUserDTO{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
