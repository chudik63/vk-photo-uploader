package service

import (
	"vk-photo-uploader/internal/entity"
	"vk-photo-uploader/internal/repository"
)

type UserService struct {
	vkRepo repository.VkRepository
}

func NewUserService(vkRepo repository.VkRepository) *UserService {
	return &UserService{
		vkRepo: vkRepo,
	}
}

func (u *UserService) Register(user *entity.User) {
	u.vkRepo.SetToken(user.AccessToken)
	u.vkRepo.SetId(user.UserID)
}

func (u *UserService) Send(path string) error {
	// parallelism add
	return u.vkRepo.Upload(path)
}
