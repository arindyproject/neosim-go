-- Migration: Create kepegawaian_lampirans table
-- Timestamp: 20260623113021

CREATE TABLE IF NOT EXISTS kepegawaian_lampirans (
    id          BIGSERIAL    PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    created_by  BIGINT,
    updated_by  BIGINT,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_kepegawaian_lampirans_deleted_at ON kepegawaian_lampirans(deleted_at);
