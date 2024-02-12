package models

import "time"

type Payment struct {
	Id       int64
	DateTime time.Time
	TypeId   int64
	UserID   int64
}
