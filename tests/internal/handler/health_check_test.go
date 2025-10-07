package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jampa_trip/internal/handler"
	"github.com/labstack/echo/v4"
)

func TestHealthCheckResponse_HealthCheck(t *testing.T) {
	e := echo.New()
	handler := handler.HealthCheckResponse{}

	tests := []struct {
		name           string
		method         string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "GET health check",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"nome":"Jampa-Trip","versao":"","mensagem":"Aplicação up e em pleno funcionamento"}`,
		},
		{
			name:           "HEAD health check",
			method:         http.MethodHead,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"nome":"Jampa-Trip","versao":"","mensagem":"Aplicação up e em pleno funcionamento"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/health-check", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := handler.HealthCheck(c)
			if err != nil {
				t.Errorf("HealthCheck() error = %v", err)
			}

			if rec.Code != tt.expectedStatus {
				t.Errorf("HealthCheck() status = %d, expected %d", rec.Code, tt.expectedStatus)
			}

			body := rec.Body.String()
			if tt.method == http.MethodGet {
				// Remove whitespace for comparison
				bodyTrimmed := strings.TrimSpace(body)
				expectedTrimmed := strings.TrimSpace(tt.expectedBody)
				if bodyTrimmed != expectedTrimmed {
					t.Errorf("HealthCheck() body = %s, expected %s", bodyTrimmed, expectedTrimmed)
				}
			}
		})
	}
}

func TestHealthCheckResponse_HealthCheckIntegration(t *testing.T) {
	e := echo.New()

	// Test GET endpoint
	e.GET("/health-check", handler.HealthCheckResponse{}.HealthCheck)

	req := httptest.NewRequest(http.MethodGet, "/health-check", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("GET /health-check status = %d, expected %d", rec.Code, http.StatusOK)
	}

	body := rec.Body.String()
	expectedBody := `{"nome":"Jampa-Trip","versao":"","mensagem":"Aplicação up e em pleno funcionamento"}`
	bodyTrimmed := strings.TrimSpace(body)
	expectedTrimmed := strings.TrimSpace(expectedBody)
	if bodyTrimmed != expectedTrimmed {
		t.Errorf("GET /health-check body = %s, expected %s", bodyTrimmed, expectedTrimmed)
	}
}

func TestHealthCheckResponse_HealthCheckHEAD(t *testing.T) {
	e := echo.New()

	// Test HEAD endpoint
	e.HEAD("/health-check", handler.HealthCheckResponse{}.HealthCheck)

	req := httptest.NewRequest(http.MethodHead, "/health-check", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("HEAD /health-check status = %d, expected %d", rec.Code, http.StatusOK)
	}

	// HEAD requests should have the same body as GET
	body := rec.Body.String()
	expectedBody := `{"nome":"Jampa-Trip","versao":"","mensagem":"Aplicação up e em pleno funcionamento"}`
	bodyTrimmed := strings.TrimSpace(body)
	expectedTrimmed := strings.TrimSpace(expectedBody)
	if bodyTrimmed != expectedTrimmed {
		t.Errorf("HEAD /health-check body = %s, expected %s", bodyTrimmed, expectedTrimmed)
	}
}

func TestHealthCheckResponse_HealthCheckMultipleRequests(t *testing.T) {
	e := echo.New()
	handler := handler.HealthCheckResponse{}

	// Test multiple consecutive requests
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest(http.MethodGet, "/health-check", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.HealthCheck(c)
		if err != nil {
			t.Errorf("HealthCheck() error on request %d: %v", i+1, err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("HealthCheck() status on request %d = %d, expected %d", i+1, rec.Code, http.StatusOK)
		}

		body := rec.Body.String()
		expectedBody := `{"nome":"Jampa-Trip","versao":"","mensagem":"Aplicação up e em pleno funcionamento"}`
		bodyTrimmed := strings.TrimSpace(body)
		expectedTrimmed := strings.TrimSpace(expectedBody)
		if bodyTrimmed != expectedTrimmed {
			t.Errorf("HealthCheck() body on request %d = %s, expected %s", i+1, bodyTrimmed, expectedTrimmed)
		}
	}
}

func TestHealthCheckResponse_HealthCheckWithQueryParams(t *testing.T) {
	e := echo.New()
	handler := handler.HealthCheckResponse{}

	// Test with query parameters (should still work)
	req := httptest.NewRequest(http.MethodGet, "/health-check?param=value&test=123", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.HealthCheck(c)
	if err != nil {
		t.Errorf("HealthCheck() error = %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("HealthCheck() status = %d, expected %d", rec.Code, http.StatusOK)
	}

	body := rec.Body.String()
	expectedBody := `{"nome":"Jampa-Trip","versao":"","mensagem":"Aplicação up e em pleno funcionamento"}`
	bodyTrimmed := strings.TrimSpace(body)
	expectedTrimmed := strings.TrimSpace(expectedBody)
	if bodyTrimmed != expectedTrimmed {
		t.Errorf("HealthCheck() body = %s, expected %s", bodyTrimmed, expectedTrimmed)
	}
}

func TestHealthCheckResponse_HealthCheckWithHeaders(t *testing.T) {
	e := echo.New()
	handler := handler.HealthCheckResponse{}

	// Test with custom headers
	req := httptest.NewRequest(http.MethodGet, "/health-check", nil)
	req.Header.Set("User-Agent", "Test-Agent")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer token")

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.HealthCheck(c)
	if err != nil {
		t.Errorf("HealthCheck() error = %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("HealthCheck() status = %d, expected %d", rec.Code, http.StatusOK)
	}

	body := rec.Body.String()
	expectedBody := `{"nome":"Jampa-Trip","versao":"","mensagem":"Aplicação up e em pleno funcionamento"}`
	bodyTrimmed := strings.TrimSpace(body)
	expectedTrimmed := strings.TrimSpace(expectedBody)
	if bodyTrimmed != expectedTrimmed {
		t.Errorf("HealthCheck() body = %s, expected %s", bodyTrimmed, expectedTrimmed)
	}
}
