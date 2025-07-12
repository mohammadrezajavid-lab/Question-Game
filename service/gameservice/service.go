package gameservice

import (
	"context"
	"golang.project/go-fundamentals/gameapp/adapter/redis"
	"golang.project/go-fundamentals/gameapp/contract/broker"
	"golang.project/go-fundamentals/gameapp/entity"
	"golang.project/go-fundamentals/gameapp/logger"
	"golang.project/go-fundamentals/gameapp/param/quizparam"
	"sync"
	"time"
)

type GameRepository interface {
	CreatePlayer(ctx context.Context, player entity.Player) (entity.Player, error)
	CreateGame(ctx context.Context, game entity.Game) (entity.Game, error)
}

type QuestionRepository interface {
	GetQuestions(ctx context.Context, questionIds []uint) ([]entity.Question, error)
}

type QuizClient interface {
	GetQuiz(ctx context.Context, request quizparam.GetQuizRequest) (quizparam.GetQuizResponse, error)
}

type KVStore interface {
	Set(ctx context.Context, key string, value string, ttlExpiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
}

type Config struct {
	GameEventWorkers      uint          `mapstructure:"num_workers_game_event"`
	MatchedEventWorkers   uint          `mapstructure:"num_workers_matched_event"`
	GameQuizTTLExpiration time.Duration `mapstructure:"game_quiz_ttl_expiration"`
}

type Service struct {
	brokerAdapter *redis.Adapter
	gameRepo      GameRepository
	questionRepo  QuestionRepository
	kvStore       KVStore
	quizClient    QuizClient
	publisher     broker.Publisher
	subscriber    broker.Subscriber
	config        Config
}

func New(
	brokerAdapter *redis.Adapter,
	gameRepo GameRepository,
	questionRepo QuestionRepository,
	kvStore KVStore,
	quizClient QuizClient,
	publisher broker.Publisher,
	subscriber broker.Subscriber,
	config Config,
) Service {
	return Service{
		brokerAdapter: brokerAdapter,
		gameRepo:      gameRepo,
		questionRepo:  questionRepo,
		kvStore:       kvStore,
		quizClient:    quizClient,
		publisher:     publisher,
		subscriber:    subscriber,
		config:        config,
	}
}

func (s *Service) Start(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	logger.Info("Starting game service dispatcher...")
	go s.startMatchedUsersWorkers(ctx)
	s.startGameEventWorkers(ctx)
}
