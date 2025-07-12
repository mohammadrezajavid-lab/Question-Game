package entity

type GameEvent struct {
	QuestionIds []uint
	PlayerIds   []uint
	GameId      uint
}

type CreatedGame struct {
	GameId    uint
	PlayerIds []uint
}

func NewCreatedGame(gameId uint, playerIds []uint) CreatedGame {
	return CreatedGame{
		GameId:    gameId,
		PlayerIds: playerIds,
	}
}
