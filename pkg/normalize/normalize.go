package normalize

import (
	"errors"
	"fmt"
	"regexp"
)

type normalize struct {
}


func New() normalize {
	return normalize{}
}

func (n normalize) NormalizePhoneNumber(phoneNumber string) (string, error) {

	phoneNumber = regexp.MustCompile(`[^\d\+]`).ReplaceAllString(phoneNumber, "")

	switch {
	case phoneNumber[:3] == "+98":
		return phoneNumber, nil

	case phoneNumber[:2] == "09":
		return fmt.Sprintf("+98%s", phoneNumber[1:]), nil

	case phoneNumber[:2] == "98":
		return fmt.Sprintf("+%s", phoneNumber), nil

	case phoneNumber[:1] == "9":
		return fmt.Sprintf("+98%s", phoneNumber), nil

	default:
		return "", errors.New("invalid phone number format")
	}
}
