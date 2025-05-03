package user

import (
	"fmt"
	"gocasts.ir/go-fundamentals/gameapp/entity"
	"gocasts.ir/go-fundamentals/gameapp/pkg"
	"gocasts.ir/go-fundamentals/gameapp/service/contract"
	"log"
)

type Service struct {
	userRepository contract.UserRepository
}

func NewService(userRepository contract.UserRepository) *Service {
	return &Service{userRepository: userRepository}
}

func (s *Service) Register(req *RegisterRequest) (*RegisterResponse, error) {

	// TODO - We should verify Phone Number by Verification SMS Code

	var emptyUser *entity.User = entity.NewUser("", "", "")

	// validate phone number
	if !pkg.IsPhoneNumberValid(req.PhoneNumber) {

		return NewRegisterResponse(emptyUser), fmt.Errorf("phone number is invalid")
	}

	// check uniqueness of phone number
	if isUniq, err := s.userRepository.IsPhoneNumberUniq(req.PhoneNumber); !isUniq || err != nil {
		if !isUniq {

			return NewRegisterResponse(emptyUser), fmt.Errorf("phone number is not uniq")
		}

		return NewRegisterResponse(emptyUser), fmt.Errorf("unexpected error: %w", err)
	}

	// validate name
	if !pkg.IsNameValid(req.Name) {

		return NewRegisterResponse(emptyUser), fmt.Errorf("name length should be greater than 3")
	}

	// validate password
	// TODO - It is better to use Regex for password.
	if len(req.Password) < 8 {
		return NewRegisterResponse(emptyUser), fmt.Errorf("password length should be greater than 8")
	}

	hashedPassword, hashErr := pkg.HashPassword(req.Password)
	if hashErr != nil {
		log.Println("We encountered a problem in hashing the password using bcrypt, ", hashErr)
	}

	// create new user in storage
	newUser, cErr := s.userRepository.RegisterUser(entity.NewUser(req.Name, req.PhoneNumber, hashedPassword))
	if cErr != nil {
		return NewRegisterResponse(newUser), fmt.Errorf("unexpected error: %w", cErr)
	}

	// return created user
	return NewRegisterResponse(newUser), nil
}

func (s *Service) Login(req *LoginRequest) (*LoginResponse, error) {

	user, gErr := s.userRepository.GetUser(req.PhoneNumber)
	if gErr != nil {
		log.Println(gErr.Error())
		return nil, fmt.Errorf("phoneNumber or password incorect")
	}

	if !pkg.CheckPasswordHash(req.Password, user.HashedPassword) {
		return nil, fmt.Errorf("phoneNumber or password incorrect")
	}

	// TODO - implement Me : If the user exists

	return nil, nil
}
