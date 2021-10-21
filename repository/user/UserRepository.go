package user

import (
	"gitlab.com/simateb-project/simateb-backend/domain/organization"
	"gitlab.com/simateb-project/simateb-backend/helper"
	"gitlab.com/simateb-project/simateb-backend/repository"
	mysqlQuery "gitlab.com/simateb-project/simateb-backend/repository/mysqlQuery/auth"
	"log"
	"time"
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
		&user.Info,
		&user.Tel1,
		&user.Nid,
		&user.Address,
		&user.Introducer,
		&user.Gender,
		&user.FileID,
	)
	if user.BirthDate != nil {
		year, _, _, _, _, _ := helper.TimeDiff(user.BirthDate.Time, time.Now())
		user.Birth = year
	}
	if err != nil {
		log.Println(err.Error())
	}
	return &user, nil
}