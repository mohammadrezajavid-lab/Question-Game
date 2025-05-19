package userservice

import (
	"golang.project/go-fundamentals/gameapp/param"
	"golang.project/go-fundamentals/gameapp/pkg/richerror"
)

func (s *Service) Profile(req *param.ProfileRequest) (*param.ProfileResponse, error) {

	const operation = "service.user.Profile"
	user, err := s.userRepository.GetUserById(req.UserId)
	if err != nil {

		return nil, richerror.NewRichError(operation).WithError(err)
	}

	return param.NewProfileResponse(user.Name), nil
}
