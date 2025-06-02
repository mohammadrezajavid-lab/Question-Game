package presenceparam

type GetPresenceRequest struct {
	UserIds []uint
}

func NewGetPresenceRequest(userIds []uint) GetPresenceRequest {
	return GetPresenceRequest{UserIds: userIds}
}

type PresenceItem struct {
	UserId    uint
	Timestamp int64
}

func NewPresenceItem(userId uint, timestamp int64) PresenceItem {
	return PresenceItem{UserId: userId, Timestamp: timestamp}
}

type GetPresenceResponse struct {
	Items []PresenceItem
}

func NewGetPresenceResponse(items []PresenceItem) GetPresenceResponse {
	return GetPresenceResponse{Items: items}
}
