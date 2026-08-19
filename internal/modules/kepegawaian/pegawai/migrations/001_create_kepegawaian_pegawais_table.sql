-- Migration: Create kepegawaian_pegawais table
-- Timestamp: 20260819114029

CREATE TABLE IF NOT EXISTS kepegawaian_pegawais (
    id          BIGSERIAL    PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    created_by  BIGINT,
    updated_by  BIGINT,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_kepegawaian_pegawais_deleted_at ON kepegawaian_pegawais(deleted_at);
