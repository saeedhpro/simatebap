package authController

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/simateb-project/simateb-backend/constant"
	"gitlab.com/simateb-project/simateb-backend/domain/auth"
	"gitlab.com/simateb-project/simateb-backend/domain/wallet"
	"gitlab.com/simateb-project/simateb-backend/repository"
	mysqlQuery "gitlab.com/simateb-project/simateb-backend/repository/mysqlQuery/auth"
	auth2 "gitlab.com/simateb-project/simateb-backend/utils/auth"
	"gitlab.com/simateb-project/simateb-backend/utils/errorsHandler"
	"log"
	"net/http"
	"strconv"
	"time"
)

type AuthControllerInterface interface {
	Login(c *gin.Context)
}

type AuthControllerStruct struct {
}

func NewAuthController() AuthControllerInterface {
	x := &AuthControllerStruct{}
	return x
}

func (a *AuthControllerStruct) Login(c *gin.Context) {
	var userLoginRequest auth.UserLoginRequest
	if errors := c.ShouldBindJSON(&userLoginRequest); errors != nil {
		log.Println(errors.Error(), "bind")
		errorsHandler.GinErrorResponseHandler(c, errors)
		return
	}
	stmt, err := repository.DBS.MysqlDb.Prepare(mysqlQuery.LoginQuery)
	if err != nil {
		log.Println(err.Error(), "prepare")
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	var userLoginInfo auth.UserLoginInfo
	result := stmt.QueryRow(&userLoginRequest.Tel)
	err = result.Scan(&userLoginInfo.ID,
		&userLoginInfo.FirstName,
		&userLoginInfo.LastName,
		&userLoginInfo.Tel,
		&userLoginInfo.UserGroupID,
		&userLoginInfo.OrganizationID)
	if err != nil {
		log.Println(err.Error(), "scan")
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	var response *auth.ResponseAccessToken
	response, err = createToken(&userLoginInfo)
	if err != nil {
		log.Println(err.Error(), "token")
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	c.JSON(http.StatusOK, response)
}

func createToken(user *auth.UserLoginInfo) (*auth.ResponseAccessToken, error) {
	claims := auth2.UserClaims{
		UserID:         user.ID,
		Tel:            user.Tel,
		LastName:       user.LastName,
		FirstName:      user.FirstName,
		UserGroupID:    user.UserGroupID,
		OrganizationID: user.OrganizationID,
		Wallet: wallet.GetWallet(user.ID, "user"),
	}
	claims.ExpiresAt = time.Now().Unix() + constant.ExpTime
	claims.Issuer = strconv.Itoa(int(user.ID))
	token, err := claims.GenerateToken()
	if err != nil {
		return nil, err
	}
	response := auth.ResponseAccessToken{
		AccessToken:   *token,
		ExpiresIn:     claims.ExpiresAt,
		UserLoginInfo: *user,
	}
	return &response, nil
}
