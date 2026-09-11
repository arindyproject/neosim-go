package services

import (
	"fmt"
)

// ──────────────────────────────────────────────────────────────────────────
func cacheKeyJenjangSelectList(search string) string {
	return fmt.Sprintf("kepegawaian:pendidikan:jenjang:selectlist:search%s", search)
}

// ─── Cache Prefix Constants ───────────────────────────────────────────────────────
// Digunakan untuk InvalidateList agar konsisten dan tidak typo
const (
	cachePrefixJenjangSelectList = "kepegawaian:pendidikan:jenjang:selectlist:"
)
