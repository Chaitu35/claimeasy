package role

type Role string 

const (
	RoleAdmin Role = "admin"
	RoleCoder Role = "coder"
	RoleSupport  Role = "support"
	

)

func IsValidRole(role Role) bool {
	switch role {
	case RoleAdmin, RoleCoder, RoleSupport:
		return true
	default:
		return false
	}
}

