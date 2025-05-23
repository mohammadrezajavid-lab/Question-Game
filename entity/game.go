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
type Game struct {
	Id         uint
	Category   Category
	QuestionId []uint
	PlayerIds  []uint
	WinnerId   uint
	StartTime  time.Time
}

func NewGame(id uint, category Category) *Game {
	return &Game{
		Id:         id,
		Category:   category,
		QuestionId: nil,
		PlayerIds:  nil,
		WinnerId:   0,
		StartTime:  time.Now(),
	}
}
