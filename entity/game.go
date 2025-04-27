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
	CategoryId uint
	QuestionId []uint
	PlayerIds  []uint
	WinnerId   uint
	StartTime  time.Time
}
