package repository

import (
	"cloud-render/internal/models"
	"database/sql"
	"errors"
	"fmt"
)

type PaymentTypeRepository struct {
	db *sql.DB
}

func NewPaymentTypeRepository(db *sql.DB) *PaymentTypeRepository {
	return &PaymentTypeRepository{db: db}
}

func (pt *PaymentTypeRepository) Create(status string) error {
	const fn = packagePath + "payment_type.Create"

	stmt, err := pt.db.Prepare("INSERT INTO payment_types (name) VALUES ($1)")
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

func (pt *PaymentTypeRepository) GetTypesMap() (map[string]int64, error) {
	const fn = packagePath + "payment_type.GetTypesMap"

	stmt, err := pt.db.Prepare("SELECT id, name FROM payment_types")
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement: %w", fn, err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: execute statement: %w", fn, ErrNoPaymentTypes)
		}
		return nil, fmt.Errorf("%s: execute statement: %w", fn, err)
	}
	defer rows.Close()

	types := make(map[string]int64)

	for rows.Next() {
		pType := models.PaymentType{}
		err = rows.Scan(&pType.Id, &pType.Name)
		if err != nil {
			return nil, fmt.Errorf("%s: scanning rows: %w", fn, err)
		}
		types[pType.Name] = pType.Id
	}

	return types, nil
}
