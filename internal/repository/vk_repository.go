package repository

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type VkRepository struct {
	root  string
	token string
	id    int
}

type albumCreationResponse struct {
	Response struct {
		Id int `json:"id"`
	} `json:"response"`
}

type getUploadServerResponse struct {
	Response struct {
		UploadUrl string `json:"upload_url"`
	} `json:"response"`
}

type getAlbumsResponse struct {
	Response struct {
		Count int `json:"count"`
		Array []struct {
			Id    int    `json:"id"`
			Title string `json:"title"`
		} `json:"items"`
	} `json:"response"`
}

type postUrlResponse struct {
	Server      int    `json:"server"`
	Photos_list string `json:"photos_list"`
	Aid         int    `json:"aid"`
	Hash        string `json:"hash"`
}

type savePhotoResponse struct {
	Response []interface{} `json:"response"`
}

func NewVkRepository(path string) VkRepository {
	return VkRepository{
		root:  path,
		token: "",
		id:    0,
	}
}

func (r *VkRepository) Upload(path string) error {
	dir := filepath.Dir(path)

	id, err := r.createAlbum(dir)
	if err != nil {
		return err
	}

	upload_url, err := r.getUploadServer(id)
	if err != nil {
		return err
	}

	wholePath, _ := os.Getwd()
	wholePath = filepath.Join(wholePath, r.root, path)

	caption, err := r.readCaption(wholePath)
	if err != nil {
		return err
	}

	b := &bytes.Buffer{}
	writer := multipart.NewWriter(b)

	file, err := os.Open(wholePath)
	if err != nil {
		return err
	}
	defer file.Close()

	part, err := writer.CreateFormFile("file1", filepath.Base(wholePath))
	if err != nil {
		return err
	}

	if _, err = io.Copy(part, file); err != nil {
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

	msg := &postUrlResponse{}
	if err = json.Unmarshal(body, msg); err != nil {
		return err
	}

	resp, err = http.Get(fmt.Sprintf("https://api.vk.com/method/photos.save?&album_id=%d&server=%d&photos_list=%s&hash=%s&caption=%s&access_token=%s&v=5.199", msg.Aid, msg.Server, msg.Photos_list, msg.Hash, url.QueryEscape(caption), r.token))
	if err != nil {
		return err
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	/// ПРОВЕРИТЬ загрузку
	msg1 := &savePhotoResponse{}

	err = json.Unmarshal(body, msg1)
	if err != nil {
		return err
	}

	return nil
}

func (r *VkRepository) readCaption(path string) (string, error) {
	ext := filepath.Ext(path)
	wholePathWithoutExt := strings.TrimSuffix(path, ext)

	wholePathTxt := wholePathWithoutExt + ".txt"

	caption, err := os.ReadFile(wholePathTxt)
	if err != nil {
		return "", err
	}

	return string(caption), nil
}

func (r *VkRepository) getUploadServer(id int) (string, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.vk.com/method/photos.getUploadServer?&access_token=%s&album_id=%d&v=5.199", r.token, id))
	if err != nil {
		return "", errors.New("ошибка получения сервера загрузки")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("нельзя считать ответ создания альбома")
	}

	msg := &getUploadServerResponse{}

	err = json.Unmarshal(body, msg)
	if err != nil {
		return "", errors.New("ошибка чтения JSON после получения сервера загрузки")
	}
	return msg.Response.UploadUrl, nil
}

func (r *VkRepository) createAlbum(title string) (int, error) {
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

	msg := &albumCreationResponse{}

	err = json.Unmarshal(body, msg)
	if err != nil {
		return 0, errors.New("ошибка чтения JSON после создании альбома")
	}

	return msg.Response.Id, nil
}

func (r *VkRepository) getAlbumId(title string) (int, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.vk.com/method/photos.getAlbums?access_token=%s&v=5.199", r.token))
	if err != nil {
		return 0, errors.New("нельзя получить информацию о альбомах VK")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, errors.New("нельзя считать ответ c информацией о альбомах")
	}

	msg := &getAlbumsResponse{}
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

func (r *VkRepository) SetToken(token string) {
	r.token = token
}

func (r *VkRepository) SetId(id int) {
	r.id = id
}
