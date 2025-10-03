package util

import (
	"encoding/json"
	"net/http"
	"runtime"
)

type ResponseError struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

type AppError struct {
	Msg        string
	Err        error
	File       string
	Line       int
	FuncName   string
	StatusCode int
}

// Implementação do método Error()
func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Msg
}

// MarshalJSON - serializa o erro em JSON
func (e *AppError) MarshalJSON() ([]byte, error) {
	type Alias AppError

	var errField interface{}
	if wrappedErr, ok := e.Err.(*AppError); ok {
		errField = wrappedErr
	} else if e.Err != nil {
		errField = e.Err.Error()
	}

	return json.Marshal(&struct {
		Err interface{} `json:"Err"`
		*Alias
	}{
		Err:   errField,
		Alias: (*Alias)(e),
	})
}

// WrapError - cria um erro estruturado com status HTTP
func WrapError(msg string, err error, statusCode ...int) *AppError {
	pc, file, line, _ := runtime.Caller(1)

	sc := http.StatusInternalServerError

	if len(statusCode) > 0 {
		sc = statusCode[0]
	}

	if appErr, ok := err.(*AppError); ok {
		if len(statusCode) == 0 {
			sc = appErr.StatusCode
		}
	}

	return &AppError{
		Msg:        msg,
		Err:        err,
		File:       file,
		Line:       line,
		FuncName:   runtime.FuncForPC(pc).Name(),
		StatusCode: sc,
	}
}

// HandleError - processa um erro e retorna um ResponseJSON estruturado
func HandleError(err error) ResponseError {
	appErr, ok := err.(*AppError)
	if !ok {
		return ResponseError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Erro interno do servidor",
		}
	}

	return ResponseError{
		StatusCode: appErr.StatusCode,
		Message:    appErr.Msg,
	}
}
