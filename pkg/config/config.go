package config

import (
	"net/http"
	"os"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jampa_trip/pkg/util"
)

// Config - estrutura de configuração
type Config struct {
	Debug                             string
	HTTPServerReadTimeout             string
	HTTPServerWriteTimeout            string
	HTTPServerIdleTimeout             string
	HTTPServerPort                    string
	DatabaseHost                      string
	DatabasePort                      string
	DatabaseName                      string
	DatabaseUser                      string
	DatabasePassword                  string
	DatabasePoolMaxLifetimeConnection string
	DatabaseLog                       string

	// Mercado Pago
	MercadoPagoAccessToken   string
	MercadoPagoPublicKey     string
	MercadoPagoWebhookSecret string
	MercadoPagoEnvironment   string
	MercadoPagoBaseURL       string
}

// Validate - valida os parâmetros da requisição
func (receiver Config) Validate() (err error) {
	err = validation.ValidateStruct(&receiver,
		validation.Field(&receiver.HTTPServerReadTimeout, validation.Required),
		validation.Field(&receiver.HTTPServerWriteTimeout, validation.Required),
		validation.Field(&receiver.HTTPServerIdleTimeout, validation.Required),
		validation.Field(&receiver.HTTPServerPort, validation.Required),
		validation.Field(&receiver.DatabaseHost, validation.Required),
		validation.Field(&receiver.DatabasePort, validation.Required),
		validation.Field(&receiver.DatabaseName, validation.Required),
		validation.Field(&receiver.DatabaseUser, validation.Required),
		validation.Field(&receiver.DatabasePassword, validation.Required),

		// Validações do Mercado Pago
		validation.Field(&receiver.MercadoPagoAccessToken, validation.Required),
		validation.Field(&receiver.MercadoPagoPublicKey, validation.Required),
		validation.Field(&receiver.MercadoPagoEnvironment, validation.Required, validation.In("sandbox", "production")),
		validation.Field(&receiver.MercadoPagoBaseURL, validation.Required),
	)
	return
}

// LoadConfig - carrega as configurações do arquivo .env
func LoadConfig() (config *Config, err error) {
	config = &Config{
		Debug:                             os.Getenv("DEBUG"),
		HTTPServerReadTimeout:             os.Getenv("HTTP_SERVER_READ_TIMEOUT"),
		HTTPServerWriteTimeout:            os.Getenv("HTTP_SERVER_WRITE_TIMEOUT"),
		HTTPServerIdleTimeout:             os.Getenv("HTTP_SERVER_IDLE_TIMEOUT"),
		HTTPServerPort:                    os.Getenv("HTTP_SERVER_PORT"),
		DatabaseHost:                      os.Getenv("DATABASE_POSTGRES_HOST"),
		DatabasePort:                      os.Getenv("DATABASE_POSTGRES_PORT"),
		DatabaseName:                      os.Getenv("DATABASE_POSTGRES_NAME"),
		DatabaseUser:                      os.Getenv("DATABASE_POSTGRES_USER"),
		DatabasePassword:                  os.Getenv("DATABASE_POSTGRES_PASSWORD"),
		DatabasePoolMaxLifetimeConnection: os.Getenv("DATABASE_POSTGRES_POOL_MAX_LIFETIME_CONNECTION"),
		DatabaseLog:                       os.Getenv("DATABASE_POSTGRES_LOG"),

		// Mercado Pago
		MercadoPagoAccessToken:   os.Getenv("MERCADO_PAGO_ACCESS_TOKEN"),
		MercadoPagoPublicKey:     os.Getenv("MERCADO_PAGO_PUBLIC_KEY"),
		MercadoPagoWebhookSecret: os.Getenv("MERCADO_PAGO_WEBHOOK_SECRET"),
		MercadoPagoEnvironment:   os.Getenv("MERCADO_PAGO_ENVIRONMENT"),
		MercadoPagoBaseURL:       os.Getenv("MERCADO_PAGO_BASE_URL"),
	}

	if err = config.Validate(); err != nil {
		err = util.WrapError("erro campos obrigatório(s) não informado(s) na variáveis de ambiente", err, http.StatusBadRequest)
		return
	}
	return
}
