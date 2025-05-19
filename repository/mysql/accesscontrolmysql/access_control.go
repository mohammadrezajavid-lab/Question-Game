package accesscontrolmysql

import (
	"database/sql"
	"fmt"
	"golang.project/go-fundamentals/gameapp/entity"
	"golang.project/go-fundamentals/gameapp/pkg/errormessage"
	"golang.project/go-fundamentals/gameapp/pkg/richerror"
	"golang.project/go-fundamentals/gameapp/pkg/slice"
	"golang.project/go-fundamentals/gameapp/repository/mysql"
	"log"
	"strings"
	"time"
)

func (d *DataBase) GetUserPermissionsTitle(userId uint, role entity.Role) ([]entity.PermissionTitle, error) {

	const operation = "mysql.accessControl.GetUserPermissionsTitle"

	aclRoleRows, rErr := d.dataBase.MysqlConnection.Query(
		`SELECT * FROM access_controls WHERE functor_type = ? and functor_id = ?`,
		entity.RoleFunctorType,
		role,
	)
	if rErr != nil {

		// log internal server error
		log.Println("aclRoleRows.Query(SELECT * FROM access_controls WHERE functor_type = ? and functor_id = ?)")
		log.Println("alcRoleRows, Error: ", rErr)

		return nil, richerror.NewRichError(operation).
			WithError(rErr).
			WithMessage(errormessage.ErrorMsgUnexpected).
			WithKind(richerror.KindUnexpected)
	}
	defer func(aclRoleRows *sql.Rows) {
		cErr := aclRoleRows.Close()
		if cErr != nil {
			log.Println("can't close aclRoleRows")
		}
	}(aclRoleRows)

	aclUserRows, uErr := d.dataBase.MysqlConnection.Query(
		`SELECT * FROM access_controls WHERE functor_type = ? and functor_id = ?`,
		entity.UserFunctorType,
		userId,
	)
	if uErr != nil {

		// log internal server error
		log.Println("aclUserRows.Query(SELECT * FROM access_controls WHERE functor_type = ? and functor_id = ?)")
		log.Println("aclUserRows, Error: ", uErr)

		return nil, richerror.NewRichError(operation).
			WithError(uErr).
			WithMessage(errormessage.ErrorMsgUnexpected).
			WithKind(richerror.KindUnexpected)
	}
	defer func(aclUserRows *sql.Rows) {
		cErr := aclUserRows.Close()
		if cErr != nil {
			log.Println("can't close aclUserRows")
		}
	}(aclUserRows)

	// scan aclRoleRows and aclUserRows
	aclList := make([]*entity.AccessControl, 0)
	for aclRoleRows.Next() {
		if row, sErr := scanACLRow(aclRoleRows); sErr != nil {

			return nil, richerror.NewRichError(operation).
				WithError(sErr).
				WithMessage(errormessage.ErrorMsgScanQuery).
				WithKind(richerror.KindUnexpected)
		} else {
			aclList = append(aclList, row)
		}
	}

	if aclRoleRows.Err() != nil {

		return nil, richerror.NewRichError(operation).
			WithError(aclRoleRows.Err()).
			WithMessage(errormessage.ErrorMsgScanQuery).
			WithKind(richerror.KindUnexpected)
	}

	for aclUserRows.Next() {
		if row, sErr := scanACLRow(aclUserRows); sErr != nil {

			return nil, richerror.NewRichError(operation).
				WithError(sErr).
				WithMessage(errormessage.ErrorMsgScanQuery).
				WithKind(richerror.KindUnexpected)
		} else {
			aclList = append(aclList, row)
		}
	}
	if aclUserRows.Err() != nil {

		return nil, richerror.NewRichError(operation).
			WithError(aclUserRows.Err()).
			WithMessage(errormessage.ErrorMsgScanQuery).
			WithKind(richerror.KindUnexpected)
	}

	permissionIds := make([]uint, 0)
	for _, a := range aclList {
		if !slice.DoesExist(permissionIds, a.PermissionId) {
			permissionIds = append(permissionIds, a.PermissionId)
		}
	}

	// query to database for get permissionTitle by (...permissionIds)
	if len(permissionIds) == 0 {
		return nil, nil
	}

	args := make([]any, len(permissionIds))
	for index, permId := range permissionIds {
		args[index] = permId
	}

	// warning: this query works if we have one or more permissionId
	placeHolder := "?" + strings.Repeat(",?", len(permissionIds)-1)
	var queryStr = fmt.Sprintf(`SELECT * FROM permissions WHERE id IN(%s)`, placeHolder)
	permissionRows, qErr := d.dataBase.MysqlConnection.Query(queryStr, args...)
	if qErr != nil {

		// log internal server error
		log.Println("permissionRows.Query(SELECT * FROM permissions WHERE id IN(? + strings.Repeat(\",?\", len(permissionIds)-1)))")
		log.Println("permissionRows, Error: ", qErr)

		return nil, richerror.NewRichError(operation).
			WithError(qErr).
			WithMessage(errormessage.ErrorMsgUnexpected).
			WithKind(richerror.KindUnexpected)
	}
	defer func(permissionRows *sql.Rows) {
		cErr := permissionRows.Close()
		if cErr != nil {
			log.Println("can't close permissionTitleRows")
		}
	}(permissionRows)

	var permissionTitles = make([]entity.PermissionTitle, 0)
	for permissionRows.Next() {
		if permission, sErr := scanPermissionRow(permissionRows); sErr != nil {

			return nil, richerror.NewRichError(operation).
				WithError(sErr).
				WithMessage(errormessage.ErrorMsgScanQuery).
				WithKind(richerror.KindUnexpected)
		} else {
			permissionTitles = append(permissionTitles, permission.Title)
		}
	}
	if permissionRows.Err() != nil {

		return nil, richerror.NewRichError(operation).
			WithError(permissionRows.Err()).
			WithMessage(errormessage.ErrorMsgScanQuery).
			WithKind(richerror.KindUnexpected)
	}
	return permissionTitles, nil
}

func scanACLRow(scanner mysql.Scanner) (*entity.AccessControl, error) {
	var createdAt time.Time
	a := entity.NewAccessControl()

	err := scanner.Scan(&a.Id, &a.FunctorId, &a.FunctorType, &a.PermissionId, &createdAt)

	return a, err
}
