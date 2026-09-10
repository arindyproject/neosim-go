-- Migration: Create kepegawaian_alamats table
-- Timestamp: 20260910111246

CREATE TABLE IF NOT EXISTS kepegawaian_alamats (
    id          BIGSERIAL    PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    created_by  BIGINT,
    updated_by  BIGINT,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_kepegawaian_alamats_deleted_at ON kepegawaian_alamats(deleted_at);
