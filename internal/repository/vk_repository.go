package repository

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"sync"
	"vk-photo-uploader/internal/entity"
	"vk-photo-uploader/internal/repository/responses"
)

type PhotoRepository interface {
	CreateAlbum(title, token string) (int, error)
	GetUploadServer(id int, token string) (string, error)
	Upload(url, token string, photos ...*entity.Photo) error
	DeleteAlbum(name, token string) error
}

type vkRepository struct {
	mu sync.Mutex
}

func NewVkRepository() PhotoRepository {
	return &vkRepository{}
}

func (r *vkRepository) CreateAlbum(title, token string) (int, error) {
	r.mu.Lock()

	id, err := r.getAlbumId(title, token)
	if err == nil {
		r.mu.Unlock()
		return id, nil
	}

	resp, err := http.Get(fmt.Sprintf("https://api.vk.com/method/photos.createAlbum?title=%s&access_token=%s&v=5.199", title, token))
	if err != nil {
		return 0, errors.New("нельзя создать альбом VK")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, errors.New("нельзя считать ответ создания альбома")
	}

	msg := &responses.AlbumCreateResponse{}

	err = json.Unmarshal(body, msg)
	if err != nil {
		return 0, errors.New("ошибка чтения JSON после создании альбома")
	}

	r.mu.Unlock()

	return msg.Response.Id, nil
}

func (r *vkRepository) GetUploadServer(id int, token string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.vk.com/method/photos.getUploadServer?&access_token=%s&album_id=%d&v=5.199", token, id))
	if err != nil {
		return "", errors.New("ошибка получения сервера загрузки")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("нельзя считать ответ создания альбома")
	}

	msg := &responses.GetUploadServerResponse{}

	err = json.Unmarshal(body, msg)
	if err != nil {
		return "", errors.New("ошибка чтения JSON после получения сервера загрузки")
	}
	return msg.Response.UploadUrl, nil
}

func (r *vkRepository) Upload(url, token string, photos ...*entity.Photo) error {
	b := &bytes.Buffer{}
	writer := multipart.NewWriter(b)

	for i, photo := range photos {
		part, err := writer.CreateFormFile(fmt.Sprintf("file%d", i+1), photo.Name)
		if err != nil {
			return err
		}

		if _, err = io.Copy(part, photo.File); err != nil {
			return err
		}

		writer.Close()
	}

	req, err := http.NewRequest("POST", url, b)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	msg := &responses.PostUrlResponse{}
	if err = json.Unmarshal(body, msg); err != nil {
		return err
	}

	err = r.savePhoto(msg.Aid, msg.Server, msg.Photos_list, msg.Hash, token)

	if err != nil {
		return err
	}

	return nil
}

func (r *vkRepository) DeleteAlbum(title, token string) error {
	id, err := r.getAlbumId(title, token)
	if err != nil {
		return err
	}

	http.Get(fmt.Sprintf("https://api.vk.com/method/photos.deleteAlbum?album_id=%d&access_token=%s&v=5.199", id, token))

	return nil
}

func (r *vkRepository) savePhoto(aid, server int, photos_list, hash, token string) error {
	resp, err := http.Get(fmt.Sprintf("https://api.vk.com/method/photos.save?album_id=%d&server=%d&photos_list=%s&hash=%s&access_token=%s&v=5.199", aid, server, photos_list, hash, token))
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	msg := &responses.SavePhotoResponse{}

	err = json.Unmarshal(body, msg)
	if err != nil {
		return err
	}

	if len(msg.Response) == 0 {
		return errors.New("ошибка загрузки фотографий")
	}

	return nil
}

func (r *vkRepository) getAlbumId(title, token string) (int, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.vk.com/method/photos.getAlbums?access_token=%s&v=5.199", token))
	if err != nil {
		return 0, errors.New("нельзя получить информацию о альбомах VK")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, errors.New("нельзя считать ответ c информацией о альбомах")
	}

	msg := &responses.GetAlbumsResponse{}
	if err := json.Unmarshal(body, msg); err != nil {
		return 0, errors.New("ошибка чтения JSON после проверки альбома")
	}

	for i := 0; i < msg.Response.Count; i++ {
		if msg.Response.Array[i].Title == title {
			return msg.Response.Array[i].Id, nil
		}
	}

	return 0, errors.New("альбом не существует")
}
