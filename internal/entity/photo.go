package entity

import (
	"mime/multipart"
)

type Photo struct {
	File   multipart.File
	Name   string
	Folder string
	Size   int64
}
