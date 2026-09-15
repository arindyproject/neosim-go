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
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk membuat KepegawaianAlamat baru.", nil)
	}

	// validasi keberadaan Pegawai
	pegawaiMaster, err := s.pegawaiRepo.GetPegawaiByID(ctx, req.PegawaiID)
	if err != nil {
		return nil, appErrors.Internal("gagal mengambil data pegawai")
	}
	if pegawaiMaster == nil {
		return nil, appErrors.Wrap(http.StatusUnprocessableEntity, "ID Pegawai tidak ditemukan.", nil)
	}

	// validasi keberadaan master Tipe
	tipeMaster, err := s.repo.GetTipeByID(ctx, req.TipeID)
	if err != nil {
		return nil, appErrors.Internal("gagal mengambil Tipe Alamat")
	}
	if tipeMaster == nil {
		return nil, appErrors.Wrap(http.StatusUnprocessableEntity, "Tipe Alamat tidak ditemukan.", nil)
	}

	// ── VALIDASI HIRARKI ALAMAT (OPSIONAL) ───────────────────────────────────

	// 1. Validasi Negara (jika diisi)
	if req.NegaraID != nil {
		negaraMaster, err := s.masterAlamatRepo.GetByIDNegara(ctx, *req.NegaraID)
		if err != nil {
			return nil, appErrors.Internal("gagal mengambil Data Negara")
		}
		if negaraMaster == nil {
			return nil, appErrors.Wrap(http.StatusUnprocessableEntity, "ID Negara tidak ditemukan.", nil)
		}

		// 2. Validasi Provinsi (hanya jika Negara dan Provinsi diisi)
		if req.ProvinsiID != nil {
			provinsiMaster, err := s.masterAlamatRepo.CheckProvinsi(ctx, *req.NegaraID, *req.ProvinsiID)
			if err != nil {
				return nil, appErrors.Internal("gagal mengambil Data Provinsi")
			}
			if !provinsiMaster {
				return nil, appErrors.Wrap(http.StatusUnprocessableEntity, "ID Provinsi tidak ditemukan atau tidak cocok dengan Negara yang dipilih.", nil)
			}

			// 3. Validasi Kota/Kabupaten (hanya jika Provinsi dan KotaKabupaten diisi)
			if req.KotaKabupatenID != nil {
				kotaKabupatenMaster, err := s.masterAlamatRepo.CheckKotaKabupaten(ctx, *req.ProvinsiID, *req.KotaKabupatenID)
				if err != nil {
					return nil, appErrors.Internal("gagal mengambil Data Kota/Kabupaten")
				}
				if !kotaKabupatenMaster {
					return nil, appErrors.Wrap(http.StatusUnprocessableEntity, "ID Kota/Kabupaten tidak ditemukan atau tidak cocok dengan Provinsi yang dipilih.", nil)
				}

				// 4. Validasi Kecamatan (hanya jika KotaKabupaten dan Kecamatan diisi)
				if req.KecamatanID != nil {
					kecamatanMaster, err := s.masterAlamatRepo.CheckKecamatan(ctx, *req.KotaKabupatenID, *req.KecamatanID)
					if err != nil {
						return nil, appErrors.Internal("gagal mengambil Data Kecamatan")
					}
					if !kecamatanMaster {
						return nil, appErrors.Wrap(http.StatusUnprocessableEntity, "ID Kecamatan tidak ditemukan atau tidak cocok dengan Kota/Kabupaten yang dipilih.", nil)
					}

					// 5. Validasi Kelurahan/Desa (hanya jika Kecamatan dan KelurahanDesa diisi)
					if req.KelurahanDesaID != nil {
						kelurahanDesaMaster, err := s.masterAlamatRepo.CheckKelurahanDesa(ctx, *req.KecamatanID, *req.KelurahanDesaID)
						if err != nil {
							return nil, appErrors.Internal("gagal mengambil Data Kelurahan/Desa")
						}
						if !kelurahanDesaMaster {
							return nil, appErrors.Wrap(http.StatusUnprocessableEntity, "ID Kelurahan/Desa tidak ditemukan atau tidak cocok dengan Kecamatan yang dipilih.", nil)
						}
					}
				}
			}
		}
	} else {
		// Guard optional: Jika Negara nil, tapi user nekad kirim ID Provinsi/Kabupaten/dll.
		if req.ProvinsiID != nil || req.KotaKabupatenID != nil || req.KecamatanID != nil || req.KelurahanDesaID != nil {
			return nil, appErrors.Wrap(http.StatusBadRequest, "ID Negara wajib diisi jika ingin memilih Wilayah/Provinsi.", nil)
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
	if err := s.repo.CreateAlamat(ctx, m); err != nil {
		return nil, err
	}

	creator := s.buildCreator(ctx, m.CreatedBy)

	return dto.ToKepegawaianAlamatResponse(dto.KepegawaianAlamatResponseParams{
		KepegawaianAlamat: m,
		Creator:           creator,
		Updater:           creator,
	}), nil
}

// ── GetByID ───────────────────────────────────────────────────────────────────
func (s *service) GetAlamatByID(ctx context.Context, id int64, actor he.AuthContext) (*dto.KepegawaianAlamatResponse, error) {
	can, err := s.canReadKepegawaianAlamat(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
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

	return dto.ToKepegawaianAlamatResponse(dto.KepegawaianAlamatResponseParams{
		KepegawaianAlamat: m,
		Creator:           creator,
		Updater:           updater,
	}), nil
}

// ─── GetByPegawaiID ───────────────────────────────────────────────────────────
func (s *service) GetAlamatByPegawaiID(ctx context.Context, pegawaiID int64, page, pageSize int, actor he.AuthContext) ([]dto.KepegawaianAlamatResponse, int64, error) {
	can, err := s.canReadKepegawaianAlamat(ctx, actor)
	if err != nil {
		return nil, 0, appErrors.Internal("gagal cek akses")
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
	return dto.ToKepegawaianAlamatListResponse(items, creatorsMap, updatersMap), total, nil
}

// ── List ──────────────────────────────────────────────────────────────────────
func (s *service) ListAlamat(ctx context.Context, page, pageSize int, filter *dto.FilterKepegawaianAlamatRequest, actor he.AuthContext) ([]dto.KepegawaianAlamatResponse, int64, error) {
	can, err := s.canReadKepegawaianAlamat(ctx, actor)
	if err != nil {
		return nil, 0, appErrors.Internal("gagal cek akses")
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
		return nil, 0, err
	}

	creatorsMap, updatersMap := s.buildAuditMaps(ctx, items)
	return dto.ToKepegawaianAlamatListResponse(items, creatorsMap, updatersMap), total, nil
}

// ── Update ────────────────────────────────────────────────────────────────────
func (s *service) UpdateAlamat(ctx context.Context, id int64, req *dto.UpdateKepegawaianAlamatRequest, actor he.AuthContext) (*dto.KepegawaianAlamatResponse, error) {
	if req == nil {
		return nil, appErrors.Wrap(http.StatusBadRequest, "Data request tidak boleh kosong.", nil)
	}

	can, err := s.canUpdateKepegawaianAlamat(ctx, actor)
	if err != nil {
		return nil, appErrors.Internal("gagal cek akses")
	}
	if !can {
		return nil, appErrors.Wrap(http.StatusForbidden,
			"Akses ditolak. Anda tidak memiliki hak akses untuk mengubah KepegawaianAlamat.", nil)
	}

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
			return nil, appErrors.Internal("gagal mengambil Tipe Alamat")
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

	// 3. Validasi Hirarki Alamat berdasarkan nilai efektif
	if effectiveNegaraID != nil {
		negaraMaster, err := s.masterAlamatRepo.GetByIDNegara(ctx, *effectiveNegaraID)
		if err != nil {
			return nil, appErrors.Internal("gagal mengambil Data Negara")
		}
		if negaraMaster == nil {
			return nil, appErrors.Wrap(http.StatusUnprocessableEntity, "ID Negara tidak ditemukan.", nil)
		}

		if effectiveProvinsiID != nil {
			provinsiMaster, err := s.masterAlamatRepo.CheckProvinsi(ctx, *effectiveNegaraID, *effectiveProvinsiID)
			if err != nil {
				return nil, appErrors.Internal("gagal mengambil Data Provinsi")
			}
			if !provinsiMaster {
				return nil, appErrors.Wrap(http.StatusUnprocessableEntity, "ID Provinsi tidak ditemukan atau tidak cocok dengan Negara yang dipilih.", nil)
			}

			if effectiveKotaKabupatenID != nil {
				kotaKabupatenMaster, err := s.masterAlamatRepo.CheckKotaKabupaten(ctx, *effectiveProvinsiID, *effectiveKotaKabupatenID)
				if err != nil {
					return nil, appErrors.Internal("gagal mengambil Data Kota/Kabupaten")
				}
				if !kotaKabupatenMaster {
					return nil, appErrors.Wrap(http.StatusUnprocessableEntity, "ID Kota/Kabupaten tidak ditemukan atau tidak cocok dengan Provinsi yang dipilih.", nil)
				}

				if effectiveKecamatanID != nil {
					kecamatanMaster, err := s.masterAlamatRepo.CheckKecamatan(ctx, *effectiveKotaKabupatenID, *effectiveKecamatanID)
					if err != nil {
						return nil, appErrors.Internal("gagal mengambil Data Kecamatan")
					}
					if !kecamatanMaster {
						return nil, appErrors.Wrap(http.StatusUnprocessableEntity, "ID Kecamatan tidak ditemukan atau tidak cocok dengan Kota/Kabupaten yang dipilih.", nil)
					}

					if effectiveKelurahanDesaID != nil {
						kelurahanDesaMaster, err := s.masterAlamatRepo.CheckKelurahanDesa(ctx, *effectiveKecamatanID, *effectiveKelurahanDesaID)
						if err != nil {
							return nil, appErrors.Internal("gagal mengambil Data Kelurahan/Desa")
						}
						if !kelurahanDesaMaster {
							return nil, appErrors.Wrap(http.StatusUnprocessableEntity, "ID Kelurahan/Desa tidak ditemukan atau tidak cocok dengan Kecamatan yang dipilih.", nil)
						}
					}
				}
			}
		}
	} else {
		if effectiveProvinsiID != nil || effectiveKotaKabupatenID != nil || effectiveKecamatanID != nil || effectiveKelurahanDesaID != nil {
			return nil, appErrors.Wrap(http.StatusBadRequest, "ID Negara wajib diisi jika ingin memilih Wilayah/Provinsi.", nil)
		}
	}

	// 4. Update data entity
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

	if err := s.repo.UpdateAlamat(ctx, m); err != nil {
		return nil, err
	}

	creator := s.buildCreator(ctx, m.CreatedBy)
	updater := s.buildCreator(ctx, m.UpdatedBy)

	return dto.ToKepegawaianAlamatResponse(dto.KepegawaianAlamatResponseParams{
		KepegawaianAlamat: m,
		Creator:           creator,
		Updater:           updater,
	}), nil
}

// ── Delete ────────────────────────────────────────────────────────────────────
func (s *service) DeleteAlamat(ctx context.Context, id int64, actor he.AuthContext) error {
	can, err := s.canDeleteKepegawaianAlamat(ctx, actor)
	if err != nil {
		return appErrors.Internal("gagal cek akses")
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
	return s.repo.DeleteAlamat(ctx, id, actor.UserID)
}
