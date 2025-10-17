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
	now := time.Now()
	futureDate := now.AddDate(0, 0, 7)

	tests := []struct {
		name     string
		reserva  *model.Reserva
		hasError bool
	}{
		{
			name: "Valid reservation with pending status",
			reserva: &model.Reserva{
				TourID:                 1,
				ClienteID:              1,
				PagamentoID:            0,
				Status:                 "pendente",
				DataReserva:            now,
				DataPasseioSelecionada: futureDate,
				QuantidadePessoas:      2,
				ValorTotal:             150.50,
				Observacoes:            "Teste",
				MomentoCriacao:         now,
				MomentoAtualizacao:     now,
			},
			hasError: false,
		},
		{
			name: "Reservation with confirmed status",
			reserva: &model.Reserva{
				TourID:                 2,
				ClienteID:              1,
				PagamentoID:            1,
				Status:                 "confirmada",
				DataReserva:            now,
				DataPasseioSelecionada: futureDate,
				QuantidadePessoas:      3,
				ValorTotal:             300.00,
				Observacoes:            "",
				MomentoCriacao:         now,
				MomentoAtualizacao:     now,
			},
			hasError: false,
		},
		{
			name: "Reservation with aguardando_pagamento status",
			reserva: &model.Reserva{
				TourID:                 1,
				ClienteID:              2,
				PagamentoID:            0,
				Status:                 "aguardando_pagamento",
				DataReserva:            now,
				DataPasseioSelecionada: futureDate,
				QuantidadePessoas:      1,
				ValorTotal:             100.00,
				Observacoes:            "",
				MomentoCriacao:         now,
				MomentoAtualizacao:     now,
			},
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectBegin()
			mock.ExpectQuery(`INSERT INTO "reservations"`).
				WithArgs(
					tt.reserva.TourID,
					tt.reserva.ClienteID,
					tt.reserva.PagamentoID,
					tt.reserva.Status,
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					tt.reserva.QuantidadePessoas,
					tt.reserva.ValorTotal,
					tt.reserva.Observacoes,
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
					sqlmock.AnyArg(),
				).
				WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			mock.ExpectCommit()

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
	now := time.Now()
	futureDate := now.AddDate(0, 0, 7)

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
			mockRows: sqlmock.NewRows([]string{
				"id", "tour_id", "cliente_id", "pagamento_id", "status",
				"data_reserva", "data_passeio_selecionada", "quantidade_pessoas",
				"valor_total", "observacoes", "momento_criacao", "momento_atualizacao",
				"momento_cancelamento",
			}).AddRow(1, 1, 1, 0, "pendente", now, futureDate, 2, 150.50, "Teste", now, now, nil),
			expected: &model.Reserva{
				ID:                1,
				TourID:            1,
				ClienteID:         1,
				PagamentoID:       0,
				Status:            "pendente",
				QuantidadePessoas: 2,
				ValorTotal:        150.50,
			},
			hasError: false,
		},
		{
			name: "Non-existent reservation ID",
			id:   999,
			mockRows: sqlmock.NewRows([]string{
				"id", "tour_id", "cliente_id", "pagamento_id", "status",
				"data_reserva", "data_passeio_selecionada", "quantidade_pessoas",
				"valor_total", "observacoes", "momento_criacao", "momento_atualizacao",
				"momento_cancelamento",
			}),
			expected: nil,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock para SELECT principal (executada primeiro)
			mock.ExpectQuery(`SELECT \* FROM "reservations"`).
				WithArgs(tt.id, 1).
				WillReturnRows(tt.mockRows)

			if !tt.hasError {
				// Mock para Preloads (ordem real: Cliente, Tour, Pagamento)
				mock.ExpectQuery(`SELECT \* FROM "clients"`).
					WillReturnRows(sqlmock.NewRows([]string{"id"}))

				mock.ExpectQuery(`SELECT \* FROM "tours"`).
					WillReturnRows(sqlmock.NewRows([]string{"id"}))

				mock.ExpectQuery(`SELECT \* FROM "pagamentos"`).
					WillReturnRows(sqlmock.NewRows([]string{"id"}))
			}

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
			name:      "Client with no reservations",
			clienteID: 999,
			page:      1,
			limit:     10,
			mockRows: sqlmock.NewRows([]string{
				"id", "tour_id", "cliente_id", "pagamento_id", "status",
				"data_reserva", "data_passeio_selecionada", "quantidade_pessoas",
				"valor_total", "observacoes", "momento_criacao", "momento_atualizacao",
				"momento_cancelamento",
			}),
			expected: 0,
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock para COUNT
			mock.ExpectQuery(`SELECT count\(\*\) FROM "reservations"`).
				WithArgs(tt.clienteID).
				WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(tt.expected))

			// Mock para SELECT principal
			mock.ExpectQuery(`SELECT \* FROM "reservations"`).
				WithArgs(tt.clienteID, tt.limit).
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

func TestReservaRepository_Update(t *testing.T) {
	db, mock := setupMockDB(t)
	defer mock.ExpectationsWereMet()

	repo := repository.ReservaRepositoryNew(db)
	now := time.Now()
	futureDate := now.AddDate(0, 0, 7)

	reserva := &model.Reserva{
		ID:                     1,
		TourID:                 1,
		ClienteID:              1,
		PagamentoID:            0,
		Status:                 "confirmada",
		DataReserva:            now,
		DataPasseioSelecionada: futureDate,
		QuantidadePessoas:      3,
		ValorTotal:             450.00,
		Observacoes:            "Atualizado",
		MomentoCriacao:         now,
		MomentoAtualizacao:     now,
	}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "reservations"`).
		WithArgs(
			reserva.TourID,
			reserva.ClienteID,
			reserva.PagamentoID,
			reserva.Status,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			reserva.QuantidadePessoas,
			reserva.ValorTotal,
			reserva.Observacoes,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			reserva.ID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Update(reserva)
	if err != nil {
		t.Errorf("Update() error = %v", err)
	}
}

func TestReservaRepository_Cancel(t *testing.T) {
	db, mock := setupMockDB(t)
	defer mock.ExpectationsWereMet()

	repo := repository.ReservaRepositoryNew(db)

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "reservations"`).
		WithArgs(
			sqlmock.AnyArg(), // momento_atualizacao
			sqlmock.AnyArg(), // momento_cancelamento
			"cancelada",      // status
			1,                // id
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Cancel(1)
	if err != nil {
		t.Errorf("Cancel() error = %v", err)
	}
}

func TestReservaRepository_UpdateStatus(t *testing.T) {
	db, mock := setupMockDB(t)
	defer mock.ExpectationsWereMet()

	repo := repository.ReservaRepositoryNew(db)

	tests := []struct {
		name     string
		id       int
		status   string
		hasError bool
	}{
		{
			name:     "Update to confirmada",
			id:       1,
			status:   "confirmada",
			hasError: false,
		},
		{
			name:     "Update to aguardando_pagamento",
			id:       2,
			status:   "aguardando_pagamento",
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectBegin()
			mock.ExpectExec(`UPDATE "reservations"`).
				WithArgs(
					sqlmock.AnyArg(), // momento_atualizacao
					tt.status,        // status
					tt.id,            // id
				).
				WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()

			err := repo.UpdateStatus(tt.id, tt.status)
			if (err != nil) != tt.hasError {
				t.Errorf("UpdateStatus() error = %v, hasError = %v", err, tt.hasError)
			}
		})
	}
}

func TestReservaRepository_GetUpcoming(t *testing.T) {
	db, mock := setupMockDB(t)
	defer mock.ExpectationsWereMet()

	repo := repository.ReservaRepositoryNew(db)
	now := time.Now()
	futureDate := now.AddDate(0, 0, 7)

	clienteID := 1

	// Mock para COUNT
	mock.ExpectQuery(`SELECT count\(\*\) FROM "reservations"`).
		WithArgs(clienteID, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	// Mock para SELECT principal
	mock.ExpectQuery(`SELECT \* FROM "reservations"`).
		WithArgs(clienteID, sqlmock.AnyArg(), 10).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "tour_id", "cliente_id", "pagamento_id", "status",
			"data_reserva", "data_passeio_selecionada", "quantidade_pessoas",
			"valor_total", "observacoes", "momento_criacao", "momento_atualizacao",
			"momento_cancelamento",
		}).AddRow(1, 1, clienteID, 0, "confirmada", now, futureDate, 2, 200.00, "", now, now, nil))

	// Mocks para Preloads (ordem real: Cliente, Tour, Pagamento)
	mock.ExpectQuery(`SELECT \* FROM "clients"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	mock.ExpectQuery(`SELECT \* FROM "tours"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	mock.ExpectQuery(`SELECT \* FROM "pagamentos"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	reservas, total, err := repo.GetUpcoming(clienteID, 1, 10)
	if err != nil {
		t.Errorf("GetUpcoming() error = %v", err)
	}

	if len(reservas) != 1 {
		t.Errorf("GetUpcoming() returned %d reservations, expected 1", len(reservas))
	}

	if total != 1 {
		t.Errorf("GetUpcoming() total = %d, expected 1", total)
	}
}

func TestReservaRepository_GetHistory(t *testing.T) {
	db, mock := setupMockDB(t)
	defer mock.ExpectationsWereMet()

	repo := repository.ReservaRepositoryNew(db)
	now := time.Now()
	pastDate := now.AddDate(0, 0, -7)

	clienteID := 1

	// Mock para COUNT
	mock.ExpectQuery(`SELECT count\(\*\) FROM "reservations"`).
		WithArgs(clienteID, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	// Mock para SELECT principal
	mock.ExpectQuery(`SELECT \* FROM "reservations"`).
		WithArgs(clienteID, sqlmock.AnyArg(), 10).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "tour_id", "cliente_id", "pagamento_id", "status",
			"data_reserva", "data_passeio_selecionada", "quantidade_pessoas",
			"valor_total", "observacoes", "momento_criacao", "momento_atualizacao",
			"momento_cancelamento",
		}).AddRow(1, 1, clienteID, 0, "concluida", now, pastDate, 2, 200.00, "", now, now, nil))

	// Mocks para Preloads (ordem real: Cliente, Tour, Pagamento)
	mock.ExpectQuery(`SELECT \* FROM "clients"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	mock.ExpectQuery(`SELECT \* FROM "tours"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	mock.ExpectQuery(`SELECT \* FROM "pagamentos"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	reservas, total, err := repo.GetHistory(clienteID, 1, 10)
	if err != nil {
		t.Errorf("GetHistory() error = %v", err)
	}

	if len(reservas) != 1 {
		t.Errorf("GetHistory() returned %d reservations, expected 1", len(reservas))
	}

	if total != 1 {
		t.Errorf("GetHistory() total = %d, expected 1", total)
	}
}
