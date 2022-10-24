package adapter

import (
	"github.com/giovanibrioni/audit-server/audit"
)

func StorageFactory(storageType string) audit.AuditRepo {
	switch storageType {
	case "redis":
		return NewRedisAuditRepository()
	default:
		return NewStdoutAuditRepository()
	}
}
