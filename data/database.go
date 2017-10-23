package data

import (
	"log"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
)

const (
	DBHost  = "10.10.10.10:3306"
	DBUser  = "root"
	DBPass  = "root"
	DBDbase = "chitchat"
	Port    = ":8080"
)

var Db *sql.DB

func init() {
	var err error

	dbConn := fmt.Sprintf("%s:%s@tcp(%s)/%s", DBUser, DBPass, DBHost, DBDbase)
	Db, err = sql.Open("mysql", dbConn)

	if err != nil {
		log.Fatal(err)
	}
}
