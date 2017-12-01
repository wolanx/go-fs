/**
 * 图片服务器
 */
package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"runtime/debug"
	"strings"
)

const (
	ListDir      = 0x0001
	TEMP_DIR     = "./temp"
	UPLOAD_DIR   = "./uploads"
	TEMPLATE_DIR = "./views"
)

var templates map[string]*template.Template

func init() {
	fileInfoArr, err := ioutil.ReadDir(TEMPLATE_DIR)
	check(err)

	templates = make(map[string]*template.Template)
	var templateName, templatePath string
	for _, fileInfo := range fileInfoArr {
		templateName = fileInfo.Name()
		if ext := path.Ext(templateName); ext != ".html" {
			continue
		}
		templatePath = TEMPLATE_DIR + "/" + templateName
		//log.Println("Loading template: ", templatePath)
		t := template.Must(template.ParseFiles(templatePath))
		templates[templatePath] = t
	}
}

func main() {
	mux := http.NewServeMux()
	staticDirHandler(mux, "/assets/", "./public", 0)
	mux.HandleFunc("/", safeHandler(listHandler))
	mux.HandleFunc("/view", safeHandler(viewHandler))
	mux.HandleFunc("/demo", safeHandler(demoHandler))
	mux.HandleFunc("/upload", safeHandler(uploadHandler))
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func readerHtml(w http.ResponseWriter, tmpl string, locals map[string]interface{}) {
	tmpl = TEMPLATE_DIR + "/" + tmpl + ".html"
	err := templates[tmpl].Execute(w, locals)
	check(err)
}

// func readerHtml(w http.ResponseWriter, tmpl string, locals map[string]interface{}) (err error) {
// 	t, err := template.ParseFiles(TEMPLATE_DIR + "/" + tmpl + ".html")
// 	if err != nil {
// 		return
// 	}
// 	err = t.Execute(w, locals)
// 	return
// }

func isExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return os.IsExist(err)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	imageId := r.FormValue("id")
	imagePath := UPLOAD_DIR + "/" + imageId
	if exists := isExists(imagePath); !exists {
		http.NotFound(w, r)
		return
	}
	//w.Header().Set("Content-Type", "application/pdf")
	http.ServeFile(w, r, imagePath)
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	fileInfoArr, err := ioutil.ReadDir("./uploads")
	check(err)
	locals := make(map[string]interface{})
	images := []string{}
	for _, fileInfo := range fileInfoArr {
		images = append(images, fileInfo.Name())
	}
	locals["hostname"], _ = os.Hostname()
	locals["images"] = images
	readerHtml(w, "list", locals)
}

func safeHandler(fn http.HandlerFunc) http.HandlerFunc {
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

func staticDirHandler(mux *http.ServeMux, prefix string, staticDir string, flags int) {
	mux.HandleFunc(prefix, func(w http.ResponseWriter, r *http.Request) {
		file := staticDir + r.URL.Path[len(prefix)-1:]
		if (flags & ListDir) == 0 {
			if exists := isExists(file); !exists {
				http.NotFound(w, r)
				return
			}
		}
		http.ServeFile(w, r, file)
		return
	})
}

//@see https://developer.qiniu.com/kodo/manual/1208/upload-token
func checkToken(token string) bool {
	tokenArr := strings.Split(token, ":")
	if len(tokenArr) != 3 {
		return false
	}
	accessKey := tokenArr[0]
	sign := tokenArr[1]
	policy := tokenArr[2]
	log.Println(accessKey, sign, policy)

	return true
}

func demoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		readerHtml(w, "upload", nil)
	}
	if r.Method == "POST" {
		w.Header().Set("content-type", "application/json")

		accessKey := "asdsad"
		sign := "123"
		policy := "xzc"

		ret := map[string]interface{}{}
		ret["token"] = accessKey + ":" + sign + ":" + policy
		str, _ := json.Marshal(ret)
		w.Write(str)
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		origin := r.Header.Get("Origin")
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("content-type", "application/json")

		ret := map[string]interface{}{}

		token := r.PostFormValue("token")
		if !checkToken(token) {
			ret["msg"] = "token false"
			str, _ := json.Marshal(ret)
			w.WriteHeader(500)
			w.Write(str)
			return
		}

		file, handle, err := r.FormFile("file")
		defer file.Close()
		check(err)

		upload_name := handle.Filename
		ext := path.Ext(upload_name) // .png
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
		temp_file, err := ioutil.TempFile(TEMP_DIR, upload_name)
		defer temp_file.Close()
		//defer os.Remove(temp_file.Name()) // temp/favicon.ico395854444
		check(err)
		_, err = io.Copy(temp_file, file)
		check(err)
		temp_file.Seek(0, 0)
		temp_file.Sync()

		// md5
		m := md5.New()
		io.Copy(m, temp_file)
		md5_hex := m.Sum([]byte(""))
		md5_name := fmt.Sprintf("%x", md5_hex)

		temp_file.Seek(0, 0)

		new_name := string(md5_name) + ext
		log.Println(new_name)
		// 新建文件
		new_file, err := os.Create(UPLOAD_DIR + "/" + new_name)
		check(err)
		defer new_file.Close()
		_, err = io.Copy(new_file, temp_file)
		check(err)
		err = new_file.Sync()
		check(err)

		ret["name"] = new_name
		str, _ := json.Marshal(ret)
		w.Write(str)
	}
}
