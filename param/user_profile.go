package param

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
