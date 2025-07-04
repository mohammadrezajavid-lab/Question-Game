package protobufencodedecode

import (
	"encoding/base64"
	"github.com/prometheus/client_golang/prometheus"
	"golang.project/go-fundamentals/gameapp/contract/goprotobuf/matching"
	"golang.project/go-fundamentals/gameapp/entity"
	"golang.project/go-fundamentals/gameapp/logger"
	"golang.project/go-fundamentals/gameapp/metrics"
	"golang.project/go-fundamentals/gameapp/pkg/slice"
	"google.golang.org/protobuf/proto"
)

func EncodeMatchingWaitedUsersEvent(mu entity.MatchedUsers) string {
	const operation = "protobufencodedecode.EncodeMatchingWaitedUsersEvent"

	pbMu := matching.MatchedUsers{
		Category: string(mu.Category),
		UserIds:  slice.MapFromUintToUint64(mu.UserIds),
	}

	payload, mErr := proto.Marshal(&pbMu)
	if mErr != nil {
		metrics.EncodeToProtobufFailedCounter.With(prometheus.Labels{"encoder_name": operation}).Inc()
		logger.Error(mErr, "Proto buf Marshal MatchedUsers Error")

		return ""
	}

	var payloadStr = base64.StdEncoding.EncodeToString(payload)
	return payloadStr
}

func DecodeMatchingWaitedUsersEvent(payload string) entity.MatchedUsers {
	const operation = "protobufencodedecode.DecodeMatchingWaitedUsersEvent"

	payloadByte, dErr := base64.StdEncoding.DecodeString(payload)
	if dErr != nil {
		metrics.DecodeFromProtobufFailedCounter.With(prometheus.Labels{"decoder_name": operation}).Inc()
		logger.Error(dErr, "Decode Match Waited users Event Error")

		return entity.MatchedUsers{}
	}

	pbMu := matching.MatchedUsers{
		Category: "",
		UserIds:  nil,
	}
	if uErr := proto.Unmarshal(payloadByte, &pbMu); uErr != nil {
		metrics.DecodeFromProtobufFailedCounter.With(prometheus.Labels{"decoder_name": operation}).Inc()
		logger.Error(uErr, "Unmarshal payload proto buf to MatchedUsers Error")

		return entity.MatchedUsers{}
	}

	return entity.NewMatchedUsers(entity.Category(pbMu.Category), slice.MapFromUint64ToUint(pbMu.UserIds))
}
