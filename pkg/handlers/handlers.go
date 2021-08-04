package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
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

func (h *postHandler) NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/post", h.CreatePost).Methods("POST")
	r.HandleFunc("/post/{Id}", h.GetPost).Methods("GET")
	r.HandleFunc("/post", h.GetAll).Methods("GET")
	r.HandleFunc("/post/{Id}", h.DeletePost).Methods("DELETE")
	r.HandleFunc("/post/{Id}", h.UpdatePost).Methods("PUT")

	return r
}

func New(s *storage.Storage) *postHandler {
	ps := services.NewStore(*s)
	return &postHandler{
		*ps,
	}
}

func (h *postHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var post model.Post
	json.NewDecoder(r.Body).Decode(&post)
	res, _ := h.Services.CreateId(&post)
	json.NewEncoder(w).Encode(&res)
}

func (h *postHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("It is working!!!") //debug not working!
	vars := mux.Vars(r)
	key := vars["Id"]
	idInt, _ := strconv.Atoi(key)
	res, _ := h.Services.GetId(idInt)
	json.NewEncoder(w).Encode(&res)
}
func (h *postHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	res, _ := h.Services.GetALL()
	json.NewEncoder(w).Encode(&res)
}

func (h *postHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)
	idInt, _ := strconv.Atoi(id["id"])
	res, _ := h.Services.DeleteId(idInt)
	json.NewEncoder(w).Encode(&res)
}

func (h *postHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	var post model.Post
	json.NewDecoder(r.Body).Decode(&post)
	res, _ := h.Services.UpdateId(&post)
	json.NewEncoder(w).Encode(&res)
}
