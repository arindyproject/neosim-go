-- Migration: Create kepegawaian_identifikasis table
-- Timestamp: 20260623112951

CREATE TABLE IF NOT EXISTS kepegawaian_identifikasis (
    id          BIGSERIAL    PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    created_by  BIGINT,
    updated_by  BIGINT,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_kepegawaian_identifikasis_deleted_at ON kepegawaian_identifikasis(deleted_at);
