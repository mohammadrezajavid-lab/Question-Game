package accesscontrolmysql

import "golang.project/go-fundamentals/gameapp/repository/mysql"

type DataBase struct {
	dataBase *mysql.DB
}

func NewDataBase(dataBase *mysql.DB) *DataBase {
	return &DataBase{dataBase: dataBase}
}
