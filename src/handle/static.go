package handle

import (
	"net/http"

	"github.com/zx5435/go-fs/src/config"
	"github.com/zx5435/go-fs/src/lib"
)

func StaticHandler(mux *http.ServeMux, prefix string, staticDir string, flags int) {
	mux.HandleFunc(prefix, func(w http.ResponseWriter, r *http.Request) {
		file := staticDir + r.URL.Path[len(prefix)-1:]
		if (flags & config.ListDir) == 0 {
			if exists := lib.IsExists(file); !exists {
				http.NotFound(w, r)
				return
			}
		}
		http.ServeFile(w, r, file)
		return
	})
}
