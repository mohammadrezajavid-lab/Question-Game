package entity

type Role uint8

const (
	UserRole = iota + 1
	AdminRole
)

func (r Role) String() string {

	switch r {
	case UserRole:
		return "user"
	case AdminRole:
		return "admin"
	default:
		return ""
	}
}
