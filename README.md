

## deploy

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
