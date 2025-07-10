package quizserver

import (
	"context"
	"fmt"
	"golang.project/go-fundamentals/gameapp/contract/goprotobuf/quiz"
	"golang.project/go-fundamentals/gameapp/entity"
	"golang.project/go-fundamentals/gameapp/logger"
	"golang.project/go-fundamentals/gameapp/param/quizparam"
	"golang.project/go-fundamentals/gameapp/pkg/slice"
	"golang.project/go-fundamentals/gameapp/service/quizservice"
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

type QuizGrpcServer struct {
	grpcCfg *Config
	quiz.UnimplementedQuizServiceServer
	quizSvc    *quizservice.Service
	grpcServer *grpc.Server
}

func NewQuizGrpcServer(quizSvc *quizservice.Service, grpcCfg *Config) *QuizGrpcServer {
	return &QuizGrpcServer{
		grpcCfg:                        grpcCfg,
		UnimplementedQuizServiceServer: quiz.UnimplementedQuizServiceServer{},
		quizSvc:                        quizSvc,
		grpcServer:                     nil,
	}
}

func (q *QuizGrpcServer) GetQuiz(ctx context.Context, req *quiz.GetQuizRequest) (*quiz.GetQuizResponse, error) {
	res, err := q.quizSvc.GetQuiz(ctx, quizparam.GetQuizRequest{
		Category:   entity.Category(req.Category),
		Difficulty: entity.QuestionDifficulty(req.Difficulty),
	})
	if err != nil {
		logger.Warn(err, err.Error())
		return nil, status.Errorf(codes.Internal, "An unexpected error occured: %s", err.Error())
	}

	return &quiz.GetQuizResponse{QuestionIds: slice.MapFromUintToUint64(res.QuestionIds)}, nil
}

func (q *QuizGrpcServer) Start() {
	addr := fmt.Sprintf("%s:%d", q.grpcCfg.Host, q.grpcCfg.Port)

	listener, lErr := net.Listen(q.grpcCfg.Network, addr)
	if lErr != nil {
		logger.Panic(lErr, "initial listener for quiz grpc server Failed")
	}

	grpcSrv := grpc.NewServer()
	q.grpcServer = grpcSrv

	quiz.RegisterQuizServiceServer(grpcSrv, q)

	logger.Info(fmt.Sprintf("quiz grpc server started on %s", addr))
	if sErr := grpcSrv.Serve(listener); sErr != nil {
		logger.Fatal(sErr, "couldn't serve quiz grpc server")
	}
}

func (q *QuizGrpcServer) Shutdown() {
	if q != nil && q.grpcServer != nil {
		logger.Info("quiz grpc server shutting down gracefully")
		q.grpcServer.Stop()
	}
}
