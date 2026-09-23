package services

import (
	"fmt"
)

// ──────────────────────────────────────────────────────────────────────────
func cacheKeyPositionsSelectList(search string) string {
	return fmt.Sprintf("kepegawaian:jabatan:positions:selectlist:search%s", search)
}
func cacheKeyPositionsKategoriSelectList(search string) string {
	return fmt.Sprintf("kepegawaian:jabatan:positions:kategori:selectlist:search%s", search)
}

// ---------------------------------------------------------------------------
func cacheKeySpecializationsSelectList(search string) string {
	return fmt.Sprintf("kepegawaian:jabatan:specializations:selectlist:search%s", search)
}

// ---------------------------------------------------------------------------
func cacheKeyJobTitlesSelectList(search string) string {
	return fmt.Sprintf("kepegawaian:jabatan:job_titles:selectlist:search%s", search)
}
func cacheKeyJobTitleKategoriSelectList(search string) string {
	return fmt.Sprintf("kepegawaian:jabatan:job_titles:kategori:selectlist:search%s", search)
}
func cacheKeyJobTitleRumpunProfesiSelectList(search string) string {
	return fmt.Sprintf("kepegawaian:jabatan:job_titles:rumpun_profesi:selectlist:search%s", search)
}

// ─── Cache Prefix Constants ───────────────────────────────────────────────────────
// Digunakan untuk InvalidateList agar konsisten dan tidak typo
const (
	cachePrefixPositionsSelectList         = "kepegawaian:jabatan:positions:selectlist:"
	cachePrefixPositionsKategoriSelectList = "kepegawaian:jabatan:positions:kategori:selectlist:"

	cachePrefixSpecializationsSelectList = "kepegawaian:jabatan:specializations:selectlist:"

	cachePrefixJobTitlesSelectList             = "kepegawaian:jabatan:job_titles:selectlist:"
	cachePrefixJobTitleKategoriSelectList      = "kepegawaian:jabatan:job_titles:kategori:selectlist:"
	cachePrefixJobTitleRumpunProfesiSelectList = "kepegawaian:jabatan:job_titles:rumpun_profesi:selectlist:"
)
