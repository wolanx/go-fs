package lib

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/zx5435/go-fs/src/config"
)

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

type SaveConfig struct {
	Engine string `json:"engine"`
}

type Policy struct {
	Filename   string     `json:"filename"`
	Deadline   int64      `json:"deadline"`
	SaveConfig SaveConfig `json:"saveConfig"`
}

var MyAccessKey string
var MySecretKey string

var keyArr map[string]interface{}
var templates map[string]*template.Template

func init() {
	MyAccessKey = os.Getenv("ACCESS_KEY")
	MySecretKey = os.Getenv("SECRET_KEY")
	if MyAccessKey == "" {
		MyAccessKey = "MY_test1"
	}
	if MySecretKey == "" {
		MySecretKey = "MY_test2"
	}

	log.Println(MyAccessKey, MySecretKey)

	keyArr = make(map[string]interface{})
	//keyArr["ACCESS_KEY"] = "SECRET_KEY"
	keyArr[MyAccessKey] = MySecretKey

	fileInfoArr, err := ioutil.ReadDir(config.TemplateDir)
	Check(err)

	templates = make(map[string]*template.Template)
	var templateName, templatePath string
	for _, fileInfo := range fileInfoArr {
		templateName = fileInfo.Name()
		if ext := path.Ext(templateName); ext != ".html" {
			continue
		}
		templatePath = config.TemplateDir + "/" + templateName
		//log.Println("Loading template: ", templatePath)
		t := template.Must(template.ParseFiles(templatePath))
		templates[templatePath] = t
	}
}

func ReaderHtml(w http.ResponseWriter, tmpl string, locals map[string]interface{}) {
	tmpl = config.TemplateDir + "/" + tmpl + ".html"
	err := templates[tmpl].Execute(w, locals)
	Check(err)
}

func IsExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return os.IsExist(err)
}

func CheckToken(token string) (policy Policy, err error) {
	tokenArr := strings.Split(token, ":")
	if len(tokenArr) != 3 {
		err = errors.New("no 3")
		return
	}
	accessKey := tokenArr[0]
	sign := tokenArr[1]
	policyStr, _ := base64.StdEncoding.DecodeString(tokenArr[2])
	json.Unmarshal(policyStr, &policy)

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
