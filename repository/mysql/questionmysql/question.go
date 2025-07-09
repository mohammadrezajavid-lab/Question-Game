package questionmysql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/prometheus/client_golang/prometheus"
	"golang.project/go-fundamentals/gameapp/entity"
	"golang.project/go-fundamentals/gameapp/logger"
	"golang.project/go-fundamentals/gameapp/metrics"
	"golang.project/go-fundamentals/gameapp/pkg/errormessage"
	"golang.project/go-fundamentals/gameapp/pkg/richerror"
	"golang.project/go-fundamentals/gameapp/repository/mysql"
)

func (d *DataBase) GetRandomQuestions(ctx context.Context, category entity.Category, difficulty entity.QuestionDifficulty, numberOfQuestion uint) ([]uint, error) {
	const operation = "gamemysql.GetRandomQuestions"
	const queryType = "select"

	questionIdRows, err := d.dataBase.MysqlConnection.QueryContext(
		ctx,
		`SELECT id FROM questions WHERE category=? AND difficulty=? ORDER BY RAND() LIMIT ?`,
		category, difficulty, numberOfQuestion,
	)
	metrics.DBQueryCounter.With(prometheus.Labels{"query_type": queryType}).Inc()

	if err != nil {
		logger.Warn(err, errormessage.ErrorMsgFailedExecuteQuery)
		metrics.DBFailedQueryCounter.With(prometheus.Labels{"query_type": queryType}).Inc()

		return make([]uint, 0), richerror.NewRichError(operation).
			WithError(err).
			WithMessage(errormessage.ErrorMsgUnexpected).
			WithKind(richerror.KindUnexpected)
	}

	questionIds, sErr := scanQuestionIds(questionIdRows, numberOfQuestion)
	if sErr != nil {
		return make([]uint, 0), richerror.NewRichError(operation).
			WithError(sErr).
			WithMessage(errormessage.ErrorMsgUnexpected).
			WithKind(richerror.KindUnexpected)
	}

	return questionIds, nil
}

func scanQuestionIds(questionIdRows *sql.Rows, numberOfQuestion uint) ([]uint, error) {
	questionIds := make([]uint, numberOfQuestion)
	for questionIdRows.Next() {
		questionId, sErr := scanQuestionId(questionIdRows)
		if sErr != nil {
			return nil, sErr
		}
		questionIds = append(questionIds, questionId)
	}
	if cErr := questionIdRows.Close(); cErr != nil {
		return nil, cErr
	}
	if eErr := questionIdRows.Err(); eErr != nil {
		return nil, eErr
	}

	return questionIds, nil
}

func scanQuestionId(scanner mysql.Scanner) (uint, error) {
	var id uint
	if err := scanner.Scan(&id); err != nil {
		return 0, errors.New(errormessage.ErrorMsgScanQuery)
	}
	return id, nil
}
