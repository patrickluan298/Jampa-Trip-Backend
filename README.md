# Jampa Trip Backend

[![Go Version](https://img.shields.io/badge/Go-1.23.5-blue.svg)](https://golang.org/)
[![Echo Framework](https://img.shields.io/badge/Echo-v4.13.4-green.svg)](https://echo.labstack.com/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-blue.svg)](https://postgresql.org/)
[![Docker](https://img.shields.io/badge/Docker-Supported-blue.svg)](https://docker.com/)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

##  Descrição

O **Jampa Trip Backend** é uma API REST desenvolvida em Go que serve como backend para uma aplicação mobile de turismo. O projeto foi desenvolvido como parte de um TCC (Trabalho de Conclusão de Curso) do curso de Ciência da Computação, focando na gestão de fornecedores de serviços turísticos e clientes.

A aplicação oferece funcionalidades de autenticação, cadastro de fornecedores e clientes, com uma arquitetura limpa e escalável utilizando o framework Echo, GORM para ORM e PostgreSQL como banco de dados.

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
# Copie o arquivo de exemplo (se existir) ou configure manualmente
export DEBUG=false
export HTTP_SERVER_READ_TIMEOUT=20
export HTTP_SERVER_WRITE_TIMEOUT=60
export HTTP_SERVER_IDLE_TIMEOUT=120
export HTTP_SERVER_PORT=:1450

export DATABASE_POSTGRES_HOST=localhost
export DATABASE_POSTGRES_PORT=5432
export DATABASE_POSTGRES_NAME=jampa_trip_db
export DATABASE_POSTGRES_USER=jampa_trip_user
export DATABASE_POSTGRES_PASSWORD=jampa_trip_password
export DATABASE_POSTGRES_POOL_MAX_LIFETIME_CONNECTION=300
export DATABASE_POSTGRES_LOG=""
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
├── build                 # Arquivos de build
├── cmd                   # Ponto de entrada da aplicação
│   └── app
├── deployments           # Configurações de deploy (Docker Compose, scripts SQL)
├── docs                  # Documentação da API (OpenAPI/Swagger)
│   ├── components
│   └── paths
│       └── fornecedor
├── internal              # Código interno da aplicação
│   ├── app               # Lógica de negócio
│   │   ├── contract
│   │   ├── handler
│   │   ├── middleware
│   │   ├── model
│   │   ├── query
│   │   ├── repository
│   │   └── service
│   └── pkg               # Pacotes utilitários
│       ├── config
│       ├── database
│       ├── util
│       └── webserver
├── tests                 # Testes automatizados
├── go.mod
├── go.sum
├── Makefile              # Comandos de automação
└── run.sh                # Script de execução
```

### Arquitetura

O projeto segue os princípios da **Clean Architecture** com as seguintes camadas:

- **Handler:** Recebe requisições HTTP e valida dados de entrada
- **Service:** Contém a lógica de negócio
- **Repository:** Gerencia acesso aos dados
- **Model:** Define as entidades do domínio

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

### Configuração do Banco de Dados

O banco PostgreSQL é configurado automaticamente via Docker Compose com:

- **Database:** `jampa_trip_db`
- **User:** `jampa_trip_user`
- **Password:** `jampa_trip_password`
- **Port:** `6432` (mapeada para `5432` no container)

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

## ️ Tecnologias Utilizadas

- **[Go 1.23.5](https://golang.org/)** - Linguagem de programação
- **[Echo v4](https://echo.labstack.com/)** - Framework web
- **[GORM](https://gorm.io/)** - ORM para Go
- **[PostgreSQL 15](https://postgresql.org/)** - Banco de dados
- **[Docker](https://docker.com/)** - Containerização
- **[Swagger/OpenAPI 3.0.3](https://swagger.io/)** - Documentação da API
- **[Ozzo Validation](https://github.com/go-ozzo/ozzo-validation)** - Validação de dados

---

**Desenvolvido como parte do TCC do curso de Ciência da Computação** 🎓
