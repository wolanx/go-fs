#!/bin/bash

docker run --rm -v "$PWD":/usr/src/myapp -w /usr/src/myapp golang:1.8.1-alpine go build -v
docker build -t go-fs-img-1:v1 .

printf '🍻 %.0s' {1..30}
echo "end"
