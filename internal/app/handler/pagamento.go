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

type PagamentoHandler struct{}

// AutorizarCartao - autoriza um pagamento com cartão de crédito
func (h PagamentoHandler) AutorizarCartao(ctx echo.Context) error {
	request := &contract.AutorizarCartaoRequest{}

	if err := ctx.Bind(request); err != nil {
		if erro := util.ValidarTipoBody(err); erro != nil {
			return webserver.ErrorResponse(ctx, erro)
		}
		return webserver.BadJSONResponse(ctx, err)
	}

	if err := request.Validate(); err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	servicePagamento := service.PagamentoServiceNew(app.DB)
	response, err := servicePagamento.AutorizarCartao(ctx.Request().Context(), request)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}

// CapturarPagamento - captura um pagamento autorizado
// func (h PagamentoHandler) CapturarPagamento(ctx echo.Context) error {
// 	request := &contract.CapturarPagamentoRequest{}

// 	if err := ctx.Bind(request); err != nil {
// 		if erro := util.ValidarTipoBody(err); erro != nil {
// 			return webserver.ErrorResponse(ctx, erro)
// 		}
// 		return webserver.BadJSONResponse(ctx, err)
// 	}

// 	if err := request.Validate(); err != nil {
// 		return webserver.ErrorResponse(ctx, err)
// 	}

// 	servicePagamento := service.PagamentoServiceNew(app.DB)
// 	response, err := servicePagamento.CapturarPagamento(ctx.Request().Context(), request)
// 	if err != nil {
// 		return webserver.ErrorResponse(ctx, err)
// 	}

// 	return ctx.JSON(http.StatusOK, response)
// }

// CancelarPagamento - cancela um pagamento autorizado
// func (h PagamentoHandler) CancelarPagamento(ctx echo.Context) error {
// 	request := &contract.CancelarPagamentoRequest{}

// 	if err := ctx.Bind(request); err != nil {
// 		if erro := util.ValidarTipoBody(err); erro != nil {
// 			return webserver.ErrorResponse(ctx, erro)
// 		}
// 		return webserver.BadJSONResponse(ctx, err)
// 	}

// 	if err := request.Validate(); err != nil {
// 		return webserver.ErrorResponse(ctx, err)
// 	}

// 	servicePagamento := service.PagamentoServiceNew(app.DB)
// 	response, err := servicePagamento.CancelarPagamento(ctx.Request().Context(), request)
// 	if err != nil {
// 		return webserver.ErrorResponse(ctx, err)
// 	}

// 	return ctx.JSON(http.StatusOK, response)
// }

// ReembolsarPagamento - reembolsa um pagamento capturado
// func (h PagamentoHandler) ReembolsarPagamento(ctx echo.Context) error {
// 	request := &contract.ReembolsarPagamentoRequest{}

// 	if err := ctx.Bind(request); err != nil {
// 		if erro := util.ValidarTipoBody(err); erro != nil {
// 			return webserver.ErrorResponse(ctx, erro)
// 		}
// 		return webserver.BadJSONResponse(ctx, err)
// 	}

// 	if err := request.Validate(); err != nil {
// 		return webserver.ErrorResponse(ctx, err)
// 	}

// 	servicePagamento := service.PagamentoServiceNew(app.DB)
// 	response, err := servicePagamento.ReembolsarPagamento(ctx.Request().Context(), request)
// 	if err != nil {
// 		return webserver.ErrorResponse(ctx, err)
// 	}

// 	return ctx.JSON(http.StatusOK, response)
// }

// ObterPagamento - obtém informações de um pagamento pelo ID do Mercado Pago
// func (h PagamentoHandler) ObterPagamento(ctx echo.Context) error {
// 	paymentIDStr := ctx.Param("payment_id")
// 	paymentID, err := strconv.ParseInt(paymentIDStr, 10, 64)
// 	if err != nil {
// 		return webserver.ErrorResponse(ctx, util.WrapError("ID do pagamento inválido", err, http.StatusBadRequest))
// 	}

// 	servicePagamento := service.PagamentoServiceNew(app.DB)
// 	response, err := servicePagamento.ObterPagamento(ctx.Request().Context(), paymentID)
// 	if err != nil {
// 		return webserver.ErrorResponse(ctx, err)
// 	}

// 	return ctx.JSON(http.StatusOK, response)
// }

// ListarMeiosPagamento - lista os meios de pagamento disponíveis
// func (h PagamentoHandler) ListarMeiosPagamento(ctx echo.Context) error {
// 	servicePagamento := service.PagamentoServiceNew(app.DB)
// 	response, err := servicePagamento.ListarMeiosPagamento(ctx.Request().Context())
// 	if err != nil {
// 		return webserver.ErrorResponse(ctx, err)
// 	}

// 	return ctx.JSON(http.StatusOK, response)
// }
