package audit

import (
	"encoding/base64"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Handler handles audit-related endpoints
type Handler struct {
	logger *Logger
}

func NewHandler(logger *Logger) *Handler {
	return &Handler{logger: logger}
}

func (h *Handler) GetAuditLogs(c *gin.Context) {
	entity := c.Param("entity")
	entityIDStr := c.Param("id")
	entityID, err := strconv.Atoi(entityIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid entity id"})
		return
	}

	var filters AuditFilters

	if actorIDStr := c.Query("actor_id"); actorIDStr != "" {
		if actorID, err := strconv.Atoi(actorIDStr); err == nil {
			id := uint(actorID)
			filters.ActorID = &id
		}
	}
	filters.Action = c.Query("action")
	if fromStr := c.Query("from"); fromStr != "" {
		if from, err := time.Parse(time.RFC3339, fromStr); err == nil {
			filters.From = from
		}
	}
	if toStr := c.Query("to"); toStr != "" {
		if to, err := time.Parse(time.RFC3339, toStr); err == nil {
			filters.To = to
		}
	}
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			filters.Limit = limit
		}
	}
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			filters.Offset = offset
		}
	}

	// Fetch logs
	logs, err := h.logger.storage.List(c.Request.Context(), entity, uint(entityID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Decode FullData safely
	for i := range logs {
		if len(logs[i].FullData) > 0 {
			decoded, err := base64.StdEncoding.DecodeString(string(logs[i].FullData))
			if err == nil {
				logs[i].FullDataStr = string(decoded) // new string field
			} else {
				logs[i].FullDataStr = string(logs[i].FullData) // fallback
			}
		}
	}

	c.JSON(http.StatusOK, logs)
}
