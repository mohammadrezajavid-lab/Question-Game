package uservalidator

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation"
	"golang.project/go-fundamentals/gameapp/entity"
	"golang.project/go-fundamentals/gameapp/pkg/errormessage"
	"regexp"
)

type Repository interface {
	IsPhoneNumberUniq(phoneNumber string) (bool, error)
	GetUserByPhoneNumber(phoneNumber string) (*entity.User, error)
}

type Validator struct {
	repository Repository
}

func NewValidator(repository Repository) Validator {
	return Validator{repository: repository}
}

func checkPasswordRegex() validation.RuleFunc {

	return func(value interface{}) error {
		password, ok := value.(string)
		if !ok {
			return errors.New(errormessage.ErrorMsgInvalidPhoneType)
		}

		var (
			hasUpper   bool = regexp.MustCompile(`[A-Z]`).MatchString(password)
			hasLower   bool = regexp.MustCompile(`[a-z]`).MatchString(password)
			hasNumber  bool = regexp.MustCompile(`[0-9]`).MatchString(password)
			hasSpecial bool = regexp.MustCompile(`[@%!%*?&#]`).MatchString(password)
		)

		if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
			return errors.New(errormessage.ErrorMsgInvalidPasswordRegex)
		}
		return nil
	}
}

func checkPhoneNumberRegex() validation.RuleFunc {

	return func(value interface{}) error {
		phoneNumber, ok := value.(string)
		if !ok {
			return errors.New(errormessage.ErrorMsgInvalidPhoneType)
		}

		var (
			hasValidPrefix     = regexp.MustCompile(`^(?:\+989|09|9)`).MatchString(phoneNumber)
			hasNineDigitsAfter = regexp.MustCompile(`^(?:\+989|09|9)\d{9}$`).MatchString(phoneNumber)
		)

		if !hasValidPrefix {
			return errors.New(errormessage.ErrorMsgInvalidPhoneNumberRegex1)
		}

		if !hasNineDigitsAfter {
			return errors.New(errormessage.ErrorMsgInvalidPhoneNumberRegex2)
		}

		return nil
	}
}
