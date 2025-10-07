package mercadopago

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jampa_trip/pkg/mercadopago"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name        string
		accessToken string
		baseURL     string
		expected    *mercadopago.Client
	}{
		{
			name:        "Valid client creation",
			accessToken: "test-token",
			baseURL:     "https://api.mercadopago.com",
			expected: &mercadopago.Client{
				AccessToken: "test-token",
				BaseURL:     "https://api.mercadopago.com",
				HTTPClient:  &http.Client{Timeout: 30 * time.Second},
			},
		},
		{
			name:        "Client with empty token",
			accessToken: "",
			baseURL:     "https://api.mercadopago.com",
			expected: &mercadopago.Client{
				AccessToken: "",
				BaseURL:     "https://api.mercadopago.com",
				HTTPClient:  &http.Client{Timeout: 30 * time.Second},
			},
		},
		{
			name:        "Client with custom base URL",
			accessToken: "test-token",
			baseURL:     "https://custom-api.com",
			expected: &mercadopago.Client{
				AccessToken: "test-token",
				BaseURL:     "https://custom-api.com",
				HTTPClient:  &http.Client{Timeout: 30 * time.Second},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := mercadopago.NewClient(tt.accessToken, tt.baseURL)

			if client == nil {
				t.Errorf("NewClient() returned nil")
				return
			}

			if client.AccessToken != tt.expected.AccessToken {
				t.Errorf("NewClient() AccessToken = %s, expected %s", client.AccessToken, tt.expected.AccessToken)
			}

			if client.BaseURL != tt.expected.BaseURL {
				t.Errorf("NewClient() BaseURL = %s, expected %s", client.BaseURL, tt.expected.BaseURL)
			}

			if client.HTTPClient == nil {
				t.Errorf("NewClient() HTTPClient is nil")
				return
			}

			if client.HTTPClient.Timeout != 30*time.Second {
				t.Errorf("NewClient() HTTPClient timeout = %v, expected %v", client.HTTPClient.Timeout, 30*time.Second)
			}
		})
	}
}

func TestClient_CreateOrder(t *testing.T) {
	tests := []struct {
		name          string
		orderReq      *mercadopago.OrderRequest
		mockResponse  string
		mockStatus    int
		expectedError bool
	}{
		{
			name: "Valid order creation",
			orderReq: &mercadopago.OrderRequest{
				ExternalReference: "order-123",
				TotalAmount:       100.50,
				Items: []mercadopago.OrderItem{
					{
						ID:          "item-1",
						Title:       "Test Item",
						Description: "Test Description",
						Quantity:    1,
						CurrencyID:  "BRL",
						UnitPrice:   100.50,
					},
				},
				Payer: mercadopago.Payer{
					Name:  "Test User",
					Email: "test@example.com",
				},
			},
			mockResponse: `{
				"id": "order-123",
				"external_reference": "order-123",
				"total_amount": 100.50,
				"status": "pending",
				"status_detail": "pending",
				"date_created": "2024-01-01T00:00:00Z",
				"date_last_updated": "2024-01-01T00:00:00Z"
			}`,
			mockStatus:    http.StatusCreated,
			expectedError: false,
		},
		{
			name: "Order creation with error response",
			orderReq: &mercadopago.OrderRequest{
				ExternalReference: "order-123",
				TotalAmount:       100.50,
			},
			mockResponse: `{
				"message": "Invalid request",
				"error": "validation_error",
				"status": 400
			}`,
			mockStatus:    http.StatusBadRequest,
			expectedError: true,
		},
		{
			name: "Order creation with server error",
			orderReq: &mercadopago.OrderRequest{
				ExternalReference: "order-123",
				TotalAmount:       100.50,
			},
			mockResponse: `{
				"message": "Internal server error",
				"error": "server_error",
				"status": 500
			}`,
			mockStatus:    http.StatusInternalServerError,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("Expected POST request, got %s", r.Method)
				}

				if r.URL.Path != "/v1/orders" {
					t.Errorf("Expected /v1/orders path, got %s", r.URL.Path)
				}

				authHeader := r.Header.Get("Authorization")
				if authHeader != "Bearer test-token" {
					t.Errorf("Expected Authorization header 'Bearer test-token', got %s", authHeader)
				}

				contentType := r.Header.Get("Content-Type")
				if contentType != "application/json" {
					t.Errorf("Expected Content-Type 'application/json', got %s", contentType)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.mockStatus)
				w.Write([]byte(tt.mockResponse))
			}))
			defer server.Close()

			client := mercadopago.NewClient("test-token", server.URL)

			result, err := client.CreateOrder(tt.orderReq)

			if (err != nil) != tt.expectedError {
				t.Errorf("CreateOrder() error = %v, expectedError = %v", err, tt.expectedError)
			}

			if !tt.expectedError && result != nil {
				if result.ID == "" {
					t.Errorf("CreateOrder() returned order with empty ID")
				}
				if result.ExternalReference != tt.orderReq.ExternalReference {
					t.Errorf("CreateOrder() returned order with external_reference %s, expected %s", result.ExternalReference, tt.orderReq.ExternalReference)
				}
			}
		})
	}
}

func TestClient_GetOrder(t *testing.T) {
	tests := []struct {
		name          string
		orderID       string
		mockResponse  string
		mockStatus    int
		expectedError bool
	}{
		{
			name:    "Valid order retrieval",
			orderID: "order-123",
			mockResponse: `{
				"id": "order-123",
				"external_reference": "order-123",
				"total_amount": 100.50,
				"status": "pending",
				"status_detail": "pending",
				"date_created": "2024-01-01T00:00:00Z",
				"date_last_updated": "2024-01-01T00:00:00Z"
			}`,
			mockStatus:    http.StatusOK,
			expectedError: false,
		},
		{
			name:    "Order not found",
			orderID: "nonexistent",
			mockResponse: `{
				"message": "Order not found",
				"error": "not_found",
				"status": 404
			}`,
			mockStatus:    http.StatusNotFound,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("Expected GET request, got %s", r.Method)
				}

				expectedPath := "/v1/orders/" + tt.orderID
				if r.URL.Path != expectedPath {
					t.Errorf("Expected %s path, got %s", expectedPath, r.URL.Path)
				}

				authHeader := r.Header.Get("Authorization")
				if authHeader != "Bearer test-token" {
					t.Errorf("Expected Authorization header 'Bearer test-token', got %s", authHeader)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.mockStatus)
				w.Write([]byte(tt.mockResponse))
			}))
			defer server.Close()

			client := mercadopago.NewClient("test-token", server.URL)

			result, err := client.GetOrder(tt.orderID)

			if (err != nil) != tt.expectedError {
				t.Errorf("GetOrder() error = %v, expectedError = %v", err, tt.expectedError)
			}

			if !tt.expectedError && result != nil {
				if result.ID != tt.orderID {
					t.Errorf("GetOrder() returned order with ID %s, expected %s", result.ID, tt.orderID)
				}
			}
		})
	}
}

func TestClient_CreatePayment(t *testing.T) {
	tests := []struct {
		name          string
		paymentReq    *mercadopago.PaymentRequest
		mockResponse  string
		mockStatus    int
		expectedError bool
	}{
		{
			name: "Valid payment creation",
			paymentReq: &mercadopago.PaymentRequest{
				TransactionAmount: 100.50,
				Description:       "Test Payment",
				PaymentMethodID:   "credit_card",
				Payer: mercadopago.PaymentPayer{
					Email: "test@example.com",
				},
			},
			mockResponse: `{
				"id": 123456789,
				"status": "pending",
				"status_detail": "pending",
				"transaction_amount": 100.50,
				"description": "Test Payment",
				"payment_method_id": "credit_card",
				"date_created": "2024-01-01T00:00:00Z",
				"date_last_updated": "2024-01-01T00:00:00Z"
			}`,
			mockStatus:    http.StatusCreated,
			expectedError: false,
		},
		{
			name: "Payment creation with error",
			paymentReq: &mercadopago.PaymentRequest{
				TransactionAmount: 100.50,
				Description:       "Test Payment",
				PaymentMethodID:   "invalid_method",
			},
			mockResponse: `{
				"message": "Invalid payment method",
				"error": "validation_error",
				"status": 400
			}`,
			mockStatus:    http.StatusBadRequest,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("Expected POST request, got %s", r.Method)
				}

				if r.URL.Path != "/v1/payments" {
					t.Errorf("Expected /v1/payments path, got %s", r.URL.Path)
				}

				authHeader := r.Header.Get("Authorization")
				if authHeader != "Bearer test-token" {
					t.Errorf("Expected Authorization header 'Bearer test-token', got %s", authHeader)
				}

				contentType := r.Header.Get("Content-Type")
				if contentType != "application/json" {
					t.Errorf("Expected Content-Type 'application/json', got %s", contentType)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.mockStatus)
				w.Write([]byte(tt.mockResponse))
			}))
			defer server.Close()

			client := mercadopago.NewClient("test-token", server.URL)

			result, err := client.CreatePayment(tt.paymentReq)

			if (err != nil) != tt.expectedError {
				t.Errorf("CreatePayment() error = %v, expectedError = %v", err, tt.expectedError)
			}

			if !tt.expectedError && result != nil {
				if result.ID == 0 {
					t.Errorf("CreatePayment() returned payment with zero ID")
				}
				if result.TransactionAmount != tt.paymentReq.TransactionAmount {
					t.Errorf("CreatePayment() returned payment with amount %f, expected %f", result.TransactionAmount, tt.paymentReq.TransactionAmount)
				}
			}
		})
	}
}

func TestClient_CreatePIXPayment(t *testing.T) {
	tests := []struct {
		name          string
		pixReq        *mercadopago.PIXRequest
		mockResponse  string
		mockStatus    int
		expectedError bool
	}{
		{
			name: "Valid PIX payment creation",
			pixReq: &mercadopago.PIXRequest{
				TransactionAmount: 100.50,
				Description:       "Test PIX Payment",
				PaymentMethodID:   "pix",
				Payer: mercadopago.PaymentPayer{
					Email: "test@example.com",
				},
			},
			mockResponse: `{
				"id": 123456789,
				"status": "pending",
				"status_detail": "pending",
				"transaction_amount": 100.50,
				"description": "Test PIX Payment",
				"payment_method_id": "pix",
				"date_created": "2024-01-01T00:00:00Z",
				"date_last_updated": "2024-01-01T00:00:00Z",
				"point_of_interaction": {
					"type": "PIX",
					"application_data": {
						"name": "Mercado Pago",
						"version": "1.0"
					}
				}
			}`,
			mockStatus:    http.StatusCreated,
			expectedError: false,
		},
		{
			name: "PIX payment creation with error",
			pixReq: &mercadopago.PIXRequest{
				TransactionAmount: 100.50,
				Description:       "Test PIX Payment",
				PaymentMethodID:   "invalid_pix",
			},
			mockResponse: `{
				"message": "Invalid PIX configuration",
				"error": "validation_error",
				"status": 400
			}`,
			mockStatus:    http.StatusBadRequest,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("Expected POST request, got %s", r.Method)
				}

				if r.URL.Path != "/v1/payments" {
					t.Errorf("Expected /v1/payments path, got %s", r.URL.Path)
				}

				authHeader := r.Header.Get("Authorization")
				if authHeader != "Bearer test-token" {
					t.Errorf("Expected Authorization header 'Bearer test-token', got %s", authHeader)
				}

				contentType := r.Header.Get("Content-Type")
				if contentType != "application/json" {
					t.Errorf("Expected Content-Type 'application/json', got %s", contentType)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.mockStatus)
				w.Write([]byte(tt.mockResponse))
			}))
			defer server.Close()

			client := mercadopago.NewClient("test-token", server.URL)

			result, err := client.CreatePIXPayment(tt.pixReq)

			if (err != nil) != tt.expectedError {
				t.Errorf("CreatePIXPayment() error = %v, expectedError = %v", err, tt.expectedError)
			}

			if !tt.expectedError && result != nil {
				if result.ID == 0 {
					t.Errorf("CreatePIXPayment() returned payment with zero ID")
				}
				if result.TransactionAmount != tt.pixReq.TransactionAmount {
					t.Errorf("CreatePIXPayment() returned payment with amount %f, expected %f", result.TransactionAmount, tt.pixReq.TransactionAmount)
				}
				if result.PaymentMethodID != tt.pixReq.PaymentMethodID {
					t.Errorf("CreatePIXPayment() returned payment with method %s, expected %s", result.PaymentMethodID, tt.pixReq.PaymentMethodID)
				}
			}
		})
	}
}

func TestClient_GetPayment(t *testing.T) {
	tests := []struct {
		name          string
		paymentID     string
		mockResponse  string
		mockStatus    int
		expectedError bool
	}{
		{
			name:      "Valid payment retrieval",
			paymentID: "123456789",
			mockResponse: `{
				"id": 123456789,
				"status": "approved",
				"status_detail": "accredited",
				"transaction_amount": 100.50,
				"description": "Test Payment",
				"payment_method_id": "credit_card",
				"date_created": "2024-01-01T00:00:00Z",
				"date_approved": "2024-01-01T00:01:00Z",
				"date_last_updated": "2024-01-01T00:01:00Z"
			}`,
			mockStatus:    http.StatusOK,
			expectedError: false,
		},
		{
			name:      "Payment not found",
			paymentID: "nonexistent",
			mockResponse: `{
				"message": "Payment not found",
				"error": "not_found",
				"status": 404
			}`,
			mockStatus:    http.StatusNotFound,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("Expected GET request, got %s", r.Method)
				}

				expectedPath := "/v1/payments/" + tt.paymentID
				if r.URL.Path != expectedPath {
					t.Errorf("Expected %s path, got %s", expectedPath, r.URL.Path)
				}

				authHeader := r.Header.Get("Authorization")
				if authHeader != "Bearer test-token" {
					t.Errorf("Expected Authorization header 'Bearer test-token', got %s", authHeader)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.mockStatus)
				w.Write([]byte(tt.mockResponse))
			}))
			defer server.Close()

			client := mercadopago.NewClient("test-token", server.URL)

			result, err := client.GetPayment(tt.paymentID)

			if (err != nil) != tt.expectedError {
				t.Errorf("GetPayment() error = %v, expectedError = %v", err, tt.expectedError)
			}

			if !tt.expectedError && result != nil {
				if result.ID == 0 {
					t.Errorf("GetPayment() returned payment with zero ID")
				}
			}
		})
	}
}

func TestClient_HTTPErrorHandling(t *testing.T) {
	t.Run("Network error", func(t *testing.T) {
		client := mercadopago.NewClient("test-token", "http://invalid-url:9999")

		orderReq := &mercadopago.OrderRequest{
			ExternalReference: "order-123",
			TotalAmount:       100.50,
		}

		_, err := client.CreateOrder(orderReq)
		if err == nil {
			t.Errorf("CreateOrder() should have failed with network error")
		}
	})

	t.Run("Timeout error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(2 * time.Second)
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		client := mercadopago.NewClient("test-token", server.URL)
		client.HTTPClient.Timeout = 100 * time.Millisecond

		orderReq := &mercadopago.OrderRequest{
			ExternalReference: "order-123",
			TotalAmount:       100.50,
		}

		_, err := client.CreateOrder(orderReq)
		if err == nil {
			t.Errorf("CreateOrder() should have failed with timeout error")
		}
	})
}

func TestClient_ContextHandling(t *testing.T) {
	t.Run("Context cancellation", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(2 * time.Second)
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		client := mercadopago.NewClient("test-token", server.URL)

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		paymentReq := &mercadopago.PaymentRequest{
			TransactionAmount: 100.50,
			Description:       "Test Payment",
			PaymentMethodID:   "credit_card",
		}

		_, err := client.CreateCreditCardPayment(ctx, &mercadopago.CreditCardPaymentRequest{
			TransactionAmount: paymentReq.TransactionAmount,
			Description:       paymentReq.Description,
			PaymentMethodID:   paymentReq.PaymentMethodID,
			Token:             "test-token-12345678",
		})

		if err == nil {
			t.Errorf("CreateCreditCardPayment() should have failed with context timeout")
		}
	})
}
