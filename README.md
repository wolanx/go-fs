## 上传step

1. server端获取token
1. 将token+file上传给upload
1. 获得name
1. 通过配置设置后续事件

## build
```text
make
docker tag go-fs:v1 zx5435/go-fs:v1
docker push zx5435/go-fs:v1
```

## deploy

### manual
```text
make
docker run -it -d -p 22016:8080 --name go-fs go-fs:v1
```

### deploy yml
```text
version: "3"
services:
  upload:
    image: zx5435/go-fs:v1
    volumes:
      - ./upload:/app/uploads
    environment:
      - ACCESS_KEY=MDev1
      - SECRET_KEY=MDev2
      - URL_PATH=/creatives
    networks:
      - mynet
networks:
  mynet:
volumes:
  mydir:
```
