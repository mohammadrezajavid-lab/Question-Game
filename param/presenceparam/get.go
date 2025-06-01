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

func NewPresenceItem() PresenceItem {
	return PresenceItem{}
}

type GetPresenceResponse struct {
	Items []PresenceItem
}

func NewGetPresenceResponse(items []PresenceItem) GetPresenceResponse {
	return GetPresenceResponse{Items: items}
}
