package matchingservice

import (
	"context"
	"golang.project/go-fundamentals/gameapp/entity"
	"golang.project/go-fundamentals/gameapp/param/matchingparam"
	"golang.project/go-fundamentals/gameapp/pkg/richerror"
	"log"
	"time"
)

type Repository interface {
	AddToWaitingList(ctx context.Context, userId uint, category entity.Category) error
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

func (s *Service) AddToWaitingList(ctx context.Context, req *matchingparam.AddToWaitingListRequest) (*matchingparam.AddToWaitingListResponse, error) {
	const operation = richerror.Operation("matchingservice.AddToWaitingList")

	// req.Category should be sanitized before sent to service layer
	err := s.repo.AddToWaitingList(ctx, req.UserId, req.Category)
	if err != nil {
		return nil, richerror.NewRichError(operation).WithError(err)
	}

	return matchingparam.NewAddToWaitingListResponse(s.config.WaitingTimeOut), nil
}

func (s *Service) MatchWaitedUser(req *matchingparam.MatchWaitedUserRequest) (*matchingparam.MatchWaitedUserResponse, error) {
	log.Println("run MatchWaitedUser", time.Now())
	return nil, nil
}
