package matchinghandler

import (
	"github.com/labstack/echo/v4"
	"golang.project/go-fundamentals/gameapp/delivery/httpserver/parsericherror"
	"golang.project/go-fundamentals/gameapp/param/matchingparam"
	"golang.project/go-fundamentals/gameapp/pkg/claim"
	"net/http"
)

func (h *MatchingHandler) addToWaitingList(ctx echo.Context) error {

	var request = matchingparam.NewAddToWaitingListRequest()
	if err := ctx.Bind(request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// bind UserId from claims jwt token
	claims := claim.GetClaimsFromEchoContext(ctx)
	request.UserId = claims.UserId

	// validate request
	if fieldErrors, validateErr := h.matchingValidator.ValidateAddToWaitingListRequest(request); validateErr != nil {

		message, statusCode := parsericherror.New().ParseRichError(validateErr)

		return ctx.JSON(statusCode, echo.Map{
			"message": message,
			"errors":  fieldErrors,
		})
	}

	response, aErr := h.matchingService.AddToWaitingList(ctx.Request().Context(), request)
	if aErr != nil {

		message, statusCode := parsericherror.New().ParseRichError(aErr)
		return echo.NewHTTPError(statusCode, echo.Map{
			"message": message,
		})
	}

	return ctx.JSON(http.StatusOK, response)
}
