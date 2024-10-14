package repository

import (
	"bufio"
	"errors"
	"fmt"
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
	dir := filepath.Join(r.path, filepath.Dir(photo.Path))

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return errors.New("ошибка создания директории")
	}

	dst, err := os.Create(filepath.Join(dir, photo.Name))
	if err != nil {
		return errors.New("ошибка создания файла")
	}
	defer dst.Close()

	if _, err := io.Copy(dst, photo.File); err != nil {
		return errors.New("ошибка сохранения файла")
	}

	photoExt := filepath.Ext(photo.Name)
	metaName := photo.Name[0 : len(photo.Name)-len(photoExt)]
	metaExt := ".txt"
	meta, err := os.Create(filepath.Join(dir, metaName+metaExt))
	if err != nil {
		return errors.New("ошибка создания метаданных")
	}
	defer meta.Close()

	writer := bufio.NewWriter(meta)
	writer.WriteString(fmt.Sprintf("Путь: %s\nРазмер: %d Byte\nДата последнего изменения: %v", photo.Path, photo.Size, photo.LastModified))

	writer.Flush()

	return nil
}

func (r *photoRepository) Delete() error {
	if err := os.RemoveAll(r.path); err != nil {
		return errors.New("ошибка удаления репозитория")
	}
	return nil
}
