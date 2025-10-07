package service

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jampa_trip/internal/contract"
	"github.com/jampa_trip/internal/service"
	"github.com/jampa_trip/tests/testutils"
)

// TestLoginServiceNew_Refactored demonstrates the refactored test using testutils
func TestLoginServiceNew_Refactored(t *testing.T) {
	// Arrange
	factory := testutils.NewMockFactory(t)
	defer factory.Cleanup()
	defer factory.ExpectationsWereMet()

	// Act
	loginService := service.LoginServiceNew(factory.DB)

	// Assert
	testutils.AssertNotNil(t, loginService)
	testutils.AssertNotNil(t, loginService.CompanyRepository)
	testutils.AssertNotNil(t, loginService.ClientRepository)
}

// TestLoginService_Login_Refactored demonstrates the refactored login test
func TestLoginService_Login_Refactored(t *testing.T) {
	// Setup
	factory := testutils.NewMockFactory(t)
	defer factory.Cleanup()
	defer factory.ExpectationsWereMet()

	loginService := service.LoginServiceNew(factory.DB)

	tests := []struct {
		name      string
		request   *contract.LoginRequest
		setupMock func(mock sqlmock.Sqlmock)
		hasError  bool
	}{
		{
			name: "Valid company login",
			request: &contract.LoginRequest{
				Email:    "empresa@example.com",
				Password: "CompanyPassword123!",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				// Mock company repository call
				mock.ExpectQuery(`SELECT`).
					WithArgs("empresa@example.com").
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password", "cnpj", "phone", "address", "created_at", "updated_at"}).
						AddRow(1, "Empresa ABC", "empresa@example.com", "$2a$10$test.hash.password", "12345678000195", "11999999999", "Rua A, 123", time.Now(), time.Now()))
			},
			hasError: false,
		},
		{
			name: "Valid client login",
			request: &contract.LoginRequest{
				Email:    "cliente@example.com",
				Password: "ClientPassword123!",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				// Mock company repository call (not found)
				mock.ExpectQuery(`SELECT`).
					WithArgs("cliente@example.com").
					WillReturnError(sqlmock.ErrCancelled)

				// Mock client repository call
				mock.ExpectQuery(`SELECT`).
					WithArgs("cliente@example.com").
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password", "cpf", "phone", "birth_date", "created_at", "updated_at"}).
						AddRow(1, "João Silva", "cliente@example.com", "$2a$10$test.hash.password", "12345678901", "11999999999", time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC), time.Now(), time.Now()))
			},
			hasError: false,
		},
		{
			name: "Invalid email - user not found",
			request: &contract.LoginRequest{
				Email:    "nonexistent@example.com",
				Password: "Password123!",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				// Mock company repository call (not found)
				mock.ExpectQuery(`SELECT`).
					WithArgs("nonexistent@example.com").
					WillReturnError(sqlmock.ErrCancelled)

				// Mock client repository call (not found)
				mock.ExpectQuery(`SELECT`).
					WithArgs("nonexistent@example.com").
					WillReturnError(sqlmock.ErrCancelled)
			},
			hasError: true,
		},
		{
			name: "Empty email",
			request: &contract.LoginRequest{
				Email:    "",
				Password: "Password123!",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				// No mock needed - validation should fail before DB call
			},
			hasError: true,
		},
		{
			name: "Empty password",
			request: &contract.LoginRequest{
				Email:    "empresa@example.com",
				Password: "",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				// No mock needed - validation should fail before DB call
			},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: Still skipping due to JWT/Redis dependencies
			// TODO: Implement JWT/Redis mocking to enable these tests
			t.Skip("Skipping due to JWT/Redis dependencies - requires proper mocking")

			// Setup mocks
			if tt.setupMock != nil {
				tt.setupMock(factory.SQLMock)
			}

			// Act
			response, err := loginService.Login(tt.request)

			// Assert
			if tt.hasError {
				testutils.AssertError(t, err)
				testutils.AssertNil(t, response)
			} else {
				testutils.AssertNoError(t, err)
				testutils.AssertNotNil(t, response)
				testutils.AssertNotEqual(t, response.AccessToken, "")
				testutils.AssertNotEqual(t, response.RefreshToken, "")
			}
		})
	}
}

// Example of a fully functional test without dependencies
func TestLoginService_ValidateRequest_Refactored(t *testing.T) {
	fixtures := testutils.NewFixtures()

	tests := []struct {
		name     string
		request  *contract.LoginRequest
		hasError bool
	}{
		{
			name:     "Valid request",
			request:  &contract.LoginRequest{Email: "test@example.com", Password: "Password123!"},
			hasError: false,
		},
		{
			name:     "Empty email",
			request:  &contract.LoginRequest{Email: "", Password: "Password123!"},
			hasError: true,
		},
		{
			name:     "Empty password",
			request:  &contract.LoginRequest{Email: "test@example.com", Password: ""},
			hasError: true,
		},
		{
			name:     "Invalid email format",
			request:  &contract.LoginRequest{Email: fixtures.InvalidEmail(), Password: "Password123!"},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			err := tt.request.Validate()

			// Assert
			if tt.hasError {
				testutils.AssertError(t, err)
			} else {
				testutils.AssertNoError(t, err)
			}
		})
	}
}

