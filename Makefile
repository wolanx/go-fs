default:
	cat Makefile

build:
	docker run --rm -v $(shell pwd):/usr/local/go/src/github.com/zx5435/go-fs -w /usr/local/go/src/github.com/zx5435/go-fs golang:1.10.2-alpine go build -v
	docker build -t zx5435/go-fs:v1 .

clear:
	ls
