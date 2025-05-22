package repository

import (
	"context"

	"github.com/Chaitu35/claimeasy/backend/ent"
	"github.com/Chaitu35/claimeasy/backend/ent/user"
)

type userEntRepository struct {

client *ent.Client

}

//constructor
func NewUserEntRepository(client *ent.Client)*userEntRepository {
	return &userEntRepository{
		client: client,
	}
}


func (r *userEntRepository) GetByUsername(ctx context.Context,username  string ) (*ent.User,error){
	return r.client.User.Query().Where(user.Username(username)).Only(ctx)	
}
func (r *userEntRepository) UpdatePassword(ctx context.Context,userId int,password string)error{
	_,err := r.client.User.UpdateOneID(userId).SetPasswordHash(password).Save(ctx)
    
	return err
}