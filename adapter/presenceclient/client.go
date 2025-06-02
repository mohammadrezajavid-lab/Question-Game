package presenceclient

import (
	"context"
	"fmt"
	"golang.project/go-fundamentals/gameapp/contract/golang/presence"
	"golang.project/go-fundamentals/gameapp/param/presenceparam"
	"golang.project/go-fundamentals/gameapp/pkg/protobufmapper"
	"golang.project/go-fundamentals/gameapp/pkg/slice"
	"google.golang.org/grpc"
	"log"
)

type Config struct {
	Host    string `mapstructure:"host"`
	Port    int    `mapstructure:"port"`
	Network string `mapstructure:"network"`
}
type Client struct {
	config Config
}

func NewClient(config Config) Client {
	return Client{
		config: config,
	}
}

func (c Client) GetPresence(ctx context.Context, request presenceparam.GetPresenceRequest) (presenceparam.GetPresenceResponse, error) {

	client := c.definitionGrpcClient()

	res, err := client.GetPresence(ctx, &presence.GetPresenceRequest{UserIds: slice.MapFromUintToUint64(request.UserIds)})
	if err != nil {
		log.Println(err)

		return presenceparam.GetPresenceResponse{}, err
	}

	return protobufmapper.MapProtobufToGetPresenceResponse(res), nil
}

func (c Client) Upsert(ctx context.Context, request presenceparam.UpsertPresenceRequest) (presenceparam.UpsertPresenceResponse, error) {

	client := c.definitionGrpcClient()

	res, err := client.Upsert(ctx, &presence.UpsertPresenceRequest{UserId: uint64(request.UserId), Timestamp: request.TimeStamp})
	if err != nil {
		log.Println(err)

		return presenceparam.UpsertPresenceResponse{}, err
	}

	return presenceparam.NewUpsertPresenceResponse(res.Timestamp), nil
}

func (c Client) definitionGrpcClient() presence.PresenceServiceClient {
	target := fmt.Sprintf("%s:%d", c.config.Host, c.config.Port)
	grpcConnectionClient, err := grpc.Dial(target, grpc.WithInsecure())
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer grpcConnectionClient.Close()

	client := presence.NewPresenceServiceClient(grpcConnectionClient)

	return client
}
