package service

import (
	"cloud-render/internal/dto"
	"cloud-render/internal/lib/password"
	"cloud-render/internal/models"
	"cloud-render/internal/repository"
	"errors"
	"strconv"

	"github.com/dgrijalva/jwt-go"
)

type UserProvider interface {
	CreateUser(user models.User) error
	GetOneUser(uid int64) (*models.User, error)
	UpdateUser(user models.User) error
	GetHashedPassword(loginOrEmail, password string) ([]models.User, error)
	UpdateRefreshToken(uid int64, refreshToken string) error
	GetRefreshToken(uid int64) (string, error)
}

type TokenManager interface {
	NewJWT(int64) (string, error)
	NewRT(int64) (string, error)
	Parse(string) (*jwt.StandardClaims, error)
}

type UserService struct {
	userProvider UserProvider
	tokenManager TokenManager
}

func NewUserService(userProvider UserProvider, tokenManager TokenManager) *UserService {
	return &UserService{
		userProvider: userProvider,
		tokenManager: tokenManager,
	}
}

func (s *UserService) CreateUser(userDTO dto.CreateUserDTO) error {
	hash, err := password.HashPassword(userDTO.Password)
	if err != nil {
		return err
	}

	err = s.userProvider.CreateUser(models.User{
		Login:    userDTO.Login,
		Email:    userDTO.Email,
		Password: hash,
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

func (s *UserService) EditUser(userDTO dto.EditUserDTO) error {
	hash, err := password.HashPassword(userDTO.Password)
	if err != nil {
		return err
	}

	err = s.userProvider.UpdateUser(models.User{
		Id:       userDTO.Id,
		Login:    userDTO.Login,
		Email:    userDTO.Email,
		Password: hash,
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
	hash, err := password.HashPassword(userDTO.Password)
	if err != nil {
		return nil, err
	}

	users, err := s.userProvider.GetHashedPassword(userDTO.LoginOrEmail, hash)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	for _, v := range users {
		if password.CheckPasswordHash(userDTO.Password, v.Password) {
			accessToken, refreshToken, err := s.updateTokens(v.Id)
			if err != nil {
				return nil, err
			}
			return &dto.AuthUserDTO{AccessToken: accessToken, RefreshToken: refreshToken}, nil
		}
	}

	return nil, ErrInvalidCredentials
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
