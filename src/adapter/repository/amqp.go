package repository

import (
	"context"
	"encoding/json"
	"log"

	"github.com/giovanibrioni/audit-server/audit"
	"github.com/giovanibrioni/audit-server/helper"
	"github.com/streadway/amqp"
)

type AmqpAuditRepository struct {
	channel *amqp.Channel
	queue   string
	ctx     context.Context
}

func NewAmqpAuditRepository() audit.AuditRepo {
	amqpServerURL := helper.GetEnvOrDefault("AMQP_SERVER_URL", "amqp://guest:guest@localhost:5672/")
	queue := helper.GetEnvOrDefault("AMQP_QUEUE", "audit_logs")

	// Create a new RabbitMQ connection.
	connectRabbitMQ, err := amqp.Dial(amqpServerURL)
	if err != nil {
		panic(err)
	}

	channelRabbitMQ, err := connectRabbitMQ.Channel()
	if err != nil {
		panic(err)
	}
	_, err = channelRabbitMQ.QueueDeclare(
		queue, // queue name
		true,  // durable
		false, // auto delete
		false, // exclusive
		false, // no wait
		nil,   // arguments
	)
	if err != nil {
		panic(err)
	}

	return &AmqpAuditRepository{
		channel: channelRabbitMQ,
		queue:   queue,
		ctx:     context.Background(),
	}
}

func (a *AmqpAuditRepository) SaveBatch(auditLogs []*audit.AuditEntity) error {
	for _, auditLog := range auditLogs {
		encoded, err := json.Marshal(auditLog)
		if err != nil {
			log.Fatal("Unable to marshal auditLogs")
			return err
		}
		message := amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(encoded),
		}
		err = a.channel.Publish(
			"",      // exchange
			a.queue, // queue name
			false,   // mandatory
			false,   // immediate
			message,
		)
		if err != nil {
			log.Fatal("Failed to Publish message: ", err)
			return err
		}
	}

	return nil
}
