package repository_test

import (
	"cloud-render/internal/models"
	"cloud-render/internal/repository"
	"database/sql"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (r *RepositoryTestSuite) TestOrderRepository_Create() {
	repo := repository.NewOrderRepository(r.api.db)
	sRepo := repository.NewOrderStatusRepository(r.api.db)

	order, err := repo.GetOne(1)
	require.Error(r.T(), err)
	require.Equal(r.T(), repository.ErrNoOrdersFound, err)
	require.Nil(r.T(), order)

	require.NoError(r.T(), sRepo.Create("status"))

	testTime := time.Now()

	tests := []struct {
		name  string
		input models.Order
		isErr bool
	}{
		{
			name: "correct",
			input: models.Order{
				FileName:     "filename",
				StoringName:  "storingname",
				CreationDate: testTime,
				UserId:       1,
				StatusId:     1,
				IsDeleted:    false,
			},
			isErr: false,
		},
		{
			name: "no filename name",
			input: models.Order{
				StoringName:  "storingname",
				CreationDate: testTime,
				UserId:       1,
				StatusId:     1,
				IsDeleted:    false,
			},
			isErr: true,
		},
		{
			name: "no storing name",
			input: models.Order{
				FileName:     "filename",
				CreationDate: testTime,
				UserId:       1,
				StatusId:     1,
				IsDeleted:    false,
			},
			isErr: true,
		},
		{
			name: "wrong status id",
			input: models.Order{
				FileName:     "filename",
				StoringName:  "storingname",
				CreationDate: testTime,
				UserId:       1,
				StatusId:     999,
				IsDeleted:    false,
			},
			isErr: true,
		},
	}

	for _, t := range tests {
		if t.isErr {
			assert.Error(r.T(), repo.Create(t.input), t.name)
		} else {
			assert.NoError(r.T(), repo.Create(t.input), t.name)
		}
	}

	order, err = repo.GetOne(1)
	require.NoError(r.T(), err)
	require.NotNil(r.T(), order)
	assert.Equal(r.T(), int64(1), order.Id)
	assert.Equal(r.T(), int64(1), order.UserId)
	assert.Equal(r.T(), int64(1), order.StatusId)
	assert.Equal(r.T(), "filename", order.FileName)
	assert.Equal(r.T(), "storingname", order.StoringName)
	assert.Equal(r.T(), testTime.Year(), order.CreationDate.Year())
	assert.Equal(r.T(), testTime.Month(), order.CreationDate.Month())
	assert.Equal(r.T(), testTime.Day(), order.CreationDate.Day())
	assert.False(r.T(), order.IsDeleted)
}

func (r *RepositoryTestSuite) TestOrderRepository_GetOne() {
	repo := repository.NewOrderRepository(r.api.db)
	sRepo := repository.NewOrderStatusRepository(r.api.db)

	order, err := repo.GetOne(1)
	require.Error(r.T(), err)
	require.Equal(r.T(), repository.ErrNoOrdersFound, err)
	require.Nil(r.T(), order)

	require.NoError(r.T(), sRepo.Create("status"))

	testTime := time.Now()

	require.NoError(r.T(), repo.Create(models.Order{
		FileName:     "filename",
		StoringName:  "storingname",
		CreationDate: testTime,
		UserId:       1,
		StatusId:     1,
		IsDeleted:    false,
	}))

	order, err = repo.GetOne(1)
	require.NoError(r.T(), err)
	require.NotNil(r.T(), order)
	assert.Equal(r.T(), int64(1), order.Id)
	assert.Equal(r.T(), int64(1), order.UserId)
	assert.Equal(r.T(), int64(1), order.StatusId)
	assert.Equal(r.T(), "filename", order.FileName)
	assert.Equal(r.T(), "storingname", order.StoringName)
	assert.Equal(r.T(), testTime.Year(), order.CreationDate.Year())
	assert.Equal(r.T(), testTime.Month(), order.CreationDate.Month())
	assert.Equal(r.T(), testTime.Day(), order.CreationDate.Day())
	assert.False(r.T(), order.IsDeleted)
}

func (r *RepositoryTestSuite) TestOrderRepository_GetMany() {
	repo := repository.NewOrderRepository(r.api.db)
	sRepo := repository.NewOrderStatusRepository(r.api.db)

	orders, err := repo.GetMany(1)
	require.NoError(r.T(), err)
	require.Equal(r.T(), 0, len(orders))

	require.NoError(r.T(), sRepo.Create("status"))

	testTime := time.Now()

	require.NoError(r.T(), repo.Create(models.Order{
		FileName:     "filename1",
		StoringName:  "storingname1",
		CreationDate: testTime,
		UserId:       1,
		StatusId:     1,
		IsDeleted:    false,
	}))
	require.NoError(r.T(), repo.Create(models.Order{
		FileName:     "filename2",
		StoringName:  "storingname2",
		CreationDate: testTime,
		UserId:       1,
		StatusId:     1,
		IsDeleted:    false,
	}))

	orders, err = repo.GetMany(1)
	require.NoError(r.T(), err)
	require.Equal(r.T(), 2, len(orders))

	require.NotNil(r.T(), orders[0])
	assert.Equal(r.T(), int64(1), orders[0].Id)
	assert.Equal(r.T(), int64(1), orders[0].UserId)
	assert.Equal(r.T(), int64(1), orders[0].StatusId)
	assert.Equal(r.T(), "filename1", orders[0].FileName)
	assert.Equal(r.T(), "storingname1", orders[0].StoringName)
	assert.Equal(r.T(), testTime.Year(), orders[0].CreationDate.Year())
	assert.Equal(r.T(), testTime.Month(), orders[0].CreationDate.Month())
	assert.Equal(r.T(), testTime.Day(), orders[0].CreationDate.Day())
	assert.False(r.T(), orders[0].IsDeleted)

	require.NotNil(r.T(), orders[1])
	assert.Equal(r.T(), int64(2), orders[1].Id)
	assert.Equal(r.T(), int64(1), orders[1].UserId)
	assert.Equal(r.T(), int64(1), orders[1].StatusId)
	assert.Equal(r.T(), "filename2", orders[1].FileName)
	assert.Equal(r.T(), "storingname2", orders[1].StoringName)
	assert.Equal(r.T(), testTime.Year(), orders[1].CreationDate.Year())
	assert.Equal(r.T(), testTime.Month(), orders[1].CreationDate.Month())
	assert.Equal(r.T(), testTime.Day(), orders[1].CreationDate.Day())
	assert.False(r.T(), orders[1].IsDeleted)
}

func (r *RepositoryTestSuite) TestOrderRepository_UpdateStatus() {
	repo := repository.NewOrderRepository(r.api.db)
	sRepo := repository.NewOrderStatusRepository(r.api.db)

	order, err := repo.GetOne(1)
	require.Error(r.T(), err)
	require.Equal(r.T(), repository.ErrNoOrdersFound, err)
	require.Nil(r.T(), order)

	require.NoError(r.T(), sRepo.Create("status"))

	testTime := time.Now()

	require.NoError(r.T(), repo.Create(models.Order{
		FileName:     "filename",
		StoringName:  "storingname",
		CreationDate: testTime,
		UserId:       1,
		StatusId:     1,
		IsDeleted:    false,
	}))

	require.NoError(r.T(), sRepo.Create("new status"))

	tests := []struct {
		name        string
		storingname string
		status      int64
		userId      int64
		isErr       bool
	}{
		{
			name:        "wrong stroingname",
			storingname: "not exists",
			status:      2,
			userId:      1,
			isErr:       true,
		},
		{
			name:        "wrong status id",
			storingname: "storingname",
			status:      99,
			userId:      1,
			isErr:       true,
		},
		{
			name:        "correct",
			storingname: "storingname",
			status:      2,
			userId:      1,
			isErr:       false,
		},
	}

	for _, t := range tests {
		if t.isErr {
			assert.Error(r.T(), repo.UpdateStatus(t.storingname, t.userId, t.status), t.name)
		} else {
			assert.NoError(r.T(), repo.UpdateStatus(t.storingname, t.userId, t.status), t.name)
		}
	}

	order, err = repo.GetOne(1)
	require.NoError(r.T(), err)
	require.NotNil(r.T(), order)
	assert.Equal(r.T(), int64(1), order.Id)
	assert.Equal(r.T(), int64(1), order.UserId)
	assert.Equal(r.T(), int64(2), order.StatusId)
	assert.Equal(r.T(), "filename", order.FileName)
	assert.Equal(r.T(), "storingname", order.StoringName)
	assert.Equal(r.T(), testTime.Year(), order.CreationDate.Year())
	assert.Equal(r.T(), testTime.Month(), order.CreationDate.Month())
	assert.Equal(r.T(), testTime.Day(), order.CreationDate.Day())
	assert.False(r.T(), order.IsDeleted)
}

func (r *RepositoryTestSuite) TestOrderRepository_UpdateDownloadLink() {
	repo := repository.NewOrderRepository(r.api.db)
	sRepo := repository.NewOrderStatusRepository(r.api.db)

	order, err := repo.GetOne(1)
	require.Error(r.T(), err)
	require.Equal(r.T(), repository.ErrNoOrdersFound, err)
	require.Nil(r.T(), order)

	require.NoError(r.T(), sRepo.Create("status"))

	testTime := time.Now()

	require.NoError(r.T(), repo.Create(models.Order{
		FileName:     "filename",
		StoringName:  "storingname",
		CreationDate: testTime,
		UserId:       1,
		StatusId:     1,
		IsDeleted:    false,
	}))

	tests := []struct {
		name        string
		storingname string
		link        string
		userId      int64
		isErr       bool
	}{
		{
			name:        "wrong stroingname",
			storingname: "not exists",
			link:        "link",
			userId:      1,
			isErr:       true,
		},
		{
			name:        "wrong user id",
			storingname: "storingname",
			link:        "link",
			userId:      99,
			isErr:       true,
		},
		{
			name:        "correct",
			storingname: "storingname",
			link:        "link",
			userId:      1,
			isErr:       false,
		},
	}

	for _, t := range tests {
		if t.isErr {
			assert.Error(r.T(), repo.UpdateDownloadLink(t.userId, t.storingname, t.link), t.name)
		} else {
			assert.NoError(r.T(), repo.UpdateDownloadLink(t.userId, t.storingname, t.link), t.name)
		}
	}

	order, err = repo.GetOne(1)
	require.NoError(r.T(), err)
	require.NotNil(r.T(), order)
	assert.Equal(r.T(), int64(1), order.Id)
	assert.Equal(r.T(), int64(1), order.UserId)
	assert.Equal(r.T(), int64(1), order.StatusId)
	assert.Equal(r.T(), "filename", order.FileName)
	assert.Equal(r.T(), "storingname", order.StoringName)
	assert.Equal(r.T(), testTime.Year(), order.CreationDate.Year())
	assert.Equal(r.T(), testTime.Month(), order.CreationDate.Month())
	assert.Equal(r.T(), testTime.Day(), order.CreationDate.Day())
	assert.Equal(r.T(), sql.NullString{String: "link", Valid: true}, order.DownloadLink)
	assert.False(r.T(), order.IsDeleted)
}

func (r *RepositoryTestSuite) TestOrderRepository_SoftDelete() {
	repo := repository.NewOrderRepository(r.api.db)
	sRepo := repository.NewOrderStatusRepository(r.api.db)

	order, err := repo.GetOne(1)
	require.Error(r.T(), err)
	require.Equal(r.T(), repository.ErrNoOrdersFound, err)
	require.Nil(r.T(), order)

	require.Error(r.T(), repo.SoftDelete(1))

	require.NoError(r.T(), sRepo.Create("status"))

	testTime := time.Now()

	require.NoError(r.T(), repo.Create(models.Order{
		FileName:     "filename",
		StoringName:  "storingname",
		CreationDate: testTime,
		UserId:       1,
		StatusId:     1,
		IsDeleted:    false,
	}))

	order, err = repo.GetOne(1)
	require.NoError(r.T(), err)
	require.NotNil(r.T(), order)
	assert.Equal(r.T(), int64(1), order.Id)
	assert.Equal(r.T(), int64(1), order.UserId)
	assert.Equal(r.T(), int64(1), order.StatusId)
	assert.Equal(r.T(), "filename", order.FileName)
	assert.Equal(r.T(), "storingname", order.StoringName)
	assert.Equal(r.T(), testTime.Year(), order.CreationDate.Year())
	assert.Equal(r.T(), testTime.Month(), order.CreationDate.Month())
	assert.Equal(r.T(), testTime.Day(), order.CreationDate.Day())
	assert.False(r.T(), order.IsDeleted)

	require.NoError(r.T(), repo.SoftDelete(1))

	order, err = repo.GetOne(1)
	require.Error(r.T(), err)
	require.Equal(r.T(), repository.ErrNoOrdersFound, err)
	require.Nil(r.T(), order)
}
