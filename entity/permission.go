package entity

type PermissionTitle string

const (
	UserListPermission   PermissionTitle = "user-list"
	UserDeletePermission PermissionTitle = "user-delete"
	UserAddPermission    PermissionTitle = "user-add"
	UserEditPermission   PermissionTitle = "user-edit"
)

type Permission struct {
	Id    uint
	Title PermissionTitle // resource-action EX: user-delete
}

func NewPermission() *Permission {
	return &Permission{
		Id:    0,
		Title: "",
	}
}
