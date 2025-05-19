package entity

type Functor string

const (
	RoleFunctorType = "role"
	UserFunctorType = "user"
)

// AccessControl only keeps allowed permissions not Denied
type AccessControl struct {
	Id           uint
	FunctorType  Functor
	FunctorId    uint // if FunctorType == role --> FunctorId = 1(UserRole) or 2(AdminRole)
	PermissionId uint
}

func NewAccessControl() *AccessControl {
	return &AccessControl{
		Id:           0,
		FunctorType:  "",
		FunctorId:    0,
		PermissionId: 0,
	}
}
