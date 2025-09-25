package types

// StatusPagamento - define os possíveis status de um pagamento
type StatusPagamento string

const (
	StatusPending     StatusPagamento = "pending"
	StatusApproved    StatusPagamento = "approved"
	StatusAuthorized  StatusPagamento = "authorized"
	StatusInProcess   StatusPagamento = "in_process"
	StatusInMediation StatusPagamento = "in_mediation"
	StatusRejected    StatusPagamento = "rejected"
	StatusCancelled   StatusPagamento = "cancelled"
	StatusRefunded    StatusPagamento = "refunded"
	StatusChargedBack StatusPagamento = "charged_back"
)
