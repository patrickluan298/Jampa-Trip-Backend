package repository

import (
	"github.com/jampa_trip/internal/model"
	"github.com/jampa_trip/internal/query"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// TourRepository - objeto de contexto
type TourRepository struct {
	DB *gorm.DB
}

// TourRepositoryNew - construtor do objeto
func TourRepositoryNew(DB *gorm.DB) *TourRepository {
	return &TourRepository{
		DB: DB,
	}
}

// Create - cria um novo passeio
func (r *TourRepository) Create(tour *model.Tour) error {
	err := r.DB.Raw(query.CreateTour,
		tour.CompanyID,
		tour.Name,
		pq.Array(tour.Dates),
		tour.DepartureTime,
		tour.ArrivalTime,
		tour.MaxPeople,
		tour.Description,
		pq.Array(tour.Images),
		tour.Price,
	).Row().Scan(&tour.ID)

	return err
}

// Update - atualiza um passeio existente
func (r *TourRepository) Update(tour *model.Tour) error {
	err := r.DB.Raw(query.UpdateTour,
		tour.Name,
		pq.Array(tour.Dates),
		tour.DepartureTime,
		tour.ArrivalTime,
		tour.MaxPeople,
		tour.Description,
		pq.Array(tour.Images),
		tour.Price,
		tour.ID,
	).Row().Scan()

	return err
}

// GetByID - busca um passeio pelo ID
func (r *TourRepository) GetByID(id int) (*model.Tour, error) {
	tour := &model.Tour{}
	var companyName string

	err := r.DB.Raw(query.GetTourByID, id).Row().Scan(
		&tour.ID,
		&tour.CompanyID,
		&tour.Name,
		pq.Array(&tour.Dates),
		&tour.DepartureTime,
		&tour.ArrivalTime,
		&tour.MaxPeople,
		&tour.Description,
		pq.Array(&tour.Images),
		&tour.Price,
		&tour.CreatedAt,
		&tour.UpdatedAt,
		&companyName,
	)

	if err != nil {
		return nil, err
	}

	return tour, nil
}

// List - busca todos os passeios com filtro e paginação
func (r *TourRepository) List(search string, page, limit int) ([]*model.Tour, int64, error) {
	offset := (page - 1) * limit
	searchPattern := "%" + search + "%"

	rows, err := r.DB.Raw(query.ListTours, search, searchPattern, limit, offset).Rows()
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var tours []*model.Tour
	for rows.Next() {
		tour := &model.Tour{}
		var companyName string

		err := rows.Scan(
			&tour.ID,
			&tour.CompanyID,
			&tour.Name,
			pq.Array(&tour.Dates),
			&tour.DepartureTime,
			&tour.ArrivalTime,
			&tour.MaxPeople,
			&tour.Description,
			pq.Array(&tour.Images),
			&tour.Price,
			&tour.CreatedAt,
			&tour.UpdatedAt,
			&companyName,
		)
		if err != nil {
			return nil, 0, err
		}
		tours = append(tours, tour)
	}

	var total int64
	err = r.DB.Raw(query.CountTours, search, searchPattern).Row().Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return tours, total, nil
}

// ListByCompanyID - busca passeios de uma empresa específica
func (r *TourRepository) ListByCompanyID(companyID int, page, limit int) ([]*model.Tour, int64, error) {
	offset := (page - 1) * limit

	rows, err := r.DB.Raw(query.ListMyTours, companyID, limit, offset).Rows()
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var tours []*model.Tour
	for rows.Next() {
		tour := &model.Tour{}
		var reservationsCount int

		err := rows.Scan(
			&tour.ID,
			&tour.CompanyID,
			&tour.Name,
			pq.Array(&tour.Dates),
			&tour.DepartureTime,
			&tour.ArrivalTime,
			&tour.MaxPeople,
			&tour.Description,
			pq.Array(&tour.Images),
			&tour.Price,
			&tour.CreatedAt,
			&reservationsCount,
		)
		if err != nil {
			return nil, 0, err
		}
		tours = append(tours, tour)
	}

	var total int64
	err = r.DB.Raw(query.CountMyTours, companyID).Row().Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return tours, total, nil
}

// Delete - deleta um passeio
func (r *TourRepository) Delete(id int) error {
	err := r.DB.Raw(query.DeleteTour, id).Row().Scan()
	return err
}

// IsOwnedByCompany - verifica se o passeio pertence à empresa
func (r *TourRepository) IsOwnedByCompany(tourID, companyID int) (bool, error) {
	var count int
	err := r.DB.Raw(query.CheckTourOwnership, tourID, companyID).Row().Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CountReservationsByTourID - conta reservas ativas de um passeio
func (r *TourRepository) CountReservationsByTourID(tourID int) (int, error) {
	var count int
	err := r.DB.Raw(query.CountReservationsByTourID, tourID).Row().Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetTourWithCompanyName - busca passeio com nome da empresa
func (r *TourRepository) GetTourWithCompanyName(id int) (*model.Tour, string, error) {
	tour := &model.Tour{}
	var companyName string

	err := r.DB.Raw(query.GetTourByID, id).Row().Scan(
		&tour.ID,
		&tour.CompanyID,
		&tour.Name,
		pq.Array(&tour.Dates),
		&tour.DepartureTime,
		&tour.ArrivalTime,
		&tour.MaxPeople,
		&tour.Description,
		pq.Array(&tour.Images),
		&tour.Price,
		&tour.CreatedAt,
		&tour.UpdatedAt,
		&companyName,
	)

	if err != nil {
		return nil, "", err
	}

	return tour, companyName, nil
}
