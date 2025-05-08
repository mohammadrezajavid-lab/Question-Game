package uservalidator

import (
	"golang.project/go-fundamentals/gameapp/datatransferobject/userdto"
	"golang.project/go-fundamentals/gameapp/pkg/phonenumber"
	"golang.project/go-fundamentals/gameapp/pkg/richerror"
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

func (v Validator) ValidateRegisterRequest(req *userdto.RegisterRequest) error {

	const operation = "uservalidator.ValidateRegisterRequest"

	if !phonenumber.IsPhoneNumberValid(req.PhoneNumber) {

		return richerror.NewRichError(operation).
			WithMessage("phone number is invalid").
			WithKind(richerror.KindInvalid).
			WithMeta(map[string]interface{}{"phone_number": req.PhoneNumber})
	}

	if isUniq, err := v.repository.IsPhoneNumberUniq(req.PhoneNumber); err != nil || !isUniq {

		if err != nil {

			return richerror.NewRichError(operation).WithError(err)
		}

		return richerror.NewRichError(operation).
			WithMessage("phone number is not uniq").
			WithMeta(map[string]interface{}{"phone_number": req.PhoneNumber})
	}

	// TODO - Add 8 to config
	if len(req.Name) < 3 {

		return richerror.NewRichError(operation).
			WithMessage("name length should be greater than 3").
			WithKind(richerror.KindInvalid).
			WithMeta(map[string]interface{}{"name": req.Name})
	}

	// TODO - It is better to use Regex for password.
	// TODO - Add 8 to config
	if len(req.Password) < 8 {

		return richerror.NewRichError(operation).
			WithMessage("password length should be greater than 8").
			WithKind(richerror.KindInvalid).
			WithMeta(map[string]interface{}{"password": req.Password})
	}

	return nil
}
