package auth

import "github.com/gin-gonic/gin"

func RegisterAuthRoutes(rg *gin.RouterGroup, h *Handler) {
	authGroup := rg.Group("/auth")

	authGroup.POST("/login", h.Login)
	authGroup.POST("/forgot-password", h.ForgotPassword)
	authGroup.POST("/reset-password", h.ResetPassword)
	authGroup.POST("/verify-mfa", h.VerifyMFA)
	authGroup.POST("/enroll-totp", h.EnrollTOTP)
	authGroup.POST("/verify-totp", h.VerifyTOTP)
	authGroup.POST("/generate-recovery-codes", h.GenerateRecoveryCodes)
	authGroup.POST("/use-recovery-code", h.UseRecoveryCode)

	// NEW routes
	authGroup.POST("/refresh-token", h.RefreshToken)
	authGroup.POST("/logout", h.Logout)
}
