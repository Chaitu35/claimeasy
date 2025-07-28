package main

import (
	//"fmt"
	"net/http"

	"github.com/chaitu35/claimeasy/api/handler"
	"github.com/go-chi/chi/v5"
)

func main() {
	//fmt.Println("Api server  starting .........!")
	r:=chi.NewRouter()
	r.Get("/health", handler.HealthCheck)
	http.ListenAndServe(":8080", r)


}