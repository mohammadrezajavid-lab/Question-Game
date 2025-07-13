package gameservice

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"golang.project/go-fundamentals/gameapp/entity"
	"golang.project/go-fundamentals/gameapp/logger"
	"golang.project/go-fundamentals/gameapp/metrics"
	"golang.project/go-fundamentals/gameapp/param/quizparam"
	"golang.project/go-fundamentals/gameapp/pkg/protobufencodedecode"
	"sync"
)

func (s *Service) startMatchEventWorkers(ctx context.Context) {
	jobs, err := s.subscriber.SubscribeTopic(ctx, entity.MatchingUsersMatchedEvent)
	if err != nil {
		metrics.FailedSubscribeTopicCounter.With(prometheus.Labels{"topic": entity.MatchingUsersMatchedEvent}).Inc()
		logger.Fatal(err, "subscribe to topic failed", "topic", entity.MatchingUsersMatchedEvent)
	}

	var wg sync.WaitGroup
	for i := 0; i < int(s.config.MatchedEventWorkers); i++ {
		wg.Add(1)
		go s.runMatchWorker(ctx, jobs, &wg, i)
	}
	wg.Wait()
}

func (s *Service) runMatchWorker(ctx context.Context, jobs <-chan string, wg *sync.WaitGroup, workerId int) {
	defer wg.Done()
	logger.Info("Started MatchWorker", "worker_id", workerId)

	for {
		select {
		case <-ctx.Done():
			logger.Info("MatchWorker shutdown initiated", "worker_id", workerId)
			return
		case job, ok := <-jobs:
			if !ok {
				logger.Info("MatchWorker job channel closed", "worker_id", workerId)
				return
			}
			s.handleMatchEvent(ctx, job)
		}
	}
}

func (s *Service) handleMatchEvent(ctx context.Context, payload string) {
	match := protobufencodedecode.DecodeMatchingWaitedUsersEvent(payload)

	logger.Info("Received match event", "category", match.Category, "difficulty", match.Difficulty, "user_ids", match.UserIds)

	game, err := s.gameRepo.CreateGame(ctx, entity.NewGame(match.Category, match.Difficulty))
	if err != nil {
		metrics.FailedCreateGame.Inc()
		logger.Error(err, "Failed to create game", "category", match.Category, "difficulty", match.Difficulty, "user_ids", match.UserIds)
		return
	}
	logger.Info("Game record created", "game_id", game.Id)

	playerIds, err := s.createGamePlayers(ctx, match.UserIds, game.Id)
	if err != nil {
		metrics.FailedCreatePlayers.Inc()
		logger.Error(err, "Failed to create players", "game_id", game.Id, "user_ids", match.UserIds)
		return
	}
	game.PlayerIds = playerIds
	logger.Info("Players created and linked to game", "game_id", game.Id, "player_ids", playerIds)

	quiz, err := s.quizClient.GetQuiz(ctx, quizparam.GetQuizRequest{
		Category:   game.Category,
		Difficulty: game.Difficulty,
	})
	if err != nil {
		metrics.FailedFetchQuizFromPool.Inc()
		logger.Warn(err, "Failed to fetch quiz from pool", "game_id", game.Id, "category", game.Category, "difficulty", game.Difficulty)
		return
	}
	game.QuestionIds = quiz.QuestionIds

	gameEvent := entity.GameEvent{GameId: game.Id, PlayerIds: playerIds, QuestionIds: game.QuestionIds}
	s.publisher.PublishEvent(entity.GameSvcGameEvent, protobufencodedecode.EncodeGameSvcGameEvent(gameEvent))
	s.publisher.PublishEvent(entity.GameSvcCreatedGameEvent, protobufencodedecode.EncodeGameSvcCreatedGameEvent(
		entity.NewCreatedGame(game.Id, playerIds)))

	metrics.SuccesGameCreated.Inc()
	logger.Info("Game events published", "game_id", game.Id)
}

func (s *Service) createGamePlayers(ctx context.Context, userIds []uint, gameId uint) ([]uint, error) {
	playerIds := make([]uint, 0, len(userIds))
	for _, uid := range userIds {
		player, err := s.gameRepo.CreatePlayer(ctx, entity.NewPlayer(uid, gameId))
		if err != nil {
			logger.Error(err, "Failed to create player", "user_id", uid, "game_id", gameId)
			return nil, err
		}
		playerIds = append(playerIds, player.Id)
	}
	return playerIds, nil
}
