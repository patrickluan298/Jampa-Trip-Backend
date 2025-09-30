package mercadopago

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/jampa_trip/internal/pkg/util"
)

// Client - representa o cliente HTTP para comunicação com a API do Mercado Pago
type Client struct {
	AccessToken string
	BaseURL     string
	HTTPClient  *http.Client
}

// NewClient - cria uma nova instância do cliente Mercado Pago
func NewClient(accessToken, baseURL string) *Client {
	return &Client{
		AccessToken: accessToken,
		BaseURL:     baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// OrderRequest - representa a estrutura para criar uma order (v1/orders)
type OrderRequest struct {
	ExternalReference string            `json:"external_reference"`
	TotalAmount       float64           `json:"total_amount"`
	Items             []OrderItem       `json:"items"`
	Payer             Payer             `json:"payer"`
	NotificationURL   string            `json:"notification_url,omitempty"`
	Description       string            `json:"description,omitempty"`
	Metadata          map[string]string `json:"metadata,omitempty"`
}

// OrderItem - representa um item da order
type OrderItem struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	PictureURL  string  `json:"picture_url,omitempty"`
	CategoryID  string  `json:"category_id,omitempty"`
	Quantity    int     `json:"quantity"`
	CurrencyID  string  `json:"currency_id"`
	UnitPrice   float64 `json:"unit_price"`
}

// Payer - representa o pagador
type Payer struct {
	Name    string  `json:"name,omitempty"`
	Email   string  `json:"email,omitempty"`
	Phone   Phone   `json:"phone,omitempty"`
	Address Address `json:"address,omitempty"`
}

// Phone - representa o telefone do pagador
type Phone struct {
	AreaCode string `json:"area_code,omitempty"`
	Number   string `json:"number,omitempty"`
}

// Address - representa o endereço do pagador
type Address struct {
	StreetName   string `json:"street_name,omitempty"`
	StreetNumber int    `json:"street_number,omitempty"`
	ZipCode      string `json:"zip_code,omitempty"`
}

// OrderResponse - representa a resposta da criação de uma order
type OrderResponse struct {
	ID                string            `json:"id"`
	ExternalReference string            `json:"external_reference"`
	TotalAmount       float64           `json:"total_amount"`
	Status            string            `json:"status"`
	StatusDetail      string            `json:"status_detail"`
	Items             []OrderItem       `json:"items"`
	Payer             Payer             `json:"payer"`
	NotificationURL   string            `json:"notification_url,omitempty"`
	Description       string            `json:"description,omitempty"`
	Metadata          map[string]string `json:"metadata,omitempty"`
	DateCreated       string            `json:"date_created"`
	DateLastUpdated   string            `json:"date_last_updated"`
	DateExpiration    string            `json:"date_expiration,omitempty"`
	PaymentMethods    []PaymentMethod   `json:"payment_methods,omitempty"`
}

// PaymentMethod - representa um método de pagamento disponível
type PaymentMethod struct {
	ID               string `json:"id"`
	Type             string `json:"type"`
	Issuer           string `json:"issuer,omitempty"`
	Installments     int    `json:"installments,omitempty"`
	MinInstallments  int    `json:"min_installments,omitempty"`
	MaxInstallments  int    `json:"max_installments,omitempty"`
	SecurityCode     bool   `json:"security_code,omitempty"`
	SecurityCodeMode string `json:"security_code_mode,omitempty"`
}

// PaymentRequest - representa a estrutura para criar um pagamento
type PaymentRequest struct {
	TransactionAmount float64           `json:"transaction_amount"`
	Description       string            `json:"description"`
	PaymentMethodID   string            `json:"payment_method_id"`
	Payer             PaymentPayer      `json:"payer"`
	Metadata          map[string]string `json:"metadata,omitempty"`
}

// PaymentPayer - representa o pagador para pagamentos
type PaymentPayer struct {
	Email      string  `json:"email,omitempty"`
	FirstName  string  `json:"first_name,omitempty"`
	LastName   string  `json:"last_name,omitempty"`
	Phone      Phone   `json:"phone,omitempty"`
	Address    Address `json:"address,omitempty"`
	EntityType string  `json:"entity_type,omitempty"`
	Type       string  `json:"type,omitempty"`
	ID         string  `json:"id,omitempty"`
}

// PaymentResponse - representa a resposta da criação de um pagamento
type PaymentResponse struct {
	ID                int64             `json:"id"`
	Status            string            `json:"status"`
	StatusDetail      string            `json:"status_detail"`
	TransactionAmount float64           `json:"transaction_amount"`
	Description       string            `json:"description"`
	PaymentMethodID   string            `json:"payment_method_id"`
	Payer             PaymentPayer      `json:"payer"`
	DateCreated       string            `json:"date_created"`
	DateApproved      string            `json:"date_approved,omitempty"`
	DateLastUpdated   string            `json:"date_last_updated"`
	Metadata          map[string]string `json:"metadata,omitempty"`
}

// PIXRequest - representa a estrutura para criar um pagamento PIX
type PIXRequest struct {
	TransactionAmount float64           `json:"transaction_amount"`
	Description       string            `json:"description"`
	PaymentMethodID   string            `json:"payment_method_id"`
	Payer             PaymentPayer      `json:"payer"`
	Metadata          map[string]string `json:"metadata,omitempty"`
}

// PIXResponse - representa a resposta da criação de um pagamento PIX
type PIXResponse struct {
	ID                 int64              `json:"id"`
	Status             string             `json:"status"`
	StatusDetail       string             `json:"status_detail"`
	TransactionAmount  float64            `json:"transaction_amount"`
	Description        string             `json:"description"`
	PaymentMethodID    string             `json:"payment_method_id"`
	Payer              PaymentPayer       `json:"payer"`
	DateCreated        string             `json:"date_created"`
	DateApproved       string             `json:"date_approved,omitempty"`
	DateLastUpdated    string             `json:"date_last_updated"`
	PointOfInteraction PointOfInteraction `json:"point_of_interaction,omitempty"`
	Metadata           map[string]string  `json:"metadata,omitempty"`
}

// PointOfInteraction - representa informações do PIX
type PointOfInteraction struct {
	Type            string          `json:"type"`
	SubType         string          `json:"sub_type,omitempty"`
	ApplicationData ApplicationData `json:"application_data,omitempty"`
}

// ApplicationData - representa dados da aplicação PIX
type ApplicationData struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}

// ErrorResponse - representa uma resposta de erro da API
type ErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
	Status  int    `json:"status"`
	Cause   []struct {
		Code        string `json:"code"`
		Description string `json:"description"`
	} `json:"cause,omitempty"`
}

// CreateOrder - cria uma nova order no Mercado Pago
func (c *Client) CreateOrder(orderReq *OrderRequest) (*OrderResponse, error) {
	url := fmt.Sprintf("%s/v1/orders", c.BaseURL)

	jsonData, err := json.Marshal(orderReq)
	if err != nil {
		return nil, util.WrapError("erro ao serializar order", err, http.StatusInternalServerError)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, util.WrapError("erro ao criar requisição", err, http.StatusInternalServerError)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, util.WrapError("erro ao executar requisição", err, http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, util.WrapError("erro ao ler resposta", err, http.StatusInternalServerError)
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		var errorResp ErrorResponse
		if err := json.Unmarshal(body, &errorResp); err != nil {
			return nil, util.WrapError(fmt.Sprintf("erro na API do Mercado Pago (status %d): %s", resp.StatusCode, string(body)), err, resp.StatusCode)
		}
		return nil, util.WrapError(fmt.Sprintf("erro na API do Mercado Pago: %s", errorResp.Message), nil, resp.StatusCode)
	}

	var orderResp OrderResponse
	if err := json.Unmarshal(body, &orderResp); err != nil {
		return nil, util.WrapError("erro ao deserializar resposta", err, http.StatusInternalServerError)
	}

	return &orderResp, nil
}

// GetOrder - obtém informações de uma order específica
func (c *Client) GetOrder(orderID string) (*OrderResponse, error) {
	url := fmt.Sprintf("%s/v1/orders/%s", c.BaseURL, orderID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, util.WrapError("erro ao criar requisição", err, http.StatusInternalServerError)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, util.WrapError("erro ao executar requisição", err, http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, util.WrapError("erro ao ler resposta", err, http.StatusInternalServerError)
	}

	if resp.StatusCode != http.StatusOK {
		var errorResp ErrorResponse
		if err := json.Unmarshal(body, &errorResp); err != nil {
			return nil, util.WrapError(fmt.Sprintf("erro na API do Mercado Pago (status %d): %s", resp.StatusCode, string(body)), err, resp.StatusCode)
		}
		return nil, util.WrapError(fmt.Sprintf("erro na API do Mercado Pago: %s", errorResp.Message), nil, resp.StatusCode)
	}

	var orderResp OrderResponse
	if err := json.Unmarshal(body, &orderResp); err != nil {
		return nil, util.WrapError("erro ao deserializar resposta", err, http.StatusInternalServerError)
	}

	return &orderResp, nil
}

// CancelOrder - cancela uma order
func (c *Client) CancelOrder(orderID string) (*OrderResponse, error) {
	url := fmt.Sprintf("%s/v1/orders/%s/cancel", c.BaseURL, orderID)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, util.WrapError("erro ao criar requisição", err, http.StatusInternalServerError)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, util.WrapError("erro ao executar requisição", err, http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, util.WrapError("erro ao ler resposta", err, http.StatusInternalServerError)
	}

	if resp.StatusCode != http.StatusOK {
		var errorResp ErrorResponse
		if err := json.Unmarshal(body, &errorResp); err != nil {
			return nil, util.WrapError(fmt.Sprintf("erro na API do Mercado Pago (status %d): %s", resp.StatusCode, string(body)), err, resp.StatusCode)
		}
		return nil, util.WrapError(fmt.Sprintf("erro na API do Mercado Pago: %s", errorResp.Message), nil, resp.StatusCode)
	}

	var orderResp OrderResponse
	if err := json.Unmarshal(body, &orderResp); err != nil {
		return nil, util.WrapError("erro ao deserializar resposta", err, http.StatusInternalServerError)
	}

	return &orderResp, nil
}

// CaptureOrder - captura uma order totalmente
func (c *Client) CaptureOrder(orderID string) (*OrderResponse, error) {
	url := fmt.Sprintf("%s/v1/orders/%s/capture", c.BaseURL, orderID)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, util.WrapError("erro ao criar requisição", err, http.StatusInternalServerError)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, util.WrapError("erro ao executar requisição", err, http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, util.WrapError("erro ao ler resposta", err, http.StatusInternalServerError)
	}

	if resp.StatusCode != http.StatusOK {
		var errorResp ErrorResponse
		if err := json.Unmarshal(body, &errorResp); err != nil {
			return nil, util.WrapError(fmt.Sprintf("erro na API do Mercado Pago (status %d): %s", resp.StatusCode, string(body)), err, resp.StatusCode)
		}
		return nil, util.WrapError(fmt.Sprintf("erro na API do Mercado Pago: %s", errorResp.Message), nil, resp.StatusCode)
	}

	var orderResp OrderResponse
	if err := json.Unmarshal(body, &orderResp); err != nil {
		return nil, util.WrapError("erro ao deserializar resposta", err, http.StatusInternalServerError)
	}

	return &orderResp, nil
}

// RefundOrder - reembolsa uma order
func (c *Client) RefundOrder(orderID string) (*OrderResponse, error) {
	url := fmt.Sprintf("%s/v1/orders/%s/refund", c.BaseURL, orderID)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, util.WrapError("erro ao criar requisição", err, http.StatusInternalServerError)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, util.WrapError("erro ao executar requisição", err, http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, util.WrapError("erro ao ler resposta", err, http.StatusInternalServerError)
	}

	if resp.StatusCode != http.StatusOK {
		var errorResp ErrorResponse
		if err := json.Unmarshal(body, &errorResp); err != nil {
			return nil, util.WrapError(fmt.Sprintf("erro na API do Mercado Pago (status %d): %s", resp.StatusCode, string(body)), err, resp.StatusCode)
		}
		return nil, util.WrapError(fmt.Sprintf("erro na API do Mercado Pago: %s", errorResp.Message), nil, resp.StatusCode)
	}

	var orderResp OrderResponse
	if err := json.Unmarshal(body, &orderResp); err != nil {
		return nil, util.WrapError("erro ao deserializar resposta", err, http.StatusInternalServerError)
	}

	return &orderResp, nil
}

// CreatePayment - cria um novo pagamento no Mercado Pago
func (c *Client) CreatePayment(paymentReq *PaymentRequest) (*PaymentResponse, error) {
	url := fmt.Sprintf("%s/v1/payments", c.BaseURL)

	jsonData, err := json.Marshal(paymentReq)
	if err != nil {
		return nil, util.WrapError("erro ao serializar pagamento", err, http.StatusInternalServerError)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, util.WrapError("erro ao criar requisição", err, http.StatusInternalServerError)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, util.WrapError("erro ao executar requisição", err, http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, util.WrapError("erro ao ler resposta", err, http.StatusInternalServerError)
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		var errorResp ErrorResponse
		if err := json.Unmarshal(body, &errorResp); err != nil {
			return nil, util.WrapError(fmt.Sprintf("erro na API do Mercado Pago (status %d): %s", resp.StatusCode, string(body)), err, resp.StatusCode)
		}
		return nil, util.WrapError(fmt.Sprintf("erro na API do Mercado Pago: %s", errorResp.Message), nil, resp.StatusCode)
	}

	var paymentResp PaymentResponse
	if err := json.Unmarshal(body, &paymentResp); err != nil {
		return nil, util.WrapError("erro ao deserializar resposta", err, http.StatusInternalServerError)
	}

	return &paymentResp, nil
}

// CreatePIXPayment - cria um novo pagamento PIX no Mercado Pago
func (c *Client) CreatePIXPayment(pixReq *PIXRequest) (*PIXResponse, error) {
	url := fmt.Sprintf("%s/v1/payments", c.BaseURL)

	jsonData, err := json.Marshal(pixReq)
	if err != nil {
		return nil, util.WrapError("erro ao serializar pagamento PIX", err, http.StatusInternalServerError)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, util.WrapError("erro ao criar requisição", err, http.StatusInternalServerError)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, util.WrapError("erro ao executar requisição", err, http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, util.WrapError("erro ao ler resposta", err, http.StatusInternalServerError)
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		var errorResp ErrorResponse
		if err := json.Unmarshal(body, &errorResp); err != nil {
			return nil, util.WrapError(fmt.Sprintf("erro na API do Mercado Pago (status %d): %s", resp.StatusCode, string(body)), err, resp.StatusCode)
		}
		return nil, util.WrapError(fmt.Sprintf("erro na API do Mercado Pago: %s", errorResp.Message), nil, resp.StatusCode)
	}

	var pixResp PIXResponse
	if err := json.Unmarshal(body, &pixResp); err != nil {
		return nil, util.WrapError("erro ao deserializar resposta", err, http.StatusInternalServerError)
	}

	return &pixResp, nil
}

// GetPayment - obtém informações de um pagamento específico
func (c *Client) GetPayment(paymentID string) (*PaymentResponse, error) {
	url := fmt.Sprintf("%s/v1/payments/%s", c.BaseURL, paymentID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, util.WrapError("erro ao criar requisição", err, http.StatusInternalServerError)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, util.WrapError("erro ao executar requisição", err, http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, util.WrapError("erro ao ler resposta", err, http.StatusInternalServerError)
	}

	if resp.StatusCode != http.StatusOK {
		var errorResp ErrorResponse
		if err := json.Unmarshal(body, &errorResp); err != nil {
			return nil, util.WrapError(fmt.Sprintf("erro na API do Mercado Pago (status %d): %s", resp.StatusCode, string(body)), err, resp.StatusCode)
		}
		return nil, util.WrapError(fmt.Sprintf("erro na API do Mercado Pago: %s", errorResp.Message), nil, resp.StatusCode)
	}

	var paymentResp PaymentResponse
	if err := json.Unmarshal(body, &paymentResp); err != nil {
		return nil, util.WrapError("erro ao deserializar resposta", err, http.StatusInternalServerError)
	}

	return &paymentResp, nil
}

// CancelPayment - cancela um pagamento
func (c *Client) CancelPayment(paymentID string) (*PaymentResponse, error) {
	url := fmt.Sprintf("%s/v1/payments/%s", c.BaseURL, paymentID)

	cancelData := map[string]string{
		"status": "cancelled",
	}

	jsonData, err := json.Marshal(cancelData)
	if err != nil {
		return nil, util.WrapError("erro ao serializar dados de cancelamento", err, http.StatusInternalServerError)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, util.WrapError("erro ao criar requisição", err, http.StatusInternalServerError)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, util.WrapError("erro ao executar requisição", err, http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, util.WrapError("erro ao ler resposta", err, http.StatusInternalServerError)
	}

	if resp.StatusCode != http.StatusOK {
		var errorResp ErrorResponse
		if err := json.Unmarshal(body, &errorResp); err != nil {
			return nil, util.WrapError(fmt.Sprintf("erro na API do Mercado Pago (status %d): %s", resp.StatusCode, string(body)), err, resp.StatusCode)
		}
		return nil, util.WrapError(fmt.Sprintf("erro na API do Mercado Pago: %s", errorResp.Message), nil, resp.StatusCode)
	}

	var paymentResp PaymentResponse
	if err := json.Unmarshal(body, &paymentResp); err != nil {
		return nil, util.WrapError("erro ao deserializar resposta", err, http.StatusInternalServerError)
	}

	return &paymentResp, nil
}
