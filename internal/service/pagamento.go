package service

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/jampa_trip/internal/contract"
	"github.com/jampa_trip/internal/model"
	"github.com/jampa_trip/internal/repository"
	"github.com/jampa_trip/pkg/config"
	"github.com/jampa_trip/pkg/mercadopago"
	"github.com/jampa_trip/pkg/util"
	"gorm.io/gorm"
)

// PagamentoService - objeto de contexto
type PagamentoService struct {
	PagamentoRepository *repository.PagamentoRepository
	MPClient            *mercadopago.Client
}

// PagamentoServiceNew - construtor do objeto
func PagamentoServiceNew(DB *gorm.DB) *PagamentoService {
	cfg, _ := config.LoadConfig()

	return &PagamentoService{
		PagamentoRepository: repository.PagamentoRepositoryNew(DB),
		MPClient:            mercadopago.NewClient(cfg.MercadoPagoAccessToken, cfg.MercadoPagoBaseURL),
	}
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

// CriarPagamentoCartaoCredito - cria um pagamento com cartão de crédito
func (s *PagamentoService) CriarPagamentoCartaoCredito(ctx context.Context, req *contract.CriarPagamentoCartaoCreditoRequest) (*contract.CriarPagamentoCartaoCreditoResponse, error) {

	if err := req.Validate(); err != nil {
		return nil, util.WrapError("erro de validação", err, http.StatusBadRequest)
	}

	mpReq := &mercadopago.CreditCardPaymentRequest{
		TransactionAmount: req.TransactionAmount,
		Token:             req.Token,
		Description:       req.Description,
		Installments:      req.Installments,
		PaymentMethodID:   req.PaymentMethodID,
		IssuerID:          req.IssuerID,
		Capture:           req.Capture,
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
		return nil, err
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

	if mpResp.Status == "approved" {
		pagamento.MomentoAprovacao = &now
	}

	if mpResp.Status == "authorized" {
		pagamento.MomentoAutorizacao = &now
	}

	if err := s.PagamentoRepository.Create(pagamento); err != nil {
		return nil, util.WrapError("erro ao salvar pagamento", err, http.StatusInternalServerError)
	}

	return &contract.CriarPagamentoCartaoCreditoResponse{
		Pagamento: s.modelToResponse(pagamento),
		Message:   statusMessage,
	}, nil
}

// CriarPagamentoCartaoDebito - cria um pagamento com cartão de débito
func (s *PagamentoService) CriarPagamentoCartaoDebito(ctx context.Context, req *contract.CriarPagamentoCartaoDebitoRequest) (*contract.CriarPagamentoCartaoDebitoResponse, error) {

	if err := req.Validate(); err != nil {
		return nil, util.WrapError("erro de validação", err, http.StatusBadRequest)
	}

	mpReq := &mercadopago.CreditCardPaymentRequest{
		TransactionAmount: req.TransactionAmount,
		Token:             req.Token,
		Description:       req.Description,
		Installments:      1,
		PaymentMethodID:   req.PaymentMethodID,
		IssuerID:          req.IssuerID,
		Capture:           true,
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
		return nil, err
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
		MetodoPagamento:      "debit_card",
		Descricao:            mpResp.Description,
		NumeroParcelas:       1,
		PaymentMethodID:      mpResp.PaymentMethodID,
		IssuerID:             mpResp.IssuerID,
		LastFourDigits:       mpResp.Card.LastFourDigits,
		FirstSixDigits:       mpResp.Card.FirstSixDigits,
		CardholderName:       mpResp.Card.Cardholder.Name,
		Captured:             mpResp.Captured,
		MomentoCriacao:       now,
		MomentoAtualizacao:   now,
	}

	if mpResp.Status == "approved" {
		pagamento.MomentoAprovacao = &now
	}

	if err := s.PagamentoRepository.Create(pagamento); err != nil {
		return nil, util.WrapError("erro ao salvar pagamento", err, http.StatusInternalServerError)
	}

	return &contract.CriarPagamentoCartaoDebitoResponse{
		Pagamento: s.modelToResponse(pagamento),
		Message:   statusMessage,
	}, nil
}

// CriarPagamentoPIX - cria um pagamento com PIX
func (s *PagamentoService) CriarPagamentoPIX(ctx context.Context, req *contract.CriarPagamentoPIXRequest) (*contract.CriarPagamentoPIXResponse, error) {

	if err := req.Validate(); err != nil {
		return nil, util.WrapError("erro de validação", err, http.StatusBadRequest)
	}

	mpReq := &mercadopago.PIXRequest{
		TransactionAmount: req.TransactionAmount,
		Description:       req.Description,
		PaymentMethodID:   "pix",
		Payer: mercadopago.PaymentPayer{
			Email: req.Payer.Email,
		},
		Metadata: map[string]string{
			"cliente_id": strconv.Itoa(req.ClienteID),
			"empresa_id": strconv.Itoa(req.EmpresaID),
		},
	}

	mpResp, err := s.MPClient.CreatePIXPayment(mpReq)
	if err != nil {
		return nil, err
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
		Moeda:                "BRL",
		MetodoPagamento:      "pix",
		Descricao:            mpResp.Description,
		NumeroParcelas:       1,
		MomentoCriacao:       now,
		MomentoAtualizacao:   now,
	}

	if mpResp.Status == "approved" {
		pagamento.MomentoAprovacao = &now
	}

	if err := s.PagamentoRepository.Create(pagamento); err != nil {
		return nil, util.WrapError("erro ao salvar pagamento", err, http.StatusInternalServerError)
	}

	qrCode := ""
	qrCodeBase64 := ""
	ticketURL := ""

	if mpResp.PointOfInteraction.Type == "PIX" {
		// Aqui você pode extrair os dados do PIX da resposta
		// Dependendo da estrutura da resposta do Mercado Pago
	}

	return &contract.CriarPagamentoPIXResponse{
		Pagamento:    s.modelToResponse(pagamento),
		Message:      statusMessage,
		QRCode:       qrCode,
		QRCodeBase64: qrCodeBase64,
		TicketURL:    ticketURL,
	}, nil
}

// BuscarPagamentos - busca pagamentos com filtros
func (s *PagamentoService) BuscarPagamentos(ctx context.Context, req *contract.BuscarPagamentosRequest) (*contract.BuscarPagamentosResponse, error) {

	if err := req.Validate(); err != nil {
		return nil, util.WrapError("erro de validação", err, http.StatusBadRequest)
	}

	if req.Limit == 0 {
		req.Limit = 20
	}

	// Por enquanto, vamos usar uma busca simples por cliente ou empresa
	// Em uma implementação real, você criaria um método Search no repository
	var pagamentos []model.Pagamento
	var total int64

	pagamentos, err := s.PagamentoRepository.GetByClienteID(1)
	if err != nil {
		return nil, util.WrapError("erro ao buscar pagamentos", err, http.StatusInternalServerError)
	}

	total = int64(len(pagamentos))

	var pagamentosResponse []contract.PagamentoResponse
	for _, pagamento := range pagamentos {
		pagamentosResponse = append(pagamentosResponse, s.modelToResponse(&pagamento))
	}

	hasMore := (req.Offset + len(pagamentos)) < int(total)

	return &contract.BuscarPagamentosResponse{
		Pagamentos: pagamentosResponse,
		Total:      int(total),
		Offset:     req.Offset,
		Limit:      req.Limit,
		HasMore:    hasMore,
	}, nil
}

// ObterPagamentoPorID - obtém um pagamento por ID
func (s *PagamentoService) ObterPagamentoPorID(ctx context.Context, paymentID int64) (*contract.ObterPagamentoResponse, error) {

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

	return &contract.ObterPagamentoResponse{
		Pagamento: s.modelToResponse(pagamento),
	}, nil
}

// AtualizarPagamento - atualiza um pagamento
func (s *PagamentoService) AtualizarPagamento(ctx context.Context, req *contract.AtualizarPagamentoRequest) (*contract.AtualizarPagamentoResponse, error) {

	if err := req.Validate(); err != nil {
		return nil, util.WrapError("erro de validação", err, http.StatusBadRequest)
	}

	paymentIDStr := strconv.FormatInt(req.ID, 10)
	pagamento, err := s.PagamentoRepository.GetByMercadoPagoPaymentID(paymentIDStr)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, util.WrapError("pagamento não encontrado", err, http.StatusNotFound)
		}
		return nil, util.WrapError("erro ao buscar pagamento", err, http.StatusInternalServerError)
	}

	if req.Description != "" {
		pagamento.Descricao = req.Description
	}

	if req.Metadata != nil {
		// Aqui você pode atualizar metadados se necessário
		// Dependendo da estrutura do seu modelo
	}

	if req.NotificationURL != "" {
		// Aqui você pode atualizar a URL de notificação se necessário
		// Dependendo da estrutura do seu modelo
	}

	pagamento.MomentoAtualizacao = time.Now()

	if err := s.PagamentoRepository.Update(pagamento); err != nil {
		return nil, util.WrapError("erro ao atualizar pagamento", err, http.StatusInternalServerError)
	}

	return &contract.AtualizarPagamentoResponse{
		Pagamento: s.modelToResponse(pagamento),
		Message:   "Pagamento atualizado com sucesso",
	}, nil
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
