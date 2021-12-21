package cases

import (
	"database/sql"
	"gitlab.com/simateb-project/simateb-backend/repository"
	"log"
)

type CaseInfo struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	ParentId   int64  `json:"parent_id"`
	ParentName int64  `json:"parent_name"`
	IsMain     bool   `json:"is_main"`
}

type ProfessionCaseInfo struct {
	PhotographyCases string `json:"photography_cases"`
	RadiologyCases   string  `json:"radiology_cases"`
}

type CreateCaseRequest struct {
	Name     string `json:"name"`
	ParentId int64  `json:"parent_id"`
	IsMain   bool   `json:"is_main"`
}

type UpdateCaseRequest struct {
	Name     string `json:"name"`
	ParentId int64  `json:"parent_id"`
	IsMain   bool   `json:"is_main"`
}

func GetCaseByName(name string) (*CaseInfo, error) {
	query := "SELECT cases.id id, `cases`.`name` `name`, parent.id parent_id, parent.name parent_name, cases.is_main is_main FROM cases LEFT JOIN cases as parent on cases.parent_id = parent.id WHERE cases.name = ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error(), "prepare")
		return nil, err
	}
	var cases CaseInfo
	result := stmt.QueryRow(name)
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
			return nil, err
		}
		return nil, err
	}
	return &cases, nil
}
