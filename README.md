# Utility for easier uploading
This app allows to upload folders with photos from a computer to VK automatically creating albums with the same titles as folders` ones

# How to use it
Run:
```
git clone github.com/chudik63/vk-photo-uploader
cd vk-photo-uploader
docker build -t uploader .
docker run uploader
``` 

Then go to localhost:80 and enjoy


# Technologies used:
- Golang (Gin)
- JavaScript, HTML, CSS
- Docker
