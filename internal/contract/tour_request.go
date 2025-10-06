package contract

import (
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jampa_trip/pkg/util"
)

// CreateTourRequest - objeto de request do endpoint de criação de passeio
type CreateTourRequest struct {
	Name          string   `json:"name"`
	Dates         []string `json:"dates"`
	DepartureTime string   `json:"departure_time"`
	ArrivalTime   string   `json:"arrival_time"`
	MaxPeople     int      `json:"max_people"`
	Description   string   `json:"description"`
	Images        []string `json:"images"`
	Price         float64  `json:"price"`
}

// Validate - valida os campos da requisição de criação
func (receiver CreateTourRequest) Validate() error {
	err := validation.ValidateStruct(&receiver,
		validation.Field(&receiver.Name, validation.Required, validation.Length(3, 255)),
		validation.Field(&receiver.Dates, validation.Required, validation.Length(1, 0)),
		validation.Field(&receiver.DepartureTime, validation.Required, util.TimeValidator()),
		validation.Field(&receiver.ArrivalTime, validation.Required, util.TimeValidator()),
		validation.Field(&receiver.MaxPeople, validation.Required, validation.Min(1)),
		validation.Field(&receiver.Description, validation.Length(0, 1000)),
		validation.Field(&receiver.Price, validation.Required, validation.Min(0.01)),
	)

	if err != nil {
		return util.WrapError(util.FormatarErroValidacao(err).Error(), err, http.StatusUnprocessableEntity)
	}

	if err := util.ValidateDates(receiver.Dates); err != nil {
		return err
	}

	if err := util.ValidateImageURLs(receiver.Images); err != nil {
		return err
	}

	return nil
}

// UpdateTourRequest - objeto de request do endpoint de atualização de passeio
type UpdateTourRequest struct {
	ID            int      `json:"id"`
	Name          string   `json:"name"`
	Dates         []string `json:"dates"`
	DepartureTime string   `json:"departure_time"`
	ArrivalTime   string   `json:"arrival_time"`
	MaxPeople     int      `json:"max_people"`
	Description   string   `json:"description"`
	Images        []string `json:"images"`
	Price         float64  `json:"price"`
}

// Validate - valida os campos da requisição de atualização
func (receiver UpdateTourRequest) Validate() error {
	err := validation.ValidateStruct(&receiver,
		validation.Field(&receiver.ID, validation.Required, validation.Min(1)),
		validation.Field(&receiver.Name, validation.Required, validation.Length(3, 255)),
		validation.Field(&receiver.Dates, validation.Required, validation.Length(1, 0)),
		validation.Field(&receiver.DepartureTime, validation.Required, util.TimeValidator()),
		validation.Field(&receiver.ArrivalTime, validation.Required, util.TimeValidator()),
		validation.Field(&receiver.MaxPeople, validation.Required, validation.Min(1)),
		validation.Field(&receiver.Description, validation.Length(0, 1000)),
		validation.Field(&receiver.Price, validation.Required, validation.Min(0.01)),
	)

	if err != nil {
		return util.WrapError(util.FormatarErroValidacao(err).Error(), err, http.StatusUnprocessableEntity)
	}

	if err := util.ValidateDates(receiver.Dates); err != nil {
		return err
	}

	if err := util.ValidateImageURLs(receiver.Images); err != nil {
		return err
	}

	return nil
}

// ListToursRequest - objeto de request do endpoint de listagem de passeios
type ListToursRequest struct {
	Search string `json:"search"`
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
}

// Validate - valida os campos da requisição de listagem
func (receiver ListToursRequest) Validate() error {
	err := validation.ValidateStruct(&receiver,
		validation.Field(&receiver.Page, validation.Min(1)),
		validation.Field(&receiver.Limit, validation.Min(1), validation.Max(100)),
	)

	if err != nil {
		return util.WrapError(util.FormatarErroValidacao(err).Error(), err, http.StatusUnprocessableEntity)
	}

	config := util.DefaultPagination()
	if receiver.Page == 0 {
		receiver.Page = config.Page
	}
	if receiver.Limit == 0 {
		receiver.Limit = config.Limit
	}

	return nil
}
