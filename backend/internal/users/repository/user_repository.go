package repository

import (
	"context"

	"github.com/Chaitu35/claimeasy/backend/ent"
)

type UserRepository interface{
	GetByUsername(ctx context.Context, username string) (ent.User, error)
    UpdatePassword(ctx context.Context, userID int, newPassword string) error		
}