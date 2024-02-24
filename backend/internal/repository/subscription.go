package repository

import (
	"cloud-render/internal/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type SubscriptionRepository struct {
	db *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (s *SubscriptionRepository) Create(subscription models.Subscription, payment models.Payment) error {
	const fn = packagePath + "subscription.Create"

	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s: prepare transaction: %w", fn, err)
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO payments (date, type_id, user_id) VALUES ($1, $2, $3)",
		payment.DateTime, payment.TypeId, payment.UserID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: execute statement: %w", fn, err)
	}

	_, err = tx.ExecContext(ctx,
		"INSERT INTO subscriptions (user_id, type_id, sub_expire_date) VALUES ($1, $2, $3)",
		subscription.UserId, subscription.TypeId, subscription.ExpireDate)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: execute statement: %w", fn, err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("%s: commit transaction: %w", fn, err)
	}

	return nil
}

func (s *SubscriptionRepository) GetOne(id int64) (*models.Subscription, error) {
	const fn = packagePath + "subscription.GetOne"

	stmt, err := s.db.Prepare("SELECT id, user_id, type_id, sub_expire_date FROM subscriptions WHERE id=$1")
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement: %w", fn, err)
	}
	defer stmt.Close()

	var sub models.Subscription

	err = stmt.QueryRow(id).Scan(&sub.Id, &sub.UserId, &sub.TypeId, &sub.ExpireDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSubscriptionNotFound
		}
		return nil, fmt.Errorf("%s: execute statement: %w", fn, err)
	}

	return &sub, nil
}

func (s *SubscriptionRepository) Update(subscription models.Subscription, payment models.Payment) error {
	const fn = packagePath + "subscription.Update"

	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s: prepare transaction: %w", fn, err)
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO payments (date, type_id, user_id) VALUES ($1, $2, $3)",
		payment.DateTime, payment.TypeId, payment.UserID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: execute statement: %w", fn, err)
	}

	_, err = tx.ExecContext(ctx,
		"UPDATE subscriptions SET type_id = $2, sub_expire_date = $3 WHERE user_id = $1",
		subscription.UserId, subscription.TypeId, subscription.ExpireDate)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: execute statement: %w", fn, err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("%s: commit transaction: %w", fn, err)
	}

	return nil
}

func (s *SubscriptionRepository) GetExpireDate(uid int64) (*time.Time, error) {
	const fn = packagePath + "subscription.GetExpireDate"

	stmt, err := s.db.Prepare("SELECT sub_expire_date FROM subscriptions WHERE user_id = $1")
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement: %w", fn, err)
	}
	defer stmt.Close()

	var expireDate time.Time

	err = stmt.QueryRow(uid).Scan(&expireDate)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSubscriptionNotFound
		}
		return nil, fmt.Errorf("%s: execute statement: %w", fn, err)
	}
	return &expireDate, nil
}
