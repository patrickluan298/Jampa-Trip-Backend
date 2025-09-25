package model

import "time"

// Pagamento - representa a entidade de pagamento
type Pagamento struct {
	ID                   int        `gorm:"column:id;primaryKey;autoIncrement"`
	ClienteID            int        `gorm:"column:cliente_id;not null"`
	EmpresaID            int        `gorm:"column:empresa_id;not null"`
	MercadoPagoOrderID   string     `gorm:"column:mercado_pago_order_id;uniqueIndex"`
	MercadoPagoPaymentID string     `gorm:"column:mercado_pago_payment_id;index"`
	Status               string     `gorm:"column:status;not null;default:'pending'"`
	StatusDetail         string     `gorm:"column:status_detail"`
	Valor                float64    `gorm:"column:valor;not null;type:decimal(10,2)"`
	Moeda                string     `gorm:"column:moeda;not null;default:'BRL'"`
	MetodoPagamento      string     `gorm:"column:metodo_pagamento;not null"`
	Descricao            string     `gorm:"column:descricao"`
	NumeroParcelas       int        `gorm:"column:numero_parcelas;default:1"`
	TokenCartao          string     `gorm:"column:token_cartao"`
	ChavePIX             string     `gorm:"column:chave_pix"`
	QRCode               string     `gorm:"column:qr_code;type:text"`
	MomentoCriacao       time.Time  `gorm:"column:momento_criacao;not null;default:CURRENT_TIMESTAMP"`
	MomentoAtualizacao   time.Time  `gorm:"column:momento_atualizacao;not null;default:CURRENT_TIMESTAMP"`
	MomentoAprovacao     *time.Time `gorm:"column:momento_aprovacao"`
	MomentoCancelamento  *time.Time `gorm:"column:momento_cancelamento"`

	// Relacionamentos
	Cliente Cliente `gorm:"foreignKey:ClienteID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Empresa Empresa `gorm:"foreignKey:EmpresaID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// TableName - especifica o nome da tabela no banco de dados
func (Pagamento) TableName() string {
	return "pagamentos"
}
