package util

// PaginationConfig - pagination configuration
type PaginationConfig struct {
	Page  int
	Limit int
}

// DefaultPagination - returns default pagination values
func DefaultPagination() PaginationConfig {
	return PaginationConfig{
		Page:  1,
		Limit: 20,
	}
}

// NormalizePagination - normalizes pagination parameters
func NormalizePagination(page, limit int) PaginationConfig {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	return PaginationConfig{Page: page, Limit: limit}
}

// CalculateTotalPages - calculates total pages
func CalculateTotalPages(total int64, limit int) int {
	return int((total + int64(limit) - 1) / int64(limit))
}
