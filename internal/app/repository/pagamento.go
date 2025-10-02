package repository

import (
	"github.com/jampa_trip/internal/app/model"
	"gorm.io/gorm"
)

// PagamentoRepository - objeto de contexto
type PagamentoRepository struct {
	DB *gorm.DB
}

// PagamentoRepositoryNew - construtor do objeto
func PagamentoRepositoryNew(DB *gorm.DB) *PagamentoRepository {
	return &PagamentoRepository{
		DB: DB,
	}
}

// Create - cria um novo pagamento
func (r *PagamentoRepository) Create(pagamento *model.Pagamento) error {
	return r.DB.Create(pagamento).Error
}

// Update - atualiza um pagamento
func (r *PagamentoRepository) Update(pagamento *model.Pagamento) error {
	return r.DB.Save(pagamento).Error
}

// GetByMercadoPagoPaymentID - busca um pagamento pelo ID do Mercado Pago
func (r *PagamentoRepository) GetByMercadoPagoPaymentID(paymentID string) (*model.Pagamento, error) {
	var pagamento model.Pagamento
	err := r.DB.Where("mercado_pago_payment_id = ?", paymentID).First(&pagamento).Error
	if err != nil {
		return nil, err
	}
	return &pagamento, nil
}

// GetByClienteID - lista pagamentos de um cliente
func (r *PagamentoRepository) GetByClienteID(clienteID int) ([]model.Pagamento, error) {
	var pagamentos []model.Pagamento
	err := r.DB.Where("cliente_id = ?", clienteID).Order("momento_criacao DESC").Find(&pagamentos).Error
	return pagamentos, err
}

// GetByEmpresaID - lista pagamentos de uma empresa
func (r *PagamentoRepository) GetByEmpresaID(empresaID int) ([]model.Pagamento, error) {
	var pagamentos []model.Pagamento
	err := r.DB.Where("empresa_id = ?", empresaID).Order("momento_criacao DESC").Find(&pagamentos).Error
	return pagamentos, err
}
