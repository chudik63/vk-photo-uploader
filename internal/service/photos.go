package service

import (
	"vk-photo-uploader/internal/entity"
	"vk-photo-uploader/internal/repository"
)

type PhotoService interface {
	UploadPhoto(photo *entity.Photo) error
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
