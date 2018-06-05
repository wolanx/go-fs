package handle

import (
	"log"
	"net/http"
	"os"

	"github.com/zx5435/go-fs/src/config"
	"github.com/zx5435/go-fs/src/lib"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	imageId := r.URL.Path // /6005f38d6f4160d3f15da8d7673102b0.json

	log.Printf("indexHandler imageId:'%s'", imageId)
	if imageId == "/" {
		locals := make(map[string]interface{})
		locals["debug"] = lib.Debug
		locals["hostname"], _ = os.Hostname()
		lib.ReaderHtml(w, "index", locals)
	} else {
		log.Println(imageId)
		imagePath := config.UploadDir + "/" + imageId
		if exists := lib.IsExists(imagePath); !exists {
			http.NotFound(w, r)
			return
		}
		//w.Header().Set("Content-Type", "application/pdf")
		http.ServeFile(w, r, imagePath)
	}
}
