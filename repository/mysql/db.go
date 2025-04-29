package mysql

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	MysqlConnection *sql.DB
}

func NewDB() *DB {

	db, err := sql.Open("mysql", "game_app:game_app_pass@(127.0.0.1:3308)/game_app_db")
	if err != nil {
		panic(fmt.Errorf("can't open mysql db: %v", err))
	}

	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return &DB{MysqlConnection: db}
}
