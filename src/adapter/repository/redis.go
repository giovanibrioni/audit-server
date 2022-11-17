package repository

import (
	"encoding/json"
	"log"
	"strconv"

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
	redisDB := helper.GetEnvOrDefault("REDIS_DB", "0")
	rcli := redisClient(dbURL, redisPassword, redisDB)
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

func redisClient(url string, password string, db string) *redis.Client {
	dbInt, err := strconv.Atoi(db)
	if err != nil {
		log.Fatal(err)
	}
	client := redis.NewClient(&redis.Options{
		Addr:     url,
		Password: password,
		DB:       dbInt,
	})
	err = client.Ping().Err()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Connected to redis database url: %s, db: %s", url, db)
	return client

}
