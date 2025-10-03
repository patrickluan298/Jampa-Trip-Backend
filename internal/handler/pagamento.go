package handler

import (
	"net/http"
	"strconv"

	"github.com/jampa_trip/internal"
	"github.com/jampa_trip/internal/contract"
	"github.com/jampa_trip/internal/service"
	"github.com/jampa_trip/pkg/util"
	"github.com/jampa_trip/pkg/webserver"
	"github.com/labstack/echo/v4"
)

type PagamentoHandler struct{}

// CriarPagamentoCartaoCredito - cria um pagamento com cartão de crédito
func (h PagamentoHandler) CriarPagamentoCartaoCredito(ctx echo.Context) error {

	request := &contract.CriarPagamentoCartaoCreditoRequest{}

	if err := ctx.Bind(request); err != nil {
		if erro := util.ValidarTipoBody(err); erro != nil {
			return webserver.ErrorResponse(ctx, erro)
		}
		return webserver.BadJSONResponse(ctx, err)
	}

	if err := request.Validate(); err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	servicePagamento := service.PagamentoServiceNew(internal.DB)
	response, err := servicePagamento.CriarPagamentoCartaoCredito(ctx.Request().Context(), request)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusCreated, response)
}

// CriarPagamentoCartaoDebito - cria um pagamento com cartão de débito
func (h PagamentoHandler) CriarPagamentoCartaoDebito(ctx echo.Context) error {

	request := &contract.CriarPagamentoCartaoDebitoRequest{}

	if err := ctx.Bind(request); err != nil {
		if erro := util.ValidarTipoBody(err); erro != nil {
			return webserver.ErrorResponse(ctx, erro)
		}
		return webserver.BadJSONResponse(ctx, err)
	}

	if err := request.Validate(); err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	servicePagamento := service.PagamentoServiceNew(internal.DB)
	response, err := servicePagamento.CriarPagamentoCartaoDebito(ctx.Request().Context(), request)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusCreated, response)
}

// CriarPagamentoPIX - cria um pagamento com PIX
func (h PagamentoHandler) CriarPagamentoPIX(ctx echo.Context) error {

	request := &contract.CriarPagamentoPIXRequest{}

	if err := ctx.Bind(request); err != nil {
		if erro := util.ValidarTipoBody(err); erro != nil {
			return webserver.ErrorResponse(ctx, erro)
		}
		return webserver.BadJSONResponse(ctx, err)
	}

	if err := request.Validate(); err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	servicePagamento := service.PagamentoServiceNew(internal.DB)
	response, err := servicePagamento.CriarPagamentoPIX(ctx.Request().Context(), request)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusCreated, response)
}

// BuscarPagamentos - busca pagamentos com filtros
func (h PagamentoHandler) BuscarPagamentos(ctx echo.Context) error {

	request := &contract.BuscarPagamentosRequest{}

	if err := ctx.Bind(request); err != nil {
		if erro := util.ValidarTipoBody(err); erro != nil {
			return webserver.ErrorResponse(ctx, erro)
		}
		return webserver.BadJSONResponse(ctx, err)
	}

	if err := request.Validate(); err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	servicePagamento := service.PagamentoServiceNew(internal.DB)
	response, err := servicePagamento.BuscarPagamentos(ctx.Request().Context(), request)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}

// ObterPagamento - obtém um pagamento por ID
func (h PagamentoHandler) ObterPagamento(ctx echo.Context) error {

	paymentIDStr := ctx.Param("id")
	paymentID, err := strconv.ParseInt(paymentIDStr, 10, 64)
	if err != nil {
		return webserver.ErrorResponse(ctx, util.WrapError("ID do pagamento inválido", err, http.StatusBadRequest))
	}

	servicePagamento := service.PagamentoServiceNew(internal.DB)
	response, err := servicePagamento.ObterPagamentoPorID(ctx.Request().Context(), paymentID)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}

// AtualizarPagamento - atualiza um pagamento
func (h PagamentoHandler) AtualizarPagamento(ctx echo.Context) error {

	request := &contract.AtualizarPagamentoRequest{}

	if err := ctx.Bind(request); err != nil {
		if erro := util.ValidarTipoBody(err); erro != nil {
			return webserver.ErrorResponse(ctx, erro)
		}
		return webserver.BadJSONResponse(ctx, err)
	}

	if err := request.Validate(); err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	servicePagamento := service.PagamentoServiceNew(internal.DB)
	response, err := servicePagamento.AtualizarPagamento(ctx.Request().Context(), request)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}
