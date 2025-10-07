package util

import (
	"testing"

	"github.com/jampa_trip/pkg/util"
)

func TestGenerateToken(t *testing.T) {
	// Test multiple token generation to ensure uniqueness
	tokens := make(map[string]bool)

	for i := 0; i < 100; i++ {
		token, err := util.GenerateToken()
		if err != nil {
			t.Errorf("GenerateToken() failed: %v", err)
		}

		// Check token length (should be 64 characters for hex encoded 32 bytes)
		if len(token) != 64 {
			t.Errorf("GenerateToken() returned token with length %d, expected 64", len(token))
		}

		// Check for uniqueness
		if tokens[token] {
			t.Errorf("GenerateToken() returned duplicate token: %s", token)
		}
		tokens[token] = true
	}
}

func TestCriptografarSenha(t *testing.T) {
	tests := []struct {
		name     string
		senha    string
		expected error
	}{
		{
			name:     "Valid password",
			senha:    "password123",
			expected: nil,
		},
		{
			name:     "Empty password",
			senha:    "",
			expected: nil,
		},
		{
			name:     "Long password",
			senha:    "this-is-a-very-long-password-with-many-characters",
			expected: nil,
		},
		{
			name:     "Password with special characters",
			senha:    "P@ssw0rd!@#$%",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := util.CriptografarSenha(tt.senha)
			if (err == nil) != (tt.expected == nil) {
				t.Errorf("CriptografarSenha(%s) error = %v, expected %v", tt.senha, err, tt.expected)
			}

			if err == nil {
				// Check that hash is not empty
				if hash == "" {
					t.Errorf("CriptografarSenha(%s) returned empty hash", tt.senha)
				}

				// Check that hash is different from original password
				if hash == tt.senha {
					t.Errorf("CriptografarSenha(%s) returned same value as input", tt.senha)
				}

				// Check that same password generates different hashes (due to salt)
				hash2, err2 := util.CriptografarSenha(tt.senha)
				if err2 != nil {
					t.Errorf("CriptografarSenha(%s) second call failed: %v", tt.senha, err2)
				}
				if hash == hash2 {
					t.Errorf("CriptografarSenha(%s) returned same hash for different calls", tt.senha)
				}
			}
		})
	}
}

func TestVerificaSenha(t *testing.T) {
	// Test with a known password
	password := "testpassword123"
	hash, err := util.CriptografarSenha(password)
	if err != nil {
		t.Fatalf("CriptografarSenha failed: %v", err)
	}

	tests := []struct {
		name           string
		password       string
		hashedPassword string
		expected       bool
	}{
		{
			name:           "Correct password",
			password:       password,
			hashedPassword: hash,
			expected:       true,
		},
		{
			name:           "Wrong password",
			password:       "wrongpassword",
			hashedPassword: hash,
			expected:       false,
		},
		{
			name:           "Empty password",
			password:       "",
			hashedPassword: hash,
			expected:       false,
		},
		{
			name:           "Empty hash",
			password:       password,
			hashedPassword: "",
			expected:       false,
		},
		{
			name:           "Both empty",
			password:       "",
			hashedPassword: "",
			expected:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := util.VerificaSenha(tt.password, tt.hashedPassword)
			if result != tt.expected {
				t.Errorf("VerificaSenha(%s, %s) = %v, expected %v", tt.password, tt.hashedPassword, result, tt.expected)
			}
		})
	}
}

func TestVerificaSenhaWithDifferentPasswords(t *testing.T) {
	// Test with multiple different passwords
	passwords := []string{
		"password1",
		"P@ssw0rd!",
		"123456",
		"",
		"very-long-password-with-many-characters",
	}

	for _, password := range passwords {
		hash, err := util.CriptografarSenha(password)
		if err != nil {
			t.Errorf("CriptografarSenha(%s) failed: %v", password, err)
			continue
		}

		// Verify correct password works
		if !util.VerificaSenha(password, hash) {
			t.Errorf("VerificaSenha(%s, hash) should return true", password)
		}

		// Verify wrong passwords don't work
		wrongPasswords := []string{
			password + "wrong",
			"wrong" + password,
			"completely different",
		}

		for _, wrongPassword := range wrongPasswords {
			if util.VerificaSenha(wrongPassword, hash) {
				t.Errorf("VerificaSenha(%s, hash) should return false for wrong password", wrongPassword)
			}
		}
	}
}
