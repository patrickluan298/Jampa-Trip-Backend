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

type ClientHandler struct{}

// Create - realiza o cadastro de um novo cliente
func (h ClientHandler) Create(ctx echo.Context) error {

	request := &contract.CreateClientRequest{}

	if err := ctx.Bind(request); err != nil {
		if erro := util.ValidateBodyType(err); erro != nil {
			return webserver.ErrorResponse(ctx, erro)
		}
		return webserver.BadJSONResponse(ctx, err)
	}

	if err := request.Validate(); err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	serviceCliente := service.ClientServiceNew(database.DB)
	response, err := serviceCliente.Create(request)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}

// Update - realiza a atualização de um cliente existente
func (h ClientHandler) Update(ctx echo.Context) error {

	ID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return webserver.InvalidIDResponse(ctx, err)
	}

	request := &contract.UpdateClientRequest{}

	if err := ctx.Bind(request); err != nil {
		if erro := util.ValidateBodyType(err); erro != nil {
			return webserver.ErrorResponse(ctx, erro)
		}
		return webserver.BadJSONResponse(ctx, err)
	}

	request.ID = ID

	if err := request.Validate(); err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	serviceCliente := service.ClientServiceNew(database.DB)
	response, err := serviceCliente.Update(request)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}

// List - realiza a listagem de todos os clientes
func (h ClientHandler) List(ctx echo.Context) error {

	Name := ctx.QueryParam("name")
	Email := ctx.QueryParam("email")
	CPF := ctx.QueryParam("cpf")
	Phone := ctx.QueryParam("phone")

	filtros := &model.Client{
		Name:  Name,
		Email: Email,
		CPF:   CPF,
		Phone: Phone,
	}

	serviceCliente := service.ClientServiceNew(database.DB)
	response, err := serviceCliente.List(filtros)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}

// Get - realiza a busca de um cliente por ID
func (h ClientHandler) Get(ctx echo.Context) error {

	ID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return webserver.InvalidIDResponse(ctx, err)
	}

	if ID < 1 {
		return webserver.ErrorResponse(ctx, util.WrapError("ID não pode ser zero ou negativo", nil, http.StatusBadRequest))
	}

	serviceCliente := service.ClientServiceNew(database.DB)
	response, err := serviceCliente.Get(ID)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}
