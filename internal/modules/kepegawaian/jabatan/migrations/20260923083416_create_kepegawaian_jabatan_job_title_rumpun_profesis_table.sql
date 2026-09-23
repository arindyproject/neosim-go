-- Migration: Create kepegawaian_jabatan_job_title_rumpun_profesis table
-- Timestamp: 20260923083416

CREATE TABLE IF NOT EXISTS kepegawaian_jabatan_job_title_rumpun_profesis (
    id          BIGSERIAL    PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    created_by  BIGINT,
    updated_by  BIGINT,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_kepegawaian_jabatan_job_title_rumpun_profesis_deleted_at ON kepegawaian_jabatan_job_title_rumpun_profesis(deleted_at);
