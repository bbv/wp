package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type DBConfig struct {
	Host                      string `yaml:"host"`
	Port                      int    `yaml:"port"`
	Login, Password, Database string
}

var DB *sql.DB

func Init(dbconfig DBConfig) error {
	var err error
	DB, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", dbconfig.Login, dbconfig.Password, dbconfig.Host, dbconfig.Port, dbconfig.Database))
	if err != nil {
		return err
	}
	return nil
}
