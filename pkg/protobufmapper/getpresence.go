package protobufmapper

import (
	"golang.project/go-fundamentals/gameapp/contract/golang/presence"
	"golang.project/go-fundamentals/gameapp/param/presenceparam"
)

func MapGetPresenceResponseToProtobuf(g presenceparam.GetPresenceResponse) *presence.GetPresenceResponse {
	r := new(presence.GetPresenceResponse)

	for _, item := range g.Items {
		r.Items = append(r.Items, &presence.Presence{
			UserId:    uint64(item.UserId),
			Timestamp: item.Timestamp,
		})
	}

	return r
}

func MapProtobufToGetPresenceResponse(g *presence.GetPresenceResponse) presenceparam.GetPresenceResponse {
	r := presenceparam.GetPresenceResponse{}

	for _, item := range g.GetItems() {
		r.Items = append(r.Items, presenceparam.PresenceItem{
			UserId:    uint(item.UserId),
			Timestamp: item.Timestamp,
		})

	}

	return r
}
