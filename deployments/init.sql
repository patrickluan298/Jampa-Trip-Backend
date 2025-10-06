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