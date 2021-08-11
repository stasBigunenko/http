package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"src/http/pkg/model"
	"src/http/pkg/services"
	"src/http/storage"
	"strconv"
	"time"
)

//Handlers with the CRUD functions

type postHandler struct {
	Services services.Store
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

func (h *postHandler) NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/post/", processTimeout(h.CreatePost, 5*time.Second)).Methods("POST")
	r.HandleFunc("/post/{id}", processTimeout(h.GetPost, 5*time.Second)).Methods("GET")
	r.HandleFunc("/posts", processTimeout(h.GetAll, 5*time.Second)).Methods("GET")
	r.HandleFunc("/post/{id}", processTimeout(h.DeletePost, 5*time.Second)).Methods("DELETE")
	r.HandleFunc("/post/{id}", processTimeout(h.UpdatePost, 5*time.Second)).Methods("PUT")
	r.HandleFunc("/post/upload", processTimeout(h.UploadPost, 5*time.Second)).Methods("POST")
	r.HandleFunc("/post/download", processTimeout(h.DownloadPost, 5*time.Second)).Methods("POST")
	r.Use(simpleLog)
	return r
}

func New(s *storage.Storage) *postHandler {
	ps := services.NewStore(*s)
	return &postHandler{
		*ps,
	}
}

//CreatePost Create post with decoding request and encoding response
func (h *postHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		msg := services.Response("Method Not Allowed")
		w.WriteHeader(405)
		w.Write(msg)
		return
	}
	var post model.Post
	json.NewDecoder(r.Body).Decode(&post)
	res, err := h.Services.CreateId(&post)
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
func (h *postHandler) GetPost(w http.ResponseWriter, r *http.Request) {
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
	res, err := h.Services.GetId(id)
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
func (h *postHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		msg := services.Response("Method Not Allowed")
		w.WriteHeader(405)
		w.Write(msg)
		return
	}
	res, err := h.Services.GetALL()
	if err != nil {
		msg := services.Response("Bad request")
		w.WriteHeader(404)
		w.Write(msg)
		return
	}

	if len(*res) == 0 {
		msg := services.Response("There is no post in the memory.")
		w.WriteHeader(200)
		w.Write(msg)
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(&res)
}

//DeletePost post by Id with decoding request and encoding response
func (h *postHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
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
	err = h.Services.DeleteId(id)
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
func (h *postHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPut {
		msg := services.Response("Method Not Allowed")
		w.WriteHeader(405)
		w.Write(msg)
		return
	}
	//id, err := strconv.Atoi(r.URL.Query().Get("Id"))
	//msg := r.URL.Query().Get("Message")
	var post model.Post
	json.NewDecoder(r.Body).Decode(&post)
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
	res, err := h.Services.UpdateId(&post)
	if err != nil {
		msg := services.Response("Couldn't update requested post.")
		w.WriteHeader(404)
		w.Write(msg)
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(&res)
}

func (h *postHandler) DownloadPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		msg := services.Response("Method Not Allowed")
		w.WriteHeader(405)
		w.Write(msg)
		return
	}
	res, err := h.Services.GetALL()
	if err != nil {
		msg := services.Response("Couldn't find posts.")
		w.WriteHeader(200)
		w.Write(msg)
		return
	}

	err = h.Services.Download(*res)
	if err != nil {
		msg := services.Response("The file couldn't be created")
		w.WriteHeader(401)
		w.Write(msg)
	}

	msg := services.Response("The file have been created")
	w.WriteHeader(200)
	w.Write(msg)
}

func (h *postHandler) UploadPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		msg := services.Response("Method Not Allowed")
		w.WriteHeader(405)
		w.Write(msg)
		return
	}

	err := h.Services.Upload()
	if err != nil {
		msg := services.Response("Couldn't upload data from the file")
		w.WriteHeader(401)
		w.Write(msg)
	}

	msg := services.Response("The data from file have been uploaded to the memory Storage")
	w.WriteHeader(200)
	w.Write(msg)
}
