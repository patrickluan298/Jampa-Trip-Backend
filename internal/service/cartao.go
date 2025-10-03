package service

import (
	"context"
	"net/http"

	"github.com/jampa_trip/internal/contract"
	"github.com/jampa_trip/pkg/config"
	"github.com/jampa_trip/pkg/mercadopago"
	"github.com/jampa_trip/pkg/util"
	"gorm.io/gorm"
)

type CartaoService struct {
	db     *gorm.DB
	client *mercadopago.Client
}

// CartaoServiceNew - cria uma nova instância do service de cartão
func CartaoServiceNew(db *gorm.DB) *CartaoService {
	cfg, _ := config.LoadConfig()
	client := mercadopago.NewClient(cfg.MercadoPagoAccessToken, cfg.MercadoPagoBaseURL)

	return &CartaoService{
		db:     db,
		client: client,
	}
}

// Create - cria um cartão para um cliente
func (s *CartaoService) Create(ctx context.Context, customerID string, req *contract.CreateCartaoRequest) (*contract.CreateCartaoResponse, error) {

	mpReq := &mercadopago.CustomerCardRequest{
		Token:           req.Token,
		PaymentMethodID: req.PaymentMethodID,
		IssuerID:        req.IssuerID,
		Cardholder: mercadopago.CustomerCardholder{
			Name: req.Cardholder.Name,
			Identification: mercadopago.CustomerIdentification{
				Type:   req.Cardholder.Identification.Type,
				Number: req.Cardholder.Identification.Number,
			},
		},
		Metadata: req.Metadata,
	}

	mpResp, err := s.client.CreateCustomerCard(ctx, customerID, mpReq)
	if err != nil {
		return nil, util.WrapError("erro ao criar cartão no Mercado Pago", err, http.StatusInternalServerError)
	}

	response := &contract.CreateCartaoResponse{
		Cartao: contract.CartaoResponse{
			ID:              mpResp.ID,
			CustomerID:      mpResp.CustomerID,
			FirstSixDigits:  mpResp.FirstSixDigits,
			LastFourDigits:  mpResp.LastFourDigits,
			ExpirationMonth: mpResp.ExpirationMonth,
			ExpirationYear:  mpResp.ExpirationYear,
			SecurityCode: contract.SecurityCodeInfo{
				Length:       mpResp.SecurityCode.Length,
				CardLocation: mpResp.SecurityCode.CardLocation,
				Mode:         mpResp.SecurityCode.Mode,
			},
			Issuer: contract.IssuerInfo{
				ID:   mpResp.Issuer.ID,
				Name: mpResp.Issuer.Name,
			},
			PaymentMethod: contract.CartaoPaymentMethodInfo{
				ID:   mpResp.PaymentMethod.ID,
				Name: mpResp.PaymentMethod.Name,
			},
			Cardholder: contract.CardholderResponse{
				Name: mpResp.Cardholder.Name,
				Identification: contract.IdentificationResponse{
					Type:   mpResp.Cardholder.Identification.Type,
					Number: mpResp.Cardholder.Identification.Number,
				},
			},
			DateCreated:     mpResp.DateCreated,
			DateLastUpdated: mpResp.DateLastUpdated,
			Metadata:        mpResp.Metadata,
		},
		Message: "Cartão criado com sucesso",
	}

	return response, nil
}

// List - lista os cartões de um cliente
func (s *CartaoService) List(ctx context.Context, customerID string) (*contract.ListCartoesResponse, error) {

	mpCards, err := s.client.ListCustomerCards(ctx, customerID)
	if err != nil {
		return nil, util.WrapError("erro ao listar cartões no Mercado Pago", err, http.StatusInternalServerError)
	}

	cartoes := make([]contract.CartaoResponse, len(mpCards))
	for i, mpCard := range mpCards {
		cartoes[i] = contract.CartaoResponse{
			ID:              mpCard.ID,
			CustomerID:      mpCard.CustomerID,
			FirstSixDigits:  mpCard.FirstSixDigits,
			LastFourDigits:  mpCard.LastFourDigits,
			ExpirationMonth: mpCard.ExpirationMonth,
			ExpirationYear:  mpCard.ExpirationYear,
			SecurityCode: contract.SecurityCodeInfo{
				Length:       mpCard.SecurityCode.Length,
				CardLocation: mpCard.SecurityCode.CardLocation,
				Mode:         mpCard.SecurityCode.Mode,
			},
			Issuer: contract.IssuerInfo{
				ID:   mpCard.Issuer.ID,
				Name: mpCard.Issuer.Name,
			},
			PaymentMethod: contract.CartaoPaymentMethodInfo{
				ID:   mpCard.PaymentMethod.ID,
				Name: mpCard.PaymentMethod.Name,
			},
			Cardholder: contract.CardholderResponse{
				Name: mpCard.Cardholder.Name,
				Identification: contract.IdentificationResponse{
					Type:   mpCard.Cardholder.Identification.Type,
					Number: mpCard.Cardholder.Identification.Number,
				},
			},
			DateCreated:     mpCard.DateCreated,
			DateLastUpdated: mpCard.DateLastUpdated,
			Metadata:        mpCard.Metadata,
		}
	}

	response := &contract.ListCartoesResponse{
		Cartoes: cartoes,
		Total:   len(cartoes),
	}

	return response, nil
}

// Get - obtém um cartão específico de um cliente
func (s *CartaoService) Get(ctx context.Context, customerID, cardID string) (*contract.CartaoResponse, error) {

	mpCard, err := s.client.GetCustomerCard(ctx, customerID, cardID)
	if err != nil {
		return nil, util.WrapError("erro ao obter cartão no Mercado Pago", err, http.StatusInternalServerError)
	}

	response := &contract.CartaoResponse{
		ID:              mpCard.ID,
		CustomerID:      mpCard.CustomerID,
		FirstSixDigits:  mpCard.FirstSixDigits,
		LastFourDigits:  mpCard.LastFourDigits,
		ExpirationMonth: mpCard.ExpirationMonth,
		ExpirationYear:  mpCard.ExpirationYear,
		SecurityCode: contract.SecurityCodeInfo{
			Length:       mpCard.SecurityCode.Length,
			CardLocation: mpCard.SecurityCode.CardLocation,
			Mode:         mpCard.SecurityCode.Mode,
		},
		Issuer: contract.IssuerInfo{
			ID:   mpCard.Issuer.ID,
			Name: mpCard.Issuer.Name,
		},
		PaymentMethod: contract.CartaoPaymentMethodInfo{
			ID:   mpCard.PaymentMethod.ID,
			Name: mpCard.PaymentMethod.Name,
		},
		Cardholder: contract.CardholderResponse{
			Name: mpCard.Cardholder.Name,
			Identification: contract.IdentificationResponse{
				Type:   mpCard.Cardholder.Identification.Type,
				Number: mpCard.Cardholder.Identification.Number,
			},
		},
		DateCreated:     mpCard.DateCreated,
		DateLastUpdated: mpCard.DateLastUpdated,
		Metadata:        mpCard.Metadata,
	}

	return response, nil
}

// Update - atualiza um cartão de um cliente
func (s *CartaoService) Update(ctx context.Context, customerID, cardID string, req *contract.UpdateCartaoRequest) (*contract.UpdateCartaoResponse, error) {

	mpReq := &mercadopago.CustomerCardUpdateRequest{
		Cardholder: mercadopago.CustomerCardholder{
			Name: req.Cardholder.Name,
			Identification: mercadopago.CustomerIdentification{
				Type:   req.Cardholder.Identification.Type,
				Number: req.Cardholder.Identification.Number,
			},
		},
		Metadata: req.Metadata,
		Default:  req.Default,
	}

	mpResp, err := s.client.UpdateCustomerCard(ctx, customerID, cardID, mpReq)
	if err != nil {
		return nil, util.WrapError("erro ao atualizar cartão no Mercado Pago", err, http.StatusInternalServerError)
	}

	response := &contract.UpdateCartaoResponse{
		Cartao: contract.CartaoResponse{
			ID:              mpResp.ID,
			CustomerID:      mpResp.CustomerID,
			FirstSixDigits:  mpResp.FirstSixDigits,
			LastFourDigits:  mpResp.LastFourDigits,
			ExpirationMonth: mpResp.ExpirationMonth,
			ExpirationYear:  mpResp.ExpirationYear,
			SecurityCode: contract.SecurityCodeInfo{
				Length:       mpResp.SecurityCode.Length,
				CardLocation: mpResp.SecurityCode.CardLocation,
				Mode:         mpResp.SecurityCode.Mode,
			},
			Issuer: contract.IssuerInfo{
				ID:   mpResp.Issuer.ID,
				Name: mpResp.Issuer.Name,
			},
			PaymentMethod: contract.CartaoPaymentMethodInfo{
				ID:   mpResp.PaymentMethod.ID,
				Name: mpResp.PaymentMethod.Name,
			},
			Cardholder: contract.CardholderResponse{
				Name: mpResp.Cardholder.Name,
				Identification: contract.IdentificationResponse{
					Type:   mpResp.Cardholder.Identification.Type,
					Number: mpResp.Cardholder.Identification.Number,
				},
			},
			DateCreated:     mpResp.DateCreated,
			DateLastUpdated: mpResp.DateLastUpdated,
			Metadata:        mpResp.Metadata,
		},
		Message: "Cartão atualizado com sucesso",
	}

	return response, nil
}

// Delete - exclui um cartão de um cliente
func (s *CartaoService) Delete(ctx context.Context, customerID, cardID string) (*contract.DeleteCartaoResponse, error) {

	err := s.client.DeleteCustomerCard(ctx, customerID, cardID)
	if err != nil {
		return nil, util.WrapError("erro ao excluir cartão no Mercado Pago", err, http.StatusInternalServerError)
	}

	response := &contract.DeleteCartaoResponse{
		Message: "Cartão excluído com sucesso",
	}

	return response, nil
}
