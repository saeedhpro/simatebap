package user

import (
	"gitlab.com/simateb-project/simateb-backend/domain/organization"
	"gitlab.com/simateb-project/simateb-backend/repository"
	mysqlQuery "gitlab.com/simateb-project/simateb-backend/repository/mysqlQuery/auth"
	"log"
)

func GetUserByID(userId int64) (*organization.OrganizationUser,  error) {
	var user organization.OrganizationUser
	stmt, err := repository.DBS.MysqlDb.Prepare(mysqlQuery.GetUserOrganizationQuery)
	if err != nil {
		log.Println(err.Error(), "prepare")
		return nil , err
	}
	result := stmt.QueryRow(userId)
	err = result.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.LastLogin,
		&user.Created,
		&user.Tel,
		&user.UserGroupID,
		&user.UserGroupName,
		&user.BirthDate,
		&user.OrganizationID,
		&user.OrganizationName,
		&user.Relation,
		&user.Description,
	)
	return &user, nil
}