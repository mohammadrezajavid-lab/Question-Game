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

type GetPresenceResponse struct {
	Items []PresenceItem
}
