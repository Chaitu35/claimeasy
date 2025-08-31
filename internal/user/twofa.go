package user

import (
	"database/sql/driver"
	"fmt"
)

// TwoFactorType represents allowed 2FA methods for a user
type TwoFactorType string

const (
	TwoFAAuthenticator TwoFactorType = "authenticator" // Google Authenticator / Authy
	TwoFASMS           TwoFactorType = "sms"           // OTP via SMS
	TwoFAEmail         TwoFactorType = "email"         // OTP via Email
	TwoFANone          TwoFactorType = ""              // No 2FA
)

// Value implements driver.Valuer → tells DB how to store enum
func (t TwoFactorType) Value() (driver.Value, error) {
	return string(t), nil
}

// Scan implements sql.Scanner → tells DB how to read enum
func (t *TwoFactorType) Scan(value interface{}) error {
	if value == nil {
		*t = TwoFANone
		return nil
	}
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("cannot scan TwoFactorType from %T", value)
	}
	*t = TwoFactorType(str)
	return nil
}
