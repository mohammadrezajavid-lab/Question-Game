package backofficeuserservice

import "golang.project/go-fundamentals/gameapp/entity"

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s Service) ListAllUsers() ([]entity.User, error) {
	// TODO - implement me
	list := make([]entity.User, 0)

	list = append(list, entity.User{
		Id:             0,
		Name:           "fake",
		PhoneNumber:    "fake",
		HashedPassword: "fake",
		Role:           entity.AdminRole,
	})

	return list, nil
}
