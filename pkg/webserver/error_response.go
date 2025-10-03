package webserver

import (
	"net/http"

	"github.com/jampa_trip/internal/contract"
	"github.com/jampa_trip/pkg/util"
	"github.com/labstack/echo/v4"
)

// ErrorResponse lida com erros do tipo AppError retornados do service, registrando-os e retornando uma resposta JSON estruturada.
func ErrorResponse(ctx echo.Context, err error) error {
	ctx.Set("errorDetails", util.WrapError("Erro ao processar a requisição", err))
	responseError := util.HandleError(err)
	return ctx.JSON(responseError.StatusCode, responseError)
}

// BadJSONResponse lida com erros relacionados a solicitações JSON malformadas.
func BadJSONResponse(ctx echo.Context, err error) error {
	ctx.Set("errorDetails", util.WrapError("Erro na formatação do JSON", err, http.StatusBadRequest))
	return ctx.JSON(http.StatusBadRequest, contract.ResponseJSON{
		StatusCode: http.StatusBadRequest,
		Message:    "Erro no formato do JSON",
	})
}

// InvalidIDResponse lida com erros na conversão do ID recebido.
func InvalidIDResponse(ctx echo.Context, err error) error {
	ctx.Set("errorDetails", util.WrapError("Erro ao converter ID recebido", err, http.StatusInternalServerError))
	return ctx.JSON(http.StatusInternalServerError, contract.ResponseJSON{
		StatusCode: http.StatusInternalServerError,
		Message:    "Erro na conversão do ID recebido",
	})
}
