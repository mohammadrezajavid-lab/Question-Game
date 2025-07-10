package quizservice

import (
	"context"
	"golang.project/go-fundamentals/gameapp/entity"
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
	GetKey(category entity.Category, difficulty entity.QuestionDifficulty) string
}

type DBRepository interface {
	GetRandomQuestions(ctx context.Context, category entity.Category, difficulty entity.QuestionDifficulty, numberOfQuestion uint) ([]uint, error)
}

type Config struct {
	ContextTimeOut    time.Duration `mapstructure:"context_time_out"`
	NumberOfQuestions uint          `mapstructure:"number_of_questions"`
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

func (s *Service) GenerateQuiz(ctx context.Context) {
	difficulties := entity.QuestionDifficulty(0).GetAllDifficulties()
	categories := entity.Category("all_categories").GetCategories()

	var wg sync.WaitGroup
	for _, dif := range difficulties {
		for _, cat := range categories {
			key := s.setRepo.GetKey(cat, dif)
			length, _ := s.setRepo.SetLength(ctx, key)
			if uint(length) < s.Config.NumberOfQuestions {
				wg.Add(1)
				/*
					Since the number of combinations of category and difficulty is limited,
					goroutines have been used. However,
					if the number of these combinations increases significantly,
					it would be better to use a worker pool.
				*/
				go s.createQuiz(ctx, cat, dif, s.Config.NumberOfQuestions-uint(length), &wg)
			}
		}
	}
}

func (s *Service) createQuiz(ctx context.Context, category entity.Category, difficulty entity.QuestionDifficulty, num uint, wg *sync.WaitGroup) {
	defer wg.Done()
	key := s.setRepo.GetKey(category, difficulty)

	for i := 0; i < int(num); i++ {
		questionIds, _ := s.dbRepo.GetRandomQuestions(ctx, category, difficulty, s.Config.NumberOfQuestions)
		quizPayload := protobufencodedecode.EncodeQuizSvcQuiz(entity.Quiz{QuestionIDs: questionIds})
		s.setRepo.SetAdd(ctx, key, quizPayload)
	}
}

func (s *Service) GetQuiz(ctx context.Context, request quizparam.GetQuizRequest) (quizparam.GetQuizResponse, error) {
	const operation = "quizservice.GetQuiz"

	key := s.setRepo.GetKey(request.Category, request.Difficulty)

	value, err := s.setRepo.SetPop(ctx, key)
	if err != nil {
		return quizparam.GetQuizResponse{},
			richerror.NewRichError(operation).
				WithError(err).
				WithMessage("pop quiz from set is failed").
				WithMeta(map[string]interface{}{"category": request.Category, "difficulty": request.Difficulty})
	}

	quiz := protobufencodedecode.DecodeQuizSvcQuiz(value)

	return quizparam.GetQuizResponse{QuestionIds: quiz.QuestionIDs}, nil
}
