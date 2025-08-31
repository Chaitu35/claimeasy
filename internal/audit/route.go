package audit

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup, h *Handler) {
	auditGroup := r.Group("/audit")
	{
		auditGroup.GET("/:entity/:id", h.GetAuditLogs)
	}
}
