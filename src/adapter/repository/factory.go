package repository

import (
	"github.com/giovanibrioni/audit-server/audit"
)

func Factory(storageType string) audit.AuditRepo {
	switch storageType {
	case "redis":
		return NewRedisAuditRepository()
	default:
		return NewStdoutAuditRepository()
	}
}
