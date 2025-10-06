package middleware

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/jampa_trip/pkg/util"
	"github.com/labstack/echo/v4"
)

// Logging - retorna um middleware do Echo que loga dados estruturados de requisição/resposta e erros.
func Logging() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			// Pula o endpoint de health check
			if c.Request().RequestURI == "/health-check" {
				return next(c)
			}

			start := time.Now()

			// Lê o corpo da requisição
			bodyBytes, err := io.ReadAll(c.Request().Body)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Erro ao processar a requisição")
			}
			c.Request().Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			// Wrap response writer para capturar o corpo da resposta
			responseBody := new(strings.Builder)
			originalWriter := c.Response().Writer
			wrapper := &responseWriterWrapper{
				ResponseWriter: originalWriter,
				body:           responseBody,
				statusCode:     http.StatusOK,
			}
			c.Response().Writer = wrapper

			// Executa o próximo handler
			err = next(c)

			// Calcula as métricas
			latency := time.Since(start).Milliseconds()
			statusCode := wrapper.statusCode
			severity := mapStatusToSeverity(statusCode)

			// Obtém o corpo da resposta (descomprime se gzipado)
			responseBodyStr := responseBody.String()
			if c.Response().Header().Get("Content-Encoding") == "gzip" {
				if decompressed, err := decompressGzip(responseBodyStr); err == nil {
					responseBodyStr = decompressed
				}
			}

			// Constroi o payload do log
			logPayload := LogPayload{
				HTTPRequest: HTTPRequestLog{
					Latency:       fmt.Sprintf("%dms", latency),
					RemoteIP:      c.RealIP(),
					RequestMethod: c.Request().Method,
					RequestURL:    c.Request().RequestURI,
					ResponseSize:  c.Response().Size,
					Status:        statusCode,
				},
				JsonPayload: LogJsonPayload{
					Request: RequestPayload{
						Body:    string(bodyBytes),
						Headers: c.Request().Header,
					},
					Response: ResponsePayload{
						Status: statusCode,
						Body:   responseBodyStr,
					},
				},
				Message:   fmt.Sprintf("[%s] %s - %d (%dms)", c.Request().Method, c.Request().RequestURI, statusCode, latency),
				Timestamp: time.Now().Format(time.RFC3339),
			}

			// Adiciona detalhes de erro se presente
			if wrappedErr, ok := c.Get("errorDetails").(*util.AppError); ok {
				logPayload.JsonPayload.Error = AppErrorDetails{
					Message:    wrappedErr.Msg,
					Err:        wrappedErr.Err,
					File:       wrappedErr.File,
					Line:       wrappedErr.Line,
					Function:   wrappedErr.FuncName,
					StatusCode: wrappedErr.StatusCode,
				}
			} else if err != nil {
				logPayload.JsonPayload.Error = GenericError{
					Message: err.Error(),
					Details: fmt.Sprintf("%+v", err),
				}
			}

			// Imprime o log estruturado no console
			jsonLog, jsonErr := json.MarshalIndent(logPayload, "", "  ")
			if jsonErr != nil {
				log.Printf("[JAMPA-TRIP] Erro ao converter log para JSON: %v", jsonErr)
			} else {
				log.Printf("[JAMPA-TRIP] [%s] %s\n%s", time.Now().Format("2006-01-02 15:04:05"), severity, string(jsonLog))
			}

			return err
		}
	}
}

// decompressGzip - descomprime o conteúdo gzipado
func decompressGzip(compressed string) (string, error) {
	reader, err := gzip.NewReader(strings.NewReader(compressed))
	if err != nil {
		return "", err
	}
	defer reader.Close()

	decompressed, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(decompressed), nil
}

// mapStatusToSeverity - mapeia os códigos de status HTTP para níveis de severidade apropriados.
func mapStatusToSeverity(status int) string {
	switch {
	case status >= 200 && status < 300:
		return "INFO"
	case status >= 400 && status < 500:
		return "WARNING"
	case status >= 500:
		return "ERROR"
	default:
		return "DEFAULT"
	}
}

// responseWriterWrapper - wrapper da resposta para capturar o corpo da resposta e o código de status.
type responseWriterWrapper struct {
	http.ResponseWriter
	body       *strings.Builder
	statusCode int
}

func (w *responseWriterWrapper) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *responseWriterWrapper) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *responseWriterWrapper) Header() http.Header {
	return w.ResponseWriter.Header()
}

// RecoveryMiddleware - cria um middleware de recuperação com logging estruturado.
func RecoveryMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			defer func() {
				if r := recover(); r != nil {
					// Captura o stack trace
					stackTrace := string(debug.Stack())

					logEntry := PanicLogPayload{
						Error:      fmt.Sprintf("%v", r),
						StackTrace: stackTrace,
						HTTPRequest: PanicRequestInfo{
							Method:  c.Request().Method,
							URL:     c.Request().URL.String(),
							Remote:  c.RealIP(),
							Headers: c.Request().Header,
						},
						Message:   "Panic recovered",
						Timestamp: time.Now().Format(time.RFC3339),
					}

					// Imprime o log de panic no console
					jsonLog, jsonErr := json.MarshalIndent(logEntry, "", "  ")
					if jsonErr != nil {
						log.Printf("[JAMPA-TRIP] Erro ao converter panic log para JSON: %v", jsonErr)
					} else {
						log.Printf("[JAMPA-TRIP] [CRITICAL] %s\n%s", time.Now().Format("2006-01-02 15:04:05"), string(jsonLog))
					}

					// Retorna uma resposta de erro genérica para o cliente
					c.JSON(http.StatusInternalServerError, map[string]string{
						"error": "Internal Server Error",
					})
					err = echo.NewHTTPError(http.StatusInternalServerError, r)
				}
			}()
			return next(c)
		}
	}
}

// LogPayload - representa a estrutura principal do log
type LogPayload struct {
	HTTPRequest HTTPRequestLog `json:"httpRequest"`
	JsonPayload LogJsonPayload `json:"jsonPayload"`
	Message     string         `json:"message"`
	Timestamp   string         `json:"timestamp"`
}

// HTTPRequestLog - contém informações da requisição HTTP
type HTTPRequestLog struct {
	Latency       string `json:"latency"`
	RemoteIP      string `json:"remoteIp"`
	RequestMethod string `json:"requestMethod"`
	RequestURL    string `json:"requestUrl"`
	ResponseSize  int64  `json:"responseSize"`
	Status        int    `json:"status"`
}

// LogJsonPayload - contém informações detalhadas da requisição/resposta
type LogJsonPayload struct {
	Request  RequestPayload  `json:"request"`
	Response ResponsePayload `json:"response"`
	Error    interface{}     `json:"error,omitempty"`
}

// RequestPayload - contém informações detalhadas da requisição
type RequestPayload struct {
	Body    string      `json:"body"`
	Headers http.Header `json:"headers"`
}

// ResponsePayload - contém informações detalhadas da resposta
type ResponsePayload struct {
	Status int    `json:"status"`
	Body   string `json:"body"`
}

// AppErrorDetails - contém informações detalhadas do erro
type AppErrorDetails struct {
	Message    string `json:"message"`
	Err        error  `json:"err,omitempty"`
	File       string `json:"file"`
	Line       int    `json:"line"`
	Function   string `json:"function"`
	StatusCode int    `json:"statusCode"`
}

// GenericError - contém informações detalhadas do erro genérico
type GenericError struct {
	Message string `json:"message"`
	Details string `json:"details"`
}

// PanicLogPayload - contém informações detalhadas do log de panic
type PanicLogPayload struct {
	Error       string           `json:"error"`
	StackTrace  string           `json:"stackTrace"`
	HTTPRequest PanicRequestInfo `json:"httpRequest"`
	Message     string           `json:"message"`
	Timestamp   string           `json:"timestamp"`
}

// PanicRequestInfo - contém informações detalhadas da requisição durante o panic
type PanicRequestInfo struct {
	Method  string      `json:"method"`
	URL     string      `json:"url"`
	Remote  string      `json:"remote"`
	Headers http.Header `json:"headers"`
}
