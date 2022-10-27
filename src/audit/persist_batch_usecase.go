package audit

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
)

type PersistBatchUseCase struct {
	Repo AuditRepo
}

func NewPersistBatchUseCase(auditRepo AuditRepo) *PersistBatchUseCase {
	return &PersistBatchUseCase{Repo: auditRepo}
}

func (p *PersistBatchUseCase) Execute(reqBody []byte) (uuid.UUID, error) {

	jobId := uuid.New()
	auditId := uuid.New()
	var rawMessages []map[string]any

	json.Unmarshal([]byte(reqBody), &rawMessages)

	auditLog := &AuditEntity{
		AuditId:     auditId,
		RawMessages: rawMessages,
	}

	err := p.Repo.Save(auditLog)
	if err != nil {
		log.Fatal(err)
	}

	return jobId, nil

}
