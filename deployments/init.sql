-- =================================================================================================
-- This file is automatically executed when the PostgreSQL container is started for the first time.
-- =================================================================================================

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
-- IMAGES TABLE
-- =============================================================================

CREATE TABLE IF NOT EXISTS images (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    tour_id INTEGER REFERENCES tours(id) ON DELETE SET NULL,
    filename VARCHAR(255) NOT NULL,
    original_name VARCHAR(255),
    url VARCHAR(500) NOT NULL,
    thumbnail_url VARCHAR(500),
    size INTEGER NOT NULL,
    width INTEGER,
    height INTEGER,
    format VARCHAR(10) NOT NULL,
    description TEXT,
    alt_text VARCHAR(255),
    is_primary BOOLEAN DEFAULT FALSE,
    sort_order INTEGER DEFAULT 0,
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- =============================================================================
-- INDEXES FOR IMAGES
-- =============================================================================

CREATE INDEX IF NOT EXISTS idx_images_user_id ON images(user_id);
CREATE INDEX IF NOT EXISTS idx_images_tour_id ON images(tour_id);
CREATE INDEX IF NOT EXISTS idx_images_format ON images(format);
CREATE INDEX IF NOT EXISTS idx_images_uploaded_at ON images(uploaded_at);
CREATE INDEX IF NOT EXISTS idx_images_is_primary ON images(is_primary);
CREATE INDEX IF NOT EXISTS idx_images_sort_order ON images(sort_order);

-- =============================================================================
-- COMPOSITE INDEXES FOR IMAGES
-- =============================================================================

CREATE INDEX IF NOT EXISTS idx_images_user_tour ON images(user_id, tour_id);
CREATE INDEX IF NOT EXISTS idx_images_user_uploaded ON images(user_id, uploaded_at);
CREATE INDEX IF NOT EXISTS idx_images_tour_sort ON images(tour_id, sort_order);

-- =============================================================================
-- CONSTRAINTS FOR IMAGES
-- =============================================================================

ALTER TABLE images ADD CONSTRAINT chk_images_size CHECK (size > 0);
ALTER TABLE images ADD CONSTRAINT chk_images_width CHECK (width > 0);
ALTER TABLE images ADD CONSTRAINT chk_images_height CHECK (height > 0);
ALTER TABLE images ADD CONSTRAINT chk_images_format CHECK (format IN ('jpg', 'jpeg', 'png', 'gif', 'webp'));
ALTER TABLE images ADD CONSTRAINT chk_images_sort_order CHECK (sort_order >= 0);

-- =============================================================================
-- FUNCTIONS FOR IMAGES
-- =============================================================================

CREATE OR REPLACE FUNCTION update_images_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION ensure_single_primary_image()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.is_primary = TRUE AND NEW.tour_id IS NOT NULL THEN
        UPDATE images 
        SET is_primary = FALSE, updated_at = CURRENT_TIMESTAMP
        WHERE tour_id = NEW.tour_id 
          AND user_id = NEW.user_id 
          AND id != NEW.id;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- =============================================================================
-- TRIGGERS FOR IMAGES
-- =============================================================================

CREATE TRIGGER trigger_update_images_updated_at
    BEFORE UPDATE ON images
    FOR EACH ROW
    EXECUTE FUNCTION update_images_updated_at();

CREATE TRIGGER trigger_ensure_single_primary_image
    BEFORE INSERT OR UPDATE ON images
    FOR EACH ROW
    EXECUTE FUNCTION ensure_single_primary_image();

-- =============================================================================
-- VIEWS FOR IMAGES
-- =============================================================================

CREATE OR REPLACE VIEW image_stats AS
SELECT 
    user_id,
    COUNT(*) as total_images,
    SUM(size) as total_size,
    AVG(size) as average_size,
    COUNT(CASE WHEN is_primary = TRUE THEN 1 END) as primary_images,
    COUNT(CASE WHEN tour_id IS NULL THEN 1 END) as unused_images,
    COUNT(CASE WHEN uploaded_at >= NOW() - INTERVAL '24 hours' THEN 1 END) as recent_uploads,
    COUNT(CASE WHEN format = 'jpg' OR format = 'jpeg' THEN 1 END) as jpg_count,
    COUNT(CASE WHEN format = 'png' THEN 1 END) as png_count,
    COUNT(CASE WHEN format = 'gif' THEN 1 END) as gif_count,
    COUNT(CASE WHEN format = 'webp' THEN 1 END) as webp_count
FROM images
GROUP BY user_id;

CREATE OR REPLACE VIEW tour_images_with_metadata AS
SELECT 
    i.*,
    t.name as tour_name,
    t.company_id,
    c.name as company_name
FROM images i
LEFT JOIN tours t ON i.tour_id = t.id
LEFT JOIN companies c ON t.company_id = c.id
WHERE i.tour_id IS NOT NULL;

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

-- =============================================================================
-- RESERVATIONS TABLE
-- =============================================================================

CREATE TABLE IF NOT EXISTS reservations (
    id SERIAL PRIMARY KEY,
    tour_id INTEGER NOT NULL,
    cliente_id INTEGER NOT NULL,
    pagamento_id INTEGER,
    status VARCHAR(50) NOT NULL DEFAULT 'pendente',
    data_reserva TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    data_passeio_selecionada DATE NOT NULL,
    quantidade_pessoas INTEGER NOT NULL DEFAULT 1,
    valor_total DECIMAL(10,2) NOT NULL,
    observacoes TEXT,
    momento_criacao TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    momento_atualizacao TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    momento_cancelamento TIMESTAMP,
    
    FOREIGN KEY (tour_id) REFERENCES tours(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    FOREIGN KEY (cliente_id) REFERENCES clients(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    FOREIGN KEY (pagamento_id) REFERENCES pagamentos(id) ON UPDATE CASCADE ON DELETE SET NULL,
    
    CONSTRAINT chk_reservations_quantidade_pessoas CHECK (quantidade_pessoas >= 1),
    CONSTRAINT chk_reservations_valor_total CHECK (valor_total >= 0.01),
    CONSTRAINT chk_reservations_status CHECK (status IN ('pendente', 'aguardando_pagamento', 'confirmada', 'cancelada', 'concluida'))
);

-- =============================================================================
-- INDEXES FOR RESERVATIONS
-- =============================================================================

CREATE INDEX IF NOT EXISTS idx_reservations_tour_id ON reservations(tour_id);
CREATE INDEX IF NOT EXISTS idx_reservations_cliente_id ON reservations(cliente_id);
CREATE INDEX IF NOT EXISTS idx_reservations_pagamento_id ON reservations(pagamento_id);
CREATE INDEX IF NOT EXISTS idx_reservations_status ON reservations(status);
CREATE INDEX IF NOT EXISTS idx_reservations_data_passeio ON reservations(data_passeio_selecionada);
CREATE INDEX IF NOT EXISTS idx_reservations_momento_criacao ON reservations(momento_criacao);
CREATE INDEX IF NOT EXISTS idx_reservations_tour_data ON reservations(tour_id, data_passeio_selecionada);

-- =============================================================================
-- FUNCTIONS FOR RESERVATIONS
-- =============================================================================

CREATE OR REPLACE FUNCTION update_reservations_momento_atualizacao()
RETURNS TRIGGER AS $$
BEGIN
    NEW.momento_atualizacao = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION validate_tour_availability()
RETURNS TRIGGER AS $$
DECLARE
    v_max_people INTEGER;
    v_reserved_count INTEGER;
BEGIN
    SELECT max_people INTO v_max_people
    FROM tours
    WHERE id = NEW.tour_id;
    
    SELECT COALESCE(SUM(quantidade_pessoas), 0) INTO v_reserved_count
    FROM reservations
    WHERE tour_id = NEW.tour_id
        AND data_passeio_selecionada = NEW.data_passeio_selecionada
        AND status NOT IN ('cancelada')
        AND id != COALESCE(NEW.id, 0);
    
    IF (v_reserved_count + NEW.quantidade_pessoas) > v_max_people THEN
        RAISE EXCEPTION 'Não há vagas disponíveis para este tour na data selecionada. Vagas disponíveis: %, Solicitadas: %',
            (v_max_people - v_reserved_count), NEW.quantidade_pessoas;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- =============================================================================
-- TRIGGERS FOR RESERVATIONS
-- =============================================================================

CREATE TRIGGER trigger_update_reservations_momento_atualizacao
    BEFORE UPDATE ON reservations
    FOR EACH ROW
    EXECUTE FUNCTION update_reservations_momento_atualizacao();

CREATE TRIGGER trigger_validate_tour_availability
    BEFORE INSERT OR UPDATE ON reservations
    FOR EACH ROW
    WHEN (NEW.status NOT IN ('cancelada'))
    EXECUTE FUNCTION validate_tour_availability();

-- =============================================================================
-- VIEWS FOR RESERVATIONS
-- =============================================================================

CREATE OR REPLACE VIEW reservation_details AS
SELECT 
    r.id,
    r.tour_id,
    r.cliente_id,
    r.pagamento_id,
    r.status,
    r.data_reserva,
    r.data_passeio_selecionada,
    r.quantidade_pessoas,
    r.valor_total,
    r.observacoes,
    r.momento_criacao,
    r.momento_atualizacao,
    r.momento_cancelamento,
    t.name as tour_name,
    t.company_id as empresa_id,
    c.name as cliente_name,
    c.email as cliente_email,
    comp.name as empresa_name,
    p.status as pagamento_status
FROM reservations r
INNER JOIN tours t ON r.tour_id = t.id
INNER JOIN clients c ON r.cliente_id = c.id
INNER JOIN companies comp ON t.company_id = comp.id
LEFT JOIN pagamentos p ON r.pagamento_id = p.id;

CREATE OR REPLACE VIEW tour_availability_stats AS
SELECT 
    t.id as tour_id,
    t.name as tour_name,
    t.max_people,
    r.data_passeio_selecionada,
    COALESCE(SUM(CASE WHEN r.status NOT IN ('cancelada') THEN r.quantidade_pessoas ELSE 0 END), 0) as reserved_spots,
    t.max_people - COALESCE(SUM(CASE WHEN r.status NOT IN ('cancelada') THEN r.quantidade_pessoas ELSE 0 END), 0) as available_spots,
    COUNT(DISTINCT r.id) FILTER (WHERE r.status NOT IN ('cancelada')) as total_reservations
FROM tours t
LEFT JOIN reservations r ON t.id = r.tour_id
GROUP BY t.id, t.name, t.max_people, r.data_passeio_selecionada;

-- =============================================================================
-- COMMENTS FOR DOCUMENTATION
-- =============================================================================

COMMENT ON TABLE reservations IS 'Tabela para armazenar reservas de tours feitas por clientes';
COMMENT ON COLUMN reservations.tour_id IS 'ID do tour reservado';
COMMENT ON COLUMN reservations.cliente_id IS 'ID do cliente que fez a reserva';
COMMENT ON COLUMN reservations.pagamento_id IS 'ID do pagamento associado à reserva';
COMMENT ON COLUMN reservations.status IS 'Status da reserva: pendente, aguardando_pagamento, confirmada, cancelada, concluida';
COMMENT ON COLUMN reservations.data_passeio_selecionada IS 'Data específica do tour escolhida pelo cliente';
COMMENT ON COLUMN reservations.quantidade_pessoas IS 'Número de pessoas incluídas na reserva';
COMMENT ON COLUMN reservations.valor_total IS 'Valor total da reserva em BRL';

-- =============================================================================
-- ADD FOREIGN KEY FROM FEEDBACKS TO RESERVATIONS
-- =============================================================================

ALTER TABLE feedbacks
ADD CONSTRAINT fk_feedbacks_reserva_id 
FOREIGN KEY (reserva_id) REFERENCES reservations(id) ON UPDATE CASCADE ON DELETE SET NULL;