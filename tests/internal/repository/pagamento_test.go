package repository

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jampa_trip/internal/model"
	"github.com/jampa_trip/internal/repository"
)

// setupMockDB is already defined in company_test.go

func TestPagamentoRepositoryNew(t *testing.T) {
	db, mock := setupMockDB(t)
	defer mock.ExpectationsWereMet()

	repo := repository.PagamentoRepositoryNew(db)
	if repo == nil {
		t.Error("PagamentoRepositoryNew() returned nil")
	}
}

func TestPagamentoRepository_Create(t *testing.T) {
	db, mock := setupMockDB(t)
	defer mock.ExpectationsWereMet()

	repo := repository.PagamentoRepositoryNew(db)

	tests := []struct {
		name      string
		pagamento *model.Pagamento
		hasError  bool
	}{
		{
			name: "Valid payment",
			pagamento: &model.Pagamento{
				ClienteID:       1,
				EmpresaID:       1,
				Status:          "pending",
				Valor:           150.50,
				Moeda:           "BRL",
				MetodoPagamento: "credit_card",
				Descricao:       "Pagamento do tour",
				MomentoCriacao:  time.Now(),
			},
			hasError: false,
		},
		{
			name: "Payment with PIX method",
			pagamento: &model.Pagamento{
				ClienteID:       1,
				EmpresaID:       1,
				Status:          "pending",
				Valor:           150.50,
				Moeda:           "BRL",
				MetodoPagamento: "pix",
				Descricao:       "Pagamento do tour via PIX",
				MomentoCriacao:  time.Now(),
			},
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectQuery(`INSERT INTO pagamentos`).
				WithArgs(
					tt.pagamento.ClienteID,
					tt.pagamento.EmpresaID,
					tt.pagamento.Status,
					tt.pagamento.Valor,
					tt.pagamento.Moeda,
					tt.pagamento.MetodoPagamento,
					tt.pagamento.Descricao,
				).
				WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

			err := repo.Create(tt.pagamento)
			if (err != nil) != tt.hasError {
				t.Errorf("Create() error = %v, hasError = %v", err, tt.hasError)
			}
		})
	}
}

func TestPagamentoRepository_GetByMercadoPagoPaymentID(t *testing.T) {
	db, mock := setupMockDB(t)
	defer mock.ExpectationsWereMet()

	repo := repository.PagamentoRepositoryNew(db)

	tests := []struct {
		name      string
		paymentID string
		mockRows  *sqlmock.Rows
		expected  *model.Pagamento
		hasError  bool
	}{
		{
			name:      "Valid payment ID",
			paymentID: "123456789",
			mockRows: sqlmock.NewRows([]string{"id", "cliente_id", "empresa_id", "mercado_pago_payment_id", "status", "valor", "metodo_pagamento", "momento_criacao"}).
				AddRow(1, 1, 1, "123456789", "approved", 150.50, "credit_card", time.Now()),
			expected: &model.Pagamento{
				ID:                   1,
				ClienteID:            1,
				EmpresaID:            1,
				MercadoPagoPaymentID: "123456789",
				Status:               "approved",
				Valor:                150.50,
				MetodoPagamento:      "credit_card",
			},
			hasError: false,
		},
		{
			name:      "Non-existent payment ID",
			paymentID: "999999999",
			mockRows:  sqlmock.NewRows([]string{"id", "cliente_id", "empresa_id", "mercado_pago_payment_id", "status", "valor", "metodo_pagamento", "momento_criacao"}),
			expected:  nil,
			hasError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectQuery(`SELECT`).
				WithArgs(tt.paymentID).
				WillReturnRows(tt.mockRows)

			result, err := repo.GetByMercadoPagoPaymentID(tt.paymentID)
			if (err != nil) != tt.hasError {
				t.Errorf("GetByMercadoPagoPaymentID() error = %v, hasError = %v", err, tt.hasError)
			}

			if !tt.hasError && result != nil {
				if result.ID != tt.expected.ID {
					t.Errorf("GetByMercadoPagoPaymentID() ID = %d, expected %d", result.ID, tt.expected.ID)
				}
				if result.Valor != tt.expected.Valor {
					t.Errorf("GetByMercadoPagoPaymentID() Valor = %f, expected %f", result.Valor, tt.expected.Valor)
				}
			}
		})
	}
}

func TestPagamentoRepository_GetByClienteID(t *testing.T) {
	db, mock := setupMockDB(t)
	defer mock.ExpectationsWereMet()

	repo := repository.PagamentoRepositoryNew(db)

	tests := []struct {
		name      string
		clienteID int
		mockRows  *sqlmock.Rows
		expected  int
		hasError  bool
	}{
		{
			name:      "Valid client ID",
			clienteID: 1,
			mockRows: sqlmock.NewRows([]string{"id", "cliente_id", "empresa_id", "status", "valor", "metodo_pagamento", "momento_criacao"}).
				AddRow(1, 1, 1, "approved", 150.50, "credit_card", time.Now()).
				AddRow(2, 1, 1, "pending", 200.00, "pix", time.Now()),
			expected: 2,
			hasError: false,
		},
		{
			name:      "Client with no payments",
			clienteID: 999,
			mockRows:  sqlmock.NewRows([]string{"id", "cliente_id", "empresa_id", "status", "valor", "metodo_pagamento", "momento_criacao"}),
			expected:  0,
			hasError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectQuery(`SELECT`).
				WithArgs(tt.clienteID).
				WillReturnRows(tt.mockRows)

			pagamentos, err := repo.GetByClienteID(tt.clienteID)
			if (err != nil) != tt.hasError {
				t.Errorf("GetByClienteID() error = %v, hasError = %v", err, tt.hasError)
			}

			if !tt.hasError {
				if len(pagamentos) != tt.expected {
					t.Errorf("GetByClienteID() returned %d payments, expected %d", len(pagamentos), tt.expected)
				}
			}
		})
	}
}
