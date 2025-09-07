package notification

import "fmt"

func SendResetSMS(phone, token string) error {
	// TODO: wire with Twilio / AWS SNS / local SMS provider
	fmt.Printf("Sending password reset SMS to %s with token %s\n", phone, token)
	return nil
}
