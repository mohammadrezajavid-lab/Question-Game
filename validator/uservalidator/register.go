package uservalidator

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation"
	"golang.project/go-fundamentals/gameapp/dto"
	"golang.project/go-fundamentals/gameapp/pkg/richerror"
	"regexp"
)

func (v *Validator) ValidateRegisterRequest(req *dto.RegisterRequest) (error, map[string]string) {

	const operation = "uservalidator.ValidateRegisterRequest"

	if err := v.validateRegisterRequest(req); err != nil {
		if internalErr, ok := err.(validation.InternalError); ok {

			return richerror.NewRichError(operation).
				WithError(internalErr).
				WithMessage("unexpected error: Try again later.").
				WithKind(richerror.KindUnexpected), nil
		} else {

			fieldErrors := make(map[string]string)

			valueErr, isOk := err.(validation.Errors)
			if isOk {
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
	}

	return nil, nil
}

func (v *Validator) validateRegisterRequest(req *dto.RegisterRequest) error {

	return validation.ValidateStruct(
		req,
		validation.Field(&req.Name, validation.Required, validation.Length(3, 50),
			validation.Match(regexp.MustCompile(`^[A-Za-z]+( [A-Za-z]+)*$`))),

		validation.Field(&req.PhoneNumber, validation.Required, validation.Length(10, 13),
			validation.By(v.checkPhoneNumberUniqueness()), validation.By(checkPasswordRegex())),

		validation.Field(&req.Password, validation.Required, validation.Length(8, 50),
			validation.By(checkPasswordRegex())),
	)

}

func (v *Validator) checkPhoneNumberUniqueness() validation.RuleFunc {

	return func(value interface{}) error {

		phoneNumber, ok := value.(string)
		if !ok {
			return errors.New("invalid phone number type")
		}

		if isUniq, err := v.repository.IsPhoneNumberUniq(phoneNumber); err != nil || !isUniq {
			if err != nil {

				return validation.NewInternalError(err)
			}

			return errors.New("phone number is not uniq")
		}

		return nil
	}
}
