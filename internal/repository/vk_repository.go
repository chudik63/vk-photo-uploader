package repository

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"vk-photo-uploader/internal/entity"
	"vk-photo-uploader/internal/repository/responses"
)

type PhotoRepository interface {
	UploadPhoto(photo *entity.Photo) error
	DeleteFolder(name string) error
}

type VkRepository interface {
	PhotoRepository
	SetToken(token string)
	SetId(id int)
}

type vkRepository struct {
	token string
	id    int
}

func NewVkRepository() VkRepository {
	return &vkRepository{
		token: "",
		id:    0,
	}
}

func (r *vkRepository) UploadPhoto(photo *entity.Photo) error {
	dir := filepath.Dir(photo.Path)

	id, err := r.createAlbum(dir)
	if err != nil {
		return err
	}

	upload_url, err := r.getUploadServer(id)
	if err != nil {
		return err
	}
	b := &bytes.Buffer{}
	writer := multipart.NewWriter(b)

	part, err := writer.CreateFormFile("file1", photo.Name)
	if err != nil {
		return err
	}

	if _, err = io.Copy(part, photo.File); err != nil {
		return err
	}

	writer.Close()

	req, _ := http.NewRequest("POST", upload_url, b)

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

	_, err = http.Get(fmt.Sprintf("https://api.vk.com/method/photos.save?album_id=%d&server=%d&photos_list=%s&hash=%s&access_token=%s&v=5.199", msg.Aid, msg.Server, msg.Photos_list, msg.Hash, r.token))
	if err != nil {
		return err
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	/// ПРОВЕРИТЬ загрузку
	msg1 := &responses.SavePhotoResponse{}

	err = json.Unmarshal(body, msg1)
	if err != nil {
		return err
	}

	return nil
}

func (r *vkRepository) DeleteFolder(path string) error {
	return nil
}

func (r *vkRepository) readCaption(path string) (string, error) {
	ext := filepath.Ext(path)
	wholePathWithoutExt := strings.TrimSuffix(path, ext)

	wholePathTxt := wholePathWithoutExt + ".txt"

	caption, err := os.ReadFile(wholePathTxt)
	if err != nil {
		return "", err
	}

	return string(caption), nil
}

func (r *vkRepository) getUploadServer(id int) (string, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.vk.com/method/photos.getUploadServer?&access_token=%s&album_id=%d&v=5.199", r.token, id))
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

func (r *vkRepository) createAlbum(title string) (int, error) {
	id, err := r.getAlbumId(title)
	if err == nil {
		return id, nil
	}

	resp, err := http.Get(fmt.Sprintf("https://api.vk.com/method/photos.createAlbum?title=%s&access_token=%s&v=5.199", title, r.token))
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

	return msg.Response.Id, nil
}

func (r *vkRepository) getAlbumId(title string) (int, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.vk.com/method/photos.getAlbums?access_token=%s&v=5.199", r.token))
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

func (r *vkRepository) SetToken(token string) {
	r.token = token
}

func (r *vkRepository) SetId(id int) {
	r.id = id
}
