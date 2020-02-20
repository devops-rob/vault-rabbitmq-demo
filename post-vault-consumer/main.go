package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/vault/api"
	"github.com/streadway/amqp"
)

func main() {
	var token = os.Getenv("TOKEN")
	var vaultaddr = os.Getenv("VAULT_ADDR")

	if token == "" {
		token = "devopsrob"
	}

	if vaultaddr == "" {
		vaultaddr = "http://127.0.0.1:8200"
	}

	config := &api.Config{
		Address: vaultaddr,
	}

	client, err := api.NewClient(config)
	if err != nil {
		fmt.Println(err)
		return
	}

	client.SetToken(token)
	c := client.Logical()

	rbmqsecret, err := c.Read("rabbitmq/creds/rabbitrole")
	if err != nil {
		fmt.Println(err)
		return
	}

	vusername := rbmqsecret.Data["username"].(string)
	vpassword := rbmqsecret.Data["password"].(string)

	queueName := "post-vault"

	url := fmt.Sprintf("amqp://%s:%s@127.0.0.1:5672", vusername, vpassword)
	connection, err := amqp.Dial(url)

	if err != nil {
		panic("could not establish connection with RabbitMQ:" + err.Error())
	}

	channel, err := connection.Channel()

	if err != nil {
		panic("could not open RabbitMQ channel:" + err.Error())
	}

	msgs, err := channel.Consume(queueName, "", false, false, false, false, nil)

	if err != nil {
		panic("error consuming the queue: " + err.Error())
	}

	for msg := range msgs {
		fmt.Println("message received: " + string(msg.Body))
		msg.Ack(true)
	}

	defer connection.Close()

}
