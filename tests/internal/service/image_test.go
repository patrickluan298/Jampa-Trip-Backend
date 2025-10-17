package service

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"mime/multipart"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jampa_trip/internal/contract"
	"github.com/jampa_trip/internal/model"
	"github.com/jampa_trip/internal/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDBForImageService(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
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

func createTestImage() *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			img.Set(x, y, color.RGBA{255, 0, 0, 255})
		}
	}
	return img
}

func createTestImageFile(t *testing.T) (*multipart.FileHeader, []byte) {
	img := createTestImage()

	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		t.Fatalf("Failed to encode test image: %v", err)
	}

	tmpFile, err := os.CreateTemp("", "test_image_*.png")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write(buf.Bytes())
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	file, err := os.Open(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to open temp file: %v", err)
	}
	defer file.Close()

	fileHeader := &multipart.FileHeader{
		Filename: "test_image.png",
		Size:     int64(buf.Len()),
		Header:   make(map[string][]string),
	}
	fileHeader.Header.Set("Content-Type", "image/png")

	return fileHeader, buf.Bytes()
}

func TestImageServiceNew(t *testing.T) {
	db, mock := setupMockDBForImageService(t)
	defer mock.ExpectationsWereMet()

	service := service.ImageServiceNew(db)
	if service == nil {
		t.Errorf("ImageServiceNew() returned nil")
	}
	if service.ImageRepository == nil {
		t.Errorf("ImageServiceNew() returned service with nil ImageRepository")
	}
}

func TestImageService_ListImages(t *testing.T) {
	db, mock := setupMockDBForImageService(t)
	defer mock.ExpectationsWereMet()

	service := service.ImageServiceNew(db)

	tests := []struct {
		name           string
		request        *contract.ListImagesRequest
		userID         int
		mockSetup      func()
		expectedError  bool
		expectedStatus int
	}{
		{
			name: "successful list with default pagination",
			request: &contract.ListImagesRequest{
				Page:  1,
				Limit: 20,
			},
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

				mock.ExpectQuery(`SELECT COUNT`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			},
			expectedError:  false,
			expectedStatus: http.StatusOK,
		},
		{
			name: "list with tour filter",
			request: &contract.ListImagesRequest{
				TourID: func() *int { id := 1; return &id }(),
				Page:   1,
				Limit:  20,
			},
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
				mock.ExpectQuery(`SELECT COUNT`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			},
			expectedError:  false,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			response, err := service.ListImages(tt.request, tt.userID)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if response == nil {
					t.Errorf("Expected response but got nil")
				}
				if !response.Success {
					t.Errorf("Expected success=true but got false")
				}
			}
		})
	}
}

func TestImageService_DeleteImage(t *testing.T) {
	db, mock := setupMockDBForImageService(t)
	defer mock.ExpectationsWereMet()

	service := service.ImageServiceNew(db)

	tests := []struct {
		name           string
		imageID        int
		userID         int
		mockSetup      func()
		expectedError  bool
		expectedStatus int
	}{
		{
			name:    "successful deletion",
			imageID: 1,
			userID:  1,
			mockSetup: func() {
				mock.ExpectQuery(`SELECT COUNT`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
				mock.ExpectQuery(`SELECT COUNT`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
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
				mock.ExpectExec(`DELETE FROM images`).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError:  false,
			expectedStatus: http.StatusOK,
		},
		{
			name:    "image not found",
			imageID: 999,
			userID:  1,
			mockSetup: func() {
				mock.ExpectQuery(`SELECT COUNT`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
			},
			expectedError:  true,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:    "image in use",
			imageID: 1,
			userID:  1,
			mockSetup: func() {
				mock.ExpectQuery(`SELECT COUNT`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
				mock.ExpectQuery(`SELECT COUNT`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			},
			expectedError:  true,
			expectedStatus: http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			response, err := service.DeleteImage(tt.imageID, tt.userID)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if response == nil {
					t.Errorf("Expected response but got nil")
				}
				if !response.Success {
					t.Errorf("Expected success=true but got false")
				}
			}
		})
	}
}

func TestImageService_UpdateImage(t *testing.T) {
	db, mock := setupMockDBForImageService(t)
	defer mock.ExpectationsWereMet()

	service := service.ImageServiceNew(db)

	tests := []struct {
		name           string
		imageID        int
		request        *contract.UpdateImageRequest
		userID         int
		mockSetup      func()
		expectedError  bool
		expectedStatus int
	}{
		{
			name:    "successful update",
			imageID: 1,
			request: &contract.UpdateImageRequest{
				Description: "Updated description",
				AltText:     "Updated alt text",
				IsPrimary:   func() *bool { b := true; return &b }(),
			},
			userID: 1,
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{
					"id", "user_id", "tour_id", "filename", "original_name", "url", "thumbnail_url",
					"size", "width", "height", "format", "description", "alt_text", "is_primary", "sort_order",
					"uploaded_at", "updated_at",
				}).AddRow(
					1, 1, 1, "test.jpg", "original.jpg", "http://localhost/test.jpg", "http://localhost/thumb_test.jpg",
					1024, 100, 100, "jpg", "Test image", "Alt text", false, 0,
					time.Now(), time.Now(),
				)
				mock.ExpectQuery(`SELECT`).WillReturnRows(rows)
				mock.ExpectExec(`UPDATE images`).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(`UPDATE images`).WillReturnResult(sqlmock.NewResult(1, 1))
				rows2 := sqlmock.NewRows([]string{
					"id", "user_id", "tour_id", "filename", "original_name", "url", "thumbnail_url",
					"size", "width", "height", "format", "description", "alt_text", "is_primary", "sort_order",
					"uploaded_at", "updated_at",
				}).AddRow(
					1, 1, 1, "test.jpg", "original.jpg", "http://localhost/test.jpg", "http://localhost/thumb_test.jpg",
					1024, 100, 100, "jpg", "Updated description", "Updated alt text", true, 0,
					time.Now(), time.Now(),
				)
				mock.ExpectQuery(`SELECT`).WillReturnRows(rows2)
			},
			expectedError:  false,
			expectedStatus: http.StatusOK,
		},
		{
			name:    "image not found",
			imageID: 999,
			request: &contract.UpdateImageRequest{
				Description: "Updated description",
			},
			userID: 1,
			mockSetup: func() {
				mock.ExpectQuery(`SELECT`).WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedError:  true,
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			response, err := service.UpdateImage(tt.imageID, tt.request, tt.userID)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if response == nil {
					t.Errorf("Expected response but got nil")
				}
				if !response.Success {
					t.Errorf("Expected success=true but got false")
				}
			}
		})
	}
}

func TestImageService_ReorderImages(t *testing.T) {
	db, mock := setupMockDBForImageService(t)
	defer mock.ExpectationsWereMet()

	service := service.ImageServiceNew(db)

	tests := []struct {
		name           string
		request        *contract.ReorderImagesRequest
		userID         int
		mockSetup      func()
		expectedError  bool
		expectedStatus int
	}{
		{
			name: "successful reorder",
			request: &contract.ReorderImagesRequest{
				ImageIDs: []int{1, 2, 3},
				TourID:   1,
			},
			userID: 1,
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{
					"id", "user_id", "tour_id", "filename", "original_name", "url", "thumbnail_url",
					"size", "width", "height", "format", "description", "alt_text", "is_primary", "sort_order",
					"uploaded_at", "updated_at",
				}).AddRow(
					1, 1, 1, "test1.jpg", "original1.jpg", "http://localhost/test1.jpg", "http://localhost/thumb_test1.jpg",
					1024, 100, 100, "jpg", "Test image 1", "Alt text 1", false, 0,
					time.Now(), time.Now(),
				).AddRow(
					2, 1, 1, "test2.jpg", "original2.jpg", "http://localhost/test2.jpg", "http://localhost/thumb_test2.jpg",
					1024, 100, 100, "jpg", "Test image 2", "Alt text 2", false, 1,
					time.Now(), time.Now(),
				).AddRow(
					3, 1, 1, "test3.jpg", "original3.jpg", "http://localhost/test3.jpg", "http://localhost/thumb_test3.jpg",
					1024, 100, 100, "jpg", "Test image 3", "Alt text 3", false, 2,
					time.Now(), time.Now(),
				)
				mock.ExpectQuery(`SELECT`).WillReturnRows(rows)
				mock.ExpectExec(`UPDATE images`).WillReturnResult(sqlmock.NewResult(3, 3))
			},
			expectedError:  false,
			expectedStatus: http.StatusOK,
		},
		{
			name: "some images not found",
			request: &contract.ReorderImagesRequest{
				ImageIDs: []int{1, 999, 3},
				TourID:   1,
			},
			userID: 1,
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{
					"id", "user_id", "tour_id", "filename", "original_name", "url", "thumbnail_url",
					"size", "width", "height", "format", "description", "alt_text", "is_primary", "sort_order",
					"uploaded_at", "updated_at",
				}).AddRow(
					1, 1, 1, "test1.jpg", "original1.jpg", "http://localhost/test1.jpg", "http://localhost/thumb_test1.jpg",
					1024, 100, 100, "jpg", "Test image 1", "Alt text 1", false, 0,
					time.Now(), time.Now(),
				).AddRow(
					3, 1, 1, "test3.jpg", "original3.jpg", "http://localhost/test3.jpg", "http://localhost/thumb_test3.jpg",
					1024, 100, 100, "jpg", "Test image 3", "Alt text 3", false, 2,
					time.Now(), time.Now(),
				)
				mock.ExpectQuery(`SELECT`).WillReturnRows(rows)
			},
			expectedError:  true,
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			response, err := service.ReorderImages(tt.request, tt.userID)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if response == nil {
					t.Errorf("Expected response but got nil")
				}
				if !response.Success {
					t.Errorf("Expected success=true but got false")
				}
			}
		})
	}
}

func TestImageService_GetImageInfo(t *testing.T) {
	db, mock := setupMockDBForImageService(t)
	defer mock.ExpectationsWereMet()

	service := service.ImageServiceNew(db)

	tests := []struct {
		name           string
		imageID        int
		userID         int
		mockSetup      func()
		expectedError  bool
		expectedStatus int
	}{
		{
			name:    "successful get info",
			imageID: 1,
			userID:  1,
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{
					"id", "user_id", "tour_id", "filename", "original_name", "url", "thumbnail_url",
					"size", "width", "height", "format", "description", "alt_text", "is_primary", "sort_order",
					"uploaded_at", "updated_at", "tour_name",
				}).AddRow(
					1, 1, 1, "test.jpg", "original.jpg", "http://localhost/test.jpg", "http://localhost/thumb_test.jpg",
					1024, 100, 100, "jpg", "Test image", "Alt text", true, 0,
					time.Now(), time.Now(), "Test Tour",
				)
				mock.ExpectQuery(`SELECT`).WillReturnRows(rows)
				usageRows := sqlmock.NewRows([]string{"tour_name", "is_used"}).AddRow("Test Tour", true)
				mock.ExpectQuery(`SELECT`).WillReturnRows(usageRows)
			},
			expectedError:  false,
			expectedStatus: http.StatusOK,
		},
		{
			name:    "image not found",
			imageID: 999,
			userID:  1,
			mockSetup: func() {
				mock.ExpectQuery(`SELECT`).WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedError:  true,
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			response, err := service.GetImageInfo(tt.imageID, tt.userID)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if response == nil {
					t.Errorf("Expected response but got nil")
				}
				if !response.Success {
					t.Errorf("Expected success=true but got false")
				}
			}
		})
	}
}

func TestImageService_BatchDeleteImages(t *testing.T) {
	db, mock := setupMockDBForImageService(t)
	defer mock.ExpectationsWereMet()

	service := service.ImageServiceNew(db)

	tests := []struct {
		name           string
		request        *contract.BatchDeleteImagesRequest
		userID         int
		mockSetup      func()
		expectedError  bool
		expectedStatus int
	}{
		{
			name: "successful batch delete",
			request: &contract.BatchDeleteImagesRequest{
				ImageIDs: []int{1, 2, 3},
			},
			userID: 1,
			mockSetup: func() {
				mock.ExpectQuery(`SELECT COUNT`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
				mock.ExpectQuery(`SELECT COUNT`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
				mock.ExpectQuery(`SELECT COUNT`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
				mock.ExpectQuery(`SELECT COUNT`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
				mock.ExpectQuery(`SELECT COUNT`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
				mock.ExpectQuery(`SELECT COUNT`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
				rows := sqlmock.NewRows([]string{
					"id", "user_id", "tour_id", "filename", "original_name", "url", "thumbnail_url",
					"size", "width", "height", "format", "description", "alt_text", "is_primary", "sort_order",
					"uploaded_at", "updated_at",
				}).AddRow(
					1, 1, 1, "test1.jpg", "original1.jpg", "http://localhost/test1.jpg", "http://localhost/thumb_test1.jpg",
					1024, 100, 100, "jpg", "Test image 1", "Alt text 1", false, 0,
					time.Now(), time.Now(),
				).AddRow(
					2, 1, 1, "test2.jpg", "original2.jpg", "http://localhost/test2.jpg", "http://localhost/thumb_test2.jpg",
					1024, 100, 100, "jpg", "Test image 2", "Alt text 2", false, 1,
					time.Now(), time.Now(),
				).AddRow(
					3, 1, 1, "test3.jpg", "original3.jpg", "http://localhost/test3.jpg", "http://localhost/thumb_test3.jpg",
					1024, 100, 100, "jpg", "Test image 3", "Alt text 3", false, 2,
					time.Now(), time.Now(),
				)
				mock.ExpectQuery(`SELECT`).WillReturnRows(rows)
				mock.ExpectExec(`DELETE FROM images`).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(`DELETE FROM images`).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(`DELETE FROM images`).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError:  false,
			expectedStatus: http.StatusOK,
		},
		{
			name: "partial success with errors",
			request: &contract.BatchDeleteImagesRequest{
				ImageIDs: []int{1, 999, 3},
			},
			userID: 1,
			mockSetup: func() {
				mock.ExpectQuery(`SELECT COUNT`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
				mock.ExpectQuery(`SELECT COUNT`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
				mock.ExpectQuery(`SELECT COUNT`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
				mock.ExpectQuery(`SELECT COUNT`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
				mock.ExpectQuery(`SELECT COUNT`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
				rows := sqlmock.NewRows([]string{
					"id", "user_id", "tour_id", "filename", "original_name", "url", "thumbnail_url",
					"size", "width", "height", "format", "description", "alt_text", "is_primary", "sort_order",
					"uploaded_at", "updated_at",
				}).AddRow(
					1, 1, 1, "test1.jpg", "original1.jpg", "http://localhost/test1.jpg", "http://localhost/thumb_test1.jpg",
					1024, 100, 100, "jpg", "Test image 1", "Alt text 1", false, 0,
					time.Now(), time.Now(),
				).AddRow(
					3, 1, 1, "test3.jpg", "original3.jpg", "http://localhost/test3.jpg", "http://localhost/thumb_test3.jpg",
					1024, 100, 100, "jpg", "Test image 3", "Alt text 3", false, 2,
					time.Now(), time.Now(),
				)
				mock.ExpectQuery(`SELECT`).WillReturnRows(rows)
				mock.ExpectExec(`DELETE FROM images`).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(`DELETE FROM images`).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError:  false,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			response, err := service.BatchDeleteImages(tt.request, tt.userID)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if response == nil {
					t.Errorf("Expected response but got nil")
				}
				if !response.Success {
					t.Errorf("Expected success=true but got false")
				}
			}
		})
	}
}

// Test helper functions
func TestImageModel_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		image    *model.Image
		expected bool
	}{
		{
			name: "valid image",
			image: &model.Image{
				UserID:   1,
				Filename: "test.jpg",
				URL:      "http://localhost/test.jpg",
				Size:     1024,
				Format:   "jpg",
			},
			expected: true,
		},
		{
			name: "invalid - no user ID",
			image: &model.Image{
				UserID:   0,
				Filename: "test.jpg",
				URL:      "http://localhost/test.jpg",
				Size:     1024,
				Format:   "jpg",
			},
			expected: false,
		},
		{
			name: "invalid - no filename",
			image: &model.Image{
				UserID:   1,
				Filename: "",
				URL:      "http://localhost/test.jpg",
				Size:     1024,
				Format:   "jpg",
			},
			expected: false,
		},
		{
			name: "invalid - no URL",
			image: &model.Image{
				UserID:   1,
				Filename: "test.jpg",
				URL:      "",
				Size:     1024,
				Format:   "jpg",
			},
			expected: false,
		},
		{
			name: "invalid - no size",
			image: &model.Image{
				UserID:   1,
				Filename: "test.jpg",
				URL:      "http://localhost/test.jpg",
				Size:     0,
				Format:   "jpg",
			},
			expected: false,
		},
		{
			name: "invalid - no format",
			image: &model.Image{
				UserID:   1,
				Filename: "test.jpg",
				URL:      "http://localhost/test.jpg",
				Size:     1024,
				Format:   "",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.image.IsValid()
			if result != tt.expected {
				t.Errorf("IsValid() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestImageModel_HasValidFormat(t *testing.T) {
	tests := []struct {
		name     string
		image    *model.Image
		expected bool
	}{
		{
			name: "valid jpg",
			image: &model.Image{
				Format: "jpg",
			},
			expected: true,
		},
		{
			name: "valid png",
			image: &model.Image{
				Format: "png",
			},
			expected: true,
		},
		{
			name: "valid gif",
			image: &model.Image{
				Format: "gif",
			},
			expected: true,
		},
		{
			name: "valid webp",
			image: &model.Image{
				Format: "webp",
			},
			expected: true,
		},
		{
			name: "invalid format",
			image: &model.Image{
				Format: "bmp",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.image.HasValidFormat()
			if result != tt.expected {
				t.Errorf("HasValidFormat() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestImageModel_HasValidSize(t *testing.T) {
	tests := []struct {
		name     string
		image    *model.Image
		expected bool
	}{
		{
			name: "valid size",
			image: &model.Image{
				Size: 1024,
			},
			expected: true,
		},
		{
			name: "valid max size",
			image: &model.Image{
				Size: 10 * 1024 * 1024, // 10MB
			},
			expected: true,
		},
		{
			name: "invalid - too large",
			image: &model.Image{
				Size: 11 * 1024 * 1024, // 11MB
			},
			expected: false,
		},
		{
			name: "invalid - zero size",
			image: &model.Image{
				Size: 0,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.image.HasValidSize()
			if result != tt.expected {
				t.Errorf("HasValidSize() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
