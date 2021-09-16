package appointmentController

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/simateb-project/simateb-backend/controller/caseTypeController"
	"gitlab.com/simateb-project/simateb-backend/domain/appointment"
	"gitlab.com/simateb-project/simateb-backend/domain/caseType"
	"gitlab.com/simateb-project/simateb-backend/helper"
	"gitlab.com/simateb-project/simateb-backend/repository"
	appointment2 "gitlab.com/simateb-project/simateb-backend/repository/appointment"
	mysqlQuery "gitlab.com/simateb-project/simateb-backend/repository/mysqlQuery/auth"
	"gitlab.com/simateb-project/simateb-backend/utils/auth"
	"gitlab.com/simateb-project/simateb-backend/utils/errorsHandler"
	"gitlab.com/simateb-project/simateb-backend/utils/pagination"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type AppointmentControllerInterface interface {
	Create(c *gin.Context)
	Get(c *gin.Context)
	GetOperationList(c *gin.Context)
	GetAppointmentList(c *gin.Context)
	GetQueDetails(c *gin.Context)
	Update(c *gin.Context)
	ChangeStatus(c *gin.Context)
	Delete(c *gin.Context)
	SearchAppointment(c *gin.Context)
}

type AppointmentControllerStruct struct {
	r *appointment2.AppointmentRepositoryStruct
}

func NewAppointmentController(r *appointment2.AppointmentRepositoryStruct) AppointmentControllerInterface {
	x := &AppointmentControllerStruct{
		r: r,
	}
	return x
}

func (uc *AppointmentControllerStruct) Create(c *gin.Context) {
	var createAppointmentRequest appointment.CreateAppointmentRequest
	if errors := c.ShouldBindJSON(&createAppointmentRequest); errors != nil {
		log.Println(errors.Error())
		errorsHandler.GinErrorResponseHandler(c, errors)
		return
	}
	stmt, err := repository.DBS.MysqlDb.Prepare(mysqlQuery.CreateAppointmentQuery)
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

func (uc *AppointmentControllerStruct) Get(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		return
	}
	stmt, err := repository.DBS.MysqlDb.Prepare(mysqlQuery.GetAppointmentQuery)
	if err != nil {
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	var appointment appointment.SimpleAppointmentInfo
	result := stmt.QueryRow(id)
	err = result.Scan(
		&appointment.ID,
		&appointment.CaseType,
		&appointment.IsVip,
		&appointment.StartAt,
		&appointment.UserID,
		&appointment.Info,
		&appointment.Income,
		&appointment.Status,
		&appointment.UpdatedAt,
		&appointment.UserFName,
		&appointment.UserLName,
		&appointment.UserID,
		&appointment.UserGender,
	)
	if err != nil {
		log.Println(err.Error())
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, "یافت نشد")
			return
		}
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	c.JSON(http.StatusOK, appointment)
}

func (uc *AppointmentControllerStruct) GetOperationList(c *gin.Context) {
	startAt := c.Query("start_at")
	endAt := c.Query("end_at")
	if startAt == "" || endAt == "" {
		log.Println("start at or and end at are needed!")
		c.JSON(422, gin.H{
			"message": "start at or and end at are needed!",
		})
		return
	}
	startAtDate := fmt.Sprintf("%s 00:00:00", startAt)
	endAtDate := fmt.Sprintf("%s 00:00:00", endAt)
	stmt, err := repository.DBS.MysqlDb.Prepare(mysqlQuery.GetOrganizationOperationListQuery)
	if err != nil {
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	staffUser := auth.GetStaffUser(c)
	organizationID := staffUser.OrganizationID
	rows, error := stmt.Query(organizationID, startAtDate, endAtDate)
	if error != nil {
		log.Println(error.Error(), "error")
		errorsHandler.GinErrorResponseHandler(c, error)
		return
	}
	var operations []appointment.OperationInfo
	var operation appointment.OperationInfo
	for rows.Next() {
		rows.Scan(
			&operation.ID,
			&operation.UserID,
			&operation.StartAt,
			&operation.Info,
			&operation.Income,
			&operation.CaseType,
		)
		operations = append(operations, operation)
	}
	c.JSON(http.StatusOK, operations)
}

func (uc *AppointmentControllerStruct) GetAppointmentList(c *gin.Context) {
	startAt := c.Query("start_at")
	endAt := c.Query("end_at")
	status := c.Query("status")
	if startAt == "" || endAt == "" {
		log.Println("start at or and end at are needed!")
		c.JSON(422, gin.H{
			"message": "start at or and end at are needed!",
		})
		return
	}
	if status == "" {
		log.Println("status is needed!")
		c.JSON(422, gin.H{
			"message": "status is needed!",
		})
		return
	}
	startAtDate := fmt.Sprintf("%s 00:00:00", startAt)
	endAtDate := fmt.Sprintf("%s 00:00:00", endAt)
	stmt, err := repository.DBS.MysqlDb.Prepare(mysqlQuery.GetOrganizationAppointmentListQuery)
	if err != nil {
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	staffUser := auth.GetStaffUser(c)
	organizationID := staffUser.OrganizationID
	rows, error := stmt.Query(organizationID, startAtDate, endAtDate, status)
	if error != nil {
		log.Println(error.Error(), "error")
		errorsHandler.GinErrorResponseHandler(c, error)
		return
	}
	var operations []appointment.OperationInfo
	var operation appointment.OperationInfo
	for rows.Next() {
		rows.Scan(
			&operation.ID,
			&operation.UserID,
			&operation.StartAt,
			&operation.Info,
			&operation.Income,
			&operation.CaseType,
		)
		operations = append(operations, operation)
	}
	c.JSON(http.StatusOK, operations)
}

func (uc *AppointmentControllerStruct) Update(c *gin.Context) {
	var updateUserQuery = "UPDATE `appointment` SET"
	var values []interface{}
	var columns []string
	userId := c.Param("id")
	if userId == "" {
		errorsHandler.GinErrorResponseHandler(c, nil)
		return
	}
	var request appointment.UpdateAppointmentRequest
	if errors := c.ShouldBindJSON(&request); errors != nil {
		log.Println(errors.Error())
		errorsHandler.GinErrorResponseHandler(c, errors)
		return
	}
	getAppointmentUpdateColumns(&request, &columns, &values)
	columnsString := strings.Join(columns, ",")
	updateUserQuery += columnsString
	updateUserQuery += " WHERE `id` = ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(updateUserQuery)
	if err != nil {
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	values = append(values, userId)
	_, error := stmt.Exec(values...)
	if error != nil {
		log.Println(error.Error())
		errorsHandler.GinErrorResponseHandler(c, error)
		return
	}
	c.JSON(200, true)
}

func (uc *AppointmentControllerStruct) ChangeStatus(c *gin.Context) {
	var changeStatusQuery = "UPDATE `appointment` SET `status` = ? WHERE `id` = ?"
	var values []interface{}
	userId := c.Param("id")
	if userId == "" {
		errorsHandler.GinErrorResponseHandler(c, nil)
		return
	}
	var request appointment.ChangeAppointmentStatusRequest
	if errors := c.ShouldBindJSON(&request); errors != nil {
		log.Println(errors.Error())
		errorsHandler.GinErrorResponseHandler(c, errors)
		return
	}
	stmt, err := repository.DBS.MysqlDb.Prepare(changeStatusQuery)
	if err != nil {
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	values = append(values, request.Status, userId)
	_, error := stmt.Exec(values...)
	if error != nil {
		log.Println(error.Error())
		errorsHandler.GinErrorResponseHandler(c, error)
		return
	}
	c.JSON(200, true)
}

func getAppointmentUpdateColumns(o *appointment.UpdateAppointmentRequest, columns *[]string, values *[]interface{}) {
	if o.Info != "" {
		*columns = append(*columns, " `info` = ? ")
		*values = append(*values, o.Info)
	}
	if o.CaseType != "" {
		*columns = append(*columns, " `case_type` = ? ")
		*values = append(*values, o.CaseType)
	}
	*columns = append(*columns, " `income` = ? ")
	*values = append(*values, o.Income)
}

func (uc *AppointmentControllerStruct) Delete(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		return
	}
	stmt, err := repository.DBS.MysqlDb.Prepare(mysqlQuery.DeleteUserQuery)
	if err != nil {
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	stmt.QueryRow(userID)
	c.JSON(200, nil)
}

func (uc *AppointmentControllerStruct) GetQueDetails(c *gin.Context) {
	var caseTypes []caseType.CaseTypeInfo
	staffUser := auth.GetStaffUser(c)
	startAt := c.Query("start_at")
	endAt := c.Query("end_at")
	status := c.Query("status")
	if startAt == "" || endAt == "" || status == "" {
		log.Println("start at or and end at are needed!")
		c.JSON(422, gin.H{
			"message": "start at or and end at and status are needed!",
		})
		return
	}
	startAtDate := fmt.Sprintf("%s 00:00:00", startAt)
	endAtDate := fmt.Sprintf("%s 00:00:00", endAt)
	caseTypes = caseTypeController.LoadCaseTypesByOrgId(staffUser.OrganizationID)
	var limits []appointment.Limit
	for _, c := range caseTypes {
		limit := appointment.Limit{
			ID:         c.ID,
			Name:       c.Name,
			Limitation: c.Limitation,
			Total:      uc.r.LoadTotalsByDayAndCase(staffUser.OrganizationID, c.Name, startAtDate, endAtDate),
		}
		limits = append(limits, limit)
	}
	var queDetail appointment.QueDetail
	queDetail.DefaultDuration = 20
	queDetail.Ques = uc.r.LoadQueWithOrganization(staffUser.OrganizationID, startAtDate, endAtDate, status, staffUser)
	queDetail.Limits = limits
	queDetail.WorkHours = uc.r.GetOrganizationWorkHour(staffUser.OrganizationID)
	queDetail.Totals = uc.r.LoadTotals(staffUser.OrganizationID, startAtDate, endAtDate)
	c.JSON(200, queDetail)
}

func (uc *AppointmentControllerStruct) SearchAppointment(c *gin.Context) {
	query := "SELECT appointment.id id, appointment.status status, appointment.start_at start_at, user.fname user_fname, user.lname user_lname, user.id user_id, user.tel mobile, user.file_id file_id FROM appointment LEFT JOIN user ON appointment.user_id = user.id "
	startDate := c.Query("start_at")
	endDate := c.Query("end_at")
	status := c.Query("status")
	q := c.Query("q")
	page := c.Query("page")
	if page == "" {
		page = "1"
	}
	var values []interface{}
	var queries []string
	if startDate != "" {
		values = append(values, startDate)
		queries = append(queries, " appointment.start_at > ? ")
	}
	if endDate != "" {
		values = append(values, endDate)
		queries = append(queries, " appointment.start_at <= ? ")
	}
	if q != "" {
		q = "'%" + strings.TrimSpace(q) + "%'"
		queries = append(queries, fmt.Sprintf(" user.fname LIKE %s OR user.lname LIKE %s", q, q))
	}
	if status != "" {
		values = append(values, status)
		queries = append(queries, " appointment.status IN (?) ")
	}
	values = append(values, page)
	where := strings.Join(queries, " AND ")
	if where != "" {
		where = fmt.Sprintf("%s %s", "WHERE", where)
	}
	query = fmt.Sprintf("%s %s %s", query, where, " Limit 10 offset ?")
	pagination := pagination.AppointmentPaginationInfo{}
	var appointments []appointment.SimpleAppointmentInfo
	var appointment appointment.SimpleAppointmentInfo
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println("err: ", err.Error())
		c.JSON(500, pagination)
		return
	}
	rows, err := stmt.Query(values...)
	if err != nil {
		log.Println("err: ", err.Error())
		c.JSON(500, pagination)
		return
	}
	for rows.Next() {
		err := rows.Scan(
			&appointment.ID,
			&appointment.Status,
			&appointment.StartAt,
			&appointment.UserFName,
			&appointment.UserLName,
			&appointment.UserID,
			&appointment.Mobile,
			&appointment.FileID,
		)
		if err != nil {
			log.Println("err: ", err.Error())
			c.JSON(500, pagination)
			return
		}
		appointments = append(appointments, appointment)
	}
	p, _ := strconv.Atoi(page)
	pagination.Page = p
	pagination.PrevPage = p - 1
	pagination.NextPage = p + 1
	pagination.HasNextPage = true
	pagination.Data = appointments
	c.JSON(200, pagination)
}
