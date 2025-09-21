-- =============================================================================
-- SCRIPT DE INICIALIZAÇÃO DO BANCO DE DADOS JAMPA TRIP
-- =============================================================================
-- Este arquivo é executado automaticamente quando o container PostgreSQL é 
-- iniciado pela primeira vez
-- =============================================================================

-- Conectar ao banco de dados
\c jampa_trip_db;

-- =============================================================================
-- TABELA FORNECEDORES
-- =============================================================================

CREATE TABLE IF NOT EXISTS fornecedores (
    id SERIAL PRIMARY KEY,
    nome VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    senha VARCHAR(255) NOT NULL,
    cnpj VARCHAR(255) NOT NULL UNIQUE,
    telefone VARCHAR(255) NOT NULL,
    endereco VARCHAR(255) NOT NULL,
    momento_cadastro TIMESTAMP NOT NULL DEFAULT NOW(),
    momento_atualizacao TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- =============================================================================
-- ÍNDICES PARA FORNECEDORES
-- =============================================================================

CREATE INDEX IF NOT EXISTS idx_fornecedores_email ON fornecedores(email);
CREATE INDEX IF NOT EXISTS idx_fornecedores_cnpj ON fornecedores(cnpj);
CREATE INDEX IF NOT EXISTS idx_fornecedores_momento_cadastro ON fornecedores(momento_cadastro);

-- =============================================================================
-- TABELA CLIENTES
-- =============================================================================

CREATE TABLE clientes (
    id SERIAL PRIMARY KEY,
    nome VARCHAR(100) NOT NULL,
    email VARCHAR(40) UNIQUE NOT NULL,
    senha VARCHAR(255) NOT NULL,
    cpf VARCHAR(14) UNIQUE NOT NULL,
    telefone VARCHAR(15) NOT NULL,
    data_nascimento DATE NOT NULL,
    momento_cadastro TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    momento_atualizacao TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- =============================================================================
-- ÍNDICES PARA CLIENTES
-- =============================================================================

CREATE INDEX IF NOT EXISTS idx_clientes_email ON clientes(email);
CREATE INDEX IF NOT EXISTS idx_clientes_cpf ON clientes(cpf);
CREATE INDEX IF NOT EXISTS idx_clientes_momento_cadastro ON clientes(momento_cadastro);
CREATE INDEX IF NOT EXISTS idx_clientes_momento_atualizacao ON clientes(momento_atualizacao);
