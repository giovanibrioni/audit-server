package api

import (
	"encoding/json"
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
		return
	}
	if err := c.Request.Body.Close(); err != nil {
		panic(err)
	}
	var rawMessages []map[string]any
	err = json.Unmarshal([]byte(reqBody), &rawMessages)
	if err != nil || len(rawMessages) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Request body should be a list",
		})
		return
	}
	jobId, err := audit.NewPersistBatchUseCase(a.repo).Execute(rawMessages)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"job_id": jobId,
	})
}
