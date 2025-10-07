package repository

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jampa_trip/internal/model"
	"github.com/jampa_trip/internal/repository"
)

// setupMockDB is already defined in company_test.go

func TestFeedbackRepositoryNew(t *testing.T) {
	db, mock := setupMockDB(t)
	defer mock.ExpectationsWereMet()

	repo := repository.FeedbackRepositoryNew(db)
	if repo == nil {
		t.Error("FeedbackRepositoryNew() returned nil")
	}
}

func TestFeedbackRepository_Create(t *testing.T) {
	db, mock := setupMockDB(t)
	defer mock.ExpectationsWereMet()

	repo := repository.FeedbackRepositoryNew(db)

	tests := []struct {
		name     string
		feedback *model.Feedback
		hasError bool
	}{
		{
			name: "Valid feedback",
			feedback: &model.Feedback{
				ClienteID:      1,
				EmpresaID:      1,
				ReservaID:      1,
				Nota:           5,
				Comentario:     "Excelente tour!",
				Status:         "ativo",
				MomentoCriacao: time.Now(),
			},
			hasError: false,
		},
		{
			name: "Feedback with empty comment",
			feedback: &model.Feedback{
				ClienteID:      1,
				EmpresaID:      1,
				ReservaID:      1,
				Nota:           4,
				Comentario:     "",
				Status:         "ativo",
				MomentoCriacao: time.Now(),
			},
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectQuery(`INSERT INTO feedbacks`).
				WithArgs(
					tt.feedback.ClienteID,
					tt.feedback.EmpresaID,
					tt.feedback.ReservaID,
					tt.feedback.Nota,
					tt.feedback.Comentario,
					tt.feedback.Status,
				).
				WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

			err := repo.Create(tt.feedback)
			if (err != nil) != tt.hasError {
				t.Errorf("Create() error = %v, hasError = %v", err, tt.hasError)
			}
		})
	}
}

func TestFeedbackRepository_GetByID(t *testing.T) {
	db, mock := setupMockDB(t)
	defer mock.ExpectationsWereMet()

	repo := repository.FeedbackRepositoryNew(db)

	tests := []struct {
		name     string
		id       int
		mockRows *sqlmock.Rows
		expected *model.Feedback
		hasError bool
	}{
		{
			name: "Valid feedback ID",
			id:   1,
			mockRows: sqlmock.NewRows([]string{"id", "cliente_id", "empresa_id", "reserva_id", "nota", "comentario", "status", "momento_criacao"}).
				AddRow(1, 1, 1, 1, 5, "Excelente tour!", "ativo", time.Now()),
			expected: &model.Feedback{
				ID:         1,
				ClienteID:  1,
				EmpresaID:  1,
				ReservaID:  1,
				Nota:       5,
				Comentario: "Excelente tour!",
				Status:     "ativo",
			},
			hasError: false,
		},
		{
			name:     "Non-existent feedback ID",
			id:       999,
			mockRows: sqlmock.NewRows([]string{"id", "cliente_id", "empresa_id", "reserva_id", "nota", "comentario", "status", "momento_criacao"}),
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
				if result.Nota != tt.expected.Nota {
					t.Errorf("GetByID() Nota = %d, expected %d", result.Nota, tt.expected.Nota)
				}
			}
		})
	}
}

func TestFeedbackRepository_GetByEmpresaID(t *testing.T) {
	db, mock := setupMockDB(t)
	defer mock.ExpectationsWereMet()

	repo := repository.FeedbackRepositoryNew(db)

	tests := []struct {
		name     string
		empresaID int
		page     int
		limit    int
		mockRows *sqlmock.Rows
		expected int
		hasError bool
	}{
		{
			name:     "Valid empresa ID",
			empresaID: 1,
			page:     1,
			limit:    10,
			mockRows: sqlmock.NewRows([]string{"id", "cliente_id", "empresa_id", "reserva_id", "nota", "comentario", "status", "momento_criacao"}).
				AddRow(1, 1, 1, 1, 5, "Excelente tour!", "ativo", time.Now()).
				AddRow(2, 2, 1, 2, 4, "Muito bom!", "ativo", time.Now()),
			expected: 2,
			hasError: false,
		},
		{
			name:     "Empresa with no feedback",
			empresaID: 999,
			page:     1,
			limit:    10,
			mockRows: sqlmock.NewRows([]string{"id", "cliente_id", "empresa_id", "reserva_id", "nota", "comentario", "status", "momento_criacao"}),
			expected: 0,
			hasError: false,
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

			feedbacks, total, err := repo.GetByEmpresaID(tt.empresaID, tt.page, tt.limit)
			if (err != nil) != tt.hasError {
				t.Errorf("GetByEmpresaID() error = %v, hasError = %v", err, tt.hasError)
			}

			if !tt.hasError {
				if len(feedbacks) != tt.expected {
					t.Errorf("GetByEmpresaID() returned %d feedbacks, expected %d", len(feedbacks), tt.expected)
				}
				if total != int64(tt.expected) {
					t.Errorf("GetByEmpresaID() total = %d, expected %d", total, tt.expected)
				}
			}
		})
	}
}
