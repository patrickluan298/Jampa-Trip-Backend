package handler

import (
	"net/http"
	"strconv"

	"github.com/jampa_trip/internal/app/contract"
	"github.com/jampa_trip/internal/pkg/util"
	"github.com/jampa_trip/internal/pkg/webserver"
	"github.com/labstack/echo/v4"
)

type PagamentoHandler struct{}

// Create - cria um novo pagamento
func (h PagamentoHandler) Create(ctx echo.Context) error {
	request := &contract.CreatePagamentoRequest{}

	if err := ctx.Bind(request); err != nil {
		if erro := util.ValidarTipoBody(err); erro != nil {
			return webserver.ErrorResponse(ctx, erro)
		}
		return webserver.BadJSONResponse(ctx, err)
	}

	if err := request.Validate(); err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	// TODO: Implementar PagamentoService quando estiver disponível
	// pagamentoService := service.PagamentoServiceNew(app.DB)
	// response, err := pagamentoService.Create(request)
	// if err != nil {
	// 	return webserver.ErrorResponse(ctx, err)
	// }

	// Por enquanto, retornar erro de não implementado
	return webserver.ErrorResponse(ctx, util.WrapError("funcionalidade de pagamento não implementada", nil, http.StatusNotImplemented))
}

// Get - busca um pagamento pelo ID
func (h PagamentoHandler) Get(ctx echo.Context) error {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return webserver.ErrorResponse(ctx, util.WrapError("ID inválido", err, http.StatusBadRequest))
	}

	request := &contract.GetPagamentoRequest{ID: id}

	if err := request.Validate(); err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	// TODO: Implementar PagamentoService quando estiver disponível
	// pagamentoService := service.PagamentoServiceNew(app.DB)
	// response, err := pagamentoService.GetByID(request)
	// if err != nil {
	// 	return webserver.ErrorResponse(ctx, err)
	// }

	// Por enquanto, retornar erro de não implementado
	return webserver.ErrorResponse(ctx, util.WrapError("funcionalidade de pagamento não implementada", nil, http.StatusNotImplemented))
}

// List - lista pagamentos
func (h PagamentoHandler) List(ctx echo.Context) error {
	request := &contract.ListPagamentoRequest{}

	// Parse query parameters
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

	request.Status = ctx.QueryParam("status")
	request.MetodoPagamento = ctx.QueryParam("metodo_pagamento")

	if pageStr := ctx.QueryParam("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil {
			request.Page = page
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

	// TODO: Implementar PagamentoService quando estiver disponível
	// pagamentoService := service.PagamentoServiceNew(app.DB)
	// response, err := pagamentoService.List(request)
	// if err != nil {
	// 	return webserver.ErrorResponse(ctx, err)
	// }

	// Por enquanto, retornar erro de não implementado
	return webserver.ErrorResponse(ctx, util.WrapError("funcionalidade de pagamento não implementada", nil, http.StatusNotImplemented))
}

// Update - atualiza um pagamento
func (h PagamentoHandler) Update(ctx echo.Context) error {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return webserver.ErrorResponse(ctx, util.WrapError("ID inválido", err, http.StatusBadRequest))
	}

	request := &contract.UpdatePagamentoRequest{}

	if err := ctx.Bind(request); err != nil {
		if erro := util.ValidarTipoBody(err); erro != nil {
			return webserver.ErrorResponse(ctx, erro)
		}
		return webserver.BadJSONResponse(ctx, err)
	}

	if err := request.Validate(); err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	// TODO: Implementar PagamentoService quando estiver disponível
	// pagamentoService := service.PagamentoServiceNew(app.DB)
	// response, err := pagamentoService.Update(id, request)
	// if err != nil {
	// 	return webserver.ErrorResponse(ctx, err)
	// }

	// Por enquanto, retornar erro de não implementado
	return webserver.ErrorResponse(ctx, util.WrapError("funcionalidade de pagamento não implementada", nil, http.StatusNotImplemented))
}

// Delete - remove um pagamento
func (h PagamentoHandler) Delete(ctx echo.Context) error {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return webserver.ErrorResponse(ctx, util.WrapError("ID inválido", err, http.StatusBadRequest))
	}

	// TODO: Implementar PagamentoService quando estiver disponível
	// pagamentoService := service.PagamentoServiceNew(app.DB)
	// response, err := pagamentoService.Delete(id)
	// if err != nil {
	// 	return webserver.ErrorResponse(ctx, err)
	// }

	// Por enquanto, retornar erro de não implementado
	return webserver.ErrorResponse(ctx, util.WrapError("funcionalidade de pagamento não implementada", nil, http.StatusNotImplemented))
}
