package services

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"src/http/pkg/model"
	"src/http/storage"
	"strconv"
)

// FilePath :name and path of the *.csv file
const FilePath = "./static/result.csv"

//Service's functions which are working directly with Storage's functions

type Store struct {
	store storage.Storage
}

func NewStore(s storage.Storage) *Store {
	return &Store{
		store: s,
	}
}

func (s *Store) CreateId(post *model.Post) (*model.Post, error) {
	postNew, err := s.store.Create(*post)
	if err != nil {
		return nil, err
	}
	return &postNew, nil
}

func (s *Store) GetId(id int) (*model.Post, error) {
	postId, err := s.store.Get(id)
	if err != nil {
		return nil, err
	}
	return &postId, nil
}

func (s *Store) GetALL() (*[]model.Post, error) {
	var postAll []model.Post
	postAll = s.store.GetAll()
	return &postAll, nil
}

func (s *Store) DeleteId(id int) error {
	err := s.store.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) UpdateId(post *model.Post) (*model.Post, error) {
	postUpdate, err := s.store.Update(*post)
	if err != nil {
		return nil, err
	}
	return &postUpdate, nil
}

func (s *Store) CreatePost(post *model.Post) error {
	err := s.store.CreateFromFile(*post)
	if err != nil {
		return errors.New("couldn't create post from file")
	}
	return nil
}

//Upload function: open the file and save all posts in memory one by one
func (s *Store) Upload() error {
	csvFile, err := os.OpenFile(FilePath, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)

	//Indicate number of fields of our struct
	reader.FieldsPerRecord = 3
	csvData, err := reader.ReadAll()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var post model.Post

	//Go through all read data and call our function CreatePost to save the data in our Storage
	for _, record := range csvData {
		post.Id, err = strconv.Atoi(record[0])
		if err != nil {
			return err
		}
		post.Author = record[1]
		post.Message = record[2]
		s.CreatePost(&post)
	}
	return nil
}

//Download function: create a *.csv file with all our posts which have been saved in memory
func (s *Store) Download(res []model.Post) error {
	csvFile, err := os.Create(FilePath)
	if err != nil {
		return err
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)

	//Write each post on the different row in the file
	for _, val := range res {
		var row []string
		row = append(row, strconv.Itoa(val.Id))
		row = append(row, val.Author)
		row = append(row, val.Message)
		writer.Write(row)
	}
	writer.Flush()
	return nil
}
