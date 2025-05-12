package normalize

import (
	"errors"
	"regexp"
	"strings"
)

type Normalize struct {
}

func New() Normalize {
	return Normalize{}
}

func (n Normalize) NormalizePhoneNumber(phoneNumber string) (string, error) {

	phoneNumber = regexp.MustCompile(`[^\d\+]`).ReplaceAllString(phoneNumber, "")

	switch {
	case strings.HasPrefix(phoneNumber, "+98"):
		return phoneNumber, nil
	case strings.HasPrefix(phoneNumber, "09"):
		return "+98" + phoneNumber[1:], nil
	case strings.HasPrefix(phoneNumber, "9"):
		return "+98" + phoneNumber, nil
	case strings.HasPrefix(phoneNumber, "098"):
		return "+98" + phoneNumber[3:], nil
	case strings.HasPrefix(phoneNumber, "0098"):
		return "+98" + phoneNumber[4:], nil
	default:
		return "", errors.New("invalid phone number format")
	}
}
