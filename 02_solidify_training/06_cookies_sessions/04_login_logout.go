package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
)

var store = sessions.NewCookieStore([]byte("secret-password"))
var tpl *template.Template

func init() {
	tpl, _ = template.ParseGlob("D:\\golang projects\\go_training\\02_solidify_training\\06_cookies_sessions\\assets\\templates\\*.html")
}

/*登出*/
func logout(resp http.ResponseWriter, req *http.Request) {
	//获取session
	session, _ := store.Get(req, "session")
	session.Values["login"] = false
	session.Save(req, resp)
	//重定向到登陆页面
	http.Redirect(resp, req, "/login", http.StatusFound)
}

/*登陆*/
func login(resp http.ResponseWriter, req *http.Request) {
	store, _ := store.Get(req, "session")
	if req.Method == "POST" && req.FormValue("password") == "secret" {
		store.Values["login"] = true
		store.Save(req, resp)
		http.Redirect(resp, req, "/", http.StatusFound)
		return
	}
	tpl.ExecuteTemplate(resp, "login.html", nil)
}

/*首页*/
func index(resp http.ResponseWriter, req *http.Request) {
	fmt.Printf("local path:%s\n", req.URL.Path)
	session, _ := store.Get(req, "session")
	isLogin := session.Values["login"]
	if isLogin == false || isLogin == nil {
		http.Redirect(resp, req, "/login", http.StatusFound)
		return
	}
	//获取data数据
	file, header, err := req.FormFile("data")
	if req.Method == "POST" && err == nil {
		//上传照片
		uploadPhoto(file, header, session)
	}
	session.Save(req, resp)
	data := getPhotos(session)
	tpl.ExecuteTemplate(resp, "index.html", data)
}

/*上传*/
func uploadPhoto(file multipart.File, header *multipart.FileHeader, session *sessions.Session) {
	defer file.Close()
	fileName := getSha(file) + ".jpg"
	localDir, _ := os.Getwd()
	path := filepath.Join(localDir, "assets", "imgs", fileName)
	create, _ := os.Create(path)
	defer create.Close()
	create.Seek(0, 0)
	io.Copy(create, file)
	addPhoto(fileName, session)
}

/*添加*/
func addPhoto(name string, session *sessions.Session) {
	data := getPhotos(session)
	data = append(data, name)
	marshal, _ := json.Marshal(data)
	session.Values["data"] = string(marshal)
}

func getPhotos(session *sessions.Session) []string {
	var data []string
	jsonData := session.Values["data"]
	if jsonData != nil {
		json.Unmarshal([]byte(jsonData.(string)), &data)
	}
	return data
}

/*加密*/
func getSha(file multipart.File) string {
	hash := sha1.New()
	io.Copy(hash, file)
	return fmt.Sprintf("%x", hash.Sum(nil))
}
func main() {
	//处理静态资源
	localDir, _ := os.Getwd()
	assetsPath := filepath.Join(localDir, "assets", "imgs")
	/*页面展示不了图片*/
	http.Handle("/assets/imgs/", http.StripPrefix("/assets/imgs", http.FileServer(http.Dir(assetsPath))))
	http.Handle("/favicon.ico", http.NotFoundHandler())
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/", index)
	serveMux.HandleFunc("/login", login)
	serveMux.HandleFunc("/logout", logout)
	http.ListenAndServe(":8080", context.ClearHandler(serveMux))
}
