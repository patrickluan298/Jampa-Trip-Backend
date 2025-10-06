package service

import (
	"net/http"
	"time"

	"github.com/jampa_trip/internal/contract"
	"github.com/jampa_trip/internal/model"
	"github.com/jampa_trip/internal/repository"
	"github.com/jampa_trip/pkg/util"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// TourService - objeto de contexto
type TourService struct {
	TourRepository *repository.TourRepository
}

// TourServiceNew - construtor do objeto
func TourServiceNew(DB *gorm.DB) *TourService {
	return &TourService{
		TourRepository: repository.TourRepositoryNew(DB),
	}
}

// Create - cria um novo passeio
func (s *TourService) Create(request *contract.CreateTourRequest, companyID int) (*contract.CreateTourResponse, error) {

	dates := pq.StringArray(request.Dates)
	images := pq.StringArray(request.Images)

	tour := &model.Tour{
		CompanyID:     companyID,
		Name:          request.Name,
		Dates:         dates,
		DepartureTime: request.DepartureTime,
		ArrivalTime:   request.ArrivalTime,
		MaxPeople:     request.MaxPeople,
		Description:   request.Description,
		Images:        images,
		Price:         request.Price,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.TourRepository.Create(tour); err != nil {
		return nil, util.WrapError("Erro ao criar passeio", err, http.StatusInternalServerError)
	}

	tourWithCompany, companyName, err := s.TourRepository.GetTourWithCompanyName(tour.ID)
	if err != nil {
		return nil, util.WrapError("Erro ao buscar passeio criado", err, http.StatusInternalServerError)
	}

	response := &contract.CreateTourResponse{
		Success: true,
		Tour: contract.TourResponse{
			ID:            tourWithCompany.ID,
			Name:          tourWithCompany.Name,
			Dates:         tourWithCompany.Dates,
			DepartureTime: tourWithCompany.DepartureTime,
			ArrivalTime:   tourWithCompany.ArrivalTime,
			MaxPeople:     tourWithCompany.MaxPeople,
			Description:   tourWithCompany.Description,
			Images:        tourWithCompany.Images,
			Price:         tourWithCompany.Price,
			CompanyID:     tourWithCompany.CompanyID,
			CompanyName:   companyName,
			CreatedAt:     tourWithCompany.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:     tourWithCompany.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	}

	return response, nil
}

// Update - atualiza um passeio existente
func (s *TourService) Update(request *contract.UpdateTourRequest, companyID int) (*contract.UpdateTourResponse, error) {

	exists, err := s.TourRepository.IsOwnedByCompany(request.ID, companyID)
	if err != nil {
		return nil, util.WrapError("Erro ao verificar propriedade do passeio", err, http.StatusInternalServerError)
	}
	if !exists {
		return nil, util.WrapError("Passeio não encontrado ou você não tem permissão para editá-lo", nil, http.StatusForbidden)
	}

	dates := pq.StringArray(request.Dates)
	images := pq.StringArray(request.Images)

	tour := &model.Tour{
		ID:            request.ID,
		CompanyID:     companyID,
		Name:          request.Name,
		Dates:         dates,
		DepartureTime: request.DepartureTime,
		ArrivalTime:   request.ArrivalTime,
		MaxPeople:     request.MaxPeople,
		Description:   request.Description,
		Images:        images,
		Price:         request.Price,
		UpdatedAt:     time.Now(),
	}

	if err := s.TourRepository.Update(tour); err != nil {
		return nil, util.WrapError("Erro ao atualizar passeio", err, http.StatusInternalServerError)
	}

	tourWithCompany, companyName, err := s.TourRepository.GetTourWithCompanyName(tour.ID)
	if err != nil {
		return nil, util.WrapError("Erro ao buscar passeio atualizado", err, http.StatusInternalServerError)
	}

	response := &contract.UpdateTourResponse{
		Success: true,
		Tour: contract.TourResponse{
			ID:            tourWithCompany.ID,
			Name:          tourWithCompany.Name,
			Dates:         tourWithCompany.Dates,
			DepartureTime: tourWithCompany.DepartureTime,
			ArrivalTime:   tourWithCompany.ArrivalTime,
			MaxPeople:     tourWithCompany.MaxPeople,
			Description:   tourWithCompany.Description,
			Images:        tourWithCompany.Images,
			Price:         tourWithCompany.Price,
			CompanyID:     tourWithCompany.CompanyID,
			CompanyName:   companyName,
			CreatedAt:     tourWithCompany.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:     tourWithCompany.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	}

	return response, nil
}

// List - lista todos os passeios com filtro e paginação
func (s *TourService) List(request *contract.ListToursRequest) (*contract.ListToursResponse, error) {

	config := util.NormalizePagination(request.Page, request.Limit)
	page := config.Page
	limit := config.Limit

	tours, total, err := s.TourRepository.List(request.Search, page, limit)
	if err != nil {
		return nil, util.WrapError("Erro ao buscar passeios", err, http.StatusInternalServerError)
	}

	var toursResponse []contract.TourResponse
	for _, tour := range tours {
		_, companyName, err := s.TourRepository.GetTourWithCompanyName(tour.ID)
		if err != nil {
			companyName = "Empresa não encontrada"
		}

		toursResponse = append(toursResponse, contract.TourResponse{
			ID:            tour.ID,
			Name:          tour.Name,
			Dates:         tour.Dates,
			DepartureTime: tour.DepartureTime,
			ArrivalTime:   tour.ArrivalTime,
			MaxPeople:     tour.MaxPeople,
			Description:   tour.Description,
			Images:        tour.Images,
			Price:         tour.Price,
			CompanyID:     tour.CompanyID,
			CompanyName:   companyName,
			CreatedAt:     tour.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:     tour.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	totalPages := util.CalculateTotalPages(total, limit)

	response := &contract.ListToursResponse{
		Success: true,
		Tours:   toursResponse,
		Pagination: contract.PaginationResponse{
			CurrentPage:  page,
			TotalPages:   totalPages,
			TotalItems:   int(total),
			ItemsPerPage: limit,
		},
	}

	return response, nil
}

// GetMyTours - lista passeios da empresa
func (s *TourService) GetMyTours(companyID int, page, limit int) (*contract.GetMyToursResponse, error) {

	config := util.NormalizePagination(page, limit)
	page = config.Page
	limit = config.Limit

	tours, total, err := s.TourRepository.ListByCompanyID(companyID, page, limit)
	if err != nil {
		return nil, util.WrapError("Erro ao buscar passeios da empresa", err, http.StatusInternalServerError)
	}

	var toursResponse []contract.MyTourResponse
	for _, tour := range tours {
		reservationsCount, err := s.TourRepository.CountReservationsByTourID(tour.ID)
		if err != nil {
			reservationsCount = 0
		}

		toursResponse = append(toursResponse, contract.MyTourResponse{
			ID:                tour.ID,
			Name:              tour.Name,
			Dates:             tour.Dates,
			DepartureTime:     tour.DepartureTime,
			ArrivalTime:       tour.ArrivalTime,
			MaxPeople:         tour.MaxPeople,
			Description:       tour.Description,
			Images:            tour.Images,
			Price:             tour.Price,
			CreatedAt:         tour.CreatedAt.Format("2006-01-02 15:04:05"),
			ReservationsCount: reservationsCount,
		})
	}

	totalPages := util.CalculateTotalPages(total, limit)

	response := &contract.GetMyToursResponse{
		Success: true,
		Tours:   toursResponse,
		Pagination: contract.PaginationResponse{
			CurrentPage:  page,
			TotalPages:   totalPages,
			TotalItems:   int(total),
			ItemsPerPage: limit,
		},
	}

	return response, nil
}

// Delete - deleta um passeio
func (s *TourService) Delete(tourID, companyID int) (*contract.DeleteTourResponse, error) {

	exists, err := s.TourRepository.IsOwnedByCompany(tourID, companyID)
	if err != nil {
		return nil, util.WrapError("Erro ao verificar propriedade do passeio", err, http.StatusInternalServerError)
	}
	if !exists {
		return nil, util.WrapError("Passeio não encontrado ou você não tem permissão para deletá-lo", nil, http.StatusForbidden)
	}

	reservationsCount, err := s.TourRepository.CountReservationsByTourID(tourID)
	if err != nil {
		return nil, util.WrapError("Erro ao verificar reservas do passeio", err, http.StatusInternalServerError)
	}
	if reservationsCount > 0 {
		return nil, util.WrapError("Não é possível deletar passeio com reservas ativas", nil, http.StatusConflict)
	}

	if err := s.TourRepository.Delete(tourID); err != nil {
		return nil, util.WrapError("Erro ao deletar passeio", err, http.StatusInternalServerError)
	}

	response := &contract.DeleteTourResponse{
		Success: true,
		Message: "Passeio deletado com sucesso",
	}

	return response, nil
}
