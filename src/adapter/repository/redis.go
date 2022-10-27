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
	client *redis.Client
}

func NewRedisAuditRepository() audit.AuditRepo {
	dbURL := helper.GetEnvOrDefault("REDIS_URL", "localhost:6379")
	redisPassword := helper.GetEnvOrDefault("REDIS_PASSWORD", "")
	rcli := redisClient(dbURL, redisPassword)
	return &redisAuditRepository{
		client: rcli,
	}
}

func (r *redisAuditRepository) SaveBatch(auditLogs []*audit.AuditEntity) error {
	var s []string
	for _, auditLog := range auditLogs {
		encoded, err := json.Marshal(auditLog)
		if err != nil {
			log.Fatal("Unable to marshal auditLogs")
			return err
		}
		s = append(s, string(encoded))
	}
	err := r.client.RPush(redisKey, s).Err()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func redisClient(url string, password string) *redis.Client {
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
