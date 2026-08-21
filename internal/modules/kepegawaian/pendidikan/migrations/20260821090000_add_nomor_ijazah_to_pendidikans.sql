-- Migration: Add optional unique diploma number to kepegawaian_pendidikans

ALTER TABLE kepegawaian_pendidikans
    ADD COLUMN IF NOT EXISTS nomor_ijazah VARCHAR(255);

CREATE UNIQUE INDEX IF NOT EXISTS idx_kepegawaian_pendidikans_nomor_ijazah
    ON kepegawaian_pendidikans(nomor_ijazah)
    WHERE nomor_ijazah IS NOT NULL;