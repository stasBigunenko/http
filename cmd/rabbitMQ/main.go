package main

import (
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

func main() {

	path := os.Getenv("RMQPath")
	if path == "" {
		path = "localhost:5672/"
	}

	login := os.Getenv("RMQ_LOG")
	if login == "" {
		login = "guest"
	}

	pass := os.Getenv("RMQ_PASS")
	if pass == "" {
		pass = "guest"
	}

	connStr := "amqp://" + login + ":" + pass + "@" + path
	conn, err := amqp.Dial(connStr)
	if err != nil {
		fmt.Println("Failed Initializing Broker Connection")
		panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
	}
	defer ch.Close()

	if err != nil {
		fmt.Println(err)
	}

	msgs, err := ch.Consume(
		"posts",
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("Recieved Message: %s\n", d.Body)
		}
	}()

	log.Println("Successfully Connected to RabbitMQ Instance")
	log.Println(" [*] - Waiting for messages")
	<-forever
}
