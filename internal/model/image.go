package model

import (
	"fmt"
	"strings"
	"time"
)

// Image - representa a entidade de imagem
type Image struct {
	ID           int       `gorm:"column:id;primaryKey;autoIncrement"`
	UserID       int       `gorm:"column:user_id;not null"`
	TourID       *int      `gorm:"column:tour_id"`
	Filename     string    `gorm:"column:filename;not null"`
	OriginalName string    `gorm:"column:original_name"`
	URL          string    `gorm:"column:url;not null"`
	ThumbnailURL string    `gorm:"column:thumbnail_url"`
	Size         int       `gorm:"column:size;not null"`
	Width        int       `gorm:"column:width"`
	Height       int       `gorm:"column:height"`
	Format       string    `gorm:"column:format;not null"`
	Description  string    `gorm:"column:description"`
	AltText      string    `gorm:"column:alt_text"`
	IsPrimary    bool      `gorm:"column:is_primary;default:false"`
	SortOrder    int       `gorm:"column:sort_order;default:0"`
	UploadedAt   time.Time `gorm:"column:uploaded_at;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP"`
}

// TableName - especifica o nome da tabela no banco de dados
func (Image) TableName() string {
	return "images"
}

// Métodos de validação para Image
func (i *Image) IsValid() bool {
	return i.UserID > 0 &&
		len(i.Filename) > 0 &&
		len(i.URL) > 0 &&
		i.Size > 0 &&
		len(i.Format) > 0
}

func (i *Image) HasValidFormat() bool {
	validFormats := []string{"jpg", "jpeg", "png", "gif", "webp"}
	format := strings.ToLower(i.Format)

	for _, validFormat := range validFormats {
		if format == validFormat {
			return true
		}
	}
	return false
}

func (i *Image) HasValidSize() bool {
	// Máximo 10MB
	const maxSize = 10 * 1024 * 1024
	return i.Size > 0 && i.Size <= maxSize
}

func (i *Image) HasValidDimensions() bool {
	return i.Width > 0 && i.Height > 0
}

func (i *Image) GetFormattedSize() string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
	case i.Size >= GB:
		return fmt.Sprintf("%.2f GB", float64(i.Size)/GB)
	case i.Size >= MB:
		return fmt.Sprintf("%.2f MB", float64(i.Size)/MB)
	case i.Size >= KB:
		return fmt.Sprintf("%.2f KB", float64(i.Size)/KB)
	default:
		return fmt.Sprintf("%d bytes", i.Size)
	}
}

func (i *Image) GetAspectRatio() float64 {
	if i.Height == 0 {
		return 0
	}
	return float64(i.Width) / float64(i.Height)
}

func (i *Image) IsLandscape() bool {
	return i.Width > i.Height
}

func (i *Image) IsPortrait() bool {
	return i.Height > i.Width
}

func (i *Image) IsSquare() bool {
	return i.Width == i.Height
}

func (i *Image) GetThumbnailDimensions() (int, int) {
	const maxThumbnailSize = 300

	if i.Width <= maxThumbnailSize && i.Height <= maxThumbnailSize {
		return i.Width, i.Height
	}

	return i.Width, i.Height
}
