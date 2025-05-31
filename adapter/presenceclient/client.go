package presenceclient

import (
	"context"
	"golang.project/go-fundamentals/gameapp/contract/golang/presence"
	"golang.project/go-fundamentals/gameapp/param/presenceparam"
	"golang.project/go-fundamentals/gameapp/pkg/protobufmapper"
	"golang.project/go-fundamentals/gameapp/pkg/slice"
	"google.golang.org/grpc"
	"log"
)

type Client struct {
	client presence.PresenceServiceClient
}

func NewClient(cc *grpc.ClientConn) Client {
	return Client{
		client: presence.NewPresenceServiceClient(cc),
	}
}

func (c Client) GetPresence(ctx context.Context, request presenceparam.GetPresenceRequest) (presenceparam.GetPresenceResponse, error) {
	res, err := c.client.GetPresence(ctx, &presence.GetPresenceRequest{UserIds: slice.MapFromUintToUint64(request.UserIds)})
	if err != nil {
		log.Println(err)

		return presenceparam.GetPresenceResponse{}, err
	}

	return protobufmapper.MapProtobufToGetPresenceResponse(res), nil
}
