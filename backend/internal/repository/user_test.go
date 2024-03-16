package repository_test

import (
	"cloud-render/internal/models"
	"cloud-render/internal/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (r *RepositoryTestSuite) TestUserRepository_Create() {
	repo := repository.NewUserRepository(r.auth.db)

	user, err := repo.GetOneUser(1)
	require.Error(r.T(), err)
	require.ErrorIs(r.T(), err, repository.ErrUserNotFound)
	require.Nil(r.T(), user)

	correctUser := models.User{
		Login:    "login",
		Email:    "email",
		Password: "password",
	}

	tests := []struct {
		name  string
		input models.User
		isErr bool
	}{
		{
			name:  "correct",
			input: correctUser,
			isErr: false,
		},
		{
			name: "duplicate login",
			input: models.User{
				Login:    correctUser.Login,
				Email:    "email2",
				Password: "password",
			},
			isErr: true,
		},
		{
			name: "duplicate email",
			input: models.User{
				Login:    "login2",
				Email:    correctUser.Email,
				Password: "password",
			},
			isErr: true,
		},
	}

	for _, t := range tests {
		if t.isErr {
			err := repo.CreateUser(t.input)
			require.Error(r.T(), err, t.name)
			assert.ErrorIs(r.T(), err, repository.ErrUserExists, t.name)
		} else {
			assert.NoError(r.T(), repo.CreateUser(t.input), t.name)
		}
	}

	user, err = repo.GetOneUser(1)
	require.NoError(r.T(), err)
	assert.Equal(r.T(), correctUser.Login, user.Login)
	assert.Equal(r.T(), correctUser.Email, user.Email)
	assert.Equal(r.T(), correctUser.Password, user.Password)
}

func (r *RepositoryTestSuite) TestUserRepository_GetOneUser() {
	repo := repository.NewUserRepository(r.auth.db)

	user, err := repo.GetOneUser(1)
	require.Error(r.T(), err)
	require.ErrorIs(r.T(), err, repository.ErrUserNotFound)
	require.Nil(r.T(), user)

	correctUser := models.User{
		Login:    "login",
		Email:    "email",
		Password: "password",
	}

	require.NoError(r.T(), repo.CreateUser(correctUser))

	user, err = repo.GetOneUser(1)
	require.NoError(r.T(), err)
	assert.Equal(r.T(), correctUser.Login, user.Login)
	assert.Equal(r.T(), correctUser.Email, user.Email)
	assert.Equal(r.T(), correctUser.Password, user.Password)
}

func (r *RepositoryTestSuite) TestUserRepository_UpdateUser() {
	repo := repository.NewUserRepository(r.auth.db)

	user, err := repo.GetOneUser(1)
	require.Error(r.T(), err)
	require.ErrorIs(r.T(), err, repository.ErrUserNotFound)
	require.Nil(r.T(), user)

	correctUser := models.User{
		Login:    "login",
		Email:    "email",
		Password: "password",
	}

	require.NoError(r.T(), repo.CreateUser(correctUser))

	user, err = repo.GetOneUser(1)
	require.NoError(r.T(), err)
	assert.Equal(r.T(), correctUser.Login, user.Login)
	assert.Equal(r.T(), correctUser.Email, user.Email)
	assert.Equal(r.T(), correctUser.Password, user.Password)

	updatedUser := models.User{
		Id:       1,
		Login:    "new login",
		Email:    "new email",
		Password: "new password",
	}

	tests := []struct {
		name  string
		input models.User
		isErr bool
	}{
		{
			name:  "correct",
			input: updatedUser,
			isErr: false,
		},
		{
			name: "wrong id",
			input: models.User{
				Id:       99,
				Login:    updatedUser.Login,
				Email:    updatedUser.Email,
				Password: updatedUser.Password,
			},
			isErr: true,
		},
		{
			name: "no login",
			input: models.User{
				Id:       1,
				Email:    updatedUser.Email,
				Password: updatedUser.Password,
			},
			isErr: true,
		},
		{
			name: "no email",
			input: models.User{
				Id:       1,
				Login:    updatedUser.Login,
				Password: updatedUser.Password,
			},
			isErr: true,
		},
		{
			name: "no password",
			input: models.User{
				Id:    99,
				Login: updatedUser.Login,
				Email: updatedUser.Email,
			},
			isErr: true,
		},
	}

	for _, t := range tests {
		if t.isErr {
			assert.Error(r.T(), repo.UpdateUser(t.input), t.name)
		} else {
			assert.NoError(r.T(), repo.UpdateUser(t.input), t.name)
		}
	}

	user, err = repo.GetOneUser(1)
	require.NoError(r.T(), err)
	assert.Equal(r.T(), updatedUser.Login, user.Login)
	assert.Equal(r.T(), updatedUser.Email, user.Email)
	assert.Equal(r.T(), updatedUser.Password, user.Password)
}

func (r *RepositoryTestSuite) TestUserRepository_GetHashedPassword() {
	repo := repository.NewUserRepository(r.auth.db)

	correctUser := models.User{
		Login:    "login",
		Email:    "email",
		Password: "password",
	}

	users, err := repo.GetHashedPassword(correctUser.Login, correctUser.Password)
	require.NoError(r.T(), err)
	require.Equal(r.T(), 0, len(users))

	users, err = repo.GetHashedPassword(correctUser.Email, correctUser.Password)
	require.NoError(r.T(), err)
	require.Equal(r.T(), 0, len(users))

	require.NoError(r.T(), repo.CreateUser(correctUser))

	users, err = repo.GetHashedPassword(correctUser.Login, correctUser.Password)
	require.NoError(r.T(), err)
	require.Equal(r.T(), 1, len(users))
	assert.Equal(r.T(), correctUser.Password, users[0].Password)

	users, err = repo.GetHashedPassword(correctUser.Email, correctUser.Password)
	require.NoError(r.T(), err)
	require.Equal(r.T(), 1, len(users))
	assert.Equal(r.T(), correctUser.Password, users[0].Password)
}

func (r *RepositoryTestSuite) TestUserRepository_UpdateRefreshToken() {
	repo := repository.NewUserRepository(r.auth.db)

	correctUser := models.User{
		Login:    "login",
		Email:    "email",
		Password: "password",
	}

	require.NoError(r.T(), repo.CreateUser(correctUser))

	_, err := repo.GetRefreshToken(1)
	require.Error(r.T(), err)

	tests := []struct {
		name   string
		token  string
		userid int64
		isErr  bool
	}{
		{
			name:   "wrong user id",
			token:  "token",
			userid: 99,
			isErr:  true,
		},
		{
			name:   "correct",
			token:  "token",
			userid: 1,
			isErr:  false,
		},
	}

	for _, t := range tests {
		if t.isErr {
			err := repo.UpdateRefreshToken(t.userid, t.token)
			require.Error(r.T(), err, t.name)
			assert.ErrorIs(r.T(), err, repository.ErrUserNotFound)
		} else {
			assert.NoError(r.T(), repo.UpdateRefreshToken(t.userid, t.token), t.name)
		}
	}

	token, err := repo.GetRefreshToken(1)
	require.NoError(r.T(), err)
	require.Equal(r.T(), "token", token)
}

func (r *RepositoryTestSuite) TestUserRepository_GetRefreshToken() {
	repo := repository.NewUserRepository(r.auth.db)

	correctUser := models.User{
		Login:    "login",
		Email:    "email",
		Password: "password",
	}

	_, err := repo.GetRefreshToken(1)
	require.Error(r.T(), err)
	require.ErrorIs(r.T(), err, repository.ErrUserNotFound)

	require.NoError(r.T(), repo.CreateUser(correctUser))

	require.NoError(r.T(), repo.UpdateRefreshToken(1, "token"))

	token, err := repo.GetRefreshToken(1)
	require.NoError(r.T(), err)
	require.Equal(r.T(), "token", token)
}
