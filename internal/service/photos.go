package service

import (
	"vk-photo-uploader/internal/entity"
	"vk-photo-uploader/internal/repository"
)

type PhotoService interface {
	UploadPhotos(photos []*entity.Photo, token string) error
	DeleteFolder(name string, token string) error
}

type photoService struct {
	photoRepo repository.PhotoRepository
}

func NewPhotoService(photoRepo repository.PhotoRepository) PhotoService {
	return &photoService{
		photoRepo: photoRepo,
	}
}

func (p *photoService) UploadPhotos(photos []*entity.Photo, token string) error {
	id, err := p.photoRepo.CreateAlbum(photos[0].Folder, token)
	if err != nil || id == 0 {
		return err
	}

	url, err := p.photoRepo.GetUploadServer(id, token)
	if err != nil || url == "" {
		return err
	}

	err = p.photoRepo.Upload(url, token, photos...)
	if err != nil {
		return err
	}

	return nil
}

func (p *photoService) DeleteFolder(name string, token string) error {
	return p.photoRepo.DeleteAlbum(name, token)
}
