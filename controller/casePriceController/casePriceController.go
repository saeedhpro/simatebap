package casePriceController

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"gitlab.com/simateb-project/simateb-backend/domain/casePrice"
	"gitlab.com/simateb-project/simateb-backend/repository"
	"gitlab.com/simateb-project/simateb-backend/utils/errorsHandler"
	"log"
	"net/http"
	"strings"
)

type CasePriceControllerInterface interface {
	Create(c *gin.Context)
	GetList(c *gin.Context)
	Get(c *gin.Context)
	Update(c *gin.Context)
}

type CasePriceControllerStruct struct {
}

func NewCasePriceController() CasePriceControllerInterface {
	x := &CasePriceControllerStruct{
	}
	return x
}

func (uc *CasePriceControllerStruct) Create(c *gin.Context) {
	var request casePrice.CreateCasePriceRequest
	if errors := c.ShouldBindJSON(&request); errors != nil {
		log.Println(errors.Error())
		errorsHandler.GinErrorResponseHandler(c, errors)
		return
	}
	query := "INSERT INTO `case_prices`(`case_id`,`organization_id`,`price`) VALUES (?,?,?)"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error())
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	defer stmt.Close()
	result, err := stmt.Exec(
		&request.CaseId,
		&request.OrganizationId,
		&request.Price,
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

func (uc *CasePriceControllerStruct) Get(c *gin.Context) {
	caseId := c.Param("id")
	if caseId == "" {
		return
	}
	//"SELECT `case_prices`.`id`, `case_prices`.`case_id`, `case_prices`.`organization_id`, `case_prices`.`price`, `cases`.`name` `case_name` FROM `case_prices` LEFT JOIN `cases` ON `case_prices`.`case_id` = `cases`.`id`"
	query := "SELECT `id`, `case_id`, `organization_id`, `price` FROM `case_prices` WHERE id = ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error(), "prepare")
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	var casePrice casePrice.CasePriceInfo
	result := stmt.QueryRow(caseId)
	err = result.Scan(
		&casePrice.ID,
		&casePrice.CaseId,
		&casePrice.OrganizationId,
		&casePrice.Price,
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
	c.JSON(http.StatusOK, casePrice)
}

func (uc *CasePriceControllerStruct) GetList(c *gin.Context) {
	casesList := []casePrice.CasePriceInfo{}
	page := c.Query("page")
	if page == "" {
		page = "1"
	}
	casesList = LoadCasePrices(page)
	c.JSON(http.StatusOK, casesList)
}

func LoadCasePrices(page string) []casePrice.CasePriceInfo {
	casesList := []casePrice.CasePriceInfo{}
	var cases casePrice.CasePriceInfo
	query := "SELECT id, case_id, organization_id, price FROM case_prices LIMIT 10 OFFSET ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		return casesList
	}
	rows, error := stmt.Query(page)
	if error != nil {
		return casesList
	}
	for rows.Next() {
		err := rows.Scan(
			&cases.ID,
			&cases.CaseId,
			&cases.OrganizationId,
			&cases.Price,
		)
		if err != nil {
			log.Println(err.Error())
			return casesList
		}
		casesList = append(casesList, cases)
	}
	return casesList
}

func (uc *CasePriceControllerStruct) Update(c *gin.Context) {
	var updateUserQuery = "UPDATE `case_prices` SET"
	var values []interface{}
	var columns []string
	caseId := c.Param("id")
	if caseId == "" {
		log.Println(caseId, "case price id")
		errorsHandler.GinErrorResponseHandler(c, nil)
		return
	}
	var request casePrice.UpdateCasePriceRequest
	if errors := c.ShouldBindJSON(&request); errors != nil {
		log.Println(errors.Error())
		errorsHandler.GinErrorResponseHandler(c, errors)
		return
	}
	getCasePricesUpdateColumns(&request, &columns, &values)
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

func getCasePricesUpdateColumns(o *casePrice.UpdateCasePriceRequest, columns *[]string, values *[]interface{}) {
	*columns = append(*columns, " `case_id` = ? ")
	*values = append(*values, o.CaseId)
	*columns = append(*columns, " `organization_id` = ? ")
	*values = append(*values, o.OrganizationId)
	*columns = append(*columns, " `price` = ? ")
	*values = append(*values, o.Price)
}
