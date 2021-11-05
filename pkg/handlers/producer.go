package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"src/http/cmd/http/configHTTP"
	"src/http/pkg/model"
)

func connectProducer() (*amqp.Connection, error) {

	config := configHTTP.Set()

	connStr := "amqp://" + config.RMQLog + ":" + config.RMQPass + "@" + config.RMQPath
	conn, err := amqp.Dial(connStr)
	if err != nil {
		fmt.Println("Failed Initializing Broker Connection")
		panic(err)
	}

	return conn, nil
}

func pushCommentToQueue(topic string, message []byte) error {

	conn, err := connectProducer()
	if err != nil {
		return err
	}

	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
	}
	defer ch.Close()

	// with this channel open, we can then start to interact
	// with the instance and declare Queues that we can publish and
	// subscribe to
	q, err := ch.QueueDeclare(
		topic,
		false,
		false,
		false,
		false,
		nil,
	)
	// We can print out the status of our Queue here
	// this will information like the amount of messages on
	// the queue
	fmt.Println(q)
	// Handle any errors if we were unable to create the queue
	if err != nil {
		fmt.Println(err)
	}

	// attempt to publish a message to the queue!
	err = ch.Publish(
		"",
		topic,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		},
	)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Published Message to Queue")

	return nil
}

// createComment handler
func createChPost(p model.Post) error {

	var ChPost struct {
		Title   string `json:"method"`
		ID      string `json:"id"`
		Author  string `json:"author"`
		Message string `json:"message"`
	}

	pString := p.Id.String()

	if pString == "00000000-0000-0000-0000-000000000000" {
		ChPost.Title = "Create post"
	} else {
		ChPost.Title = "Update post"
	}

	ChPost.ID = pString
	ChPost.Author = p.Author
	ChPost.Message = p.Message

	pInBytes, _ := json.Marshal(ChPost)
	pushCommentToQueue("posts", pInBytes)

	return nil
}

func deleteChPost(p string) error {

	var ChPost struct {
		Title string `json:"method"`
		ID    string `json:"id"`
	}

	ChPost.Title = "Delete post"
	ChPost.ID = p

	pInBytes, _ := json.Marshal(ChPost)
	pushCommentToQueue("posts", pInBytes)

	return nil
}
