package casePrice

import (
	"database/sql"
	"gitlab.com/simateb-project/simateb-backend/repository"
	"log"
)

type CasePriceInfo struct {
	ID             int64   `json:"id"`
	CaseId         int64   `json:"case_id"`
	OrganizationId int64   `json:"organization_id"`
	Price          float64 `json:"price"`
}

type CreateCasePriceRequest struct {
	CaseId         int64   `json:"case_id"`
	OrganizationId int64   `json:"organization_id"`
	Price          float64 `json:"price"`
}

type UpdateCasePriceRequest struct {
	CaseId         int64   `json:"case_id"`
	OrganizationId int64   `json:"organization_id"`
	Price          float64 `json:"price"`
}

func GetPriceOfCase(caseID int64, organizationID int64) (float64, error) {
	query := "SELECT price form case_prices WHERE case_id = ? AND organization_id = ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error(), "prepare")
		return 0.0, err
	}
	price := 0.0
	result := stmt.QueryRow(caseID, organizationID)
	err = result.Scan(
		&price,
	)
	if err != nil {
		log.Println(err.Error())
		if err == sql.ErrNoRows {
			return 0.0, err
		}
		return 0.0, err
	}
	return price, nil
}
