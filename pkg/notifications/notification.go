package notification

type Channel string

const (
	Email Channel = "email"
	SMS   Channel = "sms"
	TOTP  Channel = "totp"
)

func Send(channel Channel, recipient, token string) error {
	switch channel {
	case Email:
		return SendResetEmail(recipient, token)
	case SMS:
		return SendResetSMS(recipient, token)
	case TOTP:
		return SendTOTPInstruction(recipient, token)
	default:
		return nil
	}
}
