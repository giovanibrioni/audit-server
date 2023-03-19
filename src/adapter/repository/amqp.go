package repository

import (
	"context"

	"github.com/goccy/go-json"

	"github.com/giovanibrioni/audit-server/audit"
	"github.com/giovanibrioni/audit-server/helper"
	rabbitmq "github.com/wagslane/go-rabbitmq"
	"go.uber.org/zap"
)

type AmqpAuditRepository struct {
	publisher *rabbitmq.Publisher
	queue     string
	ctx       context.Context
	logger    *zap.SugaredLogger
}

func NewAmqpAuditRepository(ctx context.Context, logger *zap.SugaredLogger) audit.AuditRepo {
	amqpServerURL := helper.GetEnvOrDefault("AMQP_SERVER_URL", "amqp://guest:guest@localhost:5672/")
	queue := helper.GetEnvOrDefault("AMQP_QUEUE", "audit_logs")

	publisher, err := rabbitmq.NewPublisher(
		amqpServerURL,
		rabbitmq.Config{},
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsReconnectInterval(0),
	)
	if err != nil {
		logger.Fatal(err)
		//panic(err)
	}
	returns := publisher.NotifyReturn()
	go func() {
		for r := range returns {
			logger.Errorf("failed to publish message: %s, on rabbitmq queue: %s", string(r.Body), queue)
		}
	}()
	confirmations := publisher.NotifyPublish()
	go func() {
		for c := range confirmations {
			logger.Errorf("message confirmed from server. queue: %s, tag: %v, ack: %v.", queue, c.DeliveryTag, c.Ack)
		}
	}()

	return &AmqpAuditRepository{
		publisher: publisher,
		queue:     queue,
		ctx:       ctx,
		logger:    logger,
	}
}

func (a *AmqpAuditRepository) SaveBatch(auditLogs []*audit.AuditEntity) error {
	for _, auditLog := range auditLogs {
		encoded, err := json.Marshal(auditLog)
		if err != nil {
			a.logger.Fatal("Unable to marshal auditLogs")
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
			a.logger.Fatalf("failed to publish message on RabbitMq: JobId: %s, AuditId: %s, error: %s", auditLog.JobId, auditLog.AuditId, err)
			return err
		}
	}

	return nil
}
