package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RouteRegister interface {
	RegisterRoutes(router *gin.Engine, db *gorm.DB)
}