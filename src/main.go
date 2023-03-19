package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/giovanibrioni/audit-server/adapter/api"
	"github.com/giovanibrioni/audit-server/adapter/repository"
	"github.com/giovanibrioni/audit-server/helper"
	"go.uber.org/zap"
)

var (
	logLevel    = helper.GetEnvOrDefault("LOGGER_LOG_LEVEL", "DEBUG")
	storageType = helper.GetEnvOrDefault("STORAGE_TYPE", "stdout")
	serverPort  = ":" + helper.GetEnvOrDefault("SERVER_PORT", "8080")
)

func main() {
	logger := initLogger()
	defer logger.Sync()

	logger.Info("Storage Type setting to: ", storageType)
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	auditRepo := repository.Factory(ctx, logger, storageType)
	auditHandler := api.NewAuditHandler(auditRepo)
	router := initRouter(auditHandler)
	srv := &http.Server{
		Addr:    serverPort,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to listen and serve: %v\n", err)
		}
	}()
	<-ctx.Done()

	//Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	logger.Info("Shuting down gracefully, press Ctrl+C again to force")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Info("Server exiting")

}

func initRouter(auditHandler *api.AuditHandler) *gin.Engine {
	if logLevel != "DEBUG" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.POST("/audit/batch", auditHandler.PostBatch)

	return r
}

func initLogger() *zap.SugaredLogger {
	var rootLogger *zap.Logger
	var logger *zap.SugaredLogger
	var err error
	if logLevel != "DEBUG" {
		rootLogger, err = zap.NewDevelopment()
	} else {
		rootLogger, err = zap.NewProduction()
	}
	if err != nil {
		log.Fatalf("Failed to initialize zap: %v", err)
	}
	logger = rootLogger.Sugar()
	logger = logger.With(zap.String("application", "audit-server"))
	return logger
}
