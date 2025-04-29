package entity

import "fmt"

type User struct {
	ID          uint
	Name        string
	PhoneNumber string
}

func NewUser(name string, phoneNumber string) *User {
	return &User{
		ID:          0,
		Name:        name,
		PhoneNumber: phoneNumber,
	}
}

func (u *User) String() string {
	return fmt.Sprintf("ID: %d, Name: %s, PhoneNumber: %s", u.ID, u.Name, u.PhoneNumber)
}
