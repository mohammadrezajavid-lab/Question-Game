package matchingservice

import (
	"golang.project/go-fundamentals/gameapp/entity"
	"golang.project/go-fundamentals/gameapp/param"
	"golang.project/go-fundamentals/gameapp/pkg/richerror"
	"time"
)

type Repository interface {
	AddToWaitingList(userId uint, category entity.Category) error
}

type Service struct {
	config Config
	repo   Repository
}

type Config struct {
	WaitingTimeOut time.Duration `mapstructure:"waiting_time_out"`
}

func NewService(config Config, repo Repository) *Service {
	return &Service{config: config, repo: repo}
}

func (s *Service) AddToWaitingList(req *param.AddToWaitingListRequest) (*param.AddToWaitingListResponse, error) {
	const operation = richerror.Operation("matchingservice.AddToWaitingList")

	// req.Category should be sanitized before sent to service layer
	err := s.repo.AddToWaitingList(req.UserId, req.Category)
	if err != nil {
		return nil, richerror.NewRichError(operation).WithError(err)
	}

	return param.NewAddToWaitingListResponse(s.config.WaitingTimeOut), nil
}
