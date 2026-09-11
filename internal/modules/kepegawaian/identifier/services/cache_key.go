package services

import (
	"fmt"
)

// ──────────────────────────────────────────────────────────────────────────
func cacheKeyTipeSelectList(search string) string {
	return fmt.Sprintf("kepegawaian:identifier:tipe:selectlist:search%s", search)
}

// ─── Cache Prefix Constants ───────────────────────────────────────────────────────
// Digunakan untuk InvalidateList agar konsisten dan tidak typo
const (
	cachePrefixTipeSelectList = "kepegawaian:identifier:tipe:selectlist:"
)
