package gameparam

import "golang.project/go-fundamentals/gameapp/entity"

type MatchWaitedUserRequest struct {
	Category   entity.Category           `json:"category"`
	Difficulty entity.QuestionDifficulty `json:"difficulty"`
}

func NewMatchWaitedUserRequest(category entity.Category, difficulty entity.QuestionDifficulty) *MatchWaitedUserRequest {
	return &MatchWaitedUserRequest{Category: category, Difficulty: difficulty}
}

type MatchWaitedUserResponse struct {
	WaitedUsers []WaitedUser
}
