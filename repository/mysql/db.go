package mysql

import (
	"database/sql"
	"fmt"
	"golang.project/go-fundamentals/gameapp/logger"
	"golang.project/go-fundamentals/gameapp/pkg/errormessage"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Config The following structure for the DataBase service config
type Config struct {
	DBDriverName      string        `mapstructure:"database_driver_name"`
	DBConnMaxLifetime time.Duration `mapstructure:"database_conn_max_lifetime"`
	DBMaxOpenConns    int           `mapstructure:"database_max_open_connections"`
	DBMaxIdleConns    int           `mapstructure:"database_max_idle_connections"`
	UserName          string        `mapstructure:"database_user_name"`
	Password          string        `mapstructure:"database_password"`
	DBName            string        `mapstructure:"database_name"`
	Host              string        `mapstructure:"database_host"`
	ParseTime         bool          `mapstructure:"database_parse_time"`
	Port              int           `mapstructure:"database_port"`
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
	db, oErr := sql.Open(dbCfg.DBDriverName, connectionUrl)

	if oErr != nil {
		logger.Panic(oErr, errormessage.ErrorMsgFailedOpenMysqlConn)
	}

	if pErr := db.Ping(); pErr != nil {
		logger.Warn(pErr, errormessage.ErrorMsgConnectionRefusedMysql)
	}

	db.SetConnMaxLifetime(dbCfg.DBConnMaxLifetime)
	db.SetMaxOpenConns(dbCfg.DBMaxOpenConns)
	db.SetMaxIdleConns(dbCfg.DBMaxIdleConns)

	return &DB{MysqlConnection: db}
}
