package parsericherror

import (
	"golang.project/go-fundamentals/gameapp/pkg/richerror"
	"net/http"
)

type HttpRichErrorParse struct {
}

func New() HttpRichErrorParse {
	return HttpRichErrorParse{}
}

// ParseRichError output value: message, httpStatusCode
func (rp HttpRichErrorParse) ParseRichError(err error) (string, int) {

	switch err.(type) {
	case richerror.RichError:

		re := err.(richerror.RichError)
		message := re.GetMessage()
		statusCode := rp.kindToHttpStatusCode(re.GetKind())

		// we should not expose unexpected error message
		if statusCode >= 500 {
			message = "internal server error"
		}

		return message, statusCode
	default:

		return err.Error(), http.StatusBadRequest
	}
}

func (rp HttpRichErrorParse) kindToHttpStatusCode(kind richerror.Kind) int {

	switch kind {
	case richerror.KindInvalid:
		return http.StatusUnprocessableEntity
	case richerror.KindNotFound:
		return http.StatusNotFound
	case richerror.KindUnexpected:
		return http.StatusInternalServerError
	case richerror.KindForbidden:
		return http.StatusForbidden
	default:
		return http.StatusBadRequest
	}
}
