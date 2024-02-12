package repository

import (
	"cloud-render/internal/models"
	"database/sql"
	"fmt"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (o *OrderRepository) Create(order models.Order) error {
	const fn = packagePath + ".orders.Create"

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
