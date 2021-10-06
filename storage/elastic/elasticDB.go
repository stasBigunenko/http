package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	es "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/google/uuid"
	"log"
	"src/http/pkg/model"
	"sync"
)

const index_name string = "posts"

type ElkPost struct {
	Source model.Post `json:"_source"`
}

type ElkPosts struct {
	Hits struct {
		Hits []struct {
			Source model.Post `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

type ElasticSearch struct {
	client *es.Client
	mu     sync.Mutex
}

func NewElastic(addr string) (*ElasticSearch, error) {

	a := []string{addr}
	cfg := es.Config{
		Addresses: a,
		Username:  "",
		Password:  "",
	}

	fmt.Println(cfg)

	client, err := es.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	res, err := client.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	fmt.Printf("===================================%v===================================\n", res)
	defer res.Body.Close()
	return &ElasticSearch{
		client: client,
	}, nil
}

func (e *ElasticSearch) Create(p model.Post) (model.Post, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	id := uuid.New()
	p.Id = id
	idStr := id.String()

	post, err := json.Marshal(p)
	if err != nil {
		return model.Post{}, errors.New("elastic problems")
	}

	req := esapi.CreateRequest{
		Index:      index_name,
		DocumentID: idStr,
		Body:       bytes.NewReader(post),
	}

	res, err := req.Do(context.Background(), e.client)
	if err != nil {
		return model.Post{}, errors.New("elastic problems")
	}
	log.Println(res)
	defer res.Body.Close()

	return p, nil
}

func (e *ElasticSearch) Get(id uuid.UUID) (model.Post, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	idStr := id.String()

	req := esapi.GetRequest{
		Index:      index_name,
		DocumentID: idStr,
	}

	res, err := req.Do(context.Background(), e.client)
	if err != nil {
		return model.Post{}, err
	}

	if res.StatusCode == 404 {
		return model.Post{}, errors.New("not found")
	}

	var post ElkPost
	err = json.NewDecoder(res.Body).Decode(&post)
	if err != nil {
		return model.Post{}, errors.New("internal problem")
	}

	return post.Source, nil
}

func (e *ElasticSearch) GetAll() []model.Post {
	e.mu.Lock()
	defer e.mu.Unlock()

	req := esapi.SearchRequest{
		Index: []string{index_name},
	}

	res, err := req.Do(context.Background(), e.client)
	if err != nil {
		return nil
	}

	var all ElkPosts

	err = json.NewDecoder(res.Body).Decode(&all)
	if err != nil {
		return nil
	}

	posts := []model.Post{}

	for _, p := range all.Hits.Hits {
		posts = append(posts, p.Source)
	}

	return posts
}

func (e *ElasticSearch) Update(p model.Post) (model.Post, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	var post ElkPost
	idStr := p.Id.String()

	if p.Author == "" || p.Message == "" {

		req := esapi.GetRequest{
			Index:      index_name,
			DocumentID: idStr,
		}

		res, err := req.Do(context.Background(), e.client)
		if err != nil {
			return model.Post{}, err
		}

		err = json.NewDecoder(res.Body).Decode(&post)
		if err != nil {
			return model.Post{}, errors.New("internal problem")
		}
	}

	if p.Author == "" {
		p.Author = post.Source.Author
	}

	if p.Message == "" {
		p.Message = post.Source.Message
	}

	pj, err := json.Marshal(p)
	if err != nil {
		return model.Post{}, errors.New("internal marshal problems")
	}

	req := esapi.UpdateRequest{
		Index:      index_name,
		DocumentID: idStr,
		Body:       bytes.NewReader([]byte(fmt.Sprintf(`{"doc":%s}`, pj))),
	}

	res, err := req.Do(context.Background(), e.client)
	if err != nil {
		return model.Post{}, errors.New("internal elk problems")
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return model.Post{}, errors.New("not found")
	}

	return p, nil
}

func (e *ElasticSearch) Delete(id uuid.UUID) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	idStr := id.String()

	req := esapi.DeleteRequest{
		Index:      index_name,
		DocumentID: idStr,
	}

	res, err := req.Do(context.Background(), e.client)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return errors.New("not found")
	}

	return nil
}

func (e *ElasticSearch) CreateFromFile(p model.Post) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	idStr := p.Id.String()

	post, err := json.Marshal(p)
	if err != nil {
		return errors.New("elastic problems")
	}

	req := esapi.CreateRequest{
		Index:      index_name,
		DocumentID: idStr,
		Body:       bytes.NewReader(post),
	}

	res, err := req.Do(context.Background(), e.client)
	if err != nil {
		return errors.New("elastic problems")
	}
	defer res.Body.Close()

	return nil
}
