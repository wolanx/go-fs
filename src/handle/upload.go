package handle

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/zx5435/go-fs/src/config"
	"github.com/zx5435/go-fs/src/lib"
)

/**
 * 上传 post
 */
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		origin := r.Header.Get("Origin")
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("content-type", "application/json")

		ret := map[string]interface{}{}

		token := r.PostFormValue("token")
		upConfig, err := lib.CheckToken(token)
		if err != nil {
			ret["msg"] = err.Error()
			str, _ := json.Marshal(ret)
			w.WriteHeader(500)
			w.Write(str)
			return
		}

		file, handle, err := r.FormFile("file")
		defer file.Close()
		lib.Check(err)

		uploadName := handle.Filename
		log.Println("upConfig: ", upConfig)
		if uploadName != upConfig.Filename {
			ret["msg"] = "filename not match," + uploadName + "," + upConfig.Filename
			str, _ := json.Marshal(ret)
			w.WriteHeader(500)
			w.Write(str)
			return
		}

		ext := path.Ext(uploadName) // .png
		arr := map[string]interface{}{
			".png":  "1",
			".jpg":  "1",
			".jpeg": "1",
			".pdf":  "1",
		}
		if arr[ext] == "" {
			ret["msg"] = "un support ext"
			str, _ := json.Marshal(ret)
			w.WriteHeader(500)
			w.Write(str)
			return
		}

		// 保存临时文件
		tempFile, err := ioutil.TempFile(config.TempDir, uploadName)
		defer tempFile.Close()
		//defer os.Remove(tempFile.Name()) // temp/favicon.ico395854444
		lib.Check(err)
		_, err = io.Copy(tempFile, file)
		lib.Check(err)
		tempFile.Seek(0, 0)
		tempFile.Sync()

		// md5
		m := md5.New()
		io.Copy(m, tempFile)
		md5_hex := m.Sum([]byte(""))
		md5_name := fmt.Sprintf("%x", md5_hex)

		tempFile.Seek(0, 0)

		newName := string(md5_name) + ext
		// 新建文件
		newFile, err := os.Create(config.UploadDir + "/" + newName)
		lib.Check(err)
		defer newFile.Close()
		_, err = io.Copy(newFile, tempFile)
		lib.Check(err)
		err = newFile.Sync()
		lib.Check(err)

		scheme := "http://"
		if r.TLS != nil {
			scheme = "https://"
		}

		urlPath := os.Getenv("URL_PATH")

		ret["key"] = newName
		ret["name"] = scheme + r.Host + urlPath + "/" + newName
		log.Println(ret["name"])
		str, _ := json.Marshal(ret)
		w.Write(str)
	}
}
