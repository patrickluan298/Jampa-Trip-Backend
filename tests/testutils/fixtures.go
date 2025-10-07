package testutils

import (
	"time"

	"github.com/jampa_trip/internal/model"
	"github.com/lib/pq"
)

// Fixtures provides test data for various models
type Fixtures struct{}

// NewFixtures creates a new Fixtures instance
func NewFixtures() *Fixtures {
	return &Fixtures{}
}

// Client returns a test client
func (f *Fixtures) Client() *model.Client {
	now := time.Now()
	birthDate := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)

	return &model.Client{
		ID:        1,
		Name:      "João Silva",
		Email:     "joao@example.com",
		Password:  "$2a$10$test.hash.password",
		CPF:       "12345678901",
		Phone:     "11999999999",
		BirthDate: birthDate,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Company returns a test company
func (f *Fixtures) Company() *model.Company {
	now := time.Now()

	return &model.Company{
		ID:        1,
		Name:      "Empresa Teste",
		Email:     "empresa@example.com",
		Password:  "$2a$10$test.hash.password",
		CNPJ:      "12345678901234",
		Phone:     "11999999999",
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Tour returns a test tour
func (f *Fixtures) Tour() *model.Tour {
	now := time.Now()
	tourDate := "2024-12-15" // String format for PostgreSQL array

	return &model.Tour{
		ID:            1,
		CompanyID:     1,
		Name:          "Tour Teste",
		Description:   "Descrição do tour teste",
		Price:         100.50,
		Dates:         pq.StringArray{tourDate},
		DepartureTime: "08:00",
		ArrivalTime:   "18:00",
		MaxPeople:     20,
		Images:        pq.StringArray{"https://example.com/image1.jpg"},
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

// ValidClientRequest returns a valid client creation request
func (f *Fixtures) ValidClientRequest() map[string]interface{} {
	return map[string]interface{}{
		"name":             "João Silva",
		"email":            "joao@example.com",
		"password":         "Password123!",
		"confirm_password": "Password123!",
		"cpf":              "12345678901",
		"phone":            "11999999999",
		"birth_date":       "1990-01-01",
	}
}

// ValidCompanyRequest returns a valid company creation request
func (f *Fixtures) ValidCompanyRequest() map[string]interface{} {
	return map[string]interface{}{
		"name":             "Empresa Teste",
		"email":            "empresa@example.com",
		"password":         "Password123!",
		"confirm_password": "Password123!",
		"cnpj":             "12345678901234",
		"phone":            "11999999999",
	}
}

// ValidLoginRequest returns a valid login request
func (f *Fixtures) ValidLoginRequest() map[string]interface{} {
	return map[string]interface{}{
		"email":    "joao@example.com",
		"password": "Password123!",
	}
}

// ValidTourRequest returns a valid tour creation request
func (f *Fixtures) ValidTourRequest() map[string]interface{} {
	return map[string]interface{}{
		"title":       "Tour Teste",
		"description": "Descrição do tour teste",
		"price":       100.50,
		"date":        time.Now().AddDate(0, 1, 0).Format("2006-01-02"),
		"location":    "João Pessoa, PB",
		"max_people":  20,
	}
}

// InvalidEmail returns an invalid email for testing
func (f *Fixtures) InvalidEmail() string {
	return "invalid-email"
}

// InvalidCPF returns an invalid CPF for testing
func (f *Fixtures) InvalidCPF() string {
	return "00000000000"
}

// InvalidCNPJ returns an invalid CNPJ for testing
func (f *Fixtures) InvalidCNPJ() string {
	return "00000000000000"
}

// WeakPassword returns a weak password for testing
func (f *Fixtures) WeakPassword() string {
	return "weak"
}

// StrongPassword returns a strong password for testing
func (f *Fixtures) StrongPassword() string {
	return "StrongPassword123!"
}
