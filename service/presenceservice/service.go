package presenceservice

import (
	"context"
	"fmt"
	"golang.project/go-fundamentals/gameapp/param/presenceparam"
	"golang.project/go-fundamentals/gameapp/pkg/richerror"
	"time"
)

type Repository interface {
	Upsert(ctx context.Context, key string, timestamp int64, expirationTime time.Duration) error
}

type Config struct {
	ExpirationTime time.Duration `mapstructure:"expiration_time"`
	Prefix         string        `mapstructure:"prefix"`
}

type Service struct {
	repo   Repository
	config Config
}

func New(repo Repository, config Config) Service {
	return Service{repo: repo, config: config}
}

func (s *Service) Upsert(ctx context.Context, req *presenceparam.UpsertPresenceRequest) (*presenceparam.UpsertPresenceResponse, error) {
	const operation = "presenceservice.Upsert"

	key := fmt.Sprintf("%s:%d", s.config.Prefix, req.UserId)
	if err := s.repo.Upsert(ctx, key, req.TimeStamp, s.config.ExpirationTime); err != nil {

		return nil, richerror.NewRichError(operation).WithError(err)
	}

	return nil, nil
}

func (s *Service) GetPresence(ctx context.Context, request presenceparam.GetPresenceRequest) (presenceparam.GetPresenceResponse, error) {
	// TODO - implement me

	fmt.Println("presenceservice.GetPresence: request: 	", request)

	return presenceparam.GetPresenceResponse{Items: []presenceparam.PresenceItem{
		{UserId: 1, Timestamp: 2343242342},
		{UserId: 2, Timestamp: 2343242343},
		{UserId: 3, Timestamp: 2343242344},
		{UserId: 4, Timestamp: 2343242345},
		{UserId: 5, Timestamp: 2343242346},
	}}, nil
}
