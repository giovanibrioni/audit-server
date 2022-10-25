package repository

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/giovanibrioni/audit-server/audit"
	"github.com/giovanibrioni/audit-server/helper"
	"github.com/segmentio/kafka-go"
)

const topic = "audit_logs"

type kafkaAuditRepository struct {
	writer *kafka.Writer
	ctx    context.Context
}

func NewKafkaAuditRepository() audit.AuditRepo {
	brokerAddress := helper.GetEnvOrDefault("KAFKA_URL", "localhost:9092")
	l := log.New(os.Stdout, "kafka writer: ", 0)
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{brokerAddress},
		Topic:   topic,
		// assign the logger to the writer
		Logger: l,
	})

	return &kafkaAuditRepository{
		writer: w,
		ctx:    context.Background(),
	}
}

func (k *kafkaAuditRepository) Save(auditLog *audit.AuditEntity) error {
	encoded, err := json.Marshal(auditLog)
	if err != nil {
		log.Fatal("Unable to marshal auditLogs")
		return err
	}
	err = k.writer.WriteMessages(k.ctx, kafka.Message{
		Key: []byte(strconv.Itoa(1)),
		// create an arbitrary message payload for the value
		Value: []byte(encoded),
	})
	if err != nil {
		log.Fatal("could not write message " + err.Error())
	}

	return nil
}
