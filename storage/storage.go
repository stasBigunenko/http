package storage

import (
	"src/http/pkg/model"
)

//Storage Interface
type Storage interface {
	Create(model.Post) (model.Post, error)
	Get(int) (model.Post, error)
	GetAll() []model.Post
	Update(int, model.Post) (model.Post, error)
	Delete(int) (string, error)
}
