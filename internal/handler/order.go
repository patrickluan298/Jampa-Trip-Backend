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

type OrderHandler struct{}

// Create - cria uma nova ordem
func (h OrderHandler) Create(ctx echo.Context) error {
	request := &contract.CreateOrderRequest{}

	if err := ctx.Bind(request); err != nil {
		if erro := util.ValidateBodyType(err); erro != nil {
			return webserver.ErrorResponse(ctx, erro)
		}
		return webserver.BadJSONResponse(ctx, err)
	}

	if err := request.Validate(); err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	orderService := service.OrderServiceNew(database.DB)
	response, err := orderService.Create(ctx.Request().Context(), request)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusCreated, response)
}

// Get - obtém uma ordem por ID
func (h OrderHandler) Get(ctx echo.Context) error {
	orderID := ctx.Param("id")
	if orderID == "" {
		return webserver.ErrorResponse(ctx, util.WrapError("ID da ordem é obrigatório", nil, http.StatusBadRequest))
	}

	orderService := service.OrderServiceNew(database.DB)
	response, err := orderService.Get(ctx.Request().Context(), orderID)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}

// List - lista ordens com filtros
func (h OrderHandler) List(ctx echo.Context) error {
	request := &contract.ListOrdersRequest{}

	if clienteIDStr := ctx.QueryParam("cliente_id"); clienteIDStr != "" {
		if clienteID, err := strconv.Atoi(clienteIDStr); err == nil {
			request.ClienteID = clienteID
		}
	}

	if empresaIDStr := ctx.QueryParam("empresa_id"); empresaIDStr != "" {
		if empresaID, err := strconv.Atoi(empresaIDStr); err == nil {
			request.EmpresaID = empresaID
		}
	}

	request.ExternalReference = ctx.QueryParam("external_reference")
	request.Status = ctx.QueryParam("status")

	if offsetStr := ctx.QueryParam("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			request.Offset = offset
		}
	}

	if limitStr := ctx.QueryParam("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			request.Limit = limit
		}
	}

	if err := request.Validate(); err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	orderService := service.OrderServiceNew(database.DB)
	response, err := orderService.List(ctx.Request().Context(), request)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}

// Capture - captura uma ordem autorizada
func (h OrderHandler) Capture(ctx echo.Context) error {
	orderID := ctx.Param("id")
	if orderID == "" {
		return webserver.ErrorResponse(ctx, util.WrapError("ID da ordem é obrigatório", nil, http.StatusBadRequest))
	}

	request := &contract.CaptureOrderRequest{}
	if err := ctx.Bind(request); err != nil {
		request = &contract.CaptureOrderRequest{}
	}

	if err := request.Validate(); err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	orderService := service.OrderServiceNew(database.DB)
	response, err := orderService.Capture(ctx.Request().Context(), orderID, request)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}

// Cancel - cancela uma ordem
func (h OrderHandler) Cancel(ctx echo.Context) error {
	orderID := ctx.Param("id")
	if orderID == "" {
		return webserver.ErrorResponse(ctx, util.WrapError("ID da ordem é obrigatório", nil, http.StatusBadRequest))
	}

	orderService := service.OrderServiceNew(database.DB)
	response, err := orderService.Cancel(ctx.Request().Context(), orderID)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}

// Refund - reembolsa uma ordem
func (h OrderHandler) Refund(ctx echo.Context) error {
	orderID := ctx.Param("id")
	if orderID == "" {
		return webserver.ErrorResponse(ctx, util.WrapError("ID da ordem é obrigatório", nil, http.StatusBadRequest))
	}

	request := &contract.RefundOrderRequest{}
	if err := ctx.Bind(request); err != nil {
		request = &contract.RefundOrderRequest{}
	}

	if err := request.Validate(); err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	orderService := service.OrderServiceNew(database.DB)
	response, err := orderService.Refund(ctx.Request().Context(), orderID, request)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}

// AddTransaction - adiciona uma transação à ordem
func (h OrderHandler) AddTransaction(ctx echo.Context) error {
	orderID := ctx.Param("id")
	if orderID == "" {
		return webserver.ErrorResponse(ctx, util.WrapError("ID da ordem é obrigatório", nil, http.StatusBadRequest))
	}

	request := &contract.AddTransactionRequest{}
	if err := ctx.Bind(request); err != nil {
		if erro := util.ValidateBodyType(err); erro != nil {
			return webserver.ErrorResponse(ctx, erro)
		}
		return webserver.BadJSONResponse(ctx, err)
	}

	if err := request.Validate(); err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	orderService := service.OrderServiceNew(database.DB)
	response, err := orderService.AddTransaction(ctx.Request().Context(), orderID, request)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusCreated, response)
}

// UpdateTransaction - atualiza uma transação existente
func (h OrderHandler) UpdateTransaction(ctx echo.Context) error {
	orderID := ctx.Param("id")
	if orderID == "" {
		return webserver.ErrorResponse(ctx, util.WrapError("ID da ordem é obrigatório", nil, http.StatusBadRequest))
	}

	transactionIDStr := ctx.Param("transaction_id")
	if transactionIDStr == "" {
		return webserver.ErrorResponse(ctx, util.WrapError("ID da transação é obrigatório", nil, http.StatusBadRequest))
	}

	transactionID, err := strconv.Atoi(transactionIDStr)
	if err != nil {
		return webserver.ErrorResponse(ctx, util.WrapError("ID da transação inválido", err, http.StatusBadRequest))
	}

	request := &contract.UpdateTransactionRequest{}
	if err := ctx.Bind(request); err != nil {
		if erro := util.ValidateBodyType(err); erro != nil {
			return webserver.ErrorResponse(ctx, erro)
		}
		return webserver.BadJSONResponse(ctx, err)
	}

	if err := request.Validate(); err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	orderService := service.OrderServiceNew(database.DB)
	response, err := orderService.UpdateTransaction(ctx.Request().Context(), orderID, transactionID, request)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}

// DeleteTransaction - deleta uma transação
func (h OrderHandler) DeleteTransaction(ctx echo.Context) error {
	orderID := ctx.Param("id")
	if orderID == "" {
		return webserver.ErrorResponse(ctx, util.WrapError("ID da ordem é obrigatório", nil, http.StatusBadRequest))
	}

	transactionIDStr := ctx.Param("transaction_id")
	if transactionIDStr == "" {
		return webserver.ErrorResponse(ctx, util.WrapError("ID da transação é obrigatório", nil, http.StatusBadRequest))
	}

	transactionID, err := strconv.Atoi(transactionIDStr)
	if err != nil {
		return webserver.ErrorResponse(ctx, util.WrapError("ID da transação inválido", err, http.StatusBadRequest))
	}

	orderService := service.OrderServiceNew(database.DB)
	err = orderService.DeleteTransaction(ctx.Request().Context(), orderID, transactionID)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "Transação deletada com sucesso",
	})
}

// ProcessOrder - processa uma ordem online (pagamento)
func (h OrderHandler) ProcessOrder(ctx echo.Context) error {
	orderID := ctx.Param("id")
	if orderID == "" {
		return webserver.ErrorResponse(ctx, util.WrapError("ID da ordem é obrigatório", nil, http.StatusBadRequest))
	}

	request := &contract.ProcessOrderRequest{}
	if err := ctx.Bind(request); err != nil {
		if erro := util.ValidateBodyType(err); erro != nil {
			return webserver.ErrorResponse(ctx, erro)
		}
		return webserver.BadJSONResponse(ctx, err)
	}

	if err := request.Validate(); err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	orderService := service.OrderServiceNew(database.DB)
	response, err := orderService.ProcessOrder(ctx.Request().Context(), orderID, request)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}
