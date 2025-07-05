package gameservice

import (
	"context"
	"golang.project/go-fundamentals/gameapp/adapter/redis"
	"golang.project/go-fundamentals/gameapp/contract/broker"
	"golang.project/go-fundamentals/gameapp/entity"
	"golang.project/go-fundamentals/gameapp/logger"
	"golang.project/go-fundamentals/gameapp/pkg/protobufencodedecode"
	"time"
)

type Repository interface {
	CreatePlayer(player entity.Player) (entity.Player, error)
	CreateGame(game entity.Game) (entity.Game, error)
}

type Service struct {
	rda        redis.Adapter
	gameRepo   Repository
	publisher  broker.Publisher
	subscriber broker.Subscriber
}

func (s *Service) Start() {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	msgCh, err := s.subscriber.Subscribed(ctx, entity.MatchingUsersMatchedEvent)
	if err != nil {
		logger.Fatal(err, "failed to subscribe to broker topic")
	}

	for msg := range msgCh {
		go s.handleMatchedUsers(msg)
	}
}

func (s *Service) handleMatchedUsers(payload string) {
	// decode
	mu := protobufencodedecode.DecodeMatchingWaitedUsersEvent(payload)

	// Create newGame
	newGame, rErr := s.gameRepo.CreateGame(entity.NewGame(mu.Category))
	if rErr != nil {
		logger.Error(rErr, "failed to create game")
		return
	}

	// Create players for game
	playerIds, pErr := s.createPlayers(mu.UserIds, newGame.Id)
	if pErr != nil {
		logger.Error(pErr, "failed to create player")
		return
	}
	newGame.PlayerIds = playerIds

	// TODO - انتخاب بک مجموعه سوال مثلا 15 تایی از استخر سوال ها با توجه به دسته و سطح سختی یا آسانی تعیین شده توسط کاربر

	// Published CreatedGameEvent
	cg := entity.NewCreatedGame(newGame.Id, newGame.PlayerIds, newGame.QuestionIds)
	payloadCg := protobufencodedecode.EncodeGameSvcCreatedGameEvent(cg)
	s.publisher.Published(entity.GameSvcCreatedGameEvent, payloadCg)
}

func (s *Service) createPlayers(userIds []uint, gameId uint) ([]uint, error) {

	playerIds := make([]uint, 0, 2)

	for _, userId := range userIds {
		player, err := s.gameRepo.CreatePlayer(entity.NewPlayer(userId, gameId))
		if err != nil {
			return nil, err
		}

		playerIds = append(playerIds, player.Id)
	}

	return playerIds, nil
}
