package presenceclient

import (
	"context"
	"fmt"
	"golang.project/go-fundamentals/gameapp/contract/goprotobuf/presence"
	"golang.project/go-fundamentals/gameapp/logger"
	"golang.project/go-fundamentals/gameapp/metrics"
	"golang.project/go-fundamentals/gameapp/param/presenceparam"
	"golang.project/go-fundamentals/gameapp/pkg/protobufmapper"
	"golang.project/go-fundamentals/gameapp/pkg/slice"
	"google.golang.org/grpc"
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

	// TODO - what's the best practice for reliable communication - Retry for connection time out?!
	// TODO - is it okay to create new connection for every method call?
	client, grpcClientConn := c.definitionGrpcClient()
	defer grpcClientConn.Close()

	res, err := client.GetPresence(ctx, &presence.GetPresenceRequest{UserIds: slice.MapFromUintToUint64(request.UserIds)})
	if err != nil {
		return presenceparam.GetPresenceResponse{}, err
	}

	return protobufmapper.MapProtobufToGetPresenceResponse(res), nil
}

func (c Client) Upsert(ctx context.Context, request presenceparam.UpsertPresenceRequest) (presenceparam.UpsertPresenceResponse, error) {

	// TODO - what's the best practice for reliable communication - Retry for connection time out?!
	// TODO - is it okay to create new connection for every method call?
	client, grpcClientConn := c.definitionGrpcClient()
	defer grpcClientConn.Close()

	res, err := client.Upsert(ctx, &presence.UpsertPresenceRequest{UserId: uint64(request.UserId), Timestamp: request.TimeStamp})
	if err != nil {
		return presenceparam.UpsertPresenceResponse{}, err
	}

	return presenceparam.NewUpsertPresenceResponse(res.Timestamp), nil
}

func (c Client) definitionGrpcClient() (presence.PresenceServiceClient, *grpc.ClientConn) {
	target := fmt.Sprintf("%s:%d", c.config.Host, c.config.Port)
	grpcConnection, err := grpc.Dial(target, grpc.WithInsecure())
	if err != nil {
		metrics.FailedOpenPresenceClientGRPCConnCounter.Inc()
		logger.Warn(err, "failed_open_grpc_connection")
	}

	client := presence.NewPresenceServiceClient(grpcConnection)

	return client, grpcConnection
}
