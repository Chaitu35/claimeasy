package service

import (
	"context"

	"github.com/Chaitu35/claimeasy/backend/internal/auth/dto"
)

type AuthService interface {
	Login(ctx context.Context, request dto.LoginRequest) (*dto.LoginResponse, error)
    ForgetPassword(ctx context.Context, request dto.ForgetPasswordRequest) (*dto.ForgetPasswordResponse, error)
	ChangePassword(ctx context.Context, request dto.ChangePasswordRequest) (*dto.ChangePasswordResponse, error)
}