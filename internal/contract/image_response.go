package contract

import "time"

// ImageResponse - resposta com dados de uma imagem
type ImageResponse struct {
	ID           int       `json:"id"`
	Filename     string    `json:"filename"`
	OriginalName string    `json:"original_name,omitempty"`
	URL          string    `json:"url"`
	ThumbnailURL string    `json:"thumbnail_url,omitempty"`
	Size         int       `json:"size"`
	Width        int       `json:"width"`
	Height       int       `json:"height"`
	Format       string    `json:"format"`
	Description  string    `json:"description,omitempty"`
	AltText      string    `json:"alt_text,omitempty"`
	IsPrimary    bool      `json:"is_primary"`
	TourID       *int      `json:"tour_id,omitempty"`
	UploadedAt   time.Time `json:"uploaded_at"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
}

// UploadImagesResponse - resposta do upload de imagens
type UploadImagesResponse struct {
	Success       bool            `json:"success"`
	Images        []ImageResponse `json:"images"`
	TotalUploaded int             `json:"total_uploaded"`
	TotalSize     int64           `json:"total_size"`
}

// ListImagesResponse - resposta da listagem de imagens
type ListImagesResponse struct {
	Success    bool               `json:"success"`
	Images     []ImageResponse    `json:"images"`
	Pagination PaginationResponse `json:"pagination"`
}

// UpdateImageResponse - resposta da atualização de imagem
type UpdateImageResponse struct {
	Success bool          `json:"success"`
	Image   ImageResponse `json:"image"`
}

// DeleteImageResponse - resposta da exclusão de imagem
type DeleteImageResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ReorderImagesResponse - resposta da reordenação de imagens
type ReorderImagesResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ImageInfoResponse - resposta com informações detalhadas de uma imagem
type ImageInfoResponse struct {
	Success bool          `json:"success"`
	Image   ImageResponse `json:"image"`
	Usage   ImageUsage    `json:"usage,omitempty"`
}

// ImageUsage - informações de uso da imagem
type ImageUsage struct {
	TourName string `json:"tour_name,omitempty"`
	IsUsed   bool   `json:"is_used"`
}

// BatchDeleteImagesResponse - resposta da exclusão em lote
type BatchDeleteImagesResponse struct {
	Success      bool               `json:"success"`
	DeletedCount int                `json:"deleted_count"`
	FailedCount  int                `json:"failed_count"`
	Message      string             `json:"message"`
	Errors       []BatchDeleteError `json:"errors,omitempty"`
}

// BatchDeleteError - erro individual na exclusão em lote
type BatchDeleteError struct {
	ImageID int    `json:"image_id"`
	Error   string `json:"error"`
}

// ImageUploadResult - resultado individual do upload
type ImageUploadResult struct {
	Image   ImageResponse `json:"image"`
	Success bool          `json:"success"`
	Error   string        `json:"error,omitempty"`
}

// UploadSummary - resumo do upload
type UploadSummary struct {
	TotalFiles    int   `json:"total_files"`
	SuccessCount  int   `json:"success_count"`
	ErrorCount    int   `json:"error_count"`
	TotalSize     int64 `json:"total_size"`
	ProcessedSize int64 `json:"processed_size"`
}

// ImageValidationError - erro de validação de imagem
type ImageValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   string `json:"value,omitempty"`
}

// ImageProcessingError - erro no processamento de imagem
type ImageProcessingError struct {
	Filename string `json:"filename"`
	Error    string `json:"error"`
	Code     string `json:"code,omitempty"`
}

// ThumbnailInfo - informações do thumbnail
type ThumbnailInfo struct {
	URL       string `json:"url"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Size      int    `json:"size"`
	Format    string `json:"format"`
	Generated bool   `json:"generated"`
}

// ImageMetadata - metadados da imagem
type ImageMetadata struct {
	Filename     string        `json:"filename"`
	OriginalName string        `json:"original_name"`
	Size         int           `json:"size"`
	Width        int           `json:"width"`
	Height       int           `json:"height"`
	Format       string        `json:"format"`
	AspectRatio  float64       `json:"aspect_ratio"`
	IsLandscape  bool          `json:"is_landscape"`
	IsPortrait   bool          `json:"is_portrait"`
	IsSquare     bool          `json:"is_square"`
	Thumbnail    ThumbnailInfo `json:"thumbnail"`
}

// ImageStats - estatísticas de imagens
type ImageStats struct {
	TotalImages   int            `json:"total_images"`
	TotalSize     int64          `json:"total_size"`
	AverageSize   float64        `json:"average_size"`
	FormatCounts  map[string]int `json:"format_counts"`
	PrimaryImages int            `json:"primary_images"`
	UnusedImages  int            `json:"unused_images"`
	RecentUploads int            `json:"recent_uploads"`
}
