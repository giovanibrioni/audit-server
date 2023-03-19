package repository

import (
	"context"

	"github.com/goccy/go-json"

	"github.com/giovanibrioni/audit-server/audit"
	"go.uber.org/zap"
)

type stdoutAuditRepository struct {
	ctx    context.Context
	logger *zap.SugaredLogger
}

func NewStdoutAuditRepository(ctx context.Context, logger *zap.SugaredLogger) audit.AuditRepo {
	return &stdoutAuditRepository{
		ctx:    ctx,
		logger: logger,
	}
}

func (s *stdoutAuditRepository) SaveBatch(auditLogs []*audit.AuditEntity) error {
	for _, auditLog := range auditLogs {
		encoded, err := json.Marshal(auditLog)
		if err != nil {
			s.logger.Fatal("Unable to marshal auditLogs")
			return err
		}
		s.logger.Info(string(encoded))
	}

	return nil
}
