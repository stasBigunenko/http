package services

import (
	"errors"
	"github.com/jszwec/csvutil"
	"src/http/pkg/mocks"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"src/http/pkg/model"
	"src/http/storage"
)

func TestServices_CreateId(t *testing.T) {
	s := new(mocks.Storage)
	id := uuid.New()
	s.On("Create", model.Post{Author: "Stas", Message: "first"}).Return(model.Post{id, "Stas", "first"}, nil)

	s2 := new(mocks.Storage)
	s2.On("Create", model.Post{Author: "", Message: "vfcfggf"}).Return(model.Post{}, errors.New("empty post"))

	s3 := new(mocks.Storage)
	s3.On("Create", model.Post{Author: "dsadsa", Message: ""}).Return(model.Post{}, errors.New("empty message"))

	s4 := new(mocks.Storage)
	s4.On("Create", mock.Anything).Return(model.Post{}, errors.New("err"))

	type args struct {
		storage storage.Storage
	}

	tt := []struct {
		name    string
		args    args
		insert  *model.Post
		want    *model.Post
		wantErr string
	}{
		{
			name: "Create posts Everything OK",
			args: args{
				storage: s,
			},
			insert: &model.Post{
				Author:  "Stas",
				Message: "first",
			},
			want: &model.Post{
				id,
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
			insert: &model.Post{
				Author:  "dsadsa",
				Message: "",
			},
			want:    &model.Post{},
			wantErr: "message is epmty",
		},
		{
			name: "Create posts message",
			args: args{
				storage: s3,
			},
			insert: &model.Post{
				Author:  "",
				Message: "vfcfggf",
			},
			want:    &model.Post{},
			wantErr: "err",
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			serv := NewService(tc.args.storage)
			postNew, err := serv.CreateId(tc.insert)
			if err != nil {
				assert.Error(t, err)
				return
			}

			assert.Equal(t, postNew, tc.want, "Everything cool")
		})
	}
}

func TestServices_GetId(t *testing.T) {
	s := new(mocks.Storage)
	id := uuid.New()
	idStr := id.String()
	s.On("Get", id).Return(model.Post{id, "Stas", "first"}, nil)

	s2 := new(mocks.Storage)
	s2.On("Get", mock.Anything).Return(model.Post{}, errors.New("empty post"))

	type args struct {
		storage storage.Storage
	}

	tt := []struct {
		name    string
		args    args
		insert  string
		want    *model.Post
		wantErr string
	}{
		{
			name: "Create posts Everything OK",
			args: args{
				storage: s,
			},
			insert: idStr,
			want: &model.Post{
				id,
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
			insert:  idStr,
			want:    &model.Post{},
			wantErr: "post with Id %d not found",
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			serv := NewService(tc.args.storage)
			postNew, err := serv.GetId(tc.insert)
			if err != nil {
				assert.Error(t, err, tc.wantErr)
				return
			}

			assert.Equal(t, postNew, tc.want, "Everything cool")
		})
	}
}

func TestServices_Download(t *testing.T) {
	s := new(mocks.Storage)
	id1 := uuid.New()
	id2 := uuid.New()
	x := &[]model.Post{
		{id1, "Stas", "first"},
		{id2, "Stas", "second"},
	}
	s.On("GetAll").Return(*x, nil)

	j, _ := csvutil.Marshal(x)

	type args struct {
		storage storage.Storage
	}

	tt := []struct {
		name     string
		args     args
		expected []byte
	}{
		{
			name: "Download posts Everything OK",
			args: args{
				storage: s,
			},
			expected: j,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			serv := NewService(tc.args.storage)
			res, err := serv.Download()
			if err != nil {
				t.Errorf("something wrong")
				return
			}
			assert.Equal(t, res, j, "Everything cool")
		})
	}
}

func TestServices_DeleteId(t *testing.T) {
	s := new(mocks.Storage)
	id := uuid.New()
	idStr := id.String()
	s.On("Delete", id).Return(nil)

	type args struct {
		storage storage.Storage
	}

	tt := []struct {
		name   string
		args   args
		insert string
	}{
		{
			name: "Delete post Everything OK",
			args: args{
				storage: s,
			},
			insert: idStr,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			serv := NewService(tc.args.storage)
			res := serv.DeleteId(tc.insert)
			assert.Nil(t, res)
		})
	}

}

func TestServices_UpdateId(t *testing.T) {
	s := new(mocks.Storage)
	id := uuid.New()
	s.On("Update", model.Post{id, "Stas", "the first"}).Return(model.Post{id, "Stas", "the first"}, nil)

	type args struct {
		storage storage.Storage
	}

	tt := []struct {
		name     string
		args     args
		insert   model.Post
		expected model.Post
		wantErr  error
	}{
		{
			name: "Delete post Everything OK",
			args: args{
				storage: s,
			},
			insert:   model.Post{id, "Stas", "the first"},
			expected: model.Post{id, "Stas", "the first"},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			serv := NewService(tc.args.storage)
			res, err := serv.UpdateId(&tc.insert)
			if err != nil {
				t.Errorf("error = %v", err.Error())
				return
			}
			assert.Equal(t, &tc.expected, res, "expected model.Post")
		})
	}
}

func TestServices_GetALL(t *testing.T) {
	s := new(mocks.Storage)
	id1 := uuid.New()
	id2 := uuid.New()
	x := []model.Post{
		{id1, "Stas", "first"},
		{id2, "Stas", "second"},
	}
	s.On("GetAll").Return(x, nil)

	type args struct {
		storage storage.Storage
	}

	tt := []struct {
		name     string
		args     args
		expected *[]model.Post
		wantErr  error
	}{
		{
			name: "Delete post Everything OK",
			args: args{
				storage: s,
			},
			expected: &x,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			serv := NewService(tc.args.storage)
			res := serv.GetALL()
			assert.Equal(t, tc.expected, res, "expected model.Post")
		})
	}
}

func TestServices_CreatePost(t *testing.T) {
	s := new(mocks.Storage)
	id := uuid.New()
	s.On("CreateFromFile", model.Post{Id: id, Author: "Stas", Message: "first"}).Return(nil)

	s2 := new(mocks.Storage)
	s2.On("CreateFromFile", model.Post{Author: "", Message: "vfcfggf"}).Return(errors.New("empty post"))

	type args struct {
		storage storage.Storage
	}

	tt := []struct {
		name    string
		args    args
		insert  *model.Post
		want    *model.Post
		wantErr string
	}{
		{
			name: "Create posts Everything OK",
			args: args{
				storage: s,
			},
			insert: &model.Post{
				Id:      id,
				Author:  "Stas",
				Message: "first",
			},
			want: &model.Post{
				id,
				"Stas",
				"first",
			},
		},
		{
			name: "Create posts empty",
			args: args{
				storage: s2,
			},
			insert: &model.Post{
				Author:  "",
				Message: "vfcfggf",
			},
			want: &model.Post{},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			serv := NewService(tc.args.storage)
			err := serv.CreatePost(tc.insert)
			if err != nil {
				assert.Error(t, err)
				return
			}
			assert.Nil(t, err)
		})
	}
}
