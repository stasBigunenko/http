package services

import (
	"errors"
	"src/http/pkg/model"
	"src/http/storage"
)

//Service's functions which are working directly with storage
type Store struct {
	Store storage.Storage
}

func NewStore(s storage.Storage) *Store {
	return &Store{
		Store: s,
	}
}

func (s *Store) CreateId(post *model.Post) (*model.Post, error) {
	postNew, err := s.Store.Create(*post)
	if err != nil {
		return nil, errors.New("Couldn't create post")
	}
	return &postNew, nil
}

func (s *Store) GetId(id int) (*model.Post, error) {
	postId, err := s.Store.Get(id)
	if err != nil {
		return nil, errors.New("Couldn't get Id")
	}
	return &postId, nil
}

func (s *Store) GetALL() (*[]model.Post, error) {
	var postAll []model.Post
	postAll = s.Store.GetAll()
	return &postAll, nil
}

func (s *Store) DeleteId(id int) (string, error) {
	_, err := s.Store.Delete(id)
	if err != nil {
		return "", errors.New("Couldn't get posts")
	}
	str := "Post deleted"
	return str, nil
}

func (s *Store) UpdateId(post *model.Post) (*model.Post, error) {
	postUpdate, err := s.Store.Update(*post)
	if err != nil {
		return nil, errors.New("Couldn't update post")
	}
	return &postUpdate, nil
}
