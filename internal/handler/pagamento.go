package handler

import (
	"net/http"
	"strconv"

	"github.com/jampa_trip/internal/contract"
	"github.com/jampa_trip/internal/service"
	"github.com/jampa_trip/pkg/database"
	"github.com/jampa_trip/pkg/util"
	"github.com/jampa_trip/pkg/webserver"
	"github.com/labstack/echo/v4"
)

type PaymentHandler struct{}

// CreateCreditCardPayment - cria um pagamento com cartão de crédito
func (h PaymentHandler) CreateCreditCardPayment(ctx echo.Context) error {

	request := &contract.CreateCreditCardPaymentRequest{}

	if err := ctx.Bind(request); err != nil {
		if erro := util.ValidateBodyType(err); erro != nil {
			return webserver.ErrorResponse(ctx, erro)
		}
		return webserver.BadJSONResponse(ctx, err)
	}

	if err := request.Validate(); err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	servicePagamento := service.PagamentoServiceNew(database.DB)
	response, err := servicePagamento.CreateCreditCardPayment(ctx.Request().Context(), request)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusCreated, response)
}

// CreateDebitCardPayment - cria um pagamento com cartão de débito
func (h PaymentHandler) CreateDebitCardPayment(ctx echo.Context) error {

	request := &contract.CreateDebitCardPaymentRequest{}

	if err := ctx.Bind(request); err != nil {
		if erro := util.ValidateBodyType(err); erro != nil {
			return webserver.ErrorResponse(ctx, erro)
		}
		return webserver.BadJSONResponse(ctx, err)
	}

	if err := request.Validate(); err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	servicePagamento := service.PagamentoServiceNew(database.DB)
	response, err := servicePagamento.CreateDebitCardPayment(ctx.Request().Context(), request)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusCreated, response)
}

// CreatePIXPayment - cria um pagamento com PIX
func (h PaymentHandler) CreatePIXPayment(ctx echo.Context) error {

	request := &contract.CreatePIXPaymentRequest{}

	if err := ctx.Bind(request); err != nil {
		if erro := util.ValidateBodyType(err); erro != nil {
			return webserver.ErrorResponse(ctx, erro)
		}
		return webserver.BadJSONResponse(ctx, err)
	}

	if err := request.Validate(); err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	servicePagamento := service.PagamentoServiceNew(database.DB)
	response, err := servicePagamento.CreatePIXPayment(ctx.Request().Context(), request)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusCreated, response)
}

// List - busca pagamentos com filtros
func (h PaymentHandler) List(ctx echo.Context) error {

	request := &contract.ListPaymentsRequest{}

	if err := ctx.Bind(request); err != nil {
		if erro := util.ValidateBodyType(err); erro != nil {
			return webserver.ErrorResponse(ctx, erro)
		}
		return webserver.BadJSONResponse(ctx, err)
	}

	if err := request.Validate(); err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	servicePagamento := service.PagamentoServiceNew(database.DB)
	response, err := servicePagamento.List(ctx.Request().Context(), request)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}

// Get - obtém um pagamento por ID
func (h PaymentHandler) Get(ctx echo.Context) error {

	paymentIDStr := ctx.Param("id")
	paymentID, err := strconv.ParseInt(paymentIDStr, 10, 64)
	if err != nil {
		return webserver.ErrorResponse(ctx, util.WrapError("ID do pagamento inválido", err, http.StatusBadRequest))
	}

	servicePagamento := service.PagamentoServiceNew(database.DB)
	response, err := servicePagamento.Get(ctx.Request().Context(), paymentID)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}

// Update - atualiza um pagamento
func (h PaymentHandler) Update(ctx echo.Context) error {

	request := &contract.UpdatePaymentRequest{}

	if err := ctx.Bind(request); err != nil {
		if erro := util.ValidateBodyType(err); erro != nil {
			return webserver.ErrorResponse(ctx, erro)
		}
		return webserver.BadJSONResponse(ctx, err)
	}

	if err := request.Validate(); err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	servicePagamento := service.PagamentoServiceNew(database.DB)
	response, err := servicePagamento.Update(ctx.Request().Context(), request)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}
