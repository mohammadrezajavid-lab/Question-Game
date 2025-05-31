package presenceserver

import (
	"context"
	"fmt"
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
	presenceSvc presenceservice.Service
}

func NewPresenceGrpcServer(presenceSvc presenceservice.Service) PresenceGrpcServer {
	return PresenceGrpcServer{
		UnimplementedPresenceServiceServer: presence.UnimplementedPresenceServiceServer{},
		presenceSvc:                        presenceSvc,
	}
}

func (s PresenceGrpcServer) GetPresence(ctx context.Context, req *presence.GetPresenceRequest) (*presence.GetPresenceResponse, error) {

	res, err := s.presenceSvc.GetPresence(ctx, presenceparam.NewGetPresenceRequest(slice.MapFromUint64ToUint(req.GetUserIds())))

	if err != nil {
		return nil, err
	}

	return protobufmapper.MapGetPresenceResponseToProtobuf(res), nil
}

func (s PresenceGrpcServer) Start() {
	addr := fmt.Sprintf(":%d", 8086)
	// tcp port
	listener, lErr := net.Listen("tcp", addr)
	if lErr != nil {
		panic(lErr)
	}

	// grpc-presence server
	presenceSvcSrv := PresenceGrpcServer{}

	// grpc server
	grpcSrv := grpc.NewServer()

	// presence service register to grpc server
	presence.RegisterPresenceServiceServer(grpcSrv, &presenceSvcSrv)

	// server grpcServer by lisin
	log.Printf("presence grpc server started on %s/n", addr)
	if sErr := grpcSrv.Serve(listener); sErr != nil {
		log.Fatal("couldn't server presence grpc server")
	}

}
