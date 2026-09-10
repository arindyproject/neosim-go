-- Migration: Create kepegawaian_alamat_tipes table
-- Timestamp: 20260910111246

CREATE TABLE IF NOT EXISTS kepegawaian_alamat_tipes (
    id          BIGSERIAL    PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    created_by  BIGINT,
    updated_by  BIGINT,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_kepegawaian_alamat_tipes_deleted_at ON kepegawaian_alamat_tipes(deleted_at);
