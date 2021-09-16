package dashboardRepository

import (
	"gitlab.com/simateb-project/simateb-backend/domain/organization"
	"gitlab.com/simateb-project/simateb-backend/repository"
	"log"
)

func GetTodayAppointments() int {
	var query = "SELECT COUNT(*) FROM `appointment` WHERE `appointment`.`start_at` >= now() - INTERVAL 1 DAY AND `appointment`.`status` = 2"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error(), "prepare")
		return 0
	}
	result := stmt.QueryRow()
	count := 0
	err = result.Scan(
		count,
	)
	if err != nil {
		log.Println(err.Error(), "prepare")
		return 0
	}
	return count
}

func GetLastOnLineUsers() ([]organization.LastLoginUserInfo, error) {
	var query = "SELECT id, fname, lname, tel, user_group_id, last_login FROM `user` WHERE `user`.`last_login` >= now() - INTERVAL 1 DAY"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error(), "prepare")
		return nil, err
	}
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	var user organization.LastLoginUserInfo
	var users []organization.LastLoginUserInfo
	for rows.Next() {
		rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Tel,
			&user.UserGroupID,
			&user.LastLogin,
		)
		users = append(users, user)
	}
	return users, nil
}

func GetUnknownGenderUsersCount() int {
	var query = "SELECT COUNT(*) FROM `user` WHERE `gender` IS NULL"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error(), "prepare")
		return 0
	}
	result := stmt.QueryRow()
	count := 0
	err = result.Scan(
		count,
	)
	if err != nil {
		log.Println(err.Error(), "prepare")
		return 0
	}
	return count
}
func GetFemaleUsersCount() int {
	var query = "SELECT COUNT(*) FROM `user` WHERE `gender` = 'FEMALE'"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error(), "prepare")
		return 0
	}
	result := stmt.QueryRow()
	count := 0
	err = result.Scan(
		count,
	)
	if err != nil {
		log.Println(err.Error(), "prepare")
		return 0
	}
	return count
}
func GetMaleUsersCount() int {
	var query = "SELECT COUNT(*) FROM `user` WHERE `gender` = 'MALE'"
	stmt, err := repository.DBS.MysqlDb.Prepare(query)
	if err != nil {
		log.Println(err.Error(), "prepare")
		return 0
	}
	result := stmt.QueryRow()
	count := 0
	err = result.Scan(
		count,
	)
	if err != nil {
		log.Println(err.Error(), "prepare")
		return 0
	}
	return count
}
