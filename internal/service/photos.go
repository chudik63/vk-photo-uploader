package service

import (
	"vk-photo-uploader/internal/entity"
	"vk-photo-uploader/internal/repository"
)

type PhotoService interface {
	UploadPhoto(photo *entity.Photo) error
	DeleteFolder(name string) error
}

type photoService struct {
	photoRepo repository.PhotoRepository
}

func NewPhotoService(photoRepo repository.PhotoRepository) PhotoService {
	return &photoService{
		photoRepo: photoRepo,
	}
}

func (p *photoService) UploadPhoto(photo *entity.Photo) error {
	return p.photoRepo.Upload(photo)
}

func (p *photoService) DeleteFolder(name string) error {
	return p.photoRepo.DeleteFolder(name)
}
