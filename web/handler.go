package taoweb

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/anyongjitiger/photo-backup-server/log"
	"github.com/anyongjitiger/photo-backup-server/utils"
	mux "github.com/julienschmidt/httprouter"
)

func Index(w http.ResponseWriter, r *http.Request, _ mux.Params) {
	log.Info("show index.")
	fmt.Fprintf(w, "<h1>Hello, welcome to my blog</h1>")
}

func HandleError(err error) {
	if err != nil {
		panic(err)
	}
}

func Get(w http.ResponseWriter, r *http.Request, _ mux.Params) {

	var post Post

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	HandleError(err)

	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	// Save JSON to Post struct
	if err := json.Unmarshal(body, &post); err != nil {

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)

		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
}

func PostIndex(w http.ResponseWriter, r *http.Request, _ mux.Params) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	var posts Posts = make([]Post, 0)

	post1 := Post{
		Id: 1,
		User: User{
			Id:       1001,
			Username: "kipling",
			Email:    "kipling0133@gmail.com",
		},
		Topic: "My First Post",
		Text:  "Hello everyone! .",
	}

	post2 := Post{
		Id: 2,
		User: User{
			Id:       1002,
			Username: "Markus",
			Email:    "MarkusLeeVip@gmail.com",
		},
		Topic: "My First Post",
		Text:  "Hello 世界",
	}
	posts = append(posts, post1)
	posts = append(posts, post2)
	json.NewEncoder(w).Encode(posts)
}

func PostCreate(w http.ResponseWriter, r *http.Request, _ mux.Params) {

	var post Post
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	HandleError(err)

	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	// Save JSON to Post struct
	if err := json.Unmarshal(body, &post); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
}

func TestHtml(w http.ResponseWriter, r *http.Request, _ mux.Params){
	cwd, _ := os.Getwd()
	t1, err := template.ParseFiles(filepath.Join( cwd, "../web/test.html" ))
	if err != nil {
		panic(err)
	}
	vcode := utils.GenValidateCode(4)
	log.Info(vcode)
	t1.Execute(w, vcode)
}