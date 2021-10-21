package caseTypeController

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"gitlab.com/simateb-project/simateb-backend/domain/caseType"
	"gitlab.com/simateb-project/simateb-backend/repository"
	appointment2 "gitlab.com/simateb-project/simateb-backend/repository/appointment"
	mysqlQuery "gitlab.com/simateb-project/simateb-backend/repository/mysqlQuery/auth"
	"gitlab.com/simateb-project/simateb-backend/utils/auth"
	"gitlab.com/simateb-project/simateb-backend/utils/errorsHandler"
	"log"
	"net/http"
	"strings"
)

type CaseTypeControllerInterface interface {
	Create(c *gin.Context)
	GetListByOrganization(c *gin.Context)
	Get(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type CaseTypeControllerStruct struct {
	r *appointment2.AppointmentRepositoryStruct
}

func NewCaseTypeController(r *appointment2.AppointmentRepositoryStruct) CaseTypeControllerInterface {
	x := &CaseTypeControllerStruct{
		r: r,
	}
	return x
}

func (uc *CaseTypeControllerStruct) Create(c *gin.Context) {
	var createCaseTypeRequest caseType.CreateCaseTypeRequest
	if errors := c.ShouldBindJSON(&createCaseTypeRequest); errors != nil {
		log.Println(errors.Error())
		errorsHandler.GinErrorResponseHandler(c, errors)
		return
	}
	stmt, err := repository.DBS.MysqlDb.Prepare(mysqlQuery.CreateCaseTypeQuery)
	if err != nil {
		log.Println(err.Error())
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	defer stmt.Close()
	result, err := stmt.Exec(
		createCaseTypeRequest.Name,
		createCaseTypeRequest.OrganizationId,
		createCaseTypeRequest.Duration,
		createCaseTypeRequest.IsLimited,
		createCaseTypeRequest.Limitation,
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

func (uc *CaseTypeControllerStruct) Get(c *gin.Context) {
	caseId := c.Param("id")
	if caseId == "" {
		return
	}
	stmt, err := repository.DBS.MysqlDb.Prepare(mysqlQuery.GetCaseTypeQuery)
	if err != nil {
		log.Println(err.Error(), "prepare")
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	var caseType caseType.CaseTypeInfo
	result := stmt.QueryRow(caseId)
	err = result.Scan(
		&caseType.ID,
		&caseType.Name,
		&caseType.OrganizationId,
		&caseType.Duration,
		&caseType.IsLimited,
		&caseType.Limitation,
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
	c.JSON(http.StatusOK, caseType)
}

func (uc *CaseTypeControllerStruct) GetListByOrganization(c *gin.Context) {
	staffUser := auth.GetStaffUser(c)
	var caseTypes []caseType.CaseTypeInfo
	caseTypes = LoadCaseTypesByOrgId(staffUser.OrganizationID)
	c.JSON(http.StatusOK, caseTypes)
}

func LoadCaseTypesByOrgId(id int64) []caseType.CaseTypeInfo {
	var caseTypes []caseType.CaseTypeInfo
	var caseType caseType.CaseTypeInfo
	stmt, err := repository.DBS.MysqlDb.Prepare(mysqlQuery.GetCaseTypesListForOrganizationQuery)
	if err != nil {
		return caseTypes
	}
	rows, error := stmt.Query(id)
	if error != nil {
		return caseTypes
	}
	for rows.Next() {
		err := rows.Scan(
			&caseType.ID,
			&caseType.Name,
			&caseType.OrganizationId,
			&caseType.Duration,
			&caseType.IsLimited,
			&caseType.Limitation,
		)
		if err != nil {
			log.Println(err.Error())
			return caseTypes
		}
		caseTypes = append(caseTypes, caseType)
	}
	return caseTypes
}

func (uc *CaseTypeControllerStruct) Update(c *gin.Context) {
	var updateUserQuery = "UPDATE `case_type` SET"
	var values []interface{}
	var columns []string
	caseId := c.Param("id")
	if caseId == "" {
		log.Println(caseId, "case id")
		errorsHandler.GinErrorResponseHandler(c, nil)
		return
	}
	var updateCaseTypeRequest caseType.UpdateCaseTypeRequest
	if errors := c.ShouldBindJSON(&updateCaseTypeRequest); errors != nil {
		log.Println(errors.Error())
		errorsHandler.GinErrorResponseHandler(c, errors)
		return
	}
	getCaseTypeUpdateColumns(&updateCaseTypeRequest, &columns, &values)
	columnsString := strings.Join(columns, ",")
	updateUserQuery += columnsString
	updateUserQuery += " WHERE `id` = ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(updateUserQuery)
	if err != nil {
		log.Println(err.Error(), "prepare")
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	values = append(values, caseId)
	_, error := stmt.Exec(values...)
	if error != nil {
		log.Println(error.Error())
		errorsHandler.GinErrorResponseHandler(c, error)
		return
	}
	c.JSON(200, true)
}

func getCaseTypeUpdateColumns(o *caseType.UpdateCaseTypeRequest, columns *[]string, values *[]interface{}) {
	if o.Name != "" {
		*columns = append(*columns, " `name` = ? ")
		*values = append(*values, o.Name)
	}
	*columns = append(*columns, " `organization_id` = ? ")
	*values = append(*values, o.OrganizationId)
	*columns = append(*columns, " `limitation` = ? ")
	*values = append(*values, o.Limitation)
	*columns = append(*columns, " `is_limited` = ? ")
	*values = append(*values, o.IsLimited)
	*columns = append(*columns, " `duration` = ? ")
	*values = append(*values, o.Duration)
}

func (uc *CaseTypeControllerStruct) Delete(c *gin.Context) {
	caseID := c.Param("id")
	if caseID == "" {
		return
	}
	stmt, err := repository.DBS.MysqlDb.Prepare(mysqlQuery.DeleteCaseTypeQuery)
	if err != nil {
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	stmt.Query(caseID)
	c.JSON(200, nil)
}
