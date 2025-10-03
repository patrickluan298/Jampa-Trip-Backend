package model

import "time"

// Empresa representa a entidade de empresa
type Empresa struct {
	ID                 int       `gorm:"column:id;primaryKey"`
	Nome               string    `gorm:"column:nome"`
	Email              string    `gorm:"column:email"`
	Senha              string    `gorm:"column:senha"`
	CNPJ               string    `gorm:"column:cnpj"`
	Telefone           string    `gorm:"column:telefone"`
	Endereco           string    `gorm:"column:endereco"`
	MomentoCadastro    time.Time `gorm:"column:momento_cadastro"`
	MomentoAtualizacao time.Time `gorm:"column:momento_atualizacao"`
}

// TableName especifica o nome da tabela no banco de dados
func (Empresa) TableName() string {
	return "empresas"
}
