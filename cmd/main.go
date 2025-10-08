package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/jampa_trip/pkg/config"
	"github.com/jampa_trip/pkg/database"
	"github.com/jampa_trip/pkg/middleware"
	"github.com/jampa_trip/pkg/util"
	"github.com/jampa_trip/pkg/webserver"
	"github.com/swaggo/swag"
)

var err error

const (
	VersionApplication = "v1.0.0"
)

type swagger struct{}

// ReadDoc - carrega o arquivo do swagger
func (s *swagger) ReadDoc() string {
	currentDir, _ := os.Getwd()
	doc, _ := os.ReadFile(fmt.Sprintf("%s/docs/%s", currentDir, "swagger.yaml"))
	return string(doc)
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("[JAMPA-TRIP] ")

	currentDir, _ := os.Getwd()
	util.ParseSwagger(fmt.Sprintf("%s/docs/%s", currentDir, "index.yaml"))
	swag.Register(swag.Name, &swagger{})

	os.Setenv("VERSION_APPLICATION", VersionApplication)

	database.Config, err = config.LoadConfig()
	if err != nil {
		log.Fatal(err)
		return
	}

	database.DB, err = database.GormPostgresDatabaseNew().Init(database.GormPostgresDatabaseConfig{
		Host:     database.Config.DatabaseHost,
		Port:     database.Config.DatabasePort,
		User:     database.Config.DatabaseUser,
		Password: database.Config.DatabasePassword,
		DB:       database.Config.DatabaseName,
		Logger:   database.Config.DatabaseLog,
	})
	if err != nil {
		log.Fatalf("erro ao inicializar conexão com o banco de dados: %s", err.Error())
	}

	database.RedisClientNew()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server := webserver.EchoWebServerNew().Init(webserver.EchoWebServerConfig{
		Debug:        database.Config.Debug,
		ReadTimeout:  database.Config.HTTPServerReadTimeout,
		WriteTimeout: database.Config.HTTPServerWriteTimeout,
		IDleTimeout:  database.Config.HTTPServerIdleTimeout,
	})

	middleware.SetupMiddlewares(server)

	ConfigureRoutes(server)

	log.Printf("📖 Documentação da API disponível em: http://localhost%s/docs/", database.Config.HTTPServerPort)

	go func() {
		if err := server.Start(database.Config.HTTPServerPort); err != nil {
			server.Logger.Fatalf("Finalizando servidor de aplicação: %s", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 10*time.Second)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		server.Logger.Fatalf(err.Error())
	}
}
