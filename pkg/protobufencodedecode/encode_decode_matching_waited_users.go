package protobufencodedecode

import (
	"encoding/base64"
	"golang.project/go-fundamentals/gameapp/contract/golang/matching"
	"golang.project/go-fundamentals/gameapp/entity"
	"golang.project/go-fundamentals/gameapp/pkg/slice"
	"google.golang.org/protobuf/proto"
)

func EncodeMatchingWaitedUsersEvent(mu entity.MatchedUsers) string {
	pbMu := matching.MatchedUsers{
		Category: string(mu.Category),
		UserIds:  slice.MapFromUintToUint64(mu.UserIds),
	}
	payload, mErr := proto.Marshal(&pbMu)
	if mErr != nil {
		// TODO - update metrics
		// TODO - log error
	}

	var payloadStr = base64.StdEncoding.EncodeToString(payload)
	return payloadStr
}

func DecodeMatchingWaitedUsersEvent(payload string) entity.MatchedUsers {
	payloadByte, dErr := base64.StdEncoding.DecodeString(payload)
	if dErr != nil {
		// TODO - update metrics
		// TODO - log error
	}

	pbMu := matching.MatchedUsers{
		Category: "",
		UserIds:  nil,
	}
	if uErr := proto.Unmarshal(payloadByte, &pbMu); uErr != nil {
		// TODO - update metrics
		// TODO - log error
	}

	return entity.NewMatchedUsers(entity.Category(pbMu.Category), slice.MapFromUint64ToUint(pbMu.UserIds))
}
