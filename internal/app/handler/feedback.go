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

type FeedbackHandler struct{}

// Create - cria um novo feedback
func (h FeedbackHandler) Create(ctx echo.Context) error {
	request := &contract.CreateFeedbackRequest{}

	if err := ctx.Bind(request); err != nil {
		if erro := util.ValidarTipoBody(err); erro != nil {
			return webserver.ErrorResponse(ctx, erro)
		}
		return webserver.BadJSONResponse(ctx, err)
	}

	if err := request.Validate(); err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	feedbackService := service.FeedbackServiceNew(app.DB)
	response, err := feedbackService.Create(request)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return webserver.SuccessResponse(ctx, response)
}

// Get - busca um feedback pelo ID
func (h FeedbackHandler) Get(ctx echo.Context) error {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return webserver.ErrorResponse(ctx, util.WrapError("ID inválido", err, http.StatusBadRequest))
	}

	request := &contract.GetFeedbackRequest{ID: id}

	if err := request.Validate(); err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	feedbackService := service.FeedbackServiceNew(app.DB)
	response, err := feedbackService.GetByID(request)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return webserver.SuccessResponse(ctx, response)
}

// List - lista feedbacks
func (h FeedbackHandler) List(ctx echo.Context) error {
	request := &contract.ListFeedbackRequest{}

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

	if notaStr := ctx.QueryParam("nota"); notaStr != "" {
		if nota, err := strconv.Atoi(notaStr); err == nil {
			request.Nota = nota
		}
	}

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

	feedbackService := service.FeedbackServiceNew(app.DB)
	response, err := feedbackService.List(request)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return webserver.SuccessResponse(ctx, response)
}

// Update - atualiza um feedback
func (h FeedbackHandler) Update(ctx echo.Context) error {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return webserver.ErrorResponse(ctx, util.WrapError("ID inválido", err, http.StatusBadRequest))
	}

	request := &contract.UpdateFeedbackRequest{}

	if err := ctx.Bind(request); err != nil {
		if erro := util.ValidarTipoBody(err); erro != nil {
			return webserver.ErrorResponse(ctx, erro)
		}
		return webserver.BadJSONResponse(ctx, err)
	}

	if err := request.Validate(); err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	feedbackService := service.FeedbackServiceNew(app.DB)
	response, err := feedbackService.Update(id, request)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return webserver.SuccessResponse(ctx, response)
}

// GetAverageRating - obtém a média de avaliações de uma empresa
func (h FeedbackHandler) GetAverageRating(ctx echo.Context) error {
	empresaIDStr := ctx.QueryParam("empresa_id")
	empresaID, err := strconv.Atoi(empresaIDStr)
	if err != nil {
		return webserver.ErrorResponse(ctx, util.WrapError("empresa_id inválido", err, http.StatusBadRequest))
	}

	feedbackService := service.FeedbackServiceNew(app.DB)
	average, count, err := feedbackService.GetAverageRating(empresaID)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	response := map[string]interface{}{
		"empresa_id": empresaID,
		"average":    average,
		"count":      count,
	}

	return webserver.SuccessResponse(ctx, response)
}

// GetRatingDistribution - obtém a distribuição de notas de uma empresa
func (h FeedbackHandler) GetRatingDistribution(ctx echo.Context) error {
	empresaIDStr := ctx.QueryParam("empresa_id")
	empresaID, err := strconv.Atoi(empresaIDStr)
	if err != nil {
		return webserver.ErrorResponse(ctx, util.WrapError("empresa_id inválido", err, http.StatusBadRequest))
	}

	feedbackService := service.FeedbackServiceNew(app.DB)
	distribution, err := feedbackService.GetRatingDistribution(empresaID)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	response := map[string]interface{}{
		"empresa_id":   empresaID,
		"distribution": distribution,
	}

	return webserver.SuccessResponse(ctx, response)
}

// GetRecent - busca feedbacks recentes de uma empresa
func (h FeedbackHandler) GetRecent(ctx echo.Context) error {
	empresaIDStr := ctx.QueryParam("empresa_id")
	empresaID, err := strconv.Atoi(empresaIDStr)
	if err != nil {
		return webserver.ErrorResponse(ctx, util.WrapError("empresa_id inválido", err, http.StatusBadRequest))
	}

	days := 30 // padrão
	if daysStr := ctx.QueryParam("days"); daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil {
			days = d
		}
	}

	page := 1
	if pageStr := ctx.QueryParam("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil {
			page = p
		}
	}

	limit := 10
	if limitStr := ctx.QueryParam("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	feedbackService := service.FeedbackServiceNew(app.DB)
	response, err := feedbackService.GetRecentFeedbacks(empresaID, days, page, limit)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return webserver.SuccessResponse(ctx, response)
}
