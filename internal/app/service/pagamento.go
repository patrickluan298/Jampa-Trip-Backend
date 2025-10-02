package service

import (
	"context"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/jampa_trip/internal/app/contract"
	"github.com/jampa_trip/internal/app/mercadopago"
	"github.com/jampa_trip/internal/app/model"
	"github.com/jampa_trip/internal/app/repository"
	"github.com/jampa_trip/internal/pkg/config"
	"github.com/jampa_trip/internal/pkg/util"
	"gorm.io/gorm"
)

// PagamentoService - objeto de contexto
type PagamentoService struct {
	PagamentoRepository *repository.PagamentoRepository
	MPClient            *mercadopago.Client
	paymentMethodsCache *paymentMethodsCache
}

// paymentMethodsCache - cache para meios de pagamento
type paymentMethodsCache struct {
	methods    *mercadopago.PaymentMethodsResponse
	expiration time.Time
	mu         sync.RWMutex
	ttl        time.Duration
}

// PagamentoServiceNew - construtor do objeto
func PagamentoServiceNew(DB *gorm.DB) *PagamentoService {
	cfg, _ := config.LoadConfig()

	return &PagamentoService{
		PagamentoRepository: repository.PagamentoRepositoryNew(DB),
		MPClient:            mercadopago.NewClient(cfg.MercadoPagoAccessToken, cfg.MercadoPagoBaseURL),
		paymentMethodsCache: &paymentMethodsCache{
			ttl: 10 * time.Minute,
		},
	}
}

// AutorizarCartao - autoriza um pagamento com cartão de crédito (pré-autorização)
func (s *PagamentoService) AutorizarCartao(ctx context.Context, req *contract.AutorizarCartaoRequest) (*contract.AutorizarCartaoResponse, error) {

	if req.Installments < 1 || req.Installments > 12 {
		return nil, util.WrapError("número de parcelas deve estar entre 1 e 12", nil, http.StatusBadRequest)
	}

	mpReq := &mercadopago.CreditCardPaymentRequest{
		TransactionAmount: req.TransactionAmount,
		Token:             req.Token,
		Description:       req.Description,
		Installments:      req.Installments,
		PaymentMethodID:   req.PaymentMethodID,
		IssuerID:          req.IssuerID,
		Capture:           false,
		ExternalReference: req.ExternalReference,
		Payer: mercadopago.CreditCardPayer{
			Email: req.Payer.Email,
			Identification: mercadopago.CreditCardIdentification{
				Type:   req.Payer.Identification.Type,
				Number: req.Payer.Identification.Number,
			},
			FirstName: req.Payer.FirstName,
			LastName:  req.Payer.LastName,
		},
		Metadata: map[string]string{
			"cliente_id": strconv.Itoa(req.ClienteID),
			"empresa_id": strconv.Itoa(req.EmpresaID),
		},
	}

	mpResp, err := s.MPClient.CreateCreditCardPayment(ctx, mpReq)
	if err != nil {
		return nil, s.mapMercadoPagoError(err)
	}

	statusMessage := s.getStatusDetailMessage(mpResp.StatusDetail)

	now := time.Now()
	pagamento := &model.Pagamento{
		ClienteID:            req.ClienteID,
		EmpresaID:            req.EmpresaID,
		MercadoPagoPaymentID: strconv.FormatInt(mpResp.ID, 10),
		Status:               mpResp.Status,
		StatusDetail:         mpResp.StatusDetail,
		Valor:                mpResp.TransactionAmount,
		Moeda:                mpResp.CurrencyID,
		MetodoPagamento:      "credit_card",
		Descricao:            mpResp.Description,
		NumeroParcelas:       mpResp.Installments,
		PaymentMethodID:      mpResp.PaymentMethodID,
		IssuerID:             mpResp.IssuerID,
		LastFourDigits:       mpResp.Card.LastFourDigits,
		FirstSixDigits:       mpResp.Card.FirstSixDigits,
		CardholderName:       mpResp.Card.Cardholder.Name,
		Captured:             mpResp.Captured,
		MomentoCriacao:       now,
		MomentoAtualizacao:   now,
	}

	if mpResp.Status == "authorized" {
		pagamento.MomentoAutorizacao = &now
	}

	if err := s.PagamentoRepository.Create(pagamento); err != nil {
		return nil, util.WrapError("erro ao salvar pagamento", err, http.StatusInternalServerError)
	}

	return &contract.AutorizarCartaoResponse{
		Pagamento: s.modelToResponse(pagamento),
		Message:   statusMessage,
	}, nil
}

// CapturarPagamento - captura um pagamento autorizado
func (s *PagamentoService) CapturarPagamento(ctx context.Context, req *contract.CapturarPagamentoRequest) (*contract.CapturarPagamentoResponse, error) {

	if err := req.Validate(); err != nil {
		return nil, util.WrapError("erro de validação", err, http.StatusBadRequest)
	}

	paymentIDStr := strconv.FormatInt(req.PaymentID, 10)
	pagamento, err := s.PagamentoRepository.GetByMercadoPagoPaymentID(paymentIDStr)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, util.WrapError("pagamento não encontrado", err, http.StatusNotFound)
		}
		return nil, util.WrapError("erro ao buscar pagamento", err, http.StatusInternalServerError)
	}

	if !pagamento.CanBeCaptured() {
		return nil, util.WrapError("pagamento não pode ser capturado (status: "+pagamento.Status+")", nil, http.StatusBadRequest)
	}

	captureReq := &mercadopago.CapturePaymentRequest{}
	if req.TransactionAmount != nil {
		captureReq.TransactionAmount = req.TransactionAmount
	}

	mpResp, err := s.MPClient.CapturePayment(ctx, req.PaymentID, captureReq)
	if err != nil {
		return nil, s.mapMercadoPagoError(err)
	}

	now := time.Now()
	pagamento.Status = mpResp.Status
	pagamento.StatusDetail = mpResp.StatusDetail
	pagamento.Captured = mpResp.Captured
	pagamento.MomentoAtualizacao = now

	if mpResp.Status == "approved" && mpResp.Captured {
		pagamento.MomentoCaptura = &now
		pagamento.MomentoAprovacao = &now
	}

	if err := s.PagamentoRepository.Update(pagamento); err != nil {
		return nil, util.WrapError("erro ao atualizar pagamento", err, http.StatusInternalServerError)
	}

	return &contract.CapturarPagamentoResponse{
		Pagamento: s.modelToResponse(pagamento),
		Message:   "Pagamento capturado com sucesso",
	}, nil
}

// CancelarPagamento - cancela um pagamento autorizado (void)
func (s *PagamentoService) CancelarPagamento(ctx context.Context, req *contract.CancelarPagamentoRequest) (*contract.CancelarPagamentoResponse, error) {

	if err := req.Validate(); err != nil {
		return nil, util.WrapError("erro de validação", err, http.StatusBadRequest)
	}

	paymentIDStr := strconv.FormatInt(req.PaymentID, 10)
	pagamento, err := s.PagamentoRepository.GetByMercadoPagoPaymentID(paymentIDStr)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, util.WrapError("pagamento não encontrado", err, http.StatusNotFound)
		}
		return nil, util.WrapError("erro ao buscar pagamento", err, http.StatusInternalServerError)
	}

	// if !pagamento.CanBeCancelled() {
	// 	return nil, util.WrapError("pagamento não pode ser cancelado (status: "+pagamento.Status+")", nil, http.StatusBadRequest)
	// }

	mpResp, err := s.MPClient.CancelCreditCardPayment(ctx, req.PaymentID)
	if err != nil {
		return nil, s.mapMercadoPagoError(err)
	}

	now := time.Now()
	pagamento.Status = mpResp.Status
	pagamento.StatusDetail = mpResp.StatusDetail
	pagamento.MomentoCancelamento = &now
	pagamento.MomentoAtualizacao = now

	if err := s.PagamentoRepository.Update(pagamento); err != nil {
		return nil, util.WrapError("erro ao atualizar pagamento", err, http.StatusInternalServerError)
	}

	return &contract.CancelarPagamentoResponse{
		Pagamento: s.modelToResponse(pagamento),
		Message:   "Pagamento cancelado com sucesso",
	}, nil
}

// ReembolsarPagamento - reembolsa um pagamento capturado
func (s *PagamentoService) ReembolsarPagamento(ctx context.Context, req *contract.ReembolsarPagamentoRequest) (*contract.ReembolsarPagamentoResponse, error) {

	if err := req.Validate(); err != nil {
		return nil, util.WrapError("erro de validação", err, http.StatusBadRequest)
	}

	paymentIDStr := strconv.FormatInt(req.PaymentID, 10)
	pagamento, err := s.PagamentoRepository.GetByMercadoPagoPaymentID(paymentIDStr)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, util.WrapError("pagamento não encontrado", err, http.StatusNotFound)
		}
		return nil, util.WrapError("erro ao buscar pagamento", err, http.StatusInternalServerError)
	}

	if !pagamento.CanBeRefunded() {
		return nil, util.WrapError("pagamento não pode ser reembolsado (status: "+pagamento.Status+")", nil, http.StatusBadRequest)
	}

	refundReq := &mercadopago.RefundPaymentRequest{}
	if req.Amount != nil {
		refundReq.Amount = req.Amount
	}

	refundResp, err := s.MPClient.RefundCreditCardPayment(ctx, req.PaymentID, refundReq)
	if err != nil {
		return nil, s.mapMercadoPagoError(err)
	}

	now := time.Now()
	pagamento.Status = "refunded"
	pagamento.StatusDetail = "refunded"
	pagamento.TransactionAmountRefunded = refundResp.Amount
	pagamento.MomentoAtualizacao = now

	if err := s.PagamentoRepository.Update(pagamento); err != nil {
		return nil, util.WrapError("erro ao atualizar pagamento", err, http.StatusInternalServerError)
	}

	return &contract.ReembolsarPagamentoResponse{
		Pagamento:      s.modelToResponse(pagamento),
		RefundID:       refundResp.ID,
		RefundedAmount: refundResp.Amount,
		Message:        "Pagamento reembolsado com sucesso",
	}, nil
}

// ObterPagamento - obtém informações de um pagamento
func (s *PagamentoService) ObterPagamento(ctx context.Context, paymentID int64) (*contract.PagamentoResponse, error) {

	paymentIDStr := strconv.FormatInt(paymentID, 10)
	pagamento, err := s.PagamentoRepository.GetByMercadoPagoPaymentID(paymentIDStr)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, util.WrapError("pagamento não encontrado", err, http.StatusNotFound)
		}
		return nil, util.WrapError("erro ao buscar pagamento", err, http.StatusInternalServerError)
	}

	mpResp, err := s.MPClient.GetCreditCardPayment(ctx, paymentID)
	if err == nil {
		if mpResp.Status != pagamento.Status || mpResp.StatusDetail != pagamento.StatusDetail {
			pagamento.Status = mpResp.Status
			pagamento.StatusDetail = mpResp.StatusDetail
			pagamento.Captured = mpResp.Captured
			pagamento.TransactionAmountRefunded = mpResp.TransactionAmountRefunded
			pagamento.MomentoAtualizacao = time.Now()
			s.PagamentoRepository.Update(pagamento)
		}
	}

	response := s.modelToResponse(pagamento)
	return &response, nil
}

// ListarMeiosPagamento - lista os meios de pagamento disponíveis
func (s *PagamentoService) ListarMeiosPagamento(ctx context.Context) (*contract.PaymentMethodsResponse, error) {
	s.paymentMethodsCache.mu.RLock()
	if s.paymentMethodsCache.methods != nil && time.Now().Before(s.paymentMethodsCache.expiration) {
		methods := s.paymentMethodsCache.methods
		s.paymentMethodsCache.mu.RUnlock()
		return s.mapPaymentMethodsResponse(methods), nil
	}
	s.paymentMethodsCache.mu.RUnlock()

	methods, err := s.MPClient.GetPaymentMethods(ctx)
	if err != nil {
		return nil, s.mapMercadoPagoError(err)
	}

	s.paymentMethodsCache.mu.Lock()
	s.paymentMethodsCache.methods = methods
	s.paymentMethodsCache.expiration = time.Now().Add(s.paymentMethodsCache.ttl)
	s.paymentMethodsCache.mu.Unlock()

	return s.mapPaymentMethodsResponse(methods), nil
}

// modelToResponse - converte model para response
func (s *PagamentoService) modelToResponse(p *model.Pagamento) contract.PagamentoResponse {
	return contract.PagamentoResponse{
		ID:                        p.ID,
		ClienteID:                 p.ClienteID,
		EmpresaID:                 p.EmpresaID,
		MercadoPagoOrderID:        p.MercadoPagoOrderID,
		MercadoPagoPaymentID:      p.MercadoPagoPaymentID,
		Status:                    p.Status,
		StatusDetail:              p.StatusDetail,
		Valor:                     p.Valor,
		Moeda:                     p.Moeda,
		MetodoPagamento:           p.MetodoPagamento,
		Descricao:                 p.Descricao,
		NumeroParcelas:            p.NumeroParcelas,
		LastFourDigits:            p.LastFourDigits,
		FirstSixDigits:            p.FirstSixDigits,
		PaymentMethodID:           p.PaymentMethodID,
		IssuerID:                  p.IssuerID,
		CardholderName:            p.CardholderName,
		Captured:                  p.Captured,
		TransactionAmountRefunded: p.TransactionAmountRefunded,
		MomentoCriacao:            p.MomentoCriacao,
		MomentoAtualizacao:        p.MomentoAtualizacao,
		MomentoAprovacao:          p.MomentoAprovacao,
		MomentoCancelamento:       p.MomentoCancelamento,
		MomentoAutorizacao:        p.MomentoAutorizacao,
		MomentoCaptura:            p.MomentoCaptura,
		StatusDisplay:             p.GetStatusDisplay(),
		MetodoPagamentoDisplay:    p.GetMetodoPagamentoDisplay(),
	}
}

// mapPaymentMethodsResponse - mapeia resposta de meios de pagamento
func (s *PagamentoService) mapPaymentMethodsResponse(methods *mercadopago.PaymentMethodsResponse) *contract.PaymentMethodsResponse {
	var result []contract.PaymentMethodInfo

	for _, method := range *methods {
		if method.PaymentTypeID == "credit_card" && method.Status == "active" {
			result = append(result, contract.PaymentMethodInfo{
				ID:                method.ID,
				Name:              method.Name,
				PaymentTypeID:     method.PaymentTypeID,
				Status:            method.Status,
				MinAllowedAmount:  method.MinAllowedAmount,
				MaxAllowedAmount:  method.MaxAllowedAmount,
				AccreditationTime: method.AccreditationTime,
			})
		}
	}

	return &contract.PaymentMethodsResponse{
		Methods: result,
	}
}

// mapMercadoPagoError - mapeia erros do Mercado Pago para mensagens amigáveis
func (s *PagamentoService) mapMercadoPagoError(err error) error {
	return err
}

// getStatusDetailMessage - retorna mensagens amigáveis para status_detail
func (s *PagamentoService) getStatusDetailMessage(statusDetail string) string {
	messages := map[string]string{
		"accredited":                           "Pagamento aprovado",
		"pending_contingency":                  "Aguardando confirmação",
		"pending_review_manual":                "Em revisão manual",
		"cc_rejected_bad_filled_card_number":   "Número do cartão inválido",
		"cc_rejected_bad_filled_date":          "Data de validade inválida",
		"cc_rejected_bad_filled_other":         "Dados do cartão inválidos",
		"cc_rejected_bad_filled_security_code": "Código de segurança inválido",
		"cc_rejected_blacklist":                "Cartão bloqueado",
		"cc_rejected_call_for_authorize":       "Autorização necessária - contate o banco",
		"cc_rejected_card_disabled":            "Cartão desabilitado",
		"cc_rejected_duplicated_payment":       "Pagamento duplicado",
		"cc_rejected_high_risk":                "Pagamento recusado por risco",
		"cc_rejected_insufficient_amount":      "Saldo insuficiente",
		"cc_rejected_invalid_installments":     "Número de parcelas inválido",
		"cc_rejected_max_attempts":             "Excedido número máximo de tentativas",
		"cc_rejected_other_reason":             "Pagamento recusado - contate o banco",
	}

	if msg, ok := messages[statusDetail]; ok {
		return msg
	}
	return "Status: " + statusDetail
}
