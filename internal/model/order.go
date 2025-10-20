package model

import "time"

// Order - representa uma ordem do Mercado Pago
type Order struct {
	ID                 int        `gorm:"column:id;primaryKey;autoIncrement"`
	ClienteID          int        `gorm:"column:cliente_id;not null"`
	EmpresaID          int        `gorm:"column:empresa_id;not null"`
	MercadoPagoOrderID string     `gorm:"column:mercado_pago_order_id;uniqueIndex"`
	ExternalReference  string     `gorm:"column:external_reference;index"`
	TotalAmount        float64    `gorm:"column:total_amount;not null;type:decimal(10,2)"`
	Currency           string     `gorm:"column:currency;not null;default:'BRL'"`
	Status             string     `gorm:"column:status;not null;default:'pending'"`
	StatusDetail       string     `gorm:"column:status_detail"`
	Description        string     `gorm:"column:description;type:text"`
	NotificationURL    string     `gorm:"column:notification_url"`
	MomentoCriacao     time.Time  `gorm:"column:momento_criacao;not null;default:CURRENT_TIMESTAMP"`
	MomentoAtualizacao time.Time  `gorm:"column:momento_atualizacao;not null;default:CURRENT_TIMESTAMP"`
	MomentoExpiracao   *time.Time `gorm:"column:momento_expiracao"`

	// Relacionamentos
	Cliente      Client        `gorm:"foreignKey:ClienteID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Empresa      Company       `gorm:"foreignKey:EmpresaID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Transactions []Transaction `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// TableName - especifica o nome da tabela no banco de dados
func (Order) TableName() string {
	return "orders"
}

// OrderStatus - define os possíveis status de uma ordem
type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusPaid       OrderStatus = "paid"
	OrderStatusAuthorized OrderStatus = "authorized"
	OrderStatusCancelled  OrderStatus = "cancelled"
	OrderStatusRefunded   OrderStatus = "refunded"
	OrderStatusExpired    OrderStatus = "expired"
)

// IsValid - valida se a ordem possui dados válidos
func (o *Order) IsValid() bool {
	return o.TotalAmount > 0 && o.Status != "" && o.Currency != ""
}

// IsPaid - verifica se a ordem está paga
func (o *Order) IsPaid() bool {
	return o.Status == string(OrderStatusPaid)
}

// IsCancelled - verifica se a ordem está cancelada
func (o *Order) IsCancelled() bool {
	return o.Status == string(OrderStatusCancelled)
}

// IsAuthorized - verifica se a ordem está autorizada
func (o *Order) IsAuthorized() bool {
	return o.Status == string(OrderStatusAuthorized)
}

// IsRefunded - verifica se a ordem foi reembolsada
func (o *Order) IsRefunded() bool {
	return o.Status == string(OrderStatusRefunded)
}

// IsPending - verifica se a ordem está pendente
func (o *Order) IsPending() bool {
	return o.Status == string(OrderStatusPending)
}

// UpdateStatus - atualiza o status da ordem
func (o *Order) UpdateStatus(status OrderStatus) {
	o.Status = string(status)
	o.MomentoAtualizacao = time.Now()
}

// GetStatusDisplay - retorna uma descrição amigável do status
func (o *Order) GetStatusDisplay() string {
	switch o.Status {
	case string(OrderStatusPending):
		return "Pendente"
	case string(OrderStatusProcessing):
		return "Em Processamento"
	case string(OrderStatusPaid):
		return "Pago"
	case string(OrderStatusAuthorized):
		return "Autorizado"
	case string(OrderStatusCancelled):
		return "Cancelado"
	case string(OrderStatusRefunded):
		return "Reembolsado"
	case string(OrderStatusExpired):
		return "Expirado"
	default:
		return "Desconhecido"
	}
}

// IsValidOrderStatus - valida se o status é válido
func IsValidOrderStatus(status OrderStatus) bool {
	switch status {
	case OrderStatusPending, OrderStatusProcessing, OrderStatusPaid,
		OrderStatusAuthorized, OrderStatusCancelled, OrderStatusRefunded, OrderStatusExpired:
		return true
	default:
		return false
	}
}
