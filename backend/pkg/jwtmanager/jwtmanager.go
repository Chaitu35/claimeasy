package jwtmanager


import("github.com/golang-jwt/jwt/v5"
	//"time"
	)


type PermissionClaim struct{
	View bool `json:"view"`
	Edit bool `json:"edit"`
	Delete bool `json:"delete"`
	Create bool `json:"create"`
	Print bool `json:"print"`
	Export bool `json:"export"`

}



type UserClaim struct{
	jwt.RegisteredClaims
	UserId int `json:"user_id"`	
	Username string `json:"username"`
	Role string `json:"role"`	
	Permissions map[string]PermissionClaim `json:"permissions"`
}
