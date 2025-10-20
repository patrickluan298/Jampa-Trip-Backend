# Jampa Trip - Backend

[![Go Version](https://img.shields.io/badge/Go-1.23.5-blue.svg)](https://golang.org/)
[![Echo Framework](https://img.shields.io/badge/Echo-v4.13.4-green.svg)](https://echo.labstack.com/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-blue.svg)](https://postgresql.org/)
[![Redis](https://img.shields.io/badge/Redis-7-red.svg)](https://redis.io/)
[![Docker](https://img.shields.io/badge/Docker-Supported-blue.svg)](https://docker.com/)
[![JWT](https://img.shields.io/badge/JWT-Authentication-orange.svg)](https://jwt.io/)
[![Mercado Pago](https://img.shields.io/badge/Mercado%20Pago-Integrated-green.svg)](https://mercadopago.com.br/)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

## 📋 Descrição

O **Jampa Trip Backend** é uma API REST desenvolvida em Go que serve como backend para uma aplicação mobile de turismo. O projeto foi desenvolvido como parte de um TCC (Trabalho de Conclusão de Curso) do curso de Ciência da Computação, focando na gestão de fornecedores de serviços turísticos e clientes.

A aplicação oferece uma arquitetura limpa e escalável utilizando o framework Echo, GORM para ORM, PostgreSQL como banco de dados, Redis para cache e autenticação, e integração completa com o Mercado Pago para processamento de pagamentos.

## 🚀 Funcionalidades

### 👥 Gestão de Usuários
- **Clientes**: Cadastro, login, atualização e listagem de clientes
- **Empresas**: Cadastro, login, atualização e listagem de empresas fornecedoras de serviços turísticos
- **Autenticação JWT**: Sistema completo de autenticação com access e refresh tokens
- **Segurança**: Middleware de autenticação em todas as rotas protegidas

### 🔐 Sistema de Autenticação
- **JWT Tokens**: Access tokens (15 min) e refresh tokens (7 dias)
- **Redis Storage**: Armazenamento seguro de tokens com TTL
- **Token Rotation**: Renovação automática de tokens
- **Middleware Global**: Proteção automática de rotas
- **Validação Dupla**: Tokens validados tanto no JWT quanto no Redis

### 💳 Sistema de Pagamentos
- **Integração com Mercado Pago**: Processamento completo de pagamentos
- **Múltiplos métodos**: Cartão de crédito, débito, e PIX
- **Gestão de status**: Controle completo do ciclo de vida dos pagamentos
- **Autorização de cartões**: Sistema de pré-autorização para cartões de crédito
- **Gestão de cartões**: Cadastro e gerenciamento de cartões dos clientes

### 🏗️ Arquitetura
- **Clean Architecture**: Separação clara de responsabilidades
- **API RESTful**: Endpoints bem estruturados e documentados
- **Documentação Swagger**: Interface interativa para testes da API
- **Health Checks**: Monitoramento de saúde da aplicação
- **Middleware de Logging**: Sistema completo de logs estruturados

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

# Configurações JWT
export JWT_SECRET=jampa_trip_jwt_secret_key_2024_very_secure
export JWT_ACCESS_TOKEN_EXPIRATION=15m
export JWT_REFRESH_TOKEN_EXPIRATION=168h

# Configurações Redis
export REDIS_HOST=localhost
export REDIS_PORT=6379
export REDIS_PASSWORD=""
export REDIS_DB=0

# Configurações do Mercado Pago
export MERCADO_PAGO_ACCESS_TOKEN=your_access_token_here
export MERCADO_PAGO_PUBLIC_KEY=your_public_key_here
export MERCADO_PAGO_WEBHOOK_SECRET=your_webhook_secret_here
export MERCADO_PAGO_ENVIRONMENT=sandbox
export MERCADO_PAGO_BASE_URL=https://api.mercadopago.com
```

4. **Execute os serviços:**
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
make docker-dev-up             # Inicia os containers em background
make docker-dev-build          # Faz o build e inicia os containers
make docker-dev-logs           # Exibe logs em tempo real
make docker-dev-build-logs     # Build + inicia containers + exibe logs
make docker-dev-stop           # Para containers sem removê-los
make docker-dev-down           # Para e remove containers/volumes
make docker-dev-volume-remove  # Remove o volume do banco de dados
```

## 📖 Documentação da API

A documentação Swagger está disponível em:
- **Swagger UI:** `http://localhost:1450/docs/`
- **Arquivos:** `docs/` (formato OpenAPI 3.0.3)

## 🏗️ Estrutura do Projeto

```
.
├── cmd/                      # Ponto de entrada da aplicação
├── deployments/              # Configurações de deploy
├── docs/                     # Documentação da API (OpenAPI/Swagger)
├── internal/                 # Código interno da aplicação
│   ├── contract/             # Contratos de request/response
│   ├── handler/              # Handlers HTTP
│   ├── model/                # Modelos de dados
│   ├── query/                # Queries customizadas
│   ├── repository/           # Camada de acesso a dados
│   └── service/              # Lógica de negócio
├── pkg/                      # Pacotes utilitários
│   ├── auth/                 # Autenticação JWT
│   ├── config/               # Configurações
│   ├── database/             # Conexões com banco e Redis
│   ├── middleware/           # Middlewares HTTP
│   ├── mercadopago/          # Integração Mercado Pago
│   ├── util/                 # Utilitários
│   └── webserver/            # Servidor web
├── tests/                    # Testes automatizados
└── README.md
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
| `JWT_SECRET` | Chave secreta para JWT | - | Sim |
| `JWT_ACCESS_TOKEN_EXPIRATION` | Duração do access token | `15m` | Sim |
| `JWT_REFRESH_TOKEN_EXPIRATION` | Duração do refresh token | `168h` | Sim |
| `REDIS_HOST` | Host do Redis | - | Sim |
| `REDIS_PORT` | Porta do Redis | - | Sim |
| `REDIS_PASSWORD` | Senha do Redis | - | Não |
| `REDIS_DB` | Número do banco Redis | `0` | Não |
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

### Configuração do Redis

O Redis é configurado automaticamente via Docker Compose com:

- **Host:** `redis`
- **Port:** `6379`
- **Database:** `0`
- **Persistência:** Volume Docker para dados

### Configuração do Mercado Pago

Para configurar o Mercado Pago, consulte o arquivo `MERCADO_PAGO_SETUP.md` que contém instruções detalhadas sobre:

1. Como obter as credenciais necessárias
2. Configuração das variáveis de ambiente
3. Estrutura da integração implementada
4. Status de pagamento suportados
5. Métodos de pagamento disponíveis

## 🔐 Autenticação JWT

### Visão Geral

O sistema implementa autenticação JWT completa com:

- **Access Tokens**: Válidos por 15 minutos
- **Refresh Tokens**: Válidos por 7 dias
- **Armazenamento Redis**: Tokens armazenados com TTL
- **Middleware Global**: Proteção automática de rotas
- **Token Rotation**: Renovação segura de tokens

### Arquitetura da Autenticação

#### Componentes Principais

1. **JWT Authentication (`pkg/auth/jwt.go`)**
   - Geração de tokens (access e refresh)
   - Validação e parsing de tokens
   - Estrutura de claims personalizada

2. **Redis Token Store (`pkg/auth/redis_token_store.go`)**
   - Armazenamento de tokens no Redis
   - Validação de tokens ativos
   - Gerenciamento de TTL (Time To Live)

3. **JWT Middleware (`pkg/middleware/jwt_middleware.go`)**
   - Validação automática de tokens em rotas protegidas
   - Extração de claims do contexto
   - Tratamento de erros de autenticação

4. **Redis Client (`pkg/database/redis.go`)**
   - Cliente singleton para conexão com Redis
   - Configuração e inicialização
   - Health check

### Fluxo de Autenticação

#### 1. Login e Geração de Tokens

```mermaid
sequenceDiagram
    participant Client
    participant LoginHandler
    participant LoginService
    participant JWT
    participant Redis
    participant TokenStore

    Client->>LoginHandler: POST /login {email, password}
    LoginHandler->>LoginService: Login(request)
    LoginService->>LoginService: Validar credenciais
    LoginService->>JWT: GenerateTokenPair(userID, userType, email)
    JWT-->>LoginService: {access_token, refresh_token, expires_in}
    LoginService->>TokenStore: StoreTokenPair(userID, userType, accessToken, refreshToken)
    TokenStore->>Redis: SET access_token:userID:userType
    TokenStore->>Redis: SET refresh_token:userID:userType
    LoginService-->>LoginHandler: LoginResponse com tokens
    LoginHandler-->>Client: 200 OK {access_token, refresh_token, expires_in}
```

#### 2. Validação de Tokens (Middleware)

```mermaid
sequenceDiagram
    participant Client
    participant JWTMiddleware
    participant JWT
    participant Redis
    participant TokenStore
    participant Handler

    Client->>JWTMiddleware: Request com Authorization: Bearer <token>
    JWTMiddleware->>JWT: ValidateToken(tokenString)
    JWT-->>JWTMiddleware: Claims {userID, userType, email}
    JWTMiddleware->>TokenStore: ValidateAccessToken(userID, userType, token)
    TokenStore->>Redis: GET access_token:userID:userType
    Redis-->>TokenStore: Token armazenado
    TokenStore-->>JWTMiddleware: Token válido
    JWTMiddleware->>JWTMiddleware: Injetar claims no contexto
    JWTMiddleware->>Handler: next(c)
    Handler-->>Client: Response
```

#### 3. Refresh Token

```mermaid
sequenceDiagram
    participant Client
    participant RefreshHandler
    participant RefreshService
    participant JWT
    participant Redis
    participant TokenStore

    Client->>RefreshHandler: POST /refresh {refresh_token}
    RefreshHandler->>RefreshService: RefreshToken(request)
    RefreshService->>JWT: ValidateToken(refreshToken)
    JWT-->>RefreshService: Claims {userID, userType, email}
    RefreshService->>TokenStore: ValidateRefreshToken(userID, userType, token)
    TokenStore->>Redis: GET refresh_token:userID:userType
    Redis-->>TokenStore: Token válido
    RefreshService->>JWT: GenerateTokenPair(userID, userType, email)
    JWT-->>RefreshService: Novos tokens
    RefreshService->>TokenStore: DeleteTokens(userID, userType)
    TokenStore->>Redis: DEL access_token:userID:userType
    TokenStore->>Redis: DEL refresh_token:userID:userType
    RefreshService->>TokenStore: StoreTokenPair(userID, userType, newTokens)
    TokenStore->>Redis: SET novos tokens
    RefreshService-->>RefreshHandler: RefreshTokenResponse
    RefreshHandler-->>Client: 200 OK {novos tokens}
```

### Estrutura de Dados no Redis

#### Chaves de Armazenamento

- **Access Token**: `access_token:{userID}:{userType}`
- **Refresh Token**: `refresh_token:{userID}:{userType}`

#### Exemplo de Dados

```
access_token:123:client = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
refresh_token:123:client = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

#### TTL (Time To Live)

- **Access Token**: 15 minutos (configurável via `JWT_ACCESS_TOKEN_EXPIRATION`)
- **Refresh Token**: 7 dias (configurável via `JWT_REFRESH_TOKEN_EXPIRATION`)

### Claims do JWT

#### Estrutura de Claims

```go
type JWTClaims struct {
    UserID   int    `json:"user_id"`
    UserType string `json:"user_type"`
    Email    string `json:"email"`
    jwt.RegisteredClaims
}
```

#### Campos Padrão (RegisteredClaims)

- `exp` (ExpiresAt): Data de expiração
- `iat` (IssuedAt): Data de emissão
- `nbf` (NotBefore): Data de início de validade
- `iss` (Issuer): Emissor ("jampa-trip")
- `sub` (Subject): ID do usuário

### Rotas da API

#### Rotas Públicas (sem autenticação)

- `GET /health-check` - Health check
- `POST /jampa-trip/api/v1/login` - Login
- `POST /jampa-trip/api/v1/refresh` - Renovar tokens

#### Rotas Protegidas (com autenticação JWT)

Todas as outras rotas estão protegidas pelo middleware JWT.

### Exemplos de Uso

#### 1. Login

```bash
curl -X POST http://localhost:1450/jampa-trip/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

**Resposta:**
```json
{
  "message": "Login realizado com sucesso",
  "type": "client",
  "data": {
    "id": 123,
    "name": "João Silva",
    "email": "user@example.com"
  },
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 1703123456
}
```

#### 2. Acesso a Rota Protegida

```bash
curl -X GET http://localhost:1450/jampa-trip/api/v1/companies \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

#### 3. Renovar Tokens

```bash
curl -X POST http://localhost:1450/jampa-trip/api/v1/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }'
```

**Resposta:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 1703123456
}
```

### Tratamento de Erros

#### Códigos de Status HTTP

- **401 Unauthorized**: Token inválido, expirado ou não fornecido
- **403 Forbidden**: Token não encontrado no Redis
- **422 Unprocessable Entity**: Erro de validação nos dados
- **500 Internal Server Error**: Erro interno do servidor

#### Exemplos de Erros

```json
{
  "statusCode": 401,
  "message": "token de autorização não fornecido"
}
```

```json
{
  "statusCode": 401,
  "message": "token expirado"
}
```

```json
{
  "statusCode": 401,
  "message": "token não encontrado ou inválido"
}
```

### Segurança

#### Medidas Implementadas

1. **Tokens com TTL**: Access tokens expiram em 15 minutos
2. **Refresh Token Rotation**: Novos tokens são gerados a cada refresh
3. **Armazenamento Seguro**: Tokens são armazenados no Redis com TTL
4. **Validação Dupla**: Tokens são validados tanto no JWT quanto no Redis
5. **Chave Secreta**: JWT_SECRET é configurado via variável de ambiente

#### Boas Práticas

1. **Use HTTPS** em produção
2. **Configure JWT_SECRET** com uma chave forte e única
3. **Monitore** tentativas de acesso com tokens inválidos
4. **Implemente rate limiting** para prevenir ataques de força bruta
5. **Logs de segurança** para auditoria

### Monitoramento

#### Logs Importantes

- Conexão com Redis estabelecida
- Tokens gerados com sucesso
- Tokens validados
- Erros de autenticação
- Tokens expirados

#### Métricas Recomendadas

- Número de logins por minuto
- Taxa de tokens expirados
- Tentativas de acesso com tokens inválidos
- Tempo de resposta do Redis

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

## 🚀 Deploy

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
- **Redis:** Verificação automática via `redis-cli ping`
- **Docker:** Health checks configurados nos containers

## 💳 Integração com Mercado Pago

O projeto inclui integração completa com o Mercado Pago para processamento de pagamentos:

### Funcionalidades Implementadas

- ✅ **Criação de Orders** - Gerenciamento de pedidos
- ✅ **Pagamentos com Cartão** - Crédito e débito
- ✅ **Pagamentos PIX** - Transferência instantânea
- ✅ **Autorização de Cartões** - Pré-autorização para reservas
- ✅ **Gestão de Cartões** - Cadastro e gerenciamento
- ✅ **Consulta de Pagamentos** - Busca e visualização
- ✅ **Atualização de Status** - Controle automático
- ✅ **Tratamento de Erros** - Logs e monitoramento

### Status de Pagamento Suportados

| Status | Descrição |
|--------|-----------|
| `pending` | Pendente |
| `approved` | Aprovado |
| `authorized` | Autorizado |
| `in_process` | Em Processamento |
| `in_mediation` | Em Mediação |
| `rejected` | Rejeitado |
| `cancelled` | Cancelado |
| `refunded` | Reembolsado |
| `charged_back` | Estornado |

### Métodos de Pagamento Suportados

| Método | Descrição |
|--------|-----------|
| `credit_card` | Cartão de Crédito |
| `debit_card` | Cartão de Débito |
| `pix` | PIX |

### Configuração do Mercado Pago

#### Como Obter as Credenciais

1. **Acesse o [Painel de Desenvolvedores](https://www.mercadopago.com.br/developers/panel/credentials)**
2. **Faça login** com sua conta Mercado Pago
3. **Crie uma nova aplicação** ou selecione uma existente
4. **Copie as credenciais:**
   - **Access Token**: Token de acesso para autenticação na API
   - **Public Key**: Chave pública para uso no frontend
   - **Webhook Secret**: Chave secreta para validar webhooks

#### Ambientes

- **Sandbox**: Ambiente de testes (recomendado para desenvolvimento)
- **Production**: Ambiente de produção

#### URLs da API

- **Sandbox**: `https://api.mercadopago.com`
- **Production**: `https://api.mercadopago.com`

### Estrutura do Banco de Dados

A tabela de pagamentos é criada automaticamente pelo GORM:

```sql
CREATE TABLE pagamentos (
    id SERIAL PRIMARY KEY,
    cliente_id INTEGER NOT NULL REFERENCES clientes(id),
    empresa_id INTEGER NOT NULL REFERENCES empresas(id),
    mercado_pago_order_id VARCHAR(255) UNIQUE,
    mercado_pago_payment_id VARCHAR(255),
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    status_detail VARCHAR(100),
    valor DECIMAL(10,2) NOT NULL,
    moeda VARCHAR(3) NOT NULL DEFAULT 'BRL',
    metodo_pagamento VARCHAR(50) NOT NULL,
    descricao TEXT,
    numero_parcelas INTEGER DEFAULT 1,
    token_cartao VARCHAR(255),
    chave_pix VARCHAR(255),
    qr_code TEXT,
    -- Campos específicos para cartão
    last_four_digits VARCHAR(4),
    first_six_digits VARCHAR(6),
    payment_method_id VARCHAR(50),
    issuer_id VARCHAR(50),
    cardholder_name VARCHAR(255),
    captured BOOLEAN DEFAULT FALSE,
    transaction_amount_refunded DECIMAL(10,2) DEFAULT 0,
    -- Timestamps
    momento_criacao TIMESTAMP NOT NULL DEFAULT NOW(),
    momento_atualizacao TIMESTAMP NOT NULL DEFAULT NOW(),
    momento_aprovacao TIMESTAMP,
    momento_cancelamento TIMESTAMP,
    momento_autorizacao TIMESTAMP,
    momento_captura TIMESTAMP
);
```

### Fluxos de Pagamento

#### 1. Cartão de Crédito

```mermaid
sequenceDiagram
    participant Client
    participant PaymentHandler
    participant PaymentService
    participant MercadoPago
    participant Database
    participant Redis

    Client->>PaymentHandler: POST /payments/credit-card
    PaymentHandler->>PaymentService: CreateCreditCardPayment(request)
    PaymentService->>PaymentService: Validar dados de entrada
    PaymentService->>MercadoPago: CreateCreditCardPayment(mpRequest)
    MercadoPago-->>PaymentService: {status, payment_id, card_info}
    PaymentService->>Database: Salvar pagamento
    Database-->>PaymentService: Pagamento salvo
    PaymentService->>PaymentService: Atualizar timestamps
    PaymentService-->>PaymentHandler: PaymentResponse
    PaymentHandler-->>Client: 200 OK {pagamento, message}
    
    Note over Client,Redis: Em caso de erro
    MercadoPago-->>PaymentService: Erro de pagamento
    PaymentService-->>PaymentHandler: Erro
    PaymentHandler-->>Client: 400/500 {erro}
```

#### 2. Cartão de Débito

```mermaid
sequenceDiagram
    participant Client
    participant PaymentHandler
    participant PaymentService
    participant MercadoPago
    participant Database

    Client->>PaymentHandler: POST /payments/debit-card
    PaymentHandler->>PaymentService: CreateDebitCardPayment(request)
    PaymentService->>PaymentService: Validar dados (installments=1, capture=true)
    PaymentService->>MercadoPago: CreateCreditCardPayment(mpRequest)
    MercadoPago-->>PaymentService: {status: approved, payment_id}
    PaymentService->>Database: Salvar pagamento
    Database-->>PaymentService: Pagamento salvo
    PaymentService-->>PaymentHandler: PaymentResponse
    PaymentHandler-->>Client: 200 OK {pagamento aprovado}
```

#### 3. PIX

```mermaid
sequenceDiagram
    participant Client
    participant PaymentHandler
    participant PaymentService
    participant MercadoPago
    participant Database

    Client->>PaymentHandler: POST /payments/pix
    PaymentHandler->>PaymentService: CreatePIXPayment(request)
    PaymentService->>PaymentService: Validar dados PIX
    PaymentService->>MercadoPago: CreatePIXPayment(mpRequest)
    MercadoPago-->>PaymentService: {status: pending, qr_code, ticket_url}
    PaymentService->>Database: Salvar pagamento pendente
    Database-->>PaymentService: Pagamento salvo
    PaymentService-->>PaymentHandler: PIXResponse
    PaymentHandler-->>Client: 200 OK {qr_code, ticket_url}
    
    Note over Client,Database: Cliente paga via PIX
    MercadoPago->>PaymentService: Webhook: status=approved
    PaymentService->>Database: Atualizar status
    Database-->>PaymentService: Status atualizado
```

#### 4. Gestão de Cartões

```mermaid
sequenceDiagram
    participant Client
    participant CardHandler
    participant CardService
    participant MercadoPago
    participant Database

    Client->>CardHandler: POST /clients/:id/cards
    CardHandler->>CardService: Create(customerID, request)
    CardService->>MercadoPago: CreateCustomerCard(customerID, cardData)
    MercadoPago-->>CardService: {card_id, card_info}
    CardService-->>CardHandler: CardResponse
    CardHandler-->>Client: 200 OK {cartão criado}
    
    Note over Client,Database: Listar cartões
    Client->>CardHandler: GET /clients/:id/cards
    CardHandler->>CardService: List(customerID)
    CardService->>MercadoPago: ListCustomerCards(customerID)
    MercadoPago-->>CardService: [cartões]
    CardService-->>CardHandler: ListResponse
    CardHandler-->>Client: 200 OK {lista de cartões}
```

### Exemplos de Uso

#### Pagamento com Cartão de Crédito

```bash
curl -X POST http://localhost:1450/jampa-trip/api/v1/payments/credit-card \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "cliente_id": 123,
    "empresa_id": 456,
    "valor": 100.00,
    "descricao": "Pagamento de tour",
    "numero_parcelas": 3,
    "token_cartao": "card_token_from_frontend"
  }'
```

#### Pagamento via PIX

```bash
curl -X POST http://localhost:1450/jampa-trip/api/v1/payments/pix \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "cliente_id": 123,
    "empresa_id": 456,
    "valor": 100.00,
    "descricao": "Pagamento de tour via PIX"
  }'
```

#### Consultar Pagamentos

```bash
curl -X GET http://localhost:1450/jampa-trip/api/v1/payments \
  -H "Authorization: Bearer <access_token>"
```

### Monitoramento

O sistema inclui **logs estruturados** para monitoramento de pagamentos:

- **Eventos de pagamento** registrados com detalhes completos
- **Integração com Mercado Pago** monitorada
- **Debug e auditoria** facilitados
- **Status em tempo real** dos pagamentos

## 📦 Módulo de Orders

O projeto implementa o **módulo completo de Orders do Mercado Pago**, oferecendo gerenciamento avançado de pedidos, transações e fluxos de pagamento.

### Visão Geral

O módulo Orders permite:
- Criar pedidos estruturados com múltiplos itens
- Gerenciar transações associadas a um pedido
- Capturar, cancelar e reembolsar pedidos
- Processar pagamentos online vinculados a pedidos
- Rastrear o ciclo de vida completo de pedidos e transações

### Status de Orders

| Status | Descrição |
|--------|-----------|
| `pending` | Ordem pendente de pagamento |
| `processing` | Ordem em processamento |
| `paid` | Ordem paga com sucesso |
| `authorized` | Ordem autorizada (aguardando captura) |
| `cancelled` | Ordem cancelada |
| `refunded` | Ordem reembolsada |
| `expired` | Ordem expirada |

### Tipos de Transactions

| Tipo | Descrição |
|------|-----------|
| `payment` | Transação de pagamento |
| `refund` | Transação de reembolso |

### Fluxos de Operação

#### 1. Criar Ordem

```mermaid
sequenceDiagram
    participant Client
    participant OrderHandler
    participant OrderService
    participant MercadoPago
    participant Database

    Client->>OrderHandler: POST /orders
    OrderHandler->>OrderService: Create(request)
    OrderService->>OrderService: Validar dados
    OrderService->>MercadoPago: CreateOrder(items, payer, amount)
    MercadoPago-->>OrderService: {order_id, status}
    OrderService->>Database: Salvar ordem localmente
    Database-->>OrderService: Ordem criada
    OrderService-->>OrderHandler: OrderResponse
    OrderHandler-->>Client: 201 {order, message}
```

#### 2. Processar Pagamento de Ordem

```mermaid
sequenceDiagram
    participant Client
    participant OrderHandler
    participant OrderService
    participant MercadoPago
    participant Database

    Client->>OrderHandler: POST /orders/:id/process
    OrderHandler->>OrderService: ProcessOrder(orderID, paymentData)
    OrderService->>Database: Buscar ordem
    Database-->>OrderService: Ordem encontrada
    OrderService->>MercadoPago: ProcessOrder(card, installments)
    MercadoPago-->>OrderService: {payment_id, status: approved}
    OrderService->>Database: Criar transaction (type: payment)
    OrderService->>Database: Atualizar order status = paid
    Database-->>OrderService: Atualizado
    OrderService-->>OrderHandler: ProcessOrderResponse
    OrderHandler-->>Client: 200 {order, payment_details}
```

#### 3. Capturar Ordem Autorizada

```mermaid
sequenceDiagram
    participant Client
    participant OrderHandler
    participant OrderService
    participant MercadoPago
    participant Database

    Client->>OrderHandler: POST /orders/:id/capture
    OrderHandler->>OrderService: Capture(orderID)
    OrderService->>Database: Buscar ordem
    Database-->>OrderService: Status = authorized
    OrderService->>MercadoPago: CaptureOrder(orderID)
    MercadoPago-->>OrderService: {status: paid}
    OrderService->>Database: Atualizar status
    Database-->>OrderService: Status atualizado
    OrderService-->>OrderHandler: CaptureOrderResponse
    OrderHandler-->>Client: 200 {order captured}
```

#### 4. Reembolsar Ordem

```mermaid
sequenceDiagram
    participant Client
    participant OrderHandler
    participant OrderService
    participant MercadoPago
    participant Database

    Client->>OrderHandler: POST /orders/:id/refund
    OrderHandler->>OrderService: Refund(orderID, amount?)
    OrderService->>Database: Buscar ordem (status = paid)
    Database-->>OrderService: Ordem encontrada
    OrderService->>MercadoPago: RefundOrder(orderID)
    MercadoPago-->>OrderService: {status: refunded}
    OrderService->>Database: Criar transaction (type: refund)
    OrderService->>Database: Atualizar order status
    Database-->>OrderService: Atualizado
    OrderService-->>OrderHandler: RefundOrderResponse
    OrderHandler-->>Client: 200 {order, refund_details}
```

### Sincronização com Mercado Pago

O sistema sincroniza automaticamente o status das ordens:
- Ao buscar uma ordem (GET), o status é atualizado do Mercado Pago
- Transações são criadas localmente após operações no MP
- Status Detail é sempre mantido atualizado

## 📋 Fluxo de Reservas

O sistema de reservas conecta clientes, tours e pagamentos, permitindo o agendamento de passeios turísticos com controle completo do ciclo de vida.

### Funcionalidades Implementadas

- ✅ **Criação de Reservas** - Agendamento de passeios com validações
- ✅ **Gestão de Status** - Controle de estados (pendente, confirmada, cancelada)
- ✅ **Validações de Data** - Verificação de datas futuras e disponibilidade
- ✅ **Histórico de Reservas** - Consulta de reservas passadas e futuras
- ✅ **Cancelamento** - Cancelamento com regras de negócio
- ✅ **Integração com Pagamentos** - Vinculação automática com pagamentos

### Status de Reserva Suportados

| Status | Descrição |
|--------|-----------|
| `pending` | Pendente de confirmação |
| `confirmed` | Confirmada |
| `cancelled` | Cancelada |
| `completed` | Concluída |

### Fluxos de Reserva

#### 1. Fluxo Completo de Reserva

```mermaid
sequenceDiagram
    participant Client
    participant ReservaHandler
    participant ReservaService
    participant TourService
    participant PaymentService
    participant Database

    Client->>ReservaHandler: POST /reservations
    ReservaHandler->>ReservaService: Create(request)
    ReservaService->>ReservaService: Validar data futura
    ReservaService->>ReservaService: Validar data reserva < data passeio
    ReservaService->>TourService: Verificar disponibilidade
    TourService-->>ReservaService: Tour disponível
    ReservaService->>Database: Criar reserva (status: pending)
    Database-->>ReservaService: Reserva criada
    ReservaService->>PaymentService: Processar pagamento
    PaymentService-->>ReservaService: Pagamento aprovado
    ReservaService->>Database: Atualizar status (confirmed)
    Database-->>ReservaService: Status atualizado
    ReservaService-->>ReservaHandler: ReservaResponse
    ReservaHandler-->>Client: 200 OK {reserva confirmada}
    
    Note over Client,Database: Em caso de erro
    ReservaService-->>ReservaHandler: Erro de validação
    ReservaHandler-->>Client: 400 {erro de validação}
```

#### 2. Fluxo de Cancelamento de Reserva

```mermaid
flowchart TD
    A[Cliente solicita cancelamento] --> B{Reserva existe?}
    B -->|Não| C[Erro: Reserva não encontrada]
    B -->|Sim| D{Status permite cancelamento?}
    D -->|Não| E[Erro: Não pode ser cancelada]
    D -->|Sim| F[Verificar regras de negócio]
    F --> G{Data do passeio > hoje?}
    G -->|Não| H[Erro: Passeio já ocorreu]
    G -->|Sim| I[Atualizar status para 'cancelled']
    I --> J[Processar reembolso se necessário]
    J --> K[Atualizar timestamps]
    K --> L[Reserva cancelada com sucesso]
    
    style C fill:#ffcccc
    style E fill:#ffcccc
    style H fill:#ffcccc
    style L fill:#ccffcc
```

### Exemplos de Uso

#### Criar Reserva

```bash
curl -X POST http://localhost:1450/jampa-trip/api/v1/reservations \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "cliente_id": 123,
    "empresa_id": 456,
    "tour_id": 789,
    "data_passeio": "2024-12-25T10:00:00Z",
    "quantidade_pessoas": 2,
    "valor_total": 200.00,
    "observacoes": "Dieta vegetariana"
  }'
```

#### Cancelar Reserva

```bash
curl -X PUT http://localhost:1450/jampa-trip/api/v1/reservations/123/cancel \
  -H "Authorization: Bearer <access_token>"
```

## 🛠️ Tecnologias Utilizadas

- **[Go 1.23.5](https://golang.org/)** - Linguagem de programação
- **[Echo v4](https://echo.labstack.com/)** - Framework web
- **[GORM](https://gorm.io/)** - ORM para Go
- **[PostgreSQL 15](https://postgresql.org/)** - Banco de dados
- **[Redis 7](https://redis.io/)** - Cache e armazenamento de tokens
- **[JWT](https://jwt.io/)** - Autenticação com tokens
- **[Docker](https://docker.com/)** - Containerização
- **[Swagger/OpenAPI 3.0.3](https://swagger.io/)** - Documentação da API
- **[Ozzo Validation](https://github.com/go-ozzo/ozzo-validation)** - Validação de dados
- **[Mercado Pago API](https://www.mercadopago.com.br/developers/)** - Processamento de pagamentos

## 📚 Documentação Adicional

- [Documentação da API](docs/)
- [Swagger UI](http://localhost:1450/swagger/index.html)

## 🔧 Troubleshooting

### Problemas Comuns

1. **Redis não conecta**
   ```bash
   # Verificar se o Redis está rodando
   docker exec -it jampa-trip-redis redis-cli ping
   ```

2. **Tokens não são validados**
   - Verificar se o JWT_SECRET está correto
   - Verificar se o token está no Redis
   - Verificar se o token não expirou

3. **Erro de parsing de duração**
   - Verificar formato das variáveis JWT_*_EXPIRATION
   - Usar formato Go duration (ex: "15m", "168h")

### Comandos de Debug

```bash
# Verificar conexão com Redis
docker exec -it jampa-trip-redis redis-cli ping

# Listar chaves de tokens
docker exec -it jampa-trip-redis redis-cli keys "*token*"

# Verificar TTL de uma chave
docker exec -it jampa-trip-redis redis-cli ttl access_token:123:client
```
