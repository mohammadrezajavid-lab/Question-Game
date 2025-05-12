package uservalidator

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation"
	"regexp"
)

type Repository interface {
	IsPhoneNumberUniq(phoneNumber string) (bool, error)
}

type Validator struct {
	repository Repository
}

func NewValidator(repository Repository) *Validator {
	return &Validator{repository: repository}
}

func checkPasswordRegex() validation.RuleFunc {

	return func(value interface{}) error {
		password, ok := value.(string)
		if !ok {
			return errors.New("invalid password type")
		}

		var (
			hasUpper   bool = regexp.MustCompile(`[A-Z]`).MatchString(password)
			hasLower   bool = regexp.MustCompile(`[a-z]`).MatchString(password)
			hasNumber  bool = regexp.MustCompile(`[0-9]`).MatchString(password)
			hasSpecial bool = regexp.MustCompile(`[@%!%*?&#]`).MatchString(password)
		)

		if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
			return errors.New("password must contain upper, lower, digit, and special character[@%!%*?&#]")
		}
		return nil
	}
}
