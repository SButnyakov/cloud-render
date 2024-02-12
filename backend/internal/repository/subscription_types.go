package repository

import (
	"database/sql"
	"errors"
	"fmt"
)

type SubscriptionTypeRepository struct {
	db *sql.DB
}

func NewSubscriptionTypeRepository(db *sql.DB) *SubscriptionTypeRepository {
	return &SubscriptionTypeRepository{db: db}
}

func (st *SubscriptionTypeRepository) GetTypesMap() (map[string]int64, error) {
	const fn = "postgres.repos.subscription_types.GetTypesMap"

	stmt, err := st.db.Prepare("SELECT id, name FROM subscription_types")
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement: %w", fn, err)
	}

	rows, err := stmt.Query()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: execute statement: %w", fn, storage.ErrNoSubscriptionTypes)
		}
		return nil, fmt.Errorf("%s: execute statement: %w", fn, err)
	}
	defer rows.Close()

	types := make(map[string]int64)

	for rows.Next() {
		sType := models.SubscriptionType{}
		err = rows.Scan(&sType.Id, &sType.Name)
		if err != nil {
			return nil, fmt.Errorf("%s: scanning rows: %w", fn, err)
		}
		types[sType.Name] = sType.Id
	}

	return types, nil
}
