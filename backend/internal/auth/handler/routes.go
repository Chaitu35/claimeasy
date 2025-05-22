package auth


import (
	"net/http"
    "github.com/go-chi/chi/v5"
	"github.com/Chaitu35/claimeasy/backend/ent"


)

type Handler struct {
	client *ent.Client
}

func Routes(client *ent.Client) http.Handler{
	r:= chi.NewRouter()
    h:= &Handler{
		client: client,
	  }
	r.Post("/login", h.LoginHandler)
	r.Post("/forget-password", h.ForgetPasswordHandler)
	r.Post("/change-password", h.ChangePasswordHandler)	
return r
}