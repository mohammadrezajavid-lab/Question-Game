package backofficeuserservice

import (
	"context"
	"golang.project/go-fundamentals/gameapp/entity"
)

type UserRepository interface {
	ListUsers(ctx context.Context) ([]entity.User, error)
}
type Service struct {
	userRepo UserRepository
}

func NewService(repository UserRepository) Service {
	return Service{userRepo: repository}
}

func (s Service) ListAllUsers(ctx context.Context) ([]entity.User, error) {

	listUsers, err := s.userRepo.ListUsers(ctx)
	if err != nil {
		return nil, err
	}

	//// TODO - implement me
	//list := make([]entity.User, 0)
	//
	//list = append(list, entity.User{
	//	Id:             0,
	//	Name:           "fake",
	//	PhoneNumber:    "fake",
	//	HashedPassword: "fake",
	//	Role:           entity.AdminRole,
	//})

	return listUsers, nil
}
