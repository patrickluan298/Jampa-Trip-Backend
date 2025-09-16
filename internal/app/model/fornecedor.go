package model

import "time"

// Fornecedor representa a entidade de fornecedor
type Fornecedor struct {
	ID              int       `gorm:"column:id;primaryKey"`
	Nome            string    `gorm:"column:nome"`
	Email           string    `gorm:"column:email"`
	Senha           string    `gorm:"column:senha"`
	CNPJ            string    `gorm:"column:cnpj"`
	Telefone        string    `gorm:"column:telefone"`
	Endereco        string    `gorm:"column:endereco"`
	MomentoCadastro time.Time `gorm:"column:momento_cadastro"`
}

// TableName especifica o nome da tabela no banco de dados
func (Fornecedor) TableName() string {
	return "fornecedores"
}
