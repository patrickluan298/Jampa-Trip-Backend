package repository

import (
	"github.com/jampa_trip/internal/model"
	"gorm.io/gorm"
)

// OrderRepository - objeto de contexto para operações com orders
type OrderRepository struct {
	DB *gorm.DB
}

// OrderRepositoryNew - construtor do objeto
func OrderRepositoryNew(DB *gorm.DB) *OrderRepository {
	return &OrderRepository{
		DB: DB,
	}
}

// Create - cria uma nova ordem
func (r *OrderRepository) Create(order *model.Order) error {
	return r.DB.Create(order).Error
}

// Update - atualiza uma ordem
func (r *OrderRepository) Update(order *model.Order) error {
	return r.DB.Save(order).Error
}

// GetByID - busca uma ordem pelo ID
func (r *OrderRepository) GetByID(id int) (*model.Order, error) {
	var order model.Order
	err := r.DB.Preload("Transactions").Where("id = ?", id).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// GetByMercadoPagoOrderID - busca uma ordem pelo ID do Mercado Pago
func (r *OrderRepository) GetByMercadoPagoOrderID(orderID string) (*model.Order, error) {
	var order model.Order
	err := r.DB.Preload("Transactions").Where("mercado_pago_order_id = ?", orderID).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// GetByExternalReference - busca uma ordem pela referência externa
func (r *OrderRepository) GetByExternalReference(ref string) (*model.Order, error) {
	var order model.Order
	err := r.DB.Preload("Transactions").Where("external_reference = ?", ref).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// ListByClienteID - lista ordens de um cliente com paginação
func (r *OrderRepository) ListByClienteID(clienteID int, limit, offset int) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64

	if err := r.DB.Model(&model.Order{}).Where("cliente_id = ?", clienteID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.DB.Preload("Transactions").
		Where("cliente_id = ?", clienteID).
		Order("momento_criacao DESC").
		Limit(limit).
		Offset(offset).
		Find(&orders).Error

	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// ListByEmpresaID - lista ordens de uma empresa com paginação
func (r *OrderRepository) ListByEmpresaID(empresaID int, limit, offset int) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64

	if err := r.DB.Model(&model.Order{}).Where("empresa_id = ?", empresaID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.DB.Preload("Transactions").
		Where("empresa_id = ?", empresaID).
		Order("momento_criacao DESC").
		Limit(limit).
		Offset(offset).
		Find(&orders).Error

	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// ListByStatus - lista ordens por status com paginação
func (r *OrderRepository) ListByStatus(status string, limit, offset int) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64

	if err := r.DB.Model(&model.Order{}).Where("status = ?", status).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.DB.Preload("Transactions").
		Where("status = ?", status).
		Order("momento_criacao DESC").
		Limit(limit).
		Offset(offset).
		Find(&orders).Error

	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// List - lista todas as ordens com filtros e paginação
func (r *OrderRepository) List(filters map[string]interface{}, limit, offset int) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64

	query := r.DB.Model(&model.Order{})

	for key, value := range filters {
		if value != nil && value != "" && value != 0 {
			query = query.Where(key+" = ?", value)
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("Transactions").
		Order("momento_criacao DESC").
		Limit(limit).
		Offset(offset).
		Find(&orders).Error

	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}
