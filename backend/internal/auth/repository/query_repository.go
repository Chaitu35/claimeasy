package repository

import (
	"context"

	"github.com/Chaitu35/claimeasy/backend/ent"
)


type AuthQueryRepository interface{

  GetUserByUsername(ctx context.Context,username string)(*ent.User,error)
  GetUserByEmail(ctx context.Context,email string)(*ent.User,error)
  GetUserByRefreshToken(ctx context.Context,refreshToken string)(*ent.User,error)
  VerifyUserbyResetToken(ctx context.Context,resetToken string)(*ent.User,error)
}
