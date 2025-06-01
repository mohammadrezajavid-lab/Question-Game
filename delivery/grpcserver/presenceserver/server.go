package presenceserver

import (
	"context"
	"fmt"
	"golang.project/go-fundamentals/gameapp/config/grpcconfig"
	"golang.project/go-fundamentals/gameapp/contract/golang/presence"
	"golang.project/go-fundamentals/gameapp/param/presenceparam"
	"golang.project/go-fundamentals/gameapp/pkg/protobufmapper"
	"golang.project/go-fundamentals/gameapp/pkg/slice"
	"golang.project/go-fundamentals/gameapp/service/presenceservice"
	"google.golang.org/grpc"
	"log"
	"net"
)

type PresenceGrpcServer struct {
	presence.UnimplementedPresenceServiceServer
	presenceSvc *presenceservice.Service
	grpcCfg     *grpcconfig.Config
}

func NewPresenceGrpcServer(presenceSvc *presenceservice.Service, grpcCfg *grpcconfig.Config) *PresenceGrpcServer {
	return &PresenceGrpcServer{
		UnimplementedPresenceServiceServer: presence.UnimplementedPresenceServiceServer{},
		presenceSvc:                        presenceSvc,
		grpcCfg:                            grpcCfg,
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
		return nil, err
	}

	return protobufmapper.MapGetPresenceResponseToProtobuf(res), nil
}

func (s *PresenceGrpcServer) Start() {
	addr := fmt.Sprintf("%s:%d", s.grpcCfg.GrpcCfg.Host, s.grpcCfg.GrpcCfg.Port)
	// tcp port
	listener, lErr := net.Listen(s.grpcCfg.GrpcCfg.Network, addr)
	if lErr != nil {
		panic(lErr)
	}

	// grpc-presence server
	presenceSvcSrv := NewPresenceGrpcServer(s.presenceSvc, s.grpcCfg)

	// grpc server
	grpcSrv := grpc.NewServer()

	// presence service register to grpc server
	presence.RegisterPresenceServiceServer(grpcSrv, presenceSvcSrv)

	// server grpcServer by lisin
	log.Printf("presence grpc server started on %s/n", addr)
	if sErr := grpcSrv.Serve(listener); sErr != nil {
		log.Fatal("couldn't server presence grpc server")
	}

}
