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
	InsertAll(c *gin.Context)
	GetList(c *gin.Context)
	GetListByProf(c *gin.Context)
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

func (uc *CasesControllerStruct) InsertAll(c *gin.Context) {
	query := "INSERT INTO `cases`(`name`, `val`, `is_main`, `profession_id`) VALUES "
	data := []map[string]string{
		{"name": "Periapical (Parallel)", "val": "INTRA ORAL - Periapical (Parallel)", "is_main": "1", "profession_id": "3"},
		{"name": "P.A Full mouth", "val": "INTRA ORAL - P.A Full mouth", "is_main": "1", "profession_id": "3"},
		{"name": "Occlusal", "val": "INTRA ORAL - Occlusal", "is_main": "1", "profession_id": "3"},
		{"name": "Maxilla", "val": "INTRA ORAL - Maxilla", "is_main": "1", "profession_id": "3"},
		{"name": "Mandible", "val": "INTRA ORAL - Mandible", "is_main": "1", "profession_id": "3"},
		{"name": "BiteWing", "val": "INTRA ORAL - BiteWing", "is_main": "1", "profession_id": "3"},
		{"name": "Left premolars", "val": "INTRA ORAL - Left premolars", "is_main": "1", "profession_id": "3"},
		{"name": "Right premolars", "val": "INTRA ORAL - Right premolars", "is_main": "1", "profession_id": "3"},
		{"name": "Left molars", "val": "INTRA ORAL - Left molars", "is_main": "1", "profession_id": "3"},
		{"name": "Right molars", "val": "INTRA ORAL - Right molars", "is_main": "1", "profession_id": "3"},
		{"name": "Panoramic (OPG)", "val": "EXTRA ORAL - Panoramic (OPG)", "is_main": "1", "profession_id": "3"},
		{"name": "Lateral Ceph.", "val": "EXTRA ORAL - Lateral Ceph", "is_main": "1", "profession_id": "3"},
		{"name": "Standard", "val": "EXTRA ORAL - Standard", "is_main": "1", "profession_id": "3"},
		{"name": "NHP", "val": "EXTRA ORAL - NHP", "is_main": "1", "profession_id": "3"},
		{"name": "P.A Ceph.", "val": "EXTRA ORAL - P.A Ceph", "is_main": "1", "profession_id": "3"},
		{"name": "Water's", "val": "EXTRA ORAL - Water's", "is_main": "1", "profession_id": "3"},
		{"name": "TMJ", "val": "EXTRA ORAL - TMJ", "is_main": "1", "profession_id": "3"},
		{"name": "Reverse Towne", "val": "EXTRA ORAL - Reverse Towne", "is_main": "1", "profession_id": "3"},
		{"name": "SMV", "val": "EXTRA ORAL - SMV", "is_main": "1", "profession_id": "3"},
		{"name": "Lateral Oblique", "val": "EXTRA ORAL - Lateral Oblique", "is_main": "1", "profession_id": "3"},
		{"name": "R", "val": "EXTRA ORAL - Lateral Oblique - R", "is_main": "1", "profession_id": "3"},
		{"name": "CBCT", "val": "CBCT - CBCT", "is_main": "1", "profession_id": "3"},
		{"name": "Maxilla", "val": "CBCT - Maxilla", "is_main": "1", "profession_id": "3"},
		{"name": "Implant", "val": "CBCT - Implant", "is_main": "1", "profession_id": "3"},
		{"name": "Paranasal sinus and nasal fossa", "val": "CBCT - Paranasal sinus and nasal fossa", "is_main": "1", "profession_id": "3"},
		{"name": "TMJ", "val": "CBCT - TMJ", "is_main": "1", "profession_id": "3"},
		{"name": "L", "val": "CBCT - TMJ - L", "is_main": "1", "profession_id": "3"},
		{"name": "R", "val": "CBCT - TMJ - R", "is_main": "1", "profession_id": "3"},
		{"name": "Lesion", "val": "CBCT - Lesion", "is_main": "1", "profession_id": "3"},
		{"name": "Periapical (Parallel)", "val": "INTRA ORAL - Periapical (Parallel)", "is_main": "1", "profession_id": "1"},
		{"name": "P.A Full mouth", "val": "INTRA ORAL - P.A Full mouth", "is_main": "1", "profession_id": "1"},
		{"name": "Occlusal", "val": "INTRA ORAL - Occlusal", "is_main": "false", "profession_id": "1"},
		{"name": "Maxilla", "val": "INTRA ORAL - Maxilla", "is_main": "1", "profession_id": "1"},
		{"name": "Mandible", "val": "INTRA ORAL - Mandible", "is_main": "1", "profession_id": "1"},
		{"name": "BiteWing", "val": "INTRA ORAL - BiteWing", "is_main": "false", "profession_id": "1"},
		{"name": "Left premolars", "val": "INTRA ORAL - Left premolars", "is_main": "1", "profession_id": "1"},
		{"name": "Right premolars", "val": "INTRA ORAL - Right premolars", "is_main": "1", "profession_id": "1"},
		{"name": "Left molars", "val": "INTRA ORAL - Left molars", "is_main": "1", "profession_id": "1"},
		{"name": "Right molars", "val": "INTRA ORAL - Right molars", "is_main": "1", "profession_id": "1"},
		{"name": "Panoramic (OPG)", "val": "EXTRA ORAL - Panoramic (OPG)", "is_main": "1", "profession_id": "1"},
		{"name": "Lateral Ceph", "val": "EXTRA ORAL - Lateral Ceph", "is_main": "0", "profession_id": "1"},
		{"name": "Standard", "val": "EXTRA ORAL - Standard", "is_main": "1", "profession_id": "1"},
		{"name": "NHP", "val": "EXTRA ORAL - NHP", "is_main": "1", "profession_id": "1"},
		{"name": "P.A Ceph.", "val": "EXTRA ORAL - P.A Ceph", "is_main": "1", "profession_id": "1"},
		{"name": "Water's", "val": "EXTRA ORAL - Water's", "is_main": "1", "profession_id": "1"},
		{"name": "TMJ", "val": "EXTRA ORAL - TMJ", "is_main": "1", "profession_id": "1"},
		{"name": "Reverse Towne", "val": "EXTRA ORAL - Reverse Towne", "is_main": "1", "profession_id": "1"},
		{"name": "SMV", "val": "EXTRA ORAL - SMV", "is_main": "1", "profession_id": "1"},
		{"name": "Lateral Oblique", "val": "EXTRA ORAL - Lateral Oblique", "is_main": "0", "profession_id": "1"},
		{"name": "L", "val": "EXTRA ORAL - Lateral Oblique - L", "is_main": "1", "profession_id": "1"},
		{"name": "R", "val": "EXTRA ORAL - Lateral Oblique - R", "is_main": "1", "profession_id": "1"},
		{"name": "CBCT", "val": "CBCT - CBCT", "is_main": "0", "profession_id": "1"},
		{"name": "Maxilla", "val": "CBCT - Maxilla", "is_main": "1", "profession_id": "1"},
		{"name": "Mandible", "val": "CBCT - Mandible", "is_main": "1", "profession_id": "1"},
		{"name": "Implant", "val": "CBCT - Implant", "is_main": "1", "profession_id": "1"},
		{"name": "Impacted tooth", "val": "CBCT - Impacted tooth", "is_main": "1", "profession_id": "1"},
		{"name": "Paranasal sinus and nasal fossa", "val": "CBCT - Paranasal sinus and nasal fossa", "is_main": "1", "profession_id": "1"},
		{"name": "TMJ", "val": "CBCT - TMJ", "is_main": "0", "profession_id": "1"},
		{"name": "L", "val": "CBCT - TMJ - L", "is_main": "1", "profession_id": "1"},
		{"name": "R", "val": "CBCT - TMJ - R", "is_main": "1", "profession_id": "1"},
		{"name": "Lesion", "val": "CBCT - Lesion", "is_main": "1", "profession_id": "1"},
	}
	vals := []interface{}{}
	for _, row := range data {
		query += "(?, ?, ?, ?),"
		vals = append(vals, row["name"], row["val"], row["is_main"], row["profession_id"])
	}
	query = query[0 : len(query)-1]
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error())
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(vals...)
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

func (uc *CasesControllerStruct) GetListByProf(c *gin.Context) {
	id := c.Query("prof")
	casesList := LoadCasesByProf(id)
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

func LoadCasesByProf(id string) []cases.ProfessionCaseInfo {
	casesList := []cases.ProfessionCaseInfo{}
	var cases cases.ProfessionCaseInfo
	query := "SELECT ifnull(radiology_cases, '') radiology_cases, ifnull(photography_cases, '') photography_cases FROM profession where id = ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		return casesList
	}
	rows, err := stmt.Query(id)
	if err != nil {
		return casesList
	}
	for rows.Next() {
		err := rows.Scan(
			&cases.RadiologyCases,
			&cases.PhotographyCases,
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
