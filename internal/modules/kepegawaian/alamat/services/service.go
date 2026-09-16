package services

import (
	"context"
	"neosim_go/config"
	alamatContracts "neosim_go/internal/modules/kepegawaian/alamat/contracts"
	"neosim_go/internal/modules/kepegawaian/alamat/dto"

	authContracts "neosim_go/internal/modules/auth/contracts"
	"neosim_go/internal/modules/kepegawaian/alamat/models"
	pegawaiContracts "neosim_go/internal/modules/kepegawaian/pegawai/contracts"
	masterAlamatContracts "neosim_go/internal/modules/master/alamat/contracts"
	rbacContracts "neosim_go/internal/modules/rbac/contracts"
	userContracts "neosim_go/internal/modules/users/contracts"
	"neosim_go/internal/shared/cache"
	he "neosim_go/internal/shared/httputil"
)

// service adalah satu-satunya struct service untuk sub-module ini.
// Item baru (mode add-item) TIDAK membuat struct service baru — method
// CRUD & permission-nya ditempelkan langsung ke struct ini (mis.
// services/tag_service.go, services/tag_permission.go), dan repo field
// di bawah ini otomatis mencakup method item begitu contracts.Repository
// di-embed dengan interface repository item (lihat contracts/interfaces.go).
type service struct {
	repo             alamatContracts.Repository
	rbacRepo         rbacContracts.RBACRepository
	authRepo         authContracts.AuthRepository
	userRepo         userContracts.Repository
	pegawaiRepo      pegawaiContracts.Repository
	masterAlamatRepo masterAlamatContracts.Repository
	cfg              *config.Config
	cache            *cache.Manager // <--- Gunakan Cache Manager
}

// NewKepegawaianAlamatService membuat instance service baru
func NewKepegawaianAlamatService(
	repo alamatContracts.Repository,
	rbacRepo rbacContracts.RBACRepository,
	authRepo authContracts.AuthRepository,
	userRepo userContracts.Repository,
	pegawaiRepo pegawaiContracts.Repository,
	masterAlamatRepo masterAlamatContracts.Repository,
	cfg *config.Config,
	cacheManager *cache.Manager, // <--- Terima Cache Manager
) alamatContracts.Service {
	return &service{
		repo:             repo,
		rbacRepo:         rbacRepo,
		authRepo:         authRepo,
		userRepo:         userRepo,
		pegawaiRepo:      pegawaiRepo,
		masterAlamatRepo: masterAlamatRepo,
		cfg:              cfg,
		cache:            cacheManager,
	}
}

// buildCreator mengambil data creator user
func (s *service) buildCreator(ctx context.Context, createdBy *int64) *he.UserData {
	if createdBy == nil {
		return nil
	}
	creator, err := s.userRepo.GetByID(ctx, *createdBy)
	if err != nil || creator == nil {
		return nil
	}
	return &he.UserData{
		ID:       creator.ID,
		Username: creator.Username,
		Name:     creator.Name,
	}
}

// ── helper: build creator/updater maps ───────────────────────────────────────
func (s *service) buildAuditMaps(ctx context.Context, items []models.KepegawaianAlamat) (map[int64]*he.UserData, map[int64]*he.UserData) {
	idSet := make(map[int64]struct{})
	for _, item := range items {
		if item.CreatedBy != nil {
			idSet[*item.CreatedBy] = struct{}{}
		}
		if item.UpdatedBy != nil {
			idSet[*item.UpdatedBy] = struct{}{}
		}
	}
	ids := make([]int64, 0, len(idSet))
	for id := range idSet {
		ids = append(ids, id)
	}

	users, err := s.userRepo.GetByIDs(ctx, ids) // ← 1 query total, bukan 40
	if err != nil {
		return map[int64]*he.UserData{}, map[int64]*he.UserData{}
	}

	userMap := make(map[int64]*he.UserData, len(users))
	for _, u := range users {
		userMap[u.ID] = &he.UserData{ID: u.ID, Username: u.Username, Name: u.Name}
	}
	// creator dan updater sekarang share map yang sama — reuse otomatis, kode lebih pendek juga
	return userMap, userMap
}

// buildWilayahSingle mengambil detail Negara/Provinsi/Kota/Kecamatan/Kelurahan
// untuk SATU record alamat (dipakai di Create, Update, GetByID).
func (s *service) buildWilayahSingle(ctx context.Context, m *models.KepegawaianAlamat) (
	negara, provinsi, kota, kecamatan, kelurahan *dto.WilayahSimpelResponse,
) {
	if m.NegaraID != nil {
		if x, err := s.masterAlamatRepo.GetSimpelByIDNegara(ctx, *m.NegaraID); err == nil && x != nil {
			negara = &dto.WilayahSimpelResponse{ID: x.ID, Name: x.Name}
		}
	}
	if m.ProvinsiID != nil {
		if x, err := s.masterAlamatRepo.GetSimpelByIDProvinsi(ctx, *m.ProvinsiID); err == nil && x != nil {
			provinsi = &dto.WilayahSimpelResponse{ID: x.ID, Name: x.Name}
		}
	}
	if m.KotaKabupatenID != nil {
		if x, err := s.masterAlamatRepo.GetSimpelByIDKotaKabupaten(ctx, *m.KotaKabupatenID); err == nil && x != nil {
			kota = &dto.WilayahSimpelResponse{ID: x.ID, Name: x.Name}
		}
	}
	if m.KecamatanID != nil {
		if x, err := s.masterAlamatRepo.GetSimpelByIDKecamatan(ctx, *m.KecamatanID); err == nil && x != nil {
			kecamatan = &dto.WilayahSimpelResponse{ID: x.ID, Name: x.Name}
		}
	}
	if m.KelurahanDesaID != nil {
		if x, err := s.masterAlamatRepo.GetSimpelByIDKelurahanDesa(ctx, *m.KelurahanDesaID); err == nil && x != nil {
			kelurahan = &dto.WilayahSimpelResponse{ID: x.ID, Name: x.Name}
		}
	}
	return
}

// buildWilayahMaps mengambil detail wilayah untuk BANYAK record sekaligus
// (dipakai di ListAlamat / GetAlamatByPegawaiID) supaya tidak query berulang
// untuk ID yang sama.
func (s *service) buildWilayahMaps(ctx context.Context, items []models.KepegawaianAlamat) (
	negaraMap, provinsiMap, kotaMap, kecamatanMap, kelurahanMap map[int64]*dto.WilayahSimpelResponse,
) {
	negaraMap = map[int64]*dto.WilayahSimpelResponse{}
	provinsiMap = map[int64]*dto.WilayahSimpelResponse{}
	kotaMap = map[int64]*dto.WilayahSimpelResponse{}
	kecamatanMap = map[int64]*dto.WilayahSimpelResponse{}
	kelurahanMap = map[int64]*dto.WilayahSimpelResponse{}

	for _, m := range items {
		if m.NegaraID != nil {
			if _, ok := negaraMap[*m.NegaraID]; !ok {
				if x, err := s.masterAlamatRepo.GetSimpelByIDNegara(ctx, *m.NegaraID); err == nil && x != nil {
					negaraMap[*m.NegaraID] = &dto.WilayahSimpelResponse{ID: x.ID, Name: x.Name}
				}
			}
		}
		if m.ProvinsiID != nil {
			if _, ok := provinsiMap[*m.ProvinsiID]; !ok {
				if x, err := s.masterAlamatRepo.GetSimpelByIDProvinsi(ctx, *m.ProvinsiID); err == nil && x != nil {
					provinsiMap[*m.ProvinsiID] = &dto.WilayahSimpelResponse{ID: x.ID, Name: x.Name}
				}
			}
		}
		if m.KotaKabupatenID != nil {
			if _, ok := kotaMap[*m.KotaKabupatenID]; !ok {
				if x, err := s.masterAlamatRepo.GetSimpelByIDKotaKabupaten(ctx, *m.KotaKabupatenID); err == nil && x != nil {
					kotaMap[*m.KotaKabupatenID] = &dto.WilayahSimpelResponse{ID: x.ID, Name: x.Name}
				}
			}
		}
		if m.KecamatanID != nil {
			if _, ok := kecamatanMap[*m.KecamatanID]; !ok {
				if x, err := s.masterAlamatRepo.GetSimpelByIDKecamatan(ctx, *m.KecamatanID); err == nil && x != nil {
					kecamatanMap[*m.KecamatanID] = &dto.WilayahSimpelResponse{ID: x.ID, Name: x.Name}
				}
			}
		}
		if m.KelurahanDesaID != nil {
			if _, ok := kelurahanMap[*m.KelurahanDesaID]; !ok {
				if x, err := s.masterAlamatRepo.GetSimpelByIDKelurahanDesa(ctx, *m.KelurahanDesaID); err == nil && x != nil {
					kelurahanMap[*m.KelurahanDesaID] = &dto.WilayahSimpelResponse{ID: x.ID, Name: x.Name}
				}
			}
		}
	}
	return
}
