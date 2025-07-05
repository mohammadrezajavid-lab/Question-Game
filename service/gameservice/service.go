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
	rda       redis.Adapter
	gameRepo  Repository
	publisher broker.Published
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

		// Create newGame
		newGame, rErr := s.gameRepo.CreateGame(entity.NewGame(mu.Category))
		if rErr != nil {
			logger.Error(rErr, rErr.Error())
		}

		// Create players for game
		playerIds, pErr := s.createPlayers(mu.UserIds, newGame.Id)
		if pErr != nil {
			logger.Error(pErr, pErr.Error())
		}
		newGame.PlayerIds = playerIds

		// TODO - انتخاب بک مجموعه سوال مثلا 15 تایی از استخر سوال ها با توجه به دسته و سطح سختی یا آسانی تعیین شده توسط کاربر

		// Published CreatedGameEvent
		cg := entity.NewCreatedGame(newGame.Id, newGame.PlayerIds, newGame.QuestionIds)
		payload := protobufencodedecode.EncodeGameSvcCreatedGameEvent(cg)
		s.publisher.PublishEvent(entity.GameSvcCreatedGameEvent, payload)
	}
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
