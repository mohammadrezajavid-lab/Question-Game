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

func New(repo Repository, config Config) *Service {
	return &Service{repo: repo, config: config}
}

func (s *Service) Upsert(ctx context.Context, req *presenceparam.UpsertPresenceRequest) (*presenceparam.UpsertPresenceResponse, error) {
	const operation = "presenceservice.Upsert"

	key := fmt.Sprintf("%s:%d", s.config.Prefix, req.UserId)
	if err := s.repo.Upsert(ctx, key, req.TimeStamp, s.config.ExpirationTime); err != nil {

		return nil, richerror.NewRichError(operation).WithError(err)
	}

	return nil, nil
}

func (s *Service) IsOnline() bool {
	return false
}
