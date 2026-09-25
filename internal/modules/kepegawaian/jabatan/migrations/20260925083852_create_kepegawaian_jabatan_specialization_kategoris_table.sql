-- Migration: Create kepegawaian_jabatan_specialization_kategoris table
-- Timestamp: 20260925083852

CREATE TABLE IF NOT EXISTS kepegawaian_jabatan_specialization_kategoris (
    id          BIGSERIAL    PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    created_by  BIGINT,
    updated_by  BIGINT,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_kepegawaian_jabatan_specialization_kategoris_deleted_at ON kepegawaian_jabatan_specialization_kategoris(deleted_at);
