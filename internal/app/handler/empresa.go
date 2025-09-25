package handler

import (
	"net/http"
	"strconv"

	"github.com/jampa_trip/internal/app"
	"github.com/jampa_trip/internal/app/contract"
	"github.com/jampa_trip/internal/app/service"
	"github.com/jampa_trip/internal/pkg/util"
	"github.com/jampa_trip/internal/pkg/webserver"
	"github.com/labstack/echo/v4"
)

type EmpresaHandler struct{}

// Login - realiza o login de uma empresa
func (h EmpresaHandler) Login(ctx echo.Context) error {

	request := &contract.LoginEmpresaRequest{}

	if err := ctx.Bind(request); err != nil {
		if erro := util.ValidarTipoBody(err); erro != nil {
			return webserver.ErrorResponse(ctx, erro)
		}
		return webserver.BadJSONResponse(ctx, err)
	}

	if err := request.Validate(); err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	serviceEmpresa := service.EmpresaServiceNew(app.DB)
	response, err := serviceEmpresa.Login(request)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}

// Create - realiza o cadastro de uma nova empresa
func (h EmpresaHandler) Create(ctx echo.Context) error {

	request := &contract.CadastrarEmpresaRequest{}

	if err := ctx.Bind(request); err != nil {
		if erro := util.ValidarTipoBody(err); erro != nil {
			return webserver.ErrorResponse(ctx, erro)
		}
		return webserver.BadJSONResponse(ctx, err)
	}

	if err := request.Validate(); err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	serviceEmpresa := service.EmpresaServiceNew(app.DB)
	response, err := serviceEmpresa.Create(request)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}

// Update - realiza a atualização de uma empresa existente
func (h EmpresaHandler) Update(ctx echo.Context) error {

	request := &contract.AtualizarEmpresaRequest{}

	if err := ctx.Bind(request); err != nil {
		if erro := util.ValidarTipoBody(err); erro != nil {
			return webserver.ErrorResponse(ctx, erro)
		}
		return webserver.BadJSONResponse(ctx, err)
	}

	if err := request.Validate(); err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	serviceEmpresa := service.EmpresaServiceNew(app.DB)
	response, err := serviceEmpresa.Update(request)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}

// List - realiza a listagem de todas as empresas
func (h EmpresaHandler) List(ctx echo.Context) error {

	serviceEmpresa := service.EmpresaServiceNew(app.DB)
	response, err := serviceEmpresa.List()
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}

// Get - realiza a busca de uma empresa por ID
func (h EmpresaHandler) Get(ctx echo.Context) error {

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

	serviceEmpresa := service.EmpresaServiceNew(app.DB)
	response, err := serviceEmpresa.Get(ID)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}
