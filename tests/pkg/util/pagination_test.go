package util

import (
	"testing"

	"github.com/jampa_trip/pkg/util"
)

func TestDefaultPagination(t *testing.T) {
	pagination := util.DefaultPagination()

	if pagination.Page != 1 {
		t.Errorf("DefaultPagination() Page = %d, expected 1", pagination.Page)
	}

	if pagination.Limit != 20 {
		t.Errorf("DefaultPagination() Limit = %d, expected 20", pagination.Limit)
	}
}

func TestNormalizePagination(t *testing.T) {
	tests := []struct {
		name          string
		page          int
		limit         int
		expectedPage  int
		expectedLimit int
	}{
		{
			name:          "Valid page and limit",
			page:          5,
			limit:         10,
			expectedPage:  5,
			expectedLimit: 10,
		},
		{
			name:          "Zero page should default to 1",
			page:          0,
			limit:         10,
			expectedPage:  1,
			expectedLimit: 10,
		},
		{
			name:          "Negative page should default to 1",
			page:          -1,
			limit:         10,
			expectedPage:  1,
			expectedLimit: 10,
		},
		{
			name:          "Zero limit should default to 20",
			page:          5,
			limit:         0,
			expectedPage:  5,
			expectedLimit: 20,
		},
		{
			name:          "Negative limit should default to 20",
			page:          5,
			limit:         -1,
			expectedPage:  5,
			expectedLimit: 20,
		},
		{
			name:          "Limit over 100 should be capped at 100",
			page:          5,
			limit:         150,
			expectedPage:  5,
			expectedLimit: 100,
		},
		{
			name:          "Both zero should default to 1 and 20",
			page:          0,
			limit:         0,
			expectedPage:  1,
			expectedLimit: 20,
		},
		{
			name:          "Both negative should default to 1 and 20",
			page:          -1,
			limit:         -1,
			expectedPage:  1,
			expectedLimit: 20,
		},
		{
			name:          "Limit exactly 100 should remain 100",
			page:          5,
			limit:         100,
			expectedPage:  5,
			expectedLimit: 100,
		},
		{
			name:          "Limit 101 should be capped at 100",
			page:          5,
			limit:         101,
			expectedPage:  5,
			expectedLimit: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := util.NormalizePagination(tt.page, tt.limit)

			if result.Page != tt.expectedPage {
				t.Errorf("NormalizePagination(%d, %d) Page = %d, expected %d",
					tt.page, tt.limit, result.Page, tt.expectedPage)
			}

			if result.Limit != tt.expectedLimit {
				t.Errorf("NormalizePagination(%d, %d) Limit = %d, expected %d",
					tt.page, tt.limit, result.Limit, tt.expectedLimit)
			}
		})
	}
}

func TestCalculateTotalPages(t *testing.T) {
	tests := []struct {
		name     string
		total    int64
		limit    int
		expected int
	}{
		{
			name:     "Exact division",
			total:    100,
			limit:    10,
			expected: 10,
		},
		{
			name:     "Remainder requires extra page",
			total:    101,
			limit:    10,
			expected: 11,
		},
		{
			name:     "Single item",
			total:    1,
			limit:    10,
			expected: 1,
		},
		{
			name:     "Zero total",
			total:    0,
			limit:    10,
			expected: 0,
		},
		{
			name:     "Total less than limit",
			total:    5,
			limit:    10,
			expected: 1,
		},
		{
			name:     "Large numbers",
			total:    1000000,
			limit:    100,
			expected: 10000,
		},
		{
			name:     "Limit of 1",
			total:    5,
			limit:    1,
			expected: 5,
		},
		{
			name:     "Large total with small limit",
			total:    1000,
			limit:    3,
			expected: 334, // 1000/3 = 333.33... -> 334 pages
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := util.CalculateTotalPages(tt.total, tt.limit)
			if result != tt.expected {
				t.Errorf("CalculateTotalPages(%d, %d) = %d, expected %d",
					tt.total, tt.limit, result, tt.expected)
			}
		})
	}
}

func TestPaginationConfig(t *testing.T) {
	// Test that PaginationConfig struct works correctly
	config := util.PaginationConfig{
		Page:  5,
		Limit: 25,
	}

	if config.Page != 5 {
		t.Errorf("PaginationConfig.Page = %d, expected 5", config.Page)
	}

	if config.Limit != 25 {
		t.Errorf("PaginationConfig.Limit = %d, expected 25", config.Limit)
	}
}
