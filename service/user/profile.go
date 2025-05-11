package user

import (
	"golang.project/go-fundamentals/gameapp/dto"
	"golang.project/go-fundamentals/gameapp/pkg/richerror"
)

func (s *Service) Profile(req *dto.ProfileRequest) (*dto.ProfileResponse, error) {

	const operation = "service.user.Profile"
	user, err := s.userRepository.GetUserById(req.UserId)
	if err != nil {

		return nil, richerror.NewRichError(operation).WithError(err)
	}

	return dto.NewProfileResponse(user.Name), nil
}
