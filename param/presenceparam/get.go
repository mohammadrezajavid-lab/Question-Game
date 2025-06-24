package presenceparam

import "sort"

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

func (gp *GetPresenceResponse) SortItemsByUserId() {
	sort.Slice(gp.Items, func(i, j int) bool {
		return gp.Items[i].UserId < gp.Items[j].UserId
	})
}

func (gp *GetPresenceResponse) FindByUserId(targetUserId uint) *PresenceItem {

	index := sort.Search(len(gp.Items), func(i int) bool {
		return gp.Items[i].UserId >= targetUserId
	})

	if index < len(gp.Items) && gp.Items[index].UserId == targetUserId {
		return &gp.Items[index]
	}

	return nil
}
