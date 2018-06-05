package handle

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/zx5435/go-fs/src/lib"
)

//@ref https://developer.qiniu.com/kodo/manual/1208/upload-token
func DemoHandler(w http.ResponseWriter, r *http.Request) {
	if !lib.Debug {
		w.Write([]byte("no run in debug false"))
		return
	}

	if r.Method == "GET" {
		lib.ReaderHtml(w, "upload", nil)
	}
	if r.Method == "POST" {
		w.Header().Set("content-type", "application/json")

		filename := r.PostFormValue("filename")

		accessKey := lib.MyAccessKey
		policy, _ := json.Marshal(&lib.Policy{
			Filename: filename,
			Deadline: time.Now().Unix(),
		})
		policyStr := base64.StdEncoding.EncodeToString(policy)
		mac := hmac.New(sha1.New, []byte(lib.MySecretKey))
		mac.Write([]byte(policyStr))
		sign := mac.Sum(nil)

		log.Println("lib.MyAccessKey", lib.MyAccessKey, "lib.MySecretKey", lib.MySecretKey)

		ret := map[string]interface{}{}
		ret["token"] = accessKey + ":" + base64.StdEncoding.EncodeToString(sign) + ":" + policyStr
		str, _ := json.Marshal(ret)
		w.Write(str)
	}
}
