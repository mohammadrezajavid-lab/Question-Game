package userservice

import (
	"context"
	"golang.project/go-fundamentals/gameapp/param/userparam"
	"golang.project/go-fundamentals/gameapp/pkg/richerror"
)

func (s *Service) Profile(ctx context.Context, req *userparam.ProfileRequest) (*userparam.ProfileResponse, error) {

	const operation = "service.user.Profile"
	user, err := s.userRepository.GetUserById(ctx, req.UserId)
	if err != nil {

		return nil, richerror.NewRichError(operation).WithError(err)
	}

	return userparam.NewProfileResponse(user.Name), nil
}
