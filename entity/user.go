package entity

import "fmt"

type User struct {
	ID             uint   `json:"-"`
	Name           string `json:"name"`
	PhoneNumber    string `json:"phone_number"`
	HashedPassword string `json:"-"`
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
