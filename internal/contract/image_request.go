package contract

import (
	"errors"
	"strconv"
	"strings"
)

// UploadImagesRequest - request para upload de múltiplas imagens
type UploadImagesRequest struct {
	TourID      *int   `json:"tour_id,omitempty"`
	Description string `json:"description,omitempty"`
}

func (r *UploadImagesRequest) Validate() error {
	if r.TourID != nil && *r.TourID <= 0 {
		return errors.New("tour_id deve ser um número positivo")
	}

	if len(r.Description) > 500 {
		return errors.New("descrição deve ter no máximo 500 caracteres")
	}

	return nil
}

// UpdateImageRequest - request para atualizar metadados de uma imagem
type UpdateImageRequest struct {
	Description string `json:"description,omitempty"`
	AltText     string `json:"alt_text,omitempty"`
	IsPrimary   *bool  `json:"is_primary,omitempty"`
}

func (r *UpdateImageRequest) Validate() error {
	if len(r.Description) > 500 {
		return errors.New("descrição deve ter no máximo 500 caracteres")
	}

	if len(r.AltText) > 255 {
		return errors.New("alt_text deve ter no máximo 255 caracteres")
	}

	return nil
}

// ReorderImagesRequest - request para reordenar imagens
type ReorderImagesRequest struct {
	ImageIDs []int `json:"image_ids" validate:"required,min=1"`
	TourID   int   `json:"tour_id" validate:"required,min=1"`
}

func (r *ReorderImagesRequest) Validate() error {
	if len(r.ImageIDs) == 0 {
		return errors.New("image_ids não pode estar vazio")
	}

	if len(r.ImageIDs) > 50 {
		return errors.New("máximo de 50 imagens por reordenação")
	}

	if r.TourID <= 0 {
		return errors.New("tour_id deve ser um número positivo")
	}

	seen := make(map[int]bool)
	for _, id := range r.ImageIDs {
		if id <= 0 {
			return errors.New("todos os image_ids devem ser números positivos")
		}
		if seen[id] {
			return errors.New("image_ids não pode conter IDs duplicados")
		}
		seen[id] = true
	}

	return nil
}

// BatchDeleteImagesRequest - request para deletar múltiplas imagens
type BatchDeleteImagesRequest struct {
	ImageIDs []int `json:"image_ids" validate:"required,min=1"`
}

func (r *BatchDeleteImagesRequest) Validate() error {
	if len(r.ImageIDs) == 0 {
		return errors.New("image_ids não pode estar vazio")
	}

	if len(r.ImageIDs) > 100 {
		return errors.New("máximo de 100 imagens por exclusão em lote")
	}

	seen := make(map[int]bool)
	for _, id := range r.ImageIDs {
		if id <= 0 {
			return errors.New("todos os image_ids devem ser números positivos")
		}
		if seen[id] {
			return errors.New("image_ids não pode conter IDs duplicados")
		}
		seen[id] = true
	}

	return nil
}

// ListImagesRequest - request para listar imagens
type ListImagesRequest struct {
	TourID *int   `json:"tour_id,omitempty"`
	Format string `json:"format,omitempty"`
	Page   int    `json:"page,omitempty"`
	Limit  int    `json:"limit,omitempty"`
}

func (r *ListImagesRequest) Validate() error {
	if r.TourID != nil && *r.TourID <= 0 {
		return errors.New("tour_id deve ser um número positivo")
	}

	if r.Format != "" {
		validFormats := []string{"jpg", "jpeg", "png", "gif", "webp"}
		format := strings.ToLower(r.Format)
		valid := false
		for _, vf := range validFormats {
			if format == vf {
				valid = true
				break
			}
		}
		if !valid {
			return errors.New("formato deve ser jpg, jpeg, png, gif ou webp")
		}
	}

	if r.Page < 0 {
		return errors.New("página deve ser um número não negativo")
	}

	if r.Limit < 0 {
		return errors.New("limite deve ser um número não negativo")
	}

	if r.Limit > 100 {
		return errors.New("limite máximo é 100 itens por página")
	}

	return nil
}

// ParseQueryParams - converte query parameters para o request
func (r *ListImagesRequest) ParseQueryParams(tourIDStr, format, pageStr, limitStr string) error {
	if tourIDStr != "" {
		tourID, err := strconv.Atoi(tourIDStr)
		if err != nil {
			return errors.New("tour_id deve ser um número válido")
		}
		r.TourID = &tourID
	}

	if format != "" {
		r.Format = format
	}

	if pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			return errors.New("página deve ser um número válido")
		}
		r.Page = page
	}

	if limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			return errors.New("limite deve ser um número válido")
		}
		r.Limit = limit
	}

	return nil
}
