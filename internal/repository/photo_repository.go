package repository

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"vk-photo-uploader/internal/entity"
)

type PhotoRepository interface {
	Upload(photo *entity.Photo) error
}

type photoRepository struct {
	path string
}

func NewPhotoRepository(path string) PhotoRepository {
	return &photoRepository{
		path: path,
	}
}

func (r *photoRepository) Upload(photo *entity.Photo) error {
	if err := os.MkdirAll(r.path, os.ModePerm); err != nil {
		return errors.New("ошибка создания директории")
	}

	dst, err := os.Create(filepath.Join(r.path, photo.Name))
	if err != nil {
		return errors.New("ошибка создания файла")
	}
	defer dst.Close()

	if _, err := io.Copy(dst, photo.File); err != nil {
		return errors.New("ошибка сохранения файла")
	}

	return nil
}
