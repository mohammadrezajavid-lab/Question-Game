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

func EncodeGameSvcGameEvent(event entity.GameEvent) string {
	const operation = "protobufencodedecode.EncodeGameSvcGameEvent"

	protoBufGameEvent := game.GameEvent{
		GameId:      uint64(event.GameId),
		PlayerIds:   slice.MapFromUintToUint64(event.PlayerIds),
		QuestionIds: slice.MapFromUintToUint64(event.QuestionIds),
	}

	payload, mErr := proto.Marshal(&protoBufGameEvent)
	if mErr != nil {
		metrics.EncodeToProtobufFailedCounter.With(prometheus.Labels{"encoder_name": operation}).Inc()
		logger.Error(mErr, "Proto buf Marshal GameEvent Error")

		return ""
	}

	var payloadStr = base64.StdEncoding.EncodeToString(payload)
	return payloadStr
}

func DecodeGameSvcGameEvent(payload string) entity.GameEvent {
	const operation = "protobufencodedecode.DecodeGameSvcGameEvent"

	payloadByte, dErr := base64.StdEncoding.DecodeString(payload)
	if dErr != nil {
		metrics.DecodeFromProtobufFailedCounter.With(prometheus.Labels{"decoder_name": operation}).Inc()
		logger.Error(dErr, "Decode Created Game Event Error")

		return entity.GameEvent{}
	}

	protoBufGameEvent := game.GameEvent{
		GameId:      0,
		PlayerIds:   nil,
		QuestionIds: nil,
	}

	if uErr := proto.Unmarshal(payloadByte, &protoBufGameEvent); uErr != nil {
		metrics.DecodeFromProtobufFailedCounter.With(prometheus.Labels{"decoder_name": operation}).Inc()
		logger.Error(uErr, "Unmarshal payload proto buf to GameEvent Error")

		return entity.GameEvent{}
	}

	return entity.GameEvent{
		QuestionIds: slice.MapFromUint64ToUint(protoBufGameEvent.QuestionIds),
		PlayerIds:   slice.MapFromUint64ToUint(protoBufGameEvent.PlayerIds),
		GameId:      uint(protoBufGameEvent.GameId),
	}
}
