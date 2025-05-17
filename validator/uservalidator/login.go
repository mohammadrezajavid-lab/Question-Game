package uservalidator

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation"
	"golang.project/go-fundamentals/gameapp/param"
	"golang.project/go-fundamentals/gameapp/pkg/richerror"
)

func (v *Validator) ValidateLoginRequest(req *param.LoginRequest) (error, map[string]string) {

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

func (v *Validator) validateLoginRequest(req *param.LoginRequest) error {

	return validation.ValidateStruct(
		req,
		validation.Field(&req.PhoneNumber, validation.Required, validation.Length(10, 13),
			validation.By(v.checkPhoneNumberExistence()), validation.By(checkPhoneNumberRegex())),

		validation.Field(&req.Password, validation.Required, validation.Length(8, 50)),
	)
}

func (v *Validator) checkPhoneNumberExistence() validation.RuleFunc {

	return func(value interface{}) error {

		phoneNumber, ok := value.(string)
		if !ok {
			return errors.New("invalid phone number type")
		}

		if _, err := v.repository.GetUserByPhoneNumber(phoneNumber); err != nil {
			return errors.New("phone number does not exits")
		}

		return nil
	}
}
