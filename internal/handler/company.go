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

type CompanyHandler struct{}

// Create - realiza o cadastro de uma nova empresa
func (receiver CompanyHandler) Create(ctx echo.Context) error {

	request := &contract.CreateCompanyRequest{}

	if err := ctx.Bind(request); err != nil {
		if erro := util.ValidateBodyType(err); erro != nil {
			return webserver.ErrorResponse(ctx, erro)
		}
		return webserver.BadJSONResponse(ctx, err)
	}

	if err := request.Validate(); err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	serviceEmpresa := service.CompanyServiceNew(database.DB)
	response, err := serviceEmpresa.Create(request)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}

// Update - realiza a atualização de uma empresa existente
func (receiver CompanyHandler) Update(ctx echo.Context) error {

	request := &contract.UpdateCompanyRequest{}

	if err := ctx.Bind(request); err != nil {
		if erro := util.ValidateBodyType(err); erro != nil {
			return webserver.ErrorResponse(ctx, erro)
		}
		return webserver.BadJSONResponse(ctx, err)
	}

	if err := request.Validate(); err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	serviceEmpresa := service.CompanyServiceNew(database.DB)
	response, err := serviceEmpresa.Update(request)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}

// List - realiza a listagem de todas as empresas
func (receiver CompanyHandler) List(ctx echo.Context) error {

	Name := ctx.QueryParam("name")
	Email := ctx.QueryParam("email")
	CNPJ := ctx.QueryParam("cnpj")
	Phone := ctx.QueryParam("phone")
	Address := ctx.QueryParam("address")

	filtros := &model.Company{
		Name:    Name,
		Email:   Email,
		CNPJ:    CNPJ,
		Phone:   Phone,
		Address: Address,
	}

	serviceEmpresa := service.CompanyServiceNew(database.DB)
	response, err := serviceEmpresa.List(filtros)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}

// Get - realiza a busca de uma empresa por ID
func (receiver CompanyHandler) Get(ctx echo.Context) error {

	ID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return webserver.InvalidIDResponse(ctx, err)
	}

	if ID < 1 {
		return webserver.ErrorResponse(ctx, util.WrapError("ID não pode ser zero ou negativo", nil, http.StatusBadRequest))
	}

	serviceEmpresa := service.CompanyServiceNew(database.DB)
	response, err := serviceEmpresa.Get(ID)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}
