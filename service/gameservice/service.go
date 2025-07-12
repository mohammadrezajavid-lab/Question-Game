package gameservice

import (
	"context"
	"fmt"
	"golang.project/go-fundamentals/gameapp/adapter/redis"
	"golang.project/go-fundamentals/gameapp/contract/broker"
	"golang.project/go-fundamentals/gameapp/entity"
	"golang.project/go-fundamentals/gameapp/logger"
	"golang.project/go-fundamentals/gameapp/param/quizparam"
	"golang.project/go-fundamentals/gameapp/pkg/protobufencodedecode"
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

func (s *Service) startMatchedUsersWorkers(ctx context.Context) {
	jobs, err := s.subscriber.SubscribeTopic(ctx, entity.MatchingUsersMatchedEvent)
	if err != nil {
		logger.Fatal(err, "subscribe to MatchingUsersMatchedEvent failed")
	}

	var wg sync.WaitGroup
	for i := 0; i < int(s.config.MatchedEventWorkers); i++ {
		wg.Add(1)
		go s.processMatchedUsers(ctx, jobs, &wg)
	}
	wg.Wait()
}

func (s *Service) startGameEventWorkers(ctx context.Context) {
	jobs, err := s.subscriber.SubscribeTopic(ctx, entity.GameSvcGameEvent)
	if err != nil {
		logger.Fatal(err, "subscribe to GameSvcGameEvent failed")
	}

	var wg sync.WaitGroup
	for i := 0; i < int(s.config.GameEventWorkers); i++ {
		wg.Add(1)
		go s.processGameEvent(ctx, jobs, &wg)
	}
	wg.Wait()
}

func (s *Service) processMatchedUsers(ctx context.Context, jobs <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			logger.Info("MatchedUsers worker shutdown")
			return
		case job, ok := <-jobs:
			if !ok {
				return
			}
			s.handleMatch(ctx, job)
		}
	}
}

func (s *Service) processGameEvent(ctx context.Context, jobs <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			logger.Info("GameEvent worker shutdown")
			return
		case job, ok := <-jobs:
			if !ok {
				return
			}
			s.storeQuiz(ctx, job)
		}
	}
}

func (s *Service) handleMatch(ctx context.Context, payload string) {
	match := protobufencodedecode.DecodeMatchingWaitedUsersEvent(payload)

	game, err := s.gameRepo.CreateGame(ctx, entity.NewGame(match.Category, match.Difficulty))
	if err != nil {
		logger.Error(err, "failed to create game")
		return
	}

	playerIds, err := s.createPlayers(ctx, match.UserIds, game.Id)
	if err != nil {
		logger.Error(err, "failed to create players")
		return
	}
	game.PlayerIds = playerIds

	quiz, err := s.quizClient.GetQuiz(ctx, quizparam.GetQuizRequest{
		Category:   game.Category,
		Difficulty: game.Difficulty,
	})
	if err != nil {
		logger.Warn(err, "failed to fetch quiz")
		return
	}
	game.QuestionIds = quiz.QuestionIds

	gameEvent := entity.GameEvent{GameId: game.Id, PlayerIds: playerIds, QuestionIds: game.QuestionIds}
	s.publisher.PublishEvent(entity.GameSvcGameEvent, protobufencodedecode.EncodeGameSvcGameEvent(gameEvent))
	s.publisher.PublishEvent(entity.GameSvcCreatedGameEvent, protobufencodedecode.EncodeGameSvcCreatedGameEvent(
		entity.NewCreatedGame(game.Id, playerIds)))
}

func (s *Service) storeQuiz(ctx context.Context, payload string) {
	gameEvent := protobufencodedecode.DecodeGameSvcGameEvent(payload)

	questions, err := s.questionRepo.GetQuestions(ctx, gameEvent.QuestionIds)
	if err != nil {
		logger.Warn(err, "failed to get questions")
		return
	}

	quiz := entity.GameQuiz{
		GameId:    gameEvent.GameId,
		PlayerIds: gameEvent.PlayerIds,
		Questions: questions,
	}

	key := fmt.Sprintf("game_quiz:%d", quiz.GameId)
	data := protobufencodedecode.EncodeGameSvcGameQuiz(quiz)

	if err := s.kvStore.Set(ctx, key, data, s.config.GameQuizTTLExpiration); err != nil {
		logger.Warn(err, fmt.Sprintf("failed to store quiz for game_id: %d", quiz.GameId))
	}
}

func (s *Service) createPlayers(ctx context.Context, userIds []uint, gameId uint) ([]uint, error) {
	playerIds := make([]uint, 0, len(userIds))
	for _, uid := range userIds {
		player, err := s.gameRepo.CreatePlayer(ctx, entity.NewPlayer(uid, gameId))
		if err != nil {
			return nil, err
		}
		playerIds = append(playerIds, player.Id)
	}
	return playerIds, nil
}
