package pkg

import "strconv"

func IsPhoneNumberValid(phoneNumber string) bool {
	if len(phoneNumber) != 11 {
		return false
	}

	if _, err := strconv.Atoi(phoneNumber[2:]); err != nil {
		return false
	}

	// TODO Check PhoneNumber
	if phoneNumber[:2] != "09" {
		return false
	}

	return true
}

func IsNameValid(name string) bool {
	if len(name) < 3 {

		return false
	}

	return true
}
