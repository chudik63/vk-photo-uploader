package entity

import (
	"mime/multipart"
	"time"
)

type Photo struct {
	File         multipart.File
	Name         string
	Path         string
	Size         uint64
	Type         string
	LastModified time.Time
}
