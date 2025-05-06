package entity

import "fmt"

type User struct {
	ID             uint
	Name           string
	PhoneNumber    string
	HashedPassword string
}

func NewUser(name string, phoneNumber string, password string) *User {
	return &User{
		ID:             0,
		Name:           name,
		PhoneNumber:    phoneNumber,
		HashedPassword: password,
	}
}

func (u *User) String() string {
	return fmt.Sprintf(`{"id": "%d", "name": "%s", "phone_number": "%s"}`, u.ID, u.Name, u.PhoneNumber)
}
