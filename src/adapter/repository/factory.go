package repository

import (
	"github.com/giovanibrioni/audit-server/audit"
)

func Factory(storageType string) audit.AuditRepo {
	switch storageType {
	case "kafka":
		return NewKafkaAuditRepository()
	case "redis":
		return NewRedisAuditRepository()
	case "amqp":
		return NewAmqpAuditRepository()
	default:
		return NewStdoutAuditRepository()
	}
}
