-- Migration: Create artikel_kategoris table
-- Timestamp: 20260715085659

CREATE TABLE IF NOT EXISTS artikel_kategoris (
    id          BIGSERIAL    PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    created_by  BIGINT,
    updated_by  BIGINT,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_artikel_kategoris_deleted_at ON artikel_kategoris(deleted_at);
