package uservalidator

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation"
	"golang.project/go-fundamentals/gameapp/param/userparam"
	"golang.project/go-fundamentals/gameapp/pkg/errormessage"
	"golang.project/go-fundamentals/gameapp/pkg/richerror"
	"regexp"
)

func (v *Validator) ValidateRegisterRequest(req *userparam.RegisterRequest) (map[string]string, error) {

	const operation = "uservalidator.ValidateRegisterRequest"

	if err := v.validateRegisterRequest(req); err != nil {
		if internalErr, ok := err.(validation.InternalError); ok {

			return nil, richerror.NewRichError(operation).
				WithError(internalErr).
				WithMessage(errormessage.ErrorMsgUnexpected).
				WithKind(richerror.KindUnexpected)
		} else {

			fieldErrors := make(map[string]string)

			valueErr, isOk := err.(validation.Errors)
			if isOk {
				for key, value := range valueErr {
					fieldErrors[key] = value.Error()
				}
			}

			return fieldErrors, richerror.NewRichError(operation).
				WithError(err).
				WithMessage(errormessage.ErrorMsgInvalidRequest).
				WithKind(richerror.KindInvalid)
		}
	}

	return nil, nil
}

func (v *Validator) validateRegisterRequest(req *userparam.RegisterRequest) error {

	return validation.ValidateStruct(
		req,
		validation.Field(&req.Name, validation.Required, validation.Length(3, 50),
			validation.Match(regexp.MustCompile(`^[A-Za-z]+( [A-Za-z]+)*$`))),

		validation.Field(&req.PhoneNumber, validation.Required, validation.Length(10, 13),
			validation.By(v.checkPhoneNumberUniqueness()), validation.By(checkPhoneNumberRegex())),

		validation.Field(&req.Password, validation.Required, validation.Length(8, 50),
			validation.By(checkPasswordRegex())),
	)

}

func (v *Validator) checkPhoneNumberUniqueness() validation.RuleFunc {

	return func(value interface{}) error {

		phoneNumber, ok := value.(string)
		if !ok {
			return errors.New(errormessage.ErrorMsgInvalidPhoneType)
		}

		if isUniq, err := v.repository.IsPhoneNumberUniq(phoneNumber); err != nil || !isUniq {
			if err != nil {

				return validation.NewInternalError(err)
			}

			return errors.New(errormessage.ErrorMsgPhoneNotUniq)
		}

		return nil
	}
}
