package services

import (
	"fmt"

	"neosim_go/internal/modules/master/master/dto"
)

// ─── Cache Prefix Constants ───────────────────────────────────────────────────────
// Digunakan untuk InvalidateList agar konsisten dan tidak typo
const (
	cachePrefixPekerjaanList        = "master_pekerjaan:list:"
	cachePrefixPendidikanList       = "master_pendidikan:list:"
	cachePrefixAgamaList            = "master_agama:list:"
	cachePrefixStatusPernikahanList = "master_status_pernikahan:list:"
)

// ─── Cache Keys Generator ───────────────────────────────────────────────────────
// Pekerjaan-----------------------------------------------------------------------
func cacheKeyPekerjaanDetail(id int64) string { return fmt.Sprintf("master_pekerjaan:detail:%d", id) }
func cacheKeyPekerjaanList(page, pageSize int, filter *dto.FilterMasterPekerjaanRequest) string {
	f := ""
	if filter != nil {
		f = filter.Name
	}
	return fmt.Sprintf("master_pekerjaan:list:p%d:ps%d:f%s", page, pageSize, f)
} // Pekerjaan---------------------------------------------------------------------

// Pendidikan---------------------------------------------------------------------
func cacheKeyPendidikanDetail(id int64) string { return fmt.Sprintf("master_pendidikan:detail:%d", id) }
func cacheKeyPendidikanList(page, pageSize int, filter *dto.FilterMasterPendidikanRequest) string {
	f := ""
	if filter != nil {
		f = filter.Name
	}
	return fmt.Sprintf("master_pendidikan:list:p%d:ps%d:f%s", page, pageSize, f)
} // Pendidikan---------------------------------------------------------------------

// Agama--------------------------------------------------------------------------
func cacheKeyAgamaDetail(id int64) string { return fmt.Sprintf("master_agama:detail:%d", id) }
func cacheKeyAgamaList(page, pageSize int, filter *dto.FilterMasterAgamaRequest) string {
	f := ""
	if filter != nil {
		f = filter.Name
	}
	return fmt.Sprintf("master_agama:list:p%d:ps%d:f%s", page, pageSize, f)
} // Agama--------------------------------------------------------------------------

// StatusPernikahan---------------------------------------------------------------
func cacheKeyStatusPernikahanDetail(id int64) string {
	return fmt.Sprintf("master_status_pernikahan:detail:%d", id)
}
func cacheKeyStatusPernikahanList(page, pageSize int, filter *dto.FilterMasterStatusPernikahanRequest) string {
	f := ""
	if filter != nil {
		f = filter.Name
	}
	return fmt.Sprintf("master_status_pernikahan:list:p%d:ps%d:f%s", page, pageSize, f)
} // StatusPernikahan---------------------------------------------------------------
