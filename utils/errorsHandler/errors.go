package errorsHandler

import (
	"github.com/gin-gonic/gin"
)

const (
	FA = "fa"
	EN = "en"
)

var Model ErrorsModel

func Init() {
	for k, v := range DefaultFAErrorsMap {
		FAErrorsMap[k] = v
	}
	for k, v := range DefaultENErrorsMap {
		ENErrorsMap[k] = v
	}
	for k, v := range DefaultCodeErrorsMap {
		CodeErrorsMap[k] = v
	}
	errors := ErrorsModel{
		ErrorsMap :     ErrorsMap,
		CodeErrorsMap : CodeErrorsMap,
	}
	InitErrors(errors)
}

func InitErrors(errors ErrorsModel) {
	Model = errors
}
func GinErrorResponseHandler(c *gin.Context, err error) {
	type Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	type ErrorsResponse struct {
		Errors []Error `json:"errors"`
	}
	var res ErrorsResponse
	Lang := c.GetHeader("Accept-Language")
	if Lang == "" {
		Lang = FA
	}
	errMessage, exists := Model.ErrorsMap[Lang][err.Error()]
	errCode, existsCode := Model.CodeErrorsMap[err.Error()]
	if !existsCode {
		errCode = 400
	}
	if exists {
		res.Errors = []Error{
			{
				Code:    errCode,
				Message: errMessage,
			},
		}
	} else {
		res.Errors = []Error{
			{
				Code:    400,
				Message: DefaultError,
			},
		}
	}
	c.AbortWithStatusJSON(errCode, res)
}
