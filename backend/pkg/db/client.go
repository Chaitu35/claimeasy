package db

import (
	"os"
	"context"

	"github.com/Chaitu35/claimeasy/backend/ent"
	"github.com/joho/godotenv"
)

func NewClient() *ent.Client{
	if err:=godotenv.Load(".env.dev"); err!=nil{
		panic("Error loading .env file")
	}
	dsn:=os.Getenv("DATABASE_URL")
	if dsn==""{
		panic("DATABASE_URL is not set")
	}
	client,error:=ent.Open("postgres",dsn)
	if error!=nil{
		panic("failed opening connection to postgres")
	}	
  if err := client.Schema.Create(context.Background()); err != nil {
		panic("failed creating schema resources: " + err.Error())
  }	

		return client	


}