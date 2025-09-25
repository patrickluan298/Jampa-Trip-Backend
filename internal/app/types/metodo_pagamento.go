package types

// MetodoPagamento - define os métodos de pagamento suportados
type MetodoPagamento string

const (
	MetodoCartaoCredito MetodoPagamento = "credit_card"
	MetodoCartaoDebito  MetodoPagamento = "debit_card"
	MetodoPIX           MetodoPagamento = "pix"
	MetodoBoleto        MetodoPagamento = "bolbradesco"
	MetodoPec           MetodoPagamento = "pec"
)
