package service

import (
	"net/http"
	"time"

	"github.com/jampa_trip/internal/contract"
	"github.com/jampa_trip/internal/model"
	"github.com/jampa_trip/internal/repository"
	"github.com/jampa_trip/pkg/util"
	"gorm.io/gorm"
)

// FeedbackService - objeto de contexto
type FeedbackService struct {
	FeedbackRepository *repository.FeedbackRepository
}

// FeedbackServiceNew - construtor do objeto
func FeedbackServiceNew(DB *gorm.DB) *FeedbackService {
	return &FeedbackService{
		FeedbackRepository: repository.FeedbackRepositoryNew(DB),
	}
}

// Create - cria um novo feedback
func (s *FeedbackService) Create(request *contract.CreateFeedbackRequest) (*contract.CreateFeedbackResponse, error) {
	feedback := &model.Feedback{
		ClienteID:          request.ClienteID,
		EmpresaID:          request.EmpresaID,
		ReservaID:          request.ReservaID,
		Nota:               request.Nota,
		Comentario:         request.Comentario,
		Status:             string(model.StatusFeedbackAtivo),
		MomentoCriacao:     time.Now(),
		MomentoAtualizacao: time.Now(),
	}

	if err := s.FeedbackRepository.Create(feedback); err != nil {
		return nil, util.WrapError("erro ao criar feedback", err, http.StatusInternalServerError)
	}

	response := &contract.CreateFeedbackResponse{
		Feedback: s.mapFeedbackToResponse(feedback),
		Message:  "Feedback criado com sucesso",
	}

	return response, nil
}

// GetByID - busca um feedback pelo ID
func (s *FeedbackService) GetByID(request *contract.GetFeedbackRequest) (*contract.FeedbackResponse, error) {
	feedback, err := s.FeedbackRepository.GetByID(request.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, util.WrapError("feedback não encontrado", err, http.StatusNotFound)
		}
		return nil, util.WrapError("erro ao buscar feedback", err, http.StatusInternalServerError)
	}

	response := s.mapFeedbackToResponse(feedback)
	return &response, nil
}

// List - lista feedbacks
func (s *FeedbackService) List(request *contract.ListFeedbackRequest) (*contract.ListFeedbackResponse, error) {
	var feedbacks []model.Feedback
	var total int64
	var err error

	// Definir valores padrão
	if request.Page <= 0 {
		request.Page = 1
	}
	if request.Limit <= 0 {
		request.Limit = 10
	}

	// Buscar feedbacks baseado nos filtros
	if request.ClienteID > 0 {
		feedbacks, total, err = s.FeedbackRepository.GetByClienteID(request.ClienteID, request.Page, request.Limit)
	} else if request.EmpresaID > 0 {
		feedbacks, total, err = s.FeedbackRepository.GetByEmpresaID(request.EmpresaID, request.Page, request.Limit)
	} else if request.Status != "" {
		feedbacks, total, err = s.FeedbackRepository.GetByStatus(request.Status, request.Page, request.Limit)
	} else if request.Nota > 0 {
		feedbacks, total, err = s.FeedbackRepository.GetByRating(request.Nota, request.Page, request.Limit)
	} else {
		// Buscar todos os feedbacks (implementar método GetAll se necessário)
		return nil, util.WrapError("filtros de busca não especificados", nil, http.StatusBadRequest)
	}

	if err != nil {
		return nil, util.WrapError("erro ao buscar feedbacks", err, http.StatusInternalServerError)
	}

	// Mapear para response
	responseFeedbacks := make([]contract.FeedbackResponse, len(feedbacks))
	for i, feedback := range feedbacks {
		responseFeedbacks[i] = s.mapFeedbackToResponse(&feedback)
	}

	pages := util.CalculateTotalPages(total, request.Limit)

	response := &contract.ListFeedbackResponse{
		Feedbacks: responseFeedbacks,
		Total:     int(total),
		Page:      request.Page,
		Limit:     request.Limit,
		Pages:     pages,
	}

	return response, nil
}

// Update - atualiza um feedback
func (s *FeedbackService) Update(id int, request *contract.UpdateFeedbackRequest) (*contract.UpdateFeedbackResponse, error) {
	feedback, err := s.FeedbackRepository.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, util.WrapError("feedback não encontrado", err, http.StatusNotFound)
		}
		return nil, util.WrapError("erro ao buscar feedback", err, http.StatusInternalServerError)
	}

	// Atualizar campos se fornecidos
	if request.Nota > 0 {
		feedback.Nota = request.Nota
	}
	if request.Comentario != "" {
		feedback.Comentario = request.Comentario
	}
	if request.Status != "" {
		feedback.Status = request.Status
	}

	feedback.MomentoAtualizacao = time.Now()

	if err := s.FeedbackRepository.Update(feedback); err != nil {
		return nil, util.WrapError("erro ao atualizar feedback", err, http.StatusInternalServerError)
	}

	response := &contract.UpdateFeedbackResponse{
		Feedback: s.mapFeedbackToResponse(feedback),
		Message:  "Feedback atualizado com sucesso",
	}

	return response, nil
}

// GetAverageRating - obtém a média de avaliações de uma empresa
func (s *FeedbackService) GetAverageRating(empresaID int) (float64, int, error) {
	average, count, err := s.FeedbackRepository.GetAverageRating(empresaID)
	if err != nil {
		return 0, 0, util.WrapError("erro ao calcular média de avaliações", err, http.StatusInternalServerError)
	}

	return average, count, nil
}

// GetRatingDistribution - obtém a distribuição de notas de uma empresa
func (s *FeedbackService) GetRatingDistribution(empresaID int) (map[int]int, error) {
	distribution, err := s.FeedbackRepository.GetRatingDistribution(empresaID)
	if err != nil {
		return nil, util.WrapError("erro ao obter distribuição de notas", err, http.StatusInternalServerError)
	}

	return distribution, nil
}

// GetRecentFeedbacks - busca feedbacks recentes de uma empresa
func (s *FeedbackService) GetRecentFeedbacks(empresaID int, days int, page, limit int) (*contract.ListFeedbackResponse, error) {
	config := util.NormalizePagination(page, limit)
	page = config.Page
	limit = config.Limit

	feedbacks, total, err := s.FeedbackRepository.GetRecentFeedbacks(empresaID, days, page, limit)
	if err != nil {
		return nil, util.WrapError("erro ao buscar feedbacks recentes", err, http.StatusInternalServerError)
	}

	// Mapear para response
	responseFeedbacks := make([]contract.FeedbackResponse, len(feedbacks))
	for i, feedback := range feedbacks {
		responseFeedbacks[i] = s.mapFeedbackToResponse(&feedback)
	}

	pages := util.CalculateTotalPages(total, limit)

	response := &contract.ListFeedbackResponse{
		Feedbacks: responseFeedbacks,
		Total:     int(total),
		Page:      page,
		Limit:     limit,
		Pages:     pages,
	}

	return response, nil
}

// mapFeedbackToResponse - mapeia model.Feedback para contract.FeedbackResponse
func (s *FeedbackService) mapFeedbackToResponse(feedback *model.Feedback) contract.FeedbackResponse {
	return contract.FeedbackResponse{
		ID:                 feedback.ID,
		ClienteID:          feedback.ClienteID,
		EmpresaID:          feedback.EmpresaID,
		ReservaID:          feedback.ReservaID,
		Nota:               feedback.Nota,
		Comentario:         feedback.Comentario,
		Status:             feedback.Status,
		MomentoCriacao:     feedback.MomentoCriacao,
		MomentoAtualizacao: feedback.MomentoAtualizacao,
		StatusDisplay:      feedback.GetStatusDisplay(),
		RatingDisplay:      feedback.GetRatingDisplay(),
	}
}
