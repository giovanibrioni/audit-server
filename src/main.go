package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giovanibrioni/audit-server/adapter"
	"github.com/giovanibrioni/audit-server/audit"
	"github.com/giovanibrioni/audit-server/helper"
	"github.com/google/uuid"
)

func main() {

	auditRepo := adapter.StorageFactory(helper.GetEnvOrDefault("STORAGE_TYPE", "stdout"))

	router := initRouter(auditRepo)
	router.Run(":" + helper.GetEnvOrDefault("SERVER_PORT", "8080"))

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
