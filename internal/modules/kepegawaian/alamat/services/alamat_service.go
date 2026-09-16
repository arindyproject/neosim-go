package services

import (
	"context"
	"errors"
	"net/http"
	"time"

	"neosim_go/internal/modules/kepegawaian/alamat/dto"
	"neosim_go/internal/modules/kepegawaian/alamat/models"
	appErrors "neosim_go/internal/shared/errors"
	he "neosim_go/internal/shared/httputil"
)

// ── Create ────────────────────────────────────────────────────────────────────
func (s *service) CreateAlamat(ctx context.Context, req *dto.CreateKepegawaianAlamatRequest, actor he.AuthContext) (*dto.KepegawaianAlamatResponse, error) {
	can, err := s.canCreateKepegawaianAlamat(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses : " + err.Error())
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk membuat KepegawaianAlamat baru.", nil)
	}

	// ── NORMALISASI: anggap 0 sebagai kosong/null ────────────────────────────
	normalizeZeroID(&req.NegaraID)
	normalizeZeroID(&req.ProvinsiID)
	normalizeZeroID(&req.KotaKabupatenID)
	normalizeZeroID(&req.KecamatanID)
	normalizeZeroID(&req.KelurahanDesaID)

	// validasi keberadaan Pegawai
	pegawaiMaster, err := s.pegawaiRepo.GetPegawaiByID(ctx, req.PegawaiID)
	if err != nil {
		return nil, appErrors.Internal("gagal mengambil data pegawai : " + err.Error())
	}
	if pegawaiMaster == nil {
		return nil, appErrors.Wrap(http.StatusUnprocessableEntity, "ID Pegawai tidak ditemukan.", nil)
	}

	// validasi keberadaan master Tipe
	tipeMaster, err := s.repo.GetTipeByID(ctx, req.TipeID)
	if err != nil {
		return nil, appErrors.Internal("gagal mengambil Tipe Alamat : " + err.Error())
	}
	if tipeMaster == nil {
		return nil, appErrors.Wrap(http.StatusUnprocessableEntity, "Tipe Alamat tidak ditemukan.", nil)
	}

	// ── VALIDASI HIRARKI ALAMAT (OPSIONAL, BERJENJANG) ──────────────────────
	if req.NegaraID == nil {
		if req.ProvinsiID != nil || req.KotaKabupatenID != nil || req.KecamatanID != nil || req.KelurahanDesaID != nil {
			return nil, appErrors.Wrap(http.StatusBadRequest, "ID Negara wajib diisi jika ingin memilih Wilayah/Provinsi.", nil)
		}
	} else {
		negaraMaster, err := s.masterAlamatRepo.GetByIDNegara(ctx, *req.NegaraID)
		if err != nil {
			return nil, appErrors.Internal("gagal mengambil Data Negara : " + err.Error())
		}
		if negaraMaster == nil {
			return nil, appErrors.Wrap(http.StatusUnprocessableEntity, "ID Negara tidak ditemukan.", nil)
		}

		if req.ProvinsiID == nil {
			if req.KotaKabupatenID != nil || req.KecamatanID != nil || req.KelurahanDesaID != nil {
				return nil, appErrors.Wrap(http.StatusBadRequest, "ID Provinsi wajib diisi jika ingin memilih Kota/Kabupaten.", nil)
			}
		} else {
			ok, err := s.masterAlamatRepo.CheckProvinsi(ctx, *req.NegaraID, *req.ProvinsiID)
			if err != nil {
				return nil, appErrors.Internal("gagal mengambil Data Provinsi : " + err.Error())
			}
			if !ok {
				return nil, appErrors.Wrap(http.StatusUnprocessableEntity, "ID Provinsi tidak ditemukan atau tidak cocok dengan Negara yang dipilih.", nil)
			}

			if req.KotaKabupatenID == nil {
				if req.KecamatanID != nil || req.KelurahanDesaID != nil {
					return nil, appErrors.Wrap(http.StatusBadRequest, "ID Kota/Kabupaten wajib diisi jika ingin memilih Kecamatan.", nil)
				}
			} else {
				ok, err := s.masterAlamatRepo.CheckKotaKabupaten(ctx, *req.ProvinsiID, *req.KotaKabupatenID)
				if err != nil {
					return nil, appErrors.Internal("gagal mengambil Data Kota/Kabupaten : " + err.Error())
				}
				if !ok {
					return nil, appErrors.Wrap(http.StatusUnprocessableEntity, "ID Kota/Kabupaten tidak ditemukan atau tidak cocok dengan Provinsi yang dipilih.", nil)
				}

				if req.KecamatanID == nil {
					if req.KelurahanDesaID != nil {
						return nil, appErrors.Wrap(http.StatusBadRequest, "ID Kecamatan wajib diisi jika ingin memilih Kelurahan/Desa.", nil)
					}
				} else {
					ok, err := s.masterAlamatRepo.CheckKecamatan(ctx, *req.KotaKabupatenID, *req.KecamatanID)
					if err != nil {
						return nil, appErrors.Internal("gagal mengambil Data Kecamatan : " + err.Error())
					}
					if !ok {
						return nil, appErrors.Wrap(http.StatusUnprocessableEntity, "ID Kecamatan tidak ditemukan atau tidak cocok dengan Kota/Kabupaten yang dipilih.", nil)
					}

					if req.KelurahanDesaID != nil {
						ok, err := s.masterAlamatRepo.CheckKelurahanDesa(ctx, *req.KecamatanID, *req.KelurahanDesaID)
						if err != nil {
							return nil, appErrors.Internal("gagal mengambil Data Kelurahan/Desa : " + err.Error())
						}
						if !ok {
							return nil, appErrors.Wrap(http.StatusUnprocessableEntity, "ID Kelurahan/Desa tidak ditemukan atau tidak cocok dengan Kecamatan yang dipilih.", nil)
						}
					}
				}
			}
		}
	}

	m := &models.KepegawaianAlamat{
		PegawaiID: req.PegawaiID,
		TipeID:    req.TipeID,
		Jalan:     req.Jalan,
		RT:        req.RT,
		RW:        req.RW,
		KodePos:   req.KodePos,

		NegaraID:        req.NegaraID,
		ProvinsiID:      req.ProvinsiID,
		KotaKabupatenID: req.KotaKabupatenID,
		KecamatanID:     req.KecamatanID,
		KelurahanDesaID: req.KelurahanDesaID,

		IsPrimary:   req.IsPrimary,
		Description: req.Description,
		CreatedBy:   &actor.UserID,
		UpdatedBy:   &actor.UserID,
	}

	// ── CEK DUPLIKASI ALAMAT (per pegawai, aman untuk pegawai lain) ─────────
	isDup, err := s.repo.CheckDuplicateAlamat(ctx, m, nil)
	if err != nil {
		return nil, appErrors.Internal("gagal cek duplikasi alamat : " + err.Error())
	}
	if isDup {
		return nil, appErrors.Wrap(http.StatusConflict, "Pegawai ini sudah memiliki alamat yang sama.", nil)
	}

	// ── SET PRIMARY: jika is_primary = true, unset primary lama milik pegawai ini ──
	if req.IsPrimary {
		if err := s.repo.UnsetPrimaryAlamatByPegawaiID(ctx, req.PegawaiID, actor.UserID); err != nil {
			return nil, appErrors.Internal("gagal mereset alamat primary sebelumnya : " + err.Error())
		}
	}

	if err := s.repo.CreateAlamat(ctx, m); err != nil {
		return nil, err
	}

	creator := s.buildCreator(ctx, m.CreatedBy)

	negara, provinsi, kota, kecamatan, kelurahan := s.buildWilayahSingle(ctx, m)
	return dto.ToKepegawaianAlamatResponse(dto.KepegawaianAlamatResponseParams{
		KepegawaianAlamat: m,
		Creator:           creator,
		Updater:           creator,
		Negara:            negara,
		Provinsi:          provinsi,
		KotaKabupaten:     kota,
		Kecamatan:         kecamatan,
		KelurahanDesa:     kelurahan,
	}), nil
}

// ── GetByID ───────────────────────────────────────────────────────────────────
func (s *service) GetAlamatByID(ctx context.Context, id int64, actor he.AuthContext) (*dto.KepegawaianAlamatResponse, error) {
	can, err := s.canReadKepegawaianAlamat(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses : " + err.Error())
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk Melihat KepegawaianAlamat.", nil)
	}

	m, err := s.repo.GetAlamatByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errors.New("KepegawaianAlamat tidak ditemukan")
	}

	creator := s.buildCreator(ctx, m.CreatedBy)
	updater := s.buildCreator(ctx, m.UpdatedBy)
	negara, provinsi, kota, kecamatan, kelurahan := s.buildWilayahSingle(ctx, m)
	return dto.ToKepegawaianAlamatResponse(dto.KepegawaianAlamatResponseParams{
		KepegawaianAlamat: m,
		Creator:           creator,
		Updater:           updater,
		Negara:            negara,
		Provinsi:          provinsi,
		KotaKabupaten:     kota,
		Kecamatan:         kecamatan,
		KelurahanDesa:     kelurahan,
	}), nil
}

// ─── GetByPegawaiID ───────────────────────────────────────────────────────────
func (s *service) GetAlamatByPegawaiID(ctx context.Context, pegawaiID int64, page, pageSize int, actor he.AuthContext) ([]dto.KepegawaianAlamatResponse, int64, error) {
	can, err := s.canReadKepegawaianAlamat(ctx, actor)
	if err != nil {
		return nil, 0, appErrors.Internal("gagal cek akses : " + err.Error())
	}
	if !can {
		return nil, 0, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk Melihat Kepegawaian Alamat.", nil)
	}

	items, total, err := s.repo.GetAlamatByPegawaiID(ctx, pegawaiID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	if items == nil {
		return nil, 0, errors.New("Kepegawaian Alamat tidak ditemukan")
	}

	creatorsMap, updatersMap := s.buildAuditMaps(ctx, items)
	negaraMap, provinsiMap, kotaMap, kecamatanMap, kelurahanMap := s.buildWilayahMaps(ctx, items)
	return dto.ToKepegawaianAlamatListResponse(items, creatorsMap, updatersMap,
		negaraMap, provinsiMap, kotaMap, kecamatanMap, kelurahanMap), total, nil
}

// ── List ──────────────────────────────────────────────────────────────────────
func (s *service) ListAlamat(ctx context.Context, page, pageSize int, filter *dto.FilterKepegawaianAlamatRequest, actor he.AuthContext) ([]dto.KepegawaianAlamatResponse, int64, error) {
	can, err := s.canReadKepegawaianAlamat(ctx, actor)

	if err != nil {
		return nil, 0, appErrors.Internal("gagal cek akses : " + err.Error())
	}
	if !can {
		return nil, 0, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk melihat daftar KepegawaianAlamat.", nil)
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > s.cfg.DefaultPageSizeMax {
		pageSize = s.cfg.DefaultPageSizeMax
	}

	items, total, err := s.repo.ListAlamat(ctx, page, pageSize, filter)
	if err != nil {
		return nil, 0, appErrors.Internal("gagal saat mengambil data dari Repo : " + err.Error())
	}

	creatorsMap, updatersMap := s.buildAuditMaps(ctx, items)
	negaraMap, provinsiMap, kotaMap, kecamatanMap, kelurahanMap := s.buildWilayahMaps(ctx, items)
	return dto.ToKepegawaianAlamatListResponse(items, creatorsMap, updatersMap,
		negaraMap, provinsiMap, kotaMap, kecamatanMap, kelurahanMap), total, nil
}

// ── Update ────────────────────────────────────────────────────────────────────
func (s *service) UpdateAlamat(ctx context.Context, id int64, req *dto.UpdateKepegawaianAlamatRequest, actor he.AuthContext) (*dto.KepegawaianAlamatResponse, error) {
	if req == nil {
		return nil, appErrors.Wrap(http.StatusBadRequest, "Data request tidak boleh kosong.", nil)
	}

	can, err := s.canUpdateKepegawaianAlamat(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses : " + err.Error())
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk mengubah KepegawaianAlamat.", nil)
	}

	// ── NORMALISASI: anggap 0 sebagai kosong/null ────────────────────────────
	normalizeZeroID(&req.NegaraID)
	normalizeZeroID(&req.ProvinsiID)
	normalizeZeroID(&req.KotaKabupatenID)
	normalizeZeroID(&req.KecamatanID)
	normalizeZeroID(&req.KelurahanDesaID)

	m, err := s.repo.GetAlamatByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, appErrors.Wrap(http.StatusNotFound, "KepegawaianAlamat tidak ditemukan.", nil)
	}

	// 1. Validasi TipeID jika diubah
	if req.TipeID != nil {
		tipeMaster, err := s.repo.GetTipeByID(ctx, *req.TipeID)
		if err != nil {
			return nil, appErrors.Internal("gagal mengambil Tipe Alamat : " + err.Error())
		}
		if tipeMaster == nil {
			return nil, appErrors.Wrap(http.StatusUnprocessableEntity, "Tipe Alamat tidak ditemukan.", nil)
		}
		m.TipeID = *req.TipeID
	}

	// 2. Tentukan nilai efektif (gabungan antara input request baru & data existing DB)
	effectiveNegaraID := m.NegaraID
	if req.NegaraID != nil {
		effectiveNegaraID = req.NegaraID
	}

	effectiveProvinsiID := m.ProvinsiID
	if req.ProvinsiID != nil {
		effectiveProvinsiID = req.ProvinsiID
	}

	effectiveKotaKabupatenID := m.KotaKabupatenID
	if req.KotaKabupatenID != nil {
		effectiveKotaKabupatenID = req.KotaKabupatenID
	}

	effectiveKecamatanID := m.KecamatanID
	if req.KecamatanID != nil {
		effectiveKecamatanID = req.KecamatanID
	}

	effectiveKelurahanDesaID := m.KelurahanDesaID
	if req.KelurahanDesaID != nil {
		effectiveKelurahanDesaID = req.KelurahanDesaID
	}

	// 3. Validasi Hirarki Alamat berdasarkan nilai efektif (berjenjang, tanpa "lompat" level)
	if effectiveNegaraID == nil {
		if effectiveProvinsiID != nil || effectiveKotaKabupatenID != nil || effectiveKecamatanID != nil || effectiveKelurahanDesaID != nil {
			return nil, appErrors.Wrap(http.StatusBadRequest, "ID Negara wajib diisi jika ingin memilih Wilayah/Provinsi.", nil)
		}
	} else {
		negaraMaster, err := s.masterAlamatRepo.GetByIDNegara(ctx, *effectiveNegaraID)
		if err != nil {
			return nil, appErrors.Internal("gagal mengambil Data Negara : " + err.Error())
		}
		if negaraMaster == nil {
			return nil, appErrors.Wrap(http.StatusUnprocessableEntity, "ID Negara tidak ditemukan.", nil)
		}

		if effectiveProvinsiID == nil {
			if effectiveKotaKabupatenID != nil || effectiveKecamatanID != nil || effectiveKelurahanDesaID != nil {
				return nil, appErrors.Wrap(http.StatusBadRequest, "ID Provinsi wajib diisi jika ingin memilih Kota/Kabupaten.", nil)
			}
		} else {
			ok, err := s.masterAlamatRepo.CheckProvinsi(ctx, *effectiveNegaraID, *effectiveProvinsiID)
			if err != nil {
				return nil, appErrors.Internal("gagal mengambil Data Provinsi : " + err.Error())
			}
			if !ok {
				return nil, appErrors.Wrap(http.StatusUnprocessableEntity, "ID Provinsi tidak ditemukan atau tidak cocok dengan Negara yang dipilih.", nil)
			}

			if effectiveKotaKabupatenID == nil {
				if effectiveKecamatanID != nil || effectiveKelurahanDesaID != nil {
					return nil, appErrors.Wrap(http.StatusBadRequest, "ID Kota/Kabupaten wajib diisi jika ingin memilih Kecamatan.", nil)
				}
			} else {
				ok, err := s.masterAlamatRepo.CheckKotaKabupaten(ctx, *effectiveProvinsiID, *effectiveKotaKabupatenID)
				if err != nil {
					return nil, appErrors.Internal("gagal mengambil Data Kota/Kabupaten : " + err.Error())
				}
				if !ok {
					return nil, appErrors.Wrap(http.StatusUnprocessableEntity, "ID Kota/Kabupaten tidak ditemukan atau tidak cocok dengan Provinsi yang dipilih.", nil)
				}

				if effectiveKecamatanID == nil {
					if effectiveKelurahanDesaID != nil {
						return nil, appErrors.Wrap(http.StatusBadRequest, "ID Kecamatan wajib diisi jika ingin memilih Kelurahan/Desa.", nil)
					}
				} else {
					ok, err := s.masterAlamatRepo.CheckKecamatan(ctx, *effectiveKotaKabupatenID, *effectiveKecamatanID)
					if err != nil {
						return nil, appErrors.Internal("gagal mengambil Data Kecamatan : " + err.Error())
					}
					if !ok {
						return nil, appErrors.Wrap(http.StatusUnprocessableEntity, "ID Kecamatan tidak ditemukan atau tidak cocok dengan Kota/Kabupaten yang dipilih.", nil)
					}

					if effectiveKelurahanDesaID != nil {
						ok, err := s.masterAlamatRepo.CheckKelurahanDesa(ctx, *effectiveKecamatanID, *effectiveKelurahanDesaID)
						if err != nil {
							return nil, appErrors.Internal("gagal mengambil Data Kelurahan/Desa : " + err.Error())
						}
						if !ok {
							return nil, appErrors.Wrap(http.StatusUnprocessableEntity, "ID Kelurahan/Desa tidak ditemukan atau tidak cocok dengan Kecamatan yang dipilih.", nil)
						}
					}
				}
			}
		}
	}

	// 4. Update data entity
	// ── SET PRIMARY: jika diubah jadi true, unset primary lama milik pegawai ini ──
	if req.IsPrimary != nil {
		if *req.IsPrimary {
			if err := s.repo.UnsetPrimaryAlamatByPegawaiID(ctx, m.PegawaiID, actor.UserID); err != nil {
				return nil, appErrors.Internal("gagal mereset alamat primary sebelumnya : " + err.Error())
			}
		}
		m.IsPrimary = *req.IsPrimary
	}

	if req.Jalan != nil {
		m.Jalan = *req.Jalan
	}
	if req.RT != nil {
		m.RT = req.RT
	}
	if req.RW != nil {
		m.RW = req.RW
	}
	if req.KodePos != nil {
		m.KodePos = req.KodePos
	}

	m.NegaraID = effectiveNegaraID
	m.ProvinsiID = effectiveProvinsiID
	m.KotaKabupatenID = effectiveKotaKabupatenID
	m.KecamatanID = effectiveKecamatanID
	m.KelurahanDesaID = effectiveKelurahanDesaID

	if req.Description != nil {
		m.Description = req.Description
	}
	m.UpdatedBy = &actor.UserID
	m.UpdatedAt = time.Now()

	// ── CEK DUPLIKASI ALAMAT (per pegawai, kecualikan record ini sendiri) ───
	isDup, err := s.repo.CheckDuplicateAlamat(ctx, m, &id)
	if err != nil {
		return nil, appErrors.Internal("gagal cek duplikasi alamat : " + err.Error())
	}
	if isDup {
		return nil, appErrors.Wrap(http.StatusConflict, "Pegawai ini sudah memiliki alamat yang sama.", nil)
	}

	if err := s.repo.UpdateAlamat(ctx, m); err != nil {
		return nil, err
	}

	creator := s.buildCreator(ctx, m.CreatedBy)
	updater := s.buildCreator(ctx, m.UpdatedBy)

	negara, provinsi, kota, kecamatan, kelurahan := s.buildWilayahSingle(ctx, m)
	return dto.ToKepegawaianAlamatResponse(dto.KepegawaianAlamatResponseParams{
		KepegawaianAlamat: m,
		Creator:           creator,
		Updater:           updater,
		Negara:            negara,
		Provinsi:          provinsi,
		KotaKabupaten:     kota,
		Kecamatan:         kecamatan,
		KelurahanDesa:     kelurahan,
	}), nil
}

// ── Delete ────────────────────────────────────────────────────────────────────
func (s *service) DeleteAlamat(ctx context.Context, id int64, actor he.AuthContext) error {
	can, err := s.canDeleteKepegawaianAlamat(ctx, actor)
	if err != nil {
		return appErrors.Internal("gagal cek akses : " + err.Error())
	}
	if !can {
		return appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk menghapus KepegawaianAlamat.", nil)
	}

	m, err := s.repo.GetAlamatByID(ctx, id)
	if err != nil {
		return err
	}
	if m == nil {
		return errors.New("KepegawaianAlamat tidak ditemukan")
	}

	// alamat primary tidak boleh dihapus langsung —
	// harus ada alamat lain milik pegawai yang sama sebagai primary dulu
	if m.IsPrimary {
		others, err := s.repo.FindAlamatByPegawaiID(ctx, m.PegawaiID)
		if err != nil {
			return appErrors.Internal("gagal cek alamat lain : " + err.Error())
		}
		othersCount := 0
		for _, o := range others {
			if o.ID != id {
				othersCount++
			}
		}
		if othersCount > 0 {
			return appErrors.Wrap(http.StatusUnprocessableEntity,
				"Alamat ini adalah primary. Tetapkan alamat lain sebagai primary terlebih dahulu sebelum menghapus.", nil)
		}
	}

	return s.repo.DeleteAlamat(ctx, id, actor.UserID)
}

// normalizeZeroID menganggap 0 sebagai "tidak diisi" (nil), karena client
// terkadang mengirim 0 alih-alih null / tidak mengirim field sama sekali.
func normalizeZeroID(id **int64) {
	if *id != nil && **id == 0 {
		*id = nil
	}
}
