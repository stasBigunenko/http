package inMemory

import (
	_ "errors"
	"fmt"
	"src/http/pkg/model"
	"time"
)

type Storage struct {
	Storage map[int]model.Post
	IdStor     int
}

func New() *Storage {
	return &Storage{
		Storage: make(map[int]model.Post),
		IdStor: 0,
	}
}

func (s *Storage) Create(p model.Post) (model.Post, error) {
	s.IdStor++
	p.Id = s.IdStor
	t := time.Now()
	t.Format("2006-01-02 15:04:05")
	p.Time = t
	s.Storage[s.IdStor] = p
	return p, nil
}

func (s *Storage) Get(Id int) (model.Post, error) {
	var p model.Post
	p, ok := s.Storage[Id]
	if !ok {
		return model.Post{}, fmt.Errorf("Post with Id %d not found", Id)
	}
	return p, nil
}

func (s *Storage) GetAll() ([]model.Post, error) {
	p := []model.Post{}
	for _, v := range s.Storage {
		p = append(p, v)
	}
	return p, nil
}

func (s *Storage) Update(p model.Post) (model.Post, error) {
	s.Storage[s.IdStor] = p
	return p, nil
}

func (s *Storage) Delete(IdStor int) (string, error) {
	delete(s.Storage, IdStor)
	str := "Post deleted"
	return str, nil
}
