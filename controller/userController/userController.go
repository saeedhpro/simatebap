package userController

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/simateb-project/simateb-backend/controller/organizationController"
	appointment2 "gitlab.com/simateb-project/simateb-backend/domain/appointment"
	"gitlab.com/simateb-project/simateb-backend/domain/organization"
	sms2 "gitlab.com/simateb-project/simateb-backend/domain/sms"
	wallet2 "gitlab.com/simateb-project/simateb-backend/domain/wallet"
	"gitlab.com/simateb-project/simateb-backend/helper"
	"gitlab.com/simateb-project/simateb-backend/repository"
	"gitlab.com/simateb-project/simateb-backend/repository/doc"
	"gitlab.com/simateb-project/simateb-backend/repository/medicalHistory"
	mysqlQuery "gitlab.com/simateb-project/simateb-backend/repository/mysqlQuery/auth"
	"gitlab.com/simateb-project/simateb-backend/repository/user"
	"gitlab.com/simateb-project/simateb-backend/repository/wallet"
	"gitlab.com/simateb-project/simateb-backend/utils/auth"
	"gitlab.com/simateb-project/simateb-backend/utils/errorsHandler"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type UserControllerInterface interface {
	Create(c *gin.Context)
	CreateUserWalletHistories(c *gin.Context)
	AcceptOrRejectWalletHistories(c *gin.Context)
	GetWalletHistoriesForAdmin(c *gin.Context)
	Own(c *gin.Context)
	Get(c *gin.Context)
	GetList(c *gin.Context)
	GetListForAdmin(c *gin.Context)
	GetAdminList(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	DeleteItems(c *gin.Context)
	CreateCode(c *gin.Context)
	SendDoc(c *gin.Context)
	DeleteDoc(c *gin.Context)
	UpdateDoc(c *gin.Context)
	GetUserDocs(c *gin.Context)
	ChangePassword(c *gin.Context)
	GetUserAppointmentList(c *gin.Context)
	GetUserAppointmentResultList(c *gin.Context)
	GetOrthodonticsList(c *gin.Context)
	CreateOrthodonticsList(c *gin.Context)
	GetUserWallet(c *gin.Context)
	GetUserWalletHistories(c *gin.Context)
	GetUserWalletHistoriesSum(c *gin.Context)
	GetUserWalletHistoriesDays(c *gin.Context)
	IncreaseUserWallet(c *gin.Context)
	DecreaseUserWallet(c *gin.Context)
	SetUserWallet(c *gin.Context)
	GetLastLoginUsers(c *gin.Context)
	GetLastLoginPatients(c *gin.Context)
}

type UserControllerStruct struct {
}

func NewUserController() UserControllerInterface {
	x := &UserControllerStruct{
	}
	return x
}

func (uc *UserControllerStruct) Create(c *gin.Context) {
	var createUserRequest organization.CreateUserRequest
	if errors := c.ShouldBindJSON(&createUserRequest); errors != nil {
		log.Println(errors.Error())
		errorsHandler.GinErrorResponseHandler(c, errors)
		return
	}
	stmt, err := repository.DBS.MysqlDb.Prepare(mysqlQuery.CreateUserQuery)
	if err != nil {
		log.Println(err.Error())
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	defer stmt.Close()
	staffID := auth.GetStaffUser(c).UserID
	password, err := auth.PasswordHash(createUserRequest.Password)
	if err != nil {
		log.Println(err.Error())
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	createUserRequest.Password = password
	result, err := stmt.Exec(
		createUserRequest.FirstName,
		createUserRequest.LastName,
		createUserRequest.Email,
		createUserRequest.Info,
		createUserRequest.Description,
		createUserRequest.FileID,
		createUserRequest.Gender,
		staffID,
		createUserRequest.UserGroupId,
		createUserRequest.OrganizationId,
		createUserRequest.Tel,
		createUserRequest.Tel1,
		createUserRequest.Nid,
		createUserRequest.BirthDate,
		createUserRequest.Address,
		createUserRequest.Introducer,
		createUserRequest.Password,
		createUserRequest.Relation,
		createUserRequest.Logo,
	)
	if err != nil {
		log.Println(err.Error())
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Println(err.Error())
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	var user, _ = user.GetUserByID(id)
	c.JSON(http.StatusOK, user)
}

func (uc *UserControllerStruct) CreateUserWalletHistories(c *gin.Context) {
	var request user.CreateWithdrawRequest
	if errors := c.ShouldBindJSON(&request); errors != nil {
		log.Println(errors.Error())
		errorsHandler.GinErrorResponseHandler(c, errors)
		return
	}
	result, _ := user.CreateWithdraw(request)
	c.JSON(200, result)
}

func (uc *UserControllerStruct) AcceptOrRejectWalletHistories(c *gin.Context) {
	var request user.AcceptOrRejectWithdrawRequest
	if errors := c.ShouldBindJSON(&request); errors != nil {
		log.Println(errors.Error())
		errorsHandler.GinErrorResponseHandler(c, errors)
		return
	}
	result, _ := user.AcceptOrRejectWithdraw(request)
	c.JSON(200, result)
}

func (uc *UserControllerStruct) GetWalletHistoriesForAdmin(c *gin.Context) {
	page := c.Query("page")
	if page == "" {
		 page = "1"
	}
	result, _ := user.GetWalletHistoriesForAdmin(page)
	c.JSON(200, result)
}

func (uc *UserControllerStruct) Own(c *gin.Context) {
	staff := auth.GetStaffUser(c)
	var user, err = user.GetUserByID(staff.UserID)
	if err != nil {
		log.Println(err.Error())
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, "یافت نشد")
			return
		}
		errorsHandler.GinErrorResponseHandler(c, err)
	}
	c.JSON(http.StatusOK, user)
}

func (uc *UserControllerStruct) Get(c *gin.Context) {
	userId := c.Param("id")
	if userId == "" {
		return
	}
	var uid, _ = strconv.ParseInt(userId, 10, 64)
	var user, err = user.GetUserByID(uid)
	if err != nil {
		log.Println(err.Error())
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, "یافت نشد")
			return
		}
		errorsHandler.GinErrorResponseHandler(c, err)
	}
	c.JSON(http.StatusOK, user)
}

func (uc *UserControllerStruct) GetListForAdmin(c *gin.Context) {
	group := c.Query("group")
	if group == "" {
		group = "1,2,3,4,5,6"
	}
	getUsersQuery := "SELECT user.id id, ifnull(user.fname, '') fname, ifnull(user.lname, '') lname, user.last_login last_login, user.created created, user.tel tel, user.user_group_id user_group_id, user_group.name user_group_name, user.birth_date birth_date, organization.id organization_id, organization.name organization_name, ifnull(user.relation, '') relation, ifnull(user.description, '') description FROM (user LEFT JOIN organization ON user.organization_id = organization.id) LEFT JOIN user_group ON user.user_group_id = user_group.id WHERE user.user_group_id IN (" + group + ")"
	page := c.Query("page")
	if page == "" {
		page = "1"
	}
	q := c.Query("q")
	if q != "" && q != "null" && q != "undefined" {
		getUsersQuery += " AND (user.fname LIKE '%" + q + "%' OR user.lname LIKE '%" + q + "%') "
	}
	offset, err := strconv.Atoi(page)
	offset = (offset - 1) * 10
	getUsersQuery += " ORDER BY `user`.`id` DESC LIMIT 10 OFFSET ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(getUsersQuery)
	if err != nil {
		log.Println(err.Error(), "users error")
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	staffUser := auth.GetStaffUser(c)
	if staffUser == nil {
		log.Println(err.Error(), "staff error")
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	users := []organization.OrganizationUser{}
	var user organization.OrganizationUser
	rows, err := stmt.Query(offset)
	if err != nil {
		log.Println(err.Error())
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.LastLogin,
			&user.Created,
			&user.Tel,
			&user.UserGroupID,
			&user.UserGroupName,
			&user.BirthDate,
			&user.OrganizationID,
			&user.OrganizationName,
			&user.Relation,
			&user.Description,
		)
		if err != nil {
			log.Println(err.Error())
		}
		users = append(users, user)
	}
	count := 0
	getCountQuery := "SELECT COUNT('*') count FROM (user LEFT JOIN organization ON user.organization_id = organization.id) LEFT JOIN user_group ON user.user_group_id = user_group.id WHERE user.user_group_id IN (?) "
	if q != "" && q != "null" && q != "undefined" {
		getCountQuery += " AND (user.fname LIKE '%" + q + "%' OR user.lname LIKE '%" + q + "%') "
	}
	getCountQuery += " ORDER BY `user`.`last_login` DESC"
	stmt, err = repository.DBS.MysqlDb.Prepare(getCountQuery)
	if err != nil {
		log.Println(err.Error(), "count error")
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	result := stmt.QueryRow(group)
	result.Scan(&count)
	p, _ := strconv.Atoi(page)
	paginated := organization.OrganizationUserPaginate{
		Data:       users,
		Page:       p,
		PagesCount: count,
	}
	c.JSON(http.StatusOK, paginated)
}

func (uc *UserControllerStruct) GetAdminList(c *gin.Context) {
	getUsersQuery := "SELECT user.id id, user.fname fname, user.lname lname, user.last_login last_login, user.created created, user.tel tel, user.user_group_id user_group_id, user_group.name user_group_name, user.birth_date birth_date, organization.id organization_id, organization.name organization_name, ifnull(user.relation, '') relation, ifnull(user.description, '') description FROM (user LEFT JOIN organization ON user.organization_id = organization.id) LEFT JOIN user_group ON user.user_group_id = user_group.id WHERE user_group_id = 2"
	page := c.Query("page")
	if page == "" {
		page = "1"
	}
	q := c.Query("q")
	if q != "" && q != "null" && q != "undefined" {
		getUsersQuery += " user.fname LIKE '%" + q + "%' "
	}
	offset, err := strconv.Atoi(page)
	offset = (offset - 1) * 10
	getUsersQuery += " ORDER BY id DESC"
	getUsersQuery += " LIMIT 10 OFFSET ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(getUsersQuery)
	if err != nil {
		log.Println(err.Error())
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	var users []organization.OrganizationUser
	var user organization.OrganizationUser
	rows, err := stmt.Query(offset)
	defer rows.Close()
	if err != nil {
		log.Println(err.Error())
		return
	}
	for rows.Next() {
		err = rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.LastLogin,
			&user.Created,
			&user.Tel,
			&user.UserGroupID,
			&user.UserGroupName,
			&user.BirthDate,
			&user.OrganizationID,
			&user.OrganizationName,
			&user.Relation,
			&user.Description,
		)
		if err != nil {
			log.Println(err.Error())
		}
		users = append(users, user)
	}
	c.JSON(http.StatusOK, users)
}

func (uc *UserControllerStruct) GetList(c *gin.Context) {
	getUsersQuery := "SELECT user.id id, user.fname fname, user.lname lname, user.last_login last_login, user.created created, user.tel tel, user.user_group_id user_group_id, user_group.name user_group_name, user.birth_date birth_date, organization.id organization_id, organization.name organization_name, ifnull(user.relation, '') relation, ifnull(user.description, '') description FROM (user LEFT JOIN organization ON user.organization_id = organization.id) LEFT JOIN user_group ON user.user_group_id = user_group.id WHERE user.organization_id = ? AND user.user_group_id != 2 "
	page := c.Query("page")
	if page == "" {
		page = "1"
	}
	q := c.Query("q")
	if q != "" && q != "null" && q != "undefined" {
		getUsersQuery += " LIKE '%" + q + "%' "
	}
	offset, err := strconv.Atoi(page)
	offset = (offset - 1) * 10
	getUsersQuery += " ORDER BY id DESC"
	getUsersQuery += " LIMIT 10 OFFSET ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(getUsersQuery)
	if err != nil {
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	staffUser := auth.GetStaffUser(c)
	if staffUser == nil {
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	var users []organization.OrganizationUser
	var user organization.OrganizationUser
	rows, err := stmt.Query(staffUser.OrganizationID, offset)
	defer rows.Close()
	if err != nil {
		log.Println(err.Error())
		return
	}
	for rows.Next() {
		err = rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.LastLogin,
			&user.Created,
			&user.Tel,
			&user.UserGroupID,
			&user.UserGroupName,
			&user.BirthDate,
			&user.OrganizationID,
			&user.OrganizationName,
			&user.Relation,
			&user.Description,
		)
		if user.BirthDate != nil {
			year, _, _, _, _, _ := helper.TimeDiff(user.BirthDate.Time, time.Now())
			user.Birth = year
		}
		if err != nil {
			log.Println(err.Error())
		}
		users = append(users, user)
	}
	count := 0
	getCountQuery := "SELECT COUNT('*') count FROM (user LEFT JOIN organization ON user.organization_id = organization.id) LEFT JOIN user_group ON user.user_group_id = user_group.id WHERE user.organization_id = ? AND user.user_group_id != 2 "
	if q != "" && q != "null" && q != "undefined" {
		getCountQuery += " LIKE '%" + q + "%' "
	}
	stmt, err = repository.DBS.MysqlDb.Prepare(getCountQuery)
	if err != nil {
		log.Println(err.Error(), "count error")
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	result := stmt.QueryRow(staffUser.OrganizationID)
	err = result.Scan(&count)
	if err != nil {
		log.Println(err.Error(), "err")
	}
	p, _ := strconv.Atoi(page)
	paginated := organization.OrganizationUserPaginate{
		Data:       users,
		Page:       p,
		PagesCount: count,
	}
	c.JSON(http.StatusOK, paginated)
}

func (uc *UserControllerStruct) GetLastLoginUsers(c *gin.Context) {
	getUsersQuery := "SELECT user.id user_id, user.fname user_fname, user.lname user_lname, user.tel tel, user.last_login last_login, user_group.name user_group_name from user left join user_group on user.user_group_id = user_group.id WHERE user.organization_id = ? AND user.last_login IS NOT NULL ORDER by user.last_login DESC LIMIT 10"
	stmt, err := repository.DBS.MysqlDb.Prepare(getUsersQuery)
	if err != nil {
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	staffUser := auth.GetStaffUser(c)
	if staffUser == nil {
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	users := []organization.LastLoginUser{}
	var user organization.LastLoginUser
	rows, err := stmt.Query(staffUser.OrganizationID)
	defer rows.Close()
	if err != nil {
		log.Println(err.Error())
		return
	}
	for rows.Next() {
		err = rows.Scan(
			&user.ID,
			&user.UserFirstName,
			&user.UserLastName,
			&user.Tel,
			&user.LastLogin,
			&user.UserGroupName,
		)
		if err != nil {
			log.Println(err.Error())
		}
		users = append(users, user)
	}
	c.JSON(http.StatusOK, users)
}
func (uc *UserControllerStruct) GetLastLoginPatients(c *gin.Context) {
	getUsersQuery := "SELECT user.id user_id, user.fname user_fname, user.lname user_lname, user.tel tel, user.last_login last_login, organization.id user_organization_id, organization.name user_organization_name from user left join organization on user.organization_id = organization.id WHERE user.organization_id = ? AND user.last_login IS NOT NULL ORDER by user.last_login DESC LIMIT 10"
	stmt, err := repository.DBS.MysqlDb.Prepare(getUsersQuery)
	if err != nil {
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	staffUser := auth.GetStaffUser(c)
	if staffUser == nil {
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	users := []organization.LastLoginUser{}
	var user organization.LastLoginUser
	rows, err := stmt.Query(staffUser.OrganizationID)
	defer rows.Close()
	if err != nil {
		log.Println(err.Error())
		return
	}
	for rows.Next() {
		err = rows.Scan(
			&user.ID,
			&user.UserFirstName,
			&user.UserLastName,
			&user.Tel,
			&user.LastLogin,
			&user.OrganizationID,
			&user.OrganizationName,
		)
		if err != nil {
			log.Println(err.Error())
		}
		users = append(users, user)
	}
	c.JSON(http.StatusOK, users)
}

func (uc *UserControllerStruct) GetUserAppointmentList(c *gin.Context) {
	query := "SELECT id, user_id, created_at, ifnull(info, '') info, staff_id, start_at, end_at, status, ifnull(director_id, -1) director_id, updated_at, income, ifnull(subject, '') subject, ifnull(case_type, '') case_type, ifnull(laboratory_cases, '') laboratory_cases, ifnull(photography_cases, '') photography_cases, ifnull(radiology_cases, '') radiology_cases, ifnull(prescription, '') prescription, ifnull(future_prescription, '') future_prescription, ifnull(laboratory_msg, '') laboratory_msg, ifnull(photography_msg, '') photography_msg, ifnull(radiology_msg, '') radiology_msg, organization_id, ifnull(director_id, -1) laboratory_id, ifnull(photography_id, -1) photography_id, ifnull(radiology_id, -1) radiology_id, l_admission_at, r_admission_at, p_admission_at, l_result_at, r_result_at, p_result_at, ifnull(l_rnd_img, '') l_rnd_img, ifnull(r_rnd_img, '') r_rnd_img, ifnull(p_rnd_img, '') p_rnd_img, l_imgs, r_imgs, p_imgs, ifnull(code, '') code, is_vip, ifnull(vip_introducer, 0) vip_introducer, absence, ifnull(file_id, '') file_id FROM appointment WHERE user_id = ? "
	userID := c.Param("id")
	staffUser := auth.GetStaffUser(c)
	if staffUser == nil {
		errorsHandler.GinErrorResponseHandler(c, nil)
		return
	}
	staffOrg := organizationController.GetOrganization(fmt.Sprintf("%d", staffUser.OrganizationID))
	if staffOrg == nil {
		errorsHandler.GinErrorResponseHandler(c, nil)
		return
	}
	if staffOrg.ProfessionID == "1" {
		query += " AND photography_id = ? "
	} else if staffOrg.ProfessionID == "2" {
		query += " AND laboratory_id = ? "
	} else if staffOrg.ProfessionID == "3" {
		query += " AND radiology_id = ? "
	} else {
		query += " AND organization_id = ? "
	}
	appointments := []appointment2.UserAppointmentInfo{}
	var appointment appointment2.UserAppointmentInfo
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error(), "prepare")
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	rows, err := stmt.Query(userID, staffUser.OrganizationID)
	if err != nil {
		log.Println(err.Error())
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	for rows.Next() {
		err = rows.Scan(
			&appointment.ID,
			&appointment.UserID,
			&appointment.CreatedAt,
			&appointment.Info,
			&appointment.StaffID,
			&appointment.StartAt,
			&appointment.EndAt,
			&appointment.Status,
			&appointment.DirectorID,
			&appointment.UpdatedAt,
			&appointment.Income,
			&appointment.Subject,
			&appointment.CaseType,
			&appointment.LaboratoryCases,
			&appointment.PhotographyCases,
			&appointment.RadiologyCases,
			&appointment.Prescription,
			&appointment.FuturePrescription,
			&appointment.LaboratoryMsg,
			&appointment.PhotographyMsg,
			&appointment.RadiologyMsg,
			&appointment.OrganizationID,
			&appointment.LaboratoryID,
			&appointment.PhotographyID,
			&appointment.RadiologyID,
			&appointment.LAdmissionAt,
			&appointment.RAdmissionAt,
			&appointment.PAdmissionAt,
			&appointment.LResultAt,
			&appointment.RResultAt,
			&appointment.PResultAt,
			&appointment.LRndImg,
			&appointment.RRndImg,
			&appointment.PRndImg,
			&appointment.LImgs,
			&appointment.RImgs,
			&appointment.PImgs,
			&appointment.Code,
			&appointment.IsVip,
			&appointment.VipIntroducer,
			&appointment.Absence,
			&appointment.FileID,
		)
		if err != nil {
			log.Println(err.Error(), "err")
			c.JSON(http.StatusOK, appointments)
			return
		}
		appointment.Last = GetAppointmentLastPrescription(fmt.Sprintf("%d", appointment.ID))
		appointment.PhotographyImages = GetResultImages(appointment.ID, "photo", c.Request.Host)
		appointment.RadiologyImages = GetResultImages(appointment.ID, "radio", c.Request.Host)
		if len(appointment.Logos) > 0 {
			appointments = append(appointments, appointment)
		}
		appointments = append(appointments, appointment)
	}
	c.JSON(http.StatusOK, appointments)
}
func GetAppointmentLastPrescription(id string) appointment2.LastPrescription {
	last := appointment2.LastPrescription{}
	query := "SELECT `prescription`, `future_prescription` FROM `last_prescription` WHERE appointment_id = ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error(), "err")
		return last
	}
	result := stmt.QueryRow(id)
	if err = result.Err(); err != nil {
		log.Println(err.Error(), "err")
		return last
	}
	err = result.Scan(
		last.Prescription,
		last.FuturePrescription,
	)
	if err != nil {
		log.Println(err.Error(), "err")
		return last
	}
	log.Println(last)
	return last
}

func (uc *UserControllerStruct) GetUserAppointmentResultList(c *gin.Context) {
	query := "SELECT id, user_id, created_at, ifnull(info, '') info, staff_id, start_at, end_at, status, ifnull(director_id, -1) director_id, updated_at, income, ifnull(subject, '') subject, ifnull(case_type, '') case_type, ifnull(laboratory_cases, '') laboratory_cases, ifnull(photography_cases, '') photography_cases, ifnull(radiology_cases, '') radiology_cases, ifnull(prescription, '') prescription, ifnull(future_prescription, '') future_prescription, ifnull(laboratory_msg, '') laboratory_msg, ifnull(photography_msg, '') photography_msg, ifnull(radiology_msg, '') radiology_msg, organization_id, ifnull(director_id, -1) laboratory_id, ifnull(photography_id, -1) photography_id, ifnull(radiology_id, -1) radiology_id, l_admission_at, r_admission_at, p_admission_at, l_result_at, r_result_at, p_result_at, ifnull(l_rnd_img, '') l_rnd_img, ifnull(r_rnd_img, '') r_rnd_img, ifnull(p_rnd_img, '') p_rnd_img, l_imgs, r_imgs, p_imgs, ifnull(code, '') code, is_vip, ifnull(vip_introducer, 0) vip_introducer, absence, ifnull(file_id, '') file_id FROM appointment WHERE user_id = ? "
	userID := c.Param("id")
	appointments := []appointment2.UserAppointmentInfo{}
	var appointment appointment2.UserAppointmentInfo
	staffUser := auth.GetStaffUser(c)
	staffOrg := organizationController.GetOrganization(fmt.Sprintf("%d", staffUser.OrganizationID))
	if staffOrg == nil {
		errorsHandler.GinErrorResponseHandler(c, nil)
		return
	}
	if staffOrg.ProfessionID == "1" {
		query += " AND photography_id = ? "
	} else if staffOrg.ProfessionID == "2" {
		query += " AND laboratory_id = ? "
	} else if staffOrg.ProfessionID == "3" {
		query += " AND radiology_id = ? "
	} else {
		query += " AND organization_id = ? "
	}
	prof := c.Query("prof")
	if prof == "photo" {
		query += " AND p_result_at IS NOT NULL "
	}
	if prof == "lab" {
		query += " AND l_result_at IS NOT NULL "
	}
	if prof == "radio" {
		query += " AND r_result_at IS NOT NULL "
	}
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error(), "prepare")
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	rows, err := stmt.Query(userID, staffUser.OrganizationID)
	defer rows.Close()
	if err != nil {
		log.Println(err.Error())
		return
	}
	for rows.Next() {
		err = rows.Scan(
			&appointment.ID,
			&appointment.UserID,
			&appointment.CreatedAt,
			&appointment.Info,
			&appointment.StaffID,
			&appointment.StartAt,
			&appointment.EndAt,
			&appointment.Status,
			&appointment.DirectorID,
			&appointment.UpdatedAt,
			&appointment.Income,
			&appointment.Subject,
			&appointment.CaseType,
			&appointment.LaboratoryCases,
			&appointment.PhotographyCases,
			&appointment.RadiologyCases,
			&appointment.Prescription,
			&appointment.FuturePrescription,
			&appointment.LaboratoryMsg,
			&appointment.PhotographyMsg,
			&appointment.RadiologyMsg,
			&appointment.OrganizationID,
			&appointment.LaboratoryID,
			&appointment.PhotographyID,
			&appointment.RadiologyID,
			&appointment.LAdmissionAt,
			&appointment.RAdmissionAt,
			&appointment.PAdmissionAt,
			&appointment.LResultAt,
			&appointment.RResultAt,
			&appointment.PResultAt,
			&appointment.LRndImg,
			&appointment.RRndImg,
			&appointment.PRndImg,
			&appointment.LImgs,
			&appointment.RImgs,
			&appointment.PImgs,
			&appointment.Code,
			&appointment.IsVip,
			&appointment.VipIntroducer,
			&appointment.Absence,
			&appointment.FileID,
		)
		if err != nil {
			log.Println(err.Error(), "err")
			c.JSON(http.StatusOK, appointments)
			return
		}
		appointment.Logos = GetResultImages(appointment.ID, prof, c.Request.Host)
		if len(appointment.Logos) > 0 {
			appointments = append(appointments, appointment)
		}
	}
	c.JSON(http.StatusOK, appointments)
}

func GetResultImages(id int64, prof string, host string) []string {
	logos := []string{}
	location := fmt.Sprintf("./images/results/%d", id)
	if _, err := os.Stat(location); os.IsNotExist(err) {
		return logos
	}
	switch prof {
	case "radio":
		location += "/radio"
		break
	case "photo":
		location += "/photo"
		break
	case "lab":
		location += "/lab"
		break
	}
	files, err := ioutil.ReadDir(location)
	if err != nil {
		log.Println(err.Error(), "system")
		return logos
	}
	for _, f := range files {
		if f.Name() != "." || f.Name() != ".." {
			logos = append(logos, fmt.Sprintf("http://%s/images/results/%d/%s/%s", host, id, prof, f.Name()))
		}
	}
	return logos
}

func (uc *UserControllerStruct) GetOrthodonticsList(c *gin.Context) {
	userID := c.Param("id")
	histories, _ := medicalHistory.GetMedicalHistory(userID)
	c.JSON(http.StatusOK, histories)
}

func (uc *UserControllerStruct) CreateOrthodonticsList(c *gin.Context) {
	var request medicalHistory.CreateMedicalHistoryStruct
	if errors := c.ShouldBindJSON(&request); errors != nil {
		log.Println(errors.Error())
		errorsHandler.GinErrorResponseHandler(c, errors)
		return
	}
	res := medicalHistory.CreateMedicalHistory(request)
	c.JSON(http.StatusOK, res)
}

func (uc *UserControllerStruct) Update(c *gin.Context) {
	var updateUserQuery = "UPDATE `user` SET"
	var values []interface{}
	var columns []string
	userId := c.Param("id")
	if userId == "" {
		log.Println(userId, "userid")
		errorsHandler.GinErrorResponseHandler(c, nil)
		return
	}
	var updateUserRequest organization.UpdateUserRequest
	if errors := c.ShouldBindJSON(&updateUserRequest); errors != nil {
		log.Println(errors.Error())
		errorsHandler.GinErrorResponseHandler(c, errors)
		return
	}
	getUserUpdateColumns(&updateUserRequest, &columns, &values)
	columnsString := strings.Join(columns, ",")
	if len(columnsString) == 0 {
		c.JSON(422, false)
		return
	}
	updateUserQuery += columnsString
	log.Println(updateUserQuery, "update")
	updateUserQuery += " WHERE `id` = ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(updateUserQuery)
	if err != nil {
		log.Println(err.Error())
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

func getUserUpdateColumns(o *organization.UpdateUserRequest, columns *[]string, values *[]interface{}) {
	if o.FirstName != "" {
		*columns = append(*columns, " `fname` = ? ")
		*values = append(*values, o.FirstName)
	}
	if o.LastName != "" {
		*columns = append(*columns, " `lname` = ? ")
		*values = append(*values, o.LastName)
	}
	if o.Info != "" {
		*columns = append(*columns, " `info` = ? ")
		*values = append(*values, o.Info)
	}
	if o.Description != "" {
		*columns = append(*columns, " `description` = ? ")
		*values = append(*values, o.Description)
	}
	if o.Relation != "" {
		*columns = append(*columns, " `relation` = ? ")
		*values = append(*values, o.Relation)
	}
	if o.FileID != "" {
		*columns = append(*columns, " `file_id` = ? ")
		*values = append(*values, o.FileID)
	}
	if o.Email != "" {
		*columns = append(*columns, " `email` = ? ")
		*values = append(*values, o.Email)
	}
	if o.Gender != "" && (strings.ToUpper(o.Gender) == "MALE" || strings.ToUpper(o.Gender) == "FEMALE") {
		*columns = append(*columns, " `gender` = ? ")
		*values = append(*values, o.Gender)
	}
	if o.Tel != "" {
		*columns = append(*columns, " `tel` = ? ")
		*values = append(*values, o.Tel)
	}
	if o.Logo != "" {
		*columns = append(*columns, " `logo` = ? ")
		*values = append(*values, o.Logo)
	}
	if o.Tel1 != "" {
		*columns = append(*columns, " `tel1` = ? ")
		*values = append(*values, o.Tel1)
	}
	if o.Nid != "" {
		*columns = append(*columns, " `nid` = ? ")
		*values = append(*values, o.Nid)
	}
	if o.BirthDate != "" {
		*columns = append(*columns, " `birth_date` = ? ")
		*values = append(*values, o.BirthDate)
	}
	if o.Address != "" {
		*columns = append(*columns, " `address` = ? ")
		*values = append(*values, o.Address)
	}
	if o.Password != "" {
		password, err := auth.PasswordHash(o.Password)
		if err != nil {
			log.Println(err.Error())
			return
		}
		*columns = append(*columns, " `pass` = ? ")
		*values = append(*values, password)
	}
}

func (uc *UserControllerStruct) Delete(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		return
	}
	//staff := auth.GetStaffUser(c)
	//if staff.UserGroupID != 2 {
	//	c.JSON(403, gin.H{})
	//	return
	//}
	stmt, err := repository.DBS.MysqlDb.Prepare(mysqlQuery.DeleteUserQuery)
	if err != nil {
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	stmt.QueryRow(userID)
	c.JSON(200, nil)
}

func (uc *UserControllerStruct) DeleteItems(c *gin.Context) {
	//staff := auth.GetStaffUser(c)
	//if staff.UserGroupID != 2 {
	//	c.JSON(403, gin.H{})
	//	return
	//}
	query := "DELETE FROM `user` WHERE id in ("
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
	log.Println(query, "q")
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

func (uc *UserControllerStruct) CreateCode(c *gin.Context) {
	var updateUserQuery = "UPDATE `user` SET `appcode` = ? WHERE `id` = ?"
	id := c.Param("id")
	if id == "" {
		errorsHandler.GinErrorResponseHandler(c, nil)
		return
	}
	stmt, err := repository.DBS.MysqlDb.Prepare(updateUserQuery)
	if err != nil {
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	code := helper.RandomString(6)
	_, err = stmt.Exec(code, id)
	if err != nil {
		log.Println(err.Error())
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	c.JSON(200, code)
}

func (uc *UserControllerStruct) SendDoc(c *gin.Context) {
	staff := auth.GetStaffUser(c)
	request := doc.CreateDocStruct{}
	if errors := c.ShouldBindJSON(&request); errors != nil {
		log.Println(errors.Error())
		errorsHandler.GinErrorResponseHandler(c, errors)
		return
	}
	oid, _ := strconv.ParseInt(fmt.Sprintf("%d", staff.OrganizationID), 10, 32)
	userID := c.Param("id")
	uid, _ := strconv.ParseInt(userID, 10, 32)
	request.OrganizationID = oid
	request.UserID = uid
	err := doc.CreateUserDoc(request)
	if err != nil {
		log.Println(err.Error(), "err")
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	c.JSON(http.StatusOK, nil)
}

func (uc *UserControllerStruct) DeleteDoc(c *gin.Context) {
	id := c.Param("id")
	err := doc.DeleteDoc(id)
	if err != nil {
		log.Println(err.Error(), "err")
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	c.JSON(http.StatusOK, nil)
}

func (uc *UserControllerStruct) UpdateDoc(c *gin.Context) {
	id := c.Param("id")
	request := doc.UpdateDocStruct{}
	if errors := c.ShouldBindJSON(&request); errors != nil {
		log.Println(errors.Error())
		errorsHandler.GinErrorResponseHandler(c, errors)
		return
	}
	err := doc.UpdateDoc(id, request.DoctorDesc, request.UserDesc)
	if err != nil {
		log.Println(err.Error(), "err")
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	doc, _ := doc.GetDoc(id)
	c.JSON(http.StatusOK, doc)
}

func (uc *UserControllerStruct) GetUserDocs(c *gin.Context) {
	staff := auth.GetStaffUser(c)
	userID := c.Param("id")
	docs, _ := doc.GetUserDocs(userID, fmt.Sprintf("%d", staff.OrganizationID))
	c.JSON(http.StatusOK, docs)
}

func (uc *UserControllerStruct) ChangePassword(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		return
	}
	var request organization.ChangeUserPasswordRequest
	if errors := c.ShouldBindJSON(&request); errors != nil {
		log.Println(errors.Error())
		errorsHandler.GinErrorResponseHandler(c, errors)
		return
	}
	stmt, err := repository.DBS.MysqlDb.Prepare(mysqlQuery.ChangePasswordQuery)
	if err != nil {
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	password, err := auth.PasswordHash(request.Password)
	if err != nil {
		log.Println(err.Error())
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	request.Password = password
	stmt.QueryRow(password, userID)
	c.JSON(200, nil)
}

func (uc *UserControllerStruct) GetUserWallet(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(500, gin.H{
			"message": "آی دی صحیح نیست",
		})
		return
	}
	uID, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	wallet := wallet2.GetWallet(uID, "user")
	c.JSON(200, wallet)
}

func (uc *UserControllerStruct) GetUserWalletHistories(c *gin.Context) {
	userID := c.Param("id")
	start_date := c.Query("start_date")
	end_date := c.Query("end_date")
	page := c.Query("page")
	q := c.Query("q")
	if userID == "" {
		c.JSON(500, gin.H{
			"message": "آی دی صحیح نیست",
		})
		return
	}
	uID, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	histories, _ := wallet.GetUserAllWalletHistories(uID, start_date, end_date, q, page)
	c.JSON(200, histories)
}

func (uc *UserControllerStruct) GetUserWalletHistoriesSum(c *gin.Context) {
	userID := c.Param("id")
	start_date := c.Query("start_date")
	end_date := c.Query("end_date")
	if userID == "" {
		c.JSON(500, gin.H{
			"message": "آی دی صحیح نیست",
		})
		return
	}
	var sum int64
	sum = 0
	if start_date == "" || end_date == "" || start_date == "null" || end_date == "null" {
		c.JSON(200, gin.H{
			"date": sum,
		})
	}
	uID, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	sum, _ = wallet.GetUserWalletHistoriesSum(uID, start_date, end_date)
	c.JSON(200, sum)
}

func (uc *UserControllerStruct) GetUserWalletHistoriesDays(c *gin.Context) {
	userID := c.Param("id")
	start_date := c.Query("start_date")
	end_date := c.Query("end_date")
	if userID == "" {
		c.JSON(500, gin.H{
			"message": "آی دی صحیح نیست",
		})
		return
	}
	histories := []wallet.WalletHistoryStruct{}
	if start_date == "" || end_date == "" || start_date == "null" || end_date == "null" {
		c.JSON(200, gin.H{
			"date": histories,
		})
	}
	uID, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	histories, _ = wallet.GetUserWalletHistoriesDays(uID, start_date, end_date)
	c.JSON(200, histories)
}

func (uc *UserControllerStruct) IncreaseUserWallet(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		return
	}
	var request wallet2.ChangeUserWalletBalance
	if errors := c.ShouldBindJSON(&request); errors != nil {
		log.Println(errors.Error())
		errorsHandler.GinErrorResponseHandler(c, errors)
		return
	}
	uID, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		c.JSON(500, nil)
		return
	}
	wallet := wallet2.GetWallet(uID, "user")
	result, balance := wallet.Increase(request.Amount)
	if result {
		c.JSON(200, balance)
		return
	}
	c.JSON(500, nil)
}

func (uc *UserControllerStruct) DecreaseUserWallet(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		return
	}
	var request wallet2.ChangeUserWalletBalance
	if errors := c.ShouldBindJSON(&request); errors != nil {
		log.Println(errors.Error())
		errorsHandler.GinErrorResponseHandler(c, errors)
		return
	}
	uID, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		c.JSON(500, nil)
		return
	}
	wallet := wallet2.GetWallet(uID, "user")
	result, balance := wallet.Decrease(request.Amount, false)
	if result {
		c.JSON(200, balance)
		return
	}
	c.JSON(500, nil)
}

func (uc *UserControllerStruct) SetUserWallet(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		return
	}
	var request wallet2.ChangeUserWalletBalance
	if errors := c.ShouldBindJSON(&request); errors != nil {
		log.Println(errors.Error())
		errorsHandler.GinErrorResponseHandler(c, errors)
		return
	}
	uID, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		c.JSON(500, nil)
		return
	}
	wallet := wallet2.GetWallet(uID, "user")
	result := wallet.SetBalance(request.Amount)
	if result {
		c.JSON(200, nil)
		return
	}
	c.JSON(500, nil)
}
