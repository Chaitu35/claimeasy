package user

import (
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(rg *gin.RouterGroup, h *Handler) {
	user := rg.Group("/users")
	{
		user.POST("/", h.CreateUser)
		user.GET("/:id", h.GetUserByID)
		user.GET("/", h.SearchUsers)
		user.PUT("/:id", h.UpdateUser)
		user.PATCH("/:id/status", h.UpdateUserStatus)
	}
}

