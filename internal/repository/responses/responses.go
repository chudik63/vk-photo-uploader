package responses

type AlbumCreateResponse struct {
	Response struct {
		Id int `json:"id"`
	} `json:"response"`
}

type GetUploadServerResponse struct {
	Response struct {
		UploadUrl string `json:"upload_url"`
	} `json:"response"`
}

type GetAlbumsResponse struct {
	Response struct {
		Count int `json:"count"`
		Array []struct {
			Id    int    `json:"id"`
			Title string `json:"title"`
		} `json:"items"`
	} `json:"response"`
}

type PostUrlResponse struct {
	Server      int    `json:"server"`
	Photos_list string `json:"photos_list"`
	Aid         int    `json:"aid"`
	Hash        string `json:"hash"`
}

type SavePhotoResponse struct {
	Response []struct {
		Album_id int `json:"album_id"`
		Id       int `json:"id"`
	} `json:"response"`
}
