package repository

import (
	"time"

	"github.com/jampa_trip/internal/model"
	"github.com/jampa_trip/internal/query"
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
	err := r.DB.Raw(query.CreateFeedback,
		feedback.ClienteID,
		feedback.EmpresaID,
		feedback.ReservaID,
		feedback.Nota,
		feedback.Comentario,
		feedback.Status,
		feedback.MomentoCriacao,
		feedback.MomentoAtualizacao,
	).Row().Scan(&feedback.ID)

	return err
}

// GetByID - busca um feedback pelo ID
func (r *FeedbackRepository) GetByID(id int) (*model.Feedback, error) {
	feedback := &model.Feedback{}

	err := r.DB.Raw(query.GetFeedbackByID, id).Row().Scan(
		&feedback.ID,
		&feedback.ClienteID,
		&feedback.EmpresaID,
		&feedback.ReservaID,
		&feedback.Nota,
		&feedback.Comentario,
		&feedback.Status,
		&feedback.MomentoCriacao,
		&feedback.MomentoAtualizacao,
	)

	if err != nil {
		return nil, err
	}

	return feedback, nil
}

// GetByClienteID - busca feedbacks por cliente
func (r *FeedbackRepository) GetByClienteID(clienteID int, page, limit int) ([]model.Feedback, int64, error) {
	offset := (page - 1) * limit

	// Buscar registros
	rows, err := r.DB.Raw(query.GetFeedbacksByClienteID, clienteID, limit, offset).Rows()
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var feedbacks []model.Feedback
	for rows.Next() {
		feedback := model.Feedback{}
		err := rows.Scan(
			&feedback.ID,
			&feedback.ClienteID,
			&feedback.EmpresaID,
			&feedback.ReservaID,
			&feedback.Nota,
			&feedback.Comentario,
			&feedback.Status,
			&feedback.MomentoCriacao,
			&feedback.MomentoAtualizacao,
		)
		if err != nil {
			return nil, 0, err
		}
		feedbacks = append(feedbacks, feedback)
	}

	// Contar total
	var total int64
	err = r.DB.Raw(query.CountFeedbacksByClienteID, clienteID).Row().Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return feedbacks, total, nil
}

// GetByEmpresaID - busca feedbacks por empresa
func (r *FeedbackRepository) GetByEmpresaID(empresaID int, page, limit int) ([]model.Feedback, int64, error) {
	offset := (page - 1) * limit

	// Buscar registros
	rows, err := r.DB.Raw(query.GetFeedbacksByEmpresaID, empresaID, limit, offset).Rows()
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var feedbacks []model.Feedback
	for rows.Next() {
		feedback := model.Feedback{}
		err := rows.Scan(
			&feedback.ID,
			&feedback.ClienteID,
			&feedback.EmpresaID,
			&feedback.ReservaID,
			&feedback.Nota,
			&feedback.Comentario,
			&feedback.Status,
			&feedback.MomentoCriacao,
			&feedback.MomentoAtualizacao,
		)
		if err != nil {
			return nil, 0, err
		}
		feedbacks = append(feedbacks, feedback)
	}

	// Contar total
	var total int64
	err = r.DB.Raw(query.CountFeedbacksByEmpresaID, empresaID).Row().Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return feedbacks, total, nil
}

// GetByStatus - busca feedbacks por status
func (r *FeedbackRepository) GetByStatus(status string, page, limit int) ([]model.Feedback, int64, error) {
	offset := (page - 1) * limit

	// Buscar registros
	rows, err := r.DB.Raw(query.GetFeedbacksByStatus, status, limit, offset).Rows()
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var feedbacks []model.Feedback
	for rows.Next() {
		feedback := model.Feedback{}
		err := rows.Scan(
			&feedback.ID,
			&feedback.ClienteID,
			&feedback.EmpresaID,
			&feedback.ReservaID,
			&feedback.Nota,
			&feedback.Comentario,
			&feedback.Status,
			&feedback.MomentoCriacao,
			&feedback.MomentoAtualizacao,
		)
		if err != nil {
			return nil, 0, err
		}
		feedbacks = append(feedbacks, feedback)
	}

	// Contar total
	var total int64
	err = r.DB.Raw(query.CountFeedbacksByStatus, status).Row().Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return feedbacks, total, nil
}

// GetByRating - busca feedbacks por nota
func (r *FeedbackRepository) GetByRating(rating int, page, limit int) ([]model.Feedback, int64, error) {
	offset := (page - 1) * limit

	// Buscar registros
	rows, err := r.DB.Raw(query.GetFeedbacksByRating, rating, limit, offset).Rows()
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var feedbacks []model.Feedback
	for rows.Next() {
		feedback := model.Feedback{}
		err := rows.Scan(
			&feedback.ID,
			&feedback.ClienteID,
			&feedback.EmpresaID,
			&feedback.ReservaID,
			&feedback.Nota,
			&feedback.Comentario,
			&feedback.Status,
			&feedback.MomentoCriacao,
			&feedback.MomentoAtualizacao,
		)
		if err != nil {
			return nil, 0, err
		}
		feedbacks = append(feedbacks, feedback)
	}

	// Contar total
	var total int64
	err = r.DB.Raw(query.CountFeedbacksByRating, rating).Row().Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return feedbacks, total, nil
}

// Update - atualiza um feedback
func (r *FeedbackRepository) Update(feedback *model.Feedback) error {
	err := r.DB.Raw(query.UpdateFeedback,
		feedback.Nota,
		feedback.Comentario,
		feedback.Status,
		feedback.MomentoAtualizacao,
		feedback.ID,
	).Row().Scan()

	return err
}

// UpdateStatus - atualiza apenas o status de um feedback
func (r *FeedbackRepository) UpdateStatus(id int, status string) error {
	err := r.DB.Raw(query.UpdateFeedbackStatus,
		status,
		time.Now(),
		id,
	).Row().Scan()

	return err
}

// Delete - remove um feedback
func (r *FeedbackRepository) Delete(id int) error {
	err := r.DB.Raw(query.DeleteFeedback, id).Row().Scan()
	return err
}

// GetAverageRating - calcula a média de avaliações de uma empresa
func (r *FeedbackRepository) GetAverageRating(empresaID int) (float64, int, error) {
	var average float64
	var count int

	err := r.DB.Raw(query.GetAverageRating, empresaID, string(model.StatusFeedbackAtivo)).Row().Scan(&average, &count)
	if err != nil {
		return 0, 0, err
	}

	return average, count, nil
}

// GetRatingDistribution - obtém a distribuição de notas de uma empresa
func (r *FeedbackRepository) GetRatingDistribution(empresaID int) (map[int]int, error) {
	rows, err := r.DB.Raw(query.GetRatingDistribution, empresaID, string(model.StatusFeedbackAtivo)).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	distribution := make(map[int]int)
	for rows.Next() {
		var nota int
		var count int
		err := rows.Scan(&nota, &count)
		if err != nil {
			return nil, err
		}
		distribution[nota] = count
	}

	return distribution, nil
}

// GetRecentFeedbacks - busca feedbacks recentes
func (r *FeedbackRepository) GetRecentFeedbacks(empresaID int, days int, page, limit int) ([]model.Feedback, int64, error) {
	offset := (page - 1) * limit
	since := time.Now().AddDate(0, 0, -days)

	// Buscar registros
	rows, err := r.DB.Raw(query.GetRecentFeedbacks, empresaID, since, limit, offset).Rows()
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var feedbacks []model.Feedback
	for rows.Next() {
		feedback := model.Feedback{}
		err := rows.Scan(
			&feedback.ID,
			&feedback.ClienteID,
			&feedback.EmpresaID,
			&feedback.ReservaID,
			&feedback.Nota,
			&feedback.Comentario,
			&feedback.Status,
			&feedback.MomentoCriacao,
			&feedback.MomentoAtualizacao,
		)
		if err != nil {
			return nil, 0, err
		}
		feedbacks = append(feedbacks, feedback)
	}

	// Contar total
	var total int64
	err = r.DB.Raw(query.CountRecentFeedbacks, empresaID, since).Row().Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return feedbacks, total, nil
}
