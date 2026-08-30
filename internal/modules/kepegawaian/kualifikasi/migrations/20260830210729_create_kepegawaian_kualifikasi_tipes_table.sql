-- Migration: Create kepegawaian_kualifikasi_tipes table
-- Timestamp: 20260830210729

CREATE TABLE IF NOT EXISTS kepegawaian_kualifikasi_tipes (
    id          BIGSERIAL    PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    created_by  BIGINT,
    updated_by  BIGINT,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_kepegawaian_kualifikasi_tipes_deleted_at ON kepegawaian_kualifikasi_tipes(deleted_at);
