-- =============================================================================
-- JAMPA TRIP DATABASE INITIALIZATION SCRIPT
-- =============================================================================
-- This file is automatically executed when the PostgreSQL container is 
-- started for the first time
-- =============================================================================

-- Connect to database
\c jampa_trip_db;

-- =============================================================================
-- COMPANIES TABLE
-- =============================================================================

CREATE TABLE IF NOT EXISTS companies (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    cnpj VARCHAR(255) NOT NULL UNIQUE,
    phone VARCHAR(255) NOT NULL,
    address VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- =============================================================================
-- INDEXES FOR COMPANIES
-- =============================================================================

CREATE INDEX IF NOT EXISTS idx_companies_email ON companies(email);
CREATE INDEX IF NOT EXISTS idx_companies_cnpj ON companies(cnpj);
CREATE INDEX IF NOT EXISTS idx_companies_created_at ON companies(created_at);

-- =============================================================================
-- CLIENTS TABLE
-- =============================================================================

CREATE TABLE clients (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(40) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    cpf VARCHAR(14) UNIQUE NOT NULL,
    phone VARCHAR(15) NOT NULL,
    birth_date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- =============================================================================
-- INDEXES FOR CLIENTS
-- =============================================================================

CREATE INDEX IF NOT EXISTS idx_clients_email ON clients(email);
CREATE INDEX IF NOT EXISTS idx_clients_cpf ON clients(cpf);
CREATE INDEX IF NOT EXISTS idx_clients_created_at ON clients(created_at);
CREATE INDEX IF NOT EXISTS idx_clients_updated_at ON clients(updated_at);

-- =============================================================================
-- TOURS TABLE
-- =============================================================================

CREATE TABLE IF NOT EXISTS tours (
    id SERIAL PRIMARY KEY,
    company_id INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    dates TIMESTAMP[],
    departure_time VARCHAR(10),
    arrival_time VARCHAR(10),
    max_people INTEGER DEFAULT 1,
    description TEXT,
    images TEXT[],
    price DECIMAL(10,2) DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE
);

-- =============================================================================
-- INDEXES FOR TOURS
-- =============================================================================

CREATE INDEX IF NOT EXISTS idx_tours_company_id ON tours(company_id);
CREATE INDEX IF NOT EXISTS idx_tours_created_at ON tours(created_at);
CREATE INDEX IF NOT EXISTS idx_tours_price ON tours(price);
CREATE INDEX IF NOT EXISTS idx_tours_name ON tours(name);

-- =============================================================================
-- FEEDBACKS TABLE
-- =============================================================================

CREATE TABLE IF NOT EXISTS feedbacks (
    id SERIAL PRIMARY KEY,
    cliente_id INTEGER NOT NULL,
    empresa_id INTEGER NOT NULL,
    reserva_id INTEGER,
    nota INTEGER NOT NULL CHECK (nota >= 1 AND nota <= 5),
    comentario TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'ativo',
    momento_criacao TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    momento_atualizacao TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (cliente_id) REFERENCES clients(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    FOREIGN KEY (empresa_id) REFERENCES companies(id) ON UPDATE CASCADE ON DELETE RESTRICT
);

-- =============================================================================
-- INDEXES FOR FEEDBACKS
-- =============================================================================

CREATE INDEX IF NOT EXISTS idx_feedbacks_cliente_id ON feedbacks(cliente_id);
CREATE INDEX IF NOT EXISTS idx_feedbacks_empresa_id ON feedbacks(empresa_id);
CREATE INDEX IF NOT EXISTS idx_feedbacks_reserva_id ON feedbacks(reserva_id);
CREATE INDEX IF NOT EXISTS idx_feedbacks_nota ON feedbacks(nota);
CREATE INDEX IF NOT EXISTS idx_feedbacks_status ON feedbacks(status);
CREATE INDEX IF NOT EXISTS idx_feedbacks_momento_criacao ON feedbacks(momento_criacao);

-- =============================================================================
-- COMMENTS FOR DOCUMENTATION
-- =============================================================================

COMMENT ON TABLE feedbacks IS 'Tabela para armazenar feedbacks e avaliações de clientes sobre empresas';
COMMENT ON COLUMN feedbacks.nota IS 'Nota de avaliação de 1 a 5 estrelas';
COMMENT ON COLUMN feedbacks.status IS 'Status do feedback: ativo, inativo ou moderado';