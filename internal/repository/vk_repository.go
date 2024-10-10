package repository

import "vk-photo-uploader/internal/entity"

type VkRepository interface {
	PhotoRepository
	CreateAlbum(title string) error
}

type vkRepository struct {
	token string
}

func NewVkRepository() VkRepository {
	return &vkRepository{
		token: "",
	}
}

func (r *vkRepository) Upload(photo *entity.Photo) error {

	return nil
}

func (r *vkRepository) CreateAlbum(title string) error {
	return nil
}
