package services

import (
	"mime/multipart"
	"src/http/pkg/model"
)

//ServicesInteface

type ServicesInterface interface {
	CreateId(post *model.Post) (*model.Post, error)
	GetId(id string) (*model.Post, error)
	GetALL() *[]model.Post
	DeleteId(id string) error
	UpdateId(post *model.Post) (*model.Post, error)
	CreatePost(post *model.Post) error
	Upload(file multipart.File) error
	Download() ([]byte, error)
}
