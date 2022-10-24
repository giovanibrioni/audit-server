package adapter

import (
	"encoding/json"
	"log"

	"github.com/giovanibrioni/audit-server/audit"
)

type stdoutAuditRepository struct{}

func NewStdoutAuditRepository() audit.AuditRepo {
	return &stdoutAuditRepository{}
}

func (r *stdoutAuditRepository) PersistLogs(auditLog *audit.AuditEntity) error {

	encoded, err := json.Marshal(auditLog)
	if err != nil {
		log.Fatal("Unable to marshal auditLogs")
		return err
	}
	log.Print(string(encoded))

	return nil
}