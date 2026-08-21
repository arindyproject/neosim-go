-- Migration: Create kepegawaian_pendidikans table
-- Timestamp: 20260821085532

CREATE TABLE IF NOT EXISTS kepegawaian_pendidikans (
    id          BIGSERIAL    PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    created_by  BIGINT,
    updated_by  BIGINT,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_kepegawaian_pendidikans_deleted_at ON kepegawaian_pendidikans(deleted_at);
