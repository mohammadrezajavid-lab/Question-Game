package entity

type Role uint8

const (
	UserRole Role = iota + 1
	AdminRole
)

const (
	UserRoleStr  = "user"
	AdminRoleStr = "admin"
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

func MapToRoleEntity(roleStr string) Role {
	switch roleStr {
	case UserRoleStr:
		return UserRole
	case AdminRoleStr:
		return AdminRole
	default:
		return 0
	}
}
