package contract

// CartaoResponse - representa a resposta de um cartão
type CartaoResponse struct {
	ID              string                  `json:"id"`
	CustomerID      string                  `json:"customer_id"`
	FirstSixDigits  string                  `json:"first_six_digits"`
	LastFourDigits  string                  `json:"last_four_digits"`
	ExpirationMonth int                     `json:"expiration_month"`
	ExpirationYear  int                     `json:"expiration_year"`
	SecurityCode    SecurityCodeInfo        `json:"security_code"`
	Issuer          IssuerInfo              `json:"issuer"`
	PaymentMethod   CartaoPaymentMethodInfo `json:"payment_method"`
	Cardholder      CardholderResponse      `json:"cardholder"`
	DateCreated     string                  `json:"date_created"`
	DateLastUpdated string                  `json:"date_last_updated"`
	Metadata        map[string]string       `json:"metadata,omitempty"`
}

// SecurityCodeInfo - representa informações do código de segurança
type SecurityCodeInfo struct {
	Length       int    `json:"length"`
	CardLocation string `json:"card_location"`
	Mode         string `json:"mode"`
}

// IssuerInfo - representa informações do emissor
type IssuerInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// CartaoPaymentMethodInfo - representa informações do método de pagamento
type CartaoPaymentMethodInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// CardholderResponse - representa o portador do cartão na resposta
type CardholderResponse struct {
	Name           string                 `json:"name"`
	Identification IdentificationResponse `json:"identification"`
}

// IdentificationResponse - representa a identificação do portador na resposta
type IdentificationResponse struct {
	Type   string `json:"type"`
	Number string `json:"number"`
}

// ListCartoesResponse - representa a resposta da listagem de cartões
type ListCartoesResponse struct {
	Cartoes []CartaoResponse `json:"cartoes"`
	Total   int              `json:"total"`
}

// CreateCartaoResponse - representa a resposta da criação de cartão
type CreateCartaoResponse struct {
	Cartao  CartaoResponse `json:"cartao"`
	Message string         `json:"message"`
}

// UpdateCartaoResponse - representa a resposta da atualização de cartão
type UpdateCartaoResponse struct {
	Cartao  CartaoResponse `json:"cartao"`
	Message string         `json:"message"`
}

// DeleteCartaoResponse - representa a resposta da exclusão de cartão
type DeleteCartaoResponse struct {
	Message string `json:"message"`
}
