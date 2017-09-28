// @APIVersion 1.0.0
// @APITitle Small Proxy Application for Booming Games - Code Challenge
// @APIDescription Code Challenge
// @Contact bruno.milani@gmail.com

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/satori/go.uuid"
	"github.com/streadway/amqp"
)

var (
	res string
)

type Client struct {
	ID string `json:"client_id"`
}

func init() {
	flag.Parse()
}

func main() {

	RabbitHost := os.Getenv("RABBIT_HOST")
	QueueReplyTo := uuid.NewV4().String()

	conn, err := amqp.Dial("amqp://guest:guest@" + RabbitHost + ":5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// queue declare
	q, err := ch.QueueDeclare(
		QueueReplyTo, // name
		false,        // durable
		true,         // delete when unused
		true,         // exclusive
		false,        // noWait
		nil,          // arguments
	)

	failOnError(err, "Error declare queue")

	//channel consume
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	failOnError(err, "Error Listener")

	http.HandleFunc("/clients.json", func(w http.ResponseWriter, r *http.Request) {

		// correlation id
		correlationID := uuid.NewV4().String()

		err = ch.Publish(
			"",                // exchange
			"backend.clients", // routing key
			false,             // mandatory
			false,             // immediate
			amqp.Publishing{
				ContentType:   "text/plain",
				CorrelationId: correlationID,
				ReplyTo:       QueueReplyTo,
			})

		failOnError(err, "Error publish backend.clients queue")

		for d := range msgs {
			if correlationID == d.CorrelationId {
				res = string(d.Body)
				break
			}
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, res)

	})

	http.HandleFunc("/invoices.json", func(w http.ResponseWriter, r *http.Request) {

		// correlation id
		correlationID := uuid.NewV4().String()

		// client id
		paramClientID := r.URL.Query().Get("client_id")

		client := &Client{ID: paramClientID}
		b, err := json.Marshal(client)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(b))

		if paramClientID != "" {

			err = ch.Publish(
				"",                 // exchange
				"backend.invoices", // routing key
				false,              // mandatory
				false,              // immediate
				amqp.Publishing{
					ContentType:   "application/json",
					CorrelationId: correlationID,
					ReplyTo:       QueueReplyTo,
					Body:          []byte(string(b)),
				})

			failOnError(err, "Error publish backend.invoices queue")

			for d := range msgs {
				if correlationID == d.CorrelationId {
					res = string(d.Body)
					break
				}
			}

			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, res)

		} else {

		}

	})

	http.ListenAndServe(":8020", nil)

}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
