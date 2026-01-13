package auth

import (
	"net/http"
	//"time"

	"github.com/chaitu35/claimeasy/internal/user"
	//security "github.com/chaitu35/claimeasy/pkg/security"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	service Service
 userService *user.Service
	db      *gorm.DB
}

func NewHandler(s Service, u *user.Service,db *gorm.DB) *Handler {
	return &Handler{service: s,userService: u, db: db}
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

    accessToken, newRefreshToken, err := h.service.RefreshToken(c.Request.Context(), h.db, req.RefreshToken)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "access_token":  accessToken,
        "refresh_token": newRefreshToken,
    })
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

func (h *Handler) Logout(c *gin.Context) {
	// If using a token blacklist, you would add the refresh token here
	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}