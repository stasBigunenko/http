package storage

import "src/http/pkg/model"

var badStorage = MockStorage{}

type MockStorage struct {
	MockCreate         func(p model.Post) (model.Post, error)
	MockGet            func(id int) (model.Post, error)
	MockGetAll         func() []model.Post
	MockUpdate         func(p model.Post) (model.Post, error)
	MockDelete         func(id int) error
	MockCreateFromFile func(p model.Post) error
}

func (m MockStorage) Create(p model.Post) (model.Post, error) {
	return m.MockCreate(p)
}

func (m MockStorage) Get(id int) (model.Post, error) {
	return m.MockGet(id)
}

func (m MockStorage) GetAll() []model.Post {
	return m.MockGetAll()
}

func (m MockStorage) Update(p model.Post) (model.Post, error) {
	return m.MockUpdate(p)
}

func (m MockStorage) Delete(id int) error {
	return m.MockDelete(id)
}

func (m MockStorage) CreateFromFile(p model.Post) error {
	return m.MockCreateFromFile(p)
}
