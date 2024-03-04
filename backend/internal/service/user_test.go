package service_test

import (
	"cloud-render/internal/dto"
	"cloud-render/internal/lib/password"
	"cloud-render/internal/lib/tokenManager"
	mocks "cloud-render/internal/mocks/repository"
	"cloud-render/internal/models"
	"cloud-render/internal/repository"
	"cloud-render/internal/service"
	"errors"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserService_CreateUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUserProvider := mocks.NewMockUserProvider(mockCtrl)

	tm, err := tokenManager.New("secretKey")
	require.NoError(t, err)

	userService := service.NewUserService(mockUserProvider, tm)

	userDTO := dto.CreateUserDTO{
		Login:    "newUser",
		Email:    "newUser@example.com",
		Password: "password",
	}

	mockUserProvider.EXPECT().
		CreateUser(gomock.AssignableToTypeOf(models.User{})).
		Return(nil).
		Times(1)

	err = userService.CreateUser(userDTO)
	assert.NoError(t, err)

	mockUserProvider.EXPECT().
		CreateUser(gomock.AssignableToTypeOf(models.User{})).
		Return(repository.ErrUserExists).
		Times(1)

	err = userService.CreateUser(userDTO)
	assert.NotNil(t, err)
	assert.Equal(t, err, service.ErrUserAlreadyExists)

	unknownError := errors.New("unknown error")

	mockUserProvider.EXPECT().
		CreateUser(gomock.AssignableToTypeOf(models.User{})).
		Return(unknownError).
		Times(1)

	err = userService.CreateUser(userDTO)
	assert.NotNil(t, err)
	assert.Equal(t, err, unknownError)
}

func TestUserService_GetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserProvider := mocks.NewMockUserProvider(ctrl)
	tm, err := tokenManager.New("secretKey")
	require.NoError(t, err)

	userService := service.NewUserService(mockUserProvider, tm)

	userID := int64(1)
	expectedUserDTO := &dto.GetUserDTO{
		Login: "testuser",
		Email: "testuser@example.com",
	}

	mockUserProvider.EXPECT().
		GetOneUser(userID).
		Return(&models.User{Login: expectedUserDTO.Login, Email: expectedUserDTO.Email}, nil).
		Times(1)

	userDTO, err := userService.GetUser(userID)
	assert.NoError(t, err)
	require.NotNil(t, userDTO)
	assert.Equal(t, expectedUserDTO.Login, userDTO.Login)
	assert.Equal(t, expectedUserDTO.Email, userDTO.Email)

	mockUserProvider.EXPECT().
		GetOneUser(userID).
		Return(nil, repository.ErrUserNotFound).
		Times(1)

	userDTO, err = userService.GetUser(userID)
	assert.Error(t, err)
	assert.Nil(t, userDTO)
	assert.Equal(t, service.ErrUserNotFound, err)

	unknownError := errors.New("unknown error")

	mockUserProvider.EXPECT().
		GetOneUser(userID).
		Return(nil, unknownError).
		Times(1)

	userDTO, err = userService.GetUser(userID)
	assert.Error(t, err)
	assert.Nil(t, userDTO)
	assert.Equal(t, err, unknownError)
}

func TestUserService_EditUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserProvider := mocks.NewMockUserProvider(ctrl)
	tm, err := tokenManager.New("secretKey")
	require.NoError(t, err)

	userService := service.NewUserService(mockUserProvider, tm)

	userDTO := dto.EditUserDTO{
		Id:       1,
		Login:    "newlogin",
		Email:    "newemail@example.com",
		Password: "newpassword",
	}

	mockUserProvider.EXPECT().
		UpdateUser(gomock.AssignableToTypeOf(models.User{})).
		Return(nil).
		Times(1)

	err = userService.EditUser(userDTO)
	assert.NoError(t, err)

	mockUserProvider.EXPECT().
		UpdateUser(gomock.AssignableToTypeOf(models.User{})).
		Return(repository.ErrUserNotFound).
		Times(1)

	err = userService.EditUser(userDTO)
	assert.Error(t, err)
	assert.Equal(t, service.ErrInvalidCredentials, err)

	unknownError := errors.New("unknown error")

	mockUserProvider.EXPECT().
		UpdateUser(gomock.AssignableToTypeOf(models.User{})).
		Return(unknownError).
		Times(1)

	err = userService.EditUser(userDTO)
	assert.Error(t, err)
	assert.Equal(t, err, unknownError)
}

func TestUserService_AuthUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserProvider := mocks.NewMockUserProvider(ctrl)
	mockTokenManager := mocks.NewMockTokenManager(ctrl)

	userService := service.NewUserService(mockUserProvider, mockTokenManager)

	userDTO := dto.AuthUserDTO{
		LoginOrEmail: "testuser",
		Password:     "password123",
	}

	hashedPassword, _ := password.HashPassword(userDTO.Password)
	mockUserProvider.EXPECT().
		GetHashedPassword(userDTO.LoginOrEmail, gomock.AssignableToTypeOf("")).
		Return([]models.User{{Id: 1, Password: hashedPassword}}, nil).
		Times(1)
	mockUserProvider.EXPECT().
		UpdateRefreshToken(int64(1), "refreshToken").
		Return(nil).
		Times(1)
	mockTokenManager.EXPECT().
		NewJWT(int64(1)).
		Return("accessToken", nil).
		Times(1)
	mockTokenManager.EXPECT().
		NewRT(int64(1)).
		Return("refreshToken", nil).
		Times(1)

	authUserDTO, err := userService.AuthUser(userDTO)
	assert.NoError(t, err)
	require.NotNil(t, authUserDTO)
	assert.Equal(t, "accessToken", authUserDTO.AccessToken)
	assert.Equal(t, "refreshToken", authUserDTO.RefreshToken)

	mockUserProvider.EXPECT().
		GetHashedPassword(userDTO.LoginOrEmail, gomock.AssignableToTypeOf("")).
		Return(nil, repository.ErrUserNotFound).
		Times(1)

	authUserDTO, err = userService.AuthUser(userDTO)
	assert.Error(t, err)
	assert.Nil(t, authUserDTO)
	assert.Equal(t, service.ErrInvalidCredentials, err)

	mockUserProvider.EXPECT().
		GetHashedPassword(userDTO.LoginOrEmail, gomock.AssignableToTypeOf("")).
		Return([]models.User{{Id: 1, Password: "wrongPassword"}}, nil).
		Times(1)

	authUserDTO, err = userService.AuthUser(userDTO)
	assert.Error(t, err)
	assert.Nil(t, authUserDTO)
	assert.Equal(t, service.ErrInvalidCredentials, err)

	mockUserProvider.EXPECT().
		GetHashedPassword(userDTO.LoginOrEmail, gomock.AssignableToTypeOf("")).
		Return(nil, errors.New("some repository error")).
		Times(1)

	authUserDTO, err = userService.AuthUser(userDTO)
	assert.Error(t, err)
	assert.Nil(t, authUserDTO)
	assert.EqualError(t, err, "some repository error")

	mockUserProvider.EXPECT().
		GetHashedPassword(userDTO.LoginOrEmail, gomock.AssignableToTypeOf("")).
		Return([]models.User{{Id: 1, Password: hashedPassword}}, nil).
		Times(1)
	mockUserProvider.EXPECT().
		UpdateRefreshToken(int64(1), "refreshToken").
		Return(errors.New("failed to update refresh token")).
		Times(1)
	mockTokenManager.EXPECT().
		NewJWT(int64(1)).
		Return("accessToken", nil).
		Times(1)
	mockTokenManager.EXPECT().
		NewRT(int64(1)).
		Return("refreshToken", nil).
		Times(1)

	authUserDTO, err = userService.AuthUser(userDTO)
	assert.Error(t, err)
	assert.Nil(t, authUserDTO)
	assert.EqualError(t, err, "failed to update refresh token")

	mockUserProvider.EXPECT().
		GetHashedPassword(userDTO.LoginOrEmail, gomock.AssignableToTypeOf("")).
		Return([]models.User{{Id: 1, Password: hashedPassword}}, nil).
		Times(1)
	mockTokenManager.EXPECT().
		NewJWT(int64(1)).
		Return("", errors.New("failed to generate jwt")).
		Times(1)

	authUserDTO, err = userService.AuthUser(userDTO)
	assert.Error(t, err)
	assert.Nil(t, authUserDTO)
	assert.EqualError(t, err, "failed to generate jwt")

	mockUserProvider.EXPECT().
		GetHashedPassword(userDTO.LoginOrEmail, gomock.AssignableToTypeOf("")).
		Return([]models.User{{Id: 1, Password: hashedPassword}}, nil).
		Times(1)
	mockTokenManager.EXPECT().
		NewJWT(int64(1)).
		Return("accessToken", nil).
		Times(1)
	mockTokenManager.EXPECT().
		NewRT(int64(1)).
		Return("", errors.New("failed to generate rt")).
		Times(1)

	authUserDTO, err = userService.AuthUser(userDTO)
	assert.Error(t, err)
	assert.Nil(t, authUserDTO)
	assert.EqualError(t, err, "failed to generate rt")
}

func TestUserService_ReauthUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserProvider := mocks.NewMockUserProvider(ctrl)
	mockTokenManager := mocks.NewMockTokenManager(ctrl)

	userService := service.NewUserService(mockUserProvider, mockTokenManager)

	refreshToken := "validRefreshToken"
	userDTO := dto.ReAuthUserDTO{
		RefreshToken: refreshToken,
	}

	// Проверка успешной переаутентификации пользователя
	mockTokenManager.EXPECT().
		Parse(refreshToken).
		Return(&jwt.StandardClaims{Subject: "1"}, nil).
		Times(1)
	mockUserProvider.EXPECT().
		GetRefreshToken(int64(1)).
		Return(refreshToken, nil).
		Times(1)
	mockUserProvider.EXPECT().
		UpdateRefreshToken(int64(1), "refreshToken").
		Return(nil).
		Times(1)
	mockTokenManager.EXPECT().
		NewJWT(int64(1)).
		Return("accessToken", nil).
		Times(1)
	mockTokenManager.EXPECT().
		NewRT(int64(1)).
		Return("refreshToken", nil).
		Times(1)

	reAuthUserDTO, err := userService.ReauthUser(userDTO)
	assert.NoError(t, err)
	assert.NotNil(t, reAuthUserDTO)
	assert.Equal(t, "accessToken", reAuthUserDTO.AccessToken)
	assert.Equal(t, "refreshToken", reAuthUserDTO.RefreshToken)

	// Проверка ошибки при переаутентификации (токен не удалось распарсить)
	mockTokenManager.EXPECT().
		Parse(refreshToken).
		Return(nil, errors.New("token parse error")).
		Times(1)

	reAuthUserDTO, err = userService.ReauthUser(userDTO)
	assert.Error(t, err)
	assert.Nil(t, reAuthUserDTO)
	assert.EqualError(t, err, "token parse error")

	// Проверка ошибки при переаутентификации (пользователь не найден)
	mockTokenManager.EXPECT().
		Parse(refreshToken).
		Return(&jwt.StandardClaims{Subject: "1"}, nil).
		Times(1)
	mockUserProvider.EXPECT().
		GetRefreshToken(int64(1)).
		Return("", repository.ErrUserNotFound).
		Times(1)

	reAuthUserDTO, err = userService.ReauthUser(userDTO)
	assert.Error(t, err)
	assert.Nil(t, reAuthUserDTO)
	assert.Equal(t, service.ErrInvalidCredentials, err)

	// Проверка ошибки при переаутентификации (несовпадение токенов)
	mockTokenManager.EXPECT().
		Parse(refreshToken).
		Return(&jwt.StandardClaims{Subject: "1"}, nil).
		Times(1)
	mockUserProvider.EXPECT().
		GetRefreshToken(int64(1)).
		Return("differentToken", nil).
		Times(1)

	reAuthUserDTO, err = userService.ReauthUser(userDTO)
	assert.Error(t, err)
	assert.Nil(t, reAuthUserDTO)
	assert.Equal(t, service.ErrInvalidCredentials, err)

	// Проверка другой ошибки при переаутентификации
	mockTokenManager.EXPECT().
		Parse(refreshToken).
		Return(&jwt.StandardClaims{Subject: "1"}, nil).
		Times(1)
	mockUserProvider.EXPECT().
		GetRefreshToken(int64(1)).
		Return("", errors.New("some repository error")).
		Times(1)

	reAuthUserDTO, err = userService.ReauthUser(userDTO)
	assert.Error(t, err)
	assert.Nil(t, reAuthUserDTO)
	assert.EqualError(t, err, "some repository error")

	mockTokenManager.EXPECT().
		Parse(refreshToken).
		Return(&jwt.StandardClaims{Subject: "not int"}, nil).
		Times(1)

	reAuthUserDTO, err = userService.ReauthUser(userDTO)
	assert.Error(t, err)
	assert.Nil(t, reAuthUserDTO)
}
