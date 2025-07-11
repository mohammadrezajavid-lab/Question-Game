package presenceserver

import (
	"context"
	"fmt"
	"golang.project/go-fundamentals/gameapp/contract/goprotobuf/presence"
	"golang.project/go-fundamentals/gameapp/logger"
	"golang.project/go-fundamentals/gameapp/param/presenceparam"
	"golang.project/go-fundamentals/gameapp/pkg/protobufmapper"
	"golang.project/go-fundamentals/gameapp/pkg/richerror"
	"golang.project/go-fundamentals/gameapp/pkg/slice"
	"golang.project/go-fundamentals/gameapp/service/presenceservice"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"
)

type Config struct {
	Host    string `mapstructure:"host"`
	Port    int    `mapstructure:"port"`
	Network string `mapstructure:"network"`
}

type PresenceGrpcServer struct {
	grpcCfg *Config
	presence.UnimplementedPresenceServiceServer
	presenceSvc *presenceservice.Service
	grpcServer  *grpc.Server
}

func NewPresenceGrpcServer(presenceSvc *presenceservice.Service, grpcCfg *Config) *PresenceGrpcServer {
	return &PresenceGrpcServer{
		grpcCfg:                            grpcCfg,
		UnimplementedPresenceServiceServer: presence.UnimplementedPresenceServiceServer{},
		presenceSvc:                        presenceSvc,
		grpcServer:                         nil,
	}
}

func (s *PresenceGrpcServer) Upsert(ctx context.Context, req *presence.UpsertPresenceRequest) (*presence.UpsertPresenceResponse, error) {
	res, err := s.presenceSvc.Upsert(ctx, presenceparam.NewUpsertPresenceRequest(uint(req.UserId), req.Timestamp))

	if err != nil {
		return nil, err
	}

	return &presence.UpsertPresenceResponse{Timestamp: res.TimeStamp}, nil
}

func (s *PresenceGrpcServer) GetPresence(ctx context.Context, req *presence.GetPresenceRequest) (*presence.GetPresenceResponse, error) {

	res, err := s.presenceSvc.GetPresence(ctx, presenceparam.NewGetPresenceRequest(slice.MapFromUint64ToUint(req.GetUserIds())))

	if err != nil {

		richErr, ok := err.(*richerror.RichError)
		if ok {
			if richErr.GetKind() == richerror.KindNotFound {
				return nil, status.Errorf(codes.NotFound, richErr.Error())
			}

			return nil, status.Errorf(codes.Internal, "An internal error occurred: %s", richErr.Error())
		}

		return nil, status.Errorf(codes.Internal, "An unexpected error occurred: %s", err.Error())
	}

	return protobufmapper.MapGetPresenceResponseToProtobuf(res), nil
}

func (s *PresenceGrpcServer) Start() {
	addr := fmt.Sprintf("%s:%d", s.grpcCfg.Host, s.grpcCfg.Port)

	listener, lErr := net.Listen(s.grpcCfg.Network, addr)
	if lErr != nil {
		panic(lErr)
	}

	// grpc server
	grpcSrv := grpc.NewServer()
	s.grpcServer = grpcSrv

	// presence service register to grpc server
	presence.RegisterPresenceServiceServer(grpcSrv, s)

	// server grpcServer by listen
	logger.Info(fmt.Sprintf("presence grpc server started on %s", addr))
	if sErr := grpcSrv.Serve(listener); sErr != nil {
		logger.Fatal(sErr, "couldn't serve presence grpc server")
	}
}

func (s *PresenceGrpcServer) Shutdown() {
	if s != nil && s.grpcServer != nil {
		logger.Info("presence grpc server shutting down gracefully")
		s.grpcServer.Stop()
	}
}
