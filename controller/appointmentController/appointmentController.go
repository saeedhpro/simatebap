package appointmentController

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/simateb-project/simateb-backend/controller/caseTypeController"
	"gitlab.com/simateb-project/simateb-backend/domain/admission"
	"gitlab.com/simateb-project/simateb-backend/domain/appointment"
	"gitlab.com/simateb-project/simateb-backend/domain/caseType"
	"gitlab.com/simateb-project/simateb-backend/helper"
	"gitlab.com/simateb-project/simateb-backend/repository"
	appointment2 "gitlab.com/simateb-project/simateb-backend/repository/appointment"
	mysqlQuery "gitlab.com/simateb-project/simateb-backend/repository/mysqlQuery/auth"
	"gitlab.com/simateb-project/simateb-backend/repository/user"
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
	AcceptAppointment(c *gin.Context)
	GetAppointmentByCode(c *gin.Context)
	GetAppointmentAdmissions(c *gin.Context)
	FinishAppointmentAdmissions(c *gin.Context)
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
		&appointment.Price,
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
		operation.User, _ = user.GetUserByID(operation.UserID)
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
		operation.User, _ = user.GetUserByID(operation.UserID)
		operations = append(operations, operation)
	}
	c.JSON(http.StatusOK, operations)
}

func (uc *AppointmentControllerStruct) GetAppointmentAdmissions(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		return
	}
	query := "SELECT id, appointment_id, case_name, file_name FROM appointment_admissions WHERE appointment_id = ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	rows, error := stmt.Query(id)
	if error != nil {
		log.Println(error.Error(), "error")
		errorsHandler.GinErrorResponseHandler(c, error)
		return
	}
	var admissions = []admission.AdmissionInfo{}
	var admission admission.AdmissionInfo
	for rows.Next() {
		rows.Scan(
			&admission.ID,
			&admission.AppointmentID,
			&admission.CaseName,
			&admission.FileName,
		)
		admissions = append(admissions, admission)
	}
	c.JSON(http.StatusOK, admissions)
}

func (uc *AppointmentControllerStruct) FinishAppointmentAdmissions(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		return
	}
	query := "SELECT id, status FROM `appointment` WHERE id = ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	row := stmt.QueryRow(id)
	if row.Err() != nil {
		return
	}
	var appointmentInfo appointment.SimpleAppointmentInfo
	err = row.Scan(
		&appointmentInfo.ID,
		&appointmentInfo.Status,
	)
	if err != nil {
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	if appointmentInfo.Status != 2 {
		c.JSON(http.StatusForbidden, nil)
		return
	}
	
	c.JSON(http.StatusOK, true)
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
	query := "SELECT appointment.id id, appointment.status status, ifnull(appointment.case_type, '') case_type, appointment.start_at start_at, user.fname user_fname, user.lname user_lname, user.id user_id, user.tel mobile, user.file_id file_id, appointment.price price FROM appointment LEFT JOIN user ON appointment.user_id = user.id "
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
	paginationInfo := pagination.AppointmentPaginationInfo{}
	appointments := []appointment.SimpleAppointmentInfo{}
	var appointmentInfo appointment.SimpleAppointmentInfo
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println("err: ", err.Error())
		c.JSON(500, paginationInfo)
		return
	}
	rows, err := stmt.Query(values...)
	if err != nil {
		log.Println("err: ", err.Error())
		c.JSON(500, paginationInfo)
		return
	}
	for rows.Next() {
		err := rows.Scan(
			&appointmentInfo.ID,
			&appointmentInfo.Status,
			&appointmentInfo.CaseType,
			&appointmentInfo.StartAt,
			&appointmentInfo.UserFName,
			&appointmentInfo.UserLName,
			&appointmentInfo.UserID,
			&appointmentInfo.Mobile,
			&appointmentInfo.FileID,
			&appointmentInfo.Price,
		)
		if err != nil {
			log.Println("err: ", err.Error())
			c.JSON(500, paginationInfo)
			return
		}
		appointments = append(appointments, appointmentInfo)
	}
	p, _ := strconv.Atoi(page)
	paginationInfo.Page = p
	paginationInfo.PrevPage = p - 1
	paginationInfo.NextPage = p + 1
	paginationInfo.HasNextPage = true
	paginationInfo.Data = appointments
	c.JSON(200, paginationInfo)
}

func (uc *AppointmentControllerStruct) AcceptAppointment(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		return
	}
	var request appointment.AcceptAppointmentRequest
	if errors := c.ShouldBindJSON(&request); errors != nil {
		log.Println(errors.Error())
		errorsHandler.GinErrorResponseHandler(c, errors)
		return
	}
	code := helper.RandomString(6)
	price := 0.0
	if len(request.RadiologyCases) > 0 {
		price += calcRadiologyCasesPrice(request.RadiologyCases)
	}
	status := acceptAppointment(id, code, price, request)
	c.JSON(200 , status)
}

func acceptAppointment(id string, code string, price float64, request appointment.AcceptAppointmentRequest) bool {
	query := "UPDATE `appointment` SET `status`= ?, `photography_cases`= ? ,`radiology_cases`= ? ,`prescription`= ? ," +
		"`future_prescription`= ? , `photography_msg`= ?, `radiology_msg`= ?, `radiology_id`= ?," +
		"`photography_id`= ?, `code`= ?, price = ? WHERE id = ? "
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		return false
	}
	_, error := stmt.Exec(
		2,
		strings.Join(request.PhotographyCases, ","),
		strings.Join(request.RadiologyCases, ","),
		request.Prescription,
		request.FuturePrescription,
		request.PhotographyMsg,
		request.RadiologyMsg,
		request.RadiologyID,
		request.PhotographyID,
		code,
		id,
		price,
	)
	if error != nil {
		log.Println(error.Error())
		return false
	}
	return true
}

func calcRadiologyCasesPrice(cases []string) float64 {
	price := 0.0
	ids := getRadiologyCasesID(cases)
	idsStr := strings.Trim(strings.Replace(fmt.Sprint(ids), " ", ",", -1), "[]")
	price += getRadiologyCasesPrices(idsStr)
	return price
}

func getRadiologyCasesPrices(ids string) float64 {
	price := 0.0
	query := "SELECT SUM(price) price FROM `case_prices` WHERE id in (" + ids + ") and organization_id = ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		return price
	}
	rows, error := stmt.Query()
	if error != nil {
		return price
	}
	for rows.Next() {
		rows.Scan(&price)
	}
	return price
}

func getRadiologyCasesID(cases []string) []int64 {
	query := "SELECT cases.id, cases.name FROM cases WHERE cases.name in (" + strings.Join(cases, ",") + ")"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	ids := []int64{}
	var id int64
	if err != nil {
		return ids
	}
	rows, error := stmt.Query()
	if error != nil {
		return ids
	}
	for rows.Next() {
		rows.Scan(&id)
		ids = append(ids, id)
	}
	return ids
}

func (uc *AppointmentControllerStruct) GetAppointmentByCode(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		return
	}
	query := "SELECT appointment.id id, appointment.case_type case_type, appointment.is_vip is_vip, appointment.start_at start_at, appointment.user_id user_id, ifnull(appointment.info, '') info, appointment.income income, appointment.status appointment_status,  appointment.updated_at updated_at, user.fname user_fname, user.lname user_lname, user.id user_id, user.gender user_gender, appointment.price price from appointment LEFT JOIN user on appointment.user_id = user.id WHERE `appointment`.`status` = 1 and `appointment`.`code` = ? ORDER BY `appointment`.`id` DESC LIMIT 1"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error())
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	var appointment appointment.SimpleAppointmentInfo
	result := stmt.QueryRow(code)
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
		&appointment.Price,
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
