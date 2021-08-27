package services

import (
	"mime/multipart"
	"src/http/pkg/model"
)

type ServicesInterface interface {
	CreateId(post *model.Post) (*model.Post, error)
	GetId(id int) (*model.Post, error)
	GetALL() (*[]model.Post, error)
	DeleteId(id int) error
	UpdateId(post *model.Post) (*model.Post, error)
	CreatePost(post *model.Post) error
	Upload(file multipart.File) error
	Download() ([]byte, error)
}
