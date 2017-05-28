default:
	docker run --rm -v $(shell pwd):/usr/src/myapp -w /usr/src/myapp golang:1.8.1-alpine go build -v
	docker build -t go-fs-img-1:v1 .