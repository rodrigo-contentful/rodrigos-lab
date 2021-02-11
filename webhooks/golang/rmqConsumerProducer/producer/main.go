package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/streadway/amqp"
)

//SysInner ...
type SysInner struct {
	ID       string `json:"id"`
	LinkType string `json:"linkType"`
	Type     string `json:"type"`
}

//ContentType ...
type ContentType struct {
	Sys SysInner `json:"sys"`
}

// Environment ...
type Environment struct {
	Sys SysInner `json:"sys"`
}

// Space ..
type Space struct {
	Sys SysInner `json:"sys"`
}

// Tags ...
type Tags struct {
	Tags interface{} `json:"tags"`
}

// Sys define sthe sys part of the payload
type Sys struct {
	ID        string      `json:"id"`
	Revision  int         `json:"revision"`
	Type      string      `json:"type"`
	Space     Space       `json:"space"`
	CT        ContentType `json:"contentType"`
	CreatedAt string      `json:"createdAt"`
	UpdatedAt string      `json:"updatedAt"`
	Env       Environment `json:"environment"`
}

// Payload defines the JSON payload from getContentType
type Payload struct {
	Fields   interface{} `json:"fields"`
	Metadata Tags        `json:"metadata"`
	Sys      Sys         `json:"sys"`
}

const (
	// config for RMQ
	numDeliveries = 100000000
	batchSize     = 10000
	// RMQ connection
	rmqUsername  = "guest"
	rmqPassword  = "guest"
	rmqHost      = "localhost"
	rmqPort      = "5672"
	rmqProtocol  = "amqp"
	rmqQueueName = "TestQueue"
	// config for endpoint
	listenAndServePort = "8081"
)

func main() {

	connRMQ, err := connectRMQ()
	if err != nil {
		fmt.Println("Failed Initializing Broker Connection")
		panic(err)
	}
	fmt.Printf("Loading listener - 'localhost: %s'", listenAndServePort)
	listener(connRMQ)

}

// listener creates POST endpoint for requests
func listener(connRMQ *amqp.Connection) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		decoder := json.NewDecoder(r.Body)
		var t Payload
		err := decoder.Decode(&t)
		if err != nil {
			panic(err)
		}

		rmqPost(t, connRMQ)
	})

	lsport := fmt.Sprintf(":%s", listenAndServePort)
	http.ListenAndServe(lsport, nil)
}

// connectRMQ creates a RMQ connection
func connectRMQ() (*amqp.Connection, error) {
	dial := fmt.Sprintf("%s://%s:%s@%s:%s/", rmqProtocol, rmqUsername, rmqPassword, rmqHost, rmqPort)
	fmt.Println("Connecting RabbitMQ")
	fmt.Println(dial)
	return amqp.Dial(dial)
}

// rmqPost will post a message in RMQ
func rmqPost(t Payload, conn *amqp.Connection) {
	// Let's start by opening a channel to our RabbitMQ instance
	// over the connection we have already established
	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
	}
	defer ch.Close()

	// with this channel open, we can then start to interact
	// with the instance and declare Queues that we can publish and
	// subscribe to
	q, err := ch.QueueDeclare(
		rmqQueueName,
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

	//
	j, err := json.Marshal(t)
	if err != nil {
		log.Fatalf("Error occured during marshaling. Error: %s", err.Error())
	}
	fmt.Printf("Message JSON: %s\n", string(j))

	// attempt to publish a message to the queue!
	err = ch.Publish(
		"",
		rmqQueueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        j,
		},
	)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Published Message to Queue")

}
