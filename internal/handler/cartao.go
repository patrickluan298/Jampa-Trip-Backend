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

type CardHandler struct{}

// Create - cria um cartão para um cliente
func (h CardHandler) Create(ctx echo.Context) error {

	customerID := ctx.Param("customer_id")
	if customerID == "" {
		return webserver.ErrorResponse(ctx, util.WrapError("customer_id é obrigatório", nil, http.StatusBadRequest))
	}

	request := &contract.CreateCartaoRequest{}

	if err := ctx.Bind(request); err != nil {
		if erro := util.ValidateBodyType(err); erro != nil {
			return webserver.ErrorResponse(ctx, erro)
		}
		return webserver.BadJSONResponse(ctx, err)
	}

	if err := request.Validate(); err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	serviceCartao := service.CartaoServiceNew(database.DB)
	response, err := serviceCartao.Create(ctx.Request().Context(), customerID, request)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusCreated, response)
}

// List - lista os cartões de um cliente
func (h CardHandler) List(ctx echo.Context) error {
	customerID := ctx.Param("customer_id")
	if customerID == "" {
		return webserver.ErrorResponse(ctx, util.WrapError("customer_id é obrigatório", nil, http.StatusBadRequest))
	}

	serviceCartao := service.CartaoServiceNew(database.DB)
	response, err := serviceCartao.List(ctx.Request().Context(), customerID)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}

// Get - obtém um cartão específico de um cliente
func (h CardHandler) Get(ctx echo.Context) error {
	customerID := ctx.Param("customer_id")
	if customerID == "" {
		return webserver.ErrorResponse(ctx, util.WrapError("customer_id é obrigatório", nil, http.StatusBadRequest))
	}

	cardID := ctx.Param("card_id")
	if cardID == "" {
		return webserver.ErrorResponse(ctx, util.WrapError("card_id é obrigatório", nil, http.StatusBadRequest))
	}

	serviceCartao := service.CartaoServiceNew(database.DB)
	response, err := serviceCartao.Get(ctx.Request().Context(), customerID, cardID)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}

// Update - atualiza um cartão de um cliente
func (h CardHandler) Update(ctx echo.Context) error {
	customerID := ctx.Param("customer_id")
	if customerID == "" {
		return webserver.ErrorResponse(ctx, util.WrapError("customer_id é obrigatório", nil, http.StatusBadRequest))
	}

	cardID := ctx.Param("card_id")
	if cardID == "" {
		return webserver.ErrorResponse(ctx, util.WrapError("card_id é obrigatório", nil, http.StatusBadRequest))
	}

	request := &contract.UpdateCartaoRequest{}

	if err := ctx.Bind(request); err != nil {
		if erro := util.ValidateBodyType(err); erro != nil {
			return webserver.ErrorResponse(ctx, erro)
		}
		return webserver.BadJSONResponse(ctx, err)
	}

	if err := request.Validate(); err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	serviceCartao := service.CartaoServiceNew(database.DB)
	response, err := serviceCartao.Update(ctx.Request().Context(), customerID, cardID, request)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}

// Delete - exclui um cartão de um cliente
func (h CardHandler) Delete(ctx echo.Context) error {
	customerID := ctx.Param("customer_id")
	if customerID == "" {
		return webserver.ErrorResponse(ctx, util.WrapError("customer_id é obrigatório", nil, http.StatusBadRequest))
	}

	cardID := ctx.Param("card_id")
	if cardID == "" {
		return webserver.ErrorResponse(ctx, util.WrapError("card_id é obrigatório", nil, http.StatusBadRequest))
	}

	serviceCartao := service.CartaoServiceNew(database.DB)
	response, err := serviceCartao.Delete(ctx.Request().Context(), customerID, cardID)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}
