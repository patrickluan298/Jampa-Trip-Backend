package repository

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jampa_trip/internal/model"
	"github.com/jampa_trip/internal/repository"
)

// setupMockDB is already defined in company_test.go

func TestTourRepositoryNew(t *testing.T) {
	db, mock := setupMockDB(t)
	defer mock.ExpectationsWereMet()

	repo := repository.TourRepositoryNew(db)
	if repo == nil {
		t.Errorf("TourRepositoryNew() returned nil")
	}
	if repo.DB == nil {
		t.Errorf("TourRepositoryNew() returned repository with nil DB")
	}
}

func TestTourRepository_Create(t *testing.T) {
	db, mock := setupMockDB(t)
	defer mock.ExpectationsWereMet()

	repo := repository.TourRepositoryNew(db)

	tests := []struct {
		name     string
		tour     *model.Tour
		mockID   int
		hasError bool
	}{
		{
			name: "Valid tour",
			tour: &model.Tour{
				CompanyID:     1,
				Name:          "Passeio Histórico",
				Dates:         []string{"2024-01-15", "2024-01-22"},
				DepartureTime: "08:00",
				ArrivalTime:   "18:00",
				MaxPeople:     20,
				Description:   "Passeio pela cidade histórica",
				Images:        []string{"image1.jpg", "image2.jpg"},
				Price:         150.50,
			},
			mockID:   1,
			hasError: false,
		},
		{
			name: "Tour with empty dates",
			tour: &model.Tour{
				CompanyID:     1,
				Name:          "Passeio Especial",
				Dates:         []string{},
				DepartureTime: "09:00",
				ArrivalTime:   "17:00",
				MaxPeople:     15,
				Description:   "Passeio especial",
				Images:        []string{},
				Price:         200.00,
			},
			mockID:   2,
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectQuery(`INSERT INTO tours`).
				WithArgs(
					tt.tour.CompanyID,
					tt.tour.Name,
					tt.tour.Dates,
					tt.tour.DepartureTime,
					tt.tour.ArrivalTime,
					tt.tour.MaxPeople,
					tt.tour.Description,
					tt.tour.Images,
					tt.tour.Price,
				).
				WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(tt.mockID))

			err := repo.Create(tt.tour)
			if (err != nil) != tt.hasError {
				t.Errorf("Create() error = %v, hasError = %v", err, tt.hasError)
			}

			if !tt.hasError {
				if tt.tour.ID != tt.mockID {
					t.Errorf("Create() ID = %d, expected %d", tt.tour.ID, tt.mockID)
				}
			}
		})
	}
}

func TestTourRepository_GetByID(t *testing.T) {
	db, mock := setupMockDB(t)
	defer mock.ExpectationsWereMet()

	repo := repository.TourRepositoryNew(db)

	tests := []struct {
		name     string
		id       int
		mockRows *sqlmock.Rows
		expected *model.Tour
		hasError bool
	}{
		{
			name: "Valid tour ID",
			id:   1,
			mockRows: sqlmock.NewRows([]string{"id", "company_id", "name", "dates", "departure_time", "arrival_time", "max_people", "description", "images", "price", "created_at", "updated_at", "company_name"}).
				AddRow(1, 1, "Passeio Histórico", []string{"2024-01-15", "2024-01-22"}, "08:00", "18:00", 20, "Passeio pela cidade histórica", []string{"image1.jpg", "image2.jpg"}, 150.50, time.Now(), time.Now(), "Empresa ABC"),
			expected: &model.Tour{
				ID:            1,
				CompanyID:     1,
				Name:          "Passeio Histórico",
				Dates:         []string{"2024-01-15", "2024-01-22"},
				DepartureTime: "08:00",
				ArrivalTime:   "18:00",
				MaxPeople:     20,
				Description:   "Passeio pela cidade histórica",
				Images:        []string{"image1.jpg", "image2.jpg"},
				Price:         150.50,
			},
			hasError: false,
		},
		{
			name:     "Non-existent tour ID",
			id:       999,
			mockRows: sqlmock.NewRows([]string{"id", "company_id", "name", "dates", "departure_time", "arrival_time", "max_people", "description", "images", "price", "created_at", "updated_at", "company_name"}),
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
			}
		})
	}
}

func TestTourRepository_List(t *testing.T) {
	db, mock := setupMockDB(t)
	defer mock.ExpectationsWereMet()

	repo := repository.TourRepositoryNew(db)

	tests := []struct {
		name          string
		search        string
		page          int
		limit         int
		mockRows      *sqlmock.Rows
		mockTotal     int64
		expectedLen   int
		expectedTotal int64
		hasError      bool
	}{
		{
			name:   "List tours with search",
			search: "histórico",
			page:   1,
			limit:  10,
			mockRows: sqlmock.NewRows([]string{"id", "company_id", "name", "dates", "departure_time", "arrival_time", "max_people", "description", "images", "price", "created_at", "updated_at", "company_name"}).
				AddRow(1, 1, "Passeio Histórico", []string{"2024-01-15"}, "08:00", "18:00", 20, "Passeio pela cidade histórica", []string{"image1.jpg"}, 150.50, time.Now(), time.Now(), "Empresa ABC"),
			mockTotal:     1,
			expectedLen:   1,
			expectedTotal: 1,
			hasError:      false,
		},
		{
			name:   "List tours without search",
			search: "",
			page:   1,
			limit:  10,
			mockRows: sqlmock.NewRows([]string{"id", "company_id", "name", "dates", "departure_time", "arrival_time", "max_people", "description", "images", "price", "created_at", "updated_at", "company_name"}).
				AddRow(1, 1, "Passeio Histórico", []string{"2024-01-15"}, "08:00", "18:00", 20, "Passeio pela cidade histórica", []string{"image1.jpg"}, 150.50, time.Now(), time.Now(), "Empresa ABC").
				AddRow(2, 1, "Passeio Cultural", []string{"2024-01-20"}, "09:00", "17:00", 15, "Passeio cultural", []string{"image2.jpg"}, 120.00, time.Now(), time.Now(), "Empresa ABC"),
			mockTotal:     2,
			expectedLen:   2,
			expectedTotal: 2,
			hasError:      false,
		},
		{
			name:          "Empty result",
			search:        "nonexistent",
			page:          1,
			limit:         10,
			mockRows:      sqlmock.NewRows([]string{"id", "company_id", "name", "dates", "departure_time", "arrival_time", "max_people", "description", "images", "price", "created_at", "updated_at", "company_name"}),
			mockTotal:     0,
			expectedLen:   0,
			expectedTotal: 0,
			hasError:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock the main query
			mock.ExpectQuery(`SELECT`).
				WithArgs(tt.search, "%"+tt.search+"%", tt.limit, (tt.page-1)*tt.limit).
				WillReturnRows(tt.mockRows)

			// Mock the count query
			mock.ExpectQuery(`SELECT count`).
				WithArgs(tt.search, "%"+tt.search+"%").
				WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(tt.mockTotal))

			result, total, err := repo.List(tt.search, tt.page, tt.limit)
			if (err != nil) != tt.hasError {
				t.Errorf("List() error = %v, hasError = %v", err, tt.hasError)
			}

			if !tt.hasError {
				if len(result) != tt.expectedLen {
					t.Errorf("List() returned %d tours, expected %d", len(result), tt.expectedLen)
				}
				if total != tt.expectedTotal {
					t.Errorf("List() total = %d, expected %d", total, tt.expectedTotal)
				}
			}
		})
	}
}

func TestTourRepository_ListByCompanyID(t *testing.T) {
	db, mock := setupMockDB(t)
	defer mock.ExpectationsWereMet()

	repo := repository.TourRepositoryNew(db)

	tests := []struct {
		name          string
		companyID     int
		page          int
		limit         int
		mockRows      *sqlmock.Rows
		mockTotal     int64
		expectedLen   int
		expectedTotal int64
		hasError      bool
	}{
		{
			name:      "List tours for company",
			companyID: 1,
			page:      1,
			limit:     10,
			mockRows: sqlmock.NewRows([]string{"id", "company_id", "name", "dates", "departure_time", "arrival_time", "max_people", "description", "images", "price", "created_at", "reservations_count"}).
				AddRow(1, 1, "Passeio Histórico", []string{"2024-01-15"}, "08:00", "18:00", 20, "Passeio pela cidade histórica", []string{"image1.jpg"}, 150.50, time.Now(), 5),
			mockTotal:     1,
			expectedLen:   1,
			expectedTotal: 1,
			hasError:      false,
		},
		{
			name:          "Empty result for company",
			companyID:     999,
			page:          1,
			limit:         10,
			mockRows:      sqlmock.NewRows([]string{"id", "company_id", "name", "dates", "departure_time", "arrival_time", "max_people", "description", "images", "price", "created_at", "reservations_count"}),
			mockTotal:     0,
			expectedLen:   0,
			expectedTotal: 0,
			hasError:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock the main query
			mock.ExpectQuery(`SELECT`).
				WithArgs(tt.companyID, tt.limit, (tt.page-1)*tt.limit).
				WillReturnRows(tt.mockRows)

			// Mock the count query
			mock.ExpectQuery(`SELECT count`).
				WithArgs(tt.companyID).
				WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(tt.mockTotal))

			result, total, err := repo.ListByCompanyID(tt.companyID, tt.page, tt.limit)
			if (err != nil) != tt.hasError {
				t.Errorf("ListByCompanyID() error = %v, hasError = %v", err, tt.hasError)
			}

			if !tt.hasError {
				if len(result) != tt.expectedLen {
					t.Errorf("ListByCompanyID() returned %d tours, expected %d", len(result), tt.expectedLen)
				}
				if total != tt.expectedTotal {
					t.Errorf("ListByCompanyID() total = %d, expected %d", total, tt.expectedTotal)
				}
			}
		})
	}
}

func TestTourRepository_Update(t *testing.T) {
	db, mock := setupMockDB(t)
	defer mock.ExpectationsWereMet()

	repo := repository.TourRepositoryNew(db)

	tests := []struct {
		name     string
		tour     *model.Tour
		hasError bool
	}{
		{
			name: "Valid tour update",
			tour: &model.Tour{
				ID:            1,
				Name:          "Passeio Histórico Atualizado",
				Dates:         []string{"2024-01-15", "2024-01-22", "2024-01-29"},
				DepartureTime: "09:00",
				ArrivalTime:   "19:00",
				MaxPeople:     25,
				Description:   "Passeio pela cidade histórica atualizado",
				Images:        []string{"image1.jpg", "image2.jpg", "image3.jpg"},
				Price:         180.00,
			},
			hasError: false,
		},
		{
			name: "Update non-existent tour",
			tour: &model.Tour{
				ID:            999,
				Name:          "Non-existent",
				Dates:         []string{"2024-01-15"},
				DepartureTime: "08:00",
				ArrivalTime:   "18:00",
				MaxPeople:     20,
				Description:   "Non-existent tour",
				Images:        []string{},
				Price:         100.00,
			},
			hasError: false, // Update operation doesn't fail if record doesn't exist
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectExec(`UPDATE tours`).
				WithArgs(
					tt.tour.Name,
					tt.tour.Dates,
					tt.tour.DepartureTime,
					tt.tour.ArrivalTime,
					tt.tour.MaxPeople,
					tt.tour.Description,
					tt.tour.Images,
					tt.tour.Price,
					tt.tour.ID,
				).
				WillReturnResult(sqlmock.NewResult(1, 1))

			err := repo.Update(tt.tour)
			if (err != nil) != tt.hasError {
				t.Errorf("Update() error = %v, hasError = %v", err, tt.hasError)
			}
		})
	}
}

func TestTourRepository_Delete(t *testing.T) {
	db, mock := setupMockDB(t)
	defer mock.ExpectationsWereMet()

	repo := repository.TourRepositoryNew(db)

	tests := []struct {
		name     string
		id       int
		hasError bool
	}{
		{
			name:     "Valid tour deletion",
			id:       1,
			hasError: false,
		},
		{
			name:     "Delete non-existent tour",
			id:       999,
			hasError: false, // Delete operation doesn't fail if record doesn't exist
		},
		{
			name:     "Delete zero ID",
			id:       0,
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectExec(`DELETE FROM tours`).
				WithArgs(tt.id).
				WillReturnResult(sqlmock.NewResult(1, 1))

			err := repo.Delete(tt.id)
			if (err != nil) != tt.hasError {
				t.Errorf("Delete() error = %v, hasError = %v", err, tt.hasError)
			}
		})
	}
}

func TestTourRepository_IsOwnedByCompany(t *testing.T) {
	db, mock := setupMockDB(t)
	defer mock.ExpectationsWereMet()

	repo := repository.TourRepositoryNew(db)

	tests := []struct {
		name      string
		tourID    int
		companyID int
		mockCount int
		expected  bool
		hasError  bool
	}{
		{
			name:      "Tour owned by company",
			tourID:    1,
			companyID: 1,
			mockCount: 1,
			expected:  true,
			hasError:  false,
		},
		{
			name:      "Tour not owned by company",
			tourID:    1,
			companyID: 2,
			mockCount: 0,
			expected:  false,
			hasError:  false,
		},
		{
			name:      "Non-existent tour",
			tourID:    999,
			companyID: 1,
			mockCount: 0,
			expected:  false,
			hasError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectQuery(`SELECT count`).
				WithArgs(tt.tourID, tt.companyID).
				WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(tt.mockCount))

			result, err := repo.IsOwnedByCompany(tt.tourID, tt.companyID)
			if (err != nil) != tt.hasError {
				t.Errorf("IsOwnedByCompany() error = %v, hasError = %v", err, tt.hasError)
			}

			if !tt.hasError {
				if result != tt.expected {
					t.Errorf("IsOwnedByCompany() = %v, expected %v", result, tt.expected)
				}
			}
		})
	}
}

func TestTourRepository_CountReservationsByTourID(t *testing.T) {
	db, mock := setupMockDB(t)
	defer mock.ExpectationsWereMet()

	repo := repository.TourRepositoryNew(db)

	tests := []struct {
		name      string
		tourID    int
		mockCount int
		expected  int
		hasError  bool
	}{
		{
			name:      "Tour with reservations",
			tourID:    1,
			mockCount: 5,
			expected:  5,
			hasError:  false,
		},
		{
			name:      "Tour without reservations",
			tourID:    2,
			mockCount: 0,
			expected:  0,
			hasError:  false,
		},
		{
			name:      "Non-existent tour",
			tourID:    999,
			mockCount: 0,
			expected:  0,
			hasError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectQuery(`SELECT count`).
				WithArgs(tt.tourID).
				WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(tt.mockCount))

			result, err := repo.CountReservationsByTourID(tt.tourID)
			if (err != nil) != tt.hasError {
				t.Errorf("CountReservationsByTourID() error = %v, hasError = %v", err, tt.hasError)
			}

			if !tt.hasError {
				if result != tt.expected {
					t.Errorf("CountReservationsByTourID() = %d, expected %d", result, tt.expected)
				}
			}
		})
	}
}

func TestTourRepository_GetTourWithCompanyName(t *testing.T) {
	db, mock := setupMockDB(t)
	defer mock.ExpectationsWereMet()

	repo := repository.TourRepositoryNew(db)

	tests := []struct {
		name         string
		id           int
		mockRows     *sqlmock.Rows
		expectedTour *model.Tour
		expectedName string
		hasError     bool
	}{
		{
			name: "Valid tour with company name",
			id:   1,
			mockRows: sqlmock.NewRows([]string{"id", "company_id", "name", "dates", "departure_time", "arrival_time", "max_people", "description", "images", "price", "created_at", "updated_at", "company_name"}).
				AddRow(1, 1, "Passeio Histórico", []string{"2024-01-15"}, "08:00", "18:00", 20, "Passeio pela cidade histórica", []string{"image1.jpg"}, 150.50, time.Now(), time.Now(), "Empresa ABC"),
			expectedTour: &model.Tour{
				ID:            1,
				CompanyID:     1,
				Name:          "Passeio Histórico",
				Dates:         []string{"2024-01-15"},
				DepartureTime: "08:00",
				ArrivalTime:   "18:00",
				MaxPeople:     20,
				Description:   "Passeio pela cidade histórica",
				Images:        []string{"image1.jpg"},
				Price:         150.50,
			},
			expectedName: "Empresa ABC",
			hasError:     false,
		},
		{
			name:         "Non-existent tour",
			id:           999,
			mockRows:     sqlmock.NewRows([]string{"id", "company_id", "name", "dates", "departure_time", "arrival_time", "max_people", "description", "images", "price", "created_at", "updated_at", "company_name"}),
			expectedTour: nil,
			expectedName: "",
			hasError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectQuery(`SELECT`).
				WithArgs(tt.id).
				WillReturnRows(tt.mockRows)

			tour, companyName, err := repo.GetTourWithCompanyName(tt.id)
			if (err != nil) != tt.hasError {
				t.Errorf("GetTourWithCompanyName() error = %v, hasError = %v", err, tt.hasError)
			}

			if !tt.hasError {
				if tour.ID != tt.expectedTour.ID {
					t.Errorf("GetTourWithCompanyName() tour ID = %d, expected %d", tour.ID, tt.expectedTour.ID)
				}
				if companyName != tt.expectedName {
					t.Errorf("GetTourWithCompanyName() company name = %s, expected %s", companyName, tt.expectedName)
				}
			}
		})
	}
}
