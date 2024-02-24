package repository_test

import (
	"cloud-render/internal/repository"

	"github.com/stretchr/testify/assert"
)

func (r *RepositoryTestSuite) TestOrderStatusRepository_CreateOrderStatus() {
	repo := repository.NewOrderStatusRepository(r.api.db)

	tests := []struct {
		input string
		isErr bool
	}{
		{"new status", false},
		{"", true},
	}

	for _, t := range tests {
		if t.isErr {
			assert.Error(r.T(), repo.Create(t.input))
		} else {
			assert.NoError(r.T(), repo.Create(t.input))
		}
	}
}

func (r *RepositoryTestSuite) TestOrderStatusRepository_GetStatusesMapStringToInt() {
	repo := repository.NewOrderStatusRepository(r.api.db)

	m, _ := repo.GetStatusesMapStringToInt()
	assert.Equal(r.T(), 0, len(m))

	status1 := "status 1"
	status2 := "status 2"
	assert.NoError(r.T(), repo.Create(status1))
	assert.NoError(r.T(), repo.Create(status2))

	m, err := repo.GetStatusesMapStringToInt()
	assert.NoError(r.T(), err)

	assert.Equal(r.T(), int64(1), m[status1])
	assert.Equal(r.T(), int64(2), m[status2])
}

func (r *RepositoryTestSuite) TestOrderStatusRepository_GetStatusesMapIntToString() {
	repo := repository.NewOrderStatusRepository(r.api.db)

	m, _ := repo.GetStatusesMapIntToString()
	assert.Equal(r.T(), 0, len(m))

	status1 := "status 1"
	status2 := "status 2"
	assert.NoError(r.T(), repo.Create(status1))
	assert.NoError(r.T(), repo.Create(status2))

	m, err := repo.GetStatusesMapIntToString()
	assert.NoError(r.T(), err)

	assert.Equal(r.T(), status1, m[1])
	assert.Equal(r.T(), status2, m[2])
}
