package types

type Role string
const(
	Admin Role = "admin"
	Manager Role = "manager"
	Coder Role = "coder"
	Support Role = "support"

)  

func AllRoles() []Role {
	return []Role{Admin, Manager, Coder, Support}
}
