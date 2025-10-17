package service

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jampa_trip/internal/contract"
	"github.com/jampa_trip/internal/model"
	"github.com/jampa_trip/internal/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDBForService(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to create GORM DB: %v", err)
	}

	return gormDB, mock
}

func TestClientServiceNew(t *testing.T) {
	db, mock := setupMockDBForService(t)
	defer mock.ExpectationsWereMet()

	service := service.ClientServiceNew(db)
	if service == nil {
		t.Errorf("ClientServiceNew() returned nil")
	}
	if service.ClientRepository == nil {
		t.Errorf("ClientServiceNew() returned service with nil ClientRepository")
	}
}

func TestClientService_Create(t *testing.T) {
	db, mock := setupMockDBForService(t)
	defer mock.ExpectationsWereMet()

	service := service.ClientServiceNew(db)

	tests := []struct {
		name     string
		request  *contract.CreateClientRequest
		hasError bool
	}{
		{
			name: "Valid client creation",
			request: &contract.CreateClientRequest{
				Name:            "João Silva",
				Email:           "joao@example.com",
				Password:        "Password123!",
				ConfirmPassword: "Password123!",
				CPF:             "12345678901",
				Phone:           "11999999999",
				BirthDate:       "1990-01-01",
			},
			hasError: false,
		},
		{
			name: "Password mismatch",
			request: &contract.CreateClientRequest{
				Name:            "João Silva",
				Email:           "joao@example.com",
				Password:        "Password123!",
				ConfirmPassword: "DifferentPassword123!",
				CPF:             "12345678901",
				Phone:           "11999999999",
				BirthDate:       "1990-01-01",
			},
			hasError: true,
		},
		{
			name: "Invalid date format",
			request: &contract.CreateClientRequest{
				Name:            "João Silva",
				Email:           "joao@example.com",
				Password:        "Password123!",
				ConfirmPassword: "Password123!",
				CPF:             "12345678901",
				Phone:           "11999999999",
				BirthDate:       "invalid-date",
			},
			hasError: true,
		},
		{
			name: "Empty name",
			request: &contract.CreateClientRequest{
				Name:            "",
				Email:           "joao@example.com",
				Password:        "Password123!",
				ConfirmPassword: "Password123!",
				CPF:             "12345678901",
				Phone:           "11999999999",
				BirthDate:       "1990-01-01",
			},
			hasError: false, // Should still work, validation happens at handler level
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock email existence check
			mock.ExpectQuery(`SELECT count`).
				WithArgs(tt.request.Email).
				WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

			// Mock client creation
			if !tt.hasError {
				mock.ExpectQuery(`INSERT INTO clients`).
					WithArgs(
						tt.request.Name,
						tt.request.Email,
						sqlmock.AnyArg(), // hashed password
						tt.request.CPF,
						tt.request.Phone,
						sqlmock.AnyArg(), // parsed birth date
						sqlmock.AnyArg(), // created_at
						sqlmock.AnyArg(), // updated_at
					).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			}

			result, err := service.Create(tt.request)
			if (err != nil) != tt.hasError {
				t.Errorf("Create() error = %v, hasError = %v", err, tt.hasError)
			}

			if !tt.hasError && result != nil {
				if result.Message == "" {
					t.Errorf("Create() returned empty message")
				}
				if result.Data.ID == 0 {
					t.Errorf("Create() returned client with zero ID")
				}
				if result.Data.Name != tt.request.Name {
					t.Errorf("Create() returned client with name %s, expected %s", result.Data.Name, tt.request.Name)
				}
			}
		})
	}
}

func TestClientService_Update(t *testing.T) {
	db, mock := setupMockDBForService(t)
	defer mock.ExpectationsWereMet()

	service := service.ClientServiceNew(db)

	// Helper function to create string pointers
	strPtr := func(s string) *string { return &s }

	tests := []struct {
		name     string
		request  *contract.UpdateClientRequest
		hasError bool
	}{
		{
			name: "Valid partial update - only name",
			request: &contract.UpdateClientRequest{
				ID:   1,
				Name: strPtr("João Silva Updated"),
			},
			hasError: false,
		},
		{
			name: "Valid partial update - name and email",
			request: &contract.UpdateClientRequest{
				ID:    1,
				Name:  strPtr("João Silva Updated"),
				Email: strPtr("joao.updated@example.com"),
			},
			hasError: false,
		},
		{
			name: "Valid partial update - with password",
			request: &contract.UpdateClientRequest{
				ID:              1,
				Password:        strPtr("NewPassword123!"),
				ConfirmPassword: strPtr("NewPassword123!"),
			},
			hasError: false,
		},
		{
			name: "Password mismatch",
			request: &contract.UpdateClientRequest{
				ID:              1,
				Password:        strPtr("Password123!"),
				ConfirmPassword: strPtr("DifferentPassword123!"),
			},
			hasError: true,
		},
		{
			name: "Non-existent client",
			request: &contract.UpdateClientRequest{
				ID:   999,
				Name: strPtr("Non-existent"),
			},
			hasError: true,
		},
		{
			name: "Email already exists for another client",
			request: &contract.UpdateClientRequest{
				ID:    1,
				Email: strPtr("existing@example.com"),
			},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock existing client check
			if tt.request.ID != 999 {
				mock.ExpectQuery(`SELECT`).
					WithArgs(tt.request.ID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password", "cpf", "phone", "birth_date", "created_at", "updated_at"}).
						AddRow(tt.request.ID, "Old Name", "old@example.com", "old_hash", "12345678901", "11999999999", time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC), time.Now(), time.Now()))
			} else {
				mock.ExpectQuery(`SELECT`).
					WithArgs(tt.request.ID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password", "cpf", "phone", "birth_date", "created_at", "updated_at"}))
			}

			// Mock email existence check if email is being updated
			if tt.request.Email != nil && !tt.hasError && tt.request.ID != 999 {
				emailCount := 0
				if *tt.request.Email == "existing@example.com" {
					emailCount = 1
				}
				mock.ExpectQuery(`SELECT count`).
					WithArgs(*tt.request.Email, tt.request.ID).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(emailCount))
			}

			// Mock GORM update operation
			if !tt.hasError && tt.request.ID != 999 {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE "clients"`).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()

				// Mock fetching updated client
				mock.ExpectQuery(`SELECT`).
					WithArgs(tt.request.ID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password", "cpf", "phone", "birth_date", "created_at", "updated_at"}).
						AddRow(tt.request.ID, "João Silva Updated", "joao.updated@example.com", "new_hash", "12345678901", "11999999999", time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC), time.Now(), time.Now()))
			}

			result, err := service.Update(tt.request)
			if (err != nil) != tt.hasError {
				t.Errorf("Update() error = %v, hasError = %v", err, tt.hasError)
			}

			if !tt.hasError && result != nil {
				if result.Message == "" {
					t.Errorf("Update() returned empty message")
				}
				if result.Data.ID != tt.request.ID {
					t.Errorf("Update() returned client with ID %d, expected %d", result.Data.ID, tt.request.ID)
				}
			}
		})
	}
}

func TestClientService_List(t *testing.T) {
	db, mock := setupMockDBForService(t)
	defer mock.ExpectationsWereMet()

	service := service.ClientServiceNew(db)

	tests := []struct {
		name     string
		filtros  *model.Client
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
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock the list query
			mock.ExpectQuery(`SELECT`).
				WithArgs(
					tt.filtros.Name, tt.filtros.Name,
					tt.filtros.Email, tt.filtros.Email,
					tt.filtros.CPF, tt.filtros.CPF,
					tt.filtros.Phone, tt.filtros.Phone,
					tt.filtros.BirthDate, tt.filtros.BirthDate,
				).
				WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "cpf", "phone", "birth_date", "created_at", "updated_at"}).
					AddRow(1, "João Silva", "joao@example.com", "12345678901", "11999999999", time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC), time.Now(), time.Now()).
					AddRow(2, "Maria Santos", "maria@example.com", "98765432100", "11888888888", time.Date(1985, 5, 15, 0, 0, 0, 0, time.UTC), time.Now(), time.Now()))

			result, err := service.List(tt.filtros)
			if (err != nil) != tt.hasError {
				t.Errorf("List() error = %v, hasError = %v", err, tt.hasError)
			}

			if !tt.hasError && result != nil {
				if result.Message == "" {
					t.Errorf("List() returned empty message")
				}
				if len(result.Data) == 0 {
					t.Errorf("List() returned empty data")
				}
			}
		})
	}
}

func TestClientService_Get(t *testing.T) {
	db, mock := setupMockDBForService(t)
	defer mock.ExpectationsWereMet()

	service := service.ClientServiceNew(db)

	tests := []struct {
		name     string
		id       int
		hasError bool
	}{
		{
			name:     "Valid client ID",
			id:       1,
			hasError: false,
		},
		{
			name:     "Non-existent client ID",
			id:       999,
			hasError: true,
		},
		{
			name:     "Zero client ID",
			id:       0,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.id == 1 {
				// Mock successful client retrieval
				mock.ExpectQuery(`SELECT`).
					WithArgs(tt.id).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password", "cpf", "phone", "birth_date", "created_at", "updated_at"}).
						AddRow(tt.id, "João Silva", "joao@example.com", "hashed_password", "12345678901", "11999999999", time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC), time.Now(), time.Now()))
			} else {
				// Mock client not found
				mock.ExpectQuery(`SELECT`).
					WithArgs(tt.id).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password", "cpf", "phone", "birth_date", "created_at", "updated_at"}))
			}

			result, err := service.Get(tt.id)
			if (err != nil) != tt.hasError {
				t.Errorf("Get() error = %v, hasError = %v", err, tt.hasError)
			}

			if !tt.hasError && result != nil {
				if result.Message == "" {
					t.Errorf("Get() returned empty message")
				}
				if result.Data.ID != tt.id {
					t.Errorf("Get() returned client with ID %d, expected %d", result.Data.ID, tt.id)
				}
			}
		})
	}
}

func TestClientService_EmailValidation(t *testing.T) {
	db, mock := setupMockDBForService(t)
	defer mock.ExpectationsWereMet()

	service := service.ClientServiceNew(db)

	// Test email already exists scenario
	t.Run("Email already exists", func(t *testing.T) {
		request := &contract.CreateClientRequest{
			Name:     "João Silva",
			Email:    "existing@example.com",
			Password: "Password123!",
			// ConfirmPassword field doesn't exist in the struct
			CPF:       "12345678901",
			Phone:     "11999999999",
			BirthDate: "1990-01-01",
		}

		// Mock email exists check
		mock.ExpectQuery(`SELECT count`).
			WithArgs(request.Email).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		result, err := service.Create(request)
		if err == nil {
			t.Errorf("Create() should have failed with existing email")
		}
		if result != nil {
			t.Errorf("Create() should have returned nil result on error")
		}
	})

	// Test email exists for another client scenario
	t.Run("Email exists for another client", func(t *testing.T) {
		strPtr := func(s string) *string { return &s }

		request := &contract.UpdateClientRequest{
			ID:    1,
			Email: strPtr("existing@example.com"),
		}

		// Mock existing client check
		mock.ExpectQuery(`SELECT`).
			WithArgs(request.ID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password", "cpf", "phone", "birth_date", "created_at", "updated_at"}).
				AddRow(request.ID, "Old Name", "old@example.com", "old_hash", "12345678901", "11999999999", time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC), time.Now(), time.Now()))

		// Mock email exists for another client
		mock.ExpectQuery(`SELECT count`).
			WithArgs(*request.Email, request.ID).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		result, err := service.Update(request)
		if err == nil {
			t.Errorf("Update() should have failed with existing email for another client")
		}
		if result != nil {
			t.Errorf("Update() should have returned nil result on error")
		}
	})
}
