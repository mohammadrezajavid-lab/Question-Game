package entity

import "time"

type PlayerAnswer struct {
	Id         uint
	PlayerId   uint
	QuestionId uint
	Choice     PossibleAnswerChoice
}
type Player struct {
	Id      uint
	UserId  uint
	GameId  uint
	Score   uint
	Answers []PlayerAnswer
}

func NewPlayer(userId uint, gameId uint) Player {
	return Player{
		Id:      0,
		UserId:  userId,
		GameId:  gameId,
		Score:   0,
		Answers: nil,
	}
}

type Game struct {
	Id          uint
	Category    Category
	QuestionIds []uint
	PlayerIds   []uint
	WinnerId    uint
	StartTime   time.Time
}

func NewGame(id uint, category Category) Game {
	return Game{
		Id:          id,
		Category:    category,
		QuestionIds: nil,
		PlayerIds:   nil,
		WinnerId:    0,
		StartTime:   time.Now(),
	}
}
