package doc

import (
	"gitlab.com/simateb-project/simateb-backend/controller/organizationController"
	"gitlab.com/simateb-project/simateb-backend/domain/organization"
	"gitlab.com/simateb-project/simateb-backend/repository"
	"log"
	"strconv"
)

type DocStruct struct {
	ID             int64                              `json:"id"`
	UserID         *int64                             `json:"user_id"`
	OrganizationID *int64                             `json:"organization_id"`
	Name           string                             `json:"name"`
	Path           string                             `json:"path"`
	UserDesc       string                             `json:"user_desc"`
	DoctorDesc     string                             `json:"doctor_desc"`
	Type           string                             `json:"type"`
	Profession     *organization.SimpleProfessionInfo `json:"profession"`
}

type CreateDocStruct struct {
	UserID         int64 `json:"user_id"`
	OrganizationID int64 `json:"organization_id"`
	Name           string `json:"name"`
	Path           string `json:"path"`
	UserDesc       string `json:"user_desc"`
	DoctorDesc     string `json:"doctor_desc"`
	Type           string `json:"type"`
	ProfessionID   int64 `json:"profession_id"`
}

func GetUserDocs(userID string, organizationID string) ([]DocStruct, error) {
	docs := []DocStruct{}
	doc := DocStruct{}
	query := "SELECT id, ifnull(name, '') name, ifnull(path, '') path, ifnull(doctor_desc, '') doctor_desc, ifnull(user_desc, '') user_desc, ifnull(type, '') type FROM docs WHERE user_id = ? AND organization_id = ?"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error(), "prepare")
		return docs, err
	}
	rows, err := stmt.Query(userID, organizationID)
	if err != nil {
		log.Println(err.Error(), "prepare")
		return docs, err
	}
	for rows.Next() {
		err = rows.Scan(
			&doc.ID,
			&doc.Name,
			&doc.Path,
			&doc.UserDesc,
			&doc.DoctorDesc,
			&doc.Type,
		)
		if err != nil {
			log.Println(err.Error(), " :err: ")
			return docs, err
		}
		uid, _ := strconv.ParseInt(userID, 10, 64)
		oid, _ := strconv.ParseInt(organizationID, 10, 64)
		doc.UserID = &uid
		doc.OrganizationID = &oid
		doc.Profession = organizationController.GetProfession(organizationID)
		docs = append(docs, doc)
	}
	return docs, nil
}

func CreateUserDoc(request CreateDocStruct) error {
	query := "INSERT INTO `docs`(`user_id`, `name`, `path`, `organization_id`, `doctor_desc`, `user_desc`, `type`, `proffession_id`) VALUES (?,?,?,?,?,?,?,?) "
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(
		&request.UserID,
		&request.Name,
		&request.Path,
		&request.OrganizationID,
		&request.DoctorDesc,
		&request.UserDesc,
		&request.Type,
		&request.ProfessionID,
	)
	if err != nil {
		return err
	}
	return nil
}

func DeleteDoc(id string) error {
	query := "DELETE FROM `docs` WHERE id = ? "
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}
