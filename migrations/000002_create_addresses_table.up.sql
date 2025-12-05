-- Create enum types
CREATE TYPE entity_type AS ENUM ('user');
CREATE TYPE address_type AS ENUM ('shipping', 'billing');

CREATE TABLE addresses (
    id SERIAL PRIMARY KEY,
    entity_type entity_type NOT NULL,
    entity_id INTEGER NOT NULL,
    address_type address_type NOT NULL,
    street_line1 VARCHAR(255) NOT NULL,
    street_line2 VARCHAR(255),
    city VARCHAR(100) NOT NULL,
    state VARCHAR(100) NOT NULL,
    postal_code VARCHAR(20) NOT NULL,
    country VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

-- Index for looking up addresses by entity
CREATE INDEX idx_addresses_entity ON addresses(entity_type, entity_id);

-- Index for looking up by address type
CREATE INDEX idx_addresses_type ON addresses(entity_type, entity_id, address_type);
