package vip

import (
	"database/sql"
	"gitlab.com/simateb-project/simateb-backend/repository"
	"log"
)

type VipSchedule struct {
	ID             int64         `json:"id"`
	AppCount       int64         `json:"app_count"`
	SiteCount      int64         `json:"site_count"`
	DoctorCount    int64         `json:"doctor_count"`
	StartAt        *sql.NullTime `json:"start_at"`
	EndAt          *sql.NullTime `json:"end_at"`
	OrganizationID int64         `json:"organization_id"`
	CaseType       string        `json:"case_type"`
}

type VipCase struct {
	ID             int64  `json:"id"`
	OrganizationID int64  `json:"organization_id"`
	Name           string `json:"name"`
}

func GetVipScheduleByID(id int64) (*VipSchedule, error) {
	query := "SELECT `id`, `app_count`, `doctor_count`, `site_count`, `start_at`, `end_at`, `organization_id`, `case_type` FROM `vip_schedule` WHERE id = ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		return nil, err
	}
	result := stmt.QueryRow(id)
	var vip VipSchedule
	err = result.Scan(
		&vip.ID,
		&vip.AppCount,
		&vip.DoctorCount,
		&vip.SiteCount,
		&vip.StartAt,
		&vip.EndAt,
		&vip.OrganizationID,
		&vip.CaseType,
	)
	return &vip, nil
}

func GetVipScheduleCaseByID(id int64) (*VipCase, error) {
	query := "SELECT `id`, `name`, `organization_id` FROM `vip_cases` WHERE id = ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		return nil, err
	}
	result := stmt.QueryRow(id)
	var vip VipCase
	err = result.Scan(
		&vip.ID,
		&vip.Name,
		&vip.OrganizationID,
	)
	return &vip, nil
}

func GetVipScheduleCaseByName(name string) (*VipCase, error) {
	query := "SELECT `id`, `name`, `organization_id` FROM `vip_cases` WHERE `name` = ? LIMIT 1"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	result := stmt.QueryRow(name)
	var vip VipCase
	err = result.Scan(
		&vip.ID,
		&vip.Name,
		&vip.OrganizationID,
	)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return &vip, nil
}
