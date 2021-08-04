package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"src/http/pkg/model"
	"src/http/pkg/services"
	"src/http/storage"
	"strconv"
)

//Handlers with the CRUD functions
//decoding requests from []byte and coding response in []byte
type postHandler struct {
	Services services.Store
}

func (h *postHandler) NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/post", h.CreatePost).Methods("POST")
	r.HandleFunc("/post", h.GetPost).Methods("GET")
	r.HandleFunc("/posts", h.GetAll).Methods("GET")
	r.HandleFunc("/post", h.DeletePost).Methods("DELETE")
	r.HandleFunc("/post", h.UpdatePost).Methods("PUT")

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
	res, err := h.Services.CreateId(&post)
	if err != nil {
		w.Write([]byte("could not create empty post"))
		return
	}
	json.NewEncoder(w).Encode(&res)

}

func (h *postHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("Id"))
	res, err := h.Services.GetId(id)
	if err != nil {
		w.Write([]byte("could not get post"))
		return
	}
	json.NewEncoder(w).Encode(&res)
}
func (h *postHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	res, err := h.Services.GetALL()
	if err != nil {
		w.Write([]byte("could not get all posts"))
		return
	}
	json.NewEncoder(w).Encode(&res)
}

func (h *postHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("Id"))
	res, err := h.Services.DeleteId(id)
	if err != nil {
		w.Write([]byte("could not delete post"))
		return
	}
	json.NewEncoder(w).Encode(&res)
}

func (h *postHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	var post model.Post
	id, err := strconv.Atoi(r.URL.Query().Get("Id"))
	json.NewDecoder(r.Body).Decode(&post)
	res, err := h.Services.UpdateId(id, &post)
	if err != nil {
		w.Write([]byte("could not update post"))
		return
	}
	json.NewEncoder(w).Encode(&res)
}
