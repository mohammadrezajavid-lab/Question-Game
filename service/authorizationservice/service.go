package authorizationservice

import (
	"golang.project/go-fundamentals/gameapp/entity"
	"golang.project/go-fundamentals/gameapp/pkg/richerror"
)

type Repository interface {
	GetUserPermissionsTitle(userId uint, role entity.Role) ([]entity.PermissionTitle, error)
}

type Service struct {
	accessControlRepo Repository
}

func NewService(accessControlRepo Repository) *Service {
	return &Service{accessControlRepo: accessControlRepo}
}

func (s *Service) CheckAccess(userId uint, role entity.Role, permissions ...entity.PermissionTitle) (bool, error) {

	const operation = "authorizationService.CheckAccess"

	permissionTitles, gErr := s.accessControlRepo.GetUserPermissionsTitle(userId, role)
	if gErr != nil {
		return false, richerror.NewRichError(operation).WithError(gErr)
	}

	for _, perm := range permissions {
		for _, permTitle := range permissionTitles {
			if perm == permTitle {
				return true, nil
			}
		}
	}

	return false, nil
}
