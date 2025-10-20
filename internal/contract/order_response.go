package contract

import "time"

// OrderResponse - representa a resposta de uma ordem
type OrderResponse struct {
	ID                 int                   `json:"id"`
	ClienteID          int                   `json:"cliente_id"`
	EmpresaID          int                   `json:"empresa_id"`
	MercadoPagoOrderID string                `json:"mercado_pago_order_id"`
	ExternalReference  string                `json:"external_reference"`
	TotalAmount        float64               `json:"total_amount"`
	Currency           string                `json:"currency"`
	Status             string                `json:"status"`
	StatusDetail       string                `json:"status_detail"`
	Description        string                `json:"description"`
	NotificationURL    string                `json:"notification_url"`
	MomentoCriacao     time.Time             `json:"momento_criacao"`
	MomentoAtualizacao time.Time             `json:"momento_atualizacao"`
	MomentoExpiracao   *time.Time            `json:"momento_expiracao,omitempty"`
	StatusDisplay      string                `json:"status_display"`
	Transactions       []TransactionResponse `json:"transactions,omitempty"`
}

// TransactionResponse - representa a resposta de uma transação
type TransactionResponse struct {
	ID                       int       `json:"id"`
	OrderID                  int       `json:"order_id"`
	MercadoPagoTransactionID string    `json:"mercado_pago_transaction_id"`
	PaymentID                int64     `json:"payment_id"`
	Type                     string    `json:"type"`
	Status                   string    `json:"status"`
	StatusDetail             string    `json:"status_detail"`
	Amount                   float64   `json:"amount"`
	Currency                 string    `json:"currency"`
	PaymentMethodID          string    `json:"payment_method_id"`
	PaymentTypeID            string    `json:"payment_type_id"`
	MomentoCriacao           time.Time `json:"momento_criacao"`
	MomentoAtualizacao       time.Time `json:"momento_atualizacao"`
	StatusDisplay            string    `json:"status_display"`
	TypeDisplay              string    `json:"type_display"`
}

// CreateOrderResponse - representa a resposta da criação de uma ordem
type CreateOrderResponse struct {
	Order   OrderResponse `json:"order"`
	Message string        `json:"message"`
}

// GetOrderResponse - representa a resposta de obter uma ordem
type GetOrderResponse struct {
	Order OrderResponse `json:"order"`
}

// ListOrdersResponse - representa a resposta da listagem de ordens
type ListOrdersResponse struct {
	Orders  []OrderResponse `json:"orders"`
	Total   int64           `json:"total"`
	Offset  int             `json:"offset"`
	Limit   int             `json:"limit"`
	HasMore bool            `json:"has_more"`
}

// CaptureOrderResponse - representa a resposta da captura de uma ordem
type CaptureOrderResponse struct {
	Order   OrderResponse `json:"order"`
	Message string        `json:"message"`
}

// CancelOrderResponse - representa a resposta do cancelamento de uma ordem
type CancelOrderResponse struct {
	Order   OrderResponse `json:"order"`
	Message string        `json:"message"`
}

// RefundOrderResponse - representa a resposta do reembolso de uma ordem
type RefundOrderResponse struct {
	Order         OrderResponse         `json:"order"`
	RefundDetails RefundDetailsResponse `json:"refund_details"`
	Message       string                `json:"message"`
}

// RefundDetailsResponse - representa os detalhes do reembolso
type RefundDetailsResponse struct {
	RefundID    int64   `json:"refund_id"`
	Amount      float64 `json:"amount"`
	Status      string  `json:"status"`
	DateCreated string  `json:"date_created"`
}

// AddTransactionResponse - representa a resposta de adicionar uma transação
type AddTransactionResponse struct {
	Transaction TransactionResponse `json:"transaction"`
	Order       OrderResponse       `json:"order"`
	Message     string              `json:"message"`
}

// UpdateTransactionResponse - representa a resposta de atualizar uma transação
type UpdateTransactionResponse struct {
	Transaction TransactionResponse `json:"transaction"`
	Message     string              `json:"message"`
}

// ProcessOrderResponse - representa a resposta do processamento de uma ordem
type ProcessOrderResponse struct {
	Order          OrderResponse         `json:"order"`
	PaymentDetails ProcessPaymentDetails `json:"payment_details"`
	Message        string                `json:"message"`
}

// ProcessPaymentDetails - representa os detalhes do pagamento processado
type ProcessPaymentDetails struct {
	PaymentID         int64   `json:"payment_id"`
	Status            string  `json:"status"`
	StatusDetail      string  `json:"status_detail"`
	PaymentMethodID   string  `json:"payment_method_id"`
	PaymentTypeID     string  `json:"payment_type_id"`
	TransactionAmount float64 `json:"transaction_amount"`
	DateCreated       string  `json:"date_created"`
	DateApproved      string  `json:"date_approved,omitempty"`
}
