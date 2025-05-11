package dto

import "golang.project/go-fundamentals/gameapp/entity"

type RegisterRequest struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

func NewRegisterRequest() *RegisterRequest {
	return &RegisterRequest{Name: "name", PhoneNumber: "phoneNumber", Password: "password"}
}

type RegisterResponse struct {
	User struct {
		Id          uint   `json:"id"`
		PhoneNumber string `json:"phone_number"`
		Name        string `json:"name"`
	} `json:"user"`
}

func NewRegisterResponse(user *entity.User) *RegisterResponse {
	return &RegisterResponse{User: struct {
		Id          uint   `json:"id"`
		PhoneNumber string `json:"phone_number"`
		Name        string `json:"name"`
	}{Id: user.ID, PhoneNumber: user.PhoneNumber, Name: user.Name}}
}
