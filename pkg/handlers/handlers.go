package handlers

import (
	"context"
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"src/http/pkg/model"
	"src/http/pkg/services"
)

//Handlers with the CRUD functions and Middleware

type PostHandler struct {
	service services.ServicesInterface
}

const MaxRequestSize = 2 * 1024

func NewHandler(service services.ServicesInterface) *PostHandler {
	return &PostHandler{
		service: service,
	}
}

//Middleware: timeout to the handler process
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

func (h *PostHandler) Routes(sub *mux.Router) *mux.Router {
	sub.Path("/metrics").Handler(promhttp.Handler())

	sub.HandleFunc("/login", processTimeout(h.Login, 5*time.Second)).Methods("GET")

	sub.HandleFunc("/", processTimeout(h.GetAll, 5*time.Second)).Methods("GET")
	sub.HandleFunc("/download", processTimeout(h.DownloadPost, 5*time.Second)).Methods("GET")
	sub.HandleFunc("/upload", processTimeout(h.UploadPost, 5*time.Second)).Methods("POST")
	sub.HandleFunc("/create", processTimeout(h.CreatePost, 5*time.Second)).Methods("POST")
	sub.HandleFunc("/{id}", processTimeout(h.GetPost, 5*time.Second)).Methods("GET")
	sub.HandleFunc("/{id}", processTimeout(h.DeletePost, 5*time.Second)).Methods("DELETE")
	sub.HandleFunc("/{id}", processTimeout(h.UpdatePost, 5*time.Second)).Methods("PUT")

	sub.Use(h.prometheusMiddleware)
	sub.Use(h.VerifyUser)

	return sub
}

//CreatePost Create post with decoding request and encoding response
func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		msg := services.Response("Method Not Allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(msg)
		return
	}

	ctx := r.Context()
	cl, ok := ctx.Value("claims").(model.Claims)
	if !ok {
		log.Println("cannot identify user name and id")
	}
	log.Printf("user with id '%s', name '%s' send command CreatePost", cl.Id, cl.Name)

	//receive requests only with size limit 2024 bytes
	r.Body = http.MaxBytesReader(w, r.Body, MaxRequestSize)

	var post model.Post

	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		msg := services.Response("Too big request.")
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write(msg)
		return
	}

	createChPost(post)

	res, err := h.service.CreateId(&post)
	if err != nil {
		msg := services.Response("internal problem")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(msg)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&res)
}

//GetPost Get post by Id with decoding request and encoding response
func (h *PostHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		msg := services.Response("Method Not Allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(msg)
		return
	}

	ctx := r.Context()
	cl, ok := ctx.Value("claims").(model.Claims)
	if !ok {
		log.Println("cannot identify user name")
	}
	log.Printf("user with id '%s', name '%s' send command GetPost", cl.Id, cl.Name)

	//id, err := strconv.Atoi(r.URL.Query().Get("id"))
	vars := mux.Vars(r)
	key := vars["id"]

	res, err := h.service.GetId(key)
	if err != nil {
		msg := services.Response("This id doesn't exist")
		w.WriteHeader(http.StatusNotFound)
		w.Write(msg)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&res)
}

//GetAll posts by Id with decoding request and encoding response
func (h *PostHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		msg := services.Response("Method Not Allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(msg)
		return
	}

	ctx := r.Context()
	cl, ok := ctx.Value("claims").(model.Claims)
	if !ok {
		log.Println("cannot identify user name")
	}
	log.Printf("user with id '%s', name '%s' send command GetAll", cl.Id, cl.Name)

	res := h.service.GetALL()

	if len(*res) == 0 {
		msg := services.Response("There is no posts in the memory.")
		w.WriteHeader(http.StatusOK)
		w.Write(msg)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&res)
}

//DeletePost post by Id with decoding request and encoding response
func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodDelete {
		msg := services.Response("Method Not Allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(msg)
		return
	}

	ctx := r.Context()
	cl, ok := ctx.Value("claims").(model.Claims)
	if !ok {
		log.Println("cannot identify user name")
	}
	log.Printf("user with id '%s', name '%s' send command DeletePost", cl.Id, cl.Name)

	//id, err := strconv.Atoi(r.URL.Query().Get("Id"))
	vars := mux.Vars(r)
	key := vars["id"]

	deleteChPost(key)

	err := h.service.DeleteId(key)
	if err != nil {
		msg := services.Response("This id doesn't exist")
		w.WriteHeader(http.StatusNotFound)
		w.Write(msg)
		return
	}

	msg := services.Response("The post have been deleted")
	w.WriteHeader(http.StatusOK)
	w.Write(msg)
}

//UpdatePost post by Id with decoding request and encoding response
func (h *PostHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPut {
		msg := services.Response("Method Not Allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(msg)
		return
	}

	ctx := r.Context()
	cl, ok := ctx.Value("claims").(model.Claims)
	if !ok {
		log.Println("cannot identify user name")
	}
	log.Printf("user with id '%s', name '%s' send command UpdatePost", cl.Id, cl.Name)

	r.Body = http.MaxBytesReader(w, r.Body, MaxRequestSize)

	//id, err := strconv.Atoi(r.URL.Query().Get("Id"))
	//msg := r.URL.Query().Get("Message")
	var post model.Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		msg := services.Response("Too big file.")
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write(msg)
		return
	}

	vars := mux.Vars(r)
	key := vars["id"]

	post.Id, err = uuid.Parse(key)
	if err != nil {
		msg := services.Response("Couldn't parse id.")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(msg)
		return
	}

	createChPost(post)

	res, err := h.service.UpdateId(&post)
	if err != nil {
		msg := services.Response("Couldn't update requested post.")
		w.WriteHeader(http.StatusNotFound)
		w.Write(msg)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&res)
}

func (h *PostHandler) DownloadPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		msg := services.Response("Method Not Allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(msg)
		return
	}

	ctx := r.Context()
	cl, ok := ctx.Value("claims").(model.Claims)
	if !ok {
		log.Println("cannot identify user name")
	}
	log.Printf("user with id '%s', name '%s' send command DownloadPost", cl.Id, cl.Name)

	res, err := h.service.Download()
	if err != nil {
		msg := services.Response("The file couldn't be created")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(msg)
		return
	}

	t := time.Now()
	time := (t.Format("2006_01_02_15-04"))

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=result_"+time+".csv")
	//msg := Service.Response("The file downloaded to the memory")
	w.Write(res)
}

func (h *PostHandler) UploadPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		msg := services.Response("Method Not Allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(msg)
		return
	}

	ctx := r.Context()
	cl, ok := ctx.Value("claims").(model.Claims)
	if !ok {
		log.Println("cannot identify user name")
	}
	log.Printf("user with id '%s', name '%s' send command UploadPost", cl.Id, cl.Name)

	err := r.ParseMultipartForm(50) // limit input length!
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	r.Body = http.MaxBytesReader(w, r.Body, MaxRequestSize)

	// upload file
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fileSize := fileHeader.Size
	if fileSize > MaxRequestSize {
		msg := services.Response("Too big request.")
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write(msg)
		return
	}

	err = h.service.Upload(file)
	if err != nil {
		msg := services.Response("Couldn't upload data from the file")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(msg)
		return
	}

	msg := services.Response("The data from file have been uploaded to the memory Storage")
	w.WriteHeader(http.StatusOK)
	w.Write(msg)
}

func (h *PostHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		msg := services.Response("Method Not Allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(msg)
		return
	}

	var user struct {
		Name string
	}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		msg := services.Response("Invalid user name")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(msg)
		return
	}

	token, err := h.generateTokenStringForUser(user.Name)
	if err != nil {
		msg := services.Response("Invalid name")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(msg)
		return
	}

	w.Header().Set("token", token)
	w.Header().Set("Access-Control-Expose-Headers", "token")
	w.WriteHeader(http.StatusCreated)
}
