package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/labstack/echo/v4"
)

// Regex patterns for validation
var (
	TimePattern = regexp.MustCompile(`^([0-1]?[0-9]|2[0-3]):[0-5][0-9]$`)
	URLPattern  = regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)
)

// ValidateTimeFormat - validates HH:MM format
func ValidateTimeFormat(timeStr string) bool {
	return TimePattern.MatchString(timeStr)
}

// ValidateURL - validates URL format
func ValidateURL(url string) bool {
	return URLPattern.MatchString(url)
}

// ValidateDateFormat - validates YYYY-MM-DD format
func ValidateDateFormat(dateStr string) error {
	_, err := time.Parse("2006-01-02", dateStr)
	return err
}

// ValidateDates - validates slice of dates
func ValidateDates(dates []string) error {
	for _, date := range dates {
		if err := ValidateDateFormat(date); err != nil {
			return WrapError("Formato de data inválido. Use YYYY-MM-DD", err, http.StatusUnprocessableEntity)
		}
	}
	return nil
}

// ValidateImageURLs - validates slice of image URLs
func ValidateImageURLs(urls []string) error {
	for _, url := range urls {
		if url != "" && !ValidateURL(url) {
			return WrapError("URL de imagem inválida", nil, http.StatusUnprocessableEntity)
		}
	}
	return nil
}

// ValidateBodyType - valida o tipo esperado para os campos da struct
func ValidateBodyType(err error) error {
	if unmarshalTypeError := new(json.UnmarshalTypeError); errors.As(err, &unmarshalTypeError) {
		message := fmt.Sprintf("O campo '%s' espera um valor do tipo %s, mas recebeu um valor inválido", unmarshalTypeError.Field, unmarshalTypeError.Type.String())
		return WrapError(message, errors.New(message), http.StatusUnprocessableEntity)
	}
	return nil
}

// ValidateParameterType valida o tipo esperado para os campos dos parâmetros
func ValidateParameterType(ctx echo.Context, expectedTypes map[string]string) error {

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

// TimeValidator - validador customizado para horários
func TimeValidator() validation.Rule {
	return validation.By(func(value interface{}) error {
		if !ValidateTimeFormat(value.(string)) {
			return errors.New("formato de horário inválido")
		}
		return nil
	})
}
