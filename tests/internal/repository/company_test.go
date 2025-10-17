package repository

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jampa_trip/internal/model"
	"github.com/jampa_trip/internal/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
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

func TestCompanyRepositoryNew(t *testing.T) {
	db, mock := setupMockDB(t)
	defer mock.ExpectationsWereMet()

	repo := repository.CompanyRepositoryNew(db)
	if repo == nil {
		t.Errorf("CompanyRepositoryNew() returned nil")
	}
	if repo.DB == nil {
		t.Errorf("CompanyRepositoryNew() returned repository with nil DB")
	}
}

func TestCompanyRepository_GetByID(t *testing.T) {
	db, mock := setupMockDB(t)
	defer mock.ExpectationsWereMet()

	repo := repository.CompanyRepositoryNew(db)

	tests := []struct {
		name     string
		id       int
		mockRows *sqlmock.Rows
		expected *model.Company
		hasError bool
	}{
		{
			name: "Valid company ID",
			id:   1,
			mockRows: sqlmock.NewRows([]string{"id", "name", "email", "password", "cnpj", "phone", "address", "created_at", "updated_at"}).
				AddRow(1, "Empresa ABC", "empresa@example.com", "hashed_password", "12345678000195", "11999999999", "Rua A, 123", time.Now(), time.Now()),
			expected: &model.Company{
				ID:       1,
				Name:     "Empresa ABC",
				Email:    "empresa@example.com",
				Password: "hashed_password",
				CNPJ:     "12345678000195",
				Phone:    "11999999999",
				Address:  "Rua A, 123",
			},
			hasError: false,
		},
		{
			name:     "Non-existent company ID",
			id:       999,
			mockRows: sqlmock.NewRows([]string{"id", "name", "email", "password", "cnpj", "phone", "address", "created_at", "updated_at"}),
			expected: nil,
			hasError: true,
		},
		{
			name:     "Zero company ID",
			id:       0,
			mockRows: sqlmock.NewRows([]string{"id", "name", "email", "password", "cnpj", "phone", "address", "created_at", "updated_at"}),
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

func TestCompanyRepository_GetByEmail(t *testing.T) {
	db, mock := setupMockDB(t)
	defer mock.ExpectationsWereMet()

	repo := repository.CompanyRepositoryNew(db)

	tests := []struct {
		name     string
		email    string
		mockRows *sqlmock.Rows
		expected *model.Company
		hasError bool
	}{
		{
			name:  "Valid email",
			email: "empresa@example.com",
			mockRows: sqlmock.NewRows([]string{"id", "name", "email", "password", "cnpj", "phone", "address", "created_at", "updated_at"}).
				AddRow(1, "Empresa ABC", "empresa@example.com", "hashed_password", "12345678000195", "11999999999", "Rua A, 123", time.Now(), time.Now()),
			expected: &model.Company{
				ID:       1,
				Name:     "Empresa ABC",
				Email:    "empresa@example.com",
				Password: "hashed_password",
				CNPJ:     "12345678000195",
				Phone:    "11999999999",
				Address:  "Rua A, 123",
			},
			hasError: false,
		},
		{
			name:     "Non-existent email",
			email:    "nonexistent@example.com",
			mockRows: sqlmock.NewRows([]string{"id", "name", "email", "password", "cnpj", "phone", "address", "created_at", "updated_at"}),
			expected: nil,
			hasError: true,
		},
		{
			name:     "Empty email",
			email:    "",
			mockRows: sqlmock.NewRows([]string{"id", "name", "email", "password", "cnpj", "phone", "address", "created_at", "updated_at"}),
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

func TestCompanyRepository_Create(t *testing.T) {
	db, mock := setupMockDB(t)
	defer mock.ExpectationsWereMet()

	repo := repository.CompanyRepositoryNew(db)

	tests := []struct {
		name     string
		company  *model.Company
		mockID   int
		hasError bool
	}{
		{
			name: "Valid company",
			company: &model.Company{
				Name:      "Empresa ABC",
				Email:     "empresa@example.com",
				Password:  "hashed_password",
				CNPJ:      "12345678000195",
				Phone:     "11999999999",
				Address:   "Rua A, 123",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			mockID:   1,
			hasError: false,
		},
		{
			name: "Company with empty name",
			company: &model.Company{
				Name:      "",
				Email:     "empresa@example.com",
				Password:  "hashed_password",
				CNPJ:      "12345678000195",
				Phone:     "11999999999",
				Address:   "Rua A, 123",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			mockID:   2,
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectQuery(`INSERT INTO companies`).
				WithArgs(
					tt.company.Name,
					tt.company.Email,
					tt.company.Password,
					tt.company.CNPJ,
					tt.company.Phone,
					tt.company.Address,
					tt.company.CreatedAt,
					tt.company.UpdatedAt,
				).
				WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(tt.mockID))

			err := repo.Create(tt.company)
			if (err != nil) != tt.hasError {
				t.Errorf("Create() error = %v, hasError = %v", err, tt.hasError)
			}

			if !tt.hasError {
				if tt.company.ID != tt.mockID {
					t.Errorf("Create() ID = %d, expected %d", tt.company.ID, tt.mockID)
				}
			}
		})
	}
}

func TestCompanyRepository_Update(t *testing.T) {
	db, mock := setupMockDB(t)
	defer mock.ExpectationsWereMet()

	repo := repository.CompanyRepositoryNew(db)

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
				"name":       "Empresa ABC Updated",
				"updated_at": time.Now(),
			},
			hasError: false,
		},
		{
			name: "Update only email",
			id:   1,
			updates: map[string]interface{}{
				"email":      "empresa.updated@example.com",
				"updated_at": time.Now(),
			},
			hasError: false,
		},
		{
			name: "Update multiple fields",
			id:   1,
			updates: map[string]interface{}{
				"name":       "Empresa ABC Updated",
				"email":      "empresa.updated@example.com",
				"phone":      "11988888888",
				"address":    "Rua Nova, 999",
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
			name: "Update non-existent company",
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
			mock.ExpectExec(`UPDATE "companies"`).
				WillReturnResult(sqlmock.NewResult(0, 1))
			mock.ExpectCommit()

			err := repo.Update(tt.id, tt.updates)
			if (err != nil) != tt.hasError {
				t.Errorf("Update() error = %v, hasError = %v", err, tt.hasError)
			}
		})
	}
}

func TestCompanyRepository_List(t *testing.T) {
	db, mock := setupMockDB(t)
	defer mock.ExpectationsWereMet()

	repo := repository.CompanyRepositoryNew(db)

	tests := []struct {
		name     string
		filtros  *model.Company
		mockRows *sqlmock.Rows
		expected int
		hasError bool
	}{
		{
			name: "List all companies",
			filtros: &model.Company{
				Name:    "",
				Email:   "",
				CNPJ:    "",
				Phone:   "",
				Address: "",
			},
			mockRows: sqlmock.NewRows([]string{"id", "name", "email", "cnpj", "phone", "address", "created_at", "updated_at"}).
				AddRow(1, "Empresa ABC", "empresa1@example.com", "12345678000195", "11999999999", "Rua A, 123", time.Now(), time.Now()).
				AddRow(2, "Empresa XYZ", "empresa2@example.com", "98765432000123", "11888888888", "Rua B, 456", time.Now(), time.Now()),
			expected: 2,
			hasError: false,
		},
		{
			name: "List companies with name filter",
			filtros: &model.Company{
				Name:    "ABC",
				Email:   "",
				CNPJ:    "",
				Phone:   "",
				Address: "",
			},
			mockRows: sqlmock.NewRows([]string{"id", "name", "email", "cnpj", "phone", "address", "created_at", "updated_at"}).
				AddRow(1, "Empresa ABC", "empresa1@example.com", "12345678000195", "11999999999", "Rua A, 123", time.Now(), time.Now()),
			expected: 1,
			hasError: false,
		},
		{
			name: "List companies with email filter",
			filtros: &model.Company{
				Name:    "",
				Email:   "empresa1@example.com",
				CNPJ:    "",
				Phone:   "",
				Address: "",
			},
			mockRows: sqlmock.NewRows([]string{"id", "name", "email", "cnpj", "phone", "address", "created_at", "updated_at"}).
				AddRow(1, "Empresa ABC", "empresa1@example.com", "12345678000195", "11999999999", "Rua A, 123", time.Now(), time.Now()),
			expected: 1,
			hasError: false,
		},
		{
			name: "Empty result",
			filtros: &model.Company{
				Name:    "NonExistent",
				Email:   "",
				CNPJ:    "",
				Phone:   "",
				Address: "",
			},
			mockRows: sqlmock.NewRows([]string{"id", "name", "email", "cnpj", "phone", "address", "created_at", "updated_at"}),
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
					tt.filtros.CNPJ, tt.filtros.CNPJ,
					tt.filtros.Phone, tt.filtros.Phone,
					tt.filtros.Address, tt.filtros.Address,
				).
				WillReturnRows(tt.mockRows)

			result, err := repo.List(tt.filtros)
			if (err != nil) != tt.hasError {
				t.Errorf("List() error = %v, hasError = %v", err, tt.hasError)
			}

			if !tt.hasError {
				if len(result) != tt.expected {
					t.Errorf("List() returned %d companies, expected %d", len(result), tt.expected)
				}
			}
		})
	}
}

func TestCompanyRepository_EmailExiste(t *testing.T) {
	db, mock := setupMockDB(t)
	defer mock.ExpectationsWereMet()

	repo := repository.CompanyRepositoryNew(db)

	tests := []struct {
		name      string
		email     string
		mockCount int64
		expected  bool
		hasError  bool
	}{
		{
			name:      "Email exists",
			email:     "empresa@example.com",
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

func TestCompanyRepository_EmailExisteParaOutraEmpresa(t *testing.T) {
	db, mock := setupMockDB(t)
	defer mock.ExpectationsWereMet()

	repo := repository.CompanyRepositoryNew(db)

	tests := []struct {
		name      string
		email     string
		id        int
		mockCount int64
		expected  bool
		hasError  bool
	}{
		{
			name:      "Email exists for another company",
			email:     "empresa@example.com",
			id:        1,
			mockCount: 1,
			expected:  true,
			hasError:  false,
		},
		{
			name:      "Email does not exist for another company",
			email:     "empresa@example.com",
			id:        1,
			mockCount: 0,
			expected:  false,
			hasError:  false,
		},
		{
			name:      "Same company ID",
			email:     "empresa@example.com",
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

			result, err := repo.EmailExisteParaOutraEmpresa(tt.email, tt.id)
			if (err != nil) != tt.hasError {
				t.Errorf("EmailExisteParaOutraEmpresa() error = %v, hasError = %v", err, tt.hasError)
			}

			if !tt.hasError {
				if result != tt.expected {
					t.Errorf("EmailExisteParaOutraEmpresa() = %v, expected %v", result, tt.expected)
				}
			}
		})
	}
}
