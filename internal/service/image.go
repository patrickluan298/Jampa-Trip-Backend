package service

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jampa_trip/internal/contract"
	"github.com/jampa_trip/internal/model"
	"github.com/jampa_trip/internal/repository"
	"github.com/jampa_trip/pkg/util"
	"golang.org/x/image/draw"
	"golang.org/x/image/webp"
	"gorm.io/gorm"
)

// ImageService - objeto de contexto
type ImageService struct {
	ImageRepository *repository.ImageRepository
}

// ImageServiceNew - construtor do objeto
func ImageServiceNew(DB *gorm.DB) *ImageService {
	return &ImageService{
		ImageRepository: repository.ImageRepositoryNew(DB),
	}
}

// UploadImages - faz upload de múltiplas imagens
func (s *ImageService) UploadImages(files []*multipart.FileHeader, request *contract.UploadImagesRequest, userID int) (*contract.UploadImagesResponse, error) {
	if len(files) == 0 {
		return nil, util.WrapError("Nenhum arquivo enviado", nil, http.StatusBadRequest)
	}

	if len(files) > 10 {
		return nil, util.WrapError("Máximo de 10 imagens por upload", nil, http.StatusBadRequest)
	}

	uploadDir := fmt.Sprintf("uploads/images/%d", userID)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, util.WrapError("Erro ao criar diretório de upload", err, http.StatusInternalServerError)
	}

	var uploadedImages []contract.ImageResponse
	var totalSize int64
	successCount := 0

	for _, fileHeader := range files {
		if err := s.validateImageFile(fileHeader); err != nil {
			continue
		}

		imageData, err := s.processImageFile(fileHeader, uploadDir, userID, request.TourID)
		if err != nil {
			continue
		}

		if err := s.ImageRepository.Create(imageData); err != nil {
			s.cleanupFiles(imageData.URL, imageData.ThumbnailURL)
			continue
		}

		uploadedImages = append(uploadedImages, contract.ImageResponse{
			ID:           imageData.ID,
			Filename:     imageData.Filename,
			OriginalName: imageData.OriginalName,
			URL:          imageData.URL,
			ThumbnailURL: imageData.ThumbnailURL,
			Size:         imageData.Size,
			Width:        imageData.Width,
			Height:       imageData.Height,
			Format:       imageData.Format,
			Description:  imageData.Description,
			IsPrimary:    imageData.IsPrimary,
			TourID:       imageData.TourID,
			UploadedAt:   imageData.UploadedAt,
		})

		totalSize += int64(imageData.Size)
		successCount++
	}

	if successCount == 0 {
		return nil, util.WrapError("Nenhuma imagem foi processada com sucesso", nil, http.StatusBadRequest)
	}

	response := &contract.UploadImagesResponse{
		Success:       true,
		Images:        uploadedImages,
		TotalUploaded: successCount,
		TotalSize:     totalSize,
	}

	return response, nil
}

// ListImages - lista imagens do usuário
func (s *ImageService) ListImages(request *contract.ListImagesRequest, userID int) (*contract.ListImagesResponse, error) {
	config := util.NormalizePagination(request.Page, request.Limit)
	page := config.Page
	limit := config.Limit

	images, total, err := s.ImageRepository.List(userID, request.TourID, request.Format, "uploaded_at", page, limit)
	if err != nil {
		return nil, util.WrapError("Erro ao buscar imagens", err, http.StatusInternalServerError)
	}

	var imagesResponse []contract.ImageResponse
	for _, img := range images {
		imagesResponse = append(imagesResponse, contract.ImageResponse{
			ID:           img.ID,
			Filename:     img.Filename,
			OriginalName: img.OriginalName,
			URL:          img.URL,
			ThumbnailURL: img.ThumbnailURL,
			Size:         img.Size,
			Width:        img.Width,
			Height:       img.Height,
			Format:       img.Format,
			Description:  img.Description,
			AltText:      img.AltText,
			IsPrimary:    img.IsPrimary,
			TourID:       img.TourID,
			UploadedAt:   img.UploadedAt,
			UpdatedAt:    img.UpdatedAt,
		})
	}

	totalPages := util.CalculateTotalPages(total, limit)

	response := &contract.ListImagesResponse{
		Success: true,
		Images:  imagesResponse,
		Pagination: contract.PaginationResponse{
			CurrentPage:  page,
			TotalPages:   totalPages,
			TotalItems:   int(total),
			ItemsPerPage: limit,
		},
	}

	return response, nil
}

// DeleteImage - deleta uma imagem
func (s *ImageService) DeleteImage(imageID, userID int) (*contract.DeleteImageResponse, error) {
	exists, err := s.ImageRepository.IsOwnedByUser(imageID, userID)
	if err != nil {
		return nil, util.WrapError("Erro ao verificar propriedade da imagem", err, http.StatusInternalServerError)
	}
	if !exists {
		return nil, util.WrapError("Imagem não encontrada ou você não tem permissão para deletá-la", nil, http.StatusNotFound)
	}

	isUsed, err := s.ImageRepository.IsUsedInActiveTour(imageID)
	if err != nil {
		return nil, util.WrapError("Erro ao verificar uso da imagem", err, http.StatusInternalServerError)
	}
	if isUsed {
		return nil, util.WrapError("Não é possível deletar imagem que está sendo usada em passeio ativo", nil, http.StatusConflict)
	}

	image, err := s.ImageRepository.GetByIDAndUser(imageID, userID)
	if err != nil {
		return nil, util.WrapError("Erro ao buscar dados da imagem", err, http.StatusInternalServerError)
	}

	if err := s.ImageRepository.Delete(imageID); err != nil {
		return nil, util.WrapError("Erro ao deletar imagem do banco", err, http.StatusInternalServerError)
	}

	s.cleanupFiles(image.URL, image.ThumbnailURL)

	response := &contract.DeleteImageResponse{
		Success: true,
		Message: "Imagem deletada com sucesso",
	}

	return response, nil
}

// UpdateImage - atualiza metadados de uma imagem
func (s *ImageService) UpdateImage(imageID int, request *contract.UpdateImageRequest, userID int) (*contract.UpdateImageResponse, error) {
	image, err := s.ImageRepository.GetByIDAndUser(imageID, userID)
	if err != nil {
		return nil, util.WrapError("Imagem não encontrada ou você não tem permissão para editá-la", err, http.StatusNotFound)
	}

	if request.Description != "" {
		image.Description = request.Description
	}
	if request.AltText != "" {
		image.AltText = request.AltText
	}
	if request.IsPrimary != nil && *request.IsPrimary {
		if image.TourID != nil {
			if err := s.ImageRepository.RemovePrimaryFromTour(*image.TourID, userID, imageID); err != nil {
				return nil, util.WrapError("Erro ao atualizar outras imagens", err, http.StatusInternalServerError)
			}
		}
		image.IsPrimary = true
	}

	if err := s.ImageRepository.Update(image); err != nil {
		return nil, util.WrapError("Erro ao atualizar imagem", err, http.StatusInternalServerError)
	}

	updatedImage, err := s.ImageRepository.GetByIDAndUser(imageID, userID)
	if err != nil {
		return nil, util.WrapError("Erro ao buscar imagem atualizada", err, http.StatusInternalServerError)
	}

	response := &contract.UpdateImageResponse{
		Success: true,
		Image: contract.ImageResponse{
			ID:           updatedImage.ID,
			Filename:     updatedImage.Filename,
			OriginalName: updatedImage.OriginalName,
			URL:          updatedImage.URL,
			ThumbnailURL: updatedImage.ThumbnailURL,
			Size:         updatedImage.Size,
			Width:        updatedImage.Width,
			Height:       updatedImage.Height,
			Format:       updatedImage.Format,
			Description:  updatedImage.Description,
			AltText:      updatedImage.AltText,
			IsPrimary:    updatedImage.IsPrimary,
			TourID:       updatedImage.TourID,
			UploadedAt:   updatedImage.UploadedAt,
			UpdatedAt:    updatedImage.UpdatedAt,
		},
	}

	return response, nil
}

// ReorderImages - reordena imagens de um passeio
func (s *ImageService) ReorderImages(request *contract.ReorderImagesRequest, userID int) (*contract.ReorderImagesResponse, error) {
	images, err := s.ImageRepository.GetByIDs(request.ImageIDs, userID)
	if err != nil {
		return nil, util.WrapError("Erro ao buscar imagens", err, http.StatusInternalServerError)
	}

	if len(images) != len(request.ImageIDs) {
		return nil, util.WrapError("Algumas imagens não foram encontradas ou não pertencem a você", nil, http.StatusForbidden)
	}

	for _, img := range images {
		if img.TourID == nil || *img.TourID != request.TourID {
			return nil, util.WrapError("Todas as imagens devem pertencer ao mesmo passeio", nil, http.StatusBadRequest)
		}
	}

	if err := s.ImageRepository.BatchUpdateSortOrder(request.ImageIDs, userID); err != nil {
		return nil, util.WrapError("Erro ao reordenar imagens", err, http.StatusInternalServerError)
	}

	response := &contract.ReorderImagesResponse{
		Success: true,
		Message: "Ordem das imagens atualizada com sucesso",
	}

	return response, nil
}

// GetImageInfo - obtém informações detalhadas de uma imagem
func (s *ImageService) GetImageInfo(imageID, userID int) (*contract.ImageInfoResponse, error) {
	image, tourName, err := s.ImageRepository.GetWithTourInfo(imageID)
	if err != nil {
		return nil, util.WrapError("Imagem não encontrada", err, http.StatusNotFound)
	}

	if image.UserID != userID {
		return nil, util.WrapError("Você não tem permissão para acessar esta imagem", nil, http.StatusForbidden)
	}

	_, isUsed, err := s.ImageRepository.GetImageUsage(imageID)
	if err != nil {
		return nil, util.WrapError("Erro ao buscar informações de uso", err, http.StatusInternalServerError)
	}

	response := &contract.ImageInfoResponse{
		Success: true,
		Image: contract.ImageResponse{
			ID:           image.ID,
			Filename:     image.Filename,
			OriginalName: image.OriginalName,
			URL:          image.URL,
			ThumbnailURL: image.ThumbnailURL,
			Size:         image.Size,
			Width:        image.Width,
			Height:       image.Height,
			Format:       image.Format,
			Description:  image.Description,
			AltText:      image.AltText,
			IsPrimary:    image.IsPrimary,
			TourID:       image.TourID,
			UploadedAt:   image.UploadedAt,
			UpdatedAt:    image.UpdatedAt,
		},
		Usage: contract.ImageUsage{
			TourName: tourName,
			IsUsed:   isUsed,
		},
	}

	return response, nil
}

// BatchDeleteImages - deleta múltiplas imagens
func (s *ImageService) BatchDeleteImages(request *contract.BatchDeleteImagesRequest, userID int) (*contract.BatchDeleteImagesResponse, error) {
	var deletedCount, failedCount int
	var errors []contract.BatchDeleteError

	for _, imageID := range request.ImageIDs {
		exists, err := s.ImageRepository.IsOwnedByUser(imageID, userID)
		if err != nil || !exists {
			failedCount++
			errors = append(errors, contract.BatchDeleteError{
				ImageID: imageID,
				Error:   "Imagem não encontrada ou sem permissão",
			})
			continue
		}

		isUsed, err := s.ImageRepository.IsUsedInActiveTour(imageID)
		if err != nil || isUsed {
			failedCount++
			errors = append(errors, contract.BatchDeleteError{
				ImageID: imageID,
				Error:   "Imagem está sendo usada em passeio ativo",
			})
			continue
		}

		image, err := s.ImageRepository.GetByIDAndUser(imageID, userID)
		if err != nil {
			failedCount++
			errors = append(errors, contract.BatchDeleteError{
				ImageID: imageID,
				Error:   "Erro ao buscar dados da imagem",
			})
			continue
		}

		if err := s.ImageRepository.Delete(imageID); err != nil {
			failedCount++
			errors = append(errors, contract.BatchDeleteError{
				ImageID: imageID,
				Error:   "Erro ao deletar do banco",
			})
			continue
		}

		s.cleanupFiles(image.URL, image.ThumbnailURL)
		deletedCount++
	}

	message := fmt.Sprintf("%d imagens deletadas", deletedCount)
	if failedCount > 0 {
		message += fmt.Sprintf(", %d falharam", failedCount)
	}

	response := &contract.BatchDeleteImagesResponse{
		Success:      true,
		DeletedCount: deletedCount,
		FailedCount:  failedCount,
		Message:      message,
	}

	if len(errors) > 0 {
		response.Errors = errors
	}

	return response, nil
}

// validateImageFile - valida um arquivo de imagem
func (s *ImageService) validateImageFile(fileHeader *multipart.FileHeader) error {
	const maxSize = 10 * 1024 * 1024
	if fileHeader.Size > maxSize {
		return fmt.Errorf("arquivo muito grande: %d bytes (máximo: %d bytes)", fileHeader.Size, maxSize)
	}

	contentType := fileHeader.Header.Get("Content-Type")
	validTypes := []string{
		"image/jpeg", "image/jpg", "image/png", "image/gif", "image/webp",
	}

	valid := false
	for _, validType := range validTypes {
		if contentType == validType {
			valid = true
			break
		}
	}

	if !valid {
		return fmt.Errorf("tipo de arquivo não suportado: %s", contentType)
	}

	return nil
}

// processImageFile - processa um arquivo de imagem
func (s *ImageService) processImageFile(fileHeader *multipart.FileHeader, uploadDir string, userID int, tourID *int) (*model.Image, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileData, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	format, err := s.detectImageFormat(fileData)
	if err != nil {
		return nil, err
	}

	img, err := s.decodeImage(fileData, format)
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	uniqueID := uuid.New().String()
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%d_%s_%s", timestamp, uniqueID[:8], fileHeader.Filename)

	originalPath := filepath.Join(uploadDir, filename)
	if err := os.WriteFile(originalPath, fileData, 0644); err != nil {
		return nil, err
	}

	thumbnailPath := filepath.Join(uploadDir, "thumb_"+filename)
	thumbnailURL, err := s.generateThumbnail(img, thumbnailPath, format)
	if err != nil {
		thumbnailURL = ""
	}

	baseURL := "http://localhost:8080" // TODO: configurar via variável de ambiente
	imageURL := fmt.Sprintf("%s/uploads/images/%d/%s", baseURL, userID, filename)
	thumbnailImageURL := ""
	if thumbnailURL != "" {
		thumbnailImageURL = fmt.Sprintf("%s/uploads/images/%d/thumb_%s", baseURL, userID, filename)
	}

	image := &model.Image{
		UserID:       userID,
		TourID:       tourID,
		Filename:     filename,
		OriginalName: fileHeader.Filename,
		URL:          imageURL,
		ThumbnailURL: thumbnailImageURL,
		Size:         int(fileHeader.Size),
		Width:        width,
		Height:       height,
		Format:       format,
		Description:  "",
		AltText:      "",
		IsPrimary:    false,
		SortOrder:    0,
		UploadedAt:   time.Now(),
		UpdatedAt:    time.Now(),
	}

	return image, nil
}

// detectImageFormat - detecta o formato da imagem
func (s *ImageService) detectImageFormat(data []byte) (string, error) {
	if len(data) < 8 {
		return "", fmt.Errorf("arquivo muito pequeno")
	}

	if len(data) >= 8 {
		if data[0] == 0xFF && data[1] == 0xD8 {
			return "jpg", nil
		}
		if data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47 {
			return "png", nil
		}
		if len(data) >= 6 && string(data[0:6]) == "GIF87a" || string(data[0:6]) == "GIF89a" {
			return "gif", nil
		}
		if len(data) >= 12 && string(data[0:4]) == "RIFF" && string(data[8:12]) == "WEBP" {
			return "webp", nil
		}
	}

	return "", fmt.Errorf("formato de imagem não suportado")
}

// decodeImage - decodifica uma imagem
func (s *ImageService) decodeImage(data []byte, format string) (image.Image, error) {
	switch format {
	case "jpg", "jpeg":
		return jpeg.Decode(strings.NewReader(string(data)))
	case "png":
		return png.Decode(strings.NewReader(string(data)))
	case "gif":
		return gif.Decode(strings.NewReader(string(data)))
	case "webp":
		return webp.Decode(strings.NewReader(string(data)))
	default:
		return nil, fmt.Errorf("formato não suportado: %s", format)
	}
}

// generateThumbnail - gera thumbnail da imagem
func (s *ImageService) generateThumbnail(img image.Image, thumbnailPath, format string) (string, error) {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	thumbWidth, thumbHeight := s.CalculateThumbnailDimensions(width, height)

	// Criar nova imagem redimensionada
	thumb := image.NewRGBA(image.Rect(0, 0, thumbWidth, thumbHeight))

	// Redimensionar usando interpolação Lanczos para melhor qualidade
	draw.CatmullRom.Scale(thumb, thumb.Bounds(), img, img.Bounds(), draw.Over, nil)

	file, err := os.Create(thumbnailPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	switch format {
	case "jpg", "jpeg":
		return thumbnailPath, jpeg.Encode(file, thumb, &jpeg.Options{Quality: 80})
	case "png":
		return thumbnailPath, png.Encode(file, thumb)
	default:
		return thumbnailPath, jpeg.Encode(file, thumb, &jpeg.Options{Quality: 80})
	}
}

// CalculateThumbnailDimensions - calcula as dimensões do thumbnail
func (s *ImageService) CalculateThumbnailDimensions(width, height int) (int, int) {
	const maxThumbnailSize = 300

	if width <= maxThumbnailSize && height <= maxThumbnailSize {
		return width, height
	}

	ratio := float64(maxThumbnailSize) / float64(util.Max(width, height))
	newWidth := int(float64(width) * ratio)
	newHeight := int(float64(height) * ratio)

	return newWidth, newHeight
}

// cleanupFiles - remove arquivos do sistema de arquivos
func (s *ImageService) cleanupFiles(imageURL, thumbnailURL string) {
	if imageURL != "" {
		s.RemoveFileFromURL(imageURL)
	}
	if thumbnailURL != "" {
		s.RemoveFileFromURL(thumbnailURL)
	}
}

// RemoveFileFromURL - remove arquivo do sistema de arquivos baseado na URL
func (s *ImageService) RemoveFileFromURL(fileURL string) {
	if fileURL == "" {
		return
	}

	// Extrair o caminho do arquivo a partir da URL
	parsedURL, err := url.Parse(fileURL)
	if err != nil {
		// Log silencioso - não queremos falhar a operação principal
		return
	}

	// Remover o prefixo da URL base para obter o caminho relativo
	// URL: http://localhost:8080/uploads/images/{user_id}/{filename}
	// Path: uploads/images/{user_id}/{filename}
	filePath := strings.TrimPrefix(parsedURL.Path, "/")

	// Verificar se o arquivo existe antes de tentar removê-lo
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Arquivo não existe, nada a fazer
		return
	}

	// Remover o arquivo
	if err := os.Remove(filePath); err != nil {
		// Log silencioso - não queremos falhar a operação principal
		// Em produção, poderia ser logado para monitoramento
	}
}
