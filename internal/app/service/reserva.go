package service

import (
	"net/http"
	"time"

	"github.com/jampa_trip/internal/app/contract"
	"github.com/jampa_trip/internal/app/model"
	"github.com/jampa_trip/internal/app/repository"
	"github.com/jampa_trip/internal/pkg/util"
	"gorm.io/gorm"
)

// ReservaService - objeto de contexto
type ReservaService struct {
	ReservaRepository *repository.ReservaRepository
}

// ReservaServiceNew - construtor do objeto
func ReservaServiceNew(DB *gorm.DB) *ReservaService {
	return &ReservaService{
		ReservaRepository: repository.ReservaRepositoryNew(DB),
	}
}

// Create - cria uma nova reserva
func (s *ReservaService) Create(request *contract.CreateReservaRequest) (*contract.CreateReservaResponse, error) {
	// Validar se a data do passeio é futura
	if request.DataPasseio.Before(time.Now()) {
		return nil, util.WrapError("data do passeio deve ser futura", nil, http.StatusBadRequest)
	}

	// Validar se a data de reserva é anterior à data do passeio
	if request.DataReserva.After(request.DataPasseio) {
		return nil, util.WrapError("data de reserva deve ser anterior à data do passeio", nil, http.StatusBadRequest)
	}

	reserva := &model.Reserva{
		ClienteID:          request.ClienteID,
		EmpresaID:          request.EmpresaID,
		PagamentoID:        request.PagamentoID,
		Status:             string(model.StatusReservaPendente),
		DataReserva:        request.DataReserva,
		DataPasseio:        request.DataPasseio,
		QuantidadePessoas:  request.QuantidadePessoas,
		ValorTotal:         request.ValorTotal,
		Observacoes:        request.Observacoes,
		MomentoCriacao:     time.Now(),
		MomentoAtualizacao: time.Now(),
	}

	if err := s.ReservaRepository.Create(reserva); err != nil {
		return nil, util.WrapError("erro ao criar reserva", err, http.StatusInternalServerError)
	}

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

	// Definir valores padrão
	if request.Page <= 0 {
		request.Page = 1
	}
	if request.Limit <= 0 {
		request.Limit = 10
	}

	// Buscar reservas baseado nos filtros
	if request.ClienteID > 0 {
		reservas, total, err = s.ReservaRepository.GetByClienteID(request.ClienteID, request.Page, request.Limit)
	} else if request.EmpresaID > 0 {
		reservas, total, err = s.ReservaRepository.GetByEmpresaID(request.EmpresaID, request.Page, request.Limit)
	} else if request.Status != "" {
		reservas, total, err = s.ReservaRepository.GetByStatus(request.Status, request.Page, request.Limit)
	} else {
		// Buscar todas as reservas (implementar método GetAll se necessário)
		return nil, util.WrapError("filtros de busca não especificados", nil, http.StatusBadRequest)
	}

	if err != nil {
		return nil, util.WrapError("erro ao buscar reservas", err, http.StatusInternalServerError)
	}

	// Mapear para response
	responseReservas := make([]contract.ReservaResponse, len(reservas))
	for i, reserva := range reservas {
		responseReservas[i] = s.mapReservaToResponse(&reserva)
	}

	pages := int((total + int64(request.Limit) - 1) / int64(request.Limit))

	response := &contract.ListReservaResponse{
		Reservas: responseReservas,
		Total:    int(total),
		Page:     request.Page,
		Limit:    request.Limit,
		Pages:    pages,
	}

	return response, nil
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

	// Atualizar campos se fornecidos
	if request.Status != "" {
		reserva.Status = request.Status
	}
	if !request.DataPasseio.IsZero() {
		reserva.DataPasseio = request.DataPasseio
	}
	if request.QuantidadePessoas > 0 {
		reserva.QuantidadePessoas = request.QuantidadePessoas
	}
	if request.ValorTotal > 0 {
		reserva.ValorTotal = request.ValorTotal
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

	// Verificar se a reserva pode ser cancelada
	if !reserva.CanBeCancelled() {
		return nil, util.WrapError("reserva não pode ser cancelada", nil, http.StatusBadRequest)
	}

	if err := s.ReservaRepository.Cancel(request.ID); err != nil {
		return nil, util.WrapError("erro ao cancelar reserva", err, http.StatusInternalServerError)
	}

	// Buscar reserva atualizada
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
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	reservas, total, err := s.ReservaRepository.GetUpcoming(clienteID, page, limit)
	if err != nil {
		return nil, util.WrapError("erro ao buscar reservas futuras", err, http.StatusInternalServerError)
	}

	// Mapear para response
	responseReservas := make([]contract.ReservaResponse, len(reservas))
	for i, reserva := range reservas {
		responseReservas[i] = s.mapReservaToResponse(&reserva)
	}

	pages := int((total + int64(limit) - 1) / int64(limit))

	response := &contract.ListReservaResponse{
		Reservas: responseReservas,
		Total:    int(total),
		Page:     page,
		Limit:    limit,
		Pages:    pages,
	}

	return response, nil
}

// GetHistory - busca histórico de reservas de um cliente
func (s *ReservaService) GetHistory(clienteID int, page, limit int) (*contract.ListReservaResponse, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	reservas, total, err := s.ReservaRepository.GetHistory(clienteID, page, limit)
	if err != nil {
		return nil, util.WrapError("erro ao buscar histórico de reservas", err, http.StatusInternalServerError)
	}

	// Mapear para response
	responseReservas := make([]contract.ReservaResponse, len(reservas))
	for i, reserva := range reservas {
		responseReservas[i] = s.mapReservaToResponse(&reserva)
	}

	pages := int((total + int64(limit) - 1) / int64(limit))

	response := &contract.ListReservaResponse{
		Reservas: responseReservas,
		Total:    int(total),
		Page:     page,
		Limit:    limit,
		Pages:    pages,
	}

	return response, nil
}

// mapReservaToResponse - mapeia model.Reserva para contract.ReservaResponse
func (s *ReservaService) mapReservaToResponse(reserva *model.Reserva) contract.ReservaResponse {
	return contract.ReservaResponse{
		ID:                  reserva.ID,
		ClienteID:           reserva.ClienteID,
		EmpresaID:           reserva.EmpresaID,
		PagamentoID:         reserva.PagamentoID,
		Status:              reserva.Status,
		DataReserva:         reserva.DataReserva,
		DataPasseio:         reserva.DataPasseio,
		QuantidadePessoas:   reserva.QuantidadePessoas,
		ValorTotal:          reserva.ValorTotal,
		Observacoes:         reserva.Observacoes,
		MomentoCriacao:      reserva.MomentoCriacao,
		MomentoAtualizacao:  reserva.MomentoAtualizacao,
		MomentoCancelamento: reserva.MomentoCancelamento,
		StatusDisplay:       reserva.GetStatusDisplay(),
	}
}
