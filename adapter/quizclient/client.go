package quizclient

import (
	"context"
	"fmt"
	"golang.project/go-fundamentals/gameapp/contract/goprotobuf/quiz"
	"golang.project/go-fundamentals/gameapp/logger"
	"golang.project/go-fundamentals/gameapp/metrics"
	"golang.project/go-fundamentals/gameapp/param/quizparam"
	"golang.project/go-fundamentals/gameapp/pkg/slice"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	Host    string `mapstructure:"host"`
	Port    int    `mapstructure:"port"`
	Network string `mapstructure:"network"`
}
type Client struct {
	config Config
	conn   *grpc.ClientConn
	client quiz.QuizServiceClient
}

func NewClient(config Config) (Client, error) {
	target := fmt.Sprintf("%s:%d", config.Host, config.Port)

	retryPolicy := `{
	  "methodConfig": [{
	    "name": [{"service": "quiz.QuizService"}],
	    "retryPolicy": {
	      "MaxAttempts": 4,
	      "InitialBackoff": "0.1s",
	      "MaxBackoff": "1s",
	      "BackoffMultiplier": 2.0,
	      "RetryableStatusCodes": [ "UNAVAILABLE" ]
	    }
	  }]
	}`

	conn, err := grpc.Dial(
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(retryPolicy),
	)

	if err != nil {
		metrics.FailedOpenQuizClientGRPCConnCounter.Inc()
		logger.Error(err, "failed to dial  Quiz gRPC server")
		return Client{}, err
	}

	client := quiz.NewQuizServiceClient(conn)

	return Client{
		config: config,
		conn:   conn,
		client: client,
	}, nil
}

func (c *Client) GetQuiz(ctx context.Context, request quizparam.GetQuizRequest) (quizparam.GetQuizResponse, error) {
	res, err := c.client.GetQuiz(ctx, &quiz.GetQuizRequest{
		Category:   string(request.Category),
		Difficulty: uint32(request.Difficulty),
	})
	if err != nil {
		return quizparam.GetQuizResponse{}, err
	}
	return quizparam.GetQuizResponse{QuestionIds: slice.MapFromUint64ToUint(res.QuestionIds)}, nil
}

func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}
