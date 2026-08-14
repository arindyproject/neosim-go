-- Migration: Create kepegawaian_identifiers table
-- Timestamp: 20260814105026

CREATE TABLE IF NOT EXISTS kepegawaian_identifiers (
    id          BIGSERIAL    PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    created_by  BIGINT,
    updated_by  BIGINT,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_kepegawaian_identifiers_deleted_at ON kepegawaian_identifiers(deleted_at);
