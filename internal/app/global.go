package app

import (
	"github.com/jampa_trip/internal/pkg/config"
	"gorm.io/gorm"
)

var (
	// Conexão com o banco de dados
	DB *gorm.DB

	// Configurações
	Config *config.Config
)
