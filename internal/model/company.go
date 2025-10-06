package model

import "time"

// Company representa a entidade de empresa
type Company struct {
	ID        int       `gorm:"column:id;primaryKey"`
	Name      string    `gorm:"column:name"`
	Email     string    `gorm:"column:email"`
	Password  string    `gorm:"column:password"`
	CNPJ      string    `gorm:"column:cnpj"`
	Phone     string    `gorm:"column:phone"`
	Address   string    `gorm:"column:address"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

// TableName especifica o nome da tabela no banco de dados
func (Company) TableName() string {
	return "companies"
}
