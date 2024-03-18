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

func (o *OrderRepository) Create(order models.Order) (int64, error) {
	const fn = packagePath + "order.Create"

	stmt, err := o.db.Prepare("INSERT INTO orders (filename, storingname, creation_date, user_id, status_id, is_deleted) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id")
	if err != nil {
		return 0, fmt.Errorf("%s: prepare statement: %w", fn, err)
	}
	defer stmt.Close()

	var lastInsertedId int64

	err = stmt.QueryRow(order.FileName, order.StoringName, order.CreationDate, order.UserId, order.StatusId, false).Scan(&lastInsertedId)
	if err != nil {
		return 0, fmt.Errorf("%s: execute statement: %w", fn, err)
	}

	return lastInsertedId, nil
}

func (o *OrderRepository) GetOne(id int64) (*models.Order, error) {
	const fn = packagePath + "order.GetOne"

	stmt, err := o.db.Prepare("SELECT id, fileName, storingName, creation_date, status_id, user_id, download_link FROM orders WHERE is_deleted = FALSE AND id = $1")
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement: %w", fn, err)
	}
	defer stmt.Close()

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

func (o *OrderRepository) GetMany(id int64) ([]models.Order, error) {
	const fn = packagePath + "order.Orders"

	stmt, err := o.db.Prepare("SELECT id, fileName, storingName, creation_date, status_id, user_id, download_link FROM orders WHERE is_deleted = FALSE AND user_id = $1")
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement: %w", fn, err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(id)
	if err != nil {
		return nil, fmt.Errorf("%s: execute statement: %w", fn, err)
	}
	defer rows.Close()

	orders := make([]models.Order, 0)

	for rows.Next() {
		order := models.Order{}
		err = rows.Scan(&order.Id, &order.FileName, &order.StoringName, &order.CreationDate, &order.StatusId, &order.UserId, &order.DownloadLink)
		if err != nil {
			return nil, fmt.Errorf("%s: scanning rows: %w", fn, err)
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (o *OrderRepository) UpdateStatus(orderId, statusId int64) error {
	const fn = packagePath + "order.UpdateStatus"

	stmt, err := o.db.Prepare("UPDATE orders SET status_id = $2 WHERE id = $1")
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", fn, err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(orderId, statusId)
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", fn, err)
	}

	affectedRows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: calc affected rows: %w", fn, err)
	}

	if affectedRows == 0 {
		return ErrNoOrdersFound
	}

	return nil
}

func (o *OrderRepository) UpdateDownloadLink(orderId int64, downloadLink string) error {
	const fn = packagePath + "order.UpdateDownloadLink"

	stmt, err := o.db.Prepare("UPDATE orders SET download_link = $2 WHERE id = $1")
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", fn, err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(orderId, downloadLink)
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", fn, err)
	}

	affectedRows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: calc affected rows: %w", fn, err)
	}

	if affectedRows == 0 {
		return ErrNoOrdersFound
	}

	return nil
}

func (o *OrderRepository) SoftDelete(id int64) error {
	const fn = packagePath + "order.SoftDelete"

	stmt, err := o.db.Prepare("UPDATE orders SET is_deleted=TRUE WHERE id=$1")
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", fn, err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", fn, err)
	}

	if affected, _ := res.RowsAffected(); affected == 0 {
		return ErrNoOrdersFound
	}

	return nil
}
