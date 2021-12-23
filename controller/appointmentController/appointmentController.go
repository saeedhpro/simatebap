package appointmentController

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/simateb-project/simateb-backend/controller/caseTypeController"
	"gitlab.com/simateb-project/simateb-backend/controller/organizationController"
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
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type AppointmentControllerInterface interface {
	Get(c *gin.Context)
	Create(c *gin.Context)
	GetOperationList(c *gin.Context)
	GetAppointmentList(c *gin.Context)
	GetRadiosAppointmentList(c *gin.Context)
	GetPhotosAppointmentList(c *gin.Context)
	GetResultImages(c *gin.Context)
	GetOffsAppointmentList(c *gin.Context)
	GetQueDetails(c *gin.Context)
	Update(c *gin.Context)
	ChangeStatus(c *gin.Context)
	SendResult(c *gin.Context)
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
	PRndImg := helper.RandomString(6)
	RRndImg := helper.RandomString(6)
	LRndImg := helper.RandomString(6)
	CreateAppointmentQuery := "INSERT INTO `appointment`(`user_id`, `info`, `start_at`, `case_type`, `income`, `staff_id`, `is_vip`, `code`, `p_rnd_img`, `r_rnd_img`, `l_rnd_img`, `organization_id`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)"
	stmt, err := repository.DBS.MysqlDb.Prepare(CreateAppointmentQuery)
	if err != nil {
		log.Println(err.Error())
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	log.Println(createAppointmentRequest, "sa")
	randomCode := helper.RandomString(6)
	staff := auth.GetStaffUser(c)
	_, err = stmt.Exec(
		&createAppointmentRequest.UserID,
		&createAppointmentRequest.Info,
		&createAppointmentRequest.StartAt,
		&createAppointmentRequest.CaseType,
		&createAppointmentRequest.Income,
		staff.UserID,
		&createAppointmentRequest.IsVip,
		&randomCode,
		&PRndImg,
		&RRndImg,
		&LRndImg,
		&staff.OrganizationID,
	)
	if err != nil {
		log.Println(err.Error())
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	c.JSON(http.StatusOK, true)
}

func (uc *AppointmentControllerStruct) SendResult(c *gin.Context) {
	staff := auth.GetStaffUser(c)
	id := c.Param("id")
	_, header, err := c.Request.FormFile("file")
	if err != nil {
		log.Println(err.Error(), "read file")
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	app, err := appointment.GetAppointmentById(id)
	if err != nil {
		log.Println(err.Error(), "get app")
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	newLocation := fmt.Sprintf("./images/results/%d", app.ID)
	location := fmt.Sprintf("/images/results/%d", app.ID)
	if _, err := os.Stat(newLocation); os.IsNotExist(err) {
		err := os.Mkdir(newLocation, 0755)
		if err != nil {
			log.Println(err.Error(), "make first dir")
			errorsHandler.GinErrorResponseHandler(c, err)
			return
		}
	}
	staffOrg := organizationController.GetOrganization(fmt.Sprintf("%d", staff.OrganizationID))
	prof := "/photo"
	switch staffOrg.ProfessionID {
	case "3":
		newLocation += "/radio"
		prof = "radio"
		location += "/radio"
		break
	case "1":
		newLocation += "/photo"
		prof = "photo"
		location += "/photo"
		break
	case "2":
		newLocation += "/lab"
		prof = "lab"
		location += "/lab"
		break
	}
	if _, err := os.Stat(newLocation); os.IsNotExist(err) {
		err := os.Mkdir(newLocation, 0755)
		if err != nil {
			log.Println(err.Error(), "make second dir")
			errorsHandler.GinErrorResponseHandler(c, err)
			return
		}
	}
	t := time.Now().UnixNano()
	fileName := fmt.Sprintf("%d%s", t, filepath.Ext(header.Filename))
	err = c.SaveUploadedFile(header, newLocation+"/"+fileName)
	if err != nil {
		log.Println("err", err.Error())
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	updateApp(id, prof)
	c.JSON(http.StatusAccepted, gin.H{
		"name": fileName,
		"id":   id,
		"prof": prof,
		"path": fmt.Sprintf("%s/%s", location, fileName),
	})
}

func updateApp(id string, prof string) {
	query := "UPDATE `appointment` SET "
	switch prof {
	case "radio":
		query += "`r_result_at`= ? "
		break
	case "photo":
		query += "`p_result_at`= ? "
		break
	case "lab":
		query += "`l_rnd_img`= ? "
		break
	}
	query += " WHERE id = ?"
	t := time.Now()
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error())
		return
	}
	_, err = stmt.Exec(t, id)
	if err != nil {
		log.Println(err.Error())
		return
	}

}

func (uc *AppointmentControllerStruct) GetResultImages(c *gin.Context) {
	id := c.Param("id")
	prof := c.Param("prof")
	logos := []string{}
	location := fmt.Sprintf("./images/results/%s", id)
	switch prof {
	case "3":
		location += "/radio"
		break
	case "1":
		location += "/photo"
		break
	case "2":
		location += "/lab"
		break
	}
	files, err := ioutil.ReadDir(location)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if f.Name() != "." || f.Name() != ".." {
			logos = append(logos, fmt.Sprintf("http://%s/images/results/%s/%s/%s", c.Request.Host, id, prof, f.Name()))
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"logos": logos,
	})
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
		&appointment.UserID,
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
	appointment.User, _ = user.GetUserByID(appointment.UserID)
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
	startAt := c.Query("start_date")
	endAt := c.Query("end_date")
	status := c.Query("status")
	vip := c.Query("vip")
	q := c.Query("q")
	var values []interface{}
	if status == "" {
		log.Println("status is needed!")
		c.JSON(422, gin.H{
			"message": "status is needed!",
		})
		return
	}
	staffUser := auth.GetStaffUser(c)
	organizationID := staffUser.OrganizationID
	values = append(values, organizationID)
	query := "SELECT appointment.id id, appointment.user_id user_id, appointment.start_at start_at, ifnull(appointment.info, '') info, appointment.income, ifnull(appointment.case_type, '') case_type, ifnull(appointment.code, '') code, appointment.photography_status photography_status, appointment.radiology_status radiology_status, appointment.status status, ifnull(user.fname, '') fname, ifnull(user.fname, '') lname, ifnull(user.tel, '') tel, ifnull(user.file_id, '') file_id, ifnull(user.logo, '') logo FROM `appointment` LEFT JOIN `user` on `appointment`.`user_id` = `user`.id WHERE appointment.organization_id = ? "
	if startAt != "" && startAt != "null" && startAt != "undefined" {
		query += " AND appointment.start_at >= ? "
		values = append(values, startAt)
	}
	if endAt != "" && endAt != "null" && endAt != "undefined" {
		query += " AND appointment.start_at <= ? "
		values = append(values, endAt)
	}
	query += " AND appointment.status in ("
	ss := strings.Split(status, ",")
	for i := 0; i < len(ss); i++ {
		values = append(values, ss[i])
		if i != len(ss)-1 {
			query += "?,"
		} else {
			query += "?"
		}
	}
	query += ")"
	if vip != "" {
		query += " AND vip = ? "
		values = append(values, vip)
	}
	if q != "" && q != "null" && q != "undefined" {
		query += " AND (user.fname LIKE '%" + q + "%' OR user.lname LIKE '%" + q + "%' OR appointment.code LIKE '%" + q + "%')"
	}
	query += " ORDER BY `id` DESC"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error(), "first")
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	var rows *sql.Rows
	rows, err = stmt.Query(values...)
	if err != nil {
		log.Println(err.Error(), "second")
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	operations := []appointment.OperationInfo{}
	var operation appointment.OperationInfo
	for rows.Next() {
		err = rows.Scan(
			&operation.ID,
			&operation.UserID,
			&operation.StartAt,
			&operation.Info,
			&operation.Income,
			&operation.CaseType,
			&operation.Code,
			&operation.PhotographyStatus,
			&operation.RadiologyStatus,
			&operation.Status,
			&operation.FName,
			&operation.LName,
			&operation.Tel,
			&operation.FileID,
			&operation.Logo,
		)
		if err != nil {
			log.Println(err.Error(), "third")
			errorsHandler.GinErrorResponseHandler(c, err)
			return
		}
		operation.User, _ = user.GetUserByID(operation.UserID)
		operations = append(operations, operation)
	}
	c.JSON(http.StatusOK, operations)
}
func (uc *AppointmentControllerStruct) GetRadiosAppointmentList(c *gin.Context) {
	var values []interface{}
	organizationID := c.Param("id")
	page := c.Query("page")
	if page != "" {
		page = "1"
	}
	values = append(values, organizationID, page)
	query := "SELECT appointment.id id, appointment.user_id user_id, ifnull(appointment.prescription, '') prescription,ifnull(appointment.radiology_cases, '') radiology_cases,ifnull(appointment.photography_cases, '') photography_cases, appointment.start_at start_at, ifnull(appointment.info, '') info, appointment.income, ifnull(appointment.case_type, '') case_type, ifnull(appointment.code, '') code, appointment.status status, ifnull(user.fname, '') fname, ifnull(user.lname, '') lname, ifnull(user.tel, '') tel, ifnull(user.file_id, '') file_id FROM `appointment` LEFT JOIN `user` on `appointment`.`user_id` = `user`.id WHERE appointment.radiology_id = ? AND appointment.status = 2 AND appointment.r_result_at IS NULL LIMIT 10 OFFSET ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error(), "first")
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	var rows *sql.Rows
	rows, err = stmt.Query(values...)
	if err != nil {
		log.Println(err.Error(), "second")
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	paginated := pagination.SendDocAppointmentPaginationInfo{}
	apps := []appointment.UserAppointmentInfo{}
	var app appointment.UserAppointmentInfo
	for rows.Next() {
		err = rows.Scan(
			&app.ID,
			&app.UserID,
			&app.Prescription,
			&app.RadiologyCases,
			&app.PhotographyCases,
			&app.StartAt,
			&app.Info,
			&app.Income,
			&app.CaseType,
			&app.Code,
			&app.Status,
			&app.FName,
			&app.LName,
			&app.Tel,
			&app.FileID,
		)
		if err != nil {
			log.Println(err.Error(), "third")
			errorsHandler.GinErrorResponseHandler(c, err)
			return
		}
		apps = append(apps, app)
	}
	paginated.Data = apps
	count := getRadiosPageCount(organizationID)
	paginated.PagesCount = count
	paginated.HasNextPage = count/10 > 10
	c.JSON(http.StatusOK, paginated)
}
func getRadiosPageCount(organizationID string) int {
	count := 0
	var values []interface{}
	values = append(values, organizationID)
	query := "SELECT COUNT(*) `count` FROM `appointment` LEFT JOIN `user` on `appointment`.`user_id` = `user`.id WHERE appointment.radiology_id = ? AND appointment.status = 2 AND appointment.r_result_at IS NULL ORDER BY appointment.`id` DESC"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error(), "first")
		return count
	}
	row := stmt.QueryRow(values...)
	err = row.Scan(&count)
	if err != nil {
		log.Println(err.Error(), "second")
		return count
	}
	return count
}
func getPhotosPageCount(organizationID string) int {
	count := 0
	var values []interface{}
	values = append(values, organizationID)
	query := "SELECT COUNT(*) `count` FROM `appointment` LEFT JOIN `user` on `appointment`.`user_id` = `user`.id WHERE appointment.photography_id = ? AND appointment.status = 2 AND appointment.p_result_at IS NULL ORDER BY appointment.`id` DESC"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error(), "first")
		return count
	}
	row := stmt.QueryRow(values...)
	err = row.Scan(&count)
	if err != nil {
		log.Println(err.Error(), "second")
		return count
	}
	return count
}
func getOffsPageCount(organizationID string) int {
	count := 0
	var values []interface{}
	values = append(values, organizationID)
	query := "SELECT COUNT(*) `count` FROM `appointment` LEFT JOIN `user` on `appointment`.`user_id` = `user`.id WHERE appointment.office_id = ? AND appointment.status = 2 AND appointment.p_result_at IS NULL ORDER BY appointment.`id` DESC"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error(), "first")
		return count
	}
	row := stmt.QueryRow(values...)
	err = row.Scan(&count)
	if err != nil {
		log.Println(err.Error(), "second")
		return count
	}
	return count
}

func (uc *AppointmentControllerStruct) GetPhotosAppointmentList(c *gin.Context) {
	var values []interface{}
	organizationID := c.Param("id")
	page := c.Query("page")
	values = append(values, organizationID, page)
	query := "SELECT appointment.id id, appointment.user_id user_id, ifnull(appointment.prescription, '') prescription,appointment.radiology_status, appointment.photography_status,ifnull(appointment.radiology_cases, '') radiology_cases,ifnull(appointment.photography_cases, '') photography_cases, appointment.start_at start_at, ifnull(appointment.info, '') info, appointment.income, ifnull(appointment.case_type, '') case_type, ifnull(appointment.code, '') code, appointment.status status, ifnull(user.fname, '') fname, ifnull(user.lname, '') lname, ifnull(user.tel, '') tel, ifnull(user.file_id, '') file_id FROM `appointment` LEFT JOIN `user` on `appointment`.`user_id` = `user`.id WHERE appointment.photography_id = ? AND appointment.status = 2 AND appointment.p_result_at IS NULL LIMIT 10 OFFSET ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error(), "first")
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	var rows *sql.Rows
	rows, err = stmt.Query(values...)
	paginated := pagination.SendDocAppointmentPaginationInfo{}
	apps := []appointment.UserAppointmentInfo{}
	var app appointment.UserAppointmentInfo
	for rows.Next() {
		err = rows.Scan(
			&app.ID,
			&app.UserID,
			&app.Prescription,
			&app.RadiologyStatus,
			&app.PhotographyStatus,
			&app.RadiologyCases,
			&app.PhotographyCases,
			&app.StartAt,
			&app.Info,
			&app.Income,
			&app.CaseType,
			&app.Code,
			&app.Status,
			&app.FName,
			&app.LName,
			&app.Tel,
			&app.FileID,
		)
		if err != nil {
			log.Println(err.Error(), "third")
			errorsHandler.GinErrorResponseHandler(c, err)
			return
		}
		apps = append(apps, app)
	}
	paginated.Data = apps
	count := getPhotosPageCount(organizationID)
	paginated.PagesCount = count
	paginated.HasNextPage = count/10 > 10
	c.JSON(http.StatusOK, paginated)
}
func (uc *AppointmentControllerStruct) GetOffsAppointmentList(c *gin.Context) {
	var values []interface{}
	organizationID := c.Param("id")
	page := c.Query("page")
	if page != "" {
		page = "1"
	}
	values = append(values, organizationID, page)
	query := "SELECT appointment.id id, appointment.user_id user_id, ifnull(appointment.prescription, '') prescription,ifnull(appointment.radiology_cases, '') radiology_cases,ifnull(appointment.photography_cases, '') photography_cases, appointment.start_at start_at, ifnull(appointment.info, '') info, appointment.income, ifnull(appointment.case_type, '') case_type, ifnull(appointment.code, '') code, appointment.status status, ifnull(user.fname, '') fname, ifnull(user.lname, '') lname, ifnull(user.tel, '') tel, ifnull(user.file_id, '') file_id FROM `appointment` LEFT JOIN `user` on `appointment`.`user_id` = `user`.id WHERE appointment.office_id = ? AND appointment.status = 2 AND appointment.p_result_at IS NOT NULL LIMIT 10 OFFSET ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error(), "first")
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	var rows *sql.Rows
	rows, err = stmt.Query(values...)
	paginated := pagination.SendDocAppointmentPaginationInfo{}
	apps := []appointment.UserAppointmentInfo{}
	var app appointment.UserAppointmentInfo
	for rows.Next() {
		err = rows.Scan(
			&app.ID,
			&app.UserID,
			&app.Prescription,
			&app.RadiologyCases,
			&app.PhotographyCases,
			&app.StartAt,
			&app.Info,
			&app.Income,
			&app.CaseType,
			&app.Code,
			&app.Status,
			&app.FName,
			&app.LName,
			&app.Tel,
			&app.FileID,
		)
		if err != nil {
			log.Println(err.Error(), "third")
			errorsHandler.GinErrorResponseHandler(c, err)
			return
		}
		apps = append(apps, app)
	}
	paginated.Data = apps
	count := getOffsPageCount(organizationID)
	paginated.PagesCount = count
	paginated.HasNextPage = count/10 > 10
	c.JSON(http.StatusOK, paginated)
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
	id := c.Param("id")
	if id == "" {
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
	values = append(values, id)
	_, error := stmt.Exec(values...)
	if error != nil {
		log.Println(error.Error())
		errorsHandler.GinErrorResponseHandler(c, error)
		return
	}
	InsertLastAppointmentPrescription(id, request.FuturePrescription, request.Prescription)
	staff := auth.GetStaffUser(c)
	calcAppointmentNewPrice(id, request.PhotographyCases, request.RadiologyCases, staff.OrganizationID)
	c.JSON(200, true)
}
func calcAppointmentNewPrice(id string, photos []string, radios []string, orgID int64) {
	lastPhotos, lastRadios, price := getLastPhotosRadios(id)
	newPhotos := getNewFromArray(lastPhotos, photos)
	newRadios := getNewFromArray(lastRadios, radios)
	newPrice := calcNewPrice(newPhotos, newRadios, orgID)
	price += newPrice
	ps := getAllItems(lastPhotos, newPhotos)
	rs := getAllItems(lastRadios, newRadios)
	updateAppPrice(id, price, ps, rs)
}

func getAllItems(last string, news []string) string {
	items := ""
	if len(last) > 0 {
		items = fmt.Sprintf("%s,%s", last, strings.Join(news, ","))
	}
	return items
}

func updateAppPrice(id string, price float64, photos string, radios string) {
	var query = "UPDATE `appointment` SET last_radiology_cases = ? , last_photography_cases = ?, price = ? WHERE id = ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		return
	}
	_, err = stmt.Exec(
		radios,
		photos,
		price,
		id,
	)

}

func calcNewPrice(photos []string, radios []string, orgID int64) float64 {
	price := 0.0
	if len(radios) > 0 {
		price += calcRadiologyCasesPrice(radios, orgID)
	}
	if len(photos) > 0 {
		price += calcPhotographyCasesPrice(photos, orgID)
	}
	return price
}

func getNewFromArray(last string, array []string) []string {
	lastArray := strings.Split(last, ",")
	tempArray := []string{}
	for i := 0; i < len(array); i++ {
		if stringInSlice(array[i], lastArray) {
			tempArray = append(tempArray, array[i])
		}
	}
	return tempArray
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func getLastPhotosRadios(id string) (string, string, float64) {
	lastPhotos := ""
	lastRadios := ""
	price := 0.0
	query := "SELECT last_radiology_cases, last_photography_cases, price FORM appointment WHERE id = ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error())
		return lastPhotos, lastRadios, price
	}
	result := stmt.QueryRow(id)
	if err = result.Err(); err != nil {
		log.Println(err.Error())
		return lastPhotos, lastRadios, price
	}
	err = result.Scan(
		lastRadios,
		lastPhotos,
		price,
	)
	if err = result.Err(); err != nil {
		log.Println(err.Error())
		return lastPhotos, lastRadios, price
	}
	return lastPhotos, lastRadios, price
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
	query := "SELECT appointment.id id, appointment.status `status`, ifnull(appointment.case_type, '') case_type, appointment.start_at start_at, appointment.price price, appointment.user_id user_id FROM appointment LEFT JOIN user ON appointment.user_id = user.id"
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	status := c.Query("status")
	vip := c.Query("vip")
	q := c.Query("q")
	page := c.Query("page")
	if page == "" {
		page = "1"
	}
	var values []interface{}
	var queries []string
	if startDate != "" {
		startDate = fmt.Sprintf("%s 00:00:00", startDate)
		values = append(values, startDate)
		queries = append(queries, " appointment.start_at >= ? ")
	}
	if endDate != "" {
		endDate = fmt.Sprintf("%s 23:59:59", endDate)
		values = append(values, endDate)
		queries = append(queries, " appointment.start_at <= ? ")
	}
	if vip != "" {
		values = append(values, vip)
		queries = append(queries, " appointment.is_vip = ? ")
	}
	if q != "" && q != "null" && q != "undefined" {
		q = "'%" + strings.TrimSpace(q) + "%'"
		queries = append(queries, fmt.Sprintf(" user.fname LIKE %s OR user.lname LIKE %s", q, q))
	}
	if status != "" {
		values = append(values, status)
		queries = append(queries, " appointment.status IN (?) ")
	}
	staff := auth.GetStaffUser(c)
	values = append(values, staff.OrganizationID)
	queries = append(queries, " user.organization_id = ? ")
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
	log.Println(err, "error")
	for rows.Next() {
		err := rows.Scan(
			&appointmentInfo.ID,
			&appointmentInfo.Status,
			&appointmentInfo.CaseType,
			&appointmentInfo.StartAt,
			&appointmentInfo.Price,
			&appointmentInfo.UserID,
		)
		if err != nil {
			log.Println("err: ", err.Error())
			c.JSON(500, paginationInfo)
			return
		}
		appointmentInfo.User, _ = user.GetUserByID(appointmentInfo.UserID)
		appointments = append(appointments, appointmentInfo)
	}
	count := 0
	getCountQuery := "SELECT COUNT('*') FROM appointment LEFT JOIN user ON appointment.user_id = user.id "
	if q != "" && q != "null" && q != "undefined" {
		getCountQuery += " AND (user.fname LIKE '%" + q + "%' OR user.lname LIKE '%" + q + "%') "
	}
	var values2 []interface{}
	var queries2 []string
	if startDate != "" {
		values2 = append(values2, startDate)
		queries2 = append(queries2, " appointment.start_at > ? ")
	}
	if endDate != "" {
		values2 = append(values2, endDate)
		queries2 = append(queries2, " appointment.start_at <= ? ")
	}
	if vip != "" {
		values2 = append(values2, vip)
		queries2 = append(queries2, " appointment.is_vip = ? ")
	}
	if q != "" && q != "null" && q != "undefined" {
		q = "'%" + strings.TrimSpace(q) + "%'"
		queries2 = append(queries2, fmt.Sprintf(" user.fname LIKE %s OR user.lname LIKE %s", q, q))
	}
	if status != "" {
		values2 = append(values2, status)
		queries2 = append(queries2, " appointment.status IN (?) ")
	}
	where2 := strings.Join(queries2, " AND ")
	if where2 != "" {
		where2 = fmt.Sprintf("%s %s", "WHERE", where2)
	}
	getCountQuery = fmt.Sprintf("%s %s", getCountQuery, where2)
	stmt, err = repository.DBS.MysqlDb.Prepare(getCountQuery)
	if err != nil {
		log.Println(err.Error(), "count error")
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	result := stmt.QueryRow(values2...)
	err = result.Scan(&count)
	if err != nil {
		log.Println(err.Error(), "count")
	}
	p, _ := strconv.Atoi(page)
	paginationInfo.HasNextPage = count > 10 && count > (p*10)
	paginationInfo.HasPreviousPage = count > 10 && p > 1
	paginationInfo.PagesCount = count
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
	staff := auth.GetStaffUser(c)
	code := helper.RandomString(6)
	price := 0.0
	if len(request.RadiologyCases) > 0 {
		price += calcRadiologyCasesPrice(request.RadiologyCases, staff.OrganizationID)
	}
	if len(request.PhotographyCases) > 0 {
		price += calcPhotographyCasesPrice(request.PhotographyCases, staff.OrganizationID)
	}
	status := acceptAppointment(id, code, price, request)
	c.JSON(200, status)
}

func acceptAppointment(id string, code string, price float64, request appointment.AcceptAppointmentRequest) bool {
	query := "UPDATE `appointment` SET `status`= ?, `photography_cases`= ? ,`radiology_cases`= ? ," +
		"`last_photography_cases`= ? ,`last_radiology_cases`= ? ,`prescription`= ? ," +
		"`future_prescription`= ? , `photography_msg`= ?, `radiology_msg`= ?, `radiology_id`= ?," +
		"`photography_id`= ?, `code`= ?, price = ? WHERE id = ? "
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		return false
	}
	_, err = stmt.Exec(
		2,
		strings.Join(request.PhotographyCases, ","),
		strings.Join(request.RadiologyCases, ","),
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
	if err != nil {
		log.Println(err.Error())
		return false
	}
	if request.Prescription != "" || request.FuturePrescription != "" {
		InsertLastAppointmentPrescription(id, request.FuturePrescription, request.Prescription)
	}
	return true
}

func InsertLastAppointmentPrescription(id string, futurePrescription string, prescription string) bool {
	count := getAppointmentLastPrescription(id)
	if count > 0 {
		return false
	}
	query := "UPDATE `last_prescription` SET `future_prescription`= ?,`prescription`= ? WHERE `appointment_id` = ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		return false
	}
	_, error := stmt.Exec(
		futurePrescription,
		prescription,
		id,
	)
	if error != nil {
		log.Println(error.Error())
		return false
	}
	return true
}

func getAppointmentLastPrescription(id string) int {
	count := 0
	query := "SELECT COUNT(*) count FROM `last_prescription` WHERE appointment_id = ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		return count
	}
	result := stmt.QueryRow(id)
	if result.Err() != nil {
		return count
	}
	err = result.Scan(count)
	if err != nil {
		return count
	}
	return count
}

func calcRadiologyCasesPrice(cases []string, orgID int64) float64 {
	price := 0.0
	ids := getRadiologyCasesID(cases)
	idsStr := strings.Trim(strings.Replace(fmt.Sprint(ids), " ", ",", -1), "[]")
	price += getRadiologyCasesPrices(idsStr, orgID)
	return price
}

func calcPhotographyCasesPrice(cases []string, orgID int64) float64 {
	price := 0.0
	ids := getPhotographyCasesID(cases)
	idsStr := strings.Trim(strings.Replace(fmt.Sprint(ids), " ", ",", -1), "[]")
	price += getPhotographyCasesPrices(idsStr, orgID)
	return price
}

func getRadiologyCasesPrices(ids string, orgID int64) float64 {
	price := 0.0
	query := "SELECT SUM(price) price FROM `case_prices` WHERE id in (" + ids + ") and organization_id = ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		return price
	}
	rows, error := stmt.Query(orgID)
	if error != nil {
		return price
	}
	for rows.Next() {
		rows.Scan(&price)
	}
	return price
}
func getPhotographyCasesPrices(ids string, orgID int64) float64 {
	price := 0.0
	query := "SELECT SUM(price) price FROM `case_prices` WHERE id in (" + ids + ") and organization_id = ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		return price
	}
	rows, error := stmt.Query(orgID)
	if error != nil {
		return price
	}
	for rows.Next() {
		rows.Scan(&price)
	}
	return price
}

func getRadiologyCasesID(cases []string) []int64 {
	query := "SELECT cases.id, cases.name FROM cases WHERE cases.val in (" + strings.Join(cases, ",") + ") AND profession_id = 3"
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
func getPhotographyCasesID(cases []string) []int64 {
	query := "SELECT cases.id, cases.name FROM cases WHERE cases.val in (" + strings.Join(cases, ",") + ") AND profession_id = 1"
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
		&appointment.UserID,
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
	appointment.User, _ = user.GetUserByID(appointment.UserID)
	c.JSON(http.StatusOK, appointment)
}
