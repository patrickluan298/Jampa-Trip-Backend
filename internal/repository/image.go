package repository

import (
	"database/sql"
	"strconv"
	"strings"

	"github.com/jampa_trip/internal/model"
	"github.com/jampa_trip/internal/query"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// ImageRepository - objeto de contexto
type ImageRepository struct {
	DB *gorm.DB
}

// ImageRepositoryNew - construtor do objeto
func ImageRepositoryNew(DB *gorm.DB) *ImageRepository {
	return &ImageRepository{
		DB: DB,
	}
}

// Create - cria uma nova imagem
func (r *ImageRepository) Create(image *model.Image) error {
	err := r.DB.Raw(query.CreateImage,
		image.UserID,
		image.TourID,
		image.Filename,
		image.OriginalName,
		image.URL,
		image.ThumbnailURL,
		image.Size,
		image.Width,
		image.Height,
		image.Format,
		image.Description,
		image.AltText,
		image.IsPrimary,
		image.SortOrder,
	).Row().Scan(&image.ID, &image.UploadedAt, &image.UpdatedAt)

	return err
}

// GetByID - busca uma imagem pelo ID
func (r *ImageRepository) GetByID(id int) (*model.Image, error) {
	image := &model.Image{}

	err := r.DB.Raw(query.GetImageByID, id).Row().Scan(
		&image.ID,
		&image.UserID,
		&image.TourID,
		&image.Filename,
		&image.OriginalName,
		&image.URL,
		&image.ThumbnailURL,
		&image.Size,
		&image.Width,
		&image.Height,
		&image.Format,
		&image.Description,
		&image.AltText,
		&image.IsPrimary,
		&image.SortOrder,
		&image.UploadedAt,
		&image.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return image, nil
}

// GetByIDAndUser - busca uma imagem pelo ID e usuário
func (r *ImageRepository) GetByIDAndUser(id, userID int) (*model.Image, error) {
	image := &model.Image{}

	err := r.DB.Raw(query.GetImageByIDAndUser, id, userID).Row().Scan(
		&image.ID,
		&image.UserID,
		&image.TourID,
		&image.Filename,
		&image.OriginalName,
		&image.URL,
		&image.ThumbnailURL,
		&image.Size,
		&image.Width,
		&image.Height,
		&image.Format,
		&image.Description,
		&image.AltText,
		&image.IsPrimary,
		&image.SortOrder,
		&image.UploadedAt,
		&image.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return image, nil
}

// Update - atualiza uma imagem
func (r *ImageRepository) Update(image *model.Image) error {
	err := r.DB.Raw(query.UpdateImage,
		image.ID,
		image.TourID,
		image.Description,
		image.AltText,
		image.IsPrimary,
	).Row().Scan(&image.UpdatedAt)

	return err
}

// Delete - deleta uma imagem
func (r *ImageRepository) Delete(id int) error {
	err := r.DB.Raw(query.DeleteImage, id).Row().Scan()
	return err
}

// List - lista imagens do usuário com filtros
func (r *ImageRepository) List(userID int, tourID *int, format string, sortBy string, page, limit int) ([]*model.Image, int64, error) {
	offset := (page - 1) * limit

	rows, err := r.DB.Raw(query.ListImages, userID, tourID, format, sortBy, limit, offset).Rows()
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var images []*model.Image
	for rows.Next() {
		image := &model.Image{}

		err := rows.Scan(
			&image.ID,
			&image.UserID,
			&image.TourID,
			&image.Filename,
			&image.OriginalName,
			&image.URL,
			&image.ThumbnailURL,
			&image.Size,
			&image.Width,
			&image.Height,
			&image.Format,
			&image.Description,
			&image.AltText,
			&image.IsPrimary,
			&image.SortOrder,
			&image.UploadedAt,
			&image.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		images = append(images, image)
	}

	var total int64
	err = r.DB.Raw(query.CountImages, userID, tourID, format).Row().Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return images, total, nil
}

// ListByTour - lista imagens de um passeio específico
func (r *ImageRepository) ListByTour(tourID, userID int) ([]*model.Image, error) {
	rows, err := r.DB.Raw(query.ListImagesByTour, tourID, userID).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []*model.Image
	for rows.Next() {
		image := &model.Image{}

		err := rows.Scan(
			&image.ID,
			&image.UserID,
			&image.TourID,
			&image.Filename,
			&image.OriginalName,
			&image.URL,
			&image.ThumbnailURL,
			&image.Size,
			&image.Width,
			&image.Height,
			&image.Format,
			&image.Description,
			&image.AltText,
			&image.IsPrimary,
			&image.SortOrder,
			&image.UploadedAt,
			&image.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		images = append(images, image)
	}

	return images, nil
}

// IsOwnedByUser - verifica se a imagem pertence ao usuário
func (r *ImageRepository) IsOwnedByUser(imageID, userID int) (bool, error) {
	var exists bool
	err := r.DB.Raw(query.IsImageOwnedByUser, imageID, userID).Row().Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// IsUsedInActiveTour - verifica se a imagem está sendo usada em passeio ativo
func (r *ImageRepository) IsUsedInActiveTour(imageID int) (bool, error) {
	var isUsed bool
	err := r.DB.Raw(query.IsImageUsedInActiveTour, imageID).Row().Scan(&isUsed)
	if err != nil {
		return false, err
	}
	return isUsed, nil
}

// GetImageUsage - obtém informações de uso da imagem
func (r *ImageRepository) GetImageUsage(imageID int) (string, bool, error) {
	var tourName sql.NullString
	var isUsed bool

	err := r.DB.Raw(query.GetImageUsage, imageID).Row().Scan(&tourName, &isUsed)
	if err != nil {
		return "", false, err
	}

	if tourName.Valid {
		return tourName.String, isUsed, nil
	}

	return "", isUsed, nil
}

// UpdateSortOrder - atualiza a ordem de uma imagem
func (r *ImageRepository) UpdateSortOrder(imageID, sortOrder, userID int) error {
	err := r.DB.Raw(query.UpdateImageSortOrder, imageID, sortOrder, userID).Row().Scan()
	return err
}

// BatchUpdateSortOrder - atualiza ordem de múltiplas imagens
func (r *ImageRepository) BatchUpdateSortOrder(imageIDs []int, userID int) error {
	var caseStatements []string
	for i, id := range imageIDs {
		caseStatements = append(caseStatements, "WHEN "+strconv.Itoa(id)+" THEN "+strconv.Itoa(i+1))
	}

	caseQuery := strings.Join(caseStatements, " ")

	query := strings.Replace(query.BatchUpdateSortOrder, "$1", caseQuery, 1)

	err := r.DB.Raw(query, pq.Array(imageIDs), userID).Row().Scan()
	return err
}

// RemovePrimaryFromTour - remove flag primary de outras imagens do mesmo passeio
func (r *ImageRepository) RemovePrimaryFromTour(tourID, userID, excludeImageID int) error {
	err := r.DB.Raw(query.RemovePrimaryFromTour, tourID, userID, excludeImageID).Row().Scan()
	return err
}

// SetImageAsPrimary - define uma imagem como primary
func (r *ImageRepository) SetImageAsPrimary(imageID, userID int) error {
	err := r.DB.Raw(query.SetImageAsPrimary, imageID, userID).Row().Scan()
	return err
}

// GetImageStats - obtém estatísticas das imagens do usuário
func (r *ImageRepository) GetImageStats(userID int) (int, int64, float64, int, int, int, error) {
	var totalImages, primaryImages, unusedImages, recentUploads int
	var totalSize int64
	var averageSize float64

	err := r.DB.Raw(query.GetImageStats, userID).Row().Scan(
		&totalImages,
		&totalSize,
		&averageSize,
		&primaryImages,
		&unusedImages,
		&recentUploads,
	)

	return totalImages, totalSize, averageSize, primaryImages, unusedImages, recentUploads, err
}

// GetImageFormatCounts - conta imagens por formato
func (r *ImageRepository) GetImageFormatCounts(userID int) (map[string]int, error) {
	rows, err := r.DB.Raw(query.GetImageFormatCounts, userID).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	formatCounts := make(map[string]int)
	for rows.Next() {
		var format string
		var count int

		err := rows.Scan(&format, &count)
		if err != nil {
			return nil, err
		}
		formatCounts[format] = count
	}

	return formatCounts, nil
}

// GetByIDs - busca imagens por IDs
func (r *ImageRepository) GetByIDs(imageIDs []int, userID int) ([]*model.Image, error) {
	rows, err := r.DB.Raw(query.GetImagesByIDs, pq.Array(imageIDs), userID).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []*model.Image
	for rows.Next() {
		image := &model.Image{}

		err := rows.Scan(
			&image.ID,
			&image.UserID,
			&image.TourID,
			&image.Filename,
			&image.OriginalName,
			&image.URL,
			&image.ThumbnailURL,
			&image.Size,
			&image.Width,
			&image.Height,
			&image.Format,
			&image.Description,
			&image.AltText,
			&image.IsPrimary,
			&image.SortOrder,
			&image.UploadedAt,
			&image.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		images = append(images, image)
	}

	return images, nil
}

// BatchDelete - deleta múltiplas imagens
func (r *ImageRepository) BatchDelete(imageIDs []int, userID int) error {
	err := r.DB.Raw(query.BatchDeleteImages, pq.Array(imageIDs), userID).Row().Scan()
	return err
}

// GetByTourID - busca imagens de um passeio específico
func (r *ImageRepository) GetByTourID(tourID int) ([]*model.Image, error) {
	rows, err := r.DB.Raw(query.GetImagesByTourID, tourID).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []*model.Image
	for rows.Next() {
		image := &model.Image{}

		err := rows.Scan(
			&image.ID,
			&image.UserID,
			&image.TourID,
			&image.Filename,
			&image.OriginalName,
			&image.URL,
			&image.ThumbnailURL,
			&image.Size,
			&image.Width,
			&image.Height,
			&image.Format,
			&image.Description,
			&image.AltText,
			&image.IsPrimary,
			&image.SortOrder,
			&image.UploadedAt,
			&image.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		images = append(images, image)
	}

	return images, nil
}

// Exists - verifica se a imagem existe
func (r *ImageRepository) Exists(id int) (bool, error) {
	var exists bool
	err := r.DB.Raw(query.CheckImageExists, id).Row().Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// GetWithTourInfo - busca imagem com informações do passeio
func (r *ImageRepository) GetWithTourInfo(id int) (*model.Image, string, error) {
	image := &model.Image{}
	var tourName sql.NullString

	err := r.DB.Raw(query.GetImageWithTourInfo, id).Row().Scan(
		&image.ID,
		&image.UserID,
		&image.TourID,
		&image.Filename,
		&image.OriginalName,
		&image.URL,
		&image.ThumbnailURL,
		&image.Size,
		&image.Width,
		&image.Height,
		&image.Format,
		&image.Description,
		&image.AltText,
		&image.IsPrimary,
		&image.SortOrder,
		&image.UploadedAt,
		&image.UpdatedAt,
		&tourName,
	)

	if err != nil {
		return nil, "", err
	}

	if tourName.Valid {
		return image, tourName.String, nil
	}

	return image, "", nil
}

// GetRecent - busca imagens recentes do usuário
func (r *ImageRepository) GetRecent(userID, limit int) ([]*model.Image, error) {
	rows, err := r.DB.Raw(query.GetRecentImages, userID, limit).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []*model.Image
	for rows.Next() {
		image := &model.Image{}

		err := rows.Scan(
			&image.ID,
			&image.UserID,
			&image.TourID,
			&image.Filename,
			&image.OriginalName,
			&image.URL,
			&image.ThumbnailURL,
			&image.Size,
			&image.Width,
			&image.Height,
			&image.Format,
			&image.Description,
			&image.AltText,
			&image.IsPrimary,
			&image.SortOrder,
			&image.UploadedAt,
			&image.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		images = append(images, image)
	}

	return images, nil
}

// Search - busca imagens por termo
func (r *ImageRepository) Search(userID int, searchTerm string, page, limit int) ([]*model.Image, int64, error) {
	offset := (page - 1) * limit
	searchPattern := "%" + searchTerm + "%"

	rows, err := r.DB.Raw(query.SearchImages, userID, searchPattern, limit, offset).Rows()
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var images []*model.Image
	for rows.Next() {
		image := &model.Image{}

		err := rows.Scan(
			&image.ID,
			&image.UserID,
			&image.TourID,
			&image.Filename,
			&image.OriginalName,
			&image.URL,
			&image.ThumbnailURL,
			&image.Size,
			&image.Width,
			&image.Height,
			&image.Format,
			&image.Description,
			&image.AltText,
			&image.IsPrimary,
			&image.SortOrder,
			&image.UploadedAt,
			&image.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		images = append(images, image)
	}

	var total int64
	err = r.DB.Raw(query.CountSearchImages, userID, searchPattern).Row().Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return images, total, nil
}
