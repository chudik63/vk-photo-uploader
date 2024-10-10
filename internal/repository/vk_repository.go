package repository

import "vk-photo-uploader/internal/entity"

type VkRepository interface {
	PhotoRepository
	CreateAlbum(title string) error
	SetToken(token string)
	SetId(id int)
}

type vkRepository struct {
	token string
	id    int
}

func NewVkRepository() VkRepository {
	return &vkRepository{}
}

func (r *vkRepository) Upload(photo *entity.Photo) error {
	return nil
}

func (r *vkRepository) CreateAlbum(title string) error {
	return nil
}

func (r *vkRepository) SetToken(token string) {
	r.token = token
}

func (r *vkRepository) SetId(id int) {
	r.id = id
}
