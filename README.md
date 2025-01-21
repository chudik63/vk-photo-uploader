# Utility for conveniently uploading
This app allows you to upload folders with photos from your computer to a VK profile, automatically creating albums with the same titles.

# How to use it
## Installation
```
git clone https://github.com/chudik63/vk-photo-uploader.git
cd vk-photo-uploader
``` 
## Building
```
docker build -t uploader .
docker run --name uploader -p 80:80 uploader 
```

Then go to localhost and enjoy

# Integrated with:
- OAuth VK ID

# Technologies used:
- Golang (Gin)
- JavaScript, HTML, CSS
- Docker
