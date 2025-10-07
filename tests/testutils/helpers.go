package testutils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// CreateTestRequest creates a test HTTP request with optional body
func CreateTestRequest(t *testing.T, method, url string, body interface{}) *http.Request {
	t.Helper()

	var bodyReader io.Reader
	if body != nil {
		switch v := body.(type) {
		case string:
			bodyReader = bytes.NewBufferString(v)
		case []byte:
			bodyReader = bytes.NewBuffer(v)
		default:
			jsonData, err := json.Marshal(body)
			if err != nil {
				t.Fatalf("Failed to marshal request body: %v", err)
			}
			bodyReader = bytes.NewBuffer(jsonData)
		}
	}

	req := httptest.NewRequest(method, url, bodyReader)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req
}

// CreateJSONRequest creates a test HTTP request with JSON body
func CreateJSONRequest(t *testing.T, method, url string, body interface{}) *http.Request {
	t.Helper()

	jsonData, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	req := httptest.NewRequest(method, url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	return req
}

// ParseJSONResponse parses JSON response from httptest.ResponseRecorder
func ParseJSONResponse(t *testing.T, rec *httptest.ResponseRecorder, target interface{}) {
	t.Helper()

	if err := json.Unmarshal(rec.Body.Bytes(), target); err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}
}

// AssertJSONResponse asserts that the response matches the expected JSON
func AssertJSONResponse(t *testing.T, rec *httptest.ResponseRecorder, expected interface{}) {
	t.Helper()

	expectedJSON, err := json.Marshal(expected)
	if err != nil {
		t.Fatalf("Failed to marshal expected JSON: %v", err)
	}

	actualJSON := rec.Body.Bytes()

	// Compare JSON by unmarshaling and remarshaling to normalize formatting
	var expectedMap, actualMap map[string]interface{}
	if err := json.Unmarshal(expectedJSON, &expectedMap); err != nil {
		t.Fatalf("Failed to unmarshal expected JSON: %v", err)
	}
	if err := json.Unmarshal(actualJSON, &actualMap); err != nil {
		t.Fatalf("Failed to unmarshal actual JSON: %v", err)
	}

	expectedNormalized, _ := json.Marshal(expectedMap)
	actualNormalized, _ := json.Marshal(actualMap)

	if string(expectedNormalized) != string(actualNormalized) {
		t.Errorf("JSON response mismatch:\nExpected: %s\nActual: %s", expectedNormalized, actualNormalized)
	}
}

// AssertStatusCode asserts that the response has the expected status code
func AssertStatusCode(t *testing.T, rec *httptest.ResponseRecorder, expected int) {
	t.Helper()

	if rec.Code != expected {
		t.Errorf("Status code mismatch: got %d, expected %d", rec.Code, expected)
	}
}

// AssertNoError asserts that the error is nil
func AssertNoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

// AssertError asserts that the error is not nil
func AssertError(t *testing.T, err error) {
	t.Helper()

	if err == nil {
		t.Error("Expected an error, got nil")
	}
}

// AssertEqual asserts that two values are equal
func AssertEqual(t *testing.T, actual, expected interface{}) {
	t.Helper()

	if actual != expected {
		t.Errorf("Values not equal:\nActual: %v\nExpected: %v", actual, expected)
	}
}

// AssertNotEqual asserts that two values are not equal
func AssertNotEqual(t *testing.T, actual, expected interface{}) {
	t.Helper()

	if actual == expected {
		t.Errorf("Values should not be equal: %v", actual)
	}
}

// AssertNil asserts that a value is nil
func AssertNil(t *testing.T, value interface{}) {
	t.Helper()

	if value != nil {
		t.Errorf("Expected nil, got: %v", value)
	}
}

// AssertNotNil asserts that a value is not nil
func AssertNotNil(t *testing.T, value interface{}) {
	t.Helper()

	if value == nil {
		t.Error("Expected non-nil value, got nil")
	}
}

// AssertContains asserts that a string contains a substring
func AssertContains(t *testing.T, haystack, needle string) {
	t.Helper()

	if !contains(haystack, needle) {
		t.Errorf("String %q does not contain %q", haystack, needle)
	}
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || findSubstring(s, substr))
}

// findSubstring finds a substring in a string
func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// MockPQArray converte slice para formato string PostgreSQL
func MockPQArray(items []string) string {
	if len(items) == 0 {
		return "{}"
	}
	return "{" + strings.Join(items, ",") + "}"
}
