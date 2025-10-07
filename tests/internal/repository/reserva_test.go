package repository

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jampa_trip/internal/model"
	"github.com/jampa_trip/internal/repository"
)

// setupMockDB is already defined in company_test.go

func TestReservaRepositoryNew(t *testing.T) {
	db, mock := setupMockDB(t)
	defer mock.ExpectationsWereMet()

	repo := repository.ReservaRepositoryNew(db)
	if repo == nil {
		t.Error("ReservaRepositoryNew() returned nil")
	}
}

func TestReservaRepository_Create(t *testing.T) {
	db, mock := setupMockDB(t)
	defer mock.ExpectationsWereMet()

	repo := repository.ReservaRepositoryNew(db)

	tests := []struct {
		name     string
		reserva  *model.Reserva
		hasError bool
	}{
		{
			name: "Valid reservation",
			reserva: &model.Reserva{
				ClienteID:         1,
				EmpresaID:         1,
				Status:            "pendente",
				DataReserva:       time.Now(),
				DataPasseio:       time.Now().AddDate(0, 0, 7),
				QuantidadePessoas: 1,
				ValorTotal:        150.50,
				MomentoCriacao:    time.Now(),
			},
			hasError: false,
		},
		{
			name: "Reservation with confirmed status",
			reserva: &model.Reserva{
				ClienteID:         1,
				EmpresaID:         1,
				Status:            "confirmada",
				DataReserva:       time.Now(),
				DataPasseio:       time.Now().AddDate(0, 0, 7),
				QuantidadePessoas: 2,
				ValorTotal:        300.00,
				MomentoCriacao:    time.Now(),
			},
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectQuery(`INSERT INTO reservas`).
				WithArgs(
					tt.reserva.ClienteID,
					tt.reserva.EmpresaID,
					tt.reserva.Status,
					tt.reserva.DataReserva,
					tt.reserva.DataPasseio,
					tt.reserva.QuantidadePessoas,
					tt.reserva.ValorTotal,
				).
				WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

			err := repo.Create(tt.reserva)
			if (err != nil) != tt.hasError {
				t.Errorf("Create() error = %v, hasError = %v", err, tt.hasError)
			}
		})
	}
}

func TestReservaRepository_GetByID(t *testing.T) {
	db, mock := setupMockDB(t)
	defer mock.ExpectationsWereMet()

	repo := repository.ReservaRepositoryNew(db)

	tests := []struct {
		name     string
		id       int
		mockRows *sqlmock.Rows
		expected *model.Reserva
		hasError bool
	}{
		{
			name: "Valid reservation ID",
			id:   1,
			mockRows: sqlmock.NewRows([]string{"id", "cliente_id", "empresa_id", "status", "data_reserva", "data_passeio", "quantidade_pessoas", "valor_total", "momento_criacao"}).
				AddRow(1, 1, 1, "pendente", time.Now(), time.Now().AddDate(0, 0, 7), 1, 150.50, time.Now()),
			expected: &model.Reserva{
				ID:                 1,
				ClienteID:          1,
				EmpresaID:          1,
				Status:             "pendente",
				DataReserva:        time.Now(),
				DataPasseio:        time.Now().AddDate(0, 0, 7),
				QuantidadePessoas:  1,
				ValorTotal:         150.50,
			},
			hasError: false,
		},
		{
			name:     "Non-existent reservation ID",
			id:       999,
			mockRows: sqlmock.NewRows([]string{"id", "cliente_id", "empresa_id", "status", "data_reserva", "data_passeio", "quantidade_pessoas", "valor_total", "momento_criacao"}),
			expected: nil,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectQuery(`SELECT`).
				WithArgs(tt.id).
				WillReturnRows(tt.mockRows)

			result, err := repo.GetByID(tt.id)
			if (err != nil) != tt.hasError {
				t.Errorf("GetByID() error = %v, hasError = %v", err, tt.hasError)
			}

			if !tt.hasError && result != nil {
				if result.ID != tt.expected.ID {
					t.Errorf("GetByID() ID = %d, expected %d", result.ID, tt.expected.ID)
				}
				if result.Status != tt.expected.Status {
					t.Errorf("GetByID() Status = %s, expected %s", result.Status, tt.expected.Status)
				}
			}
		})
	}
}

func TestReservaRepository_GetByClienteID(t *testing.T) {
	db, mock := setupMockDB(t)
	defer mock.ExpectationsWereMet()

	repo := repository.ReservaRepositoryNew(db)

	tests := []struct {
		name      string
		clienteID int
		page      int
		limit     int
		mockRows  *sqlmock.Rows
		expected  int
		hasError  bool
	}{
		{
			name:      "Valid client ID",
			clienteID: 1,
			page:      1,
			limit:     10,
			mockRows: sqlmock.NewRows([]string{"id", "cliente_id", "empresa_id", "status", "data_reserva", "data_passeio", "quantidade_pessoas", "valor_total", "momento_criacao"}).
				AddRow(1, 1, 1, "pendente", time.Now(), time.Now().AddDate(0, 0, 7), 1, 150.50, time.Now()).
				AddRow(2, 1, 1, "confirmada", time.Now(), time.Now().AddDate(0, 0, 14), 2, 300.00, time.Now()),
			expected: 2,
			hasError: false,
		},
		{
			name:      "Client with no reservations",
			clienteID: 999,
			page:      1,
			limit:     10,
			mockRows:  sqlmock.NewRows([]string{"id", "cliente_id", "empresa_id", "status", "data_reserva", "data_passeio", "quantidade_pessoas", "valor_total", "momento_criacao"}),
			expected:  0,
			hasError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock para COUNT
			mock.ExpectQuery(`SELECT count`).
				WithArgs(tt.clienteID).
				WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(tt.expected))

			// Mock para SELECT principal
			mock.ExpectQuery(`SELECT`).
				WithArgs(tt.clienteID).
				WillReturnRows(tt.mockRows)

			reservas, total, err := repo.GetByClienteID(tt.clienteID, tt.page, tt.limit)
			if (err != nil) != tt.hasError {
				t.Errorf("GetByClienteID() error = %v, hasError = %v", err, tt.hasError)
			}

			if !tt.hasError {
				if len(reservas) != tt.expected {
					t.Errorf("GetByClienteID() returned %d reservations, expected %d", len(reservas), tt.expected)
				}
				if total != int64(tt.expected) {
					t.Errorf("GetByClienteID() total = %d, expected %d", total, tt.expected)
				}
			}
		})
	}
}

func TestReservaRepository_GetByEmpresaID(t *testing.T) {
	db, mock := setupMockDB(t)
	defer mock.ExpectationsWereMet()

	repo := repository.ReservaRepositoryNew(db)

	tests := []struct {
		name      string
		empresaID int
		page      int
		limit     int
		mockRows  *sqlmock.Rows
		expected  int
		hasError  bool
	}{
		{
			name:      "Valid empresa ID",
			empresaID: 1,
			page:      1,
			limit:     10,
			mockRows: sqlmock.NewRows([]string{"id", "cliente_id", "empresa_id", "status", "data_reserva", "data_passeio", "quantidade_pessoas", "valor_total", "momento_criacao"}).
				AddRow(1, 1, 1, "pendente", time.Now(), time.Now().AddDate(0, 0, 7), 1, 150.50, time.Now()).
				AddRow(2, 2, 1, "confirmada", time.Now(), time.Now().AddDate(0, 0, 14), 2, 300.00, time.Now()),
			expected: 2,
			hasError: false,
		},
		{
			name:      "Empresa with no reservations",
			empresaID: 999,
			page:      1,
			limit:     10,
			mockRows:  sqlmock.NewRows([]string{"id", "cliente_id", "empresa_id", "status", "data_reserva", "data_passeio", "quantidade_pessoas", "valor_total", "momento_criacao"}),
			expected:  0,
			hasError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock para COUNT
			mock.ExpectQuery(`SELECT count`).
				WithArgs(tt.empresaID).
				WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(tt.expected))

			// Mock para SELECT principal
			mock.ExpectQuery(`SELECT`).
				WithArgs(tt.empresaID).
				WillReturnRows(tt.mockRows)

			reservas, total, err := repo.GetByEmpresaID(tt.empresaID, tt.page, tt.limit)
			if (err != nil) != tt.hasError {
				t.Errorf("GetByEmpresaID() error = %v, hasError = %v", err, tt.hasError)
			}

			if !tt.hasError {
				if len(reservas) != tt.expected {
					t.Errorf("GetByEmpresaID() returned %d reservations, expected %d", len(reservas), tt.expected)
				}
				if total != int64(tt.expected) {
					t.Errorf("GetByEmpresaID() total = %d, expected %d", total, tt.expected)
				}
			}
		})
	}
}
