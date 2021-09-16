package organizationController

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/simateb-project/simateb-backend/domain/organization"
	"gitlab.com/simateb-project/simateb-backend/helper"
	"gitlab.com/simateb-project/simateb-backend/repository"
	mysqlQuery "gitlab.com/simateb-project/simateb-backend/repository/mysqlQuery/auth"
	"gitlab.com/simateb-project/simateb-backend/utils/auth"
	"gitlab.com/simateb-project/simateb-backend/utils/errorsHandler"
	"gitlab.com/simateb-project/simateb-backend/utils/pagination"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type OrganizationControllerInterface interface {
	Create(c *gin.Context)
	Get(c *gin.Context)
	Update(c *gin.Context)
	GetList(c *gin.Context)
	GetListForAdmin(c *gin.Context)
	GetUsers(c *gin.Context)
}

type OrganizationControllerStruct struct {
}

func NewOrganizationController() OrganizationControllerInterface {
	x := &OrganizationControllerStruct{
	}
	return x
}

func (oc *OrganizationControllerStruct) Create(c *gin.Context) {
	var createOrganizationRequest organization.CreateOrganizationRequest
	if errors := c.ShouldBindJSON(&createOrganizationRequest); errors != nil {
		log.Println(errors.Error())
		errorsHandler.GinErrorResponseHandler(c, errors)
		return
	}
	stmt, err := repository.DBS.MysqlDb.Prepare(mysqlQuery.CreateOrganizationQuery)
	if err != nil {
		log.Println(err.Error())
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	defer stmt.Close()
	staffID := auth.GetStaffUser(c).UserID
	result, err := stmt.Exec(
		createOrganizationRequest.Name,
		createOrganizationRequest.KnownAs,
		createOrganizationRequest.ProfessionID,
		createOrganizationRequest.Logo,
		createOrganizationRequest.Phone,
		createOrganizationRequest.Phone1,
		staffID,
		createOrganizationRequest.Info,
		createOrganizationRequest.CaseTypes,
		createOrganizationRequest.SmsPrice,
		createOrganizationRequest.SmsCredit,
		createOrganizationRequest.Website,
		createOrganizationRequest.Instagram,
	)
	if err != nil {
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	setOrganizationRelations(id, createOrganizationRequest.RelRadiologies, createOrganizationRequest.RelLaboratories, createOrganizationRequest.RelDoctorOffices)
	c.JSON(http.StatusOK, true)
}

func setOrganizationRelations(id int64, radiologies []organization.RelOrganizationType, laboratories []organization.RelOrganizationType, offices []organization.RelOrganizationType) {
	var ids []int64
	stmt, err := repository.DBS.MysqlDb.Prepare(mysqlQuery.GetOrganizationRelations)
	if err != nil {
		return
	}
	rows, error := stmt.Query(id)
	if error != nil {
		return
	}
	for rows.Next() {
		var rel_organization_id int64
		rows.Scan(&rel_organization_id)
		ids = append(ids, rel_organization_id)
	}
	insertQuery := "INSERT INTO `rel_organization`(`organization_id`, `profession_id`, `rel_organization_id`) VALUES "
	var radValues []organization.RelOrganizationType
	var labValues []organization.RelOrganizationType
	var offValues []organization.RelOrganizationType
	var queryStr []string
	var allValues []interface{}
	for _, n := range radiologies {
		if exists := helper.ItemExists(ids, n); !exists {
			radValues = append(radValues, n)
		}
	}
	for _, n := range laboratories {
		if !helper.ItemExists(ids, n) {
			labValues = append(labValues, n)
		}
	}
	for _, n := range offices {
		if !helper.ItemExists(ids, n) {
			offValues = append(offValues, n)
		}
	}
	for _, i := range radValues {
		queryStr = append(queryStr, "(?,?,?)")
		allValues = append(allValues, id, i.ProfessionID, i.ID)
	}
	for _, i := range labValues {
		queryStr = append(queryStr, "(?,?,?)")
		allValues = append(allValues, id, i.ProfessionID, i.ID)
	}
	for _, i := range offValues {
		queryStr = append(queryStr, "(?,?,?)")
		allValues = append(allValues, id, i.ProfessionID, i.ID)
	}
	insertQuery = fmt.Sprintf("%s%s", insertQuery, strings.Join(queryStr, ","))
	stmt, error = repository.DBS.MysqlDb.Prepare(insertQuery)
	if error != nil {
		log.Println(error.Error())
		return
	}
	_, error = stmt.Exec(allValues...)
	if error != nil {
		log.Println(error.Error())
	}
}

func (oc *OrganizationControllerStruct) Get(c *gin.Context) {
	organizationID := c.Param("id")
	if organizationID == "" {
		return
	}
	stmt, err := repository.DBS.MysqlDb.Prepare(mysqlQuery.GetOrganizationQuery)
	if err != nil {
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	var organizationInfo organization.OrganizationInfo
	result := stmt.QueryRow(organizationID)
	err = result.Scan(
		&organizationInfo.ID,
		&organizationInfo.Name,
		&organizationInfo.Phone,
		&organizationInfo.Phone1,
		&organizationInfo.ProfessionID,
		&organizationInfo.KnownAs,
		&organizationInfo.CaseTypes,
		&organizationInfo.StaffID,
		&organizationInfo.Info,
		&organizationInfo.Website,
		&organizationInfo.Instagram,
		&organizationInfo.SmsPrice,
		&organizationInfo.SmsCredit,
		&organizationInfo.CreatedAt,
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
	c.JSON(http.StatusOK, organizationInfo)
}

func (oc *OrganizationControllerStruct) GetList(c *gin.Context) {
	var query = "SELECT id, ifnull(name, ''), ifnull(phone, ''), ifnull(phone1, ''), ifnull(profession_id, '')," +
		" ifnull(known_as, ''), ifnull(case_types, ''), ifnull(staff_id, ''), ifnull(info, ''), ifnull(website, '')," +
		" ifnull(instagram, ''), sms_price, sms_credit FROM organization "
	var values []interface{}
	userGroupID := c.Query("group")
	q := c.Query("q")
	var query2 = ""
	if q != "" {
		q = "'%" + q + "%'"
		query2 += fmt.Sprintf(" WHERE name LIKE %s ", q)
	}
	var err error
	page := c.Query("page")
	if page == "" {
		page = "1"
	}
	offset, err := strconv.Atoi(page)
	offset = (offset - 1) * 10
	values = append(values, offset)
	if userGroupID != "" {
		if userGroupID == "2" || userGroupID == "3" {
			if query2 != "" {
				query += query2
				query += fmt.Sprintf(" AND profession_id = ? ")
			} else {
				query += fmt.Sprintf(" WHERE profession_id = ? ")
			}
			values = append(values,userGroupID)
		} else {
			query += fmt.Sprintf(" AND profession_id != 2 AND profession_id != 3 ")
		}
	}
	query += "LIMIT 10 OFFSET ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error(), "log")
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	var organizations []organization.OrganizationInfo
	var organizationInfo organization.OrganizationInfo
	rows, error := stmt.Query(values...)
	if error != nil {
		log.Println(error.Error(), "error")
		errorsHandler.GinErrorResponseHandler(c, error)
		return
	}
	for rows.Next() {
		err := rows.Scan(
			&organizationInfo.ID,
			&organizationInfo.Name,
			&organizationInfo.Phone,
			&organizationInfo.Phone1,
			&organizationInfo.ProfessionID,
			&organizationInfo.KnownAs,
			&organizationInfo.CaseTypes,
			&organizationInfo.StaffID,
			&organizationInfo.Info,
			&organizationInfo.Website,
			&organizationInfo.Instagram,
			&organizationInfo.SmsPrice,
			&organizationInfo.SmsCredit,
			&organizationInfo.CreatedAt,
		)
		if err != nil {
			log.Println(err.Error())
			return
		}
		profession := getProfession(organizationInfo.ProfessionID)
		if profession != nil {
			organizationInfo.Profession = *profession
		}
		staff := getStaff(organizationInfo.StaffID)
		if staff != nil {
			organizationInfo.Staff = *staff
		}
		organizations = append(organizations, organizationInfo)
	}
	i, err := strconv.Atoi(page)
	paginated := pagination.OrganizationPaginationInfo{
		Data:        organizations,
		HasNextPage: true,
		PrevPage:    -1,
		NextPage:    2,
		Page:        i,
	}
	c.JSON(http.StatusOK, paginated)
}
func (oc *OrganizationControllerStruct) GetListForAdmin(c *gin.Context) {
	var query = "SELECT id, ifnull(name, ''), ifnull(phone, ''), ifnull(phone1, ''), ifnull(profession_id, '')," +
		" ifnull(known_as, ''), ifnull(case_types, ''), ifnull(staff_id, ''), ifnull(info, ''), ifnull(website, '')," +
		" ifnull(instagram, ''), sms_price, sms_credit FROM organization "
	var values []interface{}
	q := c.Query("q")
	if q != "" {
		q = "'%" + q + "%'"
		query += fmt.Sprintf(" WHERE name LIKE %s ", q)
	}
	var err error
	page := c.Query("page")
	if page == "" {
		page = "1"
	}
	offset, err := strconv.Atoi(page)
	offset = (offset - 1) * 10
	values = append(values, offset)
	query += "LIMIT 10 OFFSET ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error(), "log")
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	var organizations []organization.OrganizationInfo
	var organizationInfo organization.OrganizationInfo
	rows, error := stmt.Query(values...)
	if error != nil {
		log.Println(error.Error(), "error")
		errorsHandler.GinErrorResponseHandler(c, error)
		return
	}
	for rows.Next() {
		err := rows.Scan(
			&organizationInfo.ID,
			&organizationInfo.Name,
			&organizationInfo.Phone,
			&organizationInfo.Phone1,
			&organizationInfo.ProfessionID,
			&organizationInfo.KnownAs,
			&organizationInfo.CaseTypes,
			&organizationInfo.StaffID,
			&organizationInfo.Info,
			&organizationInfo.Website,
			&organizationInfo.Instagram,
			&organizationInfo.SmsPrice,
			&organizationInfo.SmsCredit,
			&organizationInfo.CreatedAt,
		)
		if err != nil {
			log.Println(err.Error())
			return
		}
		profession := getProfession(organizationInfo.ProfessionID)
		if profession != nil {
			organizationInfo.Profession = *profession
		}
		staff := getStaff(organizationInfo.StaffID)
		if staff != nil {
			organizationInfo.Staff = *staff
		}
		organizations = append(organizations, organizationInfo)
	}
	i, err := strconv.Atoi(page)
	paginated := pagination.OrganizationPaginationInfo{
		Data:        organizations,
		HasNextPage: true,
		PrevPage:    -1,
		NextPage:    2,
		Page:        i,
	}
	c.JSON(http.StatusOK, paginated)
}

func getProfession(id string) *organization.SimpleProfessionInfo {
	stmt, err := repository.DBS.MysqlDb.Prepare(mysqlQuery.GetSimpleProfessionQuery)
	var professionInfo organization.SimpleProfessionInfo
	if err != nil {
		return nil
	}
	result := stmt.QueryRow(id)
	err = result.Scan(
		&professionInfo.ID,
		&professionInfo.Name,
	)
	if err != nil {
		return nil
	}
	return &professionInfo
}

func getStaff(id int64) *organization.SimpleUserInfo {
	stmt, err := repository.DBS.MysqlDb.Prepare(mysqlQuery.GetSimpleStaffQuery)
	var userInfo organization.SimpleUserInfo
	if err != nil {
		return nil
	}
	result := stmt.QueryRow(id)
	err = result.Scan(
		&userInfo.ID,
		&userInfo.FirstName,
		&userInfo.LastName,
		&userInfo.Organization,
	)
	if err != nil {
		return nil
	}
	return &userInfo
}

func (oc *OrganizationControllerStruct) GetUsers(c *gin.Context) {
	organizationID := c.Param("id")
	userGroupID := c.Query("group")
	if userGroupID == "" {
		c.JSON(422, struct {
			Message string `json:"message"`
			Code    int    `json:"code"`
		}{
			Message: "فیلد group الزامی است",
			Code:    422,
		})
		return
	}
	page := c.Query("page")
	if page == "" {
		page = "1"
	}
	offset, err := strconv.Atoi(page)
	offset = (offset - 1) * 10
	stmt, err := repository.DBS.MysqlDb.Prepare(mysqlQuery.GetOrganizationUsersQuery)
	if err != nil {
		log.Println(err.Error())
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	var users []organization.OrganizationUser
	var user organization.OrganizationUser
	rows, err := stmt.Query(organizationID, userGroupID, offset)
	defer rows.Close()
	if err != nil {
		log.Println(err.Error())
		return
	}
	for rows.Next() {
		err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Tel,
			&user.UserGroupID,
			&user.Created,
			&user.LastLogin,
			&user.BirthDate,
			&user.UserGroupName,
			&user.OrganizationName,
			&user.OrganizationID,
		)
		if err != nil {
			log.Println(err.Error(), "user log")
			return
		}
		users = append(users, user)
	}
	c.JSON(http.StatusOK, users)
}

func (oc *OrganizationControllerStruct) Update(c *gin.Context) {
	var updateOrganizationQuery = "UPDATE `organization` SET"
	var values []interface{}
	var columns []string
	organizationId := c.Param("id")
	if organizationId == "" {
		errorsHandler.GinErrorResponseHandler(c, nil)
		return
	}
	var updateOrganizationRequest organization.UpdateOrganizationRequest
	if errors := c.ShouldBindJSON(&updateOrganizationRequest); errors != nil {
		log.Println(errors.Error())
		errorsHandler.GinErrorResponseHandler(c, errors)
		return
	}
	getOrganizationUpdateColumns(&updateOrganizationRequest, &columns, &values)
	columnsString := strings.Join(columns, ",")
	updateOrganizationQuery += columnsString
	updateOrganizationQuery += " WHERE `id` = ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(updateOrganizationQuery)
	if err != nil {
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	values = append(values, organizationId)
	_, error := stmt.Exec(values...)
	if error != nil {
		log.Println(error.Error())
		errorsHandler.GinErrorResponseHandler(c, error)
		return
	}
	id, err := strconv.ParseInt(organizationId, 10, 64)
	if err == nil {
		setOrganizationRelations(id, updateOrganizationRequest.RelRadiologies, updateOrganizationRequest.RelLaboratories, updateOrganizationRequest.RelDoctorOffices)
	}
	c.JSON(200, true)
}

func getOrganizationUpdateColumns(o *organization.UpdateOrganizationRequest, columns *[]string, values *[]interface{}) {
	if o.Name != "" {
		*columns = append(*columns, " `name` = ? ")
		*values = append(*values, o.Name)
	}
	if o.Phone != "" {
		*columns = append(*columns, " `phone` = ? ")
		*values = append(*values, o.Phone)
	}
	if o.Phone1 != "" {
		*columns = append(*columns, " `phone1` = ? ")
		*values = append(*values, o.Phone1)
	}
	if o.KnownAs != "" {
		*columns = append(*columns, " `known_as` = ? ")
		*values = append(*values, o.KnownAs)
	}
	if o.CaseTypes != "" {
		*columns = append(*columns, " `case_types` = ? ")
		*values = append(*values, o.CaseTypes)
	}
	if o.Info != "" {
		*columns = append(*columns, " `info` = ? ")
		*values = append(*values, o.Info)
	}
	if o.Website != "" {
		*columns = append(*columns, " `website` = ? ")
		*values = append(*values, o.Website)
	}
	if o.Instagram != "" {
		*columns = append(*columns, " `instagram` = ? ")
		*values = append(*values, o.Instagram)
	}
	*columns = append(*columns, " `sms_credit` = ? ")
	*values = append(*values, o.SmsCredit)
	*columns = append(*columns, " `sms_price` = ? ")
	*values = append(*values, o.SmsPrice)
	if o.Logo != "" {
		*columns = append(*columns, " `logo` = ? ")
		*values = append(*values, o.Logo)
	}
}

func GetOrganization(id string) *organization.OrganizationInfo {
	stmt, err := repository.DBS.MysqlDb.Prepare(mysqlQuery.GetSimpleOrganizationQuery)
	var organizationInfo organization.OrganizationInfo
	if err != nil {
		return nil
	}
	result := stmt.QueryRow(id)
	err = result.Scan(
		&organizationInfo.ID,
		&organizationInfo.Name,
	)
	if err != nil {
		return nil
	}
	return &organizationInfo
}

func GetGroup(id int64) *organization.UserGroup {
	stmt, err := repository.DBS.MysqlDb.Prepare(mysqlQuery.GetUserGroupQuery)
	var userGroup organization.UserGroup
	if err != nil {
		return nil
	}
	result := stmt.QueryRow(id)
	err = result.Scan(&userGroup.ID, &userGroup.Name)
	if err != nil {
		return nil
	}
	return &userGroup
}
