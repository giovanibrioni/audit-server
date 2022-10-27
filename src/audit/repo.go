package audit

import "github.com/google/uuid"

type AuditEntity struct {
	AuditId     uuid.UUID        `json:"audit_id" binding:"required"`
	RawMessages []map[string]any `json:"raw_messages" binding:"required"`
}

type AuditRepo interface {
	Save(auditEntity *AuditEntity) error
}
