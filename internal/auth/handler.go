package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	service Service
	db      *gorm.DB
}

func NewHandler(s Service, db *gorm.DB) *Handler {
	return &Handler{service: s, db: db}
}

// -------------------------
// Auth Handlers
// -------------------------

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.Login(c.Request.Context(), h.db, req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Add Refresh Token validation using service if you want
	c.JSON(http.StatusOK, gin.H{"message": "refresh not implemented"})
}

func (h *Handler) VerifyMFA(c *gin.Context) {
	var req VerifyMFARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.VerifyMFAChallenge(c.Request.Context(), h.db, req.UserID, req.Code)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) EnrollTOTP(c *gin.Context) {
	userID := c.GetUint("user_id") // assumes middleware sets user_id in context
	uri, err := h.service.EnrollTOTP(c.Request.Context(), h.db, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"qr_uri": uri})
}

func (h *Handler) VerifyTOTP(c *gin.Context) {
	userID := c.GetUint("user_id")

	var body struct {
		Code string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.VerifyTOTPEnrollment(c.Request.Context(), h.db, userID, body.Code); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "TOTP verified"})
}

func (h *Handler) ResetMFA(c *gin.Context) {
	userID := c.GetUint("user_id")
	if err := h.service.ResetMFA(c.Request.Context(), h.db, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "MFA reset successfully"})
}

func (h *Handler) GenerateRecoveryCodes(c *gin.Context) {
	userID := c.GetUint("user_id")
	codes, err := h.service.GenerateRecoveryCodes(c.Request.Context(), h.db, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"codes": codes})
}

func (h *Handler) UseRecoveryCode(c *gin.Context) {
	userID := c.GetUint("user_id")

	var body struct {
		Code string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.UseRecoveryCode(c.Request.Context(), h.db, userID, body.Code)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest				
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.ForgotPassword(c.Request.Context(), h.db, req.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}	

	c.JSON(http.StatusOK, gin.H{"message": "If the email exists, a reset link has been sent"})
}
func (h *Handler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest



	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}	
	if err := h.service.ResetPassword(c.Request.Context(), h.db, req.Token, req.NewPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}	
	c.JSON(http.StatusOK, gin.H{"message": "Password has been reset successfully"})
}