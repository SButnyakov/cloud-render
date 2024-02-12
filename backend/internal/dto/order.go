package dto

import "mime/multipart"

type CreateOrderDTO struct {
	UserId     int64
	Format     string
	Resolution string
	File       multipart.File
	Header     *multipart.FileHeader
}
