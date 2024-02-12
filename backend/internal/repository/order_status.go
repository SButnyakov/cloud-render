package repository

import (
	"cloud-render/internal/models"
	"database/sql"
	"errors"
	"fmt"
)

type OrderStatusesRepository struct {
	db *sql.DB
}

func NewOrderStatusesRepository(db *sql.DB) *OrderStatusesRepository {
	return &OrderStatusesRepository{db: db}
}

func (os *OrderStatusesRepository) GetStatusesMapStringToInt() (map[string]int64, error) {
	const fn = packagePath + "order_statuses.GetStatusesMapStringToInt"

	stmt, err := os.db.Prepare("SELECT id, name FROM order_statuses")
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement: %w", fn, err)
	}

	rows, err := stmt.Query()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: execute statement: %w", fn, ErrNoOrderStatuses)
		}
		return nil, fmt.Errorf("%s: execute statement: %w", fn, err)
	}
	defer rows.Close()

	statuses := make(map[string]int64)

	for rows.Next() {
		status := models.OrderStatus{}
		err = rows.Scan(&status.Id, &status.Name)
		if err != nil {
			return nil, fmt.Errorf("%s: scanning rows: %w", fn, err)
		}
		statuses[status.Name] = status.Id
	}

	return statuses, nil
}

func (os *OrderStatusesRepository) GetStatusesMapIntToString() (map[int64]string, error) {
	const fn = packagePath + "order_statuses.GetStatusesMapIntToString"

	stmt, err := os.db.Prepare("SELECT id, name FROM order_statuses")
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement: %w", fn, err)
	}

	rows, err := stmt.Query()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: execute statement: %w", fn, ErrNoOrderStatuses)
		}
		return nil, fmt.Errorf("%s: execute statement: %w", fn, err)
	}
	defer rows.Close()

	statuses := make(map[int64]string)

	for rows.Next() {
		status := models.OrderStatus{}
		err = rows.Scan(&status.Id, &status.Name)
		if err != nil {
			return nil, fmt.Errorf("%s: scanning rows: %w", fn, err)
		}
		statuses[status.Id] = status.Name
	}

	return statuses, nil
}
