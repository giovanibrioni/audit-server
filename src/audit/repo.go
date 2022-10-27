package audit

import "github.com/google/uuid"

type AuditEntity struct {
	JobId      uuid.UUID      `json:"job_id" binding:"required"`
	AuditId    uuid.UUID      `json:"audit_id" binding:"required"`
	RawMessage map[string]any `json:"raw_messages" binding:"required"`
}

type AuditRepo interface {
	SaveBatch(auditEntity []*AuditEntity) error
}
