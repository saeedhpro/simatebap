package userOrganizationPrice

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"gitlab.com/simateb-project/simateb-backend/domain/userOrgPrice"
	"gitlab.com/simateb-project/simateb-backend/repository"
	"gitlab.com/simateb-project/simateb-backend/utils/errorsHandler"
	"log"
	"net/http"
	"strings"
)

type UserOrganizationPriceControllerInterface interface {
	Create(c *gin.Context)
	GetList(c *gin.Context)
	Get(c *gin.Context)
	Update(c *gin.Context)
}

type UserOrganizationPriceControllerStruct struct {
}

func NewUserOrganizationPriceController() UserOrganizationPriceControllerInterface {
	x := &UserOrganizationPriceControllerStruct{
	}
	return x
}

func (uc *UserOrganizationPriceControllerStruct) Create(c *gin.Context) {
	var request userOrgPrice.CreateUserOrgPriceRequest
	if errors := c.ShouldBindJSON(&request); errors != nil {
		log.Println(errors.Error())
		errorsHandler.GinErrorResponseHandler(c, errors)
		return
	}
	query := "INSERT INTO `user_org_prices`(`user_id`,`organization_id`,`commission`) VALUES (?,?,?)"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error())
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	defer stmt.Close()
	result, err := stmt.Exec(
		&request.UserID,
		&request.OrganizationID,
		&request.Commission,
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

func (uc *UserOrganizationPriceControllerStruct) Get(c *gin.Context) {
	caseId := c.Param("id")
	if caseId == "" {
		return
	}
	query := "SELECT `id`, `user_id`, `organization_id`, `commission` FROM `user_org_prices` WHERE id = ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error(), "prepare")
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	var userOrgPrice userOrgPrice.UserOrgPriceInfo
	result := stmt.QueryRow(caseId)
	err = result.Scan(
		&userOrgPrice.ID,
		&userOrgPrice.UserID,
		&userOrgPrice.OrganizationID,
		&userOrgPrice.Commission,
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
	c.JSON(http.StatusOK, userOrgPrice)
}

func (uc *UserOrganizationPriceControllerStruct) GetList(c *gin.Context) {
	casesList := []userOrgPrice.UserOrgPriceInfo{}
	page := c.Query("page")
	id := c.Param("id")
	if page == "" {
		page = "1"
	}
	casesList = LoadUserOrganizationPrices(id, page)
	c.JSON(http.StatusOK, casesList)
}

func LoadUserOrganizationPrices(id string, page string) []userOrgPrice.UserOrgPriceInfo {
	casesList := []userOrgPrice.UserOrgPriceInfo{}
	var cases userOrgPrice.UserOrgPriceInfo
	query := "SELECT id, user_id, organization_id, commission FROM user_org_prices WHERE organization_id = ? LIMIT 10 OFFSET ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		return casesList
	}
	rows, error := stmt.Query(id, page)
	if error != nil {
		return casesList
	}
	for rows.Next() {
		err := rows.Scan(
			&cases.ID,
			&cases.UserID,
			&cases.OrganizationID,
			&cases.Commission,
		)
		if err != nil {
			log.Println(err.Error())
			return casesList
		}
		casesList = append(casesList, cases)
	}
	return casesList
}

func (uc *UserOrganizationPriceControllerStruct) Update(c *gin.Context) {
	var updateUserQuery = "UPDATE `user_org_prices` SET"
	var values []interface{}
	var columns []string
	caseId := c.Param("id")
	if caseId == "" {
		log.Println(caseId, "case price id")
		errorsHandler.GinErrorResponseHandler(c, nil)
		return
	}
	var request userOrgPrice.UpdateUserOrgPriceRequest
	if errors := c.ShouldBindJSON(&request); errors != nil {
		log.Println(errors.Error())
		errorsHandler.GinErrorResponseHandler(c, errors)
		return
	}
	getUserOrganizationPriceUpdateColumns(&request, &columns, &values)
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

func getUserOrganizationPriceUpdateColumns(o *userOrgPrice.UpdateUserOrgPriceRequest, columns *[]string, values *[]interface{}) {
	*columns = append(*columns, " `user_id` = ? ")
	*values = append(*values, o.UserID)
	*columns = append(*columns, " `organization_id` = ? ")
	*values = append(*values, o.OrganizationID)
	*columns = append(*columns, " `commission` = ? ")
	*values = append(*values, o.Commission)
}
