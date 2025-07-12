package gameservice

import (
	"context"
	"fmt"
	"golang.project/go-fundamentals/gameapp/entity"
	"golang.project/go-fundamentals/gameapp/logger"
	"golang.project/go-fundamentals/gameapp/pkg/protobufencodedecode"
	"sync"
)

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
