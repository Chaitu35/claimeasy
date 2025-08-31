package role

import "github.com/gin-gonic/gin"

func RegisterRoleRoutes(rg *gin.RouterGroup, h *Handler) {
	role := rg.Group("/roles")
	{
		role.POST("/", h.CreateRole)
		role.GET("/", h.ListRoles)
		role.GET("/:id", h.GetRoleByID)
		role.PUT("/:id", h.UpdateRole)
		role.PATCH("/:id/status", h.UpdateRoleStatus) // for status changes like active/inactive/deleted
	}
}
