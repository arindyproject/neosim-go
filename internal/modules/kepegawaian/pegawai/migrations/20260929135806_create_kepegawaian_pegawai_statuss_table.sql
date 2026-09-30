-- Migration: Create kepegawaian_pegawai_statuss table
-- Timestamp: 20260929135806

CREATE TABLE IF NOT EXISTS kepegawaian_pegawai_statuss (
    id          BIGSERIAL    PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    created_by  BIGINT,
    updated_by  BIGINT,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_kepegawaian_pegawai_statuss_deleted_at ON kepegawaian_pegawai_statuss(deleted_at);
