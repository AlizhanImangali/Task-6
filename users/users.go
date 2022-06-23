package main

import (
	"log"
	//"os"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
)

type Info struct {
	data []Data
}
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

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
func DB() *sql.DB {
	connStr := "user=postgres password=1234 dbname=Test sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	fmt.Println("Succesfully connected")
	//defer db.Close()
	return db
}
func DBInsert(response *Info) {
	var db = DB()
	var result sql.Result
	var err error
	for _, r := range response.data {
		result, err = db.Exec("Insert Into users(user_id,first_name,last_name,biin,email, phone,passwrd,type)  Values ($1,$2,$3,$4,$5,$6,$7,$8)",
			r.USER_ID, r.FIRST_NAME, r.LAST_NAME, r.BIIN, r.EMAIL, r.PHONE, r.PASSWRD, r.Type)
		if err != nil {
			panic(err)
		}
	}
	defer db.Close()
	if result != nil {
		fmt.Println(result)
	}
}

func main() {

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

	q, err := ch.QueueDeclare(
		"users_queue", // name
		false,         // durable
		false,         // delete when unused
		true,          // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,            // queue name
		"users",           // routing key
		"create_exchange", // exchange
		false,
		nil)
	failOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name,           // queue
		"users_consumer", // consumer
		true,             // auto ack
		false,            // exclusive
		false,            // no local
		false,            // no wait
		nil,              // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf(" [x] %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
