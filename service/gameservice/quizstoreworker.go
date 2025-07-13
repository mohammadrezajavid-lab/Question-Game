package gameservice

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"golang.project/go-fundamentals/gameapp/entity"
	"golang.project/go-fundamentals/gameapp/logger"
	"golang.project/go-fundamentals/gameapp/metrics"
	"golang.project/go-fundamentals/gameapp/pkg/protobufencodedecode"
	"sync"
)

func (s *Service) startQuizStoreWorkers(ctx context.Context) {
	jobs, err := s.subscriber.SubscribeTopic(ctx, entity.GameSvcGameEvent)
	if err != nil {
		metrics.FailedSubscribeTopicCounter.With(prometheus.Labels{"topic": entity.GameSvcGameEvent}).Inc()
		logger.Fatal(err, "subscribe to topic failed", "topic", entity.GameSvcGameEvent)
	}

	var wg sync.WaitGroup
	for i := 0; i < int(s.config.GameEventWorkers); i++ {
		wg.Add(1)
		go s.runQuizStoreWorker(ctx, jobs, &wg, i)
	}
	wg.Wait()
}

func (s *Service) runQuizStoreWorker(ctx context.Context, jobs <-chan string, wg *sync.WaitGroup, workerId int) {
	defer wg.Done()
	logger.Info("Started QuizStoreWorker", "worker_id", workerId)

	for {
		select {
		case <-ctx.Done():
			logger.Info("QuizStoreWorker shutdown initiated", "worker_id", workerId)
			return
		case job, ok := <-jobs:
			if !ok {
				logger.Info("QuizStoreWorker job channel closed", "worker_id", workerId)
				return
			}
			s.fetchAndStoreGameQuiz(ctx, job)
		}
	}
}

func (s *Service) fetchAndStoreGameQuiz(ctx context.Context, payload string) {
	gameEvent := protobufencodedecode.DecodeGameSvcGameEvent(payload)

	questions, err := s.questionRepo.GetQuestions(ctx, gameEvent.QuestionIds)
	if err != nil {
		metrics.FailedGetQuestionsCounter.Inc()
		logger.Warn(err, "Failed to fetch questions from Database", "game_id", gameEvent.GameId, "question_ids", gameEvent.QuestionIds)
		return
	}

	quiz := entity.GameQuiz{
		GameId:    gameEvent.GameId,
		PlayerIds: gameEvent.PlayerIds,
		Questions: questions,
	}
	data := protobufencodedecode.EncodeGameSvcGameQuiz(quiz)
	key := fmt.Sprintf("game_quiz:%d", quiz.GameId)
	if sErr := s.kvStore.Set(ctx, key, data, s.config.GameQuizTTLExpiration); err != nil {
		metrics.FailedStoreQuiz.Inc()
		logger.Error(sErr, "Failed to store quiz in key-value storage", "game_id", quiz.GameId)
	}
}
