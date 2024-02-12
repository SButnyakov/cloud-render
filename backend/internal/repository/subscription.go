package repository

import (
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

func (s *SubscriptionRepository) GetExpireDate(uid int64) (*time.Time, error) {
	const fn = "postgres.repos.subscription.GetExpireDate"

	stmt, err := s.db.Prepare("SELECT sub_expire_date FROM subscriptions WHERE user_id = $1")
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement: %w", fn, err)
	}

	var expireDate *time.Time

	_ = stmt.QueryRow(uid).Scan(&expireDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: execute statement: %w", fn, storage.ErrSubscriptionNotFound)
		}
		return nil, fmt.Errorf("%s: execute statement: %w", fn, err)
	}

	return expireDate, nil
}
