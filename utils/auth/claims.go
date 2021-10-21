package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gitlab.com/simateb-project/simateb-backend/config"
	"gitlab.com/simateb-project/simateb-backend/constant"
	"gitlab.com/simateb-project/simateb-backend/domain/wallet"
	"time"
)

type UserClaims struct {
	jwt.StandardClaims
	UserID         int64                `json:"user_id"`
	OrganizationID int64                `json:"organization_id"`
	FirstName      string               `json:"fname"`
	LastName       string               `json:"lname"`
	Tel            string               `json:"tel"`
	UserGroupID    int64                `json:"user_group_id"`
	Wallet         *wallet.WalletStruct `json:"wallet"`
}

func (u *UserClaims) GenerateToken() (*string, error) {
	u.ExpiresAt = time.Now().Unix() + constant.ExpTime
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, u)
	tokenString, err := token.SignedString([]byte(config.JwtSecret))
	if err != nil {
		return nil, err
	}
	return &tokenString, nil
}

func GetStaffUser(c *gin.Context) *UserClaims {
	resp, exist := c.Get("claims")
	if !exist {
		return nil
	}
	c2 := resp.(UserClaims)
	c2.Wallet = wallet.GetWallet(c2.UserID, "user")
	return &c2
}
