package questionmysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"golang.project/go-fundamentals/gameapp/entity"
	"golang.project/go-fundamentals/gameapp/logger"
	"golang.project/go-fundamentals/gameapp/metrics"
	"golang.project/go-fundamentals/gameapp/pkg/errormessage"
	"golang.project/go-fundamentals/gameapp/pkg/richerror"
	"golang.project/go-fundamentals/gameapp/repository/mysql"
	"strconv"
	"strings"
)

//GetQuestions
/*
	possibleAnswers pattern: (;) is SEPARATOR between any record of possibleAnswers.
	possibleAnswers.id|possibleAnswers.text|possibleAnswers.choice;possibleAnswers.id|possibleAnswers.text|possibleAnswers.choice
*/
func (d *DataBase) GetQuestions(ctx context.Context, questionIDs []uint) ([]entity.Question, error) {
	const operation = "gamemysql.GetQuestions"
	const queryType = "select"

	// Prepare query placeholders
	placeholders := strings.TrimSuffix(strings.Repeat("?,", len(questionIDs)), ",")
	fieldOrder := placeholders
	query := fmt.Sprintf(`
		SELECT 
			q.id,
			q.text,
			q.correct_answer_id,
			q.difficulty,
			q.category,
			GROUP_CONCAT(CONCAT_WS('|', pa.id, pa.text, pa.choice) SEPARATOR ';') AS answers
		FROM questions q
		LEFT JOIN possible_answers pa ON q.id = pa.question_id
		WHERE q.id IN (%s)
		GROUP BY q.id
		ORDER BY FIELD(q.id, %s)
	`, placeholders, fieldOrder)

	args := make([]interface{}, 0, len(questionIDs)*2)
	for _, id := range questionIDs {
		args = append(args, id)
	}
	for _, id := range questionIDs {
		args = append(args, id)
	}

	rows, err := d.dataBase.MysqlConnection.QueryContext(ctx, query, args...)
	metrics.DBQueryCounter.With(prometheus.Labels{"query_type": queryType}).Inc()
	if err != nil {
		metrics.DBFailedQueryCounter.With(prometheus.Labels{"query_type": queryType}).Inc()
		return nil, err
	}
	defer rows.Close()

	var questions []entity.Question

	for rows.Next() {
		var (
			question   entity.Question
			difficulty uint8
			category   string
			answersRaw sql.NullString
		)

		err := rows.Scan(
			&question.Id,
			&question.Text,
			&question.CorrectAnswer,
			&difficulty,
			&category,
			&answersRaw,
		)
		if err != nil {
			return nil, err
		}

		question.Difficulty = entity.QuestionDifficulty(difficulty)
		question.Category = entity.Category(category)

		// Parse PossibleAnswers
		if answersRaw.Valid && answersRaw.String != "" {
			answerStrings := strings.Split(answersRaw.String, ";")
			for _, ansStr := range answerStrings {
				parts := strings.Split(ansStr, "|")
				if len(parts) != 3 {
					continue
				}

				id, _ := strconv.ParseUint(parts[0], 10, 64)
				text := parts[1]
				choiceNum, _ := strconv.ParseUint(parts[2], 10, 8)
				choice := entity.PossibleAnswerChoice(choiceNum)

				if !choice.IsValid() {
					continue
				}

				pa := entity.PossibleAnswer{
					Id:     uint(id),
					Text:   text,
					Choice: choice,
				}
				question.PossibleAnswers = append(question.PossibleAnswers, pa)
			}
		}

		questions = append(questions, question)
	}

	if err := rows.Err(); err != nil {
		metrics.DBFailedQueryCounter.With(prometheus.Labels{"query_type": queryType}).Inc()
		return nil, err
	}

	return questions, nil
}

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

	questionIds, sErr := d.scanQuestionIds(questionIdRows, numberOfQuestion)
	if sErr != nil {
		return make([]uint, 0), richerror.NewRichError(operation).
			WithError(sErr).
			WithMessage(errormessage.ErrorMsgUnexpected).
			WithKind(richerror.KindUnexpected)
	}

	return questionIds, nil
}

func (d *DataBase) scanQuestionIds(questionIdRows *sql.Rows, numberOfQuestion uint) ([]uint, error) {
	questionIds := make([]uint, 0, numberOfQuestion)
	for questionIdRows.Next() {
		questionId, sErr := d.scanQuestionId(questionIdRows)
		if sErr != nil {
			return nil, sErr
		}
		questionIds = append(questionIds, questionId)
	}

	if cErr := checkRowsErr(questionIdRows); cErr != nil {
		return nil, cErr
	}

	return questionIds, nil
}

func (d *DataBase) scanQuestionId(scanner mysql.Scanner) (uint, error) {
	var id uint
	if err := scanner.Scan(&id); err != nil {
		return 0, errors.New(errormessage.ErrorMsgScanQuery)
	}
	return id, nil
}

func checkRowsErr(rows *sql.Rows) error {
	if err := rows.Close(); err != nil {
		return err
	}
	if err := rows.Err(); err != nil {
		return err
	}
	return nil
}
