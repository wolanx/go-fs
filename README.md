# go-fs
file server

```
/**
 * 文件服务器
 */
package main

import (
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./")))
	http.ListenAndServe(":8080", nil)
}
```

# deploy
```
docker build -t go-img-1:v1 .

docker run -it --name go-app-1 go-img-1:v1

docker run --rm -v "$PWD":/usr/src/myapp -w /usr/src/myapp -e GOOS=darwin golang:1.8.1-alpine go build -v
docker run --rm -v "$PWD":"$PWD" -w "$PWD" -e GOOS=darwin golang:1.8.1-alpine go build -v
```