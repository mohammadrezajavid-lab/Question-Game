package gameparam

import "golang.project/go-fundamentals/gameapp/entity"

type WaitedUser struct {
	Timestamp  int64 // score
	UserId     uint  // member
	Category   entity.Category
	Difficulty entity.QuestionDifficulty
}

func NewWaitedUser(timestamp int64, userId uint, category entity.Category, difficulty entity.QuestionDifficulty) WaitedUser {
	return WaitedUser{
		Timestamp:  timestamp,
		UserId:     userId,
		Category:   category,
		Difficulty: difficulty,
	}
}
