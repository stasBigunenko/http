package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"src/http/pkg/model"
	"src/http/pkg/services"
	"src/http/storage"
	"strconv"
	"time"
)

//Handlers with the CRUD functions and Middleware

type PostHandler struct {
	router   *mux.Router
	services *services.Store
}

const MaxRequestSize = 2 * 1024

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("./pkg/templates/*"))
}

func New(router *mux.Router, store storage.Storage) *PostHandler {
	return &PostHandler{
		router:   router,
		services: services.NewStore(store),
	}
}

//Simple middleware function which write log in the terminal requested Method and URI
func simpleLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method)
		log.Println(r.RequestURI)
		fmt.Println("-------------")
		next.ServeHTTP(w, r)
	})
}

//Middleware: timeout to handler process
func processTimeout(h http.HandlerFunc, duration time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), duration)
		defer cancel()

		r = r.WithContext(ctx)

		processDone := make(chan bool)
		go func() {
			h(w, r)
			processDone <- true
		}()

		select {
		case <-ctx.Done():
			msg := services.Response("error process timeout")
			w.Write(msg)
		case <-processDone:
		}
	}
}

func (h *PostHandler) Routes() {
	h.router.HandleFunc("/post/", processTimeout(h.CreatePost, 5*time.Second)).Methods("POST")
	h.router.HandleFunc("/post/{id}", processTimeout(h.GetPost, 5*time.Second)).Methods("GET")
	h.router.HandleFunc("/posts", processTimeout(h.GetAll, 5*time.Second)).Methods("GET")
	h.router.HandleFunc("/post/{id}", processTimeout(h.DeletePost, 5*time.Second)).Methods("DELETE")
	h.router.HandleFunc("/post/{id}", processTimeout(h.UpdatePost, 5*time.Second)).Methods("PUT")
	h.router.HandleFunc("/post/upload", processTimeout(h.UploadPost, 5*time.Second)).Methods("POST")
	h.router.HandleFunc("/post/download", processTimeout(h.DownloadPost, 5*time.Second)).Methods("POST")
	h.router.Use(simpleLog)
}

//CreatePost Create post with decoding request and encoding response
func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		msg := services.Response("Method Not Allowed")
		w.WriteHeader(405)
		w.Write(msg)
		return
	}

	//receive requests only with size limit 2024 bytes
	r.Body = http.MaxBytesReader(w, r.Body, MaxRequestSize)

	var post model.Post

	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		msg := services.Response("Too big request.")
		w.WriteHeader(404)
		w.Write(msg)
		return
	}
	res, err := h.services.CreateId(&post)
	if err != nil {
		msg := services.Response("Could not create empty post")
		w.WriteHeader(406)
		w.Write(msg)
		return
	}
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(&res)
}

//GetPost Get post by Id with decoding request and encoding response
func (h *PostHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		msg := services.Response("Method Not Allowed")
		w.WriteHeader(405)
		w.Write(msg)
		return
	}

	//id, err := strconv.Atoi(r.URL.Query().Get("id"))
	vars := mux.Vars(r)
	key := vars["id"]
	id, err := strconv.Atoi(key)
	if err != nil {
		msg := services.Response("The Id is not valid")
		w.WriteHeader(406)
		w.Write(msg)
		return
	}
	res, err := h.services.GetId(id)
	if err != nil {
		msg := services.Response("This id doesn't exist")
		w.WriteHeader(404)
		w.Write(msg)
		return
	}
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(&res)
}

//GetAll posts by Id with decoding request and encoding response
func (h *PostHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		msg := services.Response("Method Not Allowed")
		w.WriteHeader(405)
		w.Write(msg)
		return
	}
	res, err := h.services.GetALL()
	if err != nil {
		msg := services.Response("Bad request")
		w.WriteHeader(404)
		w.Write(msg)
		return
	}

	if len(*res) == 0 {
		msg := services.Response("There is no posts in the memory.")
		w.WriteHeader(200)
		w.Write(msg)
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(&res)
}

//DeletePost post by Id with decoding request and encoding response
func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodDelete {
		msg := services.Response("Method Not Allowed")
		w.WriteHeader(405)
		w.Write(msg)
		return
	}
	//id, err := strconv.Atoi(r.URL.Query().Get("Id"))
	vars := mux.Vars(r)
	key := vars["id"]
	id, err := strconv.Atoi(key)
	if err != nil {
		msg := services.Response("The Id is not valid")
		w.WriteHeader(406)
		w.Write(msg)
		return
	}
	err = h.services.DeleteId(id)
	if err != nil {
		msg := services.Response("The Id not found")
		w.WriteHeader(404)
		w.Write(msg)
		return
	}
	msg := services.Response("The post have been deleted")
	w.WriteHeader(200)
	w.Write(msg)
}

//UpdatePost post by Id with decoding request and encoding response
func (h *PostHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPut {
		msg := services.Response("Method Not Allowed")
		w.WriteHeader(405)
		w.Write(msg)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, MaxRequestSize)

	//id, err := strconv.Atoi(r.URL.Query().Get("Id"))
	//msg := r.URL.Query().Get("Message")
	var post model.Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		msg := services.Response("Too big file.")
		w.WriteHeader(401)
		w.Write(msg)
		return
	}
	vars := mux.Vars(r)
	key := vars["id"]
	id, err := strconv.Atoi(key)
	if err != nil {
		msg := services.Response("The Id is not valid")
		w.WriteHeader(406)
		w.Write(msg)
		return
	}
	post.Id = id
	res, err := h.services.UpdateId(&post)
	if err != nil {
		msg := services.Response("Couldn't update requested post.")
		w.WriteHeader(404)
		w.Write(msg)
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(&res)
}

func (h *PostHandler) DownloadPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		msg := services.Response("Method Not Allowed")
		w.WriteHeader(405)
		w.Write(msg)
		return
	}
	res, err := h.services.GetALL()
	if err != nil {
		msg := services.Response("Couldn't find posts.")
		w.WriteHeader(200)
		w.Write(msg)
		return
	}

	err = h.services.Download(*res)
	if err != nil {
		msg := services.Response("The file couldn't be created")
		w.WriteHeader(401)
		w.Write(msg)
	}

	msg := services.Response("The file have been created")
	w.WriteHeader(200)
	w.Write(msg)
}

func (h *PostHandler) UploadPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		msg := services.Response("Method Not Allowed")
		w.WriteHeader(405)
		w.Write(msg)
		return
	}

	err := h.services.Upload()
	if err != nil {
		msg := services.Response("Couldn't upload data from the file")
		w.WriteHeader(401)
		w.Write(msg)
	}

	msg := services.Response("The data from file have been uploaded to the memory Storage")
	w.WriteHeader(200)
	w.Write(msg)
}
