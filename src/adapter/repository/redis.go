package repository

import (
	"encoding/json"
	"log"

	"github.com/giovanibrioni/audit-server/audit"
	"github.com/giovanibrioni/audit-server/helper"

	"github.com/go-redis/redis"
)

var redisKey = helper.GetEnvOrDefault("REDIS_KEY", "audit_logs")

type redisAuditRepository struct {
	connection *redis.Client
}

func NewRedisAuditRepository() audit.AuditRepo {
	dbURL := helper.GetEnvOrDefault("REDIS_URL", "localhost:6379")
	redisPassword := helper.GetEnvOrDefault("REDIS_PASSWORD", "")
	rconn := redisConnect(dbURL, redisPassword)
	return &redisAuditRepository{
		connection: rconn,
	}
}

func (r *redisAuditRepository) Save(auditLog *audit.AuditEntity) error {
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

func redisConnect(url string, password string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     url,
		Password: password,
		DB:       0,
	})
	err := client.Ping().Err()

	if err != nil {
		log.Fatal(err)
	}
	return client

}
