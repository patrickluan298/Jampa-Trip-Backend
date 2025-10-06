package model

import "time"

// Client representa a entidade de cliente
type Client struct {
	ID        int       `gorm:"column:id;primaryKey"`
	Name      string    `gorm:"column:name"`
	Email     string    `gorm:"column:email"`
	Password  string    `gorm:"column:password"`
	CPF       string    `gorm:"column:cpf"`
	Phone     string    `gorm:"column:phone"`
	BirthDate time.Time `gorm:"column:birth_date"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

// TableName especifica o nome da tabela no banco de dados
func (Client) TableName() string {
	return "clients"
}
