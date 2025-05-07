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
		return re.GetMessage(), rp.kindToHttpStatusCode(re.GetKind())
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
