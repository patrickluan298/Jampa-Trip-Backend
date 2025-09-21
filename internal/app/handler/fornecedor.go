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

type FornecedorHandler struct{}

// Login - realiza o login de um fornecedor
func (h FornecedorHandler) Login(ctx echo.Context) error {

	request := &contract.LoginFornecedorRequest{}

	if err := ctx.Bind(request); err != nil {
		if erro := util.ValidarTipoBody(err); erro != nil {
			return webserver.ErrorResponse(ctx, erro)
		}
		return webserver.BadJSONResponse(ctx, err)
	}

	if err := request.Validate(); err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	serviceFornecedor := service.FornecedorServiceNew(app.DB)
	response, err := serviceFornecedor.Login(request)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}

// Cadastrar - realiza o cadastro de um novo fornecedor
func (h FornecedorHandler) Cadastrar(ctx echo.Context) error {

	request := &contract.CadastrarFornecedorRequest{}

	if err := ctx.Bind(request); err != nil {
		if erro := util.ValidarTipoBody(err); erro != nil {
			return webserver.ErrorResponse(ctx, erro)
		}
		return webserver.BadJSONResponse(ctx, err)
	}

	if err := request.Validate(); err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	serviceFornecedor := service.FornecedorServiceNew(app.DB)
	response, err := serviceFornecedor.Cadastrar(request)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}

// Atualizar - realiza a atualização de um fornecedor existente
func (h FornecedorHandler) Atualizar(ctx echo.Context) error {

	request := &contract.AtualizarFornecedorRequest{}

	if err := ctx.Bind(request); err != nil {
		if erro := util.ValidarTipoBody(err); erro != nil {
			return webserver.ErrorResponse(ctx, erro)
		}
		return webserver.BadJSONResponse(ctx, err)
	}

	if err := request.Validate(); err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	serviceFornecedor := service.FornecedorServiceNew(app.DB)
	response, err := serviceFornecedor.Atualizar(request)
	if err != nil {
		return webserver.ErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, response)
}
