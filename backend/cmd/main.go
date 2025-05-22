package main

import (
	"github.com/joho/godotenv"

	"github.com/Chaitu35/claimeasy/backend/internal/server"
	"github.com/Chaitu35/claimeasy/backend/pkg/db"
)




func main(){
	env:=godotenv.Load(".env.dev")
	if env!=nil{	
		panic("Error loading .env file")
	}
	if err:=godotenv.Load(".env.dev"); err!=nil{
		panic("Error loading .env file")
	}
	client:=db.NewClient()
	defer client.Close()

	svr:=server.New(client)

	svr.Start(":8080")
}