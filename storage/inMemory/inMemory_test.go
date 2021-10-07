package inMemory

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"src/http/pkg/model"
	"testing"
)

func TestStorage_Create(t *testing.T) {
	db := New()
	m := model.Post{
		Author:  "stas",
		Message: "the first",
	}

	tests := []struct {
		name  string
		param model.Post
	}{
		{
			name:  "everything ok",
			param: m,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			val, _ := db.Create(tc.param)
			require.NotNil(t, val, "shouldn't be error")
		})
	}
}

func TestStorage_Get(t *testing.T) {
	db := New()
	id := uuid.New()
	id1 := uuid.New()
	m := model.Post{Id: id, Author: "stas", Message: "the first"}
	db.storage[id] = m
	tests := []struct {
		name    string
		param   uuid.UUID
		want    model.Post
		wantErr string
	}{
		{
			name:  "Everything ok",
			param: id,
			want:  m,
		},
		{
			name:    "Invalid id",
			param:   id1,
			want:    model.Post{},
			wantErr: "post with Id %d not found",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := db.Get(tc.param)
			if err != nil && error.Error(err) != tc.wantErr {
				t.Errorf("error = %v, wantErr %v", err.Error(), tc.wantErr)
				return
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestStorage_GetAll(t *testing.T) {
	db := New()
	id1 := uuid.New()
	id2 := uuid.New()
	m1 := model.Post{Id: id1, Author: "stas", Message: "the first"}
	db.storage[id1] = m1
	m2 := model.Post{Id: id2, Author: "andrew", Message: "the last"}
	db.storage[id2] = m2
	m := []model.Post{m1, m2}

	tests := []struct {
		name string
		want []model.Post
	}{
		{
			name: "everything ok",
			want: m,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res := db.GetAll()
			require.NotNil(t, res, "shouldn't be error")
		})
	}
}

func TestStorage_Update(t *testing.T) {
	db := New()
	id1 := uuid.New()
	m1 := model.Post{Id: uuid.New(), Author: "gfd", Message: "gfd"}
	m2 := model.Post{Id: id1, Author: "andrew", Message: "the last"}
	db.storage[id1] = m1

	tests := []struct {
		name    string
		param   model.Post
		want    model.Post
		wantErr string
	}{
		{
			name:  "everything ok",
			param: m2,
			want:  m2,
		},
		{
			name:    "not found id",
			param:   m1,
			want:    model.Post{},
			wantErr: "post cann't be updated. The post doesn't exist",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := db.Update(tc.param)
			if err != nil && error.Error(err) != tc.wantErr {
				t.Errorf("error = %v, wantErr %v", err.Error(), tc.wantErr)
				return
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestStorage_Delete(t *testing.T) {
	db := New()
	id := uuid.New()
	id1 := uuid.New()
	m := model.Post{Id: id, Author: "stas", Message: "the first"}
	db.storage[id] = m
	tests := []struct {
		name    string
		param   uuid.UUID
		wantErr string
	}{
		{
			name:    "Everything ok",
			param:   id,
			wantErr: "",
		},
		{
			name:    "Invalid id",
			param:   id1,
			wantErr: "post can't be deleted - Id not found",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := db.Delete(tc.param)
			if err != nil && error.Error(err) != tc.wantErr {
				t.Errorf("error = %v, wantErr %v", err.Error(), tc.wantErr)
				return
			}
		})
	}
}
