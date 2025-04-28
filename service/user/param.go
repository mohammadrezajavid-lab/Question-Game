package user

import "gocasts.ir/go-fundamentals/gameapp/entity"

type RegisterRequest struct {
	Name        string
	PhoneNumber string
}

func NewRegisterRequest(name, phoneNumber string) *RegisterRequest {
	return &RegisterRequest{Name: name, PhoneNumber: phoneNumber}
}

type RegisterResponse struct {
	User *entity.User
}

func NewRegisterResponse(user *entity.User) *RegisterResponse {
	return &RegisterResponse{User: user}
}
