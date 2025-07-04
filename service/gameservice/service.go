package gameservice

import (
	"context"
	"golang.project/go-fundamentals/gameapp/adapter/redis"
	"golang.project/go-fundamentals/gameapp/entity"
	"golang.project/go-fundamentals/gameapp/logger"
	"golang.project/go-fundamentals/gameapp/pkg/protobufencodedecode"
	"time"
)

type Repository interface {
	RegisterPlayer(player entity.Player) (entity.Player, error)
	RegisterGame(game entity.Game) (entity.Game, error)
}

type Service struct {
	rda      redis.Adapter
	gameRepo Repository
}

func (s *Service) Start() {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	sub := s.rda.GetClient().Subscribe(ctx, entity.MatchingUsersMatchedEvent)
	defer sub.Close()
	//var mu = entity.NewMatchedUsers("", nil)
	for msg := range sub.Channel() {
		// decode
		mu := protobufencodedecode.DecodeMatchingWaitedUsersEvent(msg.Payload)

		newGame, rErr := s.gameRepo.RegisterGame(entity.NewGame(0, mu.Category))
		if rErr != nil {
			logger.Error(rErr, rErr.Error())
		}

		playerIds, pErr := s.createPlayers(mu.UserIds, newGame.Id)
		if pErr != nil {
			logger.Error(pErr, pErr.Error())
		}
		newGame.PlayerIds = playerIds
	}
}

func (s *Service) createPlayers(userIds []uint, gameId uint) ([]uint, error) {

	playerIds := make([]uint, 0, 2)

	for _, userId := range userIds {
		player, err := s.gameRepo.RegisterPlayer(entity.NewPlayer(userId, gameId))
		if err != nil {
			return nil, err
		}

		playerIds = append(playerIds, player.Id)
	}

	return playerIds, nil
}
