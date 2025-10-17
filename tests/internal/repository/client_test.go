package repository

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jampa_trip/internal/model"
	"github.com/jampa_trip/internal/repository"
)

// setupMockDB is already defined in the file

func TestClientRepositoryNew(t *testing.T) {
	db, mock := setupMockDB(t)
	defer mock.ExpectationsWereMet()

	repo := repository.ClientRepositoryNew(db)
	if repo == nil {
		t.Errorf("ClientRepositoryNew() returned nil")
	}
	if repo.DB == nil {
		t.Errorf("ClientRepositoryNew() returned repository with nil DB")
	}
}

func TestClientRepository_GetByID(t *testing.T) {
	db, mock := setupMockDB(t)
	defer mock.ExpectationsWereMet()

	repo := repository.ClientRepositoryNew(db)

	tests := []struct {
		name     string
		id       int
		mockRows *sqlmock.Rows
		expected *model.Client
		hasError bool
	}{
		{
			name: "Valid client ID",
			id:   1,
			mockRows: sqlmock.NewRows([]string{"id", "name", "email", "password", "cpf", "phone", "birth_date", "created_at", "updated_at"}).
				AddRow(1, "João Silva", "joao@example.com", "hashed_password", "12345678901", "11999999999", time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC), time.Now(), time.Now()),
			expected: &model.Client{
				ID:        1,
				Name:      "João Silva",
				Email:     "joao@example.com",
				Password:  "hashed_password",
				CPF:       "12345678901",
				Phone:     "11999999999",
				BirthDate: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			hasError: false,
		},
		{
			name:     "Non-existent client ID",
			id:       999,
			mockRows: sqlmock.NewRows([]string{"id", "name", "email", "password", "cpf", "phone", "birth_date", "created_at", "updated_at"}),
			expected: nil,
			hasError: true,
		},
		{
			name:     "Zero client ID",
			id:       0,
			mockRows: sqlmock.NewRows([]string{"id", "name", "email", "password", "cpf", "phone", "birth_date", "created_at", "updated_at"}),
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
				if result.Name != tt.expected.Name {
					t.Errorf("GetByID() Name = %s, expected %s", result.Name, tt.expected.Name)
				}
				if result.Email != tt.expected.Email {
					t.Errorf("GetByID() Email = %s, expected %s", result.Email, tt.expected.Email)
				}
			}
		})
	}
}

func TestClientRepository_GetByEmail(t *testing.T) {
	db, mock := setupMockDB(t)
	defer mock.ExpectationsWereMet()

	repo := repository.ClientRepositoryNew(db)

	tests := []struct {
		name     string
		email    string
		mockRows *sqlmock.Rows
		expected *model.Client
		hasError bool
	}{
		{
			name:  "Valid email",
			email: "joao@example.com",
			mockRows: sqlmock.NewRows([]string{"id", "name", "email", "password", "cpf", "phone", "birth_date", "created_at", "updated_at"}).
				AddRow(1, "João Silva", "joao@example.com", "hashed_password", "12345678901", "11999999999", time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC), time.Now(), time.Now()),
			expected: &model.Client{
				ID:        1,
				Name:      "João Silva",
				Email:     "joao@example.com",
				Password:  "hashed_password",
				CPF:       "12345678901",
				Phone:     "11999999999",
				BirthDate: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			hasError: false,
		},
		{
			name:     "Non-existent email",
			email:    "nonexistent@example.com",
			mockRows: sqlmock.NewRows([]string{"id", "name", "email", "password", "cpf", "phone", "birth_date", "created_at", "updated_at"}),
			expected: nil,
			hasError: true,
		},
		{
			name:     "Empty email",
			email:    "",
			mockRows: sqlmock.NewRows([]string{"id", "name", "email", "password", "cpf", "phone", "birth_date", "created_at", "updated_at"}),
			expected: nil,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectQuery(`SELECT`).
				WithArgs(tt.email).
				WillReturnRows(tt.mockRows)

			result, err := repo.GetByEmail(tt.email)
			if (err != nil) != tt.hasError {
				t.Errorf("GetByEmail() error = %v, hasError = %v", err, tt.hasError)
			}

			if !tt.hasError && result != nil {
				if result.ID != tt.expected.ID {
					t.Errorf("GetByEmail() ID = %d, expected %d", result.ID, tt.expected.ID)
				}
				if result.Email != tt.expected.Email {
					t.Errorf("GetByEmail() Email = %s, expected %s", result.Email, tt.expected.Email)
				}
			}
		})
	}
}

func TestClientRepository_Create(t *testing.T) {
	db, mock := setupMockDB(t)
	defer mock.ExpectationsWereMet()

	repo := repository.ClientRepositoryNew(db)

	tests := []struct {
		name     string
		client   *model.Client
		mockID   int
		hasError bool
	}{
		{
			name: "Valid client",
			client: &model.Client{
				Name:      "João Silva",
				Email:     "joao@example.com",
				Password:  "hashed_password",
				CPF:       "12345678901",
				Phone:     "11999999999",
				BirthDate: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			mockID:   1,
			hasError: false,
		},
		{
			name: "Client with empty name",
			client: &model.Client{
				Name:      "",
				Email:     "joao@example.com",
				Password:  "hashed_password",
				CPF:       "12345678901",
				Phone:     "11999999999",
				BirthDate: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			mockID:   2,
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectQuery(`INSERT INTO clients`).
				WithArgs(
					tt.client.Name,
					tt.client.Email,
					tt.client.Password,
					tt.client.CPF,
					tt.client.Phone,
					tt.client.BirthDate,
					tt.client.CreatedAt,
					tt.client.UpdatedAt,
				).
				WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(tt.mockID))

			err := repo.Create(tt.client)
			if (err != nil) != tt.hasError {
				t.Errorf("Create() error = %v, hasError = %v", err, tt.hasError)
			}

			if !tt.hasError {
				if tt.client.ID != tt.mockID {
					t.Errorf("Create() ID = %d, expected %d", tt.client.ID, tt.mockID)
				}
			}
		})
	}
}

func TestClientRepository_Update(t *testing.T) {
	db, mock := setupMockDB(t)
	defer mock.ExpectationsWereMet()

	repo := repository.ClientRepositoryNew(db)

	tests := []struct {
		name     string
		id       int
		updates  map[string]interface{}
		hasError bool
	}{
		{
			name: "Update only name",
			id:   1,
			updates: map[string]interface{}{
				"name":       "João Silva Updated",
				"updated_at": time.Now(),
			},
			hasError: false,
		},
		{
			name: "Update only email",
			id:   1,
			updates: map[string]interface{}{
				"email":      "joao.updated@example.com",
				"updated_at": time.Now(),
			},
			hasError: false,
		},
		{
			name: "Update multiple fields",
			id:   1,
			updates: map[string]interface{}{
				"name":       "João Silva Updated",
				"email":      "joao.updated@example.com",
				"phone":      "11988888888",
				"updated_at": time.Now(),
			},
			hasError: false,
		},
		{
			name: "Update with password",
			id:   1,
			updates: map[string]interface{}{
				"password":   "new_hashed_password",
				"updated_at": time.Now(),
			},
			hasError: false,
		},
		{
			name: "Update non-existent client",
			id:   999,
			updates: map[string]interface{}{
				"name":       "Non-existent",
				"updated_at": time.Now(),
			},
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectBegin()
			mock.ExpectExec(`UPDATE "clients"`).
				WillReturnResult(sqlmock.NewResult(0, 1))
			mock.ExpectCommit()

			err := repo.Update(tt.id, tt.updates)
			if (err != nil) != tt.hasError {
				t.Errorf("Update() error = %v, hasError = %v", err, tt.hasError)
			}
		})
	}
}

func TestClientRepository_List(t *testing.T) {
	db, mock := setupMockDB(t)
	defer mock.ExpectationsWereMet()

	repo := repository.ClientRepositoryNew(db)

	tests := []struct {
		name     string
		filtros  *model.Client
		mockRows *sqlmock.Rows
		expected int
		hasError bool
	}{
		{
			name: "List all clients",
			filtros: &model.Client{
				Name:  "",
				Email: "",
				CPF:   "",
				Phone: "",
			},
			mockRows: sqlmock.NewRows([]string{"id", "name", "email", "cpf", "phone", "birth_date", "created_at", "updated_at"}).
				AddRow(1, "João Silva", "joao@example.com", "12345678901", "11999999999", time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC), time.Now(), time.Now()).
				AddRow(2, "Maria Santos", "maria@example.com", "98765432100", "11888888888", time.Date(1985, 5, 15, 0, 0, 0, 0, time.UTC), time.Now(), time.Now()),
			expected: 2,
			hasError: false,
		},
		{
			name: "List clients with name filter",
			filtros: &model.Client{
				Name:  "João",
				Email: "",
				CPF:   "",
				Phone: "",
			},
			mockRows: sqlmock.NewRows([]string{"id", "name", "email", "cpf", "phone", "birth_date", "created_at", "updated_at"}).
				AddRow(1, "João Silva", "joao@example.com", "12345678901", "11999999999", time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC), time.Now(), time.Now()),
			expected: 1,
			hasError: false,
		},
		{
			name: "List clients with email filter",
			filtros: &model.Client{
				Name:  "",
				Email: "joao@example.com",
				CPF:   "",
				Phone: "",
			},
			mockRows: sqlmock.NewRows([]string{"id", "name", "email", "cpf", "phone", "birth_date", "created_at", "updated_at"}).
				AddRow(1, "João Silva", "joao@example.com", "12345678901", "11999999999", time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC), time.Now(), time.Now()),
			expected: 1,
			hasError: false,
		},
		{
			name: "Empty result",
			filtros: &model.Client{
				Name:  "NonExistent",
				Email: "",
				CPF:   "",
				Phone: "",
			},
			mockRows: sqlmock.NewRows([]string{"id", "name", "email", "cpf", "phone", "birth_date", "created_at", "updated_at"}),
			expected: 0,
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectQuery(`SELECT`).
				WithArgs(
					tt.filtros.Name, tt.filtros.Name,
					tt.filtros.Email, tt.filtros.Email,
					tt.filtros.CPF, tt.filtros.CPF,
					tt.filtros.Phone, tt.filtros.Phone,
					tt.filtros.BirthDate, tt.filtros.BirthDate,
				).
				WillReturnRows(tt.mockRows)

			result, err := repo.List(tt.filtros)
			if (err != nil) != tt.hasError {
				t.Errorf("List() error = %v, hasError = %v", err, tt.hasError)
			}

			if !tt.hasError {
				if len(result) != tt.expected {
					t.Errorf("List() returned %d clients, expected %d", len(result), tt.expected)
				}
			}
		})
	}
}

func TestClientRepository_EmailExiste(t *testing.T) {
	db, mock := setupMockDB(t)
	defer mock.ExpectationsWereMet()

	repo := repository.ClientRepositoryNew(db)

	tests := []struct {
		name      string
		email     string
		mockCount int64
		expected  bool
		hasError  bool
	}{
		{
			name:      "Email exists",
			email:     "joao@example.com",
			mockCount: 1,
			expected:  true,
			hasError:  false,
		},
		{
			name:      "Email does not exist",
			email:     "nonexistent@example.com",
			mockCount: 0,
			expected:  false,
			hasError:  false,
		},
		{
			name:      "Empty email",
			email:     "",
			mockCount: 0,
			expected:  false,
			hasError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectQuery(`SELECT count`).
				WithArgs(tt.email).
				WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(tt.mockCount))

			result, err := repo.EmailExiste(tt.email)
			if (err != nil) != tt.hasError {
				t.Errorf("EmailExiste() error = %v, hasError = %v", err, tt.hasError)
			}

			if !tt.hasError {
				if result != tt.expected {
					t.Errorf("EmailExiste() = %v, expected %v", result, tt.expected)
				}
			}
		})
	}
}

func TestClientRepository_EmailExisteParaOutroCliente(t *testing.T) {
	db, mock := setupMockDB(t)
	defer mock.ExpectationsWereMet()

	repo := repository.ClientRepositoryNew(db)

	tests := []struct {
		name      string
		email     string
		id        int
		mockCount int64
		expected  bool
		hasError  bool
	}{
		{
			name:      "Email exists for another client",
			email:     "joao@example.com",
			id:        1,
			mockCount: 1,
			expected:  true,
			hasError:  false,
		},
		{
			name:      "Email does not exist for another client",
			email:     "joao@example.com",
			id:        1,
			mockCount: 0,
			expected:  false,
			hasError:  false,
		},
		{
			name:      "Same client ID",
			email:     "joao@example.com",
			id:        1,
			mockCount: 0,
			expected:  false,
			hasError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectQuery(`SELECT count`).
				WithArgs(tt.email, tt.id).
				WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(tt.mockCount))

			result, err := repo.EmailExisteParaOutroCliente(tt.email, tt.id)
			if (err != nil) != tt.hasError {
				t.Errorf("EmailExisteParaOutroCliente() error = %v, hasError = %v", err, tt.hasError)
			}

			if !tt.hasError {
				if result != tt.expected {
					t.Errorf("EmailExisteParaOutroCliente() = %v, expected %v", result, tt.expected)
				}
			}
		})
	}
}
