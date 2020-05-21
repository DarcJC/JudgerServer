package msgq

import (
	"JudgerServer/config"

	"github.com/streadway/amqp"
)

var conn *amqp.Connection

func init() {
	amqpuri, err := config.GetConfig("RMQ_URL")
	if err != nil {
		panic(err)
	}
	c, err := amqp.Dial(amqpuri)
	if err != nil {
		panic(err)
	}
	conn = c
}
