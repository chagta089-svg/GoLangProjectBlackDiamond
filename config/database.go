package config

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() error {
	var err error
	DB, err = sql.Open("mysql", "root:1234@tcp(127.0.0.1:3306)/goproject")
	if err != nil {
		return err
	}

	err = DB.Ping()
	if err != nil {
		return err
	}

	fmt.Println("Connected to database successfully")
	return nil
}

func GetDB() *sql.DB {
	return DB
}
