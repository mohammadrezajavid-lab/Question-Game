package accesscontrolmysql

import (
	"golang.project/go-fundamentals/gameapp/entity"
	"golang.project/go-fundamentals/gameapp/repository/mysql"
	"time"
)

func scanPermissionRow(scanner mysql.Scanner) (*entity.Permission, error) {
	var createdAt time.Time
	p := entity.NewPermission()

	err := scanner.Scan(&p.Id, &p.Title, &createdAt)

	return p, err
}
