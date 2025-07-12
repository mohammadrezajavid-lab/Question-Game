package matchingvalidator

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation"
	"golang.project/go-fundamentals/gameapp/entity"
	"golang.project/go-fundamentals/gameapp/param/gameparam"
	"golang.project/go-fundamentals/gameapp/pkg/errormessage"
	"golang.project/go-fundamentals/gameapp/pkg/richerror"
)

func (v *Validator) ValidateAddToWaitingListRequest(req *gameparam.AddToWaitingListRequest) (map[string]string, error) {
	const operation = "matchingvalidator.ValidateAddToWaitingListRequest"

	if err := v.validateAddToWaitingListRequest(req); err != nil {
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

func (v *Validator) validateAddToWaitingListRequest(req *gameparam.AddToWaitingListRequest) error {

	return validation.ValidateStruct(req, validation.Field(&req.Category, validation.Required, validation.By(v.checkCategoryIsValid())))
}

func (v *Validator) checkCategoryIsValid() validation.RuleFunc {

	return func(value interface{}) error {

		category, ok := value.(entity.Category)
		if !ok {

			return errors.New(errormessage.ErrorMsgInvalidCategoryType)
		}

		if isValid := category.IsValid(); !isValid {

			return errors.New(errormessage.ErrorMsgInvalidCategory)
		}

		return nil
	}
}
