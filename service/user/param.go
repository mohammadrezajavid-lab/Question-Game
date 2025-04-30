package user

import (
	"gocasts.ir/go-fundamentals/gameapp/entity"
)

type RegisterRequest struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

func NewRegisterRequest(name string, phoneNumber string, password string) *RegisterRequest {
	return &RegisterRequest{Name: name, PhoneNumber: phoneNumber, Password: password}
}

type RegisterResponse struct {
	User *entity.User `json:"user"`
}

func NewRegisterResponse(user *entity.User) *RegisterResponse {
	return &RegisterResponse{User: user}
}
