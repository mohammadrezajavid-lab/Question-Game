package middleware

import (
	"github.com/labstack/echo/v4"
	"golang.project/go-fundamentals/gameapp/entity"
	"golang.project/go-fundamentals/gameapp/pkg/claim"
	"golang.project/go-fundamentals/gameapp/pkg/errormessage"
	"net/http"
)

func (m *Middleware) AccessCheck(permissions ...entity.PermissionTitle) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			claims := claim.GetClaimsFromEchoContext(c)
			isAccess, cErr := m.authorizationService.CheckAccess(claims.UserId, claims.Role, permissions...)
			if cErr != nil {
				// TODO - log unexpected error
				return c.JSON(http.StatusInternalServerError, echo.Map{
					"message": errormessage.ErrorMsgUnexpected,
				})
			}

			if !isAccess {
				return c.JSON(http.StatusForbidden, echo.Map{
					"message": errormessage.ErrorMsgUserNotAllowed,
				})
			}
			return next(c)
		}
	}
}
