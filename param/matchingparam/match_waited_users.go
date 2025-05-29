package matchingparam

import "golang.project/go-fundamentals/gameapp/entity"

type MatchWaitedUserRequest struct {
	Category entity.Category `json:"category"`
}

func NewMatchWaitedUserRequest(category entity.Category) *MatchWaitedUserRequest {
	return &MatchWaitedUserRequest{category}
}

type MatchWaitedUserResponse struct {
	WaitedUsers []WaitedUser
}
