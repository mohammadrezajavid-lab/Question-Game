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

func EncodeGameSvcCreatedGameEvent(cg entity.CreatedGame) string {
	const operation = "protobufencodedecode.EncodeGameSvcCreatedGameEvent"

	protoBufCg := game.CreatedGame{
		GameId:    uint64(cg.GameId),
		PlayerIds: slice.MapFromUintToUint64(cg.PlayerIds),
	}

	payload, mErr := proto.Marshal(&protoBufCg)
	if mErr != nil {
		metrics.EncodeToProtobufFailedCounter.With(prometheus.Labels{"encoder_name": operation}).Inc()
		logger.Error(mErr, "Proto buf Marshal CreatedGame Error")

		return ""
	}

	var payloadStr = base64.StdEncoding.EncodeToString(payload)
	return payloadStr
}

func DecodeGameSvcCreatedGameEvent(payload string) entity.CreatedGame {
	const operation = "protobufencodedecode.DecodeGameSvcCreatedGameEvent"

	payloadByte, dErr := base64.StdEncoding.DecodeString(payload)
	if dErr != nil {
		metrics.DecodeFromProtobufFailedCounter.With(prometheus.Labels{"decoder_name": operation}).Inc()
		logger.Error(dErr, "Decode Created Game Event Error")

		return entity.CreatedGame{}
	}

	protoBufCg := game.CreatedGame{
		GameId:    0,
		PlayerIds: nil,
	}

	if uErr := proto.Unmarshal(payloadByte, &protoBufCg); uErr != nil {
		metrics.DecodeFromProtobufFailedCounter.With(prometheus.Labels{"decoder_name": operation}).Inc()
		logger.Error(uErr, "Unmarshal payload proto buf to CreatedGame Error")

		return entity.CreatedGame{}
	}

	return entity.NewCreatedGame(uint(protoBufCg.GameId), slice.MapFromUint64ToUint(protoBufCg.PlayerIds))
}
