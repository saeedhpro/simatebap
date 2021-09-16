package auth

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"gitlab.com/simateb-project/simateb-backend/config"
	"gitlab.com/simateb-project/simateb-backend/constant"
	"time"
)

func ValidateToken(tokenString string) (*UserClaims, error){
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JwtSecret), nil
	})
	if err != nil {
		return nil, errors.New(constant.InvalidTokenError)
	}
	tokenClaims := token.Claims.(*UserClaims)
	if tokenClaims.ExpiresAt < time.Now().Unix() {
		return nil, errors.New(constant.TokenIsExpiredError)
	}
	return tokenClaims, nil
}
