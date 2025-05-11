package dto

type LoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

func NewLoginRequest() *LoginRequest {

	return &LoginRequest{
		PhoneNumber: "phoneNumber",
		Password:    "password",
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
