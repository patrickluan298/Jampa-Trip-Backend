package model

import (
	"fmt"
	"time"

	"github.com/lib/pq"
)

// Tour representa a entidade de passeio
type Tour struct {
	ID            int            `gorm:"column:id;primaryKey;autoIncrement"`
	CompanyID     int            `gorm:"column:company_id;not null"`
	Name          string         `gorm:"column:name;not null"`
	Dates         pq.StringArray `gorm:"column:dates;type:text[]"`
	DepartureTime string         `gorm:"column:departure_time"`
	ArrivalTime   string         `gorm:"column:arrival_time"`
	MaxPeople     int            `gorm:"column:max_people;default:1"`
	Description   string         `gorm:"column:description"`
	Images        pq.StringArray `gorm:"column:images;type:text[]"`
	Price         float64        `gorm:"column:price;type:decimal(10,2);default:0"`
	CreatedAt     time.Time      `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt     time.Time      `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP"`

	Company Company `gorm:"foreignKey:CompanyID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// TableName especifica o nome da tabela no banco de dados
func (Tour) TableName() string {
	return "tours"
}

// Métodos de validação para Tour
func (t *Tour) IsValid() bool {
	return t.CompanyID > 0 && len(t.Name) >= 3 && t.MaxPeople > 0 && t.Price >= 0
}

func (t *Tour) HasValidTimes() bool {
	if t.DepartureTime != "" && len(t.DepartureTime) != 5 {
		return false
	}
	if t.ArrivalTime != "" && len(t.ArrivalTime) != 5 {
		return false
	}
	return true
}

func (t *Tour) HasValidImages() bool {
	for _, img := range t.Images {
		if img == "" {
			return false
		}
	}
	return true
}

func (t *Tour) GetFormattedPrice() string {
	return fmt.Sprintf("%.2f", t.Price)
}

func (t *Tour) GetFormattedDates() []string {
	var formattedDates []string
	for _, date := range t.Dates {
		if date != "" {
			formattedDates = append(formattedDates, date)
		}
	}
	return formattedDates
}
