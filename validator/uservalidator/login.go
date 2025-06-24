package uservalidator

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation"
	"golang.project/go-fundamentals/gameapp/param/userparam"
	"golang.project/go-fundamentals/gameapp/pkg/errormessage"
	"golang.project/go-fundamentals/gameapp/pkg/richerror"
)

func (v *Validator) ValidateLoginRequest(req *userparam.LoginRequest) (map[string]string, error) {

	const operation = "uservalidator.ValidateLoginRequest"

	if err := v.validateLoginRequest(req); err != nil {

		fieldErrors := make(map[string]string)

		valueErr, ok := err.(validation.Errors)
		if ok {
			for key, value := range valueErr {
				fieldErrors[key] = value.Error()
			}
		}

		return fieldErrors, richerror.NewRichError(operation).
			WithError(err).
			WithMessage(errormessage.ErrorMsgInvalidRequest).
			WithKind(richerror.KindInvalid)
	}

	return nil, nil
}

func (v *Validator) validateLoginRequest(req *userparam.LoginRequest) error {

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
			return errors.New(errormessage.ErrorMsgInvalidPhoneType)
		}

		if _, err := v.repository.GetUserByPhoneNumber(phoneNumber); err != nil {
			return errors.New(errormessage.ErrorMsgNotExistPhoneNumber)
		}

		return nil
	}
}
