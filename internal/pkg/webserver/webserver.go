package webserver

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const PathHealthCheck = "/health-check"

// EchoWebServerConfig objeto com as configurações do webserver
type EchoWebServerConfig struct {
	Debug        string
	ReadTimeout  string
	WriteTimeout string
	IDleTimeout  string
	Routes       map[string]map[string]echo.HandlerFunc
}

// RouteConfig representa a configuração de uma rota
type RouteConfig struct {
	Handler         echo.HandlerFunc
	NeedsMiddleware bool
}

// EchoWebServer objeto de contexto
type EchoWebServer struct {
}

// EchoWebServerNew construtor do objeto
func EchoWebServerNew() *EchoWebServer {
	return &EchoWebServer{}
}

// Init inicializa o webserver do echo framework
func (receiver EchoWebServer) Init(config EchoWebServerConfig) (serverEcho *echo.Echo) {

	server := echo.New()
	server.Use(middleware.Recover())
	server.Use(middleware.Secure())

	debug, _ := strconv.ParseBool(config.Debug)
	readTimeout, _ := strconv.Atoi(config.ReadTimeout)
	writeTimeout, _ := strconv.Atoi(config.WriteTimeout)
	idleTimeout, _ := strconv.Atoi(config.IDleTimeout)

	server.Debug = debug
	server.Server.ReadTimeout = time.Duration(readTimeout) * time.Second
	server.Server.WriteTimeout = time.Duration(writeTimeout) * time.Second
	server.Server.IdleTimeout = time.Duration(idleTimeout) * time.Second
	server.HTTPErrorHandler = receiver.erroHandler

	return server
}

func (receiver EchoWebServer) erroHandler(err error, ctx echo.Context) {
	ctx.Logger().Error(err)

	he, ok := err.(*echo.HTTPError)
	if ok {
		if he.Internal != nil {
			if herr, ok := he.Internal.(*echo.HTTPError); ok {
				he = herr
			}
		}
	} else {
		he = &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: http.StatusText(http.StatusInternalServerError),
		}
	}

	code := he.Code
	message := he.Message
	if ctx.Echo().Debug {
		message = echo.Map{
			"message": err.Error(),
		}
	} else if m, ok := message.(string); ok {
		message = echo.Map{
			"message": m,
		}
	}

	if !ctx.Response().Committed {
		if ctx.Request().Method == http.MethodHead {
			err = ctx.NoContent(he.Code)
		} else {
			err = ctx.JSON(code, message)
		}

		if err != nil {
			ctx.Echo().Logger.Error(err)
		}
	}
}
