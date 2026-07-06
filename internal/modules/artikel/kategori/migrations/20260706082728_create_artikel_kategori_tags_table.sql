-- Migration: Create artikel_kategori_tags table
-- Timestamp: 20260706082728

CREATE TABLE IF NOT EXISTS artikel_kategori_tags (
    id          BIGSERIAL    PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    created_by  BIGINT,
    updated_by  BIGINT,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_artikel_kategori_tags_deleted_at ON artikel_kategori_tags(deleted_at);
