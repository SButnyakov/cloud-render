package repository

import (
	"cloud-render/internal/models"
	"database/sql"
	"errors"
	"fmt"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (o *OrderRepository) Create(order models.Order) error {
	const fn = packagePath + "orders.Create"

	stmt, err := o.db.Prepare("INSERT INTO orders (filename, storingname, creation_date, user_id, status_id, is_deleted) VALUES ($1, $2, $3, $4, $5, $6)")
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", fn, err)
	}

	_, err = stmt.Exec(order.FileName, order.StoringName, order.CreationDate, order.UserId, order.StatusId, false)
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", fn, err)
	}

	return nil
}

func (o *OrderRepository) GetOne(id int64) (*models.Order, error) {
	const fn = "postgres.repos.orders.Order"

	stmt, err := o.db.Prepare("SELECT id, fileName, storingName, creation_date, status_id, user_id, download_link FROM orders WHERE is_deleted = FALSE AND id = $1")
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement: %w", fn, err)
	}

	var order models.Order

	err = stmt.QueryRow(id).Scan(&order.Id, &order.FileName, &order.StoringName, &order.CreationDate, &order.StatusId, &order.UserId, &order.DownloadLink)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoOrdersFound
		}
		return nil, fmt.Errorf("%s: execute statement: %w", fn, err)
	}

	return &order, nil
}
