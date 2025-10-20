package contract

import validation "github.com/go-ozzo/ozzo-validation"

// CreateOrderRequest - representa a requisição para criar uma ordem
type CreateOrderRequest struct {
	ClienteID         int                `json:"cliente_id"`
	EmpresaID         int                `json:"empresa_id"`
	ExternalReference string             `json:"external_reference"`
	TotalAmount       float64            `json:"total_amount"`
	Items             []OrderItemRequest `json:"items"`
	Payer             OrderPayerRequest  `json:"payer"`
	Description       string             `json:"description"`
	NotificationURL   string             `json:"notification_url"`
	Metadata          map[string]string  `json:"metadata"`
}

// Validate - valida os campos da requisição
func (r *CreateOrderRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ClienteID, validation.Required, validation.Min(1)),
		validation.Field(&r.EmpresaID, validation.Required, validation.Min(1)),
		validation.Field(&r.ExternalReference, validation.Required, validation.Length(1, 256)),
		validation.Field(&r.TotalAmount, validation.Required, validation.Min(0.01)),
		validation.Field(&r.Items, validation.Required, validation.Length(1, 100)),
		validation.Field(&r.Description, validation.Length(0, 600)),
	)
}

// OrderItemRequest - representa um item da ordem
type OrderItemRequest struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	PictureURL  string  `json:"picture_url"`
	CategoryID  string  `json:"category_id"`
	Quantity    int     `json:"quantity"`
	CurrencyID  string  `json:"currency_id"`
	UnitPrice   float64 `json:"unit_price"`
}

// Validate - valida os campos do item
func (r *OrderItemRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ID, validation.Required, validation.Length(1, 256)),
		validation.Field(&r.Title, validation.Required, validation.Length(1, 256)),
		validation.Field(&r.Quantity, validation.Required, validation.Min(1)),
		validation.Field(&r.CurrencyID, validation.Required, validation.In("BRL", "USD", "EUR", "ARS", "CLP", "COP", "MXN", "PEN", "UYU")),
		validation.Field(&r.UnitPrice, validation.Required, validation.Min(0.01)),
	)
}

// OrderPayerRequest - representa os dados do pagador
type OrderPayerRequest struct {
	Name    string              `json:"name"`
	Email   string              `json:"email"`
	Phone   OrderPhoneRequest   `json:"phone"`
	Address OrderAddressRequest `json:"address"`
}

// OrderPhoneRequest - representa o telefone do pagador
type OrderPhoneRequest struct {
	AreaCode string `json:"area_code"`
	Number   string `json:"number"`
}

// OrderAddressRequest - representa o endereço do pagador
type OrderAddressRequest struct {
	StreetName   string `json:"street_name"`
	StreetNumber int    `json:"street_number"`
	ZipCode      string `json:"zip_code"`
}

// GetOrderRequest - representa a requisição para obter uma ordem
type GetOrderRequest struct {
	OrderID string `json:"order_id"`
}

// Validate - valida os campos da requisição
func (r *GetOrderRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.OrderID, validation.Required, validation.Length(1, 256)),
	)
}

// ListOrdersRequest - representa a requisição para listar ordens
type ListOrdersRequest struct {
	ClienteID         int    `json:"cliente_id"`
	EmpresaID         int    `json:"empresa_id"`
	ExternalReference string `json:"external_reference"`
	Status            string `json:"status"`
	Offset            int    `json:"offset"`
	Limit             int    `json:"limit"`
}

// Validate - valida os campos da requisição
func (r *ListOrdersRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ClienteID, validation.Min(0)),
		validation.Field(&r.EmpresaID, validation.Min(0)),
		validation.Field(&r.Status, validation.In("", "pending", "processing", "paid", "authorized", "cancelled", "refunded", "expired")),
		validation.Field(&r.Offset, validation.Min(0)),
		validation.Field(&r.Limit, validation.Min(1), validation.Max(100)),
	)
}

// CaptureOrderRequest - representa a requisição para capturar uma ordem
type CaptureOrderRequest struct {
	Amount *float64 `json:"amount,omitempty"`
}

// Validate - valida os campos da requisição
func (r *CaptureOrderRequest) Validate() error {
	if r.Amount != nil {
		return validation.ValidateStruct(r,
			validation.Field(&r.Amount, validation.Min(0.01)),
		)
	}
	return nil
}

// RefundOrderRequest - representa a requisição para reembolsar uma ordem
type RefundOrderRequest struct {
	Amount   *float64          `json:"amount,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// Validate - valida os campos da requisição
func (r *RefundOrderRequest) Validate() error {
	if r.Amount != nil {
		return validation.ValidateStruct(r,
			validation.Field(&r.Amount, validation.Min(0.01)),
		)
	}
	return nil
}

// AddTransactionRequest - representa a requisição para adicionar uma transação
type AddTransactionRequest struct {
	PaymentID       int64   `json:"payment_id"`
	PaymentMethodID string  `json:"payment_method_id"`
	Amount          float64 `json:"amount"`
	Currency        string  `json:"currency"`
}

// Validate - valida os campos da requisição
func (r *AddTransactionRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.PaymentID, validation.Required, validation.Min(1)),
		validation.Field(&r.PaymentMethodID, validation.Required, validation.Length(1, 100)),
		validation.Field(&r.Amount, validation.Required, validation.Min(0.01)),
		validation.Field(&r.Currency, validation.Required, validation.In("BRL", "USD", "EUR", "ARS", "CLP", "COP", "MXN", "PEN", "UYU")),
	)
}

// UpdateTransactionRequest - representa a requisição para atualizar uma transação
type UpdateTransactionRequest struct {
	Status       string  `json:"status"`
	StatusDetail string  `json:"status_detail"`
	Amount       float64 `json:"amount"`
}

// Validate - valida os campos da requisição
func (r *UpdateTransactionRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Status, validation.In("pending", "approved", "authorized", "rejected", "cancelled", "refunded")),
		validation.Field(&r.StatusDetail, validation.Length(0, 255)),
		validation.Field(&r.Amount, validation.Min(0.01)),
	)
}

// ProcessOrderRequest - representa a requisição para processar uma ordem
type ProcessOrderRequest struct {
	PaymentMethodID string                   `json:"payment_method_id"`
	Token           string                   `json:"token"`
	Installments    int                      `json:"installments"`
	IssuerID        string                   `json:"issuer_id"`
	Payer           ProcessOrderPayerRequest `json:"payer"`
	Metadata        map[string]string        `json:"metadata"`
}

// Validate - valida os campos da requisição
func (r *ProcessOrderRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.PaymentMethodID, validation.Required, validation.Length(1, 100)),
		validation.Field(&r.Token, validation.Length(0, 500)),
		validation.Field(&r.Installments, validation.Min(1), validation.Max(12)),
	)
}

// ProcessOrderPayerRequest - representa os dados do pagador no processamento
type ProcessOrderPayerRequest struct {
	Email          string                            `json:"email"`
	Identification ProcessOrderIdentificationRequest `json:"identification"`
	FirstName      string                            `json:"first_name"`
	LastName       string                            `json:"last_name"`
}

// ProcessOrderIdentificationRequest - representa a identificação do pagador
type ProcessOrderIdentificationRequest struct {
	Type   string `json:"type"`
	Number string `json:"number"`
}
