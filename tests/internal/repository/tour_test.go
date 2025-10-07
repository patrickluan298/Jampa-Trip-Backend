package repository

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jampa_trip/internal/model"
	"github.com/jampa_trip/internal/repository"
	"github.com/lib/pq"
)

// setupMockDB is already defined in company_test.go

func TestTourRepositoryNew(t *testing.T) {
	db, mock := setupMockDB(t)
	defer mock.ExpectationsWereMet()

	repo := repository.TourRepositoryNew(db)
	if repo == nil {
		t.Error("TourRepositoryNew() returned nil")
	}
}

func TestTourRepository_Create(t *testing.T) {
	db, mock := setupMockDB(t)
	defer mock.ExpectationsWereMet()

	repo := repository.TourRepositoryNew(db)

	tests := []struct {
		name     string
		tour     *model.Tour
		hasError bool
	}{
		{
			name: "Valid tour",
			tour: &model.Tour{
				CompanyID:     1,
				Name:          "Passeio Histórico",
				Dates:         pq.StringArray{"2024-01-15"},
				DepartureTime: "08:00",
				ArrivalTime:   "18:00",
				MaxPeople:     20,
				Description:   "Passeio pela cidade histórica",
				Images:        pq.StringArray{"image1.jpg"},
				Price:         150.50,
			},
			hasError: false,
		},
		{
			name: "Tour with empty dates",
			tour: &model.Tour{
				CompanyID:     1,
				Name:          "Passeio Cultural",
				Dates:         pq.StringArray{},
				DepartureTime: "09:00",
				ArrivalTime:   "17:00",
				MaxPeople:     15,
				Description:   "Passeio cultural",
				Images:        pq.StringArray{},
				Price:         120.00,
			},
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
				WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

			err := repo.Create(tt.tour)
			if (err != nil) != tt.hasError {
				t.Errorf("Create() error = %v, hasError = %v", err, tt.hasError)
			}
		})
	}
}

func TestTourRepository_GetByID(t *testing.T) {
	// Skipping due to sqlmock limitations with PostgreSQL arrays
	t.Skip("Skipping due to sqlmock limitations with PostgreSQL arrays - requires database integration")
}

func TestTourRepository_List(t *testing.T) {
	// Skipping due to sqlmock limitations with PostgreSQL arrays
	t.Skip("Skipping due to sqlmock limitations with PostgreSQL arrays - requires database integration")
}
