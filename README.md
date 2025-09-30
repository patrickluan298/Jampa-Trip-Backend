# Jampa Trip - Backend

[![Go Version](https://img.shields.io/badge/Go-1.23.5-blue.svg)](https://golang.org/)
[![Echo Framework](https://img.shields.io/badge/Echo-v4.13.4-green.svg)](https://echo.labstack.com/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-blue.svg)](https://postgresql.org/)
[![Docker](https://img.shields.io/badge/Docker-Supported-blue.svg)](https://docker.com/)
[![Mercado Pago](https://img.shields.io/badge/Mercado%20Pago-Integrated-green.svg)](https://mercadopago.com.br/)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

## 📋 Descrição

O **Jampa Trip Backend** é uma API REST desenvolvida em Go que serve como backend para uma aplicação mobile de turismo. O projeto foi desenvolvido como parte de um TCC (Trabalho de Conclusão de Curso) do curso de Ciência da Computação, focando na gestão de fornecedores de serviços turísticos e clientes.

A aplicação oferece uma arquitetura limpa e escalável utilizando o framework Echo, GORM para ORM, PostgreSQL como banco de dados e integração completa com o Mercado Pago para processamento de pagamentos.

## 🚀 Funcionalidades

### 👥 Gestão de Usuários
- **Clientes**: Cadastro, login, atualização e listagem de clientes
- **Empresas**: Cadastro, login, atualização e listagem de empresas fornecedoras de serviços turísticos

### 💳 Sistema de Pagamentos
- **Integração com Mercado Pago**: Processamento completo de pagamentos
- **Múltiplos métodos**: Cartão de crédito, débito, PIX e boleto
- **Gestão de status**: Controle completo do ciclo de vida dos pagamentos
- **Webhooks**: Notificações automáticas de mudanças de status

### 🏗️ Arquitetura
- **Clean Architecture**: Separação clara de responsabilidades
- **API RESTful**: Endpoints bem estruturados e documentados
- **Documentação Swagger**: Interface interativa para testes da API
- **Health Checks**: Monitoramento de saúde da aplicação

## 🚀 Instalação

### Pré-requisitos

- [Go 1.23.5+](https://golang.org/dl/)
- [Docker](https://docs.docker.com/get-docker/) e [Docker Compose](https://docs.docker.com/compose/install/)
- [Make](https://www.gnu.org/software/make/) (opcional, para usar os comandos do Makefile)

### Instalação Local

1. **Clone o repositório:**
```bash
git clone https://github.com/jampa_trip/Jampa-Trip-Backend.git
cd Jampa-Trip-Backend
```

2. **Instale as dependências:**
```bash
go mod download
```

3. **Configure as variáveis de ambiente:**
```bash
# Configurações básicas da aplicação
export DEBUG=false
export HTTP_SERVER_READ_TIMEOUT=20
export HTTP_SERVER_WRITE_TIMEOUT=60
export HTTP_SERVER_IDLE_TIMEOUT=120
export HTTP_SERVER_PORT=:1450

# Configurações do banco de dados
export DATABASE_POSTGRES_HOST=localhost
export DATABASE_POSTGRES_PORT=5432
export DATABASE_POSTGRES_NAME=jampa_trip_db
export DATABASE_POSTGRES_USER=jampa_trip_user
export DATABASE_POSTGRES_PASSWORD=jampa_trip_password
export DATABASE_POSTGRES_POOL_MAX_LIFETIME_CONNECTION=300
export DATABASE_POSTGRES_LOG=""

# Configurações do Mercado Pago
export MERCADO_PAGO_ACCESS_TOKEN=your_access_token_here
export MERCADO_PAGO_PUBLIC_KEY=your_public_key_here
export MERCADO_PAGO_WEBHOOK_SECRET=your_webhook_secret_here
export MERCADO_PAGO_ENVIRONMENT=sandbox
export MERCADO_PAGO_BASE_URL=https://api.mercadopago.com
```

4. **Execute o banco de dados PostgreSQL:**
```bash
make docker-dev-up
```

5. **Execute a aplicação:**
```bash
./run.sh
```

## 🐳 Instalação com Docker

### Usando Docker Compose (Recomendado)

O projeto inclui um `Makefile` com comandos pré-configurados para facilitar o desenvolvimento:

```bash
make docker-dev-up          # Inicia os containers em background
make docker-dev-build       # Faz o build e inicia os containers
make docker-dev-logs        # Exibe logs em tempo real
make docker-dev-build-logs  # Build + inicia containers + exibe logs
make docker-dev-stop        # Para containers sem removê-los
make docker-dev-down        # Para e remove containers/volumes
make docker-dev-volume-remove # Remove o volume do banco de dados
```

### Comandos Docker Manuais

```bash
docker-compose -f deployments/docker-compose.yaml up --build -d  # Build da aplicação
docker-compose -f deployments/docker-compose.yaml logs -f        # Ver logs
docker-compose -f deployments/docker-compose.yaml down           # Parar serviços
```

## 📖 Documentação da API

A documentação Swagger está disponível em:
- **Local:** `http://localhost:1450/swagger/index.html`
- **Arquivos:** `docs/` (formato OpenAPI 3.0.3)

## 🏗️ Estrutura do Projeto

```
.
├── build/                    # Arquivos de build e Dockerfile
├── cmd/                      # Ponto de entrada da aplicação
│   ├── app/
│   │   └── main.go
│   └── routes.go
├── deployments/              # Configurações de deploy
│   ├── docker-compose.yaml
│   └── init.sql
├── docs/                     # Documentação da API (OpenAPI/Swagger)
│   ├── components/
│   ├── paths/
│   └── index.yaml
├── internal/                 # Código interno da aplicação
│   ├── app/                  # Lógica de negócio
│   │   ├── contract/         # Contratos de request/response
│   │   ├── dto/              # Data Transfer Objects
│   │   ├── handler/          # Handlers HTTP
│   │   ├── middleware/       # Middlewares
│   │   ├── model/            # Modelos de dados
│   │   ├── query/            # Queries customizadas
│   │   ├── repository/       # Camada de acesso a dados
│   │   ├── service/          # Lógica de negócio
│   │   ├── types/            # Tipos customizados
│   │   └── validation/       # Validações
│   └── pkg/                  # Pacotes utilitários
│       ├── config/           # Configurações
│       ├── database/         # Conexão com banco
│       ├── mercadopago/      # Cliente Mercado Pago
│       ├── util/             # Utilitários
│       └── webserver/        # Servidor web
├── tests/                    # Testes automatizados
├── go.mod
├── go.sum
├── Makefile                  # Comandos de automação
├── run.sh                    # Script de execução
└── MERCADO_PAGO_SETUP.md     # Documentação do Mercado Pago
```

### Arquitetura

O projeto segue os princípios da **Clean Architecture** com as seguintes camadas:

- **Handler:** Recebe requisições HTTP e valida dados de entrada
- **Service:** Contém a lógica de negócio
- **Repository:** Gerencia acesso aos dados
- **Model:** Define as entidades do domínio
- **Contract:** Define contratos de entrada e saída
- **Validation:** Validação de dados de entrada

## ⚙️ Configuração

### Variáveis de Ambiente

| Variável | Descrição | Padrão | Obrigatório |
|----------|-----------|---------|-------------|
| `DEBUG` | Modo debug | `false` | Não |
| `HTTP_SERVER_READ_TIMEOUT` | Timeout de leitura HTTP (segundos) | `20` | Sim |
| `HTTP_SERVER_WRITE_TIMEOUT` | Timeout de escrita HTTP (segundos) | `60` | Sim |
| `HTTP_SERVER_IDLE_TIMEOUT` | Timeout de idle HTTP (segundos) | `120` | Sim |
| `HTTP_SERVER_PORT` | Porta do servidor HTTP | `:1450` | Sim |
| `DATABASE_POSTGRES_HOST` | Host do PostgreSQL | - | Sim |
| `DATABASE_POSTGRES_PORT` | Porta do PostgreSQL | - | Sim |
| `DATABASE_POSTGRES_NAME` | Nome do banco de dados | - | Sim |
| `DATABASE_POSTGRES_USER` | Usuário do banco | - | Sim |
| `DATABASE_POSTGRES_PASSWORD` | Senha do banco | - | Sim |
| `DATABASE_POSTGRES_POOL_MAX_LIFETIME_CONNECTION` | Tempo de vida da conexão (segundos) | `300` | Não |
| `DATABASE_POSTGRES_LOG` | Caminho do log do banco | - | Não |
| `MERCADO_PAGO_ACCESS_TOKEN` | Token de acesso do Mercado Pago | - | Sim (para pagamentos) |
| `MERCADO_PAGO_PUBLIC_KEY` | Chave pública do Mercado Pago | - | Sim (para pagamentos) |
| `MERCADO_PAGO_WEBHOOK_SECRET` | Chave secreta para webhooks | - | Não |
| `MERCADO_PAGO_ENVIRONMENT` | Ambiente (sandbox/production) | `sandbox` | Não |
| `MERCADO_PAGO_BASE_URL` | URL base da API do Mercado Pago | `https://api.mercadopago.com` | Não |

### Configuração do Banco de Dados

O banco PostgreSQL é configurado automaticamente via Docker Compose com:

- **Database:** `jampa_trip_db`
- **User:** `jampa_trip_user`
- **Password:** `jampa_trip_password`
- **Port:** `6432` (mapeada para `5432` no container)

### Configuração do Mercado Pago

Para configurar o Mercado Pago, consulte o arquivo `MERCADO_PAGO_SETUP.md` que contém instruções detalhadas sobre:

1. Como obter as credenciais necessárias
2. Configuração das variáveis de ambiente
3. Estrutura da integração implementada
4. Status de pagamento suportados
5. Métodos de pagamento disponíveis

## 🧪 Testes

```bash
# Executar todos os testes
go test ./...

# Executar testes com coverage
go test -cover ./...

# Executar testes em modo verbose
go test -v ./...

# Executar testes de um pacote específico
go test ./internal/app/service/...
```

> **Nota:** O diretório `tests/` está presente mas ainda não contém testes implementados. Esta é uma área para desenvolvimento futuro.

##  Deploy

### Deploy Local com Docker

```bash
# Build e start dos containers
make docker-dev-build

# Verificar se os serviços estão rodando
docker-compose -f deployments/docker-compose.yaml ps

# Verificar logs
make docker-dev-logs
```

### Deploy em Produção

1. **Configure as variáveis de ambiente de produção**
2. **Ajuste o `docker-compose.yaml` para produção**
3. **Execute o build:**
```bash
docker-compose -f deployments/docker-compose.yaml up --build -d
```

### Health Checks

A aplicação inclui health checks configurados:

- **Aplicação:** `http://localhost:1450/health-check`
- **Banco de dados:** Verificação automática via `pg_isready`
- **Docker:** Health checks configurados nos containers

## 💳 Integração com Mercado Pago

O projeto inclui integração completa com o Mercado Pago para processamento de pagamentos:

### Funcionalidades Implementadas
- ✅ Criação de Orders
- ✅ Criação de Pagamentos (Cartão de Crédito/Débito)
- ✅ Criação de Pagamentos PIX
- ✅ Consulta de Pagamentos
- ✅ Cancelamento de Pagamentos
- ✅ Tratamento de Erros da API

### Status de Pagamento Suportados
- `pending` - Pendente
- `approved` - Aprovado
- `authorized` - Autorizado
- `in_process` - Em Processamento
- `in_mediation` - Em Mediação
- `rejected` - Rejeitado
- `cancelled` - Cancelado
- `refunded` - Reembolsado
- `charged_back` - Estornado

### Métodos de Pagamento Suportados
- `credit_card` - Cartão de Crédito
- `debit_card` - Cartão de Débito
- `pix` - PIX
- `bolbradesco` - Boleto

## 🛠️ Tecnologias Utilizadas

- **[Go 1.23.5](https://golang.org/)** - Linguagem de programação
- **[Echo v4](https://echo.labstack.com/)** - Framework web
- **[GORM](https://gorm.io/)** - ORM para Go
- **[PostgreSQL 15](https://postgresql.org/)** - Banco de dados
- **[Docker](https://docker.com/)** - Containerização
- **[Swagger/OpenAPI 3.0.3](https://swagger.io/)** - Documentação da API
- **[Ozzo Validation](https://github.com/go-ozzo/ozzo-validation)** - Validação de dados
- **[Mercado Pago API](https://www.mercadopago.com.br/developers/)** - Processamento de pagamentos

## 📚 Documentação Adicional

- [Configuração do Mercado Pago](MERCADO_PAGO_SETUP.md)
- [Documentação da API](docs/)
- [Swagger UI](http://localhost:1450/swagger/index.html)

---

**Desenvolvido como parte do TCC do curso de Ciência da Computação** 🎓
