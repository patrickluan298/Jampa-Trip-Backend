package service

import (
	"context"
	"fmt"
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

// OrderService - objeto de contexto para operações com orders
type OrderService struct {
	OrderRepository       *repository.OrderRepository
	TransactionRepository *repository.TransactionRepository
	MPClient              *mercadopago.Client
	DB                    *gorm.DB
}

// OrderServiceNew - construtor do objeto
func OrderServiceNew(DB *gorm.DB) *OrderService {
	cfg, _ := config.LoadConfig()

	return &OrderService{
		OrderRepository:       repository.OrderRepositoryNew(DB),
		TransactionRepository: repository.TransactionRepositoryNew(DB),
		MPClient:              mercadopago.NewClient(cfg.MercadoPagoAccessToken, cfg.MercadoPagoBaseURL),
		DB:                    DB,
	}
}

// Create - cria uma nova ordem no Mercado Pago e persiste localmente
func (s *OrderService) Create(ctx context.Context, req *contract.CreateOrderRequest) (*contract.CreateOrderResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, util.WrapError("erro de validação", err, http.StatusBadRequest)
	}

	var mpItems []mercadopago.OrderItem
	for _, item := range req.Items {
		mpItems = append(mpItems, mercadopago.OrderItem{
			ID:          item.ID,
			Title:       item.Title,
			Description: item.Description,
			PictureURL:  item.PictureURL,
			CategoryID:  item.CategoryID,
			Quantity:    item.Quantity,
			CurrencyID:  item.CurrencyID,
			UnitPrice:   item.UnitPrice,
		})
	}

	mpReq := &mercadopago.OrderRequest{
		ExternalReference: req.ExternalReference,
		TotalAmount:       req.TotalAmount,
		Items:             mpItems,
		Payer: mercadopago.Payer{
			Name:  req.Payer.Name,
			Email: req.Payer.Email,
			Phone: mercadopago.Phone{
				AreaCode: req.Payer.Phone.AreaCode,
				Number:   req.Payer.Phone.Number,
			},
			Address: mercadopago.Address{
				StreetName:   req.Payer.Address.StreetName,
				StreetNumber: req.Payer.Address.StreetNumber,
				ZipCode:      req.Payer.Address.ZipCode,
			},
		},
		NotificationURL: req.NotificationURL,
		Description:     req.Description,
		Metadata: map[string]string{
			"cliente_id": strconv.Itoa(req.ClienteID),
			"empresa_id": strconv.Itoa(req.EmpresaID),
		},
	}

	for key, value := range req.Metadata {
		mpReq.Metadata[key] = value
	}

	mpResp, err := s.MPClient.MPCreateOrder(ctx, mpReq)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	order := &model.Order{
		ClienteID:          req.ClienteID,
		EmpresaID:          req.EmpresaID,
		MercadoPagoOrderID: mpResp.ID,
		ExternalReference:  mpResp.ExternalReference,
		TotalAmount:        mpResp.TotalAmount,
		Currency:           "BRL",
		Status:             mpResp.Status,
		StatusDetail:       mpResp.StatusDetail,
		Description:        mpResp.Description,
		NotificationURL:    mpResp.NotificationURL,
		MomentoCriacao:     now,
		MomentoAtualizacao: now,
	}

	if err := s.OrderRepository.Create(order); err != nil {
		return nil, util.WrapError("erro ao salvar ordem", err, http.StatusInternalServerError)
	}

	return &contract.CreateOrderResponse{
		Order:   s.modelToResponse(order),
		Message: "Ordem criada com sucesso",
	}, nil
}

// Get - obtém uma ordem por ID e sincroniza com Mercado Pago
func (s *OrderService) Get(ctx context.Context, orderID string) (*contract.GetOrderResponse, error) {
	order, err := s.OrderRepository.GetByMercadoPagoOrderID(orderID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, util.WrapError("ordem não encontrada", err, http.StatusNotFound)
		}
		return nil, util.WrapError("erro ao buscar ordem", err, http.StatusInternalServerError)
	}

	mpResp, err := s.MPClient.MPGetOrder(ctx, orderID)
	if err == nil {
		if mpResp.Status != order.Status || mpResp.StatusDetail != order.StatusDetail {
			order.Status = mpResp.Status
			order.StatusDetail = mpResp.StatusDetail
			order.MomentoAtualizacao = time.Now()
			s.OrderRepository.Update(order)
		}
	}

	return &contract.GetOrderResponse{
		Order: s.modelToResponse(order),
	}, nil
}

// List - lista ordens com filtros e paginação
func (s *OrderService) List(ctx context.Context, req *contract.ListOrdersRequest) (*contract.ListOrdersResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, util.WrapError("erro de validação", err, http.StatusBadRequest)
	}

	if req.Limit == 0 {
		req.Limit = 20
	}

	var orders []model.Order
	var total int64
	var err error

	filters := make(map[string]interface{})
	if req.ClienteID > 0 {
		filters["cliente_id"] = req.ClienteID
	}
	if req.EmpresaID > 0 {
		filters["empresa_id"] = req.EmpresaID
	}
	if req.Status != "" {
		filters["status"] = req.Status
	}
	if req.ExternalReference != "" {
		filters["external_reference"] = req.ExternalReference
	}

	orders, total, err = s.OrderRepository.List(filters, req.Limit, req.Offset)
	if err != nil {
		return nil, util.WrapError("erro ao buscar ordens", err, http.StatusInternalServerError)
	}

	var ordersResponse []contract.OrderResponse
	for _, order := range orders {
		ordersResponse = append(ordersResponse, s.modelToResponse(&order))
	}

	hasMore := (req.Offset + len(orders)) < int(total)

	return &contract.ListOrdersResponse{
		Orders:  ordersResponse,
		Total:   total,
		Offset:  req.Offset,
		Limit:   req.Limit,
		HasMore: hasMore,
	}, nil
}

// Capture - captura uma ordem autorizada
func (s *OrderService) Capture(ctx context.Context, orderID string, req *contract.CaptureOrderRequest) (*contract.CaptureOrderResponse, error) {
	order, err := s.OrderRepository.GetByMercadoPagoOrderID(orderID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, util.WrapError("ordem não encontrada", err, http.StatusNotFound)
		}
		return nil, util.WrapError("erro ao buscar ordem", err, http.StatusInternalServerError)
	}

	if order.Status != string(model.OrderStatusAuthorized) {
		return nil, util.WrapError("ordem não está autorizada para captura", nil, http.StatusBadRequest)
	}

	mpResp, err := s.MPClient.MPCaptureOrder(ctx, orderID)
	if err != nil {
		return nil, err
	}

	order.Status = mpResp.Status
	order.StatusDetail = mpResp.StatusDetail
	order.MomentoAtualizacao = time.Now()

	if err := s.OrderRepository.Update(order); err != nil {
		return nil, util.WrapError("erro ao atualizar ordem", err, http.StatusInternalServerError)
	}

	return &contract.CaptureOrderResponse{
		Order:   s.modelToResponse(order),
		Message: "Ordem capturada com sucesso",
	}, nil
}

// Cancel - cancela uma ordem
func (s *OrderService) Cancel(ctx context.Context, orderID string) (*contract.CancelOrderResponse, error) {
	order, err := s.OrderRepository.GetByMercadoPagoOrderID(orderID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, util.WrapError("ordem não encontrada", err, http.StatusNotFound)
		}
		return nil, util.WrapError("erro ao buscar ordem", err, http.StatusInternalServerError)
	}

	if order.Status == string(model.OrderStatusCancelled) {
		return nil, util.WrapError("ordem já está cancelada", nil, http.StatusBadRequest)
	}

	mpResp, err := s.MPClient.MPCancelOrder(ctx, orderID)
	if err != nil {
		return nil, err
	}

	order.Status = mpResp.Status
	order.StatusDetail = mpResp.StatusDetail
	order.MomentoAtualizacao = time.Now()

	if err := s.OrderRepository.Update(order); err != nil {
		return nil, util.WrapError("erro ao atualizar ordem", err, http.StatusInternalServerError)
	}

	return &contract.CancelOrderResponse{
		Order:   s.modelToResponse(order),
		Message: "Ordem cancelada com sucesso",
	}, nil
}

// Refund - reembolsa uma ordem
func (s *OrderService) Refund(ctx context.Context, orderID string, req *contract.RefundOrderRequest) (*contract.RefundOrderResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, util.WrapError("erro de validação", err, http.StatusBadRequest)
	}

	order, err := s.OrderRepository.GetByMercadoPagoOrderID(orderID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, util.WrapError("ordem não encontrada", err, http.StatusNotFound)
		}
		return nil, util.WrapError("erro ao buscar ordem", err, http.StatusInternalServerError)
	}

	if order.Status != string(model.OrderStatusPaid) {
		return nil, util.WrapError("apenas ordens pagas podem ser reembolsadas", nil, http.StatusBadRequest)
	}

	mpResp, err := s.MPClient.MPRefundOrder(ctx, orderID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	transaction := &model.Transaction{
		OrderID:            order.ID,
		Type:               string(model.TransactionTypeRefund),
		Status:             "approved",
		Amount:             order.TotalAmount,
		Currency:           order.Currency,
		MomentoCriacao:     now,
		MomentoAtualizacao: now,
	}

	if req.Amount != nil {
		transaction.Amount = *req.Amount
	}

	if err := s.TransactionRepository.Create(transaction); err != nil {
		return nil, util.WrapError("erro ao criar transação de reembolso", err, http.StatusInternalServerError)
	}

	order.Status = mpResp.Status
	order.StatusDetail = mpResp.StatusDetail
	order.MomentoAtualizacao = time.Now()

	if err := s.OrderRepository.Update(order); err != nil {
		return nil, util.WrapError("erro ao atualizar ordem", err, http.StatusInternalServerError)
	}

	refundDetails := contract.RefundDetailsResponse{
		RefundID:    int64(transaction.ID),
		Amount:      transaction.Amount,
		Status:      transaction.Status,
		DateCreated: transaction.MomentoCriacao.Format(time.RFC3339),
	}

	return &contract.RefundOrderResponse{
		Order:         s.modelToResponse(order),
		RefundDetails: refundDetails,
		Message:       "Ordem reembolsada com sucesso",
	}, nil
}

// AddTransaction - adiciona uma transação à ordem
func (s *OrderService) AddTransaction(ctx context.Context, orderID string, req *contract.AddTransactionRequest) (*contract.AddTransactionResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, util.WrapError("erro de validação", err, http.StatusBadRequest)
	}

	order, err := s.OrderRepository.GetByMercadoPagoOrderID(orderID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, util.WrapError("ordem não encontrada", err, http.StatusNotFound)
		}
		return nil, util.WrapError("erro ao buscar ordem", err, http.StatusInternalServerError)
	}

	mpTxReq := &mercadopago.OrderTransactionRequest{
		PaymentID:       req.PaymentID,
		PaymentMethodID: req.PaymentMethodID,
		Amount:          req.Amount,
	}

	mpResp, err := s.MPClient.MPAddTransaction(ctx, orderID, mpTxReq)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	transaction := &model.Transaction{
		OrderID:            order.ID,
		PaymentID:          req.PaymentID,
		Type:               string(model.TransactionTypePayment),
		Status:             "pending",
		Amount:             req.Amount,
		Currency:           req.Currency,
		PaymentMethodID:    req.PaymentMethodID,
		MomentoCriacao:     now,
		MomentoAtualizacao: now,
	}

	if err := s.TransactionRepository.Create(transaction); err != nil {
		return nil, util.WrapError("erro ao criar transação", err, http.StatusInternalServerError)
	}

	if mpResp.Status != order.Status {
		order.Status = mpResp.Status
		order.StatusDetail = mpResp.StatusDetail
		order.MomentoAtualizacao = time.Now()
		s.OrderRepository.Update(order)
	}

	return &contract.AddTransactionResponse{
		Transaction: s.transactionToResponse(transaction),
		Order:       s.modelToResponse(order),
		Message:     "Transação adicionada com sucesso",
	}, nil
}

// UpdateTransaction - atualiza uma transação existente
func (s *OrderService) UpdateTransaction(ctx context.Context, orderID string, transactionID int, req *contract.UpdateTransactionRequest) (*contract.UpdateTransactionResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, util.WrapError("erro de validação", err, http.StatusBadRequest)
	}

	transaction, err := s.TransactionRepository.GetByID(transactionID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, util.WrapError("transação não encontrada", err, http.StatusNotFound)
		}
		return nil, util.WrapError("erro ao buscar transação", err, http.StatusInternalServerError)
	}

	if req.Status != "" {
		transaction.Status = req.Status
	}
	if req.StatusDetail != "" {
		transaction.StatusDetail = req.StatusDetail
	}
	if req.Amount > 0 {
		transaction.Amount = req.Amount
	}
	transaction.MomentoAtualizacao = time.Now()

	if err := s.TransactionRepository.Update(transaction); err != nil {
		return nil, util.WrapError("erro ao atualizar transação", err, http.StatusInternalServerError)
	}

	return &contract.UpdateTransactionResponse{
		Transaction: s.transactionToResponse(transaction),
		Message:     "Transação atualizada com sucesso",
	}, nil
}

// DeleteTransaction - deleta uma transação
func (s *OrderService) DeleteTransaction(ctx context.Context, orderID string, transactionID int) error {
	if _, err := s.TransactionRepository.GetByID(transactionID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return util.WrapError("transação não encontrada", err, http.StatusNotFound)
		}
		return util.WrapError("erro ao buscar transação", err, http.StatusInternalServerError)
	}

	if err := s.TransactionRepository.Delete(transactionID); err != nil {
		return util.WrapError("erro ao deletar transação", err, http.StatusInternalServerError)
	}

	return nil
}

// ProcessOrder - processa uma ordem online (pagamento)
func (s *OrderService) ProcessOrder(ctx context.Context, orderID string, req *contract.ProcessOrderRequest) (*contract.ProcessOrderResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, util.WrapError("erro de validação", err, http.StatusBadRequest)
	}

	order, err := s.OrderRepository.GetByMercadoPagoOrderID(orderID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, util.WrapError("ordem não encontrada", err, http.StatusNotFound)
		}
		return nil, util.WrapError("erro ao buscar ordem", err, http.StatusInternalServerError)
	}

	mpReq := &mercadopago.OrderOnlineRequest{
		TransactionAmount: order.TotalAmount,
		Token:             req.Token,
		Description:       order.Description,
		Installments:      req.Installments,
		PaymentMethodID:   req.PaymentMethodID,
		IssuerID:          req.IssuerID,
		Payer: mercadopago.OrderOnlinePayer{
			Email:     req.Payer.Email,
			FirstName: req.Payer.FirstName,
			LastName:  req.Payer.LastName,
			Identification: mercadopago.OrderOnlineIdentification{
				Type:   req.Payer.Identification.Type,
				Number: req.Payer.Identification.Number,
			},
		},
		Metadata: req.Metadata,
	}

	paymentResp, err := s.MPClient.MPProcessOrder(ctx, orderID, mpReq)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	transaction := &model.Transaction{
		OrderID:            order.ID,
		PaymentID:          paymentResp.ID,
		Type:               string(model.TransactionTypePayment),
		Status:             paymentResp.Status,
		StatusDetail:       paymentResp.StatusDetail,
		Amount:             paymentResp.TransactionAmount,
		Currency:           order.Currency,
		PaymentMethodID:    paymentResp.PaymentMethodID,
		MomentoCriacao:     now,
		MomentoAtualizacao: now,
	}

	if err := s.TransactionRepository.Create(transaction); err != nil {
		return nil, util.WrapError("erro ao criar transação de pagamento", err, http.StatusInternalServerError)
	}

	if paymentResp.Status == "approved" {
		order.Status = string(model.OrderStatusPaid)
	} else if paymentResp.Status == "authorized" {
		order.Status = string(model.OrderStatusAuthorized)
	} else {
		order.Status = string(model.OrderStatusProcessing)
	}
	order.StatusDetail = paymentResp.StatusDetail
	order.MomentoAtualizacao = time.Now()

	if err := s.OrderRepository.Update(order); err != nil {
		return nil, util.WrapError("erro ao atualizar ordem", err, http.StatusInternalServerError)
	}

	paymentDetails := contract.ProcessPaymentDetails{
		PaymentID:         paymentResp.ID,
		Status:            paymentResp.Status,
		StatusDetail:      paymentResp.StatusDetail,
		PaymentMethodID:   paymentResp.PaymentMethodID,
		TransactionAmount: paymentResp.TransactionAmount,
		DateCreated:       paymentResp.DateCreated,
		DateApproved:      paymentResp.DateApproved,
	}

	return &contract.ProcessOrderResponse{
		Order:          s.modelToResponse(order),
		PaymentDetails: paymentDetails,
		Message:        s.getPaymentMessage(paymentResp.Status),
	}, nil
}

// modelToResponse - converte model.Order para contract.OrderResponse
func (s *OrderService) modelToResponse(order *model.Order) contract.OrderResponse {
	var transactions []contract.TransactionResponse
	for _, tx := range order.Transactions {
		transactions = append(transactions, s.transactionToResponse(&tx))
	}

	return contract.OrderResponse{
		ID:                 order.ID,
		ClienteID:          order.ClienteID,
		EmpresaID:          order.EmpresaID,
		MercadoPagoOrderID: order.MercadoPagoOrderID,
		ExternalReference:  order.ExternalReference,
		TotalAmount:        order.TotalAmount,
		Currency:           order.Currency,
		Status:             order.Status,
		StatusDetail:       order.StatusDetail,
		Description:        order.Description,
		NotificationURL:    order.NotificationURL,
		MomentoCriacao:     order.MomentoCriacao,
		MomentoAtualizacao: order.MomentoAtualizacao,
		MomentoExpiracao:   order.MomentoExpiracao,
		StatusDisplay:      order.GetStatusDisplay(),
		Transactions:       transactions,
	}
}

// transactionToResponse - converte model.Transaction para contract.TransactionResponse
func (s *OrderService) transactionToResponse(tx *model.Transaction) contract.TransactionResponse {
	return contract.TransactionResponse{
		ID:                       tx.ID,
		OrderID:                  tx.OrderID,
		MercadoPagoTransactionID: tx.MercadoPagoTransactionID,
		PaymentID:                tx.PaymentID,
		Type:                     tx.Type,
		Status:                   tx.Status,
		StatusDetail:             tx.StatusDetail,
		Amount:                   tx.Amount,
		Currency:                 tx.Currency,
		PaymentMethodID:          tx.PaymentMethodID,
		PaymentTypeID:            tx.PaymentTypeID,
		MomentoCriacao:           tx.MomentoCriacao,
		MomentoAtualizacao:       tx.MomentoAtualizacao,
		StatusDisplay:            tx.GetStatusDisplay(),
		TypeDisplay:              tx.GetTypeDisplay(),
	}
}

// getPaymentMessage - retorna mensagem amigável para status de pagamento
func (s *OrderService) getPaymentMessage(status string) string {
	messages := map[string]string{
		"approved":   "Pagamento aprovado com sucesso",
		"authorized": "Pagamento autorizado - aguardando captura",
		"pending":    "Pagamento pendente",
		"in_process": "Pagamento em processamento",
		"rejected":   "Pagamento rejeitado",
		"cancelled":  "Pagamento cancelado",
	}

	if msg, ok := messages[status]; ok {
		return msg
	}
	return fmt.Sprintf("Status do pagamento: %s", status)
}
