package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/jampa_trip/cmd"
	"github.com/jampa_trip/internal/app"
	"github.com/jampa_trip/internal/app/middleware"
	"github.com/jampa_trip/internal/pkg/config"
	"github.com/jampa_trip/internal/pkg/database"
	"github.com/jampa_trip/internal/pkg/util"
	"github.com/jampa_trip/internal/pkg/webserver"
	"github.com/labstack/echo/v4"
	"github.com/swaggo/swag"
)

var err error

const (
	VersionAplication = "v1.0.0"
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

	app.Config, err = config.LoadConfig()
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Println("Versão da Aplicação: ", VersionAplication)

	app.DB, err = database.GormPostgresDatabaseNew().Init(database.GormPostgresDatabaseConfig{
		Host:     app.Config.DatabaseHost,
		Port:     app.Config.DatabasePort,
		User:     app.Config.DatabaseUser,
		Password: app.Config.DatabasePassword,
		DB:       app.Config.DatabaseName,
	})
	if err != nil {
		log.Printf("erro ao inicializar conexão com o banco de dados escrita: %s", err.Error())
		panic("erro ao inicializar conexão com o banco de dados escrita")
	}
	log.Println("OK DB Escrita")
}

func main() {
	e := echo.New()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server := webserver.EchoWebServerNew().Init(webserver.EchoWebServerConfig{
		Debug:        app.Config.Debug,
		ReadTimeout:  app.Config.HTTPServerReadTimeout,
		WriteTimeout: app.Config.HTTPServerWriteTimeout,
		IDleTimeout:  app.Config.HTTPServerIdleTimeout,
	})

	middleware.SetupMiddlewares(e)

	cmd.ConfigureRoutes(e)

	go func() {
		if err := server.Start(app.Config.HTTPServerPort); err != nil {
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
