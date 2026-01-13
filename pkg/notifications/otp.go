package notification

import (
	"fmt"
	//"time"

	"github.com/pquerna/otp/totp"
)

// Generate secret for user during enrollment
func GenerateTOTPSecret(accountName, issuer string) (string, string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: accountName,
	})
	if err != nil {
		return "", "", err
	}
	return key.Secret(), key.URL(), nil // Secret (store in DB), QR URL (show to user)
}

// Validate a code entered by user
func ValidateTOTP(secret, code string) bool {
	return totp.Validate(code, secret)
}

// Example send (not real, just instruct user)
func SendTOTPInstruction(accountName, secret string) error {
	fmt.Printf("Provide this secret/QR to %s for authenticator setup: %s\n", accountName, secret)
	return nil
}
