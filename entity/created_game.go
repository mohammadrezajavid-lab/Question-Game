package entity

type CreatedGame struct {
	GameId      uint
	PlayerIds   []uint
	QuestionIds []uint
}

func NewCreatedGame(gameId uint, playerIds, questionIds []uint) CreatedGame {
	return CreatedGame{
		GameId:      gameId,
		PlayerIds:   playerIds,
		QuestionIds: questionIds,
	}
}
