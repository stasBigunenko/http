package storage

import (
	"src/http/pkg/model"
)

//Storage Interface

type Storage interface {
	Create(model.Post) (model.Post, error)
	Get(int) (model.Post, error)
	GetAll() []model.Post
	Update(model.Post) (model.Post, error)
	Delete(int) error
}
