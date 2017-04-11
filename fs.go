/**
 * 图片服务器
 */
package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"runtime/debug"
	"io"
	"crypto/md5"
	"fmt"
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
		log.Println("Loading template: ", templatePath)
		t := template.Must(template.ParseFiles(templatePath))
		templates[templatePath] = t
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

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		readerHtml(w, "upload", nil)
	}
	if r.Method == "POST" {
		f, h, err := r.FormFile("image")
		defer f.Close()
		check(err)

		uplode_name := h.Filename

		// 保存临时文件
		temp_file, err := ioutil.TempFile(TEMP_DIR, uplode_name) // temp_file.Name() temp\tx.jpg309941499
		defer temp_file.Close()
		check(err)
		_, err = io.Copy(temp_file, f)
		check(err)
		temp_file.Seek(0, 0)
		temp_file.Sync()

		// md5
		m := md5.New()
		io.Copy(m, temp_file)
		md5_hex := m.Sum([]byte(""))

		md5_name := fmt.Sprintf("%x", md5_hex)
		log.Printf(md5_name)

		temp_file.Seek(0, 0)

		log.Println(uplode_name)
		// 新建文件
		new_file, err := os.Create(UPLOAD_DIR + "/" + string(md5_name))
		defer new_file.Close()
		check(err)
		_, err = io.Copy(new_file, temp_file)
		check(err)
		err = new_file.Sync()
		check(err)

		http.Redirect(w, r, "/upload", http.StatusFound)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	imageId := r.FormValue("id")
	imagePath := UPLOAD_DIR + "/" + imageId
	if exists := isExists(imagePath); !exists {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "image")
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

func main() {
	mux := http.NewServeMux()
	staticDirHandler(mux, "/assets/", "./public", 0)
	mux.HandleFunc("/", safeHandler(listHandler))
	mux.HandleFunc("/view", safeHandler(viewHandler))
	mux.HandleFunc("/upload", safeHandler(uploadHandler))
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}
