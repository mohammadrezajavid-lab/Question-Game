package entity

import "fmt"

type User struct {
	Id             uint
	Name           string
	PhoneNumber    string
	HashedPassword string
	Role           Role
}

func NewUser(name string, phoneNumber string, password string) *User {
	return &User{
		Id:             0,
		Name:           name,
		PhoneNumber:    phoneNumber,
		HashedPassword: password,
		Role:           UserRole,
	}
}

func (u *User) String() string {
	return fmt.Sprintf(
		`{"id": "%d", "name": "%s", "phone_number": "%s", "role": "%s"}`,
		u.Id, u.Name, u.PhoneNumber, u.Role.String())
}
