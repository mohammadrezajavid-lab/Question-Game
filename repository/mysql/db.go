package mysql

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	UserName  string
	Password  string
	DBName    string
	Host      string
	ParseTime bool
	Port      int
}

func NewConfig(userName, password, dbName, host string, parseTime bool, port int) Config {
	return Config{
		UserName:  userName,
		Password:  password,
		DBName:    dbName,
		Host:      host,
		ParseTime: parseTime,
		Port:      port,
	}
}

type DB struct {
	config          Config
	MysqlConnection *sql.DB
}

func NewDB(dbCfg Config) *DB {

	connectionUrl := fmt.Sprintf(
		"%s:%s@(%s:%d)/%s?parseTime=%v",
		dbCfg.UserName,
		dbCfg.Password,
		dbCfg.Host,
		dbCfg.Port,
		dbCfg.DBName,
		dbCfg.ParseTime,
	)
	db, err := sql.Open("mysql", connectionUrl)
	if err != nil {
		panic(fmt.Errorf("can't open mysql db: %v", err))
	}

	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return &DB{MysqlConnection: db}
}
