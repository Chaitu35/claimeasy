package permission

import "github.com/gin-gonic/gin"

func RegisterPermissionRoutes(rg *gin.RouterGroup, h *Handler) {
	perm := rg.Group("/permissions")
	{
		perm.POST("/", h.CreatePermission)
		perm.PUT("/:id", h.UpdatePermission)
		perm.GET("/", h.ListPermissions)
		perm.GET("/:id", h.GetPermissionByID)

		perm.POST("/roles", h.AssignPermissionToRole)
		perm.DELETE("/roles/:roleID/:permID", h.RevokePermissionFromRole)
		perm.GET("/roles/:roleID", h.ListRolePermissions)

		perm.POST("/users", h.AssignPermissionToUser)
		perm.DELETE("/users/:userID/:permID", h.RevokePermissionFromUser)
		perm.GET("/users/:userID", h.ListUserPermissions)
	}
}
