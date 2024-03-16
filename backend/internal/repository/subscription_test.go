package repository_test

import (
	"cloud-render/internal/models"
	"cloud-render/internal/repository"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (r *RepositoryTestSuite) TestSubscriptionRepository_Create() {
	repo := repository.NewSubscriptionRepository(r.api.db)
	pRepo := repository.NewPaymentTypeRepository(r.api.db)
	sRepo := repository.NewSubscriptionTypeRepository(r.api.db)

	_, err := repo.GetOne(1)
	require.Error(r.T(), err)
	require.Equal(r.T(), repository.ErrSubscriptionNotFound, err)

	sType := "sType"
	pType := "pType"
	require.NoError(r.T(), sRepo.Create(sType))
	require.NoError(r.T(), pRepo.Create(pType))

	testTime := time.Now()

	tests := []struct {
		name  string
		s     models.Subscription
		p     models.Payment
		isErr bool
	}{
		{
			name: "correct",
			s: models.Subscription{
				UserId:     1,
				TypeId:     1,
				ExpireDate: testTime,
			},
			p: models.Payment{
				UserID:   1,
				TypeId:   1,
				DateTime: testTime,
			},
			isErr: false,
		},
		{
			name: "no sub type",
			s: models.Subscription{
				UserId:     1,
				ExpireDate: testTime,
			},
			p: models.Payment{
				UserID:   1,
				TypeId:   1,
				DateTime: testTime,
			},
			isErr: true,
		},
		{
			name: "no pay type",
			s: models.Subscription{
				UserId:     1,
				TypeId:     1,
				ExpireDate: testTime,
			},
			p: models.Payment{
				UserID:   1,
				DateTime: testTime,
			},
			isErr: true,
		},
	}

	for _, t := range tests {
		if t.isErr {
			assert.Error(r.T(), repo.Create(t.s, t.p), t.name)
		} else {
			assert.NoError(r.T(), repo.Create(t.s, t.p), t.name)
		}
	}

	sub, err := repo.GetOne(1)
	require.NoError(r.T(), err)
	require.NotNil(r.T(), sub)
	assert.Equal(r.T(), int64(1), sub.Id)
	assert.Equal(r.T(), int64(1), sub.UserId)
	assert.Equal(r.T(), int64(1), sub.TypeId)
	assert.Equal(r.T(), testTime.Year(), sub.ExpireDate.Year())
	assert.Equal(r.T(), testTime.Month(), sub.ExpireDate.Month())
	assert.Equal(r.T(), testTime.Day(), sub.ExpireDate.Day())
}

func (r *RepositoryTestSuite) TestSubscriptionRepository_GetOne() {
	repo := repository.NewSubscriptionRepository(r.api.db)
	pRepo := repository.NewPaymentTypeRepository(r.api.db)
	sRepo := repository.NewSubscriptionTypeRepository(r.api.db)

	_, err := repo.GetOne(1)
	require.Error(r.T(), err)
	require.Equal(r.T(), repository.ErrSubscriptionNotFound, err)

	sType := "sType"
	pType := "pType"
	require.NoError(r.T(), sRepo.Create(sType))
	require.NoError(r.T(), pRepo.Create(pType))

	testTime := time.Now()

	require.NoError(r.T(), repo.Create(models.Subscription{
		UserId:     1,
		TypeId:     1,
		ExpireDate: testTime,
	},
		models.Payment{
			UserID:   1,
			TypeId:   1,
			DateTime: testTime,
		},
	))

	sub, err := repo.GetOne(1)
	require.Equal(r.T(), nil, err)
	require.NotEqual(r.T(), nil, sub)
	assert.Equal(r.T(), int64(1), sub.Id)
	assert.Equal(r.T(), int64(1), sub.UserId)
	assert.Equal(r.T(), int64(1), sub.TypeId)
	assert.Equal(r.T(), testTime.Year(), sub.ExpireDate.Year())
	assert.Equal(r.T(), testTime.Month(), sub.ExpireDate.Month())
	assert.Equal(r.T(), testTime.Day(), sub.ExpireDate.Day())
}

func (r *RepositoryTestSuite) TestSubscriptionRepository_Update() {
	repo := repository.NewSubscriptionRepository(r.api.db)
	pRepo := repository.NewPaymentTypeRepository(r.api.db)
	sRepo := repository.NewSubscriptionTypeRepository(r.api.db)

	_, err := repo.GetOne(1)
	require.Error(r.T(), err)
	require.Equal(r.T(), repository.ErrSubscriptionNotFound, err)

	sType := "sType"
	pType := "pType"
	require.NoError(r.T(), sRepo.Create(sType))
	require.NoError(r.T(), pRepo.Create(pType))

	testTime := time.Now()

	require.NoError(r.T(), repo.Create(models.Subscription{
		UserId:     1,
		TypeId:     1,
		ExpireDate: testTime,
	}, models.Payment{
		UserID:   1,
		TypeId:   1,
		DateTime: testTime,
	}))

	tests := []struct {
		name  string
		s     models.Subscription
		p     models.Payment
		isErr bool
	}{
		{
			name: "correct",
			s: models.Subscription{
				UserId:     1,
				TypeId:     1,
				ExpireDate: testTime.AddDate(0, 1, 0),
			},
			p: models.Payment{
				UserID:   1,
				TypeId:   1,
				DateTime: testTime,
			},
			isErr: false,
		},
		{
			name: "no sub type",
			s: models.Subscription{
				UserId:     1,
				ExpireDate: testTime.AddDate(0, 1, 0),
			},
			p: models.Payment{
				UserID:   1,
				TypeId:   1,
				DateTime: testTime,
			},
			isErr: true,
		},
		{
			name: "no pay type",
			s: models.Subscription{
				UserId:     1,
				TypeId:     1,
				ExpireDate: testTime.AddDate(0, 1, 0),
			},
			p: models.Payment{
				UserID:   1,
				DateTime: testTime,
			},
			isErr: true,
		},
	}

	for _, t := range tests {
		if t.isErr {
			assert.Error(r.T(), repo.Update(t.s, t.p), t.name)
		} else {
			assert.NoError(r.T(), repo.Update(t.s, t.p), t.name)
		}
	}

	sub, err := repo.GetOne(1)
	require.NoError(r.T(), err)
	require.NotNil(r.T(), sub)
	require.Equal(r.T(), int64(1), sub.Id)
	require.Equal(r.T(), int64(1), sub.UserId)
	require.Equal(r.T(), int64(1), sub.TypeId)
	require.Equal(r.T(), testTime.Year(), sub.ExpireDate.Year())
	require.Equal(r.T(), testTime.AddDate(0, 1, 0).Month(), sub.ExpireDate.Month())
	require.Equal(r.T(), testTime.Day(), sub.ExpireDate.Day())
}

func (r *RepositoryTestSuite) TestSubscriptionRepository_GetExpireDate() {
	repo := repository.NewSubscriptionRepository(r.api.db)
	pRepo := repository.NewPaymentTypeRepository(r.api.db)
	sRepo := repository.NewSubscriptionTypeRepository(r.api.db)

	_, err := repo.GetOne(1)
	require.Error(r.T(), err)
	require.Equal(r.T(), repository.ErrSubscriptionNotFound, err)

	_, err = repo.GetExpireDate(1)
	require.Error(r.T(), err)
	require.Equal(r.T(), repository.ErrSubscriptionNotFound, err)

	sType := "sType"
	pType := "pType"
	require.NoError(r.T(), sRepo.Create(sType))
	require.NoError(r.T(), pRepo.Create(pType))

	testTime := time.Now()

	require.NoError(r.T(), repo.Create(models.Subscription{
		UserId:     1,
		TypeId:     1,
		ExpireDate: testTime,
	},
		models.Payment{
			UserID:   1,
			TypeId:   1,
			DateTime: testTime,
		},
	))

	date, err := repo.GetExpireDate(1)
	require.NoError(r.T(), err)
	require.NotNil(r.T(), date)
	require.Equal(r.T(), date.Year(), testTime.Year())
	require.Equal(r.T(), date.Month(), testTime.Month())
	require.Equal(r.T(), date.Day(), testTime.Day())
}
