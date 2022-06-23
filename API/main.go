package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
)

type Data struct {
	Id         int    `json:"id"`
	USER_ID    int    `json:"user_id"`
	FIRST_NAME string `json:"first_name"`
	LAST_NAME  string `json:"last_name"`
	BIIN       string `json:"biin"`
	EMAIL      string `json:"email"`
	PHONE      string `json:"phone"`
	PASSWRD    string `json:"passwrd"`
	Type       string `json:"type"`
}

func create(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var post Data
	err := json.Unmarshal(reqBody, &post)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(post)
	}
	Exchange(&post)
}

func handleReqs() {

	r := mux.NewRouter()
	r.HandleFunc("/post", create).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", r))
}

////////////////

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
func bodyFrom(args []string) string {
	var s string
	if (len(args) < 3) || os.Args[2] == "" {
		s = "Hello"
	} else {
		s = strings.Join(args[2:], " ")
	}
	return s
}

func Exchange(message *Data) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"create_exchange", // name
		"topic",           // type
		true,              // durable
		false,             // auto-deleted
		false,             // internal
		false,             // no-wait
		nil,               // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	body, err := json.Marshal(message)

	err = ch.Publish(
		"create_exchange", // exchange
		message.Type,      // routing key
		false,             // mandatory
		false,             // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")

	log.Printf(" [x] Sent %s", body)
}
func main() {
	handleReqs()
}
