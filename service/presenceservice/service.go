package presenceservice

import (
	"context"
	"fmt"
	"golang.project/go-fundamentals/gameapp/entity"
	"golang.project/go-fundamentals/gameapp/param/presenceparam"
	"golang.project/go-fundamentals/gameapp/pkg/richerror"
	"log"
	"sort"
	"time"
)

type Repository interface {
	Upsert(ctx context.Context, key string, timestamp int64, expirationTime time.Duration) error
	GetPresences(ctx context.Context, keys []string) ([]entity.Presence, error)
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

	key := s.getKey(req.UserId)
	if err := s.Repo.Upsert(ctx, key, req.TimeStamp, s.Config.ExpirationTime); err != nil {

		return presenceparam.UpsertPresenceResponse{}, richerror.NewRichError(operation).WithError(err)
	}

	return presenceparam.NewUpsertPresenceResponse(req.TimeStamp), nil
}

func (s Service) GetPresence(ctx context.Context, request presenceparam.GetPresenceRequest) (presenceparam.GetPresenceResponse, error) {

	const operation = "presenceservice.GetPresence"

	keys := s.generateAllKey(request.UserIds)
	usersPresence, err := s.Repo.GetPresences(ctx, keys)
	if err != nil {
		log.Println(operation, ": ", err.Error())
		return presenceparam.GetPresenceResponse{}, richerror.NewRichError(operation).WithError(err)
	}

	items := make([]presenceparam.PresenceItem, 0, len(usersPresence)/2)
	for _, pre := range usersPresence {
		items = append(items, presenceparam.NewPresenceItem(pre.UserId, pre.Timestamp))
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].UserId < items[j].UserId
	})

	return presenceparam.NewGetPresenceResponse(items), nil
}

func (s Service) generateAllKey(userIds []uint) []string {
	keys := make([]string, 0, len(userIds))

	for _, id := range userIds {
		keys = append(keys, s.getKey(id))
	}

	return keys
}
func (s Service) getKey(userId uint) string {
	return fmt.Sprintf("%s:%d", s.Config.Prefix, userId)
}
