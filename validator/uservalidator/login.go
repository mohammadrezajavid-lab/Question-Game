package uservalidator

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"golang.project/go-fundamentals/gameapp/dto"
	"golang.project/go-fundamentals/gameapp/pkg/richerror"
)

func (v *Validator) ValidateLoginRequest(req *dto.LoginRequest) (error, map[string]string) {

	const operation = "uservalidator.ValidateLoginRequest"

	if err := v.validateLoginRequest(req); err != nil {

		fieldErrors := make(map[string]string)

		valueErr, ok := err.(validation.Errors)
		if ok {
			for key, value := range valueErr {
				fieldErrors[key] = value.Error()
			}
		}

		return richerror.NewRichError(operation).
				WithError(err).
				WithMessage("invalid input").
				WithKind(richerror.KindInvalid),
			fieldErrors
	}

	return nil, nil
}

func (v *Validator) validateLoginRequest(req *dto.LoginRequest) error {

	return validation.ValidateStruct(
		req,
		validation.Field(&req.PhoneNumber, validation.Required, validation.Length(10, 13),
			validation.By(v.checkPhoneNumberUniqueness())),

		validation.Field(&req.Password, validation.Required, validation.Length(8, 50)),
	)
}
