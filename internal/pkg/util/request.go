package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

// ValidarTipoBody valida o tipo esperado para os campos da struct
func ValidarTipoBody(err error) error {
	if unmarshalTypeError := new(json.UnmarshalTypeError); errors.As(err, &unmarshalTypeError) {
		message := fmt.Sprintf("O campo '%s' espera um valor do tipo %s, mas recebeu um valor inválido", unmarshalTypeError.Field, unmarshalTypeError.Type.String())
		return WrapError(message, errors.New(message), http.StatusUnprocessableEntity)
	}
	return nil
}

// ValidarTipoParametros valida o tipo esperado para os campos dos parâmetros
func ValidarTipoParametros(ctx echo.Context, expectedTypes map[string]string) error {

	errorMessages := make(map[string]map[string][]string)
	queryParams := ctx.QueryParams()

	for paramName, values := range queryParams {
		if len(values) == 0 || values[0] == "" {
			continue
		}
		paramValue := values[0]

		expectedType, exists := expectedTypes[paramName]
		if !exists {
			continue
		}

		var receivedType string

		if expectedType == "string" {
			receivedType = "string"
		} else if _, err := strconv.Atoi(paramValue); err == nil {
			receivedType = "int"
		} else if _, err := strconv.ParseFloat(paramValue, 64); err == nil {
			receivedType = "float"
		} else if _, err := strconv.ParseBool(paramValue); err == nil {
			receivedType = "bool"
		} else {
			receivedType = "string"
		}

		if receivedType != expectedType {
			if errorMessages[expectedType] == nil {
				errorMessages[expectedType] = make(map[string][]string)
			}
			errorMessages[expectedType][receivedType] = append(errorMessages[expectedType][receivedType], paramName)
		}
	}

	var finalMessages []string
	for expectedType, receivedMap := range errorMessages {
		for receivedType, params := range receivedMap {
			finalMessages = append(finalMessages, fmt.Sprintf("Os parâmetros '%s' esperam um valor do tipo %s, mas receberam um valor do tipo %s",
				strings.Join(params, "', '"), expectedType, receivedType))
		}
	}

	if len(finalMessages) > 0 {
		message := strings.Join(finalMessages, "; ")
		return WrapError(message, errors.New(message), http.StatusUnprocessableEntity)
	}

	return nil
}
