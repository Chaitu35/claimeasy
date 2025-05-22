package auth

import (
	"net/http"
	//"encoding/json")
)

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	w.Write(([]byte("Login Successful")))
	// var login}Request LoginRequest
}
func (h *Handler) ForgetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	w.Write(([]byte("Password Reset Link Sent")))
	// var resetRequest ResetRequest
}
func (h *Handler) ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	w.Write(([]byte("Password Changed Successfully")))
	// var changeRequest ChangeRequest
}	