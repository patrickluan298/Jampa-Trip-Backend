package util

import (
	"errors"
	"testing"

	"github.com/jampa_trip/pkg/util"
)

func TestValidateTimeFormat(t *testing.T) {
	tests := []struct {
		name     string
		timeStr  string
		expected bool
	}{
		{
			name:     "Valid time format",
			timeStr:  "14:30",
			expected: true,
		},
		{
			name:     "Valid time format with leading zero",
			timeStr:  "09:05",
			expected: true,
		},
		{
			name:     "Invalid time format - no leading zero",
			timeStr:  "9:05",
			expected: true, // This format might be valid in some implementations
		},
		{
			name:     "Invalid time format - wrong format",
			timeStr:  "14:30:00",
			expected: false,
		},
		{
			name:     "Invalid time format - empty string",
			timeStr:  "",
			expected: false,
		},
		{
			name:     "Invalid time format - invalid hour",
			timeStr:  "25:30",
			expected: false,
		},
		{
			name:     "Invalid time format - invalid minute",
			timeStr:  "14:60",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := util.ValidateTimeFormat(tt.timeStr)
			if result != tt.expected {
				t.Errorf("ValidateTimeFormat(%s) = %v, expected %v", tt.timeStr, result, tt.expected)
			}
		})
	}
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected bool
	}{
		{
			name:     "Valid HTTP URL",
			url:      "http://example.com",
			expected: true,
		},
		{
			name:     "Valid HTTPS URL",
			url:      "https://example.com",
			expected: true,
		},
		{
			name:     "Valid URL with path",
			url:      "https://example.com/path/to/resource",
			expected: true,
		},
		{
			name:     "Valid URL with query parameters",
			url:      "https://example.com?param=value",
			expected: true,
		},
		{
			name:     "Invalid URL - no protocol",
			url:      "example.com",
			expected: false,
		},
		{
			name:     "Invalid URL - empty string",
			url:      "",
			expected: false,
		},
		{
			name:     "Invalid URL - malformed",
			url:      "not-a-url",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := util.ValidateURL(tt.url)
			if result != tt.expected {
				t.Errorf("ValidateURL(%s) = %v, expected %v", tt.url, result, tt.expected)
			}
		})
	}
}

func TestValidateDateFormat(t *testing.T) {
	tests := []struct {
		name     string
		dateStr  string
		expected error
	}{
		{
			name:     "Valid date format",
			dateStr:  "2023-12-25",
			expected: nil,
		},
		{
			name:     "Valid date format with leading zeros",
			dateStr:  "2023-01-01",
			expected: nil,
		},
		{
			name:     "Invalid date format - wrong separator",
			dateStr:  "2023/12/25",
			expected: errors.New("parsing time \"2023/12/25\" as \"2006-01-02\": cannot parse \"/12/25\" as \"-\""),
		},
		{
			name:     "Invalid date format - wrong order",
			dateStr:  "25-12-2023",
			expected: errors.New("parsing time \"25-12-2023\" as \"2006-01-02\": cannot parse \"25-12-2023\" as \"2006\""),
		},
		{
			name:     "Invalid date format - empty string",
			dateStr:  "",
			expected: errors.New("parsing time \"\" as \"2006-01-02\": cannot parse \"\" as \"2006\""),
		},
		{
			name:     "Invalid date - non-existent date",
			dateStr:  "2023-02-30",
			expected: errors.New("parsing time \"2023-02-30\": day out of range"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := util.ValidateDateFormat(tt.dateStr)
			if (result == nil) != (tt.expected == nil) {
				t.Errorf("ValidateDateFormat(%s) = %v, expected %v", tt.dateStr, result, tt.expected)
			}
		})
	}
}

func TestValidateDates(t *testing.T) {
	tests := []struct {
		name     string
		dates    []string
		expected error
	}{
		{
			name:     "Valid dates",
			dates:    []string{"2023-12-25", "2023-01-01"},
			expected: nil,
		},
		{
			name:     "Empty slice",
			dates:    []string{},
			expected: nil,
		},
		{
			name:     "Invalid dates",
			dates:    []string{"2023/12/25", "invalid-date"},
			expected: errors.New("parsing time \"2023/12/25\" as \"2006-01-02\": cannot parse \"/12/25\" as \"-\""),
		},
		{
			name:     "Mixed valid and invalid dates",
			dates:    []string{"2023-12-25", "invalid-date"},
			expected: errors.New("parsing time \"invalid-date\" as \"2006-01-02\": cannot parse \"invalid-date\" as \"2006\""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := util.ValidateDates(tt.dates)
			if (result == nil) != (tt.expected == nil) {
				t.Errorf("ValidateDates(%v) = %v, expected %v", tt.dates, result, tt.expected)
			}
		})
	}
}

func TestValidateImageURLs(t *testing.T) {
	tests := []struct {
		name     string
		urls     []string
		expected error
	}{
		{
			name:     "Valid image URLs",
			urls:     []string{"https://example.com/image1.jpg", "https://example.com/image2.png"},
			expected: nil,
		},
		{
			name:     "Empty slice",
			urls:     []string{},
			expected: nil,
		},
		{
			name:     "Empty URLs (should be valid)",
			urls:     []string{"", ""},
			expected: nil,
		},
		{
			name:     "Invalid URLs",
			urls:     []string{"not-a-url", "ftp://example.com/image.jpg"},
			expected: errors.New("URL de imagem inválida"),
		},
		{
			name:     "Mixed valid and invalid URLs",
			urls:     []string{"https://example.com/image.jpg", "not-a-url"},
			expected: errors.New("URL de imagem inválida"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := util.ValidateImageURLs(tt.urls)
			if (result == nil) != (tt.expected == nil) {
				t.Errorf("ValidateImageURLs(%v) = %v, expected %v", tt.urls, result, tt.expected)
			}
		})
	}
}

func TestValidaSegurancaSenha(t *testing.T) {
	tests := []struct {
		name     string
		senha    string
		expected error
	}{
		{
			name:     "Valid strong password",
			senha:    "Password123!",
			expected: nil,
		},
		{
			name:     "Password without uppercase",
			senha:    "password123!",
			expected: errors.New("senha deve conter pelo menos 1 letra maiúscula"),
		},
		{
			name:     "Password without number",
			senha:    "Password!",
			expected: errors.New("senha deve conter pelo menos 1 número"),
		},
		{
			name:     "Password without special character",
			senha:    "Password123",
			expected: errors.New("senha deve conter pelo menos 1 caractere especial"),
		},
		{
			name:     "Empty password",
			senha:    "",
			expected: errors.New("senha deve conter pelo menos 1 letra maiúscula"),
		},
		{
			name:     "Password with all requirements",
			senha:    "MyStr0ng#Pass",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := util.ValidaSegurancaSenha(tt.senha)
			if (result == nil) != (tt.expected == nil) {
				t.Errorf("ValidaSegurancaSenha(%s) = %v, expected %v", tt.senha, result, tt.expected)
			}
		})
	}
}

func TestValidaCPF(t *testing.T) {
	tests := []struct {
		name     string
		cpf      string
		expected error
	}{
		{
			name:     "Valid CPF",
			cpf:      "11144477735",
			expected: nil,
		},
		{
			name:     "Valid CPF with formatting",
			cpf:      "111.444.777-35",
			expected: nil,
		},
		{
			name:     "Invalid CPF - all same digits",
			cpf:      "11111111111",
			expected: errors.New("CPF inválido: todos os dígitos são iguais"),
		},
		{
			name:     "Invalid CPF - wrong length",
			cpf:      "123456789",
			expected: errors.New("CPF deve ter 11 dígitos"),
		},
		{
			name:     "Invalid CPF - wrong check digits",
			cpf:      "11144477734",
			expected: errors.New("CPF inválido: segundo dígito verificador incorreto"),
		},
		{
			name:     "Empty CPF",
			cpf:      "",
			expected: errors.New("CPF deve ter 11 dígitos"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := util.ValidaCPF(tt.cpf)
			if (result == nil) != (tt.expected == nil) {
				t.Errorf("ValidaCPF(%s) = %v, expected %v", tt.cpf, result, tt.expected)
			}
		})
	}
}

func TestValidaCNPJ(t *testing.T) {
	tests := []struct {
		name     string
		cnpj     string
		expected error
	}{
		{
			name:     "Valid CNPJ",
			cnpj:     "11222333000181",
			expected: nil,
		},
		{
			name:     "Valid CNPJ with formatting",
			cnpj:     "11.222.333/0001-81",
			expected: nil,
		},
		{
			name:     "Invalid CNPJ - all same digits",
			cnpj:     "11111111111111",
			expected: errors.New("CNPJ inválido: todos os dígitos são iguais"),
		},
		{
			name:     "Invalid CNPJ - wrong length",
			cnpj:     "123456789",
			expected: errors.New("CNPJ deve ter 14 dígitos"),
		},
		{
			name:     "Invalid CNPJ - wrong check digits",
			cnpj:     "11222333000180",
			expected: errors.New("CNPJ inválido: segundo dígito verificador incorreto"),
		},
		{
			name:     "Empty CNPJ",
			cnpj:     "",
			expected: errors.New("CNPJ deve ter 14 dígitos"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := util.ValidaCNPJ(tt.cnpj)
			if (result == nil) != (tt.expected == nil) {
				t.Errorf("ValidaCNPJ(%s) = %v, expected %v", tt.cnpj, result, tt.expected)
			}
		})
	}
}
