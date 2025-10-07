package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jampa_trip/internal/handler"
	"github.com/jampa_trip/tests/testutils"
	"github.com/labstack/echo/v4"
)

// TestHealthCheckResponse_HealthCheck_Refactored demonstrates the refactored health check test
func TestHealthCheckResponse_HealthCheck_Refactored(t *testing.T) {
	// Arrange
	e := echo.New()
	healthHandler := handler.HealthCheckResponse{}

	tests := []struct {
		name           string
		method         string
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:           "GET health check",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"nome":     "Jampa-Trip",
				"versao":   "",
				"mensagem": "Aplicação up e em pleno funcionamento",
			},
		},
		{
			name:           "HEAD health check",
			method:         http.MethodHead,
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"nome":     "Jampa-Trip",
				"versao":   "",
				"mensagem": "Aplicação up e em pleno funcionamento",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			req := testutils.CreateTestRequest(t, tt.method, "/health-check", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Act
			err := healthHandler.HealthCheck(c)

			// Assert
			testutils.AssertNoError(t, err)
			testutils.AssertStatusCode(t, rec, tt.expectedStatus)

			// Verify JSON response
			var response map[string]interface{}
			testutils.ParseJSONResponse(t, rec, &response)

			testutils.AssertEqual(t, response["nome"], tt.expectedBody["nome"])
			testutils.AssertEqual(t, response["mensagem"], tt.expectedBody["mensagem"])
		})
	}
}

// TestHealthCheckResponse_Integration_Refactored demonstrates integration test
func TestHealthCheckResponse_Integration_Refactored(t *testing.T) {
	// Arrange
	e := echo.New()
	e.GET("/health-check", handler.HealthCheckResponse{}.HealthCheck)
	e.HEAD("/health-check", handler.HealthCheckResponse{}.HealthCheck)

	tests := []struct {
		name           string
		method         string
		expectedStatus int
	}{
		{
			name:           "GET endpoint",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "HEAD endpoint",
			method:         http.MethodHead,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			req := testutils.CreateTestRequest(t, tt.method, "/health-check", nil)
			rec := httptest.NewRecorder()

			// Act
			e.ServeHTTP(rec, req)

			// Assert
			testutils.AssertStatusCode(t, rec, tt.expectedStatus)
		})
	}
}

// TestHealthCheckResponse_WithQueryParams_Refactored tests with query parameters
func TestHealthCheckResponse_WithQueryParams_Refactored(t *testing.T) {
	// Arrange
	e := echo.New()
	healthHandler := handler.HealthCheckResponse{}

	req := testutils.CreateTestRequest(t, http.MethodGet, "/health-check?param=value&test=123", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Act
	err := healthHandler.HealthCheck(c)

	// Assert
	testutils.AssertNoError(t, err)
	testutils.AssertStatusCode(t, rec, http.StatusOK)

	var response map[string]interface{}
	testutils.ParseJSONResponse(t, rec, &response)

	testutils.AssertEqual(t, response["nome"], "Jampa-Trip")
	testutils.AssertEqual(t, response["mensagem"], "Aplicação up e em pleno funcionamento")
}

// TestHealthCheckResponse_MultipleRequests_Refactored tests multiple consecutive requests
func TestHealthCheckResponse_MultipleRequests_Refactored(t *testing.T) {
	// Arrange
	e := echo.New()
	healthHandler := handler.HealthCheckResponse{}

	// Act & Assert
	for i := 0; i < 5; i++ {
		req := testutils.CreateTestRequest(t, http.MethodGet, "/health-check", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := healthHandler.HealthCheck(c)

		testutils.AssertNoError(t, err)
		testutils.AssertStatusCode(t, rec, http.StatusOK)

		var response map[string]interface{}
		testutils.ParseJSONResponse(t, rec, &response)

		testutils.AssertEqual(t, response["nome"], "Jampa-Trip")
	}
}

