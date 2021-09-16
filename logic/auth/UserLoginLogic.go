package auth

import (
	"context"
	"gitlab.com/simateb-project/simateb-backend/domain/auth"
)

type UserLoginInterface interface {
	UserLogin(phoneNumber string, password string) (*auth.ResponseAccessToken, error)
}

type UserLoginLogic struct {
	Context context.Context
	Self UserLoginInterface
}

func (login UserLoginLogic) UserLogin(phoneNumber string, password string) (*auth.ResponseAccessToken, error)  {
	return nil, nil
}