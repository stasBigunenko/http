package inMemory

import (
	"errors"
	"sort"
	"src/http/pkg/model"
	"sync"
)

//Storage functions

type Storage struct {
	mu      sync.Mutex
	storage map[int]model.Post
	idStore int
}

func New() *Storage {
	return &Storage{
		storage: make(map[int]model.Post),
		idStore: 0,
	}
}

//Create function: save received post to the storage and return post struct
func (s *Storage) Create(p model.Post) (model.Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.idStore++
	if p.Author == "" {
		s.idStore--
		return model.Post{}, errors.New("author is empty")
	}
	if p.Message == "" {
		s.idStore--
		return model.Post{}, errors.New("message is empty")
	}
	p.Id = s.idStore
	s.storage[p.Id] = p
	return p, nil
}

//Get function: find in storage requested Id and return Post with the same Id
func (s *Storage) Get(id int) (model.Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var p model.Post
	p, ok := s.storage[id]
	if !ok {
		return model.Post{}, errors.New("post with Id %d not found")
	}
	return p, nil
}

//GetAll function: return slice with all Posts in the Storage
func (s *Storage) GetAll() []model.Post {
	s.mu.Lock()
	defer s.mu.Unlock()
	var p []model.Post
	for _, v := range s.storage {
		p = append(p, v)
	}
	sort.Slice(p, func(i, j int) bool {
		return p[i].Id < p[j].Id
	})
	return p
}

//Update function: find in the Storage requested Id and update it according the data from request
func (s *Storage) Update(p model.Post) (model.Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.storage[p.Id]
	if !ok {
		return model.Post{}, errors.New("post cann't be updated. The post doesn't exist")
	}
	if p.Author == "" {
		p.Author = s.storage[p.Id].Author
	}
	if p.Message == "" {
		p.Message = s.storage[p.Id].Message
	}

	s.storage[p.Id] = p
	return p, nil
}

//Delete function: find in the storage requested Id and delete it from Storage
func (s *Storage) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.storage[id]
	if !ok {
		return errors.New("post can't be deleted - Id not found")
	}
	delete(s.storage, id)
	return nil
}

func (s *Storage) CreateFromFile(p model.Post) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if p.Author == "" {
		return errors.New("author is empty")
	}
	if p.Message == "" {
		return errors.New("message is empty")
	}

	//Checking if storage have the same id and if "yes" create the next one after the last id in the memory
	_, ok := s.storage[p.Id]
	if !ok {
		s.idStore = p.Id
	} else {
		s.idStore += 1
	}
	p.Id = s.idStore
	s.storage[p.Id] = p
	return nil
}
