package repository

import (
	"database/sql"
	"gitlab.com/simateb-project/simateb-backend/constant"
)

var DB *sql.DB

func init() {
	mysqlInit()
}

func mysqlInit() {
	dataSourceName := constant.MysqlUserName + ":" + constant.MysqlPassword + "@tcp(" + constant.MysqlHost + ":" + string(constant.MysqlPort) + ")/" + constant.MysqlDatabaseName
	DB, err := sql.Open("mysql", dataSourceName)

	if err != nil {
		panic(err.Error())
	}
	defer DB.Close()
}
