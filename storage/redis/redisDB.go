package redisDB

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/go-redis/redis"
	"github.com/google/uuid"

	"src/http/pkg/model"
)

type RedisDB struct {
	Client *redis.Client
	mu     sync.Mutex
}

func New(addr string, db string) *RedisDB {

	rdb, _ := strconv.Atoi(db)

	redisDB := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       rdb,
	})

	val, err := redisDB.Ping().Result()
	fmt.Println(val, err)

	return &RedisDB{Client: redisDB}
}

func (rdb *RedisDB) Create(p model.Post) (model.Post, error) {
	rdb.mu.Lock()
	rdb.mu.Unlock()

	id := uuid.New()
	p.Id = id

	jp, err := json.Marshal(p)
	if err != nil {
		return model.Post{}, errors.New("marshal problem")
	}

	idStr := id.String()

	err = rdb.Client.Set(idStr, jp, 0).Err()
	if err != nil {
		return model.Post{}, errors.New("redis problem")
	}

	return p, nil

}

func (rdb *RedisDB) Get(id uuid.UUID) (model.Post, error) {
	rdb.mu.Lock()
	rdb.mu.Unlock()

	idStr := id.String()

	val, err := rdb.Client.Get(idStr).Bytes()
	if err != nil {
		return model.Post{}, errors.New("redis problem")
	}

	jup := model.Post{}

	err = json.Unmarshal(val, &jup)
	if err != nil {
		return model.Post{}, errors.New("unmarshal problems")
	}

	return jup, nil
}

func (rdb *RedisDB) GetAll() []model.Post {
	rdb.mu.Lock()
	rdb.mu.Unlock()

	var posts []model.Post

	all, err := rdb.Client.Keys("*").Result()
	if err != nil {
		return nil
	}

	for _, val := range all {
		res, err := rdb.Client.Get(val).Bytes()
		if err != nil {
			return nil
		}
		jup := model.Post{}

		err = json.Unmarshal(res, &jup)
		if err != nil {
			return nil
		}
		posts = append(posts, jup)
	}

	return posts
}

func (rdb *RedisDB) Update(p model.Post) (model.Post, error) {
	rdb.mu.Lock()
	rdb.mu.Unlock()

	idStr := p.Id.String()

	val, err := rdb.Client.Get(idStr).Bytes()
	if err != nil {
		return model.Post{}, errors.New("post not found")
	}

	var post model.Post

	err = json.Unmarshal(val, &post)
	if err != nil {
		return model.Post{}, errors.New("marshal problem")
	}

	if post.Author == "" {
		p.Author = post.Author
	}
	if p.Message == "" {
		p.Message = post.Message
	}

	jp, err := json.Marshal(p)
	if err != nil {
		return model.Post{}, errors.New("marshal problem")
	}

	err = rdb.Client.Set(idStr, jp, 0).Err()
	if err != nil {
		return model.Post{}, errors.New("redis problem")
	}

	return p, nil
}

func (rdb *RedisDB) Delete(id uuid.UUID) error {
	rdb.mu.Lock()
	rdb.mu.Unlock()

	idStr := id.String()

	res, err := rdb.Client.Get(idStr).Result()
	fmt.Println(res)
	if err != nil {
		return errors.New("post not found")
	}

	err = rdb.Client.Del(idStr).Err()
	if err != nil {
		return errors.New("redis problem")
	}

	return nil
}

func (rdb *RedisDB) CreateFromFile(p model.Post) error {
	rdb.mu.Lock()
	rdb.mu.Unlock()

	idStr := p.Id.String()

	jp, err := json.Marshal(p)
	if err != nil {
		return errors.New("marshal problem")
	}

	err = rdb.Client.Set(idStr, jp, 0).Err()
	if err != nil {
		return errors.New("redis problem")
	}

	return nil
}
