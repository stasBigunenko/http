package handlers

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/jszwec/csvutil"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"src/http/pkg/model"
	"src/http/pkg/services"
	"src/http/pkg/services/mocks"
	"strings"
	"testing"
)

func TestCreatePost(t *testing.T) {
	s := new(mocks.ServicesInterface)
	s.On("CreateId", &model.Post{Author: "Stas", Message: "first"}).Return(&model.Post{0, "Stas", "first"}, nil)

	s2 := new(mocks.ServicesInterface)
	s2.On("CreateId", &model.Post{Author: "", Message: ""}).Return(nil, errors.New("author is empty"))

	type args struct {
		service services.ServicesInterface
	}

	tt := []struct {
		name   string
		args   args
		buf    string
		method string
		url    string
		status int
	}{
		{
			name: "POST post Everything OK",
			args: args{
				service: s,
			},
			buf:    `{"Author":"Stas","Message":"first"}`,
			method: "POST",
			url:    "localhost:8080/posts/create",
			status: http.StatusCreated,
		},
		{
			name: "POST post err: author and/or message is empty",
			args: args{
				service: s2,
			},
			buf:    `{"Author":"","Message":""}`,
			method: "POST",
			url:    "localhost:8080/posts/create",
			status: http.StatusBadRequest,
		},
		{
			name: "POST post err: wrong method",
			args: args{
				service: s2,
			},
			buf:    `{"Author":"","Message":""}`,
			method: "GET",
			url:    "localhost:8080/posts/create",
			status: http.StatusMethodNotAllowed,
		},
		{
			name: "POST post err: couldn't decode Body",
			args: args{
				service: s2,
			},
			buf:    `78787{"Author2":"","Message3":""}`,
			method: "POST",
			url:    "localhost:8080/posts/create",
			status: http.StatusNotAcceptable,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ph := NewHandler(tc.args.service)
			req, err := http.NewRequest(tc.method, tc.url, strings.NewReader(tc.buf))
			if err != nil {
				t.Fatalf("something goes wrong %v", err)
			}
			rr := httptest.NewRecorder()

			ph.CreatePost(rr, req)

			res := rr.Result()
			defer res.Body.Close()

			var p model.Post

			err = json.NewDecoder(res.Body).Decode(&p)
			if err != nil {
				t.Fatalf("coudln't read json %v", err)
			}

			assert.Equal(t, res.StatusCode, tc.status, "Everything cool")
		})
	}
}

func TestGetPost(t *testing.T) {
	s := new(mocks.ServicesInterface)
	s.On("GetId", 1).Return(&model.Post{1, "Stas", "first"}, nil)

	s2 := new(mocks.ServicesInterface)
	s2.On("GetId", 25).Return(nil, errors.New("Not found"))

	type args struct {
		service services.ServicesInterface
	}

	tt := []struct {
		name   string
		args   args
		vars   map[string]string
		buf    string
		method string
		url    string
		status int
	}{
		{
			name: "GET post Everything OK",
			args: args{
				service: s,
			},
			buf: `1`,
			vars: map[string]string{
				"id": "1",
			},
			method: "GET",
			url:    "localhost:8080/posts/1",
			status: http.StatusOK,
		},
		{
			name: "GET post wrong method",
			args: args{
				service: s,
			},
			buf: `1`,
			vars: map[string]string{
				"id": "1",
			},
			method: "POST",
			url:    "localhost:8080/posts/1",
			status: http.StatusMethodNotAllowed,
		},
		{
			name: "GET post wrong url",
			args: args{
				service: s,
			},
			buf: `100`,
			vars: map[string]string{
				"id": "1v",
			},
			method: "GET",
			url:    "localhost:8080/posts/1v",
			status: http.StatusBadRequest,
		},
		{
			name: "GET post post not found",
			args: args{
				service: s2,
			},
			buf: ``,
			vars: map[string]string{
				"id": "25",
			},
			method: "GET",
			url:    "localhost:8080/posts/25",
			status: http.StatusNotFound,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			ph := NewHandler(tc.args.service)
			req, err := http.NewRequest(tc.method, tc.url, nil)
			if err != nil {
				t.Fatalf("something goes wrong %v", err)
			}
			rr := httptest.NewRecorder()

			vars := tc.vars
			req = mux.SetURLVars(req, vars)

			ph.GetPost(rr, req)
			res := rr.Result()

			assert.Equal(t, res.StatusCode, tc.status, "Everything cool")
		})
	}
}

func TestGetAll(t *testing.T) {
	s := new(mocks.ServicesInterface)
	x := &[]model.Post{
		{1, "Stas", "first"},
		{2, "Stas", "second"},
	}
	s.On("GetALL").Return(x, nil)

	s2 := new(mocks.ServicesInterface)
	x2 := &[]model.Post{}
	s2.On("GetALL").Return(x2, nil)

	type args struct {
		service services.ServicesInterface
	}

	tt := []struct {
		name   string
		args   args
		buf    string
		len    int
		method string
		url    string
		status int
	}{
		{
			name: "GetAll post Everything OK",
			args: args{
				service: s,
			},
			buf:    `1`,
			method: "GET",
			url:    "localhost:8080/posts/",
			status: http.StatusOK,
		},
		{
			name: "GetAll post err: wrong method",
			args: args{
				service: s,
			},
			buf:    `1`,
			method: "POST",
			url:    "localhost:8080/posts/",
			status: http.StatusMethodNotAllowed,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			ph := NewHandler(tc.args.service)
			req, err := http.NewRequest(tc.method, tc.url, nil)
			if err != nil {
				t.Fatalf("something goes wrong %v", err)
			}
			rr := httptest.NewRecorder()

			ph.GetAll(rr, req)

			res := rr.Result()

			assert.Equal(t, res.StatusCode, tc.status, "Everything cool")
		})
	}
}

func TestDeletePost(t *testing.T) {
	s := new(mocks.ServicesInterface)
	s.On("DeleteId", 1).Return(nil)

	s2 := new(mocks.ServicesInterface)
	s2.On("DeleteId", 100).Return(errors.New("not found"))

	type args struct {
		service services.ServicesInterface
	}

	tt := []struct {
		name   string
		args   args
		buf    string
		method string
		vars   map[string]string
		url    string
		status int
	}{
		{
			name: "DELETE post Everything OK",
			args: args{
				service: s,
			},
			buf: `1`,
			vars: map[string]string{
				"id": "1",
			},
			method: "DELETE",
			url:    "localhost:8080/posts/1",
			status: http.StatusOK,
		},
		{
			name: "DELETE post err: wrong method",
			args: args{
				service: s,
			},
			buf: `1`,
			vars: map[string]string{
				"id": "1",
			},
			method: "GET",
			url:    "localhost:8080/posts/1",
			status: http.StatusMethodNotAllowed,
		},
		{
			name: "DELETE post err: wrong url",
			args: args{
				service: s,
			},
			buf: `1`,
			vars: map[string]string{
				"id": "1s",
			},
			method: "DELETE",
			url:    "localhost:8080/posts/1s",
			status: http.StatusBadRequest,
		},
		{
			name: "DELETE post err:post not found",
			args: args{
				service: s2,
			},
			buf: `100`,
			vars: map[string]string{
				"id": "100",
			},
			method: "DELETE",
			url:    "localhost:8080/posts/100",
			status: http.StatusNotFound,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			ph := NewHandler(tc.args.service)
			req, err := http.NewRequest(tc.method, tc.url, nil)
			if err != nil {
				t.Fatalf("something goes wrong %v", err)
			}
			rr := httptest.NewRecorder()

			vars := tc.vars
			req = mux.SetURLVars(req, vars)

			ph.DeletePost(rr, req)

			res := rr.Result()

			assert.Equal(t, res.StatusCode, tc.status, "Everything cool")
		})
	}
}

func TestUpdatePost(t *testing.T) {
	s := new(mocks.ServicesInterface)
	s.On("UpdateId", &model.Post{Id: 1, Author: "Stas", Message: "first"}).Return(&model.Post{1, "Alexey", "wrong"}, nil)

	s2 := new(mocks.ServicesInterface)
	s2.On("UpdateId", &model.Post{Id: 1000, Author: "Stas", Message: "first"}).Return(nil, errors.New("post cann't be updated. The post doesn't exist"))

	type args struct {
		service services.ServicesInterface
	}

	tt := []struct {
		name   string
		args   args
		buf    string
		method string
		vars   map[string]string
		url    string
		status int
	}{
		{
			name: "UPDATE post Everything OK",
			args: args{
				service: s,
			},
			buf:    `{"Author":"Stas","Message":"first"}`,
			method: "PUT",
			url:    "localhost:8080/posts/1",
			status: http.StatusOK,
			vars: map[string]string{
				"id": "1",
			},
		},
		{
			name: "UPDATE post err: wrong method",
			args: args{
				service: s,
			},
			buf:    `{"Author":"Stas","Message":"first"}`,
			method: "GET",
			url:    "localhost:8080/posts/1",
			status: http.StatusMethodNotAllowed,
			vars: map[string]string{
				"id": "1",
			},
		},
		{
			name: "UPDATE post err: wrong url request",
			args: args{
				service: s,
			},
			buf:    `{"Author":"Stas","Message":"first"}`,
			method: "PUT",
			url:    "localhost:8080/posts/1c",
			status: http.StatusBadRequest,
			vars: map[string]string{
				"id": "1c",
			},
		},
		{
			name: "UPDATE post err: not found",
			args: args{
				service: s2,
			},
			buf:    `{"Author":"Stas","Message":"first"}`,
			method: "PUT",
			url:    "localhost:8080/posts/1000",
			status: http.StatusNotFound,
			vars: map[string]string{
				"id": "1000",
			},
		},
		{
			name: "UPDATE post err: wrong body",
			args: args{
				service: s,
			},
			buf:    `dsadsa{"Author":"Stas","Message":"first"}`,
			method: "PUT",
			url:    "localhost:8080/posts/1",
			status: http.StatusNotAcceptable,
			vars: map[string]string{
				"id": "1",
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			ph := NewHandler(tc.args.service)
			req, err := http.NewRequest(tc.method, tc.url, strings.NewReader(tc.buf))
			if err != nil {
				t.Fatalf("something goes wrong %v", err)
			}

			vars := tc.vars
			req = mux.SetURLVars(req, vars)

			rr := httptest.NewRecorder()
			ph.UpdatePost(rr, req)

			res := rr.Result()
			defer res.Body.Close()

			assert.Equal(t, res.StatusCode, tc.status, "Everything cool")
		})
	}
}

func TestDownloadPost(t *testing.T) {
	s := new(mocks.ServicesInterface)
	x := &[]model.Post{
		{1, "Stas", "first"},
		{2, "Stas", "second"},
	}
	ap, _ := csvutil.Marshal(x)
	s.On("Download").Return(ap, nil)

	s2 := new(mocks.ServicesInterface)
	s2.On("Download").Return(nil, errors.New("couldn't be created"))

	type args struct {
		service services.ServicesInterface
	}

	tt := []struct {
		name   string
		args   args
		buf    string
		method string
		url    string
		status int
	}{
		{
			name: "Dowonload post Everything OK",
			args: args{
				service: s,
			},
			buf:    "",
			method: "GET",
			url:    "localhost:8080/posts/download",
			status: http.StatusOK,
		},
		{
			name: "Dowonload post err: wrong method",
			args: args{
				service: s,
			},
			buf:    "",
			method: "DELETE",
			url:    "localhost:8080/posts/download",
			status: http.StatusMethodNotAllowed,
		},
		{
			name: "Dowonload post err: couldn't created",
			args: args{
				service: s2,
			},
			buf:    "",
			method: "GET",
			url:    "localhost:8080/posts/download",
			status: http.StatusUnauthorized,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			ph := NewHandler(tc.args.service)
			req, err := http.NewRequest(tc.method, tc.url, nil)
			if err != nil {
				t.Fatalf("something goes wrong %v", err)
			}
			rr := httptest.NewRecorder()

			ph.DownloadPost(rr, req)

			res := rr.Result()

			assert.Equal(t, res.StatusCode, tc.status, "Everything cool")
		})
	}
}

func TestRoutes(t *testing.T) {
	s := new(mocks.ServicesInterface)
	s.On("CreateId", &model.Post{Author: "Stas", Message: "first"}).Return(&model.Post{0, "Stas", "first"}, nil)

	s2 := new(mocks.ServicesInterface)
	x := &[]model.Post{
		{1, "Stas", "first"},
		{2, "Stas", "second"},
	}
	ap, _ := csvutil.Marshal(x)
	s2.On("Download").Return(ap, nil)

	type args struct {
		service services.ServicesInterface
	}

	tt := []struct {
		name   string
		args   args
		buf    string
		method string
		url    string
		status int
	}{
		{
			name: "POST post Everything OK",
			args: args{
				service: s,
			},
			buf:    `{"Author":"Stas","Message":"first"}`,
			method: "POST",
			url:    "/create",
			status: http.StatusCreated,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ph := NewHandler(tc.args.service)
			r := mux.NewRouter()
			ph.Routes(r)
			rr := httptest.NewRecorder()
			req, err := http.NewRequest(tc.method, tc.url, strings.NewReader(tc.buf))
			if err != nil {
				t.Fatalf("smth goes wron %v", err)
			}
			r.ServeHTTP(rr, req)

			assert.Equal(t, tc.status, rr.Code, "Everything cool")
		})
	}
}
