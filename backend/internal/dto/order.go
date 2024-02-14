package dto

import (
	"mime/multipart"
	"time"
)

type CreateOrderDTO struct {
	UserId     int64
	Format     string
	Resolution string
	File       multipart.File
	Header     *multipart.FileHeader
}

type GetOrderDTO struct {
	Id           int64
	Filename     string
	Date         time.Time
	OrderStatus  string
	DownloadLink string
}

type UpdateOrderStatusDTO struct {
	UserId      int64
	StoringName string
	Status      string
}

type UpdateOrderImageDTO struct {
	UserId string
	File   multipart.File
	Header *multipart.FileHeader
}
