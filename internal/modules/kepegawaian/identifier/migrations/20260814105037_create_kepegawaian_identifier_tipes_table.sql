-- Migration: Create kepegawaian_identifier_tipes table
-- Timestamp: 20260814105037

CREATE TABLE IF NOT EXISTS kepegawaian_identifier_tipes (
    id          BIGSERIAL    PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    created_by  BIGINT,
    updated_by  BIGINT,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_kepegawaian_identifier_tipes_deleted_at ON kepegawaian_identifier_tipes(deleted_at);
