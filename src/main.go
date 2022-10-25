package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giovanibrioni/audit-server/adapter/api"
	"github.com/giovanibrioni/audit-server/adapter/repository"
	"github.com/giovanibrioni/audit-server/helper"
)

var (
	storageType = helper.GetEnvOrDefault("STORAGE_TYPE", "stdout")
	serverPort  = helper.GetEnvOrDefault("SERVER_PORT", "8080")
)

func main() {

	auditRepo := repository.Factory(storageType)
	auditHandler := api.NewAuditHandler(auditRepo)
	router := initRouter(auditHandler)
	log.Printf("\nStorate Type setting to: %s\n", storageType)
	router.Run(":" + serverPort)

}

func initRouter(auditHandler *api.AuditHandler) *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.POST("/audit/batch", auditHandler.PostBatch)

	return r
}
