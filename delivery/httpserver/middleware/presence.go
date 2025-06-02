package middleware

import (
	"github.com/labstack/echo/v4"
	"golang.project/go-fundamentals/gameapp/param/presenceparam"
	"golang.project/go-fundamentals/gameapp/pkg/claim"
	"golang.project/go-fundamentals/gameapp/pkg/errormessage"
	"golang.project/go-fundamentals/gameapp/pkg/timestamp"
	"net/http"
)

func (m *Middleware) PresenceUpsert() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims := claim.GetClaimsFromEchoContext(c)
			_, uErr := m.presenceClient.Upsert(c.Request().Context(), presenceparam.NewUpsertPresenceRequest(claims.UserId, timestamp.Now()))
			if uErr != nil {

				return c.JSON(http.StatusInternalServerError, echo.Map{
					"message": errormessage.ErrorMsgUnexpected,
				})
			}

			return next(c)
		}
	}
}
