package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"src/http/pkg/model"
	"testing"

	"src/http/pkg/services"
	"src/http/storage"
)

var p = model.Post{}

func TestPostHandler(t *testing.T) {
	type fields struct {
		s storage.Storage
	}

	type args struct {
		url    string
		method string
		header string
		body   []byte
	}

	type resp struct {
		code int
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   resp
	}{
		{
			name: "POST /post OK",
			fields: fields{
				s: storage.MockStorage{
					MockCreate: func(_ model.Post) (model.Post, error) {
						return p, nil
					},
				},
			},
			args: args{
				url:    "/post",
				method: "POST",
				body:   []byte(`{"Id": "1","Author": "Stas","Message": "The first"}`),
			},
			want: resp{code: http.StatusCreated},
		},
		{
			name:   "POST /post no name",
			fields: fields{},
			args: args{
				url:    "/post",
				method: "POST",
				body:   []byte(`{"wrong":"bad"}`),
			},
			want: resp{code: http.StatusBadRequest},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := services.NewStore(tt.fields.s)

			var s = &PostHandler{
				services: srv,
			}

			hs := httptest.NewServer(s.Routes())
			defer hs.Close()

			cl := hs.Client()
			req, _ := http.NewRequest(tt.args.method, hs.URL+tt.args.url, bytes.NewReader(tt.args.body))

			r, err := cl.Do(req)

			if err != nil || r.StatusCode != tt.want.code {
				if err != nil {
					t.Errorf("error: %s", err)
				} else {
					t.Errorf("%s %s = %v, want %v", tt.args.method, tt.args.url, r.StatusCode, tt.want.code)
				}
			}
		})
	}
}
