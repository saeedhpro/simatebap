package smsController

import (
	"fmt"
	"github.com/gin-gonic/gin"
	sms2 "gitlab.com/simateb-project/simateb-backend/domain/sms"
	"gitlab.com/simateb-project/simateb-backend/repository"
	"gitlab.com/simateb-project/simateb-backend/repository/sms"
	"gitlab.com/simateb-project/simateb-backend/utils/auth"
	"gitlab.com/simateb-project/simateb-backend/utils/errorsHandler"
	"gitlab.com/simateb-project/simateb-backend/utils/pagination"
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
	if q != "" && q != "null" && q != "undefined" {
		query += " WHERE (s.fname LIKE '%" + q + "%' OR s.lname LIKE '%" + q + "%' OR s.msg LIKE '%" + q + "%')"
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
	paginationInfo := pagination.SMSPaginationInfo{}
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
	paginationInfo.Data = SMSList
	count := 0
	p, _ := strconv.Atoi(page)
	count, _ = getSMSCountAdmin(q)
	paginationInfo.PagesCount = count
	paginationInfo.Page = p
	if p > 1 {
		paginationInfo.PrevPage = p - 1
	} else {
		paginationInfo.PrevPage = p
	}
	if p < count/10 {
		paginationInfo.NextPage = p
	} else {
		paginationInfo.NextPage = p + 1
	}
	paginationInfo.HasNextPage = (bool)(count > 10 && count > (p*10))
	c.JSON(http.StatusOK, paginationInfo)
}

func getSMSCountAdmin(q string) (int, error) {
	query := "SELECT COUNT(*) FROM (SELECT c.id, user.fname staff_fname, user.lname staff_lname, c.fname, c.lname, c.user_id, c.number, c.msg, c.created, c.sent from (SELECT sms.id id, sms.user_id, sms.staff_id, sms.number, sms.msg, sms.sent, sms.created, user.fname , user.lname  from sms LEFT join user on sms.user_id = user.id ) c left join user on c.staff_id = user.id) s "
	var values []interface{}
	count := 0
	if q != "" && q != "null" && q != "undefined" {
		query +="AND (s.fname LIKE '%" + q + "%' OR s.lname LIKE '%" + q + "%' OR s.msg LIKE '%" + q + "%')"
	}
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error())
		return count, nil
	}
	result := stmt.QueryRow(values...)
	err = result.Scan(&count)
	if err != nil {
		log.Println(err.Error(), "count")
		return count, nil
	}
	return count, nil
}

func (uc *SMSControllerStruct) GetList(c *gin.Context) {
	query := "SELECT s.id id, ifnull(s.fname, '') user_fname, ifnull(s.lname, '') user_lname, s.user_id user_id, s.number, ifnull(s.msg, '') msg, s.created, s.sent, ifnull(s.staff_fname, '') staff_fname, ifnull(s.staff_lname, '') staff_lname FROM (SELECT c.id, user.fname staff_fname, user.lname staff_lname, c.fname, c.lname, c.user_id, c.number, c.msg, c.created, c.sent, user.organization_id organization_id from (SELECT sms.id id, sms.user_id, sms.staff_id, sms.number, sms.msg, sms.sent, sms.created, user.fname , user.lname, user.organization_id organization_id from sms LEFT join user on sms.user_id = user.id ) c left join user on c.staff_id = user.id) s WHERE s.organization_id = ? "
	page := c.Query("page")
	if page == "" {
		page = "1"
	}
	q := c.Query("q")
	if q != "" && q != "null" && q != "undefined" {
		query += "AND (s.fname LIKE '%" + q + "%' OR s.lname LIKE '%" + q + "%' OR s.msg LIKE '%" + q + "%')"
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
	paginationInfo := pagination.SMSPaginationInfo{}
	rows, err := stmt.Query(staff.OrganizationID, offset)
	if err != nil {
		log.Println(err.Error())
		errorsHandler.GinErrorResponseHandler(c, err)
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
			c.JSON(http.StatusOK, paginationInfo)
			return
		}
		SMSList = append(SMSList, SMS)
	}
	paginationInfo.Data = SMSList
	count := 0
	p, _ := strconv.Atoi(page)
	count, _ = getSMSCount(fmt.Sprintf("%d", staff.OrganizationID), q)
	paginationInfo.PagesCount = count
	paginationInfo.Page = p
	if p > 1 {
		paginationInfo.PrevPage = p - 1
	} else {
		paginationInfo.PrevPage = p
	}
	if p < count/10 {
		paginationInfo.NextPage = p
	} else {
		paginationInfo.NextPage = p + 1
	}
	paginationInfo.HasNextPage = (bool)(count > 10 && count > (p*10))
	c.JSON(http.StatusOK, paginationInfo)
}

func getSMSCount(organizationID string, q string) (int, error) {
	query := "SELECT COUNT(*) FROM (SELECT c.id, user.fname staff_fname, user.lname staff_lname, c.fname, c.lname, c.user_id, c.number, c.msg, c.created, c.sent, user.organization_id organization_id from (SELECT sms.id id, sms.user_id, sms.staff_id, sms.number, sms.msg, sms.sent, sms.created, user.fname , user.lname, user.organization_id organization_id from sms LEFT join user on sms.user_id = user.id ) c left join user on c.staff_id = user.id) s WHERE s.organization_id = ? "
	var values []interface{}
	values = append(values, organizationID)
	count := 0
	if q != "" && q != "null" && q != "undefined" {
		query += " AND (s.fname LIKE '%" + q + "%' OR s.lname LIKE '%" + q + "%' OR s.msg LIKE '%" + q + "%')"
	}
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error())
		return count, nil
	}
	result := stmt.QueryRow(values...)
	err = result.Scan(&count)
	if err != nil {
		log.Println(err.Error(), "count")
		return count, nil
	}
	return count, nil
}

func (uc *SMSControllerStruct) Delete(c *gin.Context) {
	query := "DELETE FROM `sms` WHERE id in ("
	var deleteSMSRequest sms2.DeleteSMSRequest
	if errors := c.ShouldBindJSON(&deleteSMSRequest); errors != nil {
		log.Println(errors.Error())
		errorsHandler.GinErrorResponseHandler(c, errors)
		return
	}
	if len(deleteSMSRequest.IDs) == 0 {
		c.JSON(http.StatusOK, nil)
		return
	}
	var values []interface{}
	for i := 0; i < len(deleteSMSRequest.IDs); i++ {
		values = append(values, deleteSMSRequest.IDs[i])
		if i != len(deleteSMSRequest.IDs) - 1 {
			query += "?,"
		} else {
			query += "?"
		}
	}
	query += ")"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error())
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
		_, err = stmt.Exec(
		values...,
	)
	if err != nil {
		log.Println(err.Error())
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	c.JSON(http.StatusOK, nil)
}
