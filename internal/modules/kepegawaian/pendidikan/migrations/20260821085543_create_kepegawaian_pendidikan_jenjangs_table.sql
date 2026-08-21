-- Migration: Create kepegawaian_pendidikan_jenjangs table
-- Timestamp: 20260821085543

CREATE TABLE IF NOT EXISTS kepegawaian_pendidikan_jenjangs (
    id          BIGSERIAL    PRIMARY KEY,
    code        VARCHAR(100) NOT NULL UNIQUE,
    label       VARCHAR(255) NOT NULL UNIQUE,
    fhir_system TEXT,
    created_by  BIGINT,
    updated_by  BIGINT,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_kepegawaian_pendidikan_jenjangs_deleted_at ON kepegawaian_pendidikan_jenjangs(deleted_at);
