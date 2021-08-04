package services

import (
	"src/http/pkg/model"
	"src/http/storage"
)

type Store struct {
	Store storage.Storage
}

func NewStore(s storage.Storage) *Store {
	return &Store{
		Store: s,
	}
}

func (s *Store) CreateId (post *model.Post) (*model.Post, error){
	postNew, _ := s.Store.Create(*post)
	return &postNew, nil
}

func (s *Store) GetId (id int) (*model.Post, error){
	postId, _ := s.Store.Get(id)
	return &postId, nil
}

func (s *Store) GetALL () (*[]model.Post, error) {
	var postAll []model.Post
	postAll, _ = s.Store.GetAll()
	return &postAll, nil
}

func (s *Store) DeleteId (id int) (string, error) {
	s.Store.Delete(id)
	str := "Post deleted"
	return str, nil
}

func (s *Store) UpdateId (post *model.Post) (*model.Post, error) {
	postUpdate, _ := s.Store.Update(*post)
	return &postUpdate, nil
}

