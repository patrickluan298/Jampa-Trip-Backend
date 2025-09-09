package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// GormPostgresDatabaseConfig - objeto com as configurações do banco
type GormPostgresDatabaseConfig struct {
	Host                  string
	Port                  string
	User                  string
	Password              string
	DB                    string
	Logger                string
	MaxLifetimeConnection string
}

// GormPostgresDatabase - objeto de contexto
type GormPostgresDatabase struct {
}

// GormPostgresDatabaseNew - construtor do objeto
func GormPostgresDatabaseNew() *GormPostgresDatabase {
	return &GormPostgresDatabase{}
}

// Init - inicializa uma conexão com o banco de dados
func (receiver GormPostgresDatabase) Init(config GormPostgresDatabaseConfig) (database *gorm.DB, err error) {
	sql, err := sql.Open("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable binary_parameters=yes connect_timeout=5",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DB,
	))

	if err != nil {
		return
	}

	logOutput, err := os.OpenFile(config.Logger, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		return
	}

	database, err = gorm.Open(postgres.New(postgres.Config{Conn: sql}), &gorm.Config{
		AllowGlobalUpdate: false,
		Logger: logger.New(
			log.New(logOutput, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: time.Second,
				LogLevel:      logger.Error,
				Colorful:      true,
			},
		),
	})

	if err != nil {
		return
	}

	DB, err := database.DB()

	if err != nil {
		return
	}

	maxLifetimeConnection, _ := strconv.Atoi(config.MaxLifetimeConnection)
	DB.SetConnMaxLifetime(time.Duration(maxLifetimeConnection) * time.Second)

	return
}
