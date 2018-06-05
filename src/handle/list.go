package handle

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/zx5435/go-fs/src/lib"
)

func ListHandler(w http.ResponseWriter, r *http.Request) {
	if !lib.Debug {
		w.Write([]byte("no run in debug false"))
		return
	}

	fileInfoArr, err := ioutil.ReadDir("./uploads")
	lib.Check(err)
	locals := make(map[string]interface{})
	images := []string{}
	for _, fileInfo := range fileInfoArr {
		images = append(images, fileInfo.Name())
	}
	locals["hostname"], _ = os.Hostname()
	locals["images"] = images
	lib.ReaderHtml(w, "list", locals)
}
