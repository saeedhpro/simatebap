package appointment

import (
	"fmt"
	"gitlab.com/simateb-project/simateb-backend/domain/appointment"
	"gitlab.com/simateb-project/simateb-backend/repository"
	"gitlab.com/simateb-project/simateb-backend/utils/auth"
	"log"
)

type AppointmentRepositoryInterface interface {
	LoadQueWithOrganization(organizationID int64, startDate string, endDate string, status string, user *auth.UserClaims) []appointment.SimpleAppointmentInfo
	LoadTotalsByDayAndCase(organizationID int64, caseType string, startDate string, endDate string) []appointment.TotalLimit
	GetOrganizationWorkHour(organizationId int64) appointment.WorkHour
}

type AppointmentRepositoryStruct struct {
}

func NewAppointmentRepository() *AppointmentRepositoryStruct {
	x := &AppointmentRepositoryStruct{
	}
	return x
}

func (uc *AppointmentRepositoryStruct) LoadQueWithOrganization(organizationID int64, startDate string, endDate string, status string, user *auth.UserClaims) []appointment.SimpleAppointmentInfo {
	query := fmt.Sprintf("SELECT appointment.id, appointment.start_at, user.fname user_fname, user.lname user_lname, user.gender user_gender, user.id user_id, user.tel mobile, appointment.status, appointment.case_type, ct.duration duration, appointment.is_vip FROM appointment LEFT JOIN user ON appointment.user_id = user.id LEFT JOIN case_type ct on ct.organization_id = appointment.organization_id AND ct.name = appointment.case_type WHERE `status` IN (%s) AND  DATE(appointment.start_at) >= DATE(\"%s\") AND DATE(appointment.start_at) <= DATE(\"%s\") AND appointment.organization_id = ? ORDER BY appointment.start_at ASC", status, startDate, endDate)
	log.Println(query)
	log.Println(organizationID)
	var values []interface{}
	values = append(values, organizationID)
	var appointments []appointment.SimpleAppointmentInfo
	var appointment appointment.SimpleAppointmentInfo
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error())
		return appointments
	}
	rows, error := stmt.Query(values...)
	if error != nil {
		log.Println(error.Error())
		return appointments
	}
	for rows.Next() {
		log.Println("appointment")
		err = rows.Scan(
			&appointment.ID,
			&appointment.StartAt,
			&appointment.UserFName,
			&appointment.UserLName,
			&appointment.UserGender,
			&appointment.UserID,
			&appointment.Mobile,
			&appointment.Status,
			&appointment.CaseType,
			&appointment.Duration,
			&appointment.IsVip,
		)
		log.Println(appointment.ID)
		if err != nil {
			log.Println(err.Error())
		}
		appointments = append(appointments, appointment)
	}
	return appointments
}

func (uc *AppointmentRepositoryStruct) LoadTotals(organizationID int64, startDate string, endDate string) []appointment.TotalLimit {
	query := "SELECT count(appointment.id) total," +
		" DATE(appointment.start_at) `date` " +
		" FROM appointment " +
		" WHERE appointment.status != 3 " +
		" AND appointment.organization_id = ? " +
		" AND DATE(appointment.start_at) <= DATE(?)" +
		" AND DATE(appointment.start_at) >= DATE(?)" +
		" GROUP BY DATE(appointment.start_at)" +
		" ORDER BY `date` ASC"

	var values []interface{}
	values = append(values, organizationID, startDate, endDate)
	var totals []appointment.TotalLimit
	var total appointment.TotalLimit
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		return totals
	}
	rows, error := stmt.Query(values)
	if error != nil {
		return totals
	}
	for rows.Next() {
		rows.Scan(
			&total.Total,
			&total.Date,
		)
		totals = append(totals, total)
	}
	return totals
}

func (uc *AppointmentRepositoryStruct) LoadTotalsByDayAndCase(organizationID int64, caseType string, startDate string, endDate string) []appointment.TotalLimit {
	query := "SELECT count(appointment.id) total," +
		" DATE(appointment.start_at) `date` " +
		" FROM appointment " +
		" WHERE appointment.status != 3 " +
		" AND appointment.organization_id = ? " +
		" AND appointment.case_type = ?" +
		" AND DATE(appointment.start_at) <= DATE(?)" +
		" AND DATE(appointment.start_at) >= DATE(?)" +
		" GROUP BY DATE(appointment.start_at)" +
		" ORDER BY `date` ASC"

	var values []interface{}
	values = append(values, organizationID, caseType, startDate, endDate)
	var totals []appointment.TotalLimit
	var total appointment.TotalLimit
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		return totals
	}
	rows, error := stmt.Query(values)
	if error != nil {
		return totals
	}
	for rows.Next() {
		rows.Scan(
			&total.Total,
			&total.Date,
		)
		totals = append(totals, total)
	}
	return totals
}

func (uc *AppointmentRepositoryStruct) GetOrganizationWorkHour(organizationId int64) appointment.WorkHour {
	query := "SELECT work_hour_start start, work_hour_end end FROM organization WHERE organization_id = ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	var workHour appointment.WorkHour
	if err != nil {
		return workHour
	}
	result := stmt.QueryRow(organizationId)
	result.Scan(
		&workHour.Start,
		&workHour.End,
	)
	return workHour
}

