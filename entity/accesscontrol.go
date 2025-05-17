package entity

// AccessControl only keeps allowed permissions not Denied
type AccessControl struct {
	Id           uint
	FunctorId    uint
	FunctorType  Functor
	ResourceId   uint
	PermissionId uint
}

type Functor string

const (
	RoleFunctorType = "role"
	UserFunctorType = "user"
)
