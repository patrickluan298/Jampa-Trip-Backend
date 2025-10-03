package handler

import (
	"net/http"
	"strconv"

	"github.com/jampa_trip/internal/contract"
	"github.com/jampa_trip/internal/model"
	"github.com/jampa_trip/internal/service"
	"github.com/jampa_trip/pkg/database"
	"github.com/jampa_trip/pkg/util"
	"github.com/jampa_trip/pkg/webserver"
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

	serviceCliente := service.ClienteServiceNew(database.DB)
	response, err := serviceCliente.Login(request)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}

// Create - realiza o cadastro de um novo cliente
func (h ClienteHandler) Create(ctx echo.Context) error {

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

	serviceCliente := service.ClienteServiceNew(database.DB)
	response, err := serviceCliente.Create(request)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}

// Update - realiza a atualização de um cliente existente
func (h ClienteHandler) Update(ctx echo.Context) error {

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

	serviceCliente := service.ClienteServiceNew(database.DB)
	response, err := serviceCliente.Update(request)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}

// List - realiza a listagem de todos os clientes
func (h ClienteHandler) List(ctx echo.Context) error {

	Nome := ctx.QueryParam("nome")
	Email := ctx.QueryParam("email")
	CPF := ctx.QueryParam("cpf")
	Telefone := ctx.QueryParam("telefone")

	filtros := &model.Cliente{
		Nome:     Nome,
		Email:    Email,
		CPF:      CPF,
		Telefone: Telefone,
	}

	serviceCliente := service.ClienteServiceNew(database.DB)
	response, err := serviceCliente.List(filtros)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}

// Get - realiza a busca de um cliente por ID
func (h ClienteHandler) Get(ctx echo.Context) error {

	ID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return webserver.InvalidIDResponse(ctx, err)
	}

	if ID < 1 {
		return ctx.JSON(http.StatusBadRequest, contract.ResponseJSON{
			StatusCode: http.StatusBadRequest,
			Message:    "ID não pode ser zero ou negativo",
		})
	}

	serviceCliente := service.ClienteServiceNew(database.DB)
	response, err := serviceCliente.Get(ID)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}
