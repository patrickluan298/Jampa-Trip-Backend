package handler

import (
	"net/http"

	"github.com/jampa_trip/internal/app"
	"github.com/jampa_trip/internal/app/contract"
	"github.com/jampa_trip/internal/app/service"
	"github.com/jampa_trip/internal/pkg/util"
	"github.com/jampa_trip/internal/pkg/webserver"
	"github.com/labstack/echo/v4"
)

type ClienteHandler struct{}

// Login - realiza o login de um cliente
func (h ClienteHandler) Login(ctx echo.Context) error {

	request := &contract.LoginClienteRequest{}

	if err := ctx.Bind(request); err != nil {
		if erro := util.ValidarTipoBody(err); erro != nil {
			return webserver.ErrorResponse(ctx, erro)
		}
		return webserver.BadJSONResponse(ctx, err)
	}

	if err := request.Validate(); err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	serviceCliente := service.ClienteServiceNew(app.DB)
	response, err := serviceCliente.Login(request)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}

// Cadastrar - realiza o cadastro de um novo cliente
func (h ClienteHandler) Cadastrar(ctx echo.Context) error {

	request := &contract.CadastrarClienteRequest{}

	if err := ctx.Bind(request); err != nil {
		if erro := util.ValidarTipoBody(err); erro != nil {
			return webserver.ErrorResponse(ctx, erro)
		}
		return webserver.BadJSONResponse(ctx, err)
	}

	if err := request.Validate(); err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	serviceCliente := service.ClienteServiceNew(app.DB)
	response, err := serviceCliente.Cadastrar(request)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}

// Atualizar - realiza a atualização de um cliente existente
func (h ClienteHandler) Atualizar(ctx echo.Context) error {

	request := &contract.AtualizarClienteRequest{}

	if err := ctx.Bind(request); err != nil {
		if erro := util.ValidarTipoBody(err); erro != nil {
			return webserver.ErrorResponse(ctx, erro)
		}
		return webserver.BadJSONResponse(ctx, err)
	}

	if err := request.Validate(); err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	serviceCliente := service.ClienteServiceNew(app.DB)
	response, err := serviceCliente.Atualizar(request)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}

// Listar - realiza a listagem de todos os clientes
func (h ClienteHandler) Listar(ctx echo.Context) error {

	serviceCliente := service.ClienteServiceNew(app.DB)
	response, err := serviceCliente.Listar()
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}
