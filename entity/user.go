package entity

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
