package matchingservice

import (
	"context"
	"errors"
	"golang.project/go-fundamentals/gameapp/contract/broker"
	"golang.project/go-fundamentals/gameapp/entity"
	"golang.project/go-fundamentals/gameapp/logger"
	"golang.project/go-fundamentals/gameapp/metrics"
	"golang.project/go-fundamentals/gameapp/param/matchingparam"
	"golang.project/go-fundamentals/gameapp/param/presenceparam"
	"golang.project/go-fundamentals/gameapp/pkg/protobufencodedecode"
	"golang.project/go-fundamentals/gameapp/pkg/richerror"
	"golang.project/go-fundamentals/gameapp/pkg/search"
	"golang.project/go-fundamentals/gameapp/pkg/timestamp"
	"sync"
	"time"
)

type Repository interface {
	AddToWaitingList(ctx context.Context, userId uint, category entity.Category) error
	GetWaitedUserByCategory(ctx context.Context, category entity.Category) ([]matchingparam.WaitedUser, error)
	RemoveUserFromWaitingList(userIds []uint, category entity.Category)
}

type PresenceClient interface {
	GetPresence(ctx context.Context, request presenceparam.GetPresenceRequest) (presenceparam.GetPresenceResponse, error)
}

type Config struct {
	WaitingTimeOut          time.Duration `mapstructure:"waiting_time_out"`
	ContextTimeOut          time.Duration `mapstructure:"context_time_out"`
	OnlineThresholdDuration time.Duration `mapstructure:"online_threshold_duration"`
}

type Service struct {
	config         Config
	repo           Repository
	presenceClient PresenceClient
	publisher      broker.Published
}

func (s *Service) GetConfig() Config {
	if s != nil {
		return s.config
	}

	return Config{}
}

func NewService(config Config, repo Repository, presenceClient PresenceClient, publisher broker.Published) Service {
	return Service{
		config:         config,
		repo:           repo,
		presenceClient: presenceClient,
		publisher:      publisher,
	}
}

func (s *Service) AddToWaitingList(ctx context.Context, req *matchingparam.AddToWaitingListRequest) (*matchingparam.AddToWaitingListResponse, error) {
	const operation = richerror.Operation("matchingservice.AddToWaitingList")

	err := s.repo.AddToWaitingList(ctx, req.UserId, req.Category)
	if err != nil {
		metrics.FailedAddToWaitingListCounter.Inc()

		return nil, richerror.NewRichError(operation).WithError(err)
	}

	return matchingparam.NewAddToWaitingListResponse(s.config.WaitingTimeOut), nil
}

func (s *Service) MatchWaitedUsers(ctx context.Context) error {
	const operation = "matchingservice.MatchWaitedUsers"

	categories := entity.Category("all_categories").GetCategories()

	var wg sync.WaitGroup
	errCh := make(chan error, len(categories)) // buffer to prevent blocking

	/* TODO - If you have a large number of categories, you can use a worker pool instead of a goroutine for each one.*/

	for _, cat := range categories {
		wg.Add(1)
		go func(category entity.Category) {
			defer wg.Done()

			waitingListByCategory, err := s.getWaitingListByCategory(ctx, matchingparam.NewMatchWaitedUserRequest(category))
			if err != nil {
				errCh <- richerror.NewRichError(operation).WithError(err).WithMeta(map[string]interface{}{"category": category})
				return
			}

			waitedUsers := waitingListByCategory.WaitedUsers
			var waitedUserIds = make([]uint, 0, len(waitedUsers))
			for _, wu := range waitedUsers {
				waitedUserIds = append(waitedUserIds, wu.UserId)
			}

			presenceResponse, pErr := s.presenceClient.GetPresence(ctx, presenceparam.NewGetPresenceRequest(waitedUserIds))
			if pErr != nil {
				metrics.FailedGetPresenceClientCounter.Inc()
				logger.Warn(pErr, "failed_get_presence")

				errCh <- richerror.NewRichError(operation).WithError(pErr).WithMeta(map[string]interface{}{"category": category, "step": "get_presence"})
				return
			}

			var finalListWaitedUsers = make([]matchingparam.WaitedUser, 0)
			var allUsersToRemoved = make([]uint, 0)
			for _, waitedUser := range waitingListByCategory.WaitedUsers {

				userPresence, ok := search.BinarySearch(presenceResponse.Items, waitedUser.UserId)

				if ok && userPresence.Timestamp > timestamp.Add(-1*s.config.OnlineThresholdDuration) &&
					waitedUser.Timestamp > timestamp.Add(-1*s.config.WaitingTimeOut) {

					finalListWaitedUsers = append(finalListWaitedUsers, waitedUser)
					continue
				}

				allUsersToRemoved = append(allUsersToRemoved, waitedUser.UserId)
			}

			for j := 0; j+1 < len(finalListWaitedUsers); j += 2 {
				mu := entity.NewMatchedUsers(category, []uint{finalListWaitedUsers[j].UserId, finalListWaitedUsers[j+1].UserId})

				//log.Printf("user_id: [%d], user_id: [%d], for category: [%s] is matched\n", mu.UserIds[0], mu.UserIds[1], mu.Category)

				payload := protobufencodedecode.EncodeMatchingWaitedUsersEvent(mu)
				go s.publisher.PublishEvent(entity.MatchingUsersMatchedEvent, payload)

				//log.Printf("user_id: [%d], user_id: [%d], for category: [%s] is matched and published.\n", mu.UserIds[0], mu.UserIds[1], mu.Category)

				allUsersToRemoved = append(allUsersToRemoved, mu.UserIds...)
			}

			go s.repo.RemoveUserFromWaitingList(allUsersToRemoved, category)
		}(cat)
	}

	wg.Wait()
	close(errCh)

	var errs []error
	for err := range errCh {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}

func (s *Service) getWaitingListByCategory(ctx context.Context, req *matchingparam.MatchWaitedUserRequest) (matchingparam.MatchWaitedUserResponse, error) {
	const operation = richerror.Operation("matchingservice.MatchWaitedUser")

	waitedUsers, err := s.repo.GetWaitedUserByCategory(ctx, req.Category)
	if err != nil {
		return matchingparam.MatchWaitedUserResponse{},
			richerror.NewRichError(operation).WithError(err)
	}

	return matchingparam.MatchWaitedUserResponse{WaitedUsers: waitedUsers}, nil
}
