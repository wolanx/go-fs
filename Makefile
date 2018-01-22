default:

build:
	docker run --rm -v $(shell pwd):/app -w /app golang:1.9.1-alpine go build -v
	docker build -t zx5435/go-fs:v1 .
