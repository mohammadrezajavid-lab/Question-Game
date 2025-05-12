package main

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"regexp"
)

type Address struct {
	Street string
	City   string
	State  string
	Zip    string
}

func NewAddress(street string, city string, state string, zip string) *Address {
	return &Address{Street: street, City: city, State: state, Zip: zip}
}

func (addr *Address) Validate() error {
	return validation.ValidateStruct(
		addr,
		validation.Field(&addr.Street, validation.Required, validation.Length(5, 50)),
		validation.Field(&addr.City, validation.Required, validation.Length(5, 50)),
		validation.Field(&addr.State, validation.Required, validation.Match(regexp.MustCompile(`^[A-Z]{2}$`))),
		validation.Field(&addr.Zip, validation.Required, validation.Match(regexp.MustCompile(`^[0-9]{5}$`))),
	)
}

type Customer struct {
	Name    string
	Email   string
	Address Address
}

func main() {

	c := Customer{
		Name:  "Qiang Xue",
		Email: "q@d.c",
		Address: Address{
			State: "Virginia",
			Zip:   "12343",
		},
	}

	err := validation.Errors{
		"name":  validation.Validate(c.Name, validation.Required, validation.Length(5, 50)),
		"email": validation.Validate(c.Email, validation.Required, is.Email),
		"zip":   validation.Validate(c.Address.Zip, validation.Required, validation.Match(regexp.MustCompile(`^[0-9]{5}$`))),
	}.Filter()
	fmt.Println(err)
}
