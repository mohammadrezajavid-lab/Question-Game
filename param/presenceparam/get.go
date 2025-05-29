package presenceparam

type GetPresenceRequest struct {
	UserIds []uint
}

func NewGetPresenceRequest(userIds []uint) GetPresenceRequest {
	return GetPresenceRequest{UserIds: userIds}
}

type GetPresenceResponse struct {
	Items []PresenceItem
}

type PresenceItem struct {
	UserId    uint
	Timestamp int64
}
