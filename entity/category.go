package entity

type Category struct {
	Id          uint
	Name        string
	Description string
}

func NewCategory(id uint, name string, description string) *Category {
	return &Category{
		Id:          id,
		Name:        name,
		Description: description,
	}
}
