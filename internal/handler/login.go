package handler

import (
	"net/http"

	"github.com/jampa_trip/internal/contract"
	"github.com/jampa_trip/internal/service"
	"github.com/jampa_trip/pkg/database"
	"github.com/jampa_trip/pkg/util"
	"github.com/jampa_trip/pkg/webserver"
	"github.com/labstack/echo/v4"
)

type LoginHandler struct{}

// Login - realiza o login do usuário como empresa ou cliente
func (h LoginHandler) Login(ctx echo.Context) error {

	request := &contract.LoginRequest{}

	if err := ctx.Bind(request); err != nil {
		if erro := util.ValidateBodyType(err); erro != nil {
			return webserver.ErrorResponse(ctx, erro)
		}
		return webserver.BadJSONResponse(ctx, err)
	}

	if err := request.Validate(); err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	serviceLogin := service.LoginServiceNew(database.DB)
	response, err := serviceLogin.Login(request)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}
