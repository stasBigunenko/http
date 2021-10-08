package mongoDB

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"src/http/pkg/model"
)

type MongoDB struct {
	Mdb *mongo.Client
	mu  sync.Mutex
}

func NewMongo(user string, psw string, addr string) *MongoDB {

	clientOptions := options.Client().ApplyURI("mongodb://" + user + ":" + psw + "@" + addr)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Println(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		fmt.Println(err)
	}

	return &MongoDB{Mdb: client}
}

func (mdb *MongoDB) Create(p model.Post) (model.Post, error) {
	mdb.mu.Lock()
	defer mdb.mu.Unlock()

	collection := mdb.Mdb.Database("posts").Collection("posts")

	id := uuid.New()
	p.Id = id

	_, err := collection.InsertOne(context.TODO(), p)
	if err != nil {
		return model.Post{}, errors.New("internal mongo storage problem")
	}

	return p, nil
}
func (mdb *MongoDB) Get(id uuid.UUID) (model.Post, error) {
	collection := mdb.Mdb.Database("posts").Collection("posts")

	filter := bson.D{{"id", id}}

	var p model.Post

	err := collection.FindOne(context.TODO(), filter).Decode(&p)
	if err != nil {
		return model.Post{}, errors.New("couldn't find post in db")
	}

	return p, nil
}
func (mdb *MongoDB) GetAll() []model.Post {
	mdb.mu.Lock()
	defer mdb.mu.Unlock()

	collection := mdb.Mdb.Database("posts").Collection("posts")

	var posts []model.Post

	findOptions := options.Find()

	res, _ := collection.Find(context.TODO(), bson.D{{}}, findOptions)

	for res.Next(context.TODO()) {
		var p model.Post
		err := res.Decode(&p)
		if err != nil {
			return []model.Post{}
		}
		posts = append(posts, p)
	}

	if err := res.Err(); err != nil {
		log.Fatal(err)
	}

	return posts
}

func (mdb *MongoDB) Update(p model.Post) (model.Post, error) {
	mdb.mu.Lock()
	defer mdb.mu.Unlock()

	collection := mdb.Mdb.Database("posts").Collection("posts")

	filter := bson.D{{"id", p.Id}}

	var oldPost model.Post

	err := collection.FindOne(context.TODO(), filter).Decode(&oldPost)
	if err != nil {
		return model.Post{}, errors.New("couldn't find post")
	}

	if p.Author == "" {
		p.Author = oldPost.Author
	}

	if p.Message == "" {
		p.Message = oldPost.Message
	}

	update := bson.D{
		{"$set", bson.D{
			{"author", p.Author},
			{"message", p.Message},
		}},
	}

	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return model.Post{}, errors.New("couldn't update post in mongo")
	}

	return p, nil
}
func (mdb *MongoDB) Delete(id uuid.UUID) error {
	mdb.mu.Lock()
	defer mdb.mu.Unlock()

	collection := mdb.Mdb.Database("posts").Collection("posts")

	filter := bson.D{{"id", id}}

	_, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return errors.New("couldn't delete post in mongo")
	}
	return nil
}
func (mdb *MongoDB) CreateFromFile(p model.Post) error {
	mdb.mu.Lock()
	defer mdb.mu.Unlock()

	collection := mdb.Mdb.Database("posts").Collection("posts")

	_, err := collection.InsertOne(context.TODO(), p)
	if err != nil {
		return errors.New("internal mongo problem in mongo db")
	}

	return nil
}
