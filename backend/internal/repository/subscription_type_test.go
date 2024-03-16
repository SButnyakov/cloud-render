package repository_test

import (
	"cloud-render/internal/repository"

	"github.com/stretchr/testify/assert"
)

func (r *RepositoryTestSuite) TestSubscriptionTypeRepository_CreateSubscriptionType() {
	repo := repository.NewSubscriptionTypeRepository(r.api.db)

	tests := []struct {
		input string
		isErr bool
	}{
		{"new type", false},
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

func (r *RepositoryTestSuite) TestSubscriptionTypeRepository_GetTypesMap() {
	repo := repository.NewSubscriptionTypeRepository(r.api.db)

	m, _ := repo.GetTypesMap()
	assert.Equal(r.T(), 0, len(m))

	type1 := "type 1"
	type2 := "type 2"
	assert.NoError(r.T(), repo.Create(type1))
	assert.NoError(r.T(), repo.Create(type2))

	m, err := repo.GetTypesMap()
	assert.NoError(r.T(), err)

	assert.Equal(r.T(), int64(1), m[type1])
	assert.Equal(r.T(), int64(2), m[type2])
}
