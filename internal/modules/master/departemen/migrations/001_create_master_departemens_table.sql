-- Migration: Create master_departemens table
-- Timestamp: 20260819205404

CREATE TABLE IF NOT EXISTS master_departemens (
    id          BIGSERIAL    PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    created_by  BIGINT,
    updated_by  BIGINT,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_master_departemens_deleted_at ON master_departemens(deleted_at);
