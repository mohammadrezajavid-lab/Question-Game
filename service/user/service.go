package user

import (
	"fmt"
	"gocasts.ir/go-fundamentals/gameapp/entity"
	"gocasts.ir/go-fundamentals/gameapp/pkg"
	"gocasts.ir/go-fundamentals/gameapp/service/contract"
)

type Service struct {
	userRepository contract.UserRepository
}

func NewService(userRepository contract.UserRepository) *Service {
	return &Service{userRepository: userRepository}
}

func (s *Service) Register(req *RegisterRequest) (*RegisterResponse, error) {

	// TODO - We should verify Phone Number by Verification SMS Code

	// validate phone number
	if !pkg.IsPhoneNumberValid(req.PhoneNumber) {

		return NewRegisterResponse(entity.NewUser("", "")),
			fmt.Errorf("phone number is invalid")
	}

	// check uniqueness of phone number
	if isUniq, err := s.userRepository.IsPhoneNumberUniq(req.PhoneNumber); !isUniq || err != nil {
		if !isUniq {

			return NewRegisterResponse(entity.NewUser("", "")),
				fmt.Errorf("phone number is not uniq")
		}

		return NewRegisterResponse(entity.NewUser("", "")),
			fmt.Errorf("unexpected error: %w", err)
	}

	// validate name
	if !pkg.IsNameValid(req.Name) {

		return NewRegisterResponse(entity.NewUser("", "")),
			fmt.Errorf("name length should be greater than 3")
	}

	// create new user in storage
	newUser, cErr := s.userRepository.RegisterUser(entity.NewUser(req.Name, req.PhoneNumber))
	if cErr != nil {
		return NewRegisterResponse(entity.NewUser("", "")),
			fmt.Errorf("unexpected error: %w", cErr)
	}

	// return created user
	return NewRegisterResponse(newUser), nil
}
