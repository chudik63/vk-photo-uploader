package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
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

type postUrlResponse struct {
	Server      int    `json:"server"`
	Photos_list string `json:"photos_list"`
	Aid         int    `json:"aid"`
	Hash        string `json:"string"`
}

type savePhotoResponse struct {
	Response struct {
		Album_id int `json:"album_id"`
	} `json:"response"`
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

	log.Print(id)

	upload_url, err := r.getUploadServer(id)
	if err != nil {
		return err
	}

	wholePath, _ := os.Getwd()
	wholePath = filepath.Join(wholePath, r.root+path)

	b := &bytes.Buffer{}
	writer := multipart.NewWriter(b)

	file, err := os.Open(r.root + path)
	if err != nil {
		return err
	}
	defer file.Close()

	part, _ := writer.CreateFormFile("file1", wholePath)

	io.Copy(part, file)

	writer.Close()

	req, _ := http.NewRequest("POST", upload_url, b)

	client := &http.Client{}
	resp, _ := client.Do(req)

	body, _ := io.ReadAll(resp.Body)

	msg := &postUrlResponse{}
	json.Unmarshal(body, msg)

	log.Print(msg)

	resp, err = http.Get(fmt.Sprintf("https://api.vk.com/method/photos.save?&album_id=%d&server=%d&photos_list=%s&hash=%s&access_token=%s&v=5.131", id, msg.Server, msg.Photos_list, msg.Hash, r.token))
	if err != nil {
		return err
	}
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	msg1 := &savePhotoResponse{}

	err = json.Unmarshal(body, msg1)
	if err != nil {
		return err
	}
	log.Print(msg1.Response.Album_id)
	return nil
}

func (r *VkRepository) getUploadServer(id int) (string, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.vk.com/method/photos.getUploadServer?&access_token=%s&album_id=%d&v=5.199", r.token, id))
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	msg := &getUploadServerResponse{}

	err = json.Unmarshal(body, msg)
	if err != nil {
		return "", err
	}
	return msg.Response.UploadUrl, nil
}

func (r *VkRepository) createAlbum(title string) (int, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.vk.com/method/photos.createAlbum?title=%s&access_token=%s&v=5.199", title, r.token))
	if err != nil {
		return 0, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	msg := &albumCreationResponse{}

	err = json.Unmarshal(body, msg)
	if err != nil {
		return 0, err
	}

	return msg.Response.Id, nil
}

func (r *VkRepository) SetToken(token string) {
	r.token = token
	log.Print(token)
}

func (r *VkRepository) SetId(id int) {
	r.id = id
}
