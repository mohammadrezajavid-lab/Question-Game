package gameservice

import (
	"context"
	"golang.project/go-fundamentals/gameapp/adapter/redis"
	"golang.project/go-fundamentals/gameapp/contract/broker"
	"golang.project/go-fundamentals/gameapp/entity"
	"golang.project/go-fundamentals/gameapp/logger"
	"golang.project/go-fundamentals/gameapp/param/quizparam"
	"golang.project/go-fundamentals/gameapp/pkg/protobufencodedecode"
	"sync"
)

type Repository interface {
	CreatePlayer(player entity.Player) (entity.Player, error)
	CreateGame(game entity.Game) (entity.Game, error)
}

type QuizClient interface {
	GetQuiz(ctx context.Context, request quizparam.GetQuizRequest) (quizparam.GetQuizResponse, error)
}

type Config struct {
	NumWorkers uint `mapstructure:"num_workers"`
}

type Service struct {
	brokerAdapter *redis.Adapter
	gameRepo      Repository
	quizClient    QuizClient
	publisher     broker.Publisher
	subscriber    broker.Subscriber
	config        Config
}

func New(brokerAdapter *redis.Adapter, gameRepo Repository, quizClient QuizClient, publisher broker.Publisher, subscriber broker.Subscriber, config Config) Service {
	return Service{
		brokerAdapter: brokerAdapter,
		gameRepo:      gameRepo,
		quizClient:    quizClient,
		publisher:     publisher,
		subscriber:    subscriber,
		config:        config,
	}
}

func (s *Service) Start(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	logger.Info("Starting game service dispatcher...")
	s.dispatcher(ctx)
}

func (s *Service) dispatcher(ctx context.Context) {
	jobQueue, err := s.subscriber.SubscribeTopic(ctx, entity.MatchingUsersMatchedEvent)
	if err != nil {
		logger.Fatal(err, "failed to subscribe to broker topic")
	}

	var wg sync.WaitGroup
	for i := 0; i < int(s.config.NumWorkers); i++ {
		wg.Add(1)
		go s.worker(ctx, jobQueue, &wg)
	}

	wg.Wait()
}

func (s *Service) worker(ctx context.Context, jobs <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			logger.Info("Worker in GameSvc received shutdown signal")
			return
		case job, ok := <-jobs:
			if !ok {
				return
			}
			s.handleMatchedUsers(ctx, job)
		}
	}
}

func (s *Service) handleMatchedUsers(ctx context.Context, payload string) {
	select {
	case <-ctx.Done():
		return
	default:
		// decode
		mu := protobufencodedecode.DecodeMatchingWaitedUsersEvent(payload)

		// Create newGame
		newGame, rErr := s.gameRepo.CreateGame(entity.NewGame(mu.Category, mu.Difficulty))
		if rErr != nil {
			logger.Error(rErr, "failed to create game")
			return
		}

		// Create players for game
		playerIds, pErr := s.createPlayers(mu.UserIds, newGame.Id)
		if pErr != nil {
			logger.Error(pErr, "failed to create player")
			return
		}
		newGame.PlayerIds = playerIds

		// TODO - انتخاب بک مجموعه سوال مثلا 15 تایی از استخر سوال ها با توجه به دسته و سطح سختی یا آسانی تعیین شده توسط کاربر
		s.quizClient.GetQuiz(ctx, quizparam.GetQuizRequest{
			Category:   newGame.Category,
			Difficulty: newGame.Difficulty,
		})
		// Published CreatedGameEvent
		cg := entity.NewCreatedGame(newGame.Id, newGame.PlayerIds)
		payloadCg := protobufencodedecode.EncodeGameSvcCreatedGameEvent(cg)
		s.publisher.PublishEvent(entity.GameSvcCreatedGameEvent, payloadCg)
	}
}

func (s *Service) createPlayers(userIds []uint, gameId uint) ([]uint, error) {

	playerIds := make([]uint, 0, 2)

	for _, userId := range userIds {
		player, err := s.gameRepo.CreatePlayer(entity.NewPlayer(userId, gameId))
		if err != nil {
			return nil, err
		}

		playerIds = append(playerIds, player.Id)
	}

	return playerIds, nil
}
