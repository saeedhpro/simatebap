package smsController

import (
	"github.com/gin-gonic/gin"
	sms2 "gitlab.com/simateb-project/simateb-backend/domain/sms"
	"gitlab.com/simateb-project/simateb-backend/repository"
	"gitlab.com/simateb-project/simateb-backend/repository/sms"
	"gitlab.com/simateb-project/simateb-backend/utils/auth"
	"gitlab.com/simateb-project/simateb-backend/utils/errorsHandler"
	"log"
	"net/http"
	"strconv"
)

type SMSControllerInterface interface {
	Create(c *gin.Context)
	Get(c *gin.Context)
	GetList(c *gin.Context)
	GetListForAdmin(c *gin.Context)
	Delete(c *gin.Context)
}

type SMSControllerStruct struct {
}

func NewSMSController() SMSControllerInterface {
	x := &SMSControllerStruct{
	}
	return x
}

func (uc *SMSControllerStruct) Create(c *gin.Context) {
	var createSMSRequest sms2.SendSMSRequest
	if errors := c.ShouldBindJSON(&createSMSRequest); errors != nil {
		log.Println(errors.Error())
		errorsHandler.GinErrorResponseHandler(c, errors)
		return
	}
	staffID := auth.GetStaffUser(c).UserID
	sent, err := sms.SendSMS(createSMSRequest, staffID)
	if err != nil {
		log.Println(err.Error())
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	c.JSON(http.StatusOK, sent)
}

func (uc *SMSControllerStruct) Get(c *gin.Context) {

}

func (uc *SMSControllerStruct) GetListForAdmin(c *gin.Context) {
	query := "SELECT s.id id, ifnull(s.fname, '') user_fname, ifnull(s.lname, '') user_lname, ifnull(s.user_id, 0) user_id, s.number, ifnull(s.msg, '') msg, s.created, s.sent, ifnull(s.staff_fname, '') staff_fname, ifnull(s.staff_lname, '') staff_lname FROM (SELECT c.id, user.fname staff_fname, user.lname staff_lname, c.fname, c.lname, c.user_id, c.number, c.msg, c.created, c.sent from (SELECT sms.id id, sms.user_id, sms.staff_id, sms.number, sms.msg, sms.sent, sms.created, user.fname , user.lname  from sms LEFT join user on sms.user_id = user.id ) c left join user on c.staff_id = user.id) s "
	page := c.Query("page")
	if page == "" {
		page = "1"
	}
	q := c.Query("q")
	if q != "" {
		q = "WHERE (s.fname LIKE %" + q + "% OR s.lname LIKE %" + q + "% OR s.msg LIKE %" + q + "%)"
	}
	offset, err := strconv.Atoi(page)
	offset = (offset - 1) * 10
	query += " limit 10 offset ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error())
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	var SMSList []sms2.SMS
	var SMS sms2.SMS
	rows, err := stmt.Query(offset)
	defer rows.Close()
	if err != nil {
		log.Println(err.Error())
		return
	}
	for rows.Next() {
		err = rows.Scan(
			&SMS.ID,
			&SMS.UserFname,
			&SMS.UserLname,
			&SMS.UserID,
			&SMS.Number,
			&SMS.Msg,
			&SMS.Created,
			&SMS.Sent,
			&SMS.StaffFname,
			&SMS.StaffLname,
		)
		if err != nil {
			log.Println(err.Error())
		}
		SMSList = append(SMSList, SMS)
	}
	c.JSON(http.StatusOK, SMSList)
}

func (uc *SMSControllerStruct) GetList(c *gin.Context) {
	query := "SELECT s.id id, ifnull(s.fname, '') user_fname, ifnull(s.lname, '') user_lname, s.user_id user_id, s.number, ifnull(s.msg, '') msg, s.created, s.sent, ifnull(s.staff_fname, '') staff_fname, ifnull(s.staff_lname, '') staff_lname FROM (SELECT c.id, user.fname staff_fname, user.lname staff_lname, c.fname, c.lname, c.user_id, c.number, c.msg, c.created, c.sent, user.organization_id organization_id from (SELECT sms.id id, sms.user_id, sms.staff_id, sms.number, sms.msg, sms.sent, sms.created, user.fname , user.lname, user.organization_id organization_id from sms LEFT join user on sms.user_id = user.id ) c left join user on c.staff_id = user.id) s WHERE s.organization_id = ? "
	page := c.Query("page")
	if page == "" {
		page = "1"
	}
	q := c.Query("q")
	if q != "" {
		q = "AND (s.fname LIKE %" + q + "% OR s.lname LIKE %" + q + "% OR s.msg LIKE %" + q + "%)"
	}
	offset, err := strconv.Atoi(page)
	offset = (offset - 1) * 10
	query += " LIMIT 10 OFFSET ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error())
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	staff := auth.GetStaffUser(c)
	SMSList := []sms2.SMS{}
	var SMS sms2.SMS
	rows, err := stmt.Query(staff.OrganizationID, offset)
	defer rows.Close()
	if err != nil {
		log.Println(err.Error())
		return
	}
	for rows.Next() {
		err := rows.Scan(
			&SMS.ID,
			&SMS.UserFname,
			&SMS.UserLname,
			&SMS.UserID,
			&SMS.Number,
			&SMS.Msg,
			&SMS.Created,
			&SMS.Sent,
			&SMS.StaffFname,
			&SMS.StaffLname,
		)
		if err != nil {
			log.Println(err.Error())
		}
		SMSList = append(SMSList, SMS)
	}
	c.JSON(http.StatusOK, SMSList)
}

func (uc *SMSControllerStruct) Delete(c *gin.Context) {
	query := "DELETE FROM `sms` WHERE id in (?)"
	var deleteSMSRequest sms2.DeleteSMSRequest
	if errors := c.ShouldBindJSON(&deleteSMSRequest); errors != nil {
		log.Println(errors.Error())
		errorsHandler.GinErrorResponseHandler(c, errors)
		return
	}
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error())
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		&deleteSMSRequest.IDs,
	)
	if err != nil {
		log.Println(err.Error())
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	c.JSON(http.StatusOK, nil)
}
