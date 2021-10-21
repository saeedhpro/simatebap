package userOrgPrice

import (
	"database/sql"
	"gitlab.com/simateb-project/simateb-backend/repository"
	"log"
)

type UserOrgPriceInfo struct {
	ID             int64   `json:"id"`
	UserID         int64   `json:"user_id"`
	OrganizationID int64   `json:"organization_id"`
	Commission     float64 `json:"commission"`
}

type CreateUserOrgPriceRequest struct {
	UserID         int64   `json:"user_id"`
	OrganizationID int64   `json:"organization_id"`
	Commission     float64 `json:"commission"`
}

type UpdateUserOrgPriceRequest struct {
	UserID         int64   `json:"user_id"`
	OrganizationID int64   `json:"organization_id"`
	Commission     float64 `json:"commission"`
}

func GetPriceOfUser(userID int64, organizationID int64) (float64, error) {
	query := "SELECT price form user_org_prices WHERE user_id = ? AND organization_id = ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error(), "prepare")
		return 0.0, err
	}
	price := 0.0
	result := stmt.QueryRow(userID, organizationID)
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
