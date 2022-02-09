package organization

import (
	"gitlab.com/simateb-project/simateb-backend/repository"
	"log"
)

type OrgEmployeeCommission struct {
	UserID         int64 `json:"user_id"`
	OrganizationID int64 `json:"organization_id"`
	Commission     int   `json:"commission"`
}

func GetOrgEmployeeCommissionList(orgID int64) []OrgEmployeeCommission {
	list := []OrgEmployeeCommission{}
	query := "SELECT `user_id`, `organization_id`, `commission` FROM `user_org_prices` WHERE organization_id = >"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		return list
	}
	rows, err := stmt.Query(orgID)
	if err != nil {
		return list
	}
	item := OrgEmployeeCommission{}
	for rows.Next() {
		err = rows.Scan(
			item.UserID,
			item.OrganizationID,
			item.Commission,
		)
		if err != nil {
			return list
		}
		list = append(list, item)
	}
	return list
}

func AddTransfer(orgID int64, appID int64, userID int64, staffID int64, amount float64) (float64, error) {
	query := "INSERT INTO `transfer`(`organization_id`, `appointment_id`, `to_id`, `staff_id`, `amount`, `status`) VALUES (?,?,?,?,?,?)"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	status := 1
	_, err = stmt.Exec(orgID, appID, userID, staffID, amount, status)
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	return amount, nil
}