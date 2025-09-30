package repository

import (
	"time"

	"github.com/jampa_trip/internal/app/model"
	"gorm.io/gorm"
)

// FeedbackRepository - objeto de contexto
type FeedbackRepository struct {
	DB *gorm.DB
}

// FeedbackRepositoryNew - construtor do objeto
func FeedbackRepositoryNew(DB *gorm.DB) *FeedbackRepository {
	return &FeedbackRepository{
		DB: DB,
	}
}

// Create - cria um novo feedback
func (r *FeedbackRepository) Create(feedback *model.Feedback) error {
	return r.DB.Create(feedback).Error
}

// GetByID - busca um feedback pelo ID
func (r *FeedbackRepository) GetByID(id int) (*model.Feedback, error) {
	var feedback model.Feedback
	err := r.DB.Preload("Cliente").Preload("Empresa").Preload("Reserva").First(&feedback, id).Error
	if err != nil {
		return nil, err
	}
	return &feedback, nil
}

// GetByClienteID - busca feedbacks por cliente
func (r *FeedbackRepository) GetByClienteID(clienteID int, page, limit int) ([]model.Feedback, int64, error) {
	var feedbacks []model.Feedback
	var total int64

	offset := (page - 1) * limit

	query := r.DB.Model(&model.Feedback{}).Where("cliente_id = ?", clienteID)

	// Contar total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Buscar registros
	err := query.Preload("Cliente").Preload("Empresa").Preload("Reserva").
		Offset(offset).Limit(limit).
		Order("momento_criacao DESC").
		Find(&feedbacks).Error

	return feedbacks, total, err
}

// GetByEmpresaID - busca feedbacks por empresa
func (r *FeedbackRepository) GetByEmpresaID(empresaID int, page, limit int) ([]model.Feedback, int64, error) {
	var feedbacks []model.Feedback
	var total int64

	offset := (page - 1) * limit

	query := r.DB.Model(&model.Feedback{}).Where("empresa_id = ?", empresaID)

	// Contar total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Buscar registros
	err := query.Preload("Cliente").Preload("Empresa").Preload("Reserva").
		Offset(offset).Limit(limit).
		Order("momento_criacao DESC").
		Find(&feedbacks).Error

	return feedbacks, total, err
}

// GetByStatus - busca feedbacks por status
func (r *FeedbackRepository) GetByStatus(status string, page, limit int) ([]model.Feedback, int64, error) {
	var feedbacks []model.Feedback
	var total int64

	offset := (page - 1) * limit

	query := r.DB.Model(&model.Feedback{}).Where("status = ?", status)

	// Contar total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Buscar registros
	err := query.Preload("Cliente").Preload("Empresa").Preload("Reserva").
		Offset(offset).Limit(limit).
		Order("momento_criacao DESC").
		Find(&feedbacks).Error

	return feedbacks, total, err
}

// GetByRating - busca feedbacks por nota
func (r *FeedbackRepository) GetByRating(rating int, page, limit int) ([]model.Feedback, int64, error) {
	var feedbacks []model.Feedback
	var total int64

	offset := (page - 1) * limit

	query := r.DB.Model(&model.Feedback{}).Where("nota = ?", rating)

	// Contar total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Buscar registros
	err := query.Preload("Cliente").Preload("Empresa").Preload("Reserva").
		Offset(offset).Limit(limit).
		Order("momento_criacao DESC").
		Find(&feedbacks).Error

	return feedbacks, total, err
}

// Update - atualiza um feedback
func (r *FeedbackRepository) Update(feedback *model.Feedback) error {
	return r.DB.Save(feedback).Error
}

// UpdateStatus - atualiza apenas o status de um feedback
func (r *FeedbackRepository) UpdateStatus(id int, status string) error {
	return r.DB.Model(&model.Feedback{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":              status,
			"momento_atualizacao": time.Now(),
		}).Error
}

// Delete - remove um feedback
func (r *FeedbackRepository) Delete(id int) error {
	return r.DB.Delete(&model.Feedback{}, id).Error
}

// GetAverageRating - calcula a média de avaliações de uma empresa
func (r *FeedbackRepository) GetAverageRating(empresaID int) (float64, int, error) {
	var result struct {
		Average float64
		Count   int
	}

	err := r.DB.Model(&model.Feedback{}).
		Select("AVG(nota) as average, COUNT(*) as count").
		Where("empresa_id = ? AND status = ?", empresaID, string(model.StatusFeedbackAtivo)).
		Scan(&result).Error

	return result.Average, result.Count, err
}

// GetRatingDistribution - obtém a distribuição de notas de uma empresa
func (r *FeedbackRepository) GetRatingDistribution(empresaID int) (map[int]int, error) {
	var results []struct {
		Nota  int
		Count int
	}

	err := r.DB.Model(&model.Feedback{}).
		Select("nota, COUNT(*) as count").
		Where("empresa_id = ? AND status = ?", empresaID, string(model.StatusFeedbackAtivo)).
		Group("nota").
		Order("nota").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	distribution := make(map[int]int)
	for _, result := range results {
		distribution[result.Nota] = result.Count
	}

	return distribution, nil
}

// GetRecentFeedbacks - busca feedbacks recentes
func (r *FeedbackRepository) GetRecentFeedbacks(empresaID int, days int, page, limit int) ([]model.Feedback, int64, error) {
	var feedbacks []model.Feedback
	var total int64

	offset := (page - 1) * limit
	since := time.Now().AddDate(0, 0, -days)

	query := r.DB.Model(&model.Feedback{}).Where("empresa_id = ? AND momento_criacao >= ?", empresaID, since)

	// Contar total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Buscar registros
	err := query.Preload("Cliente").Preload("Empresa").Preload("Reserva").
		Offset(offset).Limit(limit).
		Order("momento_criacao DESC").
		Find(&feedbacks).Error

	return feedbacks, total, err
}
