package model

import "time"

// Cliente representa a entidade de cliente
type Cliente struct {
	ID                 int       `gorm:"column:id;primaryKey"`
	Nome               string    `gorm:"column:nome"`
	Email              string    `gorm:"column:email"`
	Senha              string    `gorm:"column:senha"`
	CPF                string    `gorm:"column:cpf"`
	Telefone           string    `gorm:"column:telefone"`
	DataNascimento     time.Time `gorm:"column:data_nascimento"`
	MomentoCadastro    time.Time `gorm:"column:momento_cadastro"`
	MomentoAtualizacao time.Time `gorm:"column:momento_atualizacao"`
}

// TableName especifica o nome da tabela no banco de dados
func (Cliente) TableName() string {
	return "clientes"
}
