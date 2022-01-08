package smsController

import (
	"fmt"
	"github.com/gin-gonic/gin"
	sms2 "gitlab.com/simateb-project/simateb-backend/domain/sms"
	"gitlab.com/simateb-project/simateb-backend/repository"
	"gitlab.com/simateb-project/simateb-backend/repository/pushe"
	"gitlab.com/simateb-project/simateb-backend/utils/auth"
	"gitlab.com/simateb-project/simateb-backend/utils/errorsHandler"
	"gitlab.com/simateb-project/simateb-backend/utils/pagination"
	"log"
	"net/http"
	"strconv"
)

type NotificationControllerInterface interface {
	Create(c *gin.Context)
	Get(c *gin.Context)
	GetList(c *gin.Context)
	GetListForAdmin(c *gin.Context)
	Delete(c *gin.Context)
}

type NotificationControllerStruct struct {
}

func NewNotificationController() NotificationControllerInterface {
	x := &NotificationControllerStruct{
	}
	return x
}

func (uc *NotificationControllerStruct) Create(c *gin.Context) {
	var request pushe.SendNotificationRequest
	if errors := c.ShouldBindJSON(&request); errors != nil {
		log.Println(errors.Error())
		errorsHandler.GinErrorResponseHandler(c, errors)
		return
	}
	staffUser := auth.GetStaffUser(c)
	insertNotification(staffUser, request)
	sent, err := pushe.SendNotification(staffUser, request)
	if err != nil {
		log.Println(err.Error())
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	c.JSON(http.StatusOK, sent)
}

func insertNotification(user *auth.UserClaims, request pushe.SendNotificationRequest) {
	query := "INSERT INTO `notifications`(`action_url`, `close_on_click`, `content`, `title`, `ids`, `type`, `user_id`, `staff_id`) VALUES (?,?,?,?,?,?,?,?)"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error())
		return
	}
	_, err = stmt.Exec(
		&request.ActionUrl,
		&request.CloseOnClick,
		&request.Content,
		&request.Title,
		&request.Ids,
		&request.Type,
		&request.UserID,
		&user.UserID,
	)
	if err != nil {
		log.Println(err.Error())
		return
	}
	return
}

func (uc *NotificationControllerStruct) Get(c *gin.Context) {

}

func (uc *NotificationControllerStruct) GetListForAdmin(c *gin.Context) {
	query := "SELECT s.id id, ifnull(s.fname, '') user_fname, ifnull(s.lname, '') user_lname, ifnull(s.user_id, 0) user_id, ifnull(s.title, ''), ifnull(s.content, '') content, s.created_at, s.type, ifnull(s.staff_fname, '') staff_fname, ifnull(s.staff_lname, '') staff_lname, ifnull(s.ids, '') ids, ifnull(s.action_url, '') action_url, ifnull(s.close_on_click, '') close_on_click, ifnull(s.close_on_click, '') close_on_click FROM (SELECT c.id, user.fname staff_fname, user.lname staff_lname, c.fname, c.lname, c.user_id, c.title, c.content, c.created_at, c.type, c.ids, c.action_url, c.close_on_click from (SELECT notifications.id id, notifications.user_id, notifications.staff_id, notifications.title, notifications.content, notifications.action_url, notifications.close_on_click, notifications.type, notifications.created_at, notifications.ids, user.fname , user.lname from notifications LEFT join user on notifications.user_id = user.id ) c left join user on c.staff_id = user.id) s "
	page := c.Query("page")
	if page == "" {
		page = "1"
	}
	q := c.Query("q")
	if q != "" && q != "null" && q != "undefined" {
		query += " WHERE (s.fname LIKE '%" + q + "%' OR s.lname LIKE '%" + q + "%' OR s.content LIKE '%" + q + "%' OR s.title LIKE '%" + q + "%')"
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
	paginationInfo := pagination.NotificationPaginationInfo{}
	var notifList []pushe.Notification
	var notif pushe.Notification
	rows, err := stmt.Query(offset)
	defer rows.Close()
	if err != nil {
		log.Println(err.Error())
		return
	}
	for rows.Next() {
		err = rows.Scan(
			&notif.ID,
			&notif.UserFName,
			&notif.UserLName,
			&notif.UserID,
			&notif.Title,
			&notif.Content,
			&notif.CreatedAt,
			&notif.Type,
			&notif.StaffFName,
			&notif.StaffLName,
			&notif.Ids,
			&notif.ActionUrl,
		)
		if err != nil {
			log.Println(err.Error())
		}
		notifList = append(notifList, notif)
	}
	paginationInfo.Data = notifList
	count := 0
	p, _ := strconv.Atoi(page)
	count, _ = getNotificationCountAdmin(q)
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

func getNotificationCountAdmin(q string) (int, error) {
	query := "SELECT COUNT(*) FROM (SELECT c.id, user.fname staff_fname, user.lname staff_lname, c.fname, c.lname, c.user_id, c.number, c.msg, c.created, c.sent from (SELECT sms.id id, sms.user_id, sms.staff_id, sms.number, sms.msg, sms.sent, sms.created, user.fname , user.lname  from sms LEFT join user on sms.user_id = user.id ) c left join user on c.staff_id = user.id) s "
	var values []interface{}
	count := 0
	if q != "" && q != "null" && q != "undefined" {
		query += "AND (s.fname LIKE '%" + q + "%' OR s.lname LIKE '%" + q + "%' OR s.msg LIKE '%" + q + "%')"
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

func (uc *NotificationControllerStruct) GetList(c *gin.Context) {
	query := "SELECT s.id id, ifnull(s.fname, '') user_fname, ifnull(s.lname, '') user_lname, ifnull(s.user_id, 0) user_id, ifnull(s.title, ''), ifnull(s.content, '') content, s.created_at, s.type, ifnull(s.staff_fname, '') staff_fname, ifnull(s.staff_lname, '') staff_lname, ifnull(s.ids, '') ids, ifnull(s.action_url, '') action_url, ifnull(s.close_on_click, '') close_on_click FROM (SELECT c.id id, c.organization_id, user.fname staff_fname, user.lname staff_lname, c.fname fname, c.lname lname, c.user_id user_id, c.title title , c.content content, c.created_at created_at, c.type, c.ids, c.action_url, c.close_on_click from (SELECT notifications.id id, notifications.user_id, notifications.staff_id, notifications.title, notifications.content, notifications.action_url, notifications.close_on_click, notifications.type, notifications.created_at, notifications.ids, user.organization_id, user.fname , user.lname from notifications LEFT join user on notifications.user_id = user.id ) c left join user on c.staff_id = user.id) s WHERE s.user_id = ? "
	page := c.Query("page")
	if page == "" {
		page = "1"
	}
	q := c.Query("q")
	if q != "" && q != "null" && q != "undefined" {
		query += " AND (s.fname LIKE '%" + q + "%' OR s.lname LIKE '%" + q + "%' OR s.content LIKE '%" + q + "%' OR s.title LIKE '%" + q + "%')"
	}
	staffUser := auth.GetStaffUser(c)
	offset, err := strconv.Atoi(page)
	offset = (offset - 1) * 10
	query += " limit 10 offset ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error())
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	paginationInfo := pagination.NotificationPaginationInfo{}
	notifList := []pushe.Notification{}
	var notif pushe.Notification
	rows, err := stmt.Query(staffUser.UserID, offset)
	defer rows.Close()
	if err != nil {
		log.Println(err.Error())
		return
	}
	for rows.Next() {
		err = rows.Scan(
			&notif.ID,
			&notif.UserFName,
			&notif.UserLName,
			&notif.UserID,
			&notif.Title,
			&notif.Content,
			&notif.CreatedAt,
			&notif.Type,
			&notif.StaffFName,
			&notif.StaffLName,
			&notif.Ids,
			&notif.ActionUrl,
			&notif.CloseOnClick,
		)
		if err != nil {
			log.Println(err.Error())
		}
		notifList = append(notifList, notif)
	}
	paginationInfo.Data = notifList
	count := 0
	p, _ := strconv.Atoi(page)
	count, _ = getNotificationCount(fmt.Sprintf("%d", staffUser.UserID), q)
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

func getNotificationCount(userID string, q string) (int, error) {
	query := "SELECT COUNT(*) FROM (SELECT c.id id, c.organization_id, user.fname staff_fname, user.lname staff_lname, c.fname fname, c.lname lname, c.user_id user_id, c.title title , c.content content, c.created_at created_at, c.type, c.ids, c.action_url, c.close_on_click from (SELECT notifications.id id, notifications.user_id, notifications.staff_id, notifications.title, notifications.content, notifications.action_url, notifications.close_on_click, notifications.type, notifications.created_at, notifications.ids, user.organization_id, user.fname , user.lname from notifications LEFT join user on notifications.user_id = user.id ) c left join user on c.staff_id = user.id) s WHERE s.user_id = ? "
	var values []interface{}
	values = append(values, userID)
	count := 0
	log.Println(userID)
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

func (uc *NotificationControllerStruct) Delete(c *gin.Context) {
	query := "DELETE FROM `notifications` WHERE id in ("
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
		if i != len(deleteSMSRequest.IDs)-1 {
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
