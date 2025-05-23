package repository

import (
	"context"
	"time"
)

type AuthCommandRepository interface {
  SaveRefreshToken(ctx context.Context, userId int, refreshToken string,expiresAt time.Time) error
  DeleteRefreshToken(ctx context.Context, token string) error
  UpdatePassword(ctx context.Context, userId int, password string) error
  CreateResetToken(ctx context.Context, userId int, resetToken string, expiresAt time.Time) error
  DeleteResetToken(ctx context.Context, token string) error
}