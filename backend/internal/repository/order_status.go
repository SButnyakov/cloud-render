package repository

import (
	"cloud-render/internal/models"
	"database/sql"
	"errors"
	"fmt"
)

type OrderStatusRepository struct {
	db *sql.DB
}

func NewOrderStatusRepository(db *sql.DB) *OrderStatusRepository {
	return &OrderStatusRepository{db: db}
}

func (os *OrderStatusRepository) Create(status string) error {
	const fn = packagePath + "order_status.Create"

	stmt, err := os.db.Prepare("INSERT INTO order_statuses (name) VALUES ($1)")
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

func (os *OrderStatusRepository) GetStatusesMapStringToInt() (map[string]int64, error) {
	const fn = packagePath + "order_statuse.GetStatusesMapStringToInt"

	stmt, err := os.db.Prepare("SELECT id, name FROM order_statuses")
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement: %w", fn, err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoOrderStatuses
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

func (os *OrderStatusRepository) GetStatusesMapIntToString() (map[int64]string, error) {
	const fn = packagePath + "order_statuse.GetStatusesMapIntToString"

	stmt, err := os.db.Prepare("SELECT id, name FROM order_statuses")
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement: %w", fn, err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
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

	if len(statuses) == 0 {
		return nil, ErrNoOrderStatuses
	}

	return statuses, nil
}
