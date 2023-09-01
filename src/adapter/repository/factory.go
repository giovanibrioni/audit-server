package repository

import (
	"context"

	"github.com/giovanibrioni/audit-server/audit"
	"go.uber.org/zap"
)

func Factory(ctx context.Context, logger *zap.SugaredLogger, storageType string) audit.AuditRepo {
	switch storageType {
	case "kafka":
		return NewKafkaAuditRepository(ctx, logger)
	case "redis":
		return NewRedisAuditRepository(ctx, logger)
	case "amqp":
		return NewAmqpAuditRepository(ctx, logger)
	case "postgres":
		return NewPostgresAuditRepository(ctx, logger)
	default:
		return NewStdoutAuditRepository(ctx, logger)
	}
}
