package repository

import (
	"time"

	"github.com/jampa_trip/internal/model"
	"gorm.io/gorm"
)

// ReservaRepository - objeto de contexto
type ReservaRepository struct {
	DB *gorm.DB
}

// ReservaRepositoryNew - construtor do objeto
func ReservaRepositoryNew(DB *gorm.DB) *ReservaRepository {
	return &ReservaRepository{
		DB: DB,
	}
}

// Create - cria uma nova reserva
func (r *ReservaRepository) Create(reserva *model.Reserva) error {
	return r.DB.Create(reserva).Error
}

// GetByID - busca uma reserva pelo ID
func (r *ReservaRepository) GetByID(id int) (*model.Reserva, error) {
	var reserva model.Reserva
	err := r.DB.Preload("Cliente").Preload("Empresa").Preload("Pagamento").First(&reserva, id).Error
	if err != nil {
		return nil, err
	}
	return &reserva, nil
}

// GetByClienteID - busca reservas por cliente
func (r *ReservaRepository) GetByClienteID(clienteID int, page, limit int) ([]model.Reserva, int64, error) {
	var reservas []model.Reserva
	var total int64

	offset := (page - 1) * limit

	query := r.DB.Model(&model.Reserva{}).Where("cliente_id = ?", clienteID)

	// Contar total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Buscar registros
	err := query.Preload("Cliente").Preload("Empresa").Preload("Pagamento").
		Offset(offset).Limit(limit).
		Order("momento_criacao DESC").
		Find(&reservas).Error

	return reservas, total, err
}

// GetByEmpresaID - busca reservas por empresa
func (r *ReservaRepository) GetByEmpresaID(empresaID int, page, limit int) ([]model.Reserva, int64, error) {
	var reservas []model.Reserva
	var total int64

	offset := (page - 1) * limit

	query := r.DB.Model(&model.Reserva{}).Where("empresa_id = ?", empresaID)

	// Contar total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Buscar registros
	err := query.Preload("Cliente").Preload("Empresa").Preload("Pagamento").
		Offset(offset).Limit(limit).
		Order("momento_criacao DESC").
		Find(&reservas).Error

	return reservas, total, err
}

// GetByStatus - busca reservas por status
func (r *ReservaRepository) GetByStatus(status string, page, limit int) ([]model.Reserva, int64, error) {
	var reservas []model.Reserva
	var total int64

	offset := (page - 1) * limit

	query := r.DB.Model(&model.Reserva{}).Where("status = ?", status)

	// Contar total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Buscar registros
	err := query.Preload("Cliente").Preload("Empresa").Preload("Pagamento").
		Offset(offset).Limit(limit).
		Order("momento_criacao DESC").
		Find(&reservas).Error

	return reservas, total, err
}

// Update - atualiza uma reserva
func (r *ReservaRepository) Update(reserva *model.Reserva) error {
	return r.DB.Save(reserva).Error
}

// UpdateStatus - atualiza apenas o status de uma reserva
func (r *ReservaRepository) UpdateStatus(id int, status string) error {
	return r.DB.Model(&model.Reserva{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":              status,
			"momento_atualizacao": time.Now(),
		}).Error
}

// Cancel - cancela uma reserva
func (r *ReservaRepository) Cancel(id int) error {
	now := time.Now()
	return r.DB.Model(&model.Reserva{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":               string(model.StatusReservaCancelada),
			"momento_cancelamento": &now,
			"momento_atualizacao":  time.Now(),
		}).Error
}

// Delete - remove uma reserva
func (r *ReservaRepository) Delete(id int) error {
	return r.DB.Delete(&model.Reserva{}, id).Error
}

// GetByDateRange - busca reservas por período
func (r *ReservaRepository) GetByDateRange(startDate, endDate time.Time, page, limit int) ([]model.Reserva, int64, error) {
	var reservas []model.Reserva
	var total int64

	offset := (page - 1) * limit

	query := r.DB.Model(&model.Reserva{}).Where("data_passeio BETWEEN ? AND ?", startDate, endDate)

	// Contar total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Buscar registros
	err := query.Preload("Cliente").Preload("Empresa").Preload("Pagamento").
		Offset(offset).Limit(limit).
		Order("data_passeio ASC").
		Find(&reservas).Error

	return reservas, total, err
}

// GetUpcoming - busca reservas futuras
func (r *ReservaRepository) GetUpcoming(clienteID int, page, limit int) ([]model.Reserva, int64, error) {
	var reservas []model.Reserva
	var total int64

	offset := (page - 1) * limit
	now := time.Now()

	query := r.DB.Model(&model.Reserva{}).Where("cliente_id = ? AND data_passeio > ?", clienteID, now)

	// Contar total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Buscar registros
	err := query.Preload("Cliente").Preload("Empresa").Preload("Pagamento").
		Offset(offset).Limit(limit).
		Order("data_passeio ASC").
		Find(&reservas).Error

	return reservas, total, err
}

// GetHistory - busca histórico de reservas
func (r *ReservaRepository) GetHistory(clienteID int, page, limit int) ([]model.Reserva, int64, error) {
	var reservas []model.Reserva
	var total int64

	offset := (page - 1) * limit
	now := time.Now()

	query := r.DB.Model(&model.Reserva{}).Where("cliente_id = ? AND data_passeio <= ?", clienteID, now)

	// Contar total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Buscar registros
	err := query.Preload("Cliente").Preload("Empresa").Preload("Pagamento").
		Offset(offset).Limit(limit).
		Order("data_passeio DESC").
		Find(&reservas).Error

	return reservas, total, err
}
