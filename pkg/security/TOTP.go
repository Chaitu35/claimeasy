package security

import (
	"github.com/pquerna/otp/totp"
)

// GenerateTOTPSecret generates a new TOTP secret and provisioning URL (for QR code).
func GenerateTOTPSecret(accountName, issuer string) (string, string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: accountName,
	})
	if err != nil {
		return "", "", err
	}
	return key.Secret(), key.URL(), nil
}

// ValidateTOTP validates the user-provided OTP code against the secret.
func ValidateTOTP(secret, code string) bool {
	return totp.Validate(code, secret)
}
