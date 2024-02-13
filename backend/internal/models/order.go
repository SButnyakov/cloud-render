package models

import (
	"database/sql"
	"time"
)

type Order struct {
	Id           int64
	FileName     string
	StoringName  string
	CreationDate time.Time
	UserId       int64
	StatusId     int64
	DownloadLink sql.NullString
	IsDeleted    bool
}
