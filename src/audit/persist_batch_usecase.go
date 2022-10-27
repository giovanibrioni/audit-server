package audit

import (
	"log"

	"github.com/google/uuid"
)

type PersistBatchUseCase struct {
	Repo AuditRepo
}

func NewPersistBatchUseCase(auditRepo AuditRepo) *PersistBatchUseCase {
	return &PersistBatchUseCase{Repo: auditRepo}
}

func (p *PersistBatchUseCase) Execute(rawMessages []map[string]any) (uuid.UUID, error) {
	var auditLogs []*AuditEntity
	jobId := uuid.New()
	for _, v := range rawMessages {
		auditId := uuid.New()
		auditLog := &AuditEntity{
			JobId:      jobId,
			AuditId:    auditId,
			RawMessage: v,
		}
		auditLogs = append(auditLogs, auditLog)
	}

	err := p.Repo.SaveBatch(auditLogs)
	if err != nil {
		log.Fatal(err)
	}

	return jobId, nil

}
