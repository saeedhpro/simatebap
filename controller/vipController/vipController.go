package vipController

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/simateb-project/simateb-backend/domain/appointment"
	"gitlab.com/simateb-project/simateb-backend/helper"
	"gitlab.com/simateb-project/simateb-backend/repository"
	"gitlab.com/simateb-project/simateb-backend/utils/auth"
	"gitlab.com/simateb-project/simateb-backend/utils/errorsHandler"
	"log"
	"net/http"
)

type VipControllerInterface interface {

}

type VipControllerStruct struct {

}

func NewVipControllerStruct() VipControllerInterface {
	x := &VipControllerStruct{
	}
	return x
}

func (uc *VipControllerStruct) Create(c *gin.Context) {
	var createAppointmentRequest appointment.CreateAppointmentRequest
	if errors := c.ShouldBindJSON(&createAppointmentRequest); errors != nil {
		log.Println(errors.Error())
		errorsHandler.GinErrorResponseHandler(c, errors)
		return
	}
	query := ""
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error())
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	defer stmt.Close()
	randomCode := helper.RandomString(6)
	staffID := auth.GetStaffUser(c).UserID
	result, err := stmt.Exec(
		&createAppointmentRequest.UserID,
		&createAppointmentRequest.Info,
		&createAppointmentRequest.StartAt,
		&createAppointmentRequest.CaseType,
		&createAppointmentRequest.Income,
		staffID,
		&createAppointmentRequest.IsVip,
		&randomCode,
	)
	if err != nil {
		log.Println(err.Error())
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	_, err = result.LastInsertId()
	if err != nil {
		log.Println(err.Error())
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	c.JSON(http.StatusOK, true)
}
