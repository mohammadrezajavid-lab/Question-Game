package user

import "gocasts.ir/go-fundamentals/gameapp/entity"

type RegisterRequest struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phoneNumber"`
}

func NewRegisterRequest(name, phoneNumber string) *RegisterRequest {
	return &RegisterRequest{Name: name, PhoneNumber: phoneNumber}
}

type RegisterResponse struct {
	User *entity.User `json:"user"`
}

func NewRegisterResponse(user *entity.User) *RegisterResponse {
	return &RegisterResponse{User: user}
}
