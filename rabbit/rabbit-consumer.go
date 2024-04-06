package main;

import (
	"fmt"
	"log"
	"github.com/streadway/amqp"
)

func main() {
	fmt.Println("Consumer starting")

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		fmt.Println(err)
		return;
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
		return;
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		"TestQueue",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err);
		return;
	}

	msgs, err := ch.Consume(
		"TestQueue",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
		return;
	}

	forever := make(chan bool);
	go func() {
		for d := range msgs {
			log.Print(string(d.Body));
		}
	}();

	fmt.Println("Successfully connected to RabbitMQ instance");
	fmt.Println(" [*] - waiting for messages");
	<-forever;
}