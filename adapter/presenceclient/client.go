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
	conn   *grpc.ClientConn
	client presence.PresenceServiceClient
}

func NewClient(config Config) (Client, error) {
	target := fmt.Sprintf("%s:%d", config.Host, config.Port)
	// Establish the connection once.
	grpcConnection, err := grpc.Dial(target, grpc.WithInsecure()) // Note: grpc.WithInsecure() is deprecated. Use credentials.NewTLS(&tls.Config{InsecureSkipVerify: true}) for testing if needed.
	if err != nil {
		metrics.FailedOpenPresenceClientGRPCConnCounter.Inc()
		logger.Warn(err, "failed_open_grpc_connection")
		return Client{}, err
	}

	client := presence.NewPresenceServiceClient(grpcConnection)

	return Client{
		config: config,
		conn:   grpcConnection,
		client: client,
	}, nil
}

// Close Add a method to close the connection during shutdown
func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *Client) Upsert(ctx context.Context, request presenceparam.UpsertPresenceRequest) (presenceparam.UpsertPresenceResponse, error) {
	res, err := c.client.Upsert(ctx, &presence.UpsertPresenceRequest{UserId: uint64(request.UserId), Timestamp: request.TimeStamp})
	if err != nil {
		return presenceparam.UpsertPresenceResponse{}, err
	}
	return presenceparam.NewUpsertPresenceResponse(res.Timestamp), nil
}

func (c *Client) GetPresence(ctx context.Context, request presenceparam.GetPresenceRequest) (presenceparam.GetPresenceResponse, error) {
	res, err := c.client.GetPresence(ctx, &presence.GetPresenceRequest{UserIds: slice.MapFromUintToUint64(request.UserIds)})
	if err != nil {
		return presenceparam.GetPresenceResponse{}, err
	}
	return protobufmapper.MapProtobufToGetPresenceResponse(res), err
}
