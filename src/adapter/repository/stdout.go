package repository

import (
	"encoding/json"
	"log"

	"github.com/giovanibrioni/audit-server/audit"
)

type stdoutAuditRepository struct{}

func NewStdoutAuditRepository() audit.AuditRepo {
	return &stdoutAuditRepository{}
}

func (r *stdoutAuditRepository) SaveBatch(auditLogs []*audit.AuditEntity) error {
	for _, auditLog := range auditLogs {
		encoded, err := json.Marshal(auditLog)
		if err != nil {
			log.Fatal("Unable to marshal auditLogs")
			return err
		}
		log.Print(string(encoded))
	}

	return nil
}
