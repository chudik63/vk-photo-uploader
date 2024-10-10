package entity

import (
	"mime/multipart"
	"time"
)

type Photo struct {
	File         multipart.File
	Name         string
	Path         string
	Size         int64
	LastModified time.Time
}
