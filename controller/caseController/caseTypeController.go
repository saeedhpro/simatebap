package casesController

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	cases "gitlab.com/simateb-project/simateb-backend/domain/case"
	"gitlab.com/simateb-project/simateb-backend/repository"
	"gitlab.com/simateb-project/simateb-backend/utils/errorsHandler"
	"log"
	"net/http"
	"strings"
)

type CasesControllerInterface interface {
	Create(c *gin.Context)
	GetList(c *gin.Context)
	Get(c *gin.Context)
	Update(c *gin.Context)
}

type CasesControllerStruct struct {
}

func NewCasesController() CasesControllerInterface {
	x := &CasesControllerStruct{
	}
	return x
}

func (uc *CasesControllerStruct) Create(c *gin.Context) {
	var request cases.CreateCaseRequest
	if errors := c.ShouldBindJSON(&request); errors != nil {
		log.Println(errors.Error())
		errorsHandler.GinErrorResponseHandler(c, errors)
		return
	}
	query := "INSERT INTO `cases`(`name`,`parent_id`,`is_main`) VALUES (?,?,?)"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error())
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	defer stmt.Close()
	result, err := stmt.Exec(
		&request.Name,
		&request.ParentId,
		&request.IsMain,
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

func (uc *CasesControllerStruct) Get(c *gin.Context) {
	caseId := c.Param("id")
	if caseId == "" {
		return
	}
	query := "SELECT cases.id id, `cases`.`name` `name`, parent.id parent_id, parent.name parent_name, cases.is_main is_main FROM cases LEFT JOIN cases as parent on cases.parent_id = parent.id WHERE cases.id = ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error(), "prepare")
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	var cases cases.CaseInfo
	result := stmt.QueryRow(caseId)
	err = result.Scan(
		&cases.ID,
		&cases.Name,
		&cases.ParentId,
		&cases.ParentName,
		&cases.IsMain,
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
	c.JSON(http.StatusOK, cases)
}

func (uc *CasesControllerStruct) GetList(c *gin.Context) {
	var casesList []cases.CaseInfo
	casesList = LoadCases()
	c.JSON(http.StatusOK, casesList)
}

func LoadCases() []cases.CaseInfo {
	casesList := []cases.CaseInfo{}
	var cases cases.CaseInfo
	query := "SELECT cases.id id, `cases`.`name` `name`, parent.id parent_id, parent.name parent_name, cases.is_main is_main FROM cases LEFT JOIN cases as parent on cases.parent_id = parent.id"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		return casesList
	}
	rows, error := stmt.Query()
	if error != nil {
		return casesList
	}
	for rows.Next() {
		err := rows.Scan(
			&cases.ID,
			&cases.Name,
			&cases.ParentId,
			&cases.ParentName,
			&cases.IsMain,
		)
		if err != nil {
			log.Println(err.Error())
			return casesList
		}
		casesList = append(casesList, cases)
	}
	return casesList
}

func (uc *CasesControllerStruct) Update(c *gin.Context) {
	var updateUserQuery = "UPDATE `case_type` SET"
	var values []interface{}
	var columns []string
	caseId := c.Param("id")
	if caseId == "" {
		log.Println(caseId, "case id")
		errorsHandler.GinErrorResponseHandler(c, nil)
		return
	}
	var request cases.UpdateCaseRequest
	if errors := c.ShouldBindJSON(&request); errors != nil {
		log.Println(errors.Error())
		errorsHandler.GinErrorResponseHandler(c, errors)
		return
	}
	getCasesUpdateColumns(&request, &columns, &values)
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

func getCasesUpdateColumns(o *cases.UpdateCaseRequest, columns *[]string, values *[]interface{}) {
	if o.Name != "" {
		*columns = append(*columns, " `name` = ? ")
		*values = append(*values, o.Name)
	}
	*columns = append(*columns, " `is_main` = ? ")
	*values = append(*values, o.IsMain)
	*columns = append(*columns, " `parent_id` = ? ")
	*values = append(*values, o.ParentId)
}
