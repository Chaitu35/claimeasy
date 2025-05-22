package regulatory

type Regulatory string

const (
	MOH Regulatory = "moh"
	DHA Regulatory = "dha"
	DOH Regulatory = "doh"
)

func IsValidRegulatory(r string) bool {
	switch Regulatory(r) {
	case MOH, DHA, DOH:
		return true
	default:
		return false
	}
}
