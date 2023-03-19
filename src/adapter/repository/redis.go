package repository

import (
	"context"
	"crypto/tls"
	"strconv"

	"github.com/giovanibrioni/audit-server/audit"
	"github.com/giovanibrioni/audit-server/helper"
	"github.com/goccy/go-json"
	"go.uber.org/zap"

	"github.com/redis/go-redis/v9"
)

var redisKey = helper.GetEnvOrDefault("REDIS_KEY", "audit_logs")

type redisAuditRepository struct {
	client *redis.Client
	ctx    context.Context
	logger *zap.SugaredLogger
}

func NewRedisAuditRepository(ctx context.Context, logger *zap.SugaredLogger) audit.AuditRepo {
	dbURL := helper.GetEnvOrDefault("REDIS_URL", "localhost:6379")
	redisPassword := helper.GetEnvOrDefault("REDIS_PASSWORD", "")
	redisDB := helper.GetEnvOrDefault("REDIS_DB", "0")
	tls := helper.GetEnvOrDefault("REDIS_ENABLE_TLS", "false")
	rcli := redisClient(ctx, dbURL, redisPassword, redisDB, tls, logger)
	return &redisAuditRepository{
		client: rcli,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *redisAuditRepository) SaveBatch(auditLogs []*audit.AuditEntity) error {
	var s []string
	for _, auditLog := range auditLogs {
		encoded, err := json.Marshal(auditLog)
		if err != nil {
			r.logger.Fatal("Unable to marshal auditLogs")
			return err
		}
		s = append(s, string(encoded))
	}
	err := r.client.RPush(r.ctx, redisKey, s).Err()
	if err != nil {
		r.logger.Fatal(err)
	}
	r.logger.Infof("AuditLog with jobId: %s, inserted on redisKey: %s", auditLogs[0].JobId, redisKey)
	return nil
}

func redisClient(ctx context.Context, url string, password string, db string, enableTLS string, logger *zap.SugaredLogger) *redis.Client {
	logger.Infof("Trying to connect on redis database url: %s, db: %s", url, db)
	dbInt, err := strconv.Atoi(db)
	if err != nil {
		logger.Fatal(err)
	}
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	if enableTLS == "false" {
		tlsConfig = nil
	}
	client := redis.NewClient(&redis.Options{
		TLSConfig: tlsConfig,
		Addr:      url,
		Password:  password,
		DB:        dbInt,
	})
	err = client.Ping(ctx).Err()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Infof("Connected to redis database url: %s, db: %s", url, db)
	return client

}
