package services

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"github.com/google/uuid"
	"github.com/jszwec/csvutil"

	"src/http/pkg/model"
	"src/http/storage"
)

//FilePath :name and path of the *.csv file
const FilePath = "./static/"

//Service's functions which are working directly with Storage's functions

type Services struct {
	store storage.Storage
}

func NewService(s storage.Storage) Services {
	return Services{
		store: s,
	}
}

func (s *Services) CreateId(post *model.Post) (*model.Post, error) {

	if post.Author == "" {
		return nil, errors.New("author is empty")
	}
	if post.Message == "" {
		return nil, errors.New("message is empty")
	}

	postNew, err := s.store.Create(*post)
	if err != nil {
		return nil, err
	}
	return &postNew, nil
}

func (s *Services) GetId(id string) (*model.Post, error) {

	val, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("couldn't parse id")
	}

	postId, err := s.store.Get(val)
	if err != nil {
		return nil, err
	}

	return &postId, nil
}

func (s *Services) GetALL() *[]model.Post {
	var postAll []model.Post
	postAll = s.store.GetAll()
	return &postAll
}

func (s *Services) DeleteId(id string) error {

	val, err := uuid.Parse(id)
	if err != nil {
		return errors.New("couldn't parse id")
	}

	err = s.store.Delete(val)
	fmt.Println(err)
	if err != nil {
		return err
	}
	return nil
}

func (s *Services) UpdateId(post *model.Post) (*model.Post, error) {
	postUpdate, err := s.store.Update(*post)
	if err != nil {
		return nil, err
	}
	return &postUpdate, nil
}

func (s *Services) CreatePost(post *model.Post) error {

	if post.Author == "" {
		return errors.New("author is empty")
	}
	if post.Message == "" {
		return errors.New("message is empty")
	}

	err := s.store.CreateFromFile(*post)
	if err != nil {
		return errors.New("couldn't create post from file")
	}
	return nil
}

//Upload function: open the file and save all posts in memory one by one
func (s *Services) Upload(file multipart.File) error {

	reader := csv.NewReader(file)

	//Indicate number of fields of our struct
	reader.FieldsPerRecord = 3
	for {
		csvData, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			os.Exit(1)
		}

		//pass the headers of the csv file
		if csvData[0] == "Id" {
			continue
		}

		var post model.Post

		//Go through read data and call our function CreatePost to save the data in our Storage
		val, err := uuid.Parse(csvData[0])
		if err != nil {
			return errors.New("couldn't parse id")
		}
		post.Id = val
		post.Author = csvData[1]
		post.Message = csvData[2]
		s.CreatePost(&post)
	}

	return nil
}

//Download function: create a *.csv file with all our posts which have been saved in memory
func (s *Services) Download() ([]byte, error) {
	allPosts := s.GetALL()

	ap, err := csvutil.Marshal(allPosts)
	if err != nil {
		return nil, err
	}

	return ap, nil
}
