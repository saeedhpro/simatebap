package repository

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var (
	DBS dbs
	//host = "simatebdb:3306"
	//host     = "localhost:3306"
	////schema = "simateb"
	//schema   = "newsimateb"
	////password = "LTIjewAWSKF9nwJGUDTgjEJ6"
	//password = ""
	//username = "root"
)

type dbs struct {
	MysqlDb *sql.DB
}

func Init() {
	mysqlInit()
}

func mysqlInit() {
	username := "doc"
	password := "sdf345kddf@#$@#dsf214"
	schema := "newsimateb"
	host := "localhost:3306"
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true", username, password, host, schema)
	db, err := sql.Open("mysql", dataSourceName)

	if err != nil {
		panic(err.Error())
	}

	if err := db.Ping(); err != nil {
		log.Println(err.Error())
		panic(err.Error())
	}

	DBS.MysqlDb = db
	log.Println("Database Connected")
}
