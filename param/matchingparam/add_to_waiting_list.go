package matchingparam

import (
	"golang.project/go-fundamentals/gameapp/entity"
	"time"
)

type AddToWaitingListRequest struct {
	UserId     uint                      `json:"-"`
	Category   entity.Category           `json:"category"`
	Difficulty entity.QuestionDifficulty `json:"difficulty"`
}

func NewAddToWaitingListRequest() *AddToWaitingListRequest {
	return &AddToWaitingListRequest{
		UserId:     0,
		Category:   "",
		Difficulty: 0,
	}
}

type AddToWaitingListResponse struct {
	TimeOut time.Duration `json:"time_out_in_nanoseconds"`
}

func NewAddToWaitingListResponse(timeOut time.Duration) *AddToWaitingListResponse {
	return &AddToWaitingListResponse{TimeOut: timeOut}
}
