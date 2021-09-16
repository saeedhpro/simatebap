package middleware

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/simateb-project/simateb-backend/constant"
	"gitlab.com/simateb-project/simateb-backend/utils/auth"
	"log"
	"net/http"
	"strings"
)

const (
	claimKey = "claims"
)

func GinJwtAuth(function gin.HandlerFunc, selfAccess, optional bool) gin.HandlerFunc {
	return func(context *gin.Context) {
		token := strings.Replace(context.GetHeader("Authorization"), "Bearer ", "", -1)
		if token == "" && optional {
			context.Set(claimKey, &auth.UserClaims{})
			function(context)
			context.Next()
		} else {
			claims, err := auth.ValidateToken(token)
			if err != nil {
				log.Println(err.Error())
				context.JSON(http.StatusUnauthorized, gin.H{
					"error": constant.UnAuthorizedError,
				})
				return
			}

			// TODO: VALIDATE USER ACCESS

			context.Set(claimKey, *claims)
			function(context)
			context.Next()
		}
	}
}
