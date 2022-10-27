package api

import (
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/giovanibrioni/audit-server/audit"
)

type AuditHandler struct {
	repo audit.AuditRepo
}

func NewAuditHandler(repo audit.AuditRepo) *AuditHandler {
	return &AuditHandler{
		repo,
	}
}

func (a *AuditHandler) PostBatch(c *gin.Context) {
	reqBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	if err := c.Request.Body.Close(); err != nil {
		panic(err)
	}
	jobId, err := audit.NewPersistBatchUseCase(a.repo).Execute(reqBody)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	c.JSON(http.StatusCreated, gin.H{
		"job_id": jobId,
	})
}
