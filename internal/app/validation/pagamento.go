package validation

import "github.com/jampa_trip/internal/app/types"

// IsValidStatus - verifica se o status é válido
func IsValidStatus(status types.StatusPagamento) bool {
	switch status {
	case types.StatusPending, types.StatusApproved, types.StatusAuthorized, types.StatusInProcess,
		types.StatusInMediation, types.StatusRejected, types.StatusCancelled, types.StatusRefunded, types.StatusChargedBack:
		return true
	default:
		return false
	}
}

// IsValidPaymentMethod - verifica se o método de pagamento é válido
func IsValidPaymentMethod(method types.MetodoPagamento) bool {
	switch method {
	case types.MetodoCartaoCredito, types.MetodoCartaoDebito, types.MetodoPIX, types.MetodoBoleto, types.MetodoPec:
		return true
	default:
		return false
	}
}

// IsValidCurrency - verifica se a moeda é válida
func IsValidCurrency(currency types.Moeda) bool {
	switch currency {
	case types.MoedaBRL, types.MoedaUSD, types.MoedaEUR, types.MoedaARS, types.MoedaCLP, types.MoedaCOP, types.MoedaMXN, types.MoedaPEN, types.MoedaUYU:
		return true
	default:
		return false
	}
}
