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

# Dockerfile
```
FROM alpine:3.6

WORKDIR /myapp
COPY . /myapp/

CMD ["./myapp"]
```

# deploy
```
make
docker run -it -d -p 8080:8080 --name go-fs-app-1 go-fs-img-1:v1
```

# hook
```
http://139.196.14.10:8080/github-webhook/
```