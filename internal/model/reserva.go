package model

import "time"

// Reserva - representa a entidade de reserva
type Reserva struct {
	ID                  int        `gorm:"column:id;primaryKey;autoIncrement"`
	ClienteID           int        `gorm:"column:cliente_id;not null"`
	EmpresaID           int        `gorm:"column:empresa_id;not null"`
	PagamentoID         int        `gorm:"column:pagamento_id"`
	Status              string     `gorm:"column:status;not null;default:'pendente'"`
	DataReserva         time.Time  `gorm:"column:data_reserva;not null"`
	DataPasseio         time.Time  `gorm:"column:data_passeio;not null"`
	QuantidadePessoas   int        `gorm:"column:quantidade_pessoas;not null;default:1"`
	ValorTotal          float64    `gorm:"column:valor_total;not null;type:decimal(10,2)"`
	Observacoes         string     `gorm:"column:observacoes"`
	MomentoCriacao      time.Time  `gorm:"column:momento_criacao;not null;default:CURRENT_TIMESTAMP"`
	MomentoAtualizacao  time.Time  `gorm:"column:momento_atualizacao;not null;default:CURRENT_TIMESTAMP"`
	MomentoCancelamento *time.Time `gorm:"column:momento_cancelamento"`

	// Relacionamentos
	Cliente   Client    `gorm:"foreignKey:ClienteID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Empresa   Company   `gorm:"foreignKey:EmpresaID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Pagamento Pagamento `gorm:"foreignKey:PagamentoID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

// TableName - especifica o nome da tabela no banco de dados
func (Reserva) TableName() string {
	return "reservas"
}

// StatusReserva - define os possíveis status de uma reserva
type StatusReserva string

const (
	StatusReservaPendente   StatusReserva = "pendente"
	StatusReservaConfirmada StatusReserva = "confirmada"
	StatusReservaCancelada  StatusReserva = "cancelada"
	StatusReservaConcluida  StatusReserva = "concluida"
)

// Métodos de validação para Reserva
func (r *Reserva) IsValid() bool {
	return r.ClienteID > 0 && r.EmpresaID > 0 && r.DataPasseio.After(time.Now()) && r.QuantidadePessoas > 0
}

func (r *Reserva) IsPending() bool {
	return r.Status == string(StatusReservaPendente)
}

func (r *Reserva) IsConfirmed() bool {
	return r.Status == string(StatusReservaConfirmada)
}

func (r *Reserva) IsCancelled() bool {
	return r.Status == string(StatusReservaCancelada)
}

func (r *Reserva) IsCompleted() bool {
	return r.Status == string(StatusReservaConcluida)
}

func (r *Reserva) CanBeCancelled() bool {
	return r.Status == string(StatusReservaPendente) || r.Status == string(StatusReservaConfirmada)
}

func (r *Reserva) UpdateStatus(status StatusReserva) {
	r.Status = string(status)
	r.MomentoAtualizacao = time.Now()
}

func (r *Reserva) GetStatusDisplay() string {
	switch r.Status {
	case string(StatusReservaPendente):
		return "Pendente"
	case string(StatusReservaConfirmada):
		return "Confirmada"
	case string(StatusReservaCancelada):
		return "Cancelada"
	case string(StatusReservaConcluida):
		return "Concluída"
	default:
		return "Desconhecido"
	}
}

// Funções de validação
func IsValidReservaStatus(status StatusReserva) bool {
	switch status {
	case StatusReservaPendente, StatusReservaConfirmada, StatusReservaCancelada, StatusReservaConcluida:
		return true
	default:
		return false
	}
}
