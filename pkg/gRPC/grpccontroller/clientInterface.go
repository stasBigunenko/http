package grpccontroller

import "src/http/pkg/model"

type ClientGRPCInterface interface {
	Get(string) (*model.Post, error)
	GetAll() []model.Post
	Create(pp model.Post) (model.Post, error)
	Update(pp model.Post) (model.Post, error)
	Delete(id string) error
	CreateFromFile(pp model.Post) error
}
