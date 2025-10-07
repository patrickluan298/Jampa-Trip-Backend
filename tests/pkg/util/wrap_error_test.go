package util

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/jampa_trip/pkg/util"
)

func TestAppError_Error(t *testing.T) {
	tests := []struct {
		name     string
		appError *util.AppError
		expected string
	}{
		{
			name: "Error with wrapped error",
			appError: &util.AppError{
				Msg: "Custom message",
				Err: errors.New("wrapped error"),
			},
			expected: "wrapped error",
		},
		{
			name: "Error without wrapped error",
			appError: &util.AppError{
				Msg: "Custom message",
				Err: nil,
			},
			expected: "Custom message",
		},
		{
			name: "Error with empty message and wrapped error",
			appError: &util.AppError{
				Msg: "",
				Err: errors.New("wrapped error"),
			},
			expected: "wrapped error",
		},
		{
			name: "Error with empty message and no wrapped error",
			appError: &util.AppError{
				Msg: "",
				Err: nil,
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.appError.Error()
			if result != tt.expected {
				t.Errorf("AppError.Error() = %s, expected %s", result, tt.expected)
			}
		})
	}
}

func TestAppError_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		appError *util.AppError
		expected map[string]interface{}
	}{
		{
			name: "Error with simple wrapped error",
			appError: &util.AppError{
				Msg:        "Custom message",
				Err:        errors.New("wrapped error"),
				File:       "test.go",
				Line:       10,
				FuncName:   "TestFunction",
				StatusCode: http.StatusBadRequest,
			},
			expected: map[string]interface{}{
				"Msg":        "Custom message",
				"File":       "test.go",
				"Line":       10,
				"FuncName":   "TestFunction",
				"StatusCode": http.StatusBadRequest,
				"Err":        "wrapped error",
			},
		},
		{
			name: "Error with AppError wrapped error",
			appError: &util.AppError{
				Msg:        "Outer message",
				Err:        &util.AppError{Msg: "Inner message", Err: nil, StatusCode: http.StatusInternalServerError},
				File:       "test.go",
				Line:       20,
				FuncName:   "TestFunction",
				StatusCode: http.StatusBadRequest,
			},
			expected: map[string]interface{}{
				"Msg":        "Outer message",
				"File":       "test.go",
				"Line":       20,
				"FuncName":   "TestFunction",
				"StatusCode": http.StatusBadRequest,
				"Err": map[string]interface{}{
					"Msg":        "Inner message",
					"Err":        nil,
					"StatusCode": http.StatusInternalServerError,
				},
			},
		},
		{
			name: "Error with nil wrapped error",
			appError: &util.AppError{
				Msg:        "Custom message",
				Err:        nil,
				File:       "test.go",
				Line:       30,
				FuncName:   "TestFunction",
				StatusCode: http.StatusInternalServerError,
			},
			expected: map[string]interface{}{
				"Msg":        "Custom message",
				"File":       "test.go",
				"Line":       30,
				"FuncName":   "TestFunction",
				"StatusCode": http.StatusInternalServerError,
				"Err":        nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := tt.appError.MarshalJSON()
			if err != nil {
				t.Errorf("AppError.MarshalJSON() failed: %v", err)
				return
			}

			var result map[string]interface{}
			if err := json.Unmarshal(jsonData, &result); err != nil {
				t.Errorf("Failed to unmarshal JSON: %v", err)
				return
			}

			for key, expectedValue := range tt.expected {
				if _, exists := result[key]; !exists {
					t.Errorf("AppError.MarshalJSON() missing field: %s", key)
					continue
				}

				if key == "Err" {
					continue
				}

				if key == "Line" || key == "StatusCode" {
					continue
				}

				if key == "Msg" || key == "File" || key == "FuncName" {
					if result[key] != expectedValue {
						t.Errorf("AppError.MarshalJSON() field %s = %v, expected %v", key, result[key], expectedValue)
					}
				}
			}
		})
	}
}

func TestWrapError(t *testing.T) {
	tests := []struct {
		name       string
		msg        string
		err        error
		statusCode []int
		expected   *util.AppError
	}{
		{
			name:       "Error with custom status code",
			msg:        "Custom error message",
			err:        errors.New("wrapped error"),
			statusCode: []int{http.StatusBadRequest},
			expected: &util.AppError{
				Msg:        "Custom error message",
				Err:        errors.New("wrapped error"),
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			name:       "Error with default status code",
			msg:        "Custom error message",
			err:        errors.New("wrapped error"),
			statusCode: []int{},
			expected: &util.AppError{
				Msg:        "Custom error message",
				Err:        errors.New("wrapped error"),
				StatusCode: http.StatusInternalServerError,
			},
		},
		{
			name:       "Error with nil wrapped error",
			msg:        "Custom error message",
			err:        nil,
			statusCode: []int{http.StatusNotFound},
			expected: &util.AppError{
				Msg:        "Custom error message",
				Err:        nil,
				StatusCode: http.StatusNotFound,
			},
		},
		{
			name:       "Error with AppError wrapped error",
			msg:        "Outer error",
			err:        &util.AppError{Msg: "Inner error", StatusCode: http.StatusBadRequest},
			statusCode: []int{http.StatusInternalServerError},
			expected: &util.AppError{
				Msg:        "Outer error",
				Err:        &util.AppError{Msg: "Inner error", StatusCode: http.StatusBadRequest},
				StatusCode: http.StatusInternalServerError,
			},
		},
		{
			name:       "Error with AppError wrapped error and no status code override",
			msg:        "Outer error",
			err:        &util.AppError{Msg: "Inner error", StatusCode: http.StatusBadRequest},
			statusCode: []int{},
			expected: &util.AppError{
				Msg:        "Outer error",
				Err:        &util.AppError{Msg: "Inner error", StatusCode: http.StatusBadRequest},
				StatusCode: http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := util.WrapError(tt.msg, tt.err, tt.statusCode...)

			if result.Msg != tt.expected.Msg {
				t.Errorf("WrapError() Msg = %s, expected %s", result.Msg, tt.expected.Msg)
			}

			if result.StatusCode != tt.expected.StatusCode {
				t.Errorf("WrapError() StatusCode = %d, expected %d", result.StatusCode, tt.expected.StatusCode)
			}

			if (result.Err == nil) != (tt.expected.Err == nil) {
				t.Errorf("WrapError() Err = %v, expected %v", result.Err, tt.expected.Err)
			}

			if result.Err != nil && tt.expected.Err != nil {
				if result.Err.Error() != tt.expected.Err.Error() {
					t.Errorf("WrapError() Err = %v, expected %v", result.Err, tt.expected.Err)
				}
			}

			if result.File == "" {
				t.Errorf("WrapError() File should not be empty")
			}

			if result.Line <= 0 {
				t.Errorf("WrapError() Line should be positive, got %d", result.Line)
			}

			if result.FuncName == "" {
				t.Errorf("WrapError() FuncName should not be empty")
			}
		})
	}
}

func TestHandleError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected util.ResponseError
	}{
		{
			name: "AppError with custom status",
			err: &util.AppError{
				Msg:        "Custom error message",
				StatusCode: http.StatusBadRequest,
			},
			expected: util.ResponseError{
				StatusCode: http.StatusBadRequest,
				Message:    "Custom error message",
			},
		},
		{
			name: "AppError with internal server error status",
			err: &util.AppError{
				Msg:        "Internal error message",
				StatusCode: http.StatusInternalServerError,
			},
			expected: util.ResponseError{
				StatusCode: http.StatusInternalServerError,
				Message:    "Internal error message",
			},
		},
		{
			name: "Regular error",
			err:  errors.New("regular error"),
			expected: util.ResponseError{
				StatusCode: http.StatusInternalServerError,
				Message:    "Erro interno do servidor",
			},
		},
		{
			name: "Nil error",
			err:  nil,
			expected: util.ResponseError{
				StatusCode: http.StatusInternalServerError,
				Message:    "Erro interno do servidor",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := util.HandleError(tt.err)

			if result.StatusCode != tt.expected.StatusCode {
				t.Errorf("HandleError() StatusCode = %d, expected %d", result.StatusCode, tt.expected.StatusCode)
			}

			if result.Message != tt.expected.Message {
				t.Errorf("HandleError() Message = %s, expected %s", result.Message, tt.expected.Message)
			}
		})
	}
}

func TestResponseError(t *testing.T) {
	responseError := util.ResponseError{
		StatusCode: http.StatusBadRequest,
		Message:    "Test error message",
	}

	if responseError.StatusCode != http.StatusBadRequest {
		t.Errorf("ResponseError.StatusCode = %d, expected %d", responseError.StatusCode, http.StatusBadRequest)
	}

	if responseError.Message != "Test error message" {
		t.Errorf("ResponseError.Message = %s, expected %s", responseError.Message, "Test error message")
	}
}
