default:
	docker run --rm -v $(shell pwd):/app -w /app golang:1.9.1-alpine go build -v
	docker build -t go-fs:v1 .
