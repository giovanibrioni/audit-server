package adapter

import (
	"encoding/json"
	"log"

	"github.com/giovanibrioni/audit-server/audit"
	"github.com/go-redis/redis"
)

const redisKey = "audit_logs"

type auditRepository struct {
	connection *redis.Client
}

func NewRedisAuditRepository(connection *redis.Client) audit.AuditRepo {
	return &auditRepository{
		connection,
	}
}

func (r *auditRepository) PersistLogs(auditLog *audit.AuditEntity) error {

	encoded, err := json.Marshal(auditLog)
	if err != nil {
		log.Fatal("Unable to marshal auditLogs")
		return err
	}
	err = r.connection.RPush(redisKey, encoded).Err()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
