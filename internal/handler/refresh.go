package handler

import (
	"net/http"

	"github.com/jampa_trip/internal/contract"
	"github.com/jampa_trip/internal/service"
	"github.com/jampa_trip/pkg/util"
	"github.com/jampa_trip/pkg/webserver"
	"github.com/labstack/echo/v4"
)

type RefreshHandler struct{}

// RefreshToken - renova o par de tokens
func (h RefreshHandler) RefreshToken(ctx echo.Context) error {
	request := &contract.RefreshTokenRequest{}

	if err := ctx.Bind(request); err != nil {
		if erro := util.ValidateBodyType(err); erro != nil {
			return webserver.ErrorResponse(ctx, erro)
		}
		return webserver.BadJSONResponse(ctx, err)
	}

	if err := request.Validate(); err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	serviceRefresh := service.RefreshServiceNew()
	response, err := serviceRefresh.RefreshToken(request)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}
