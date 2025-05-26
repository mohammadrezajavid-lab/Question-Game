package matchingvalidator

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation"
	"golang.project/go-fundamentals/gameapp/entity"
	"golang.project/go-fundamentals/gameapp/param/matchingparam"
	"golang.project/go-fundamentals/gameapp/pkg/errormessage"
	"golang.project/go-fundamentals/gameapp/pkg/richerror"
)

func (v *Validator) ValidateAddToWaitingListRequest(req *matchingparam.AddToWaitingListRequest) (error, map[string]string) {

	const operation = "matchingvalidator.ValidateAddToWaitingListRequest"
	if err := v.validateAddToWaitingListRequest(req); err != nil {
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
					WithMessage(errormessage.ErrorMsgInvalidRequest).
					WithKind(richerror.KindInvalid),
				fieldErrors
		}
	}

	return nil, nil
}

func (v *Validator) validateAddToWaitingListRequest(req *matchingparam.AddToWaitingListRequest) error {

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
