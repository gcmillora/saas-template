-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE tenant_tbl (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE user_tbl (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255),
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    auth_provider VARCHAR(50) NOT NULL DEFAULT 'email',
    auth_provider_id VARCHAR(255),
    tenant_id UUID NOT NULL REFERENCES tenant_tbl(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_user_tbl_email ON user_tbl(email);
CREATE INDEX idx_user_tbl_tenant_id ON user_tbl(tenant_id);

-- +goose Down
DROP TABLE IF EXISTS user_tbl;
DROP TABLE IF EXISTS tenant_tbl;
DROP EXTENSION IF EXISTS "uuid-ossp";
