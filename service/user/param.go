package user

import (
	"golang.project/go-fundamentals/gameapp/entity"
)

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

type LoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

func NewLoginRequest(phoneNumber string, password string) *LoginRequest {

	return &LoginRequest{
		PhoneNumber: phoneNumber,
		Password:    password,
	}
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewTokens(accessToken, refreshToken string) Tokens {
	return Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}

type UserInfo struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

func NewUserInfo(id uint, name string) UserInfo {
	return UserInfo{
		Id:   id,
		Name: name,
	}
}

type LoginResponse struct {
	User   UserInfo `json:"user"`
	Tokens Tokens   `json:"tokens"`
}

func NewLoginResponse(userInfo UserInfo, tokens Tokens) *LoginResponse {

	return &LoginResponse{
		User:   userInfo,
		Tokens: tokens,
	}
}

type ProfileRequest struct {
	UserId uint `json:"user_id"`
}

func NewProfileRequest(userId uint) *ProfileRequest {

	return &ProfileRequest{UserId: userId}
}

type ProfileResponse struct {
	Name string `json:"name"`
}

func NewProfileResponse(name string) *ProfileResponse {

	return &ProfileResponse{Name: name}
}
