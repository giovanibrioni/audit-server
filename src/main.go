package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/giovanibrioni/audit-server/adapter"
	"github.com/giovanibrioni/audit-server/audit"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
)

func GetEnvOrDefault(envVar, defaultValue string) string {
	if v, ok := os.LookupEnv(envVar); ok && len(v) > 0 {
		return v
	}
	return defaultValue
}

func main() {

	dbURL := GetEnvOrDefault("REDIS_URL", "localhost:6379")
	redisPassword := GetEnvOrDefault("REDIS_PASSWORD", "")
	db := redisConnect(dbURL, redisPassword)
	defer db.Close()

	auditRepo := adapter.NewRedisAuditRepository(db)

	router := initRouter(auditRepo)
	router.Run(":" + GetEnvOrDefault("SERVER_PORT", "8080"))

}

func initRouter(repo audit.AuditRepo) *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/batch-audit", func(c *gin.Context) {
		reqBody, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}
		if err := c.Request.Body.Close(); err != nil {
			panic(err)
		}

		audit_id := uuid.New()
		var rawMessages []map[string]any

		json.Unmarshal([]byte(reqBody), &rawMessages)

		auditLog := &audit.AuditEntity{
			ID:          audit_id,
			RawMessages: rawMessages,
		}

		err = repo.PersistLogs(auditLog)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		c.JSON(http.StatusCreated, gin.H{
			"audit_id": audit_id,
		})
	})

	return r
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
