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

func setupMockDBForImageRepository(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
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

func TestImageRepositoryNew(t *testing.T) {
	db, mock := setupMockDBForImageRepository(t)
	defer mock.ExpectationsWereMet()

	repo := repository.ImageRepositoryNew(db)
	if repo == nil {
		t.Errorf("ImageRepositoryNew() returned nil")
	}
	if repo.DB == nil {
		t.Errorf("ImageRepositoryNew() returned repository with nil DB")
	}
}

func TestImageRepository_Create(t *testing.T) {
	db, mock := setupMockDBForImageRepository(t)
	defer mock.ExpectationsWereMet()

	repo := repository.ImageRepositoryNew(db)

	tests := []struct {
		name          string
		image         *model.Image
		mockSetup     func()
		expectedError bool
	}{
		{
			name: "successful creation",
			image: &model.Image{
				UserID:       1,
				TourID:       func() *int { id := 1; return &id }(),
				Filename:     "test.jpg",
				OriginalName: "original.jpg",
				URL:          "http://localhost/test.jpg",
				ThumbnailURL: "http://localhost/thumb_test.jpg",
				Size:         1024,
				Width:        100,
				Height:       100,
				Format:       "jpg",
				Description:  "Test image",
				AltText:      "Alt text",
				IsPrimary:    false,
				SortOrder:    0,
			},
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "uploaded_at", "updated_at"}).
					AddRow(1, time.Now(), time.Now())
				mock.ExpectQuery(`INSERT INTO images`).WillReturnRows(rows)
			},
			expectedError: false,
		},
		{
			name: "creation with error",
			image: &model.Image{
				UserID:   1,
				Filename: "test.jpg",
				URL:      "http://localhost/test.jpg",
				Size:     1024,
				Format:   "jpg",
			},
			mockSetup: func() {
				mock.ExpectQuery(`INSERT INTO images`).WillReturnError(gorm.ErrInvalidData)
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := repo.Create(tt.image)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestImageRepository_GetByID(t *testing.T) {
	db, mock := setupMockDBForImageRepository(t)
	defer mock.ExpectationsWereMet()

	repo := repository.ImageRepositoryNew(db)

	tests := []struct {
		name          string
		id            int
		mockSetup     func()
		expectedError bool
	}{
		{
			name: "successful get by ID",
			id:   1,
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{
					"id", "user_id", "tour_id", "filename", "original_name", "url", "thumbnail_url",
					"size", "width", "height", "format", "description", "alt_text", "is_primary", "sort_order",
					"uploaded_at", "updated_at",
				}).AddRow(
					1, 1, 1, "test.jpg", "original.jpg", "http://localhost/test.jpg", "http://localhost/thumb_test.jpg",
					1024, 100, 100, "jpg", "Test image", "Alt text", true, 0,
					time.Now(), time.Now(),
				)
				mock.ExpectQuery(`SELECT`).WillReturnRows(rows)
			},
			expectedError: false,
		},
		{
			name: "image not found",
			id:   999,
			mockSetup: func() {
				mock.ExpectQuery(`SELECT`).WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			image, err := repo.GetByID(tt.id)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if image == nil {
					t.Errorf("Expected image but got nil")
				}
			}
		})
	}
}

func TestImageRepository_GetByIDAndUser(t *testing.T) {
	db, mock := setupMockDBForImageRepository(t)
	defer mock.ExpectationsWereMet()

	repo := repository.ImageRepositoryNew(db)

	tests := []struct {
		name          string
		id            int
		userID        int
		mockSetup     func()
		expectedError bool
	}{
		{
			name:   "successful get by ID and user",
			id:     1,
			userID: 1,
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{
					"id", "user_id", "tour_id", "filename", "original_name", "url", "thumbnail_url",
					"size", "width", "height", "format", "description", "alt_text", "is_primary", "sort_order",
					"uploaded_at", "updated_at",
				}).AddRow(
					1, 1, 1, "test.jpg", "original.jpg", "http://localhost/test.jpg", "http://localhost/thumb_test.jpg",
					1024, 100, 100, "jpg", "Test image", "Alt text", true, 0,
					time.Now(), time.Now(),
				)
				mock.ExpectQuery(`SELECT`).WillReturnRows(rows)
			},
			expectedError: false,
		},
		{
			name:   "image not found for user",
			id:     1,
			userID: 2,
			mockSetup: func() {
				mock.ExpectQuery(`SELECT`).WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			image, err := repo.GetByIDAndUser(tt.id, tt.userID)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if image == nil {
					t.Errorf("Expected image but got nil")
				}
			}
		})
	}
}

func TestImageRepository_Update(t *testing.T) {
	db, mock := setupMockDBForImageRepository(t)
	defer mock.ExpectationsWereMet()

	repo := repository.ImageRepositoryNew(db)

	tests := []struct {
		name          string
		image         *model.Image
		mockSetup     func()
		expectedError bool
	}{
		{
			name: "successful update",
			image: &model.Image{
				ID:          1,
				TourID:      func() *int { id := 1; return &id }(),
				Description: "Updated description",
				AltText:     "Updated alt text",
				IsPrimary:   true,
			},
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"updated_at"}).AddRow(time.Now())
				mock.ExpectQuery(`UPDATE images`).WillReturnRows(rows)
			},
			expectedError: false,
		},
		{
			name: "update with error",
			image: &model.Image{
				ID: 1,
			},
			mockSetup: func() {
				mock.ExpectQuery(`UPDATE images`).WillReturnError(gorm.ErrInvalidData)
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := repo.Update(tt.image)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestImageRepository_Delete(t *testing.T) {
	db, mock := setupMockDBForImageRepository(t)
	defer mock.ExpectationsWereMet()

	repo := repository.ImageRepositoryNew(db)

	tests := []struct {
		name          string
		id            int
		mockSetup     func()
		expectedError bool
	}{
		{
			name: "successful delete",
			id:   1,
			mockSetup: func() {
				mock.ExpectExec(`DELETE FROM images`).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: false,
		},
		{
			name: "delete with error",
			id:   999,
			mockSetup: func() {
				mock.ExpectExec(`DELETE FROM images`).WillReturnError(gorm.ErrInvalidData)
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := repo.Delete(tt.id)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestImageRepository_List(t *testing.T) {
	db, mock := setupMockDBForImageRepository(t)
	defer mock.ExpectationsWereMet()

	repo := repository.ImageRepositoryNew(db)

	tests := []struct {
		name          string
		userID        int
		tourID        *int
		format        string
		sortBy        string
		page          int
		limit         int
		mockSetup     func()
		expectedError bool
	}{
		{
			name:   "successful list",
			userID: 1,
			page:   1,
			limit:  20,
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{
					"id", "user_id", "tour_id", "filename", "original_name", "url", "thumbnail_url",
					"size", "width", "height", "format", "description", "alt_text", "is_primary", "sort_order",
					"uploaded_at", "updated_at",
				}).AddRow(
					1, 1, 1, "test.jpg", "original.jpg", "http://localhost/test.jpg", "http://localhost/thumb_test.jpg",
					1024, 100, 100, "jpg", "Test image", "Alt text", true, 0,
					time.Now(), time.Now(),
				)
				mock.ExpectQuery(`SELECT`).WillReturnRows(rows)
				mock.ExpectQuery(`SELECT COUNT`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			},
			expectedError: false,
		},
		{
			name:   "list with error",
			userID: 1,
			page:   1,
			limit:  20,
			mockSetup: func() {
				mock.ExpectQuery(`SELECT`).WillReturnError(gorm.ErrInvalidData)
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			images, total, err := repo.List(tt.userID, tt.tourID, tt.format, tt.sortBy, tt.page, tt.limit)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if images == nil {
					t.Errorf("Expected images but got nil")
				}
				if total < 0 {
					t.Errorf("Expected total >= 0 but got %d", total)
				}
			}
		})
	}
}

func TestImageRepository_IsOwnedByUser(t *testing.T) {
	db, mock := setupMockDBForImageRepository(t)
	defer mock.ExpectationsWereMet()

	repo := repository.ImageRepositoryNew(db)

	tests := []struct {
		name          string
		imageID       int
		userID        int
		mockSetup     func()
		expected      bool
		expectedError bool
	}{
		{
			name:    "image is owned by user",
			imageID: 1,
			userID:  1,
			mockSetup: func() {
				mock.ExpectQuery(`SELECT COUNT`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			},
			expected:      true,
			expectedError: false,
		},
		{
			name:    "image is not owned by user",
			imageID: 1,
			userID:  2,
			mockSetup: func() {
				mock.ExpectQuery(`SELECT COUNT`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			expected:      false,
			expectedError: false,
		},
		{
			name:    "query error",
			imageID: 1,
			userID:  1,
			mockSetup: func() {
				mock.ExpectQuery(`SELECT COUNT`).WillReturnError(gorm.ErrInvalidData)
			},
			expected:      false,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			result, err := repo.IsOwnedByUser(tt.imageID, tt.userID)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("IsOwnedByUser() = %v, expected %v", result, tt.expected)
				}
			}
		})
	}
}

func TestImageRepository_IsUsedInActiveTour(t *testing.T) {
	db, mock := setupMockDBForImageRepository(t)
	defer mock.ExpectationsWereMet()

	repo := repository.ImageRepositoryNew(db)

	tests := []struct {
		name          string
		imageID       int
		mockSetup     func()
		expected      bool
		expectedError bool
	}{
		{
			name:    "image is used in active tour",
			imageID: 1,
			mockSetup: func() {
				mock.ExpectQuery(`SELECT COUNT`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			},
			expected:      true,
			expectedError: false,
		},
		{
			name:    "image is not used in active tour",
			imageID: 1,
			mockSetup: func() {
				mock.ExpectQuery(`SELECT COUNT`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			expected:      false,
			expectedError: false,
		},
		{
			name:    "query error",
			imageID: 1,
			mockSetup: func() {
				mock.ExpectQuery(`SELECT COUNT`).WillReturnError(gorm.ErrInvalidData)
			},
			expected:      false,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			result, err := repo.IsUsedInActiveTour(tt.imageID)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("IsUsedInActiveTour() = %v, expected %v", result, tt.expected)
				}
			}
		})
	}
}

func TestImageRepository_GetImageUsage(t *testing.T) {
	db, mock := setupMockDBForImageRepository(t)
	defer mock.ExpectationsWereMet()

	repo := repository.ImageRepositoryNew(db)

	tests := []struct {
		name          string
		imageID       int
		mockSetup     func()
		expectedName  string
		expectedUsed  bool
		expectedError bool
	}{
		{
			name:    "image is used in tour",
			imageID: 1,
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"tour_name", "is_used"}).AddRow("Test Tour", true)
				mock.ExpectQuery(`SELECT`).WillReturnRows(rows)
			},
			expectedName:  "Test Tour",
			expectedUsed:  true,
			expectedError: false,
		},
		{
			name:    "image is not used",
			imageID: 1,
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"tour_name", "is_used"}).AddRow(nil, false)
				mock.ExpectQuery(`SELECT`).WillReturnRows(rows)
			},
			expectedName:  "",
			expectedUsed:  false,
			expectedError: false,
		},
		{
			name:    "query error",
			imageID: 1,
			mockSetup: func() {
				mock.ExpectQuery(`SELECT`).WillReturnError(gorm.ErrInvalidData)
			},
			expectedName:  "",
			expectedUsed:  false,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			tourName, isUsed, err := repo.GetImageUsage(tt.imageID)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if tourName != tt.expectedName {
					t.Errorf("GetImageUsage() tourName = %v, expected %v", tourName, tt.expectedName)
				}
				if isUsed != tt.expectedUsed {
					t.Errorf("GetImageUsage() isUsed = %v, expected %v", isUsed, tt.expectedUsed)
				}
			}
		})
	}
}

func TestImageRepository_UpdateSortOrder(t *testing.T) {
	db, mock := setupMockDBForImageRepository(t)
	defer mock.ExpectationsWereMet()

	repo := repository.ImageRepositoryNew(db)

	tests := []struct {
		name          string
		imageID       int
		sortOrder     int
		userID        int
		mockSetup     func()
		expectedError bool
	}{
		{
			name:      "successful update sort order",
			imageID:   1,
			sortOrder: 2,
			userID:    1,
			mockSetup: func() {
				mock.ExpectExec(`UPDATE images`).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: false,
		},
		{
			name:      "update sort order with error",
			imageID:   1,
			sortOrder: 2,
			userID:    1,
			mockSetup: func() {
				mock.ExpectExec(`UPDATE images`).WillReturnError(gorm.ErrInvalidData)
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := repo.UpdateSortOrder(tt.imageID, tt.sortOrder, tt.userID)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestImageRepository_Exists(t *testing.T) {
	db, mock := setupMockDBForImageRepository(t)
	defer mock.ExpectationsWereMet()

	repo := repository.ImageRepositoryNew(db)

	tests := []struct {
		name          string
		id            int
		mockSetup     func()
		expected      bool
		expectedError bool
	}{
		{
			name: "image exists",
			id:   1,
			mockSetup: func() {
				mock.ExpectQuery(`SELECT COUNT`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			},
			expected:      true,
			expectedError: false,
		},
		{
			name: "image does not exist",
			id:   999,
			mockSetup: func() {
				mock.ExpectQuery(`SELECT COUNT`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			expected:      false,
			expectedError: false,
		},
		{
			name: "query error",
			id:   1,
			mockSetup: func() {
				mock.ExpectQuery(`SELECT COUNT`).WillReturnError(gorm.ErrInvalidData)
			},
			expected:      false,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			result, err := repo.Exists(tt.id)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("Exists() = %v, expected %v", result, tt.expected)
				}
			}
		})
	}
}
