package presenceparam

type UpsertPresenceRequest struct {
	UserId    uint
	TimeStamp int64
}

func NewUpsertPresenceRequest(userId uint, timestamp int64) UpsertPresenceRequest {
	return UpsertPresenceRequest{
		UserId:    userId,
		TimeStamp: timestamp,
	}
}

type UpsertPresenceResponse struct {
	TimeStamp int64
}

func NewUpsertPresenceResponse(timestamp int64) UpsertPresenceResponse {
	return UpsertPresenceResponse{TimeStamp: timestamp}
}
