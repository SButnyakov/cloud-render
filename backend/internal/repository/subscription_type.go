package repository

import (
	"cloud-render/internal/models"
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

func (st *SubscriptionTypeRepository) Create(status string) error {
	const fn = packagePath + "subscription_type.Create"

	stmt, err := st.db.Prepare("INSERT INTO subscription_types (name) VALUES ($1)")
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", fn, err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(status)
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", fn, err)
	}

	return nil
}

func (st *SubscriptionTypeRepository) GetTypesMap() (map[string]int64, error) {
	const fn = packagePath + "subscription_type.GetTypesMap"

	stmt, err := st.db.Prepare("SELECT id, name FROM subscription_types")
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement: %w", fn, err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: execute statement: %w", fn, ErrNoSubscriptionTypes)
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
