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

type ReservationHandler struct{}

// Create - cria uma nova reserva
func (h ReservationHandler) Create(ctx echo.Context) error {
	request := &contract.CreateReservaRequest{}

	if err := ctx.Bind(request); err != nil {
		if erro := util.ValidateBodyType(err); erro != nil {
			return webserver.ErrorResponse(ctx, erro)
		}
		return webserver.BadJSONResponse(ctx, err)
	}

	if err := request.Validate(); err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	reservaService := service.ReservaServiceNew(database.DB)
	response, err := reservaService.Create(request)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}

// Get - busca uma reserva pelo ID
func (h ReservationHandler) Get(ctx echo.Context) error {
	ID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return webserver.InvalidIDResponse(ctx, err)
	}

	if ID < 1 {
		return webserver.ErrorResponse(ctx, util.WrapError("ID não pode ser zero ou negativo", nil, http.StatusBadRequest))
	}

	serviceReserva := service.ReservaServiceNew(database.DB)
	request := &contract.GetReservaRequest{ID: ID}
	response, err := serviceReserva.GetByID(request)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}

// List - lista reservas
func (h ReservationHandler) List(ctx echo.Context) error {
	request := &contract.ListReservaRequest{}

	if tourIDStr := ctx.QueryParam("tour_id"); tourIDStr != "" {
		if tourID, err := strconv.Atoi(tourIDStr); err == nil {
			request.TourID = tourID
		}
	}

	if clienteIDStr := ctx.QueryParam("cliente_id"); clienteIDStr != "" {
		if clienteID, err := strconv.Atoi(clienteIDStr); err == nil {
			request.ClienteID = clienteID
		}
	}

	if companyIDStr := ctx.QueryParam("company_id"); companyIDStr != "" {
		if companyID, err := strconv.Atoi(companyIDStr); err == nil {
			request.CompanyID = companyID
		}
	}

	request.Status = ctx.QueryParam("status")

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

	reservaService := service.ReservaServiceNew(database.DB)
	response, err := reservaService.List(request)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}

// Update - atualiza uma reserva
func (h ReservationHandler) Update(ctx echo.Context) error {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return webserver.ErrorResponse(ctx, util.WrapError("ID inválido", err, http.StatusBadRequest))
	}

	request := &contract.UpdateReservaRequest{}

	if err := ctx.Bind(request); err != nil {
		if erro := util.ValidateBodyType(err); erro != nil {
			return webserver.ErrorResponse(ctx, erro)
		}
		return webserver.BadJSONResponse(ctx, err)
	}

	if err := request.Validate(); err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	reservaService := service.ReservaServiceNew(database.DB)
	response, err := reservaService.Update(id, request)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}

// Cancel - cancela uma reserva
func (h ReservationHandler) Cancel(ctx echo.Context) error {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return webserver.ErrorResponse(ctx, util.WrapError("ID inválido", err, http.StatusBadRequest))
	}

	request := &contract.CancelarReservaRequest{ID: id}

	if err := request.Validate(); err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	reservaService := service.ReservaServiceNew(database.DB)
	response, err := reservaService.Cancel(request)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}

// GetUpcoming - busca reservas futuras de um cliente
func (h ReservationHandler) GetUpcoming(ctx echo.Context) error {
	clienteIDStr := ctx.QueryParam("cliente_id")
	clienteID, err := strconv.Atoi(clienteIDStr)
	if err != nil {
		return webserver.ErrorResponse(ctx, util.WrapError("cliente_id inválido", err, http.StatusBadRequest))
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

	reservaService := service.ReservaServiceNew(database.DB)
	response, err := reservaService.GetUpcoming(clienteID, page, limit)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}

// GetHistory - busca histórico de reservas de um cliente
func (h ReservationHandler) GetHistory(ctx echo.Context) error {
	clienteIDStr := ctx.QueryParam("cliente_id")
	clienteID, err := strconv.Atoi(clienteIDStr)
	if err != nil {
		return webserver.ErrorResponse(ctx, util.WrapError("cliente_id inválido", err, http.StatusBadRequest))
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

	reservaService := service.ReservaServiceNew(database.DB)
	response, err := reservaService.GetHistory(clienteID, page, limit)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}
