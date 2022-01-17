# 文件上传
------

## 上传step

1. server端获取token
1. 将token+file上传给upload
1. 获得name
1. 通过配置设置后续事件

## use
docker push ghcr.io/wolanx/go-fs

## deploy

### 方式1
```shell
# pre create dir to persistent storage
mkdir s1
cd s1

# test
docker run --name go-fs -it -d -p 22016:8080 -v "$PWD":/app/uploads -e DEBUG=true ghcr.io/wolanx/go-fs

# prod to set your env
docker run --name go-fs -it -d -p 22016:8080 -v "$PWD":/app/uploads \
 -e ACCESS_KEY=YourPublicKey -e SECRET_KEY=YourPrivateKey ghcr.io/wolanx/go-fs

# demo
docker run --name go-fs -it -d -p 22016:8080 -v "$PWD":/app/uploads \
 -e ACCESS_KEY=wolanx -e SECRET_KEY=wolanxkey -e URL_PATH=https://s1.wolanx.com ghcr.io/wolanx/go-fs
```

### 方式2
```text
version: "3"
services:
  upload:
    image: ghcr.io/wolanx/go-fs
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

#nginx

```text
server {
    listen       80;
    server_name  s1.wolanx.com;

    charset utf8;

    root /www/s1;
    expires 30d;

    location /.well-known {
      root /www/certbot;
    }

    location /upload {
        proxy_redirect      off;
        proxy_set_header    Host      $host;
        proxy_set_header    X-Real-IP $remote_addr;
        proxy_set_header    X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_pass          http://127.0.0.1:22016;
    }
}
server {
    listen       443 ssl;
    server_name  s1.wolanx.com;

    charset utf8;

    root /www/s1;
    expires 30d;

    ssl_certificate           /www/certbot/ssl/live/www.wolanx.com/fullchain.pem;
    ssl_certificate_key       /www/certbot/ssl/live/www.wolanx.com/privkey.pem;
    ssl_session_timeout       5m;
    ssl_ciphers               ECDHE-RSA-AES128-GCM-SHA256:ECDHE:ECDH:AES:HIGH:!NULL:!aNULL:!MD5:!ADH:!RC4;
    ssl_protocols             TLSv1 TLSv1.1 TLSv1.2;
    ssl_prefer_server_ciphers on;

    location /upload {
        proxy_redirect      off;
        proxy_set_header    Host      $host;
        proxy_set_header    X-Real-IP $remote_addr;
        proxy_set_header    X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_pass          http://127.0.0.1:22016;
    }
}
```
