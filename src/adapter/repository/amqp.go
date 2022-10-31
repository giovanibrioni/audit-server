package repository

import (
	"context"
	"encoding/json"
	"log"

	"github.com/giovanibrioni/audit-server/audit"
	"github.com/giovanibrioni/audit-server/helper"
	rabbitmq "github.com/wagslane/go-rabbitmq"
)

type AmqpAuditRepository struct {
	publisher *rabbitmq.Publisher
	queue     string
	ctx       context.Context
}

func NewAmqpAuditRepository() audit.AuditRepo {
	amqpServerURL := helper.GetEnvOrDefault("AMQP_SERVER_URL", "amqp://guest:guest@localhost:5672/")
	queue := helper.GetEnvOrDefault("AMQP_QUEUE", "audit_logs")

	publisher, err := rabbitmq.NewPublisher(
		amqpServerURL,
		rabbitmq.Config{},
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsReconnectInterval(0),
	)
	if err != nil {
		log.Fatal(err)
		//panic(err)
	}
	returns := publisher.NotifyReturn()
	go func() {
		for r := range returns {
			log.Printf("failed to publish message: %s, on rabbitmq queue: %s", string(r.Body), queue)
		}
	}()
	confirmations := publisher.NotifyPublish()
	go func() {
		for c := range confirmations {
			log.Printf("message confirmed from server. queue: %s, tag: %v, ack: %v.", queue, c.DeliveryTag, c.Ack)
		}
	}()

	return &AmqpAuditRepository{
		publisher: publisher,
		queue:     queue,
		ctx:       context.Background(),
	}
}

func (a *AmqpAuditRepository) SaveBatch(auditLogs []*audit.AuditEntity) error {
	for _, auditLog := range auditLogs {
		encoded, err := json.Marshal(auditLog)
		if err != nil {
			log.Fatal("Unable to marshal auditLogs")
			return err
		}
		err = a.publisher.Publish(
			[]byte(encoded),
			[]string{a.queue},
			rabbitmq.WithPublishOptionsContentType("application/json"),
			rabbitmq.WithPublishOptionsMandatory,
			rabbitmq.WithPublishOptionsPersistentDelivery,
		)
		if err != nil {
			log.Fatalf("failed to publish message on RabbitMq: JobId: %s, AuditId: %s, error: %s", auditLog.JobId, auditLog.AuditId, err)
			return err
		}
	}

	return nil
}
