package auth

import (
	"context"
	// "crypto/rand"
	// "encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/chaitu35/claimeasy/internal/audit"
	"github.com/chaitu35/claimeasy/internal/permission"
	"github.com/chaitu35/claimeasy/internal/user"
	security "github.com/chaitu35/claimeasy/pkg/security"
    utils "github.com/chaitu35/claimeasy/pkg/notifications"
	"gorm.io/gorm"
)

const (
    MaxLoginAttempts = 5
    LockoutDuration = 15 * time.Minute
    MFACodeExpiry = 5 * time.Minute
    MaxMFAAttempts = 3
)

// -------------------------
// Service Interface
// -------------------------
type Service interface {
    Login(ctx context.Context, db *gorm.DB, req LoginRequest) (*LoginResponse, error)
    StartMFAChallenge(ctx context.Context, db *gorm.DB, userID uint) error
    VerifyMFAChallenge(ctx context.Context, db *gorm.DB, userID uint, code string) (*LoginResponse, error)
    EnrollTOTP(ctx context.Context, db *gorm.DB, userID uint) (string, error)
    VerifyTOTPEnrollment(ctx context.Context, db *gorm.DB, userID uint, code string) error
    ResetMFA(ctx context.Context, db *gorm.DB, userID uint) error
    GenerateRecoveryCodes(ctx context.Context, db *gorm.DB, userID uint) ([]string, error)
    UseRecoveryCode(ctx context.Context, db *gorm.DB, userID uint, code string) (*LoginResponse, error)

     ForgotPassword(ctx context.Context, db *gorm.DB, email string) error
    ResetPassword(ctx context.Context, db *gorm.DB, token string, newPassword string) error

    Logout(ctx context.Context, db *gorm.DB, userID uint) error
}

// -------------------------
// Implementation
// -------------------------
type service struct {
    jwt               security.JWTService
    permissionService permission.Service
    auditLogger       *audit.Logger
     mailer            utils.Mailer
}

func NewService(jwt security.JWTService, permService permission.Service, auditor *audit.Logger,mailer utils.Mailer) Service {
    return &service{
        jwt:               jwt,
        permissionService: permService,
        auditLogger:       auditor,
        mailer:            mailer, // or use dependency injection
    }
}


// -------------------------
// Login with rate limiting and account lockout
// -------------------------
func (s *service) Login(ctx context.Context, db *gorm.DB, req LoginRequest) (*LoginResponse, error) {
    // Check for account lockout first
    var loginAttempt LoginAttempt
    err := db.WithContext(ctx).Where("identifier = ?", req.Email).
        Order("created_at desc").First(&loginAttempt).Error
    
    if err == nil && loginAttempt.IsLocked() {
        return nil, errors.New("account temporarily locked due to too many failed attempts")
    }

    var u user.User
    if err := db.WithContext(ctx).Where("email = ?", req.Email).First(&u).Error; err != nil {
        s.recordFailedAttempt(ctx, db, req.Email)

         if s.auditLogger != nil {
        _ = s.auditLogger.LogEvent(ctx, audit.ActionLogin, "User", 0, fmt.Sprintf("Failed login for email: %s", req.Email))
    }
        return nil, errors.New("invalid credentials")
    }

    if !security.CheckStringHash(req.Password, u.PasswordHash) {
        s.recordFailedAttempt(ctx, db, req.Email)
         if s.auditLogger != nil {
        _ = s.auditLogger.LogEvent(ctx, audit.ActionLogin, "User", 0, fmt.Sprintf("Failed login for email: %s", req.Email))
    }
        return nil, errors.New("invalid credentials")
    }

    // Clear failed attempts on successful password verification
    s.clearFailedAttempts(ctx, db, req.Email)

    // MFA required?
    if u.Need2FAuth {
        if err := s.StartMFAChallenge(ctx, db, u.ID); err != nil {
            log.Printf("Failed to start MFA challenge: %v", err)
              if s.auditLogger != nil {
        _ = s.auditLogger.LogEvent(ctx, audit.ActionLogin, "User", u.ID, "MFA challenge started")
    }
            return nil, errors.New("authentication service temporarily unavailable")
        }
        return &LoginResponse{
            MFARequired: true,
            UserID:      u.ID,
        }, nil
    }

    // Direct login if no MFA
    return s.generateTokenResponse(u,db)
}

// -------------------------
// Rate limiting helpers
// -------------------------
func (s *service) recordFailedAttempt(ctx context.Context, db *gorm.DB, identifier string) {
    attempt := LoginAttempt{
        Identifier: identifier,
        Success:    false,
        CreatedAt:  time.Now(),
    }
    db.WithContext(ctx).Create(&attempt)
}

func (s *service) clearFailedAttempts(ctx context.Context, db *gorm.DB, identifier string) {
    db.WithContext(ctx).Where("identifier = ?", identifier).Delete(&LoginAttempt{})
}

// -------------------------
// Start MFA Challenge with better security
// -------------------------
func (s *service) StartMFAChallenge(ctx context.Context, db *gorm.DB, userID uint) error {
    // Invalidate any existing challenges
    db.WithContext(ctx).Where("user_id = ? AND consumed = false", userID).
        Update("consumed", true)

    // Generate cryptographically secure random code
    code, err := security.GenerateNumericCode(6)
    if err != nil {
        return fmt.Errorf("failed to generate MFA code: %w", err)
    }

    hash, err := security.HashString(code)
    if err != nil {
        return fmt.Errorf("failed to hash MFA code: %w", err)
    }

    ch := MFAChallenge{
        UserID:        userID,
        Channel:       "sms", // TODO: choose dynamically based on user settings
        CodeHash:      hash,
        ExpiresAt:     time.Now().Add(MFACodeExpiry),
        MaxAttempts:   MaxMFAAttempts,
        FailedAttempts: 0,
        Consumed:      false,
    }

    if err := db.WithContext(ctx).Create(&ch).Error; err != nil {
        return fmt.Errorf("failed to save MFA challenge: %w", err)
    }

    // TODO: send `code` via SMS/email using proper service
    // Remove this debug print in production
    if os.Getenv("ENV") == "development" {
        fmt.Println("DEBUG OTP:", code)
    }

    return nil
}

// -------------------------
// Verify MFA Challenge with timing attack protection
// -------------------------
func (s *service) VerifyMFAChallenge(ctx context.Context, db *gorm.DB, userID uint, code string) (*LoginResponse, error) {
    var ch MFAChallenge
    if err := db.WithContext(ctx).
        Where("user_id = ? AND consumed = false AND expires_at > ?", userID, time.Now()).
        Order("created_at desc").
        First(&ch).Error; err != nil {
        return nil, errors.New("no active challenge found")
    }

    if ch.FailedAttempts >= ch.MaxAttempts {
        return nil, errors.New("too many failed attempts")
    }

    // Use constant-time comparison to prevent timing attacks
    if !security.CheckStringHash(code, ch.CodeHash) {
        // Increment failed attempts
        db.WithContext(ctx).Model(&ch).Update("failed_attempts", gorm.Expr("failed_attempts + 1"))
        return nil, errors.New("invalid verification code")
    }

    // Mark challenge as consumed
    if err := db.WithContext(ctx).Model(&ch).Update("consumed", true).Error; err != nil {
        return nil, errors.New("failed to process verification")
    }

    // Get user and issue tokens
    var u user.User
    if err := db.WithContext(ctx).First(&u, userID).Error; err != nil {
        return nil, errors.New("user not found")
    }

    return s.generateTokenResponse(u,db)
}

// -------------------------
// TOTP Enrollment with proper secret handling
// -------------------------
func (s *service) EnrollTOTP(ctx context.Context, db *gorm.DB, userID uint) (string, error) {
    secret, uri, err := security.GenerateTOTPSecret("ClaimEasy", fmt.Sprintf("user-%d", userID))
    if err != nil {
        return "", fmt.Errorf("failed to generate TOTP secret: %w", err)
    }

    aesKey := os.Getenv("APP_SECRET")
    if len(aesKey) != 32 {
        return "", errors.New("invalid encryption key configuration")
    }

    // Encrypt and store secret
    encryptedSecret, err := security.Encrypt(aesKey, secret)
    if err != nil {
        return "", fmt.Errorf("failed to encrypt secret: %w", err)
    }

    if err := db.WithContext(ctx).Model(&user.User{}).
        Where("id = ?", userID).
        Update("mfa_secret", encryptedSecret).Error; err != nil {
        return "", fmt.Errorf("failed to save TOTP secret: %w", err)
    }

    return uri, nil
}

// -------------------------
// Verify TOTP Enrollment
// -------------------------
func (s *service) VerifyTOTPEnrollment(ctx context.Context, db *gorm.DB, userID uint, code string) error {
    var u user.User
    if err := db.WithContext(ctx).First(&u, userID).Error; err != nil {
        return errors.New("user not found")
    }

    if u.TwoFactorSecret == "" {
        return errors.New("no TOTP secret found - please enroll first")
    }

    aesKey := os.Getenv("APP_SECRET")
    if len(aesKey) != 32 {
        return errors.New("invalid encryption key configuration")
    }

    secret, err := security.Decrypt(aesKey, u.TwoFactorSecret)
    if err != nil {
        return errors.New("failed to decrypt TOTP secret")
    }

    if !security.ValidateTOTP(secret, code) {
        return errors.New("invalid TOTP code")
    }

    return db.WithContext(ctx).Model(&u).Update("need2f_auth", true).Error
}

// -------------------------
// Reset MFA with proper cleanup
// -------------------------
func (s *service) ResetMFA(ctx context.Context, db *gorm.DB, userID uint) error {
    var u user.User
    if err := db.WithContext(ctx).First(&u, userID).Error; err != nil {
        return errors.New("user not found")
    }

    // Use transaction for atomic operations
    return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        // Reset user MFA settings
        if err := tx.Model(&u).Updates(map[string]interface{}{
            "need2f_auth":    false,
            "two_fa_method":  "",  // Using two_fa_method field
        }).Error; err != nil {
            return err
        }

        // Delete active challenges
        if err := tx.Where("user_id = ?", userID).Delete(&MFAChallenge{}).Error; err != nil {
            return err
        }

        // Mark all recovery codes as used
        if err := tx.Model(&RecoveryCode{}).
            Where("user_id = ? AND used = false", userID).
            Updates(map[string]interface{}{
                "used":    true,
                "used_at": time.Now(),
            }).Error; err != nil {
            return err
        }

        return nil
    })
}

// -------------------------
// Recovery Codes with better entropy
// -------------------------
func (s *service) GenerateRecoveryCodes(ctx context.Context, db *gorm.DB, userID uint) ([]string, error) {
    codes, err := security.GenerateRecoveryCodes(8) // 8 codes, each 12 characters
    if err != nil {
        return nil, fmt.Errorf("failed to generate recovery codes: %w", err)
    }

    var toSave []RecoveryCode
    for _, c := range codes {
        hash, err := security.HashString(c)
        if err != nil {
            return nil, fmt.Errorf("failed to hash recovery code: %w", err)
        }

        toSave = append(toSave, RecoveryCode{
            UserID:   userID,
            CodeHash: hash,
            Used:     false,
        })
    }

    if err := db.WithContext(ctx).Create(&toSave).Error; err != nil {
        return nil, fmt.Errorf("failed to save recovery codes: %w", err)
    }

    return codes, nil
}

// -------------------------
// Use Recovery Code with proper validation
// -------------------------
func (s *service) UseRecoveryCode(ctx context.Context, db *gorm.DB, userID uint, code string) (*LoginResponse, error) {
    var recoveryCodes []RecoveryCode
    if err := db.WithContext(ctx).
        Where("user_id = ? AND used = false", userID).
        Find(&recoveryCodes).Error; err != nil {
        return nil, errors.New("no recovery codes found")
    }

    var validCode *RecoveryCode
    for i := range recoveryCodes {
        if security.CheckStringHash(code, recoveryCodes[i].CodeHash) {
            validCode = &recoveryCodes[i]
            break
        }
    }

    if validCode == nil {
        return nil, errors.New("invalid recovery code")
    }

    // Mark recovery code as used
    now := time.Now()
    if err := db.WithContext(ctx).Model(validCode).Updates(map[string]interface{}{
        "used":    true,
        "used_at": &now,
    }).Error; err != nil {
        return nil, errors.New("failed to mark recovery code as used")
    }

    // Get user and issue tokens
    var u user.User
    if err := db.WithContext(ctx).First(&u, userID).Error; err != nil {
        return nil, errors.New("user not found")
    }

    return s.generateTokenResponse(u,db)
}

// -------------------------
// Helper methods
// -------------------------
func (s *service) generateTokenResponse(u user.User,db *gorm.DB) (*LoginResponse, error) {
    // Fetch user's permissions
    userPerms, err := s.permissionService.ListUserPermissions(context.Background(), db, u.ID)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch user permissions: %w", err)
    }

    // Extract permission names as []string
   
   var perms []string
   for _, p := range userPerms {
    perms = append(perms, fmt.Sprintf("%d", p.PermissionID))
                  }

    accessToken, err := s.jwt.GenerateAccessToken(u.ID, u.Email, perms)
    if err != nil {
        return nil, fmt.Errorf("failed to generate access token: %w", err)
    }

    refreshToken, err := s.jwt.GenerateRefreshToken(u.ID, u.Email)
    if err != nil {
        return nil, fmt.Errorf("failed to generate refresh token: %w", err)
    }


        if s.auditLogger != nil {
        _ = s.auditLogger.LogEvent(context.Background(), audit.ActionLogin, "User", u.ID, "Successful login")
    }
    return &LoginResponse{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        UserID:       u.ID,
        
    }, nil
}
func (s *service) Logout(ctx context.Context, db *gorm.DB, userID uint) error {
    // Record audit
    if s.auditLogger != nil {
        _ = s.auditLogger.LogEvent(ctx, audit.ActionLogout, "User", userID, "User logged out")
    }

    // If you store refresh tokens in DB, you can revoke them
    if err := db.WithContext(ctx).Where("user_id = ?", userID).Delete(&user.PasswordResetToken{}).Error; err != nil {
        return fmt.Errorf("failed to revoke refresh tokens: %w", err)
    }

    return nil
}


// Note: Using your existing security functions from password.go and totp.go:
// - security.GenerateNumericCode(6) for MFA codes
// - security.HashString(string) for hashing 
// - security.CheckStringHash(input, hash) for verification
// - security.GenerateRecoveryCodes(n) for recovery codes  
// - security.GenerateTOTPSecret(accountName, issuer) for TOTP setup
// - security.ValidateTOTP(secret, code) for TOTP verification

// -------------------------
// Models (add these to your models file)
// -------------------------

func (la *LoginAttempt) IsLocked() bool {
    cutoff := time.Now().Add(-LockoutDuration)
    // Example: mark locked if last attempt failed and is within lockout window
    return !la.Success && la.CreatedAt.After(cutoff)
}
func (s *service) ForgotPassword(ctx context.Context, db *gorm.DB, email string) error {
	var u user.User
	if err := db.WithContext(ctx).Where("email = ?", email).First(&u).Error; err != nil {
		return errors.New("user not found")
	}

	token, err := security.GenerateSecureToken(32)
	if err != nil {
		return fmt.Errorf("failed to generate reset token: %w", err)
	}

	reset := PasswordReset{
		UserID:    u.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(30 * time.Minute),
	}

	if err := db.WithContext(ctx).Create(&reset).Error; err != nil {
		return fmt.Errorf("failed to save reset token: %w", err)
	}

	// 🔹 Send email using injected mailer (SMTP/SendGrid/etc.)
	if err := s.mailer.SendResetEmail(u.Email, token); err != nil {
		return fmt.Errorf("failed to send reset email: %w", err)
	}

	// Audit log
	if s.auditLogger != nil {
		_ = s.auditLogger.LogEvent(ctx, audit.ActionResetPassword, "User", u.ID, "Password reset initiated")
	}
	return nil
}

func (s *service) ResetPassword(ctx context.Context, db *gorm.DB, token, newPassword string) error {
	var reset PasswordReset
	if err := db.WithContext(ctx).
		Where("token = ? AND expires_at > ?", token, time.Now()).
		First(&reset).Error; err != nil {
		return errors.New("invalid or expired token")
	}

	var u user.User
	if err := db.WithContext(ctx).First(&u, reset.UserID).Error; err != nil {
		return errors.New("user not found")
	}

	hash, err := security.HashString(newPassword)
	if err != nil {
		return errors.New("failed to hash new password")
	}

	if err := db.WithContext(ctx).Model(&u).Update("password_hash", hash).Error; err != nil {
        
		return errors.New("failed to update password")
        
        
	}

	_ = db.WithContext(ctx).Delete(&reset) // clean up token
	return nil
}

