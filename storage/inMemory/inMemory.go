package inMemory

import (
	"errors"
	"fmt"
	"src/http/pkg/model"
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

//Create function: save received post to the storage and return post struct
func (s *Storage) Create(p model.Post) (model.Post, error) {
	s.IdStor++
	if p.Author == "" {
		s.IdStor--
		return model.Post{}, errors.New("author is empty")
	}
	if p.Message == "" {
		s.IdStor--
		return model.Post{}, errors.New("message is empty")
	}
	p.Id = s.IdStor
	s.Storage[p.Id] = p
	return p, nil
}

//Get function: find in storage requested Id and return Post with the same Id
func (s *Storage) Get(id int) (model.Post, error) {
	var p model.Post
	p, ok := s.Storage[id]
	if !ok {
		return model.Post{}, fmt.Errorf("post with Id %d not found", id)
	}
	return p, nil
}

//GetAll function: return slice with all Posts in the Storage
func (s *Storage) GetAll() []model.Post {
	var p []model.Post
	for _, v := range s.Storage {
		p = append(p, v)
	}
	return p
}

//Update function: find in the Storage requested Id and update it according the data from request
func (s *Storage) Update(p model.Post) (model.Post, error) {
	_, ok := s.Storage[p.Id]
	if !ok {
		return model.Post{}, fmt.Errorf("post cann't be updated. The post doesn't exist")
	}
	if p.Author == "" {
		p.Author = s.Storage[p.Id].Author
	}
	if p.Message == "" {
		p.Message = s.Storage[p.Id].Message
	}

	s.Storage[p.Id] = p
	return p, nil
}

//Delete function: find in the storage requested Id and delete it from Storage
func (s *Storage) Delete(id int) error {
	_, ok := s.Storage[id]
	if !ok {
		return errors.New("post can't be deleted - Id not found")
	}
	delete(s.Storage, id)
	return nil
}

func (s *Storage) CreateFromFile(p model.Post) error {
	s.IdStor = p.Id
	if p.Author == "" {
		s.IdStor--
		return errors.New("author is empty")
	}
	if p.Message == "" {
		s.IdStor--
		return errors.New("message is empty")
	}

	s.Storage[p.Id] = p
	return nil
}
