package service

import (
	"cloud-render/internal/dto"
	"cloud-render/internal/lib/tokenManager"
	"cloud-render/internal/models"
	"cloud-render/internal/repository"
	"errors"
	"strconv"
)

type UserProvider interface {
	CreateUser(user models.User) error
	GetOneUser(uid int64) (*models.User, error)
	UpdateUser(user models.User) error
	CheckCredentials(loginOrEmail, password string) (int64, error)
	UpdateRefreshToken(uid int64, refreshToken string) error
	GetRefreshToken(uid int64) (string, error)
}

type UserService struct {
	userProvider UserProvider
	tokenManager *tokenManager.Manager
}

func NewUserService(userProvider UserProvider, tokenManager *tokenManager.Manager) *UserService {
	return &UserService{
		userProvider: userProvider,
		tokenManager: tokenManager,
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

func (s *UserService) GetUser(id int64) (*dto.GetUserDTO, error) {
	user, err := s.userProvider.GetOneUser(id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &dto.GetUserDTO{
		Login: user.Login,
		Email: user.Email,
	}, nil
}

func (s *UserService) EditUer(userDTO dto.EditUserDTO) error {
	err := s.userProvider.UpdateUser(models.User{
		Id:       userDTO.Id,
		Login:    userDTO.Login,
		Email:    userDTO.Email,
		Password: userDTO.Password,
	})
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return ErrInvalidCredentials
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

	accessToken, refreshToken, err := s.updateTokens(id)
	if err != nil {
		return nil, err
	}

	return &dto.AuthUserDTO{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (s *UserService) ReauthUser(userDTO dto.ReAuthUserDTO) (*dto.ReAuthUserDTO, error) {
	claims, err := s.tokenManager.Parse(userDTO.RefreshToken)
	if err != nil {
		return nil, err
	}

	id, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil {
		return nil, err
	}

	token, err := s.userProvider.GetRefreshToken(id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if token != userDTO.RefreshToken {
		return nil, ErrInvalidCredentials
	}

	accessToken, refreshToken, err := s.updateTokens(id)
	if err != nil {
		return nil, err
	}

	return &dto.ReAuthUserDTO{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (s *UserService) updateTokens(id int64) (string, string, error) {
	accessToken, err := s.tokenManager.NewJWT(id)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.tokenManager.NewRT(id)
	if err != nil {
		return "", "", err
	}

	err = s.userProvider.UpdateRefreshToken(id, refreshToken)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
