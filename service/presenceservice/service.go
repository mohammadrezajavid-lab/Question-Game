package presenceservice

import (
	"context"
	"fmt"
	"golang.project/go-fundamentals/gameapp/entity"
	"golang.project/go-fundamentals/gameapp/param/presenceparam"
	"golang.project/go-fundamentals/gameapp/pkg/richerror"
	"time"
)

type Repository interface {
	Upsert(ctx context.Context, key string, timestamp int64, expirationTime time.Duration) error
	GetPresence(ctx context.Context, key string) (entity.Presence, error)
}
type Config struct {
	ExpirationTime time.Duration `mapstructure:"expiration_time"`
	Prefix         string        `mapstructure:"prefix"`
}

type Service struct {
	Repo   Repository
	Config Config
}

func New(repo Repository, config Config) *Service {
	return &Service{Repo: repo, Config: config}
}

func (s Service) Upsert(ctx context.Context, req presenceparam.UpsertPresenceRequest) (presenceparam.UpsertPresenceResponse, error) {
	const operation = "presenceservice.Upsert"

	key := s.GetKey(req.UserId)
	if err := s.Repo.Upsert(ctx, key, req.TimeStamp, s.Config.ExpirationTime); err != nil {

		return presenceparam.UpsertPresenceResponse{}, richerror.NewRichError(operation).WithError(err)
	}

	return presenceparam.NewUpsertPresenceResponse(req.TimeStamp), nil
}

func (s Service) GetPresence(ctx context.Context, request presenceparam.GetPresenceRequest) (presenceparam.GetPresenceResponse, error) {

	const operation = "presenceservice.GetPresence"

	tmpPresenceItem := presenceparam.NewPresenceItem()
	presenceItems := make([]presenceparam.PresenceItem, 0)

	for _, userId := range request.UserIds {
		key := s.GetKey(userId)

		presence, err := s.Repo.GetPresence(ctx, key)
		if err != nil {

			return presenceparam.GetPresenceResponse{}, richerror.NewRichError(operation).WithError(err)
		}
		presence.UserId = userId

		tmpPresenceItem.UserId = presence.UserId
		tmpPresenceItem.Timestamp = presence.Timestamp

		presenceItems = append(presenceItems, tmpPresenceItem)
	}

	return presenceparam.NewGetPresenceResponse(presenceItems), nil
}

func (s Service) GetKey(userId uint) string {
	return fmt.Sprintf("%s:%d", s.Config.Prefix, userId)
}
