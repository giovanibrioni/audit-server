package repository

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/giovanibrioni/audit-server/audit"
	"github.com/giovanibrioni/audit-server/helper"
)

type kafkaAuditRepository struct {
	producer *kafka.Producer
	topic    string
	ctx      context.Context
}

func NewKafkaAuditRepository() audit.AuditRepo {
	bootstrapServers := helper.GetEnvOrDefault("KAFKA_URL", "localhost:9092")
	topic := helper.GetEnvOrDefault("KAFKA_TOPIC", "audit_logs")
	config := kafka.ConfigMap{"bootstrap.servers": bootstrapServers, "acks": "1", "allow.auto.create.topics": true}

	p, err := kafka.NewProducer(&config)

	if err != nil {
		log.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}

	// Listen to all the events on the default events channel
	go listenKafkaEvents(p)

	return &kafkaAuditRepository{
		producer: p,
		topic:    topic,
		ctx:      context.Background(),
	}
}

func (k *kafkaAuditRepository) Save(auditLog *audit.AuditEntity) error {
	encoded, err := json.Marshal(auditLog)
	if err != nil {
		log.Fatal("Unable to marshal auditLogs")
		return err
	}
	err = k.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &k.topic, Partition: kafka.PartitionAny},
		Value:          []byte(encoded),
	}, nil)
	if err != nil {
		if err.(kafka.Error).Code() == kafka.ErrQueueFull {
			// Producer queue is full, wait 1s for messages
			// to be delivered then try again.
			time.Sleep(time.Second)
		}
		log.Printf("Failed to produce message: %v\n", err)
	}

	return nil
}

func listenKafkaEvents(p *kafka.Producer) {
	for e := range p.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			// The message delivery report, indicating success or
			// permanent failure after retries have been exhausted.
			// Application level retries won't help since the client
			// is already configured to do that.
			m := ev
			if m.TopicPartition.Error != nil {
				log.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
			} else {
				log.Printf("Delivered message to topic %s [%d] at offset %v\n",
					*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
			}
		case kafka.Error:
			log.Printf("Error: %v\n", ev)
		default:
			log.Printf("Ignored event: %s\n", ev)
		}
	}
}
