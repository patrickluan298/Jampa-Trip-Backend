package handler

import (
	"net/http"
	"strconv"

	"github.com/jampa_trip/internal/contract"
	"github.com/jampa_trip/internal/service"
	"github.com/jampa_trip/pkg/database"
	"github.com/jampa_trip/pkg/middleware"
	"github.com/jampa_trip/pkg/util"
	"github.com/jampa_trip/pkg/webserver"
	"github.com/labstack/echo/v4"
)

type TourHandler struct{}

// Create - cria um novo passeio
func (h TourHandler) Create(ctx echo.Context) error {
	request := &contract.CreateTourRequest{}

	if err := ctx.Bind(request); err != nil {
		if erro := util.ValidateBodyType(err); erro != nil {
			return webserver.ErrorResponse(ctx, erro)
		}
		return webserver.BadJSONResponse(ctx, err)
	}

	if err := request.Validate(); err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	userType := middleware.GetUserType(ctx)
	if userType != "company" {
		return webserver.ErrorResponse(ctx, util.WrapError("Apenas empresas podem criar passeios", nil, http.StatusForbidden))
	}
	companyID := middleware.GetUserID(ctx)

	serviceTour := service.TourServiceNew(database.DB)
	response, err := serviceTour.Create(request, companyID)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusCreated, response)
}

// Update - atualiza um passeio existente
func (h TourHandler) Update(ctx echo.Context) error {
	ID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return webserver.InvalidIDResponse(ctx, err)
	}

	if ID < 1 {
		return webserver.ErrorResponse(ctx, util.WrapError("ID não pode ser zero ou negativo", nil, http.StatusBadRequest))
	}

	request := &contract.UpdateTourRequest{}

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

	userType := middleware.GetUserType(ctx)
	if userType != "company" {
		return webserver.ErrorResponse(ctx, util.WrapError("Apenas empresas podem atualizar passeios", nil, http.StatusForbidden))
	}
	companyID := middleware.GetUserID(ctx)

	serviceTour := service.TourServiceNew(database.DB)
	response, err := serviceTour.Update(request, companyID)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}

// List - lista todos os passeios
func (h TourHandler) List(ctx echo.Context) error {

	search := ctx.QueryParam("search")
	pageStr := ctx.QueryParam("page")
	limitStr := ctx.QueryParam("limit")

	page, limit := util.ParseQueryParams(pageStr, limitStr)

	request := &contract.ListToursRequest{
		Search: search,
		Page:   page,
		Limit:  limit,
	}

	serviceTour := service.TourServiceNew(database.DB)
	response, err := serviceTour.List(request)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}

// GetMyTours - lista passeios da empresa
func (h TourHandler) GetMyTours(ctx echo.Context) error {

	pageStr := ctx.QueryParam("page")
	limitStr := ctx.QueryParam("limit")

	page, limit := util.ParseQueryParams(pageStr, limitStr)

	userType := middleware.GetUserType(ctx)
	if userType != "company" {
		return webserver.ErrorResponse(ctx, util.WrapError("Apenas empresas podem acessar esta rota", nil, http.StatusForbidden))
	}
	companyID := middleware.GetUserID(ctx)

	serviceTour := service.TourServiceNew(database.DB)
	response, err := serviceTour.GetMyTours(companyID, page, limit)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}

// Delete - deleta um passeio
func (h TourHandler) Delete(ctx echo.Context) error {

	ID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return webserver.InvalidIDResponse(ctx, err)
	}

	if ID < 1 {
		return webserver.ErrorResponse(ctx, util.WrapError("ID não pode ser zero ou negativo", nil, http.StatusBadRequest))
	}

	userType := middleware.GetUserType(ctx)
	if userType != "company" {
		return webserver.ErrorResponse(ctx, util.WrapError("Apenas empresas podem deletar passeios", nil, http.StatusForbidden))
	}
	companyID := middleware.GetUserID(ctx)

	serviceTour := service.TourServiceNew(database.DB)
	response, err := serviceTour.Delete(ID, companyID)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}
