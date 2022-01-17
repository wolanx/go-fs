package handle

import (
	"net/http"

	"github.com/wolanx/go-fs/src/config"
	"github.com/wolanx/go-fs/src/lib"
)

func InfoHandler(w http.ResponseWriter, r *http.Request) {
	imageId := r.FormValue("id")
	imagePath := config.UploadDir + "/" + imageId
	if exists := lib.IsExists(imagePath); !exists {
		http.NotFound(w, r)
		return
	}
	//w.Header().Set("Content-Type", "application/pdf")
	http.ServeFile(w, r, imagePath)
}
