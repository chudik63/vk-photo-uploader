package repository

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
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

	url, err := r.getUploadServer(id)
	if err != nil {
		return err
	}

	file, _ := os.Open(r.root + path)

	_, err = http.Post(url, "file1=@"+r.root+path, file)
	if err != nil {
		return err
	}

	return nil
}

func (r *VkRepository) getUploadServer(id int) (string, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.vk.com/method/photos.getUploadServer?&access_token=%s&album_id=%d&v=5.30", r.token, id))
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

	log.Print(msg)

	return msg.Response.Id, nil
}

func (r *VkRepository) SetToken(token string) {
	r.token = token
	log.Print(token)
}

func (r *VkRepository) SetId(id int) {
	r.id = id
}
