package services

import (
	"fmt"

	"neosim_go/internal/modules/master/alamat/dto"
)

// ─── Cache Keys Generator ─────────────────────────────────────────────────────────
//
// Konvensi penamaan:
//   cacheKey<Entitas>Detail  → untuk GetByID (response sederhana)
//   cacheKey<Entitas>GetDetail → untuk GetDetail (response dengan aggregate/join)
//   cacheKey<Entitas>List    → untuk List (paginated)

// Negara ──────────────────────────────────────────────────────────────────────────
func cacheKeyNegaraDetail(id int64) string {
	return fmt.Sprintf("master_alamat:negara:detail:%d", id)
}

func cacheKeyNegaraList(page, pageSize int, filter *dto.FilterNegaraRequest) string {
	f := ""
	if filter != nil {
		f = filter.Name
	}
	return fmt.Sprintf("master_alamat:negara:list:p%d:ps%d:f%s", page, pageSize, f)
}

// Provinsi ────────────────────────────────────────────────────────────────────────
func cacheKeyProvinsiDetail(id int64) string {
	return fmt.Sprintf("master_alamat:provinsi:detail:%d", id)
}

func cacheKeyProvinsiGetDetail(id int64) string {
	return fmt.Sprintf("master_alamat:provinsi:getdetail:%d", id)
}

func cacheKeyProvinsiList(page, pageSize int, negaraID *int64, filter *dto.FilterProvinsiRequest) string {
	nID := "all"
	if negaraID != nil {
		nID = fmt.Sprintf("%d", *negaraID)
	}
	f := ""
	if filter != nil {
		f = filter.Name
	}
	return fmt.Sprintf("master_alamat:provinsi:list:p%d:ps%d:nid%s:f%s", page, pageSize, nID, f)
}

// Kota/Kabupaten ──────────────────────────────────────────────────────────────────
func cacheKeyKotaDetail(id int64) string {
	return fmt.Sprintf("master_alamat:kota:detail:%d", id)
}

func cacheKeyKotaGetDetail(id int64) string {
	return fmt.Sprintf("master_alamat:kota:getdetail:%d", id)
}

func cacheKeyKotaList(page, pageSize int, provinsiID *int64, filter *dto.FilterKotaKabupatenRequest) string {
	pID := "all"
	if provinsiID != nil {
		pID = fmt.Sprintf("%d", *provinsiID)
	}
	f := ""
	if filter != nil {
		f = filter.Name
	}
	return fmt.Sprintf("master_alamat:kota:list:p%d:ps%d:pid%s:f%s", page, pageSize, pID, f)
}

// Kecamatan ───────────────────────────────────────────────────────────────────────
func cacheKeyKecamatanDetail(id int64) string {
	return fmt.Sprintf("master_alamat:kecamatan:detail:%d", id)
}

func cacheKeyKecamatanGetDetail(id int64) string {
	return fmt.Sprintf("master_alamat:kecamatan:getdetail:%d", id)
}

func cacheKeyKecamatanList(page, pageSize int, kotaID *int64, filter *dto.FilterKecamatanRequest) string {
	kID := "all"
	if kotaID != nil {
		kID = fmt.Sprintf("%d", *kotaID)
	}
	f := ""
	if filter != nil {
		f = filter.Name
	}
	return fmt.Sprintf("master_alamat:kecamatan:list:p%d:ps%d:kid%s:f%s", page, pageSize, kID, f)
}

// Kelurahan/Desa ──────────────────────────────────────────────────────────────────
func cacheKeyDesaDetail(id int64) string {
	return fmt.Sprintf("master_alamat:desa:detail:%d", id)
}

func cacheKeyDesaGetDetail(id int64) string {
	return fmt.Sprintf("master_alamat:desa:getdetail:%d", id)
}

func cacheKeyDesaList(page, pageSize int, kecamatanID *int64, filter *dto.FilterKelurahanDesaRequest) string {
	kecID := "all"
	if kecamatanID != nil {
		kecID = fmt.Sprintf("%d", *kecamatanID)
	}
	f := ""
	if filter != nil {
		f = filter.Name
	}
	return fmt.Sprintf("master_alamat:desa:list:p%d:ps%d:kecid%s:f%s", page, pageSize, kecID, f)
}

// ─── Cache Prefix Constants ───────────────────────────────────────────────────────
// Digunakan untuk InvalidateList agar konsisten dan tidak typo
const (
	cachePrefixNegaraList    = "master_alamat:negara:list:"
	cachePrefixProvinsiList  = "master_alamat:provinsi:list:"
	cachePrefixKotaList      = "master_alamat:kota:list:"
	cachePrefixKecamatanList = "master_alamat:kecamatan:list:"
	cachePrefixDesaList      = "master_alamat:desa:list:"
)
