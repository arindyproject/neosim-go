-- Migration: Create artikel_ketegoris table
-- Timestamp: 20260619190109

CREATE TABLE IF NOT EXISTS artikel_ketegoris (
    id          BIGSERIAL    PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    created_by  BIGINT,
    updated_by  BIGINT,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_artikel_ketegoris_deleted_at ON artikel_ketegoris(deleted_at);
