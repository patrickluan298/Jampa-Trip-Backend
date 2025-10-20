package repository

import (
	"github.com/jampa_trip/internal/model"
	"gorm.io/gorm"
)

// TransactionRepository - objeto de contexto para operações com transactions
type TransactionRepository struct {
	DB *gorm.DB
}

// TransactionRepositoryNew - construtor do objeto
func TransactionRepositoryNew(DB *gorm.DB) *TransactionRepository {
	return &TransactionRepository{
		DB: DB,
	}
}

// Create - cria uma nova transação
func (r *TransactionRepository) Create(transaction *model.Transaction) error {
	return r.DB.Create(transaction).Error
}

// Update - atualiza uma transação
func (r *TransactionRepository) Update(transaction *model.Transaction) error {
	return r.DB.Save(transaction).Error
}

// GetByID - busca uma transação pelo ID
func (r *TransactionRepository) GetByID(id int) (*model.Transaction, error) {
	var transaction model.Transaction
	err := r.DB.Where("id = ?", id).First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

// GetByMercadoPagoTransactionID - busca uma transação pelo ID do Mercado Pago
func (r *TransactionRepository) GetByMercadoPagoTransactionID(txID string) (*model.Transaction, error) {
	var transaction model.Transaction
	err := r.DB.Where("mercado_pago_transaction_id = ?", txID).First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

// GetByPaymentID - busca uma transação pelo ID do pagamento
func (r *TransactionRepository) GetByPaymentID(paymentID int64) (*model.Transaction, error) {
	var transaction model.Transaction
	err := r.DB.Where("payment_id = ?", paymentID).First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

// ListByOrderID - lista todas as transações de uma ordem
func (r *TransactionRepository) ListByOrderID(orderID int) ([]model.Transaction, error) {
	var transactions []model.Transaction
	err := r.DB.Where("order_id = ?", orderID).
		Order("momento_criacao DESC").
		Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

// Delete - deleta uma transação
func (r *TransactionRepository) Delete(id int) error {
	return r.DB.Delete(&model.Transaction{}, id).Error
}

// DeleteByOrderID - deleta todas as transações de uma ordem
func (r *TransactionRepository) DeleteByOrderID(orderID int) error {
	return r.DB.Where("order_id = ?", orderID).Delete(&model.Transaction{}).Error
}
