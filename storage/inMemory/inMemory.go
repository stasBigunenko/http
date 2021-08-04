package inMemory

import (
	"errors"
	_ "errors"
	"fmt"
	"src/http/pkg/model"
	"time"
)

//Storage functions
type Storage struct {
	Storage map[int]model.Post
	IdStor  int
}

func New() *Storage {
	return &Storage{
		Storage: make(map[int]model.Post),
		IdStor:  0,
	}
}

//Create function: save to storage and return Post model
func (s *Storage) Create(p model.Post) (model.Post, error) {
	s.IdStor++
	if p.Message == "" {
		return model.Post{}, errors.New("The message is empty")
	}
	p.Id = s.IdStor
	t := time.Now()
	t.Format("2006-01-02 15:04:05")
	p.Time = t
	s.Storage[s.IdStor] = p
	return p, nil
}

//Get function: find in storage requested Id and return Post model with the same Id
func (s *Storage) Get(Id int) (model.Post, error) {
	var p model.Post
	p, ok := s.Storage[Id]
	if !ok {
		return model.Post{}, fmt.Errorf("Post with Id %d not found", Id)
	}
	return p, nil
}

//GetAll function: return slice with all Posts in the storage
func (s *Storage) GetAll() []model.Post {
	var p []model.Post
	for _, v := range s.Storage {
		p = append(p, v)
	}
	return p
}

//Update function: find in storage requested Id and update it according the data from request
func (s *Storage) Update(p model.Post) (model.Post, error) {
	_, ok := s.Storage[s.IdStor]
	if !ok {
		return model.Post{}, errors.New("Post cann't be updated - Id not found")
	}
	s.Storage[s.IdStor] = p
	return p, nil
}

//Delete function: find in storage requested Id and delete it from storage
func (s *Storage) Delete(IdStor int) (string, error) {
	_, ok := s.Storage[s.IdStor]
	if !ok {
		return "", errors.New("Post cann't be deleted - Id not found")
	}
	delete(s.Storage, IdStor)
	str := "Post deleted"
	return str, nil
}
