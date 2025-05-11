package hash

import "golang.org/x/crypto/bcrypt"

func Hash(str string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(str), 14)

	return string(bytes), err
}

func CheckHash(str string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(str))

	return err == nil
}
