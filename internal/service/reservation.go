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

// ReservaService - objeto de contexto
type ReservaService struct {
	ReservaRepository *repository.ReservaRepository
	TourRepository    *repository.TourRepository
	ClientRepository  *repository.ClientRepository
}

// ReservaServiceNew - construtor do objeto
func ReservaServiceNew(DB *gorm.DB) *ReservaService {
	return &ReservaService{
		ReservaRepository: repository.ReservaRepositoryNew(DB),
		TourRepository:    repository.TourRepositoryNew(DB),
		ClientRepository:  repository.ClientRepositoryNew(DB),
	}
}

// Create - cria uma nova reserva
func (s *ReservaService) Create(request *contract.CreateReservaRequest) (*contract.CreateReservaResponse, error) {
	cliente, err := s.ClientRepository.GetByID(request.ClienteID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, util.WrapError("cliente não encontrado", err, http.StatusNotFound)
		}
		return nil, util.WrapError("erro ao buscar cliente", err, http.StatusInternalServerError)
	}

	tour, err := s.TourRepository.GetByID(request.TourID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, util.WrapError("tour não encontrado", err, http.StatusNotFound)
		}
		return nil, util.WrapError("erro ao buscar tour", err, http.StatusInternalServerError)
	}

	if err := s.validateDataPasseio(request.DataPasseioSelecionada, tour); err != nil {
		return nil, err
	}

	valorTotal := tour.Price * float64(request.QuantidadePessoas)

	reserva := &model.Reserva{
		TourID:                 request.TourID,
		ClienteID:              request.ClienteID,
		PagamentoID:            request.PagamentoID,
		Status:                 string(model.StatusReservaPendente),
		DataReserva:            time.Now(),
		DataPasseioSelecionada: request.DataPasseioSelecionada,
		QuantidadePessoas:      request.QuantidadePessoas,
		ValorTotal:             valorTotal,
		Observacoes:            request.Observacoes,
		MomentoCriacao:         time.Now(),
		MomentoAtualizacao:     time.Now(),
	}

	if err := s.ReservaRepository.Create(reserva); err != nil {
		return nil, util.WrapError("erro ao criar reserva", err, http.StatusInternalServerError)
	}

	reserva.Tour = *tour
	reserva.Cliente = *cliente

	response := &contract.CreateReservaResponse{
		Reserva: s.mapReservaToResponse(reserva),
		Message: "Reserva criada com sucesso",
	}

	return response, nil
}

// GetByID - busca uma reserva pelo ID
func (s *ReservaService) GetByID(request *contract.GetReservaRequest) (*contract.ReservaResponse, error) {
	reserva, err := s.ReservaRepository.GetByID(request.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, util.WrapError("reserva não encontrada", err, http.StatusNotFound)
		}
		return nil, util.WrapError("erro ao buscar reserva", err, http.StatusInternalServerError)
	}

	response := s.mapReservaToResponse(reserva)
	return &response, nil
}

// List - lista reservas
func (s *ReservaService) List(request *contract.ListReservaRequest) (*contract.ListReservaResponse, error) {
	var reservas []model.Reserva
	var total int64
	var err error

	if request.Page <= 0 {
		request.Page = 1
	}
	if request.Limit <= 0 {
		request.Limit = 10
	}

	if request.TourID > 0 {
		reservas, total, err = s.ReservaRepository.GetByTourID(request.TourID, request.Page, request.Limit)
	} else if request.ClienteID > 0 {
		reservas, total, err = s.ReservaRepository.GetByClienteID(request.ClienteID, request.Page, request.Limit)
	} else if request.CompanyID > 0 {
		reservas, total, err = s.ReservaRepository.GetByCompanyID(request.CompanyID, request.Page, request.Limit)
	} else if request.Status != "" {
		reservas, total, err = s.ReservaRepository.GetByStatus(request.Status, request.Page, request.Limit)
	} else {
		return nil, util.WrapError("filtros de busca não especificados", nil, http.StatusBadRequest)
	}

	if err != nil {
		return nil, util.WrapError("erro ao buscar reservas", err, http.StatusInternalServerError)
	}

	return s.buildListResponse(reservas, total, request.Page, request.Limit), nil
}

// Update - atualiza uma reserva
func (s *ReservaService) Update(id int, request *contract.UpdateReservaRequest) (*contract.UpdateReservaResponse, error) {
	reserva, err := s.ReservaRepository.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, util.WrapError("reserva não encontrada", err, http.StatusNotFound)
		}
		return nil, util.WrapError("erro ao buscar reserva", err, http.StatusInternalServerError)
	}

	if !request.DataPasseioSelecionada.IsZero() {
		tour, err := s.TourRepository.GetByID(reserva.TourID)
		if err != nil {
			return nil, util.WrapError("erro ao buscar tour da reserva", err, http.StatusInternalServerError)
		}

		if err := s.validateDataPasseio(request.DataPasseioSelecionada, tour); err != nil {
			return nil, err
		}

		reserva.DataPasseioSelecionada = request.DataPasseioSelecionada
	}
	if request.QuantidadePessoas > 0 {
		tour, err := s.TourRepository.GetByID(reserva.TourID)
		if err != nil {
			return nil, util.WrapError("erro ao buscar tour da reserva", err, http.StatusInternalServerError)
		}

		reserva.QuantidadePessoas = request.QuantidadePessoas
		reserva.ValorTotal = tour.Price * float64(request.QuantidadePessoas)
	}

	if request.Status != "" {
		reserva.Status = request.Status
	}
	if request.Observacoes != "" {
		reserva.Observacoes = request.Observacoes
	}

	reserva.MomentoAtualizacao = time.Now()

	if err := s.ReservaRepository.Update(reserva); err != nil {
		return nil, util.WrapError("erro ao atualizar reserva", err, http.StatusInternalServerError)
	}

	response := &contract.UpdateReservaResponse{
		Reserva: s.mapReservaToResponse(reserva),
		Message: "Reserva atualizada com sucesso",
	}

	return response, nil
}

// Cancel - cancela uma reserva
func (s *ReservaService) Cancel(request *contract.CancelarReservaRequest) (*contract.CancelarReservaResponse, error) {
	reserva, err := s.ReservaRepository.GetByID(request.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, util.WrapError("reserva não encontrada", err, http.StatusNotFound)
		}
		return nil, util.WrapError("erro ao buscar reserva", err, http.StatusInternalServerError)
	}

	if !reserva.CanBeCancelled() {
		return nil, util.WrapError("reserva não pode ser cancelada", nil, http.StatusBadRequest)
	}

	if err := s.ReservaRepository.Cancel(request.ID); err != nil {
		return nil, util.WrapError("erro ao cancelar reserva", err, http.StatusInternalServerError)
	}

	reserva, err = s.ReservaRepository.GetByID(request.ID)
	if err != nil {
		return nil, util.WrapError("erro ao buscar reserva atualizada", err, http.StatusInternalServerError)
	}

	response := &contract.CancelarReservaResponse{
		Reserva: s.mapReservaToResponse(reserva),
		Message: "Reserva cancelada com sucesso",
	}

	return response, nil
}

// GetUpcoming - busca reservas futuras de um cliente
func (s *ReservaService) GetUpcoming(clienteID int, page, limit int) (*contract.ListReservaResponse, error) {
	return s.getReservasByClienteAndDate(clienteID, page, limit, true)
}

// GetHistory - busca histórico de reservas de um cliente
func (s *ReservaService) GetHistory(clienteID int, page, limit int) (*contract.ListReservaResponse, error) {
	return s.getReservasByClienteAndDate(clienteID, page, limit, false)
}

// ============================================================================
// Métodos Auxiliares Privados
// ============================================================================

// validateDataPasseio - valida se a data do passeio é futura e está disponível no tour
func (s *ReservaService) validateDataPasseio(dataPasseio time.Time, tour *model.Tour) error {
	if dataPasseio.Before(time.Now()) {
		return util.WrapError("data do passeio deve ser futura", nil, http.StatusBadRequest)
	}

	dataEncontrada := false
	dataSelecionadaStr := dataPasseio.Format("2006-01-02")
	for _, data := range tour.Dates {
		if data == dataSelecionadaStr {
			dataEncontrada = true
			break
		}
	}
	if !dataEncontrada {
		return util.WrapError("data selecionada não está disponível para este tour", nil, http.StatusBadRequest)
	}

	return nil
}

// getReservasByClienteAndDate - método auxiliar para GetUpcoming e GetHistory
func (s *ReservaService) getReservasByClienteAndDate(clienteID int, page, limit int, upcoming bool) (*contract.ListReservaResponse, error) {
	config := util.NormalizePagination(page, limit)
	page = config.Page
	limit = config.Limit

	var reservas []model.Reserva
	var total int64
	var err error

	if upcoming {
		reservas, total, err = s.ReservaRepository.GetUpcoming(clienteID, page, limit)
	} else {
		reservas, total, err = s.ReservaRepository.GetHistory(clienteID, page, limit)
	}

	if err != nil {
		tipoReserva := "histórico de reservas"
		if upcoming {
			tipoReserva = "reservas futuras"
		}
		return nil, util.WrapError("erro ao buscar "+tipoReserva, err, http.StatusInternalServerError)
	}

	return s.buildListResponse(reservas, total, page, limit), nil
}

// buildListResponse - constrói resposta de listagem paginada
func (s *ReservaService) buildListResponse(reservas []model.Reserva, total int64, page, limit int) *contract.ListReservaResponse {
	responseReservas := make([]contract.ReservaResponse, len(reservas))
	for i, reserva := range reservas {
		responseReservas[i] = s.mapReservaToResponse(&reserva)
	}

	pages := util.CalculateTotalPages(total, limit)

	return &contract.ListReservaResponse{
		Reservas: responseReservas,
		Total:    int(total),
		Page:     page,
		Limit:    limit,
		Pages:    pages,
	}
}

// mapReservaToResponse - mapeia model.Reserva para contract.ReservaResponse
func (s *ReservaService) mapReservaToResponse(reserva *model.Reserva) contract.ReservaResponse {
	response := contract.ReservaResponse{
		ID:                     reserva.ID,
		TourID:                 reserva.TourID,
		ClienteID:              reserva.ClienteID,
		PagamentoID:            reserva.PagamentoID,
		Status:                 reserva.Status,
		DataReserva:            reserva.DataReserva,
		DataPasseioSelecionada: reserva.DataPasseioSelecionada,
		QuantidadePessoas:      reserva.QuantidadePessoas,
		ValorTotal:             reserva.ValorTotal,
		Observacoes:            reserva.Observacoes,
		MomentoCriacao:         reserva.MomentoCriacao,
		MomentoAtualizacao:     reserva.MomentoAtualizacao,
		MomentoCancelamento:    reserva.MomentoCancelamento,
		StatusDisplay:          reserva.GetStatusDisplay(),
	}

	// Popular dados relacionados se disponíveis (preloaded)
	if reserva.Tour.ID != 0 {
		response.TourName = reserva.Tour.Name
		response.TourPrice = reserva.Tour.Price
		response.CompanyID = reserva.Tour.CompanyID
	}

	if reserva.Cliente.ID != 0 {
		response.ClienteName = reserva.Cliente.Name
		response.ClienteEmail = reserva.Cliente.Email
	}

	if reserva.Pagamento.ID != 0 {
		response.PagamentoStatus = reserva.Pagamento.GetStatusDisplay()
	}

	return response
}
