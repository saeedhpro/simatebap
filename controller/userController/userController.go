package userController

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"gitlab.com/simateb-project/simateb-backend/controller/organizationController"
	appointment2 "gitlab.com/simateb-project/simateb-backend/domain/appointment"
	"gitlab.com/simateb-project/simateb-backend/domain/organization"
	"gitlab.com/simateb-project/simateb-backend/repository"
	mysqlQuery "gitlab.com/simateb-project/simateb-backend/repository/mysqlQuery/auth"
	"gitlab.com/simateb-project/simateb-backend/repository/user"
	"gitlab.com/simateb-project/simateb-backend/utils/auth"
	"gitlab.com/simateb-project/simateb-backend/utils/errorsHandler"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type UserControllerInterface interface {
	Create(c *gin.Context)
	Get(c *gin.Context)
	GetList(c *gin.Context)
	GetListForAdmin(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	ChangePassword(c *gin.Context)
	GetUserAppointmentList(c *gin.Context)
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
	page := c.Query("page")
	if page == "" {
		page = "1"
	}
	offset, err := strconv.Atoi(page)
	offset = (offset - 1) * 10
	stmt, err := repository.DBS.MysqlDb.Prepare(mysqlQuery.GetAdminUsersQuery)
	if err != nil {
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
		rows.Scan(
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
		users = append(users, user)
	}
	c.JSON(http.StatusOK, users)
}

func (uc *UserControllerStruct) GetList(c *gin.Context) {
	getUsersQuery := "SELECT user.id, user.fname, ifnull(user.lname, ''), ifnull(user.last_login, ''), ifnull(user.created, ''),  ifnull(user.tel, ''), user.user_group_id user_group_id, user_group.name user_group_name, ifnull(birth_date, ''), ifnull(relation, ''), ifnull(description, '') FROM user LEFT JOIN user_group on user.user_group_id = user_group.id WHERE organization_id = ? AND user_group_id != 1 "
	page := c.Query("page")
	if page == "" {
		page = "1"
	}
	q := c.Query("q")
	if q != "" {
		getUsersQuery += " LIKE %" + q + "% "
	}
	offset, err := strconv.Atoi(page)
	offset = (offset - 1) * 10
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
		log.Println("get user first")
		rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.LastLogin,
			&user.Created,
			&user.Tel,
			&user.UserGroupID,
			&user.UserGroupName,
			&user.BirthDate,
			&user.Relation,
			&user.Description,
		)
		log.Println("get user")
		organizationID := strconv.FormatInt(staffUser.OrganizationID, 10)
		user.OrganizationID = organizationID
		user.OrganizationName = organizationController.GetOrganization(organizationID).Name
		users = append(users, user)
	}
	c.JSON(http.StatusOK, users)
}

func (uc *UserControllerStruct) GetUserAppointmentList(c *gin.Context) {
	getUserAppointmentListQuery := "SELECT * FROM appointment WHERE organization_id = ? AND user_id = ?"
	userID := c.Param("id")
	stmt, err := repository.DBS.MysqlDb.Prepare(getUserAppointmentListQuery)
	if err != nil {
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	staffUser := auth.GetStaffUser(c)
	if staffUser == nil {
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	var appointments []appointment2.UserAppointmentInfo
	var appointment appointment2.UserAppointmentInfo
	rows, err := stmt.Query(staffUser.OrganizationID, userID)
	defer rows.Close()
	if err != nil {
		log.Println(err.Error())
		return
	}
	for rows.Next() {
		rows.Scan(
			&appointment.ID,
		)
	}
	c.JSON(http.StatusOK, appointments)
}

func (uc *UserControllerStruct) Update(c *gin.Context) {
	var updateUserQuery = "UPDATE `user` SET"
	var values []interface{}
	var columns []string
	userId := c.Param("id")
	if userId == "" {
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
		*values = append(*values, o.Info)
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
	stmt, err := repository.DBS.MysqlDb.Prepare(mysqlQuery.DeleteUserQuery)
	if err != nil {
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	stmt.QueryRow(userID)
	c.JSON(200, nil)
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
