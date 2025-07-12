package gamemysql

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"golang.project/go-fundamentals/gameapp/entity"
	"golang.project/go-fundamentals/gameapp/logger"
	"golang.project/go-fundamentals/gameapp/metrics"
	"golang.project/go-fundamentals/gameapp/pkg/errormessage"
	"golang.project/go-fundamentals/gameapp/pkg/richerror"
)

func (d *DataBase) CreateGame(ctx context.Context, game entity.Game) (entity.Game, error) {

	const operation = "gamemysql.CreateGame"
	const queryType = "insert"

	var result, eErr = d.dataBase.MysqlConnection.ExecContext(
		ctx,
		`INSERT INTO game_app_db.games(category, winner_id, start_time) VALUES(?, ?, ?)`,
		game.Category,
		game.WinnerId,
		game.StartTime,
	)

	metrics.DBQueryCounter.With(prometheus.Labels{"query_type": queryType}).Inc()

	if eErr != nil {
		logger.Warn(eErr, errormessage.ErrorMsgFailedExecuteQuery)

		metrics.DBFailedQueryCounter.With(prometheus.Labels{"query_type": queryType}).Inc()

		return entity.Game{}, richerror.NewRichError(operation).
			WithError(eErr).
			WithMessage(errormessage.ErrorMsgUnexpected).
			WithKind(richerror.KindUnexpected)
	}

	var gameId, _ = result.LastInsertId()
	game.Id = uint(gameId)

	return game, nil
}

func (d *DataBase) CreatePlayer(ctx context.Context, player entity.Player) (entity.Player, error) {

	const operation = "gamemysql.CreatePlayer"
	const queryType = "insert"

	var result, eErr = d.dataBase.MysqlConnection.ExecContext(
		ctx,
		`INSERT INTO game_app_db.players(user_id, game_id, score) VALUES(?,?,?)`,
		player.UserId,
		player.GameId,
		player.Score,
	)

	metrics.DBQueryCounter.With(prometheus.Labels{"query_type": queryType}).Inc()

	if eErr != nil {
		logger.Warn(eErr, errormessage.ErrorMsgFailedExecuteQuery)

		metrics.DBFailedQueryCounter.With(prometheus.Labels{"query_type": queryType}).Inc()

		return entity.Player{}, richerror.NewRichError(operation).
			WithError(eErr).
			WithMessage(errormessage.ErrorMsgUnexpected).
			WithKind(richerror.KindUnexpected)
	}

	var playerId, _ = result.LastInsertId()
	player.Id = uint(playerId)

	return player, nil
}
