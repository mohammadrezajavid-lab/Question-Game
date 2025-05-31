package matchingservice

import (
	"context"
	"errors"
	"fmt"
	"golang.project/go-fundamentals/gameapp/entity"
	"golang.project/go-fundamentals/gameapp/param/matchingparam"
	"golang.project/go-fundamentals/gameapp/param/presenceparam"
	"golang.project/go-fundamentals/gameapp/pkg/richerror"
	"golang.project/go-fundamentals/gameapp/pkg/search"
	"golang.project/go-fundamentals/gameapp/pkg/sort"
	"golang.project/go-fundamentals/gameapp/pkg/timestamp"
	"log"
	"sync"
	"time"
)

type Repository interface {
	AddToWaitingList(ctx context.Context, userId uint, category entity.Category) error
	GetWaitedUserByCategory(ctx context.Context, category entity.Category) ([]matchingparam.WaitedUser, error)
	RemoveUserFromWaitedList(ctx context.Context, userId uint, category entity.Category) error
}

type PresenceClient interface {
	GetPresence(ctx context.Context, request presenceparam.GetPresenceRequest) (presenceparam.GetPresenceResponse, error)
}

type Service struct {
	config         Config
	repo           Repository
	presenceClient PresenceClient
}

type Config struct {
	WaitingTimeOut time.Duration `mapstructure:"waiting_time_out"`
}

func NewService(config Config, repo Repository, presenceClient PresenceClient) Service {
	return Service{config: config, repo: repo, presenceClient: presenceClient}
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

func (s *Service) MatchWaitedUsers(ctx context.Context) error {
	const operation = "matchingservice.MatchWaitedUsers"

	log.Println("Executing matchWaitedUser job at:", time.Now())

	categories := entity.Category("all_categories").GetCategories()

	var wg sync.WaitGroup
	errCh := make(chan error, len(categories)) // buffer to prevent blocking

	/* TODO - If you have a large number of categories, you can use a worker pool instead of a goroutine for each one.*/

	for _, cat := range categories {
		wg.Add(1)
		go func(category entity.Category) {
			defer wg.Done()

			waitingList, err := s.getWaitingListByCategory(ctx, matchingparam.NewMatchWaitedUserRequest(category))
			if err != nil {
				errCh <- richerror.NewRichError(operation).WithError(err).WithMeta(map[string]interface{}{"category": category})
				return
			}

			waitedUsers := waitingList.WaitedUsers
			var waitedUsersId = make([]uint, 0, len(waitedUsers))
			for _, wu := range waitedUsers {
				waitedUsersId = append(waitedUsersId, wu.UserId)
			}

			getPresenceResponse, pErr := s.presenceClient.GetPresence(ctx, presenceparam.NewGetPresenceRequest(waitedUsersId))
			if pErr != nil {
				// TODO - update metrics
				// TODO - log error
				log.Println(pErr.Error())
				return
			}

			fmt.Println(operation)
			fmt.Println("getPresenceResponse: ", getPresenceResponse)

			waitedUsersId = sort.NewQuickSort(waitedUsersId).Sort()

			// merge getPresenceResponse and waitedUsers to create finalListWaitedUsers
			var finalListWaitedUsers = make([]matchingparam.WaitedUser, 0)
			for _, item := range getPresenceResponse.Items {
				if item.Timestamp > timestamp.Add(-20*time.Second) && search.BinarySearch(waitedUsersId, item.UserId) {
					finalListWaitedUsers = append(finalListWaitedUsers, matchingparam.NewWaitedUser(item.Timestamp, item.UserId, category))

					continue
				}

				if rErr := s.repo.RemoveUserFromWaitedList(ctx, item.UserId, category); rErr != nil {
					errCh <- richerror.NewRichError(operation).WithError(rErr)

					// TODO - update metrics
					// TODO - log error

					continue
				}
			}

			fmt.Println(operation)
			fmt.Println("finalListWaitedUsers: ", finalListWaitedUsers)
			for j := 0; j+1 < len(finalListWaitedUsers); j += 2 {
				mu := entity.NewMatchedUsers(category, []uint{finalListWaitedUsers[j].UserId, finalListWaitedUsers[j+1].UserId})
				fmt.Println(mu)
				//now published a new event for mu and send to message broker

				for _, userId := range mu.UserIds {
					if rErr := s.repo.RemoveUserFromWaitedList(ctx, userId, category); rErr != nil {
						errCh <- richerror.NewRichError(operation).WithError(rErr)

						// TODO - update metrics
						// TODO - log error

						continue
					}
				}
			}
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
