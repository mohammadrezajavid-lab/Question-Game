package quizservice

import (
	"context"
	"fmt"
	"golang.project/go-fundamentals/gameapp/entity"
	"golang.project/go-fundamentals/gameapp/logger"
	"golang.project/go-fundamentals/gameapp/metrics"
	"golang.project/go-fundamentals/gameapp/param/quizparam"
	"golang.project/go-fundamentals/gameapp/pkg/protobufencodedecode"
	"golang.project/go-fundamentals/gameapp/pkg/richerror"
	"sync"
	"time"
)

type SetRepository interface {
	SetLength(ctx context.Context, key string) (int, error)
	SetAdd(ctx context.Context, key string, value string)
	SetPop(ctx context.Context, key string) (string, error)
}

type DBRepository interface {
	GetRandomQuestions(ctx context.Context, category entity.Category, difficulty entity.QuestionDifficulty, numberOfQuestion uint) ([]uint, error)
}

type Config struct {
	LocalContextTimeOut time.Duration `mapstructure:"local_context_timeout"`
	NumberOfQuestions   uint          `mapstructure:"number_of_questions"`
	Prefix              string        `mapstructure:"prefix"`
}

type Service struct {
	Config  Config
	setRepo SetRepository
	dbRepo  DBRepository
}

func New(config Config, setRepository SetRepository, dbRepository DBRepository) Service {
	return Service{
		Config:  config,
		setRepo: setRepository,
		dbRepo:  dbRepository,
	}
}

func (s *Service) GenerateQuiz() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	difficulties := entity.QuestionDifficulty(0).GetAllDifficulties()
	categories := entity.Category("all_categories").GetCategories()

	for _, dif := range difficulties {
		for _, cat := range categories {
			key := s.getKey(cat, dif)
			length, _ := s.setRepo.SetLength(ctx, key)

			if uint(length) < s.Config.NumberOfQuestions {
				wg.Add(1)
				go s.createQuiz(cat, dif, s.Config.NumberOfQuestions-uint(length), &wg)
			}
		}
	}

	wg.Wait()
}

func (s *Service) createQuiz(category entity.Category, difficulty entity.QuestionDifficulty, num uint, wg *sync.WaitGroup) {
	defer wg.Done()
	key := s.getKey(category, difficulty)

	for i := 0; i < int(num); i++ {
		localCtx, cancel := context.WithTimeout(context.Background(), s.Config.LocalContextTimeOut)
		questionIds, err := s.dbRepo.GetRandomQuestions(localCtx, category, difficulty, s.Config.NumberOfQuestions)
		if err != nil {
			logger.Warn(err, "failed to get random questions")
			metrics.FailedGetRandomQuestions.Inc()
			cancel()

			continue
		}

		payload := protobufencodedecode.EncodeQuizSvcQuiz(entity.Quiz{QuestionIDs: questionIds})
		s.setRepo.SetAdd(localCtx, key, payload)
		cancel()
	}
}

func (s *Service) GetQuiz(ctx context.Context, request quizparam.GetQuizRequest) (quizparam.GetQuizResponse, error) {
	const operation = "quizservice.GetQuiz"

	key := s.getKey(request.Category, request.Difficulty)
	value, err := s.setRepo.SetPop(ctx, key)
	if err != nil {
		logger.Warn(err, "pop quiz from set is failed")

		return quizparam.GetQuizResponse{},
			richerror.NewRichError(operation).
				WithError(err).
				WithMessage("pop quiz from set is failed").
				WithMeta(map[string]interface{}{"category": request.Category, "difficulty": request.Difficulty})
	}

	quiz := protobufencodedecode.DecodeQuizSvcQuiz(value)

	return quizparam.GetQuizResponse{QuestionIds: quiz.QuestionIDs}, nil
}

//getKey
/*
* Naming convention for key in key-value data structure,
* quiz-pool:{difficulty}:{category}
 */
func (s *Service) getKey(category entity.Category, difficulty entity.QuestionDifficulty) string {
	return fmt.Sprintf("%s:%d:%s", s.Config.Prefix, difficulty, category)
}
