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

	err = channel.ExchangeDeclare("events", "topic", true, false, false, false, nil)

	if err != nil {
		panic(err)
	}

	message := amqp.Publishing{
		Body: []byte("Please clap for the Demo Gods - IT WORKED!!!"),
	}

	err = channel.Publish("events", "random-key", false, false, message)

	if err != nil {
		panic("error publishing a message to the queue:" + err.Error())
	}

	_, err = channel.QueueDeclare(queueName, true, false, false, false, nil)

	if err != nil {
		panic("error declaring the queue: " + err.Error())
	}

	err = channel.QueueBind(queueName, "#", "events", false, nil)

	if err != nil {
		panic("error binding to the queue: " + err.Error())
	}

	defer connection.Close()

}
