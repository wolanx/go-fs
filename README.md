# 文件上传
------

## 上传step

1. server端获取token
1. 将token+file上传给upload
1. 获得name
1. 通过配置设置后续事件

## build
```text
make build
docker push zx5435/go-fs:v1
```

## deploy

### 方式1
```text
// pre
mkdir s1
cd s1

// test
docker run -it -d -p 22016:8080 -v "$PWD":/app/uploads \
 --name go-fs -e DEBUG=true zx5435/go-fs:v1

// prod to set your env
docker run -it -d -p 22016:8080 -v "$PWD":/app/uploads \
 --name go-fs -e ACCESS_KEY=YourPublicKey -e SECRET_KEY=YourPrivateKey zx5435/go-fs:v1
```

### 方式2
```text
version: "3"
services:
  upload:
    image: zx5435/go-fs:v1
    volumes:
      - ./upload:/app/uploads
    environment:
      - DEBUG=true
      - ACCESS_KEY=YourPublicKey
      - SECRET_KEY=YourPrivateKey
      - URL_PATH=/creatives
    networks:
      - mynet
networks:
  mynet:
volumes:
  mydir:
```
