package authhandler

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"golang.project/go-fundamentals/gameapp/logger"
	"golang.project/go-fundamentals/gameapp/param/authparam"
	"golang.project/go-fundamentals/gameapp/pkg/errormessage"
	"golang.project/go-fundamentals/gameapp/pkg/jwt"
	"golang.project/go-fundamentals/gameapp/service/authenticationservice"
	"net/http"
)

type AuthHandler struct {
	authenticationSvc authenticationservice.Service
	jwt               *jwt.JWT
}

func New(authenticationSvc authenticationservice.Service, jwt *jwt.JWT) AuthHandler {
	return AuthHandler{authenticationSvc: authenticationSvc, jwt: jwt}
}

func (ah *AuthHandler) SetRoute(e *echo.Echo) {
	authGroup := e.Group("/auth/")
	authGroup.POST("refresh", ah.refreshTokenHandler)
}

func (ah *AuthHandler) refreshTokenHandler(ctx echo.Context) error {
	// TODO - add log and metrics

	var request = authparam.RefreshTokenRequest{}
	if err := ctx.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims, pErr := ah.jwt.ParseJWT(request.RefreshToken)
	if pErr != nil {
		return ctx.JSON(http.StatusUnauthorized, echo.Map{
			"message": errormessage.ErrorMsgInvalidRefreshToken,
		})
	}

	if claims.Subject != ah.jwt.Config.RefreshSubject {
		return ctx.JSON(http.StatusBadRequest, echo.Map{
			"message": errormessage.ErrorMsgInvalidRefreshToken,
		})
	}

	tokens, tErr := ah.authenticationSvc.CreateTokens(&authparam.ClaimRefreshTokenRequest{
		UserId:   claims.UserId,
		UserRole: claims.Role,
	})

	if tErr != nil {
		logger.Warn(tErr, fmt.Sprintf("can't create tokens for user_id: %d", claims.UserId))

		return ctx.JSON(http.StatusUnauthorized, echo.Map{
			"message": "can't create tokens",
		})
	}

	return ctx.JSON(http.StatusOK, echo.Map{
		"tokens": tokens,
	})
}
