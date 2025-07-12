package gameservice

import (
	"context"
	"golang.project/go-fundamentals/gameapp/entity"
	"golang.project/go-fundamentals/gameapp/logger"
	"golang.project/go-fundamentals/gameapp/param/quizparam"
	"golang.project/go-fundamentals/gameapp/pkg/protobufencodedecode"
	"sync"
)

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
