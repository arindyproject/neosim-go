package services

import (
	"fmt"
)

// ──────────────────────────────────────────────────────────────────────────
func cacheKeyJenisSelectList(search string) string {
	return fmt.Sprintf("kepegawaian:pegawai:jenis:selectlist:search%s", search)
}

func cacheKeyStatusSelectList(search string) string {
	return fmt.Sprintf("kepegawaian:pegawai:status:selectlist:search%s", search)
}

// ─── Cache Prefix Constants ───────────────────────────────────────────────────────
// Digunakan untuk InvalidateList agar konsisten dan tidak typo
const (
	cachePrefixJenisSelectList  = "kepegawaian:pegawai:jenis:selectlist:"
	cachePrefixStatusSelectList = "kepegawaian:pegawai:status:selectlist:"
)
