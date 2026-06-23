-- Migration: Create kepegawaian_kontaks table
-- Timestamp: 20260623113008

CREATE TABLE IF NOT EXISTS kepegawaian_kontaks (
    id          BIGSERIAL    PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    created_by  BIGINT,
    updated_by  BIGINT,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_kepegawaian_kontaks_deleted_at ON kepegawaian_kontaks(deleted_at);
