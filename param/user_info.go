package param

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
