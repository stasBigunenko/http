package storage

import (
	"github.com/google/uuid"

	"src/http/pkg/model"
)

//Storage Interface

type Storage interface {
	Create(model.Post) (model.Post, error)
	Get(uuid.UUID) (model.Post, error)
	GetAll() []model.Post
	Update(model.Post) (model.Post, error)
	Delete(uuid.UUID) error
	CreateFromFile(model.Post) error
}
