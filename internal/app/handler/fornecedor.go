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

	return ctx.JSON(response.StatusCode, response)
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

	return ctx.JSON(response.StatusCode, response)
}

// Logout handler para logout de usuário
func (h FornecedorHandler) Logout(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "Logout realizado com sucesso",
	})
}

// Refresh handler para refresh do token
func (h FornecedorHandler) Refresh(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "Token renovado com sucesso",
		"token":   "novo_jwt_token_aqui",
	})
}
