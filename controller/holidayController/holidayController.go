package HolidayController

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/simateb-project/simateb-backend/domain/holiday"
	"gitlab.com/simateb-project/simateb-backend/repository"
	mysqlQuery "gitlab.com/simateb-project/simateb-backend/repository/mysqlQuery/auth"
	"gitlab.com/simateb-project/simateb-backend/utils/errorsHandler"
	"log"
	"net/http"
	"strings"
)

type HolidayControllerInterface interface {
	Create(c *gin.Context)
	Get(c *gin.Context)
	GetList(c *gin.Context)
	GetListForAdmin(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type HolidayControllerStruct struct {
}

func NewHolidayController() HolidayControllerInterface {
	x := &HolidayControllerStruct{
	}
	return x
}

func (uc *HolidayControllerStruct) Create(c *gin.Context) {
	var createHolidayRequest holiday.CreateHolidayRequest
	if errors := c.ShouldBindJSON(&createHolidayRequest); errors != nil {
		log.Println(errors.Error())
		errorsHandler.GinErrorResponseHandler(c, errors)
		return
	}
	stmt, err := repository.DBS.MysqlDb.Prepare(mysqlQuery.CreateHolidayQuery)
	if err != nil {
		log.Println(err.Error())
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	defer stmt.Close()
	result, err := stmt.Exec(
		createHolidayRequest.HDate,
		createHolidayRequest.OrganizationID,
		createHolidayRequest.Title,
	)
	log.Println(createHolidayRequest.OrganizationID, "org")
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

func (uc *HolidayControllerStruct) Get(c *gin.Context) {
	holidayId := c.Param("id")
	if holidayId == "" {
		return
	}
	stmt, err := repository.DBS.MysqlDb.Prepare(mysqlQuery.GetHolidayQuery)
	if err != nil {
		log.Println(err.Error(), "prepare")
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	var holiday holiday.HolidayInfo
	result := stmt.QueryRow(holidayId)
	err = result.Scan(
		&holiday.ID,
		&holiday.Title,
		&holiday.HDate,
		&holiday.OrganizationID,
		&holiday.OrganizationName,
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
	c.JSON(http.StatusOK, holiday)
}

func (uc *HolidayControllerStruct) GetList(c *gin.Context) {
	startAt := c.Query("start_date")
	//endAt := c.Query("end_date")
	id := c.Param("id")
	query := "SELECT holiday.id id, holiday.title, holiday.hdate, ifnull(holiday.organization_id, 0) organization_id, ifnull(organization.name,'') organization_name FROM holiday LEFT JOIN organization ON holiday.organization_id = organization.id WHERE (holiday.organization_id = ? OR holiday.organization_id IS NULL ) "
	var values []interface{}
	values = append(values, id)
	if startAt != "" && startAt != "null"  {
		query += " AND DATE(holiday.hdate) >= ? "
		values = append(values, startAt)
	}
	//if endAt != "" && endAt != "null"  {
	//	query += " AND DATE(holiday.hdate) <= ? "
	//	values = append(values, endAt)
	//}
	query += " ORDER BY hdate DESC"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error())
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	holidays := []holiday.HolidayInfo{}
	var hd holiday.HolidayInfo
	rows, err := stmt.Query(values...)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusOK, holidays)
		return
	}
	for rows.Next() {
		err := rows.Scan(
			&hd.ID,
			&hd.Title,
			&hd.HDate,
			&hd.OrganizationID,
			&hd.OrganizationName,
		)
		if err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusOK, holidays)
			return
		}
		holidays = append(holidays, hd)
	}
	c.JSON(http.StatusOK, holidays)
}

func (uc *HolidayControllerStruct) GetListForAdmin(c *gin.Context) {
	page := c.Query("page")
	if page == "" {
		page = "1"
	}
	query := "SELECT holiday.id id, holiday.title, holiday.hdate, ifnull(holiday.organization_id, 0) organization_id, ifnull(organization.name, '') organization_name FROM holiday LEFT JOIN organization ON holiday.organization_id = organization.id LIMIT 10 OFFSET ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	rows, err := stmt.Query(page)
	if err != nil {
		log.Println(err.Error())
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	holidays := []holiday.HolidayInfo{}
	var holiday holiday.HolidayInfo
	for rows.Next() {
		err := rows.Scan(
			&holiday.ID,
			&holiday.Title,
			&holiday.HDate,
			&holiday.OrganizationID,
			&holiday.OrganizationName,
		)
		if err != nil {
			log.Println(err.Error())
		}
		holidays = append(holidays, holiday)
	}
	c.JSON(http.StatusOK, holidays)
}

func (uc *HolidayControllerStruct) Update(c *gin.Context) {
	var updateHolidayQuery = "UPDATE `holiday` SET "
	var values []interface{}
	var columns []string
	holidayId := c.Param("id")
	if holidayId == "" {
		errorsHandler.GinErrorResponseHandler(c, nil)
		return
	}
	var updateHolidayRequest holiday.UpdateHolidayRequest
	if errors := c.ShouldBindJSON(&updateHolidayRequest); errors != nil {
		log.Println(errors.Error())
		errorsHandler.GinErrorResponseHandler(c, errors)
		return
	}
	getHolidayUpdateColumns(&updateHolidayRequest, &columns, &values)
	columnsString := strings.Join(columns, ",")
	updateHolidayQuery += columnsString
	updateHolidayQuery += " WHERE `id` = ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(updateHolidayQuery)
	if err != nil {
		log.Println(fmt.Sprintf("Update %s", err.Error()))
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	values = append(values, holidayId)
	_, error := stmt.Exec(values...)
	if error != nil {
		log.Println(error.Error())
		errorsHandler.GinErrorResponseHandler(c, error)
		return
	}
	c.JSON(200, true)
}

func getHolidayUpdateColumns(o *holiday.UpdateHolidayRequest, columns *[]string, values *[]interface{}) {
	*columns = append(*columns, " `hdate` = ? ")
	*values = append(*values, o.HDate)
	if o.Title != "" {
		*columns = append(*columns, " `title` = ? ")
		*values = append(*values, o.Title)
	}
	if o.OrganizationID != 0  {
		*columns = append(*columns, " `organization_id` = ? ")
		*values = append(*values, o.OrganizationID)
	}
}

func (uc *HolidayControllerStruct) Delete(c *gin.Context) {
	HolidayID := c.Param("id")
	if HolidayID == "" {
		return
	}
	stmt, err := repository.DBS.MysqlDb.Prepare(mysqlQuery.DeleteHolidayQuery)
	if err != nil {
		errorsHandler.GinErrorResponseHandler(c, err)
		return
	}
	stmt.QueryRow(HolidayID)
	c.JSON(200, nil)
}
