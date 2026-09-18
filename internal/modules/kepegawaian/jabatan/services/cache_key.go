package services

import (
	"fmt"
)

// ──────────────────────────────────────────────────────────────────────────
func cacheKeyPositionsSelectList(search string) string {
	return fmt.Sprintf("kepegawaian:jabatan:positions:selectlist:search%s", search)
}

func cacheKeySpecializationsSelectList(search string) string {
	return fmt.Sprintf("kepegawaian:jabatan:specializations:selectlist:search%s", search)
}

func cacheKeyJobTitlesSelectList(search string) string {
	return fmt.Sprintf("kepegawaian:jabatan:job_titles:selectlist:search%s", search)
}

// ─── Cache Prefix Constants ───────────────────────────────────────────────────────
// Digunakan untuk InvalidateList agar konsisten dan tidak typo
const (
	cachePrefixPositionsSelectList       = "kepegawaian:jabatan:positions:selectlist:"
	cachePrefixSpecializationsSelectList = "kepegawaian:jabatan:specializations:selectlist:"
	cachePrefixJobTitlesSelectList       = "kepegawaian:jabatan:job_titles:selectlist:"
)
