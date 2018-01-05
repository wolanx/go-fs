/**
 * 图片服务器
 */
package main

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
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
	"time"
)

const (
	Port        = "8080"
	UploadDir   = "./uploads"
	TemplateDir = "./views"
	TempDir     = "./temp"
	ListDir     = 0x0001
)

var templates map[string]*template.Template
var keyArr map[string]interface{}
var myAccessKey string
var mySecretKey string

func init() {
	myAccessKey = os.Getenv("ACCESS_KEY")
	mySecretKey = os.Getenv("SECRET_KEY")
	if myAccessKey == "" {
		myAccessKey = "MY_test1"
	}
	if mySecretKey == "" {
		mySecretKey = "MY_test2"
	}

	log.Println(myAccessKey, mySecretKey)

	keyArr = make(map[string]interface{})
	//keyArr["ACCESS_KEY"] = "SECRET_KEY"
	keyArr[myAccessKey] = mySecretKey

	fileInfoArr, err := ioutil.ReadDir(TemplateDir)
	check(err)

	templates = make(map[string]*template.Template)
	var templateName, templatePath string
	for _, fileInfo := range fileInfoArr {
		templateName = fileInfo.Name()
		if ext := path.Ext(templateName); ext != ".html" {
			continue
		}
		templatePath = TemplateDir + "/" + templateName
		//log.Println("Loading template: ", templatePath)
		t := template.Must(template.ParseFiles(templatePath))
		templates[templatePath] = t
	}
}

func main() {
	mux := http.NewServeMux()
	staticDirHandler(mux, "/assets/", "./assets", 0)
	mux.HandleFunc("/list", safeHandler(listHandler))
	mux.HandleFunc("/info", safeHandler(infoHandler))
	mux.HandleFunc("/demo", safeHandler(demoHandler))
	mux.HandleFunc("/upload", safeHandler(uploadHandler))
	mux.HandleFunc("/", safeHandler(indexHandler))
	err := http.ListenAndServe(":"+Port, mux)
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
	tmpl = TemplateDir + "/" + tmpl + ".html"
	err := templates[tmpl].Execute(w, locals)
	check(err)
}

func isExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return os.IsExist(err)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	imageId := r.URL.Path // /6005f38d6f4160d3f15da8d7673102b0.json

	log.Printf("indexHandler imageId:'%s'", imageId)
	if imageId == "/" {
		locals := make(map[string]interface{})
		locals["hostname"], _ = os.Hostname()
		readerHtml(w, "index", locals)
	} else {
		log.Println(imageId)
		imagePath := UploadDir + "/" + imageId
		if exists := isExists(imagePath); !exists {
			http.NotFound(w, r)
			return
		}
		//w.Header().Set("Content-Type", "application/pdf")
		http.ServeFile(w, r, imagePath)
	}
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

func infoHandler(w http.ResponseWriter, r *http.Request) {
	imageId := r.FormValue("id")
	imagePath := UploadDir + "/" + imageId
	if exists := isExists(imagePath); !exists {
		http.NotFound(w, r)
		return
	}
	//w.Header().Set("Content-Type", "application/pdf")
	http.ServeFile(w, r, imagePath)
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

type SaveConfig struct {
	Engine string `json:"engine"`
}

type Policy struct {
	Filename   string     `json:"filename"`
	Deadline   int64      `json:"deadline"`
	SaveConfig SaveConfig `json:"saveConfig"`
}

func checkToken(token string) (policy Policy, err error) {
	tokenArr := strings.Split(token, ":")
	if len(tokenArr) != 3 {
		err = errors.New("no 3")
		return
	}
	accessKey := tokenArr[0]
	sign := tokenArr[1]
	policyStr, _ := base64.StdEncoding.DecodeString(tokenArr[2])
	policy = Policy{}
	json.Unmarshal(policyStr, &policy)
	//log.Println(accessKey, sign, string(policyStr), policy)

	secretKey := keyArr[accessKey]
	if secretKey == nil {
		err = errors.New("secretKey not exist")
		return
	}

	mac := hmac.New(sha1.New, []byte(secretKey.(string)))
	mac.Write([]byte(tokenArr[2]))
	signCheck := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	if sign != signCheck {
		err = errors.New("sign not eq")
		return
	}

	return policy, nil
}

//@ref https://developer.qiniu.com/kodo/manual/1208/upload-token
func demoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		readerHtml(w, "upload", nil)
	}
	if r.Method == "POST" {
		w.Header().Set("content-type", "application/json")

		filename := r.PostFormValue("filename")

		accessKey := myAccessKey
		policy, _ := json.Marshal(&Policy{
			Filename: filename,
			Deadline: time.Now().Unix(),
		})
		policyStr := base64.StdEncoding.EncodeToString(policy)
		mac := hmac.New(sha1.New, []byte(mySecretKey))
		mac.Write([]byte(policyStr))
		sign := mac.Sum(nil)

		ret := map[string]interface{}{}
		ret["token"] = accessKey + ":" + base64.StdEncoding.EncodeToString(sign) + ":" + policyStr
		str, _ := json.Marshal(ret)
		w.Write(str)
	}
}

/**
 * 上传 post
 */
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		origin := r.Header.Get("Origin")
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("content-type", "application/json")

		ret := map[string]interface{}{}

		token := r.PostFormValue("token")
		config, err := checkToken(token)
		if err != nil {
			ret["msg"] = err.Error()
			str, _ := json.Marshal(ret)
			w.WriteHeader(500)
			w.Write(str)
			return
		}

		file, handle, err := r.FormFile("file")
		defer file.Close()
		check(err)

		uploadName := handle.Filename
		log.Println("config: ", config)
		if uploadName != config.Filename {
			ret["msg"] = "filename not match"
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
		tempFile, err := ioutil.TempFile(TempDir, uploadName)
		defer tempFile.Close()
		//defer os.Remove(tempFile.Name()) // temp/favicon.ico395854444
		check(err)
		_, err = io.Copy(tempFile, file)
		check(err)
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
		newFile, err := os.Create(UploadDir + "/" + newName)
		check(err)
		defer newFile.Close()
		_, err = io.Copy(newFile, tempFile)
		check(err)
		err = newFile.Sync()
		check(err)

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
