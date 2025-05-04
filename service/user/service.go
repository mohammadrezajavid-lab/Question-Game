package user

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"gocasts.ir/go-fundamentals/gameapp/entity"
	"gocasts.ir/go-fundamentals/gameapp/pkg"
	"gocasts.ir/go-fundamentals/gameapp/service/contract"
	"log"
	"time"
)

type Service struct {
	userRepository contract.UserRepository
	signKey        []byte
}

type Info struct {
	Name   string
	UserId uint
}

type Claims struct {
	jwt.RegisteredClaims
	Info
}

func NewService(userRepository contract.UserRepository, signKey []byte) *Service {
	return &Service{
		userRepository: userRepository,
		signKey:        signKey,
	}
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

	user, gErr := s.userRepository.GetUserByPhoneNumber(req.PhoneNumber)
	if gErr != nil {
		log.Println(gErr.Error())
		return nil, fmt.Errorf("phoneNumber or password incorect")
	}

	if !pkg.CheckPasswordHash(req.Password, user.HashedPassword) {
		return nil, fmt.Errorf("phoneNumber or password incorrect")
	}

	// TODO - implement Me : If the user exists

	accessToken, cErr := createToken(user.ID, user.Name, s.signKey)
	if cErr != nil {
		return nil, fmt.Errorf("unexpected error: %w", cErr)
	}

	return NewLoginResponse(accessToken), nil
}

// All Request Inputs for Interaction/Service Should be Sanitized.

func (s *Service) Profile(req *ProfileRequest) (*ProfileResponse, error) {

	user, err := s.userRepository.GetUserById(req.UserId)
	if err != nil {

		// I don't expect the repository call return "not found user record" error,
		// because I assume the interaction input in sanitized
		// TODO - we can use Rich Error.
		return nil, fmt.Errorf("unexpected error: %w", err)
	}

	return NewProfileResponse(user.Name), nil
}

func createToken(userId uint, name string, signKey []byte) (string, error) {

	// create a signer for SHA 256
	t := jwt.New(jwt.GetSigningMethod(jwt.SigningMethodHS256.Alg()))

	// set our claims
	t.Claims = &Claims{
		// set the expiry time
		jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7))},
		Info{name, userId},
	}
	// create token string
	return t.SignedString(signKey)
}
