package protobufencodedecode

import (
	"encoding/base64"
	"github.com/prometheus/client_golang/prometheus"
	"golang.project/go-fundamentals/gameapp/contract/goprotobuf/game"
	"golang.project/go-fundamentals/gameapp/entity"
	"golang.project/go-fundamentals/gameapp/logger"
	"golang.project/go-fundamentals/gameapp/metrics"
	"golang.project/go-fundamentals/gameapp/pkg/slice"
	"google.golang.org/protobuf/proto"
)

func EncodeGameSvcGameQuiz(quiz entity.GameQuiz) string {
	const op = "protobufencodedecode.EncodeGameSvcGameQuiz"

	protoQuiz := game.GameQuiz{
		GameId:    uint64(quiz.GameId),
		PlayerIds: slice.MapFromUintToUint64(quiz.PlayerIds),
		Questions: mapQuestionsToProto(quiz.Questions),
	}

	data, err := proto.Marshal(&protoQuiz)
	if err != nil {
		metrics.EncodeToProtobufFailedCounter.With(prometheus.Labels{"encoder_name": op}).Inc()
		logger.Error(err, "Marshal GameQuiz Error")
		return ""
	}

	return base64.StdEncoding.EncodeToString(data)
}

func DecodeGameSvcGameQuiz(data string) entity.GameQuiz {
	const op = "protobufencodedecode.DecodeGameSvcGameQuiz"

	raw, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		metrics.DecodeFromProtobufFailedCounter.With(prometheus.Labels{"decoder_name": op}).Inc()
		logger.Error(err, "Base64 decode error")
		return entity.GameQuiz{}
	}

	var protoQuiz game.GameQuiz
	if err := proto.Unmarshal(raw, &protoQuiz); err != nil {
		metrics.DecodeFromProtobufFailedCounter.With(prometheus.Labels{"decoder_name": op}).Inc()
		logger.Error(err, "Unmarshal GameQuiz Error")
		return entity.GameQuiz{}
	}

	return entity.GameQuiz{
		GameId:    uint(protoQuiz.GameId),
		PlayerIds: slice.MapFromUint64ToUint(protoQuiz.PlayerIds),
		Questions: mapProtoQuestions(protoQuiz.Questions),
	}
}

func mapProtoQuestions(pqs []*game.Question) []entity.Question {
	qs := make([]entity.Question, 0, len(pqs))
	for _, pq := range pqs {
		qs = append(qs, entity.Question{
			Id:              uint(pq.Id),
			Text:            pq.Text,
			PossibleAnswers: mapProtoAnswers(pq.PossibleAnswers),
			CorrectAnswer:   uint(pq.CorrectAnswer),
			Difficulty:      entity.QuestionDifficulty(pq.Difficulty),
			Category:        entity.Category(pq.Category),
		})
	}
	return qs
}

func mapProtoAnswers(pas []*game.PossibleAnswer) []entity.PossibleAnswer {
	answers := make([]entity.PossibleAnswer, 0, len(pas))
	for _, pa := range pas {
		answers = append(answers, entity.PossibleAnswer{
			Id:     uint(pa.Id),
			Text:   pa.Text,
			Choice: entity.PossibleAnswerChoice(pa.PossibleAnswerChoice),
		})
	}
	return answers
}

func mapQuestionsToProto(questions []entity.Question) []*game.Question {
	protoQuestions := make([]*game.Question, 0, len(questions))
	for _, q := range questions {
		protoQuestions = append(protoQuestions, &game.Question{
			Id:              uint64(q.Id),
			Text:            q.Text,
			PossibleAnswers: mapAnswersToProto(q.PossibleAnswers),
			CorrectAnswer:   uint32(q.CorrectAnswer),
			Difficulty:      uint32(q.Difficulty),
			Category:        string(q.Category),
		})
	}
	return protoQuestions
}

func mapAnswersToProto(answers []entity.PossibleAnswer) []*game.PossibleAnswer {
	protoAnswers := make([]*game.PossibleAnswer, 0, len(answers))
	for _, a := range answers {
		protoAnswers = append(protoAnswers, &game.PossibleAnswer{
			Id:                   uint64(a.Id),
			Text:                 a.Text,
			PossibleAnswerChoice: uint32(a.Choice),
		})
	}
	return protoAnswers
}
