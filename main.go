/**
 * 图片服务器
 */
package main

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/zx5435/go-fs/src/config"
	"github.com/zx5435/go-fs/src/handle"
)

func main() {
	mux := http.NewServeMux()
	handle.StaticHandler(mux, "/assets/", "./assets", 0)
	mux.HandleFunc("/list", SafeHandler(handle.ListHandler))
	mux.HandleFunc("/info", SafeHandler(handle.InfoHandler))
	mux.HandleFunc("/demo", SafeHandler(handle.DemoHandler))
	mux.HandleFunc("/upload", SafeHandler(handle.UploadHandler))
	mux.HandleFunc("/", SafeHandler(handle.IndexHandler))
	err := http.ListenAndServe(":"+config.Port, mux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}

func SafeHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e, ok := recover().(error); ok {
				http.Error(w, e.Error(), http.StatusInternalServerError)
				log.Println("Warn:panic in %v. - %v", fn, e)
				log.Println(string(debug.Stack()))
			}
		}()
		fn(w, r)
	}
}
