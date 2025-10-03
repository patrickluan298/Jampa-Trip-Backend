package handler

type ReservaHandler struct{}

// Create - cria uma nova reserva
// func (h ReservaHandler) Create(ctx echo.Context) error {
// 	request := &contract.CreateReservaRequest{}

// 	if err := ctx.Bind(request); err != nil {
// 		if erro := util.ValidarTipoBody(err); erro != nil {
// 			return webserver.ErrorResponse(ctx, erro)
// 		}
// 		return webserver.BadJSONResponse(ctx, err)
// 	}

// 	if err := request.Validate(); err != nil {
// 		return webserver.ErrorResponse(ctx, err)
// 	}

// 	reservaService := service.ReservaServiceNew(app.DB)
// 	response, err := reservaService.Create(request)
// 	if err != nil {
// 		return webserver.ErrorResponse(ctx, err)
// 	}

// 	return webserver.SuccessResponse(ctx, response)
// }

// Get - busca uma reserva pelo ID
// func (h ReservaHandler) Get(ctx echo.Context) error {

// 	ID, err := strconv.Atoi(ctx.Param("id"))
// 	if err != nil {
// 		return webserver.InvalidIDResponse(ctx, err)
// 	}

// 	if ID < 1 {
// 		return ctx.JSON(http.StatusBadRequest, contract.ResponseJSON{
// 			StatusCode: http.StatusBadRequest,
// 			Message:    "ID não pode ser zero ou negativo",
// 		})
// 	}

// 	serviceReserva := service.ReservaServiceNew(app.DB)
// 	response, err := serviceReserva.Get(ID)
// 	if err != nil {
// 		return webserver.ErrorResponse(ctx, err)
// 	}

// 	return ctx.JSON(http.StatusOK, response)
// }

// List - lista reservas
// func (h ReservaHandler) List(ctx echo.Context) error {
// 	request := &contract.ListReservaRequest{}

// 	if clienteIDStr := ctx.QueryParam("cliente_id"); clienteIDStr != "" {
// 		if clienteID, err := strconv.Atoi(clienteIDStr); err == nil {
// 			request.ClienteID = clienteID
// 		}
// 	}

// 	if empresaIDStr := ctx.QueryParam("empresa_id"); empresaIDStr != "" {
// 		if empresaID, err := strconv.Atoi(empresaIDStr); err == nil {
// 			request.EmpresaID = empresaID
// 		}
// 	}

// 	request.Status = ctx.QueryParam("status")

// 	if pageStr := ctx.QueryParam("page"); pageStr != "" {
// 		if page, err := strconv.Atoi(pageStr); err == nil {
// 			request.Page = page
// 		}
// 	}

// 	if limitStr := ctx.QueryParam("limit"); limitStr != "" {
// 		if limit, err := strconv.Atoi(limitStr); err == nil {
// 			request.Limit = limit
// 		}
// 	}

// 	if err := request.Validate(); err != nil {
// 		return webserver.ErrorResponse(ctx, err)
// 	}

// 	reservaService := service.ReservaServiceNew(app.DB)
// 	response, err := reservaService.List(request)
// 	if err != nil {
// 		return webserver.ErrorResponse(ctx, err)
// 	}

// 	return webserver.SuccessResponse(ctx, response)
// }

// Update - atualiza uma reserva
// func (h ReservaHandler) Update(ctx echo.Context) error {
// 	idStr := ctx.Param("id")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		return webserver.ErrorResponse(ctx, util.WrapError("ID inválido", err, http.StatusBadRequest))
// 	}

// 	request := &contract.UpdateReservaRequest{}

// 	if err := ctx.Bind(request); err != nil {
// 		if erro := util.ValidarTipoBody(err); erro != nil {
// 			return webserver.ErrorResponse(ctx, erro)
// 		}
// 		return webserver.BadJSONResponse(ctx, err)
// 	}

// 	if err := request.Validate(); err != nil {
// 		return webserver.ErrorResponse(ctx, err)
// 	}

// 	reservaService := service.ReservaServiceNew(app.DB)
// 	response, err := reservaService.Update(id, request)
// 	if err != nil {
// 		return webserver.ErrorResponse(ctx, err)
// 	}

// 	return webserver.SuccessResponse(ctx, response)
// }

// Cancel - cancela uma reserva
// func (h ReservaHandler) Cancel(ctx echo.Context) error {
// 	idStr := ctx.Param("id")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		return webserver.ErrorResponse(ctx, util.WrapError("ID inválido", err, http.StatusBadRequest))
// 	}

// 	request := &contract.CancelarReservaRequest{ID: id}

// 	if err := request.Validate(); err != nil {
// 		return webserver.ErrorResponse(ctx, err)
// 	}

// 	reservaService := service.ReservaServiceNew(app.DB)
// 	response, err := reservaService.Cancel(request)
// 	if err != nil {
// 		return webserver.ErrorResponse(ctx, err)
// 	}

// 	return webserver.SuccessResponse(ctx, response)
// }

// GetUpcoming - busca reservas futuras de um cliente
// func (h ReservaHandler) GetUpcoming(ctx echo.Context) error {
// 	clienteIDStr := ctx.QueryParam("cliente_id")
// 	clienteID, err := strconv.Atoi(clienteIDStr)
// 	if err != nil {
// 		return webserver.ErrorResponse(ctx, util.WrapError("cliente_id inválido", err, http.StatusBadRequest))
// 	}

// 	page := 1
// 	if pageStr := ctx.QueryParam("page"); pageStr != "" {
// 		if p, err := strconv.Atoi(pageStr); err == nil {
// 			page = p
// 		}
// 	}

// 	limit := 10
// 	if limitStr := ctx.QueryParam("limit"); limitStr != "" {
// 		if l, err := strconv.Atoi(limitStr); err == nil {
// 			limit = l
// 		}
// 	}

// 	reservaService := service.ReservaServiceNew(app.DB)
// 	response, err := reservaService.GetUpcoming(clienteID, page, limit)
// 	if err != nil {
// 		return webserver.ErrorResponse(ctx, err)
// 	}

// 	return webserver.SuccessResponse(ctx, response)
// }

// GetHistory - busca histórico de reservas de um cliente
// func (h ReservaHandler) GetHistory(ctx echo.Context) error {
// 	clienteIDStr := ctx.QueryParam("cliente_id")
// 	clienteID, err := strconv.Atoi(clienteIDStr)
// 	if err != nil {
// 		return webserver.ErrorResponse(ctx, util.WrapError("cliente_id inválido", err, http.StatusBadRequest))
// 	}

// 	page := 1
// 	if pageStr := ctx.QueryParam("page"); pageStr != "" {
// 		if p, err := strconv.Atoi(pageStr); err == nil {
// 			page = p
// 		}
// 	}

// 	limit := 10
// 	if limitStr := ctx.QueryParam("limit"); limitStr != "" {
// 		if l, err := strconv.Atoi(limitStr); err == nil {
// 			limit = l
// 		}
// 	}

// 	reservaService := service.ReservaServiceNew(app.DB)
// 	response, err := reservaService.GetHistory(clienteID, page, limit)
// 	if err != nil {
// 		return webserver.ErrorResponse(ctx, err)
// 	}

// 	return webserver.SuccessResponse(ctx, response)
// }
