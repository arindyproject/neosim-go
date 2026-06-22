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
	cachePrefixGolonganDarahList    = "master_golongan_darah:list:"
	cachePrefixSukuList             = "master_suku:list:"
	cachePrefixJenisKelaminList     = "master_jenis_kelamin:list:"
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

// GolonganDarah---------------------------------------------------------------------
func cacheKeyGolonganDarahDetail(id int64) string {
	return fmt.Sprintf("master_golongan_darah:detail:%d", id)
}
func cacheKeyGolonganDarahList(page, pageSize int, filter *dto.FilterMasterGolonganDarahRequest) string {
	f := ""
	if filter != nil {
		f = filter.Name
	}
	return fmt.Sprintf("master_golongan_darah:list:p%d:ps%d:f%s", page, pageSize, f)
} // GolonganDarah---------------------------------------------------------------------

// Suku--------------------------------------------------------------------------------
func cacheKeySukuDetail(id int64) string {
	return fmt.Sprintf("master_suku:detail:%d", id)
}
func cacheKeySukuList(page, pageSize int, filter *dto.FilterMasterSukuRequest) string {
	f := ""
	if filter != nil {
		f = filter.Name
	}
	return fmt.Sprintf("master_suku:list:p%d:ps%d:f%s", page, pageSize, f)
} // Suku------------------------------------------------------------------------------

// JenisKelamin------------------------------------------------------------------------
func cacheKeyJenisKelaminDetail(id int64) string {
	return fmt.Sprintf("master_jenis_kelamin:detail:%d", id)
}
func cacheKeyJenisKelaminList(page, pageSize int, filter *dto.FilterMasterJenisKelaminRequest) string {
	f := ""
	if filter != nil {
		f = filter.Name
	}
	return fmt.Sprintf("master_jenis_kelamin:list:p%d:ps%d:f%s", page, pageSize, f)
} // JenisKelamin------------------------------------------------------------------------
