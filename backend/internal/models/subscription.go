package models

import "time"

type Subscription struct {
	Id         int64
	UserId     int64
	TypeId     int64
	ExpireDate time.Time
}
