package services

import (
	"errors"
	"github.com/jszwec/csvutil"
	"github.com/stretchr/testify/assert"

	"src/http/pkg/model"
	"src/http/storage"
	"src/http/storage/mocks"
	"testing"
)

func TestCreateId(t *testing.T) {
	s := new(mocks.Storage)
	s.On("Create", model.Post{Author: "Stas", Message: "first"}).Return(model.Post{1, "Stas", "first"}, nil)

	s2 := new(mocks.Storage)
	s2.On("Create", model.Post{Author: "", Message: "vfcfggf"}).Return(model.Post{}, errors.New("empty post"))

	type args struct {
		storage storage.Storage
	}

	tt := []struct {
		name    string
		args    args
		insert  model.Post
		want    model.Post
		wantErr string
	}{
		{
			name: "Create posts Everything OK",
			args: args{
				storage: s,
			},
			insert: model.Post{
				Author:  "Stas",
				Message: "first",
			},
			want: model.Post{
				1,
				"Stas",
				"first",
			},
			wantErr: "",
		},
		{
			name: "Create posts empty",
			args: args{
				storage: s2,
			},
			insert: model.Post{
				Author:  "",
				Message: "vfcfggf",
			},
			want:    model.Post{},
			wantErr: "author is epmty",
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			postNew, err := tc.args.storage.Create(tc.insert)
			if err != nil {
				assert.Error(t, err)
				return
			}

			assert.Equal(t, postNew, tc.want, "Everything cool")
		})
	}
}

func TestGeteId(t *testing.T) {
	s := new(mocks.Storage)
	s.On("Get", 1).Return(model.Post{1, "Stas", "first"}, nil)

	s2 := new(mocks.Storage)
	s2.On("Get", 100).Return(model.Post{}, errors.New("empty post"))

	type args struct {
		storage storage.Storage
	}

	tt := []struct {
		name    string
		args    args
		insert  int
		want    model.Post
		wantErr string
	}{
		{
			name: "Create posts Everything OK",
			args: args{
				storage: s,
			},
			insert: 1,
			want: model.Post{
				1,
				"Stas",
				"first",
			},
			wantErr: "",
		},
		{
			name: "Create posts empty",
			args: args{
				storage: s2,
			},
			insert:  100,
			want:    model.Post{},
			wantErr: "post with Id %d not found",
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			postNew, err := tc.args.storage.Get(tc.insert)
			if err != nil {
				assert.Error(t, err, tc.wantErr)
				return
			}

			assert.Equal(t, postNew, tc.want, "Everything cool")
		})
	}
}

func TestDownload(t *testing.T) {
	s := new(mocks.Storage)
	x := &[]model.Post{
		{1, "Stas", "first"},
		{2, "Stas", "second"},
	}
	s.On("GetAll").Return(*x, nil)

	type args struct {
		storage storage.Storage
	}

	tt := []struct {
		name string
		args args
	}{
		{
			name: "Download posts Everything OK",
			args: args{
				storage: s,
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			allPosts := tc.args.storage.GetAll()

			_, err := csvutil.Marshal(allPosts)
			if err != nil {
				t.Fatalf("some error %v", err)
			}

			//assert.Equal(t, ap, , "Everything cool")
		})
	}
}
