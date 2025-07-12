package gameservice

import (
	"context"
	"fmt"
	"golang.project/go-fundamentals/gameapp/logger"
	"golang.project/go-fundamentals/gameapp/param/gameparam"
	"golang.project/go-fundamentals/gameapp/pkg/errormessage"
	"golang.project/go-fundamentals/gameapp/pkg/protobufencodedecode"
	"golang.project/go-fundamentals/gameapp/pkg/richerror"
)

func (s *Service) GetQuiz(ctx context.Context, req gameparam.GetQuizRequest) (gameparam.GetQuizResponse, error) {
	const operation = "gameservice.GetQuiz"

	key := fmt.Sprintf("game_quiz:%d", req.GameId)
	value, err := s.kvStore.Get(ctx, key)
	if err != nil {
		logger.Warn(err, fmt.Sprintf("failed to Get quiz for game_id: %d", req.GameId))
		return gameparam.GetQuizResponse{}, richerror.NewRichError(operation).
			WithError(err).
			WithMessage(errormessage.ErrorMsgRecordNotFound).
			WithKind(richerror.KindNotFound)
	}

	gameQuiz := protobufencodedecode.DecodeGameSvcGameQuiz(value)

	questionsParam := make([]gameparam.Question, 0, len(gameQuiz.Questions))
	questionParam := gameparam.Question{
		Id:              0,
		Text:            "",
		PossibleAnswers: nil,
		CorrectAnswer:   0,
		Difficulty:      0,
		Category:        "",
	}
	for _, q := range gameQuiz.Questions {
		questionParam.Id = q.Id
		questionParam.Text = q.Text
		questionParam.PossibleAnswers = gameparam.MapPossibleAnswersToParam(q.PossibleAnswers)
		questionParam.CorrectAnswer = q.CorrectAnswer
		questionParam.Difficulty = uint8(q.Difficulty)
		questionParam.Category = string(q.Category)

		questionsParam = append(questionsParam, questionParam)
	}

	return gameparam.GetQuizResponse{
		GameId:    req.GameId,
		Questions: questionsParam,
	}, nil
}
