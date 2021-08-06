package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"src/http/pkg/model"
	"src/http/pkg/services"
	"src/http/storage"
	"strconv"
)

//Handlers with the CRUD functions

type postHandler struct {
	Services services.Store
}

// Simple middleware function which write log in the terminal requested Method and URI
func simpleLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method)
		log.Println(r.RequestURI)
		fmt.Println("-------------")
		next.ServeHTTP(w, r)
	})
}

func (h *postHandler) NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/post/", h.CreatePost).Methods("POST")
	r.HandleFunc("/post/{id}", h.GetPost).Methods("GET")
	r.HandleFunc("/posts", h.GetAll).Methods("GET")
	r.HandleFunc("/post/{id}", h.DeletePost).Methods("DELETE")
	r.HandleFunc("/post/{id}", h.UpdatePost).Methods("PUT")
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
		w.WriteHeader(405)
		w.Write([]byte("Method Not Allowed"))
		return
	}
	var post model.Post
	json.NewDecoder(r.Body).Decode(&post)
	res, err := h.Services.CreateId(&post)
	if err != nil {
		w.Write([]byte("could not create empty post"))
		return
	}
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(&res)
}

//GetPost Get post by Id with decoding request and encoding response
func (h *postHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		w.WriteHeader(405)
		w.Write([]byte("Method Not Allowed"))
		return
	}
	//id, err := strconv.Atoi(r.URL.Query().Get("id"))
	vars := mux.Vars(r)
	key := vars["id"]
	id, err := strconv.Atoi(key)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("The Id is not valid"))
		return
	}
	res, err := h.Services.GetId(id)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Couldn't get requested post. This Id = %d doesn't not exist.\n", id)))
		return
	}
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(&res)
}

//GetAll posts by Id with decoding request and encoding response
func (h *postHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		w.WriteHeader(405)
		w.Write([]byte("Method Not Allowed"))
		return
	}
	res, err := h.Services.GetALL()
	if err != nil {
		w.Write([]byte("There is no posts in the memory."))
		return
	}
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(&res)
}

//DeletePost post by Id with decoding request and encoding response
func (h *postHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodDelete {
		w.WriteHeader(405)
		w.Write([]byte("Method Not Allowed"))
		return
	}
	//id, err := strconv.Atoi(r.URL.Query().Get("Id"))
	vars := mux.Vars(r)
	key := vars["id"]
	id, err := strconv.Atoi(key)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("The Id is not valid"))
		return
	}
	res, err := h.Services.DeleteId(id)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Could not delete post with this id - %d\n.", id)))
		return
	}
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(&res)
}

//UpdatePost post by Id with decoding request and encoding response
func (h *postHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPut {
		w.WriteHeader(405)
		w.Write([]byte("Method Not Allowed"))
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
		w.WriteHeader(400)
		w.Write([]byte("The Id is not valid"))
		return
	}
	post.Id = id
	res, err := h.Services.UpdateId(&post)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Couldn't update requested post. Post with id=%d doesn't exist.\n", id)))
		return
	}
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(&res)
}
