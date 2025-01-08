# Utility for conveniently uploading
This app allows you to upload folders with photos from your computer to a VK profile, automatically creating albums with the folder names.

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
