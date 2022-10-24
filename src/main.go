package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giovanibrioni/audit-server/adapter/api"
	"github.com/giovanibrioni/audit-server/adapter/repository"

	"github.com/giovanibrioni/audit-server/helper"
)

func main() {

	auditRepo := repository.Factory(helper.GetEnvOrDefault("STORAGE_TYPE", "stdout"))
	auditHandler := api.NewAuditHandler(auditRepo)
	router := initRouter(auditHandler)
	router.Run(":" + helper.GetEnvOrDefault("SERVER_PORT", "8080"))

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
