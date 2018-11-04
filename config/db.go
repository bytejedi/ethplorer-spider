package config

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type TxTable struct {
	Tx        string
	Timestamp int64
	From      string
	To        string
	Value     string
	T         string
}

func DB() *sql.DB {

	user := "root"
	password := "111111"
	host := "127.0.0.1"
	port := "3306"
	dbName := "dacc"

	db, _ := sql.Open("mysql", user+":"+password+"@tcp("+host+":"+port+")/"+dbName+"?charset=utf8")
	err := db.Ping()
	if err != nil {
		log.Panicln(err)
	}
	return db
}
