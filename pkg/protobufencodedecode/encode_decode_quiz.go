package protobufencodedecode

import (
	"encoding/base64"
	"github.com/prometheus/client_golang/prometheus"
	"golang.project/go-fundamentals/gameapp/contract/goprotobuf/quiz"
	"golang.project/go-fundamentals/gameapp/entity"
	"golang.project/go-fundamentals/gameapp/logger"
	"golang.project/go-fundamentals/gameapp/metrics"
	"golang.project/go-fundamentals/gameapp/pkg/slice"
	"google.golang.org/protobuf/proto"
)

func EncodeQuizSvcQuiz(q entity.Quiz) string {
	const operation = "protobufencodedecode.EncodeQuizSvcQuiz"

	protoBufQuiz := quiz.GenerateQuiz{QuestionIds: slice.MapFromUintToUint64(q.QuestionIDs)}

	payload, err := proto.Marshal(&protoBufQuiz)
	if err != nil {
		metrics.EncodeToProtobufFailedCounter.With(prometheus.Labels{"encoder_name": operation}).Inc()
		logger.Error(err, "Proto buf Marshal Quiz Error")

		return ""
	}

	return base64.StdEncoding.EncodeToString(payload)
}

func DecodeQuizSvcQuiz(payload string) entity.Quiz {
	const operation = "protobufencodedecode.DecodeQuizSvcQuiz"

	payloadByte, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		metrics.DecodeFromProtobufFailedCounter.With(prometheus.Labels{"decoder_name": operation}).Inc()
		logger.Error(err, "Decode Quiz payload Error")

		return entity.Quiz{}
	}

	protoBufQuiz := quiz.GenerateQuiz{QuestionIds: nil}
	if uErr := proto.Unmarshal(payloadByte, &protoBufQuiz); uErr != nil {
		metrics.DecodeFromProtobufFailedCounter.With(prometheus.Labels{"decoder_name": operation}).Inc()
		logger.Error(uErr, "Unmarshal payload proto buf to Quiz Error")

		return entity.Quiz{}
	}

	return entity.Quiz{QuestionIDs: slice.MapFromUint64ToUint(protoBufQuiz.QuestionIds)}
}
