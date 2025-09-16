-- Script de inicialização do banco de dados Jampa Trip
-- Este arquivo é executado automaticamente quando o container PostgreSQL é iniciado pela primeira vez

-- Conectar ao banco de dados
\c jampa_trip_db;

-- FORNECEDORES
CREATE TABLE IF NOT EXISTS fornecedores (
    id SERIAL PRIMARY KEY,
    nome VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    senha VARCHAR(255) NOT NULL,
    cnpj VARCHAR(255) NOT NULL UNIQUE,
    telefone VARCHAR(255) NOT NULL,
    endereco VARCHAR(255) NOT NULL,
    momento_cadastro TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_fornecedores_email ON fornecedores(email);
CREATE INDEX IF NOT EXISTS idx_fornecedores_cnpj ON fornecedores(cnpj);
CREATE INDEX IF NOT EXISTS idx_fornecedores_momento_cadastro ON fornecedores(momento_cadastro);

COMMENT ON TABLE fornecedores IS 'Tabela para armazenar informações dos fornecedores do sistema Jampa Trip';
COMMENT ON COLUMN fornecedores.id IS 'Identificador único do fornecedor';
COMMENT ON COLUMN fornecedores.nome IS 'Nome da empresa fornecedora';
COMMENT ON COLUMN fornecedores.email IS 'Email de contato do fornecedor (único)';
COMMENT ON COLUMN fornecedores.senha IS 'Senha criptografada do fornecedor';
COMMENT ON COLUMN fornecedores.cnpj IS 'CNPJ da empresa (único)';
COMMENT ON COLUMN fornecedores.telefone IS 'Telefone de contato';
COMMENT ON COLUMN fornecedores.endereco IS 'Endereço completo da empresa';
COMMENT ON COLUMN fornecedores.momento_cadastro IS 'Data e hora do cadastro do fornecedor';
