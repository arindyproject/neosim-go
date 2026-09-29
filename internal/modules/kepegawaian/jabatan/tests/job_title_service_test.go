package tests

import (
	"context"
	"fmt"
	"net/http"

	"github.com/stretchr/testify/mock"

	"neosim_go/internal/modules/kepegawaian/jabatan/dto"
	"neosim_go/internal/modules/kepegawaian/jabatan/models"
	"neosim_go/internal/modules/kepegawaian/jabatan/tests/factories"

	appErrors "neosim_go/internal/shared/errors"
)

// Catatan: suite KepegawaianJabatanServiceTestSuite, TestMain, dan helper
// (superadminActor, regularActor, mockNoPermissions, ptrTo, dst) sudah
// didefinisikan di jabatan_service_test.go. File ini HANYA menambah skenario
// test untuk JobTitle, memakai s.svc / s.repo yang SAMA.

func newCreateJobTitleReq() *dto.CreateJobTitleRequest {
	return &dto.CreateJobTitleRequest{
		Code:       "DR-SPESIALIS",
		Label:      "Dokter Spesialis",
		KategoriID: 1,
	}
}

// mockJobTitleCreateChecks men-stub urutan validasi CreateJobTitle: kategori
// valid, (opsional) rumpun profesi valid, dan kode belum dipakai.
func (s *KepegawaianJabatanServiceTestSuite) mockJobTitleCreateChecks(req *dto.CreateJobTitleRequest) {
	// PERBAIKAN: Gunakan mock.Anything untuk argumen ID agar kebal terhadap mismatch tipe (int/int64)
	// Pastikan nama method "CheckJobTitleKategori" SAMA PERSIS dengan yang ada di service.
	s.repo.On("CheckJobTitleKategori", mock.Anything, mock.Anything).
		Return(true, nil).Once()

	if req.RumpunProfesiID != nil {
		s.repo.On("CheckJobTitleRumpunProfesi", mock.Anything, mock.Anything).
			Return(true, nil).Once()
	}

	// Gunakan mock.Anything untuk argumen ketiga (exclude ID) agar aman dari mismatch tipe
	s.repo.On("ExistsByCode", req.Code, mock.Anything).
		Return(false, nil).Once()
}

// mockJobTitleCreateSaves mendaftarkan CreateJobTitle (mengisi ID) dan reload
// GetJobTitleByID(id) yang mengembalikan objek valid non-nil — service ini
// TIDAK punya fallback nil setelah reload, jadi wajib return objek asli.
func (s *KepegawaianJabatanServiceTestSuite) mockJobTitleCreateSaves(id int64, saved *models.JobTitle) {
	saved.ID = id
	s.repo.On("CreateJobTitle", mock.AnythingOfType("*models.JobTitle")).
		Run(func(args mock.Arguments) {
			args.Get(0).(*models.JobTitle).ID = id
		}).Return(nil)
	s.repo.On("GetJobTitleByID", id).Return(saved, nil)
}

// ── Create ───────────────────────────────────────────────────────────────────

func (s *KepegawaianJabatanServiceTestSuite) Test_CreateJobTitle_Superadmin_Success() {
	req := newCreateJobTitleReq()
	req.Description = ptrTo("Dokter dengan pendidikan spesialis")
	req.RumpunProfesiID = ptrTo(int64(2))
	req.MemerlukanSTR = true
	req.MemerlukanSIP = true
	req.JenjangMin = ptrTo("Sp1")
	req.FHIRCode = ptrTo("doctor")
	req.FHIRSystem = ptrTo("http://terminology.hl7.org/CodeSystem/practitioner-role")
	actor := superadminActor()

	s.mockJobTitleCreateChecks(req)
	saved := factories.NewJobTitleFactory().Make()
	saved.Code = req.Code
	saved.Label = req.Label
	saved.MemerlukanSTR = true
	saved.MemerlukanSIP = true
	saved.JenjangMin = req.JenjangMin
	saved.FHIRCode = req.FHIRCode
	s.mockJobTitleCreateSaves(1, saved)

	result, err := s.svc.CreateJobTitle(context.Background(), req, actor)

	s.NoError(err)
	s.Require().NotNil(result)
	s.Equal(req.Code, result.Code)
	s.Equal(req.Label, result.Label)
	s.True(result.MemerlukanSTR)
	s.True(result.MemerlukanSIP)
	s.Equal(req.JenjangMin, result.JenjangMin)
	s.Equal(req.FHIRCode, result.FHIRCode)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_CreateJobTitle_DefaultFlags() {
	req := newCreateJobTitleReq() // KategoriID: 1
	actor := superadminActor()

	// Pastikan di dalam method ini, mock untuk pengecekan KategoriID: 1 sudah benar
	s.mockJobTitleCreateChecks(req)

	saved := factories.NewJobTitleFactory().Make()
	saved.IsAktif = false
	saved.MemerlukanSTR = false
	saved.MemerlukanSIP = false
	s.mockJobTitleCreateSaves(1, saved)

	result, err := s.svc.CreateJobTitle(context.Background(), req, actor)

	// PERBAIKAN 1: Gunakan Require() agar test berhenti di sini jika error, mencegah panic
	s.Require().NoError(err)

	// PERBAIKAN 2: Karena Require() memastikan err == nil, result dijamin tidak nil
	s.False(result.IsAktif)
	s.False(result.MemerlukanSTR)
	s.False(result.MemerlukanSIP)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_CreateJobTitle_ExplicitIsAktifTrue() {
	req := newCreateJobTitleReq()
	req.IsAktif = true
	actor := superadminActor()

	s.mockJobTitleCreateChecks(req)
	saved := factories.NewJobTitleFactory().Make()
	saved.IsAktif = true
	s.mockJobTitleCreateSaves(1, saved)

	result, err := s.svc.CreateJobTitle(context.Background(), req, actor)

	s.NoError(err)
	s.True(result.IsAktif)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_CreateJobTitle_Forbidden() {
	req := newCreateJobTitleReq()
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.CreateJobTitle(context.Background(), req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
	s.repo.AssertNotCalled(s.T(), "CheckJobTitleKategori", mock.Anything)
	s.repo.AssertNotCalled(s.T(), "CreateJobTitle", mock.Anything)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_CreateJobTitle_KategoriInvalid() {
	req := newCreateJobTitleReq()
	actor := superadminActor()

	s.repo.On("CheckJobTitleKategori", mock.Anything, req.KategoriID).Return(false, nil)

	result, err := s.svc.CreateJobTitle(context.Background(), req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusUnprocessableEntity, appErr.Code)
	s.repo.AssertNotCalled(s.T(), "ExistsByCode", mock.Anything, mock.Anything)
	s.repo.AssertNotCalled(s.T(), "CreateJobTitle", mock.Anything)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_CreateJobTitle_KategoriCheckRepoError() {
	req := newCreateJobTitleReq()
	actor := superadminActor()

	s.repo.On("CheckJobTitleKategori", mock.Anything, req.KategoriID).Return(false, fmt.Errorf("db error"))

	result, err := s.svc.CreateJobTitle(context.Background(), req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_CreateJobTitle_RumpunProfesiInvalid() {
	req := newCreateJobTitleReq()
	req.RumpunProfesiID = ptrTo(int64(99))
	actor := superadminActor()

	s.repo.On("CheckJobTitleKategori", mock.Anything, req.KategoriID).Return(true, nil)
	s.repo.On("CheckJobTitleRumpunProfesi", mock.Anything, int64(99)).Return(false, nil)

	result, err := s.svc.CreateJobTitle(context.Background(), req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusUnprocessableEntity, appErr.Code)
	s.repo.AssertNotCalled(s.T(), "ExistsByCode", mock.Anything, mock.Anything)
	s.repo.AssertNotCalled(s.T(), "CreateJobTitle", mock.Anything)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_CreateJobTitle_DuplicateCode() {
	req := newCreateJobTitleReq()
	actor := superadminActor()

	s.repo.On("CheckJobTitleKategori", mock.Anything, req.KategoriID).Return(true, nil)
	s.repo.On("ExistsByCode", req.Code, int64(0)).Return(true, nil)

	result, err := s.svc.CreateJobTitle(context.Background(), req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusUnprocessableEntity, appErr.Code)
	s.repo.AssertNotCalled(s.T(), "CreateJobTitle", mock.Anything)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_CreateJobTitle_RepoError() {
	req := newCreateJobTitleReq()
	actor := superadminActor()

	s.mockJobTitleCreateChecks(req)
	s.repo.On("CreateJobTitle", mock.AnythingOfType("*models.JobTitle")).Return(fmt.Errorf("db error"))

	result, err := s.svc.CreateJobTitle(context.Background(), req, actor)

	s.Nil(result)
	s.Error(err)
}

// ── GetByID ──────────────────────────────────────────────────────────────────

func (s *KepegawaianJabatanServiceTestSuite) Test_GetJobTitleByID_Success() {
	actor := superadminActor()
	item := factories.NewJobTitleFactory().Make()
	item.ID = 1

	s.repo.On("GetJobTitleByID", int64(1)).Return(item, nil)

	result, err := s.svc.GetJobTitleByID(context.Background(), 1, actor)

	s.NoError(err)
	s.Require().NotNil(result)
	s.Equal(item.ID, result.ID)
	s.Equal(item.Code, result.Code)
	s.Equal(item.Label, result.Label)
	s.Equal(item.MemerlukanSTR, result.MemerlukanSTR)
	s.Equal(item.MemerlukanSIP, result.MemerlukanSIP)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_GetJobTitleByID_MapsRelations() {
	actor := superadminActor()
	item := factories.NewJobTitleFactory().Make()
	item.ID = 1
	item.Kategori = &models.JobTitleKategori{ID: item.KategoriID, Code: "MEDIS", Label: "Medis"}
	item.RumpunProfesiID = ptrTo(int64(2))
	item.RumpunProfesi = &models.JobTitleRumpunProfesi{ID: 2, Code: "DOKTER", Label: "Dokter"}

	s.repo.On("GetJobTitleByID", int64(1)).Return(item, nil)

	result, err := s.svc.GetJobTitleByID(context.Background(), 1, actor)

	s.NoError(err)
	s.Require().NotNil(result.Kategori)
	s.Equal("MEDIS", result.Kategori.Code)
	s.Require().NotNil(result.RumpunProfesi)
	s.Equal("DOKTER", result.RumpunProfesi.Code)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_GetJobTitleByID_NotFound() {
	actor := superadminActor()

	s.repo.On("GetJobTitleByID", int64(999)).Return(nil, nil)

	result, err := s.svc.GetJobTitleByID(context.Background(), 999, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianJabatanServiceTestSuite) Test_GetJobTitleByID_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.GetJobTitleByID(context.Background(), 1, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

// ── List ─────────────────────────────────────────────────────────────────────

func (s *KepegawaianJabatanServiceTestSuite) Test_ListJobTitle_Success() {
	actor := superadminActor()
	filter := &dto.FilterJobTitleRequest{}
	items := []models.JobTitle{
		*factories.NewJobTitleFactory().Make(),
		*factories.NewJobTitleFactory().Make(),
	}

	s.repo.On("ListJobTitle", 1, 10, filter).Return(items, int64(2), nil)

	result, total, err := s.svc.ListJobTitle(context.Background(), 1, 10, filter, actor)

	s.NoError(err)
	s.Equal(int64(2), total)
	s.Len(result, 2)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_ListJobTitle_Forbidden() {
	actor := regularActor()
	filter := &dto.FilterJobTitleRequest{}
	s.mockNoPermissions()

	result, total, err := s.svc.ListJobTitle(context.Background(), 1, 10, filter, actor)

	s.Nil(result)
	s.Equal(int64(0), total)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_ListJobTitle_RepoError() {
	actor := superadminActor()
	filter := &dto.FilterJobTitleRequest{}

	s.repo.On("ListJobTitle", 1, 10, filter).Return(nil, int64(0), fmt.Errorf("db error"))

	result, total, err := s.svc.ListJobTitle(context.Background(), 1, 10, filter, actor)

	s.Nil(result)
	s.Equal(int64(0), total)
	s.Error(err)
}

// ── ListSelect ───────────────────────────────────────────────────────────────
// Catatan: memakai cache (s.cache.Get/SetDefault); cache dinonaktifkan di
// SetupTest (cache.NewManager(nil, false, 0)) sehingga selalu miss dan jatuh
// ke s.repo.ListSelectJobTitle.

func (s *KepegawaianJabatanServiceTestSuite) Test_ListSelectJobTitle_Success() {
	actor := superadminActor()
	items := []models.JobTitle{*factories.NewJobTitleFactory().Make()}

	s.repo.On("ListSelectJobTitle", "").Return(items, nil)

	result, err := s.svc.ListSelectJobTitle(context.Background(), "", actor)

	s.NoError(err)
	s.Len(result, 1)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_ListSelectJobTitle_Empty_NotFound() {
	actor := superadminActor()

	s.repo.On("ListSelectJobTitle", "xxx").Return([]models.JobTitle{}, nil)

	result, err := s.svc.ListSelectJobTitle(context.Background(), "xxx", actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusNotFound, appErr.Code)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_ListSelectJobTitle_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.ListSelectJobTitle(context.Background(), "", actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

// ── Update ───────────────────────────────────────────────────────────────────

func (s *KepegawaianJabatanServiceTestSuite) Test_UpdateJobTitle_Success() {
	actor := superadminActor()
	existing := factories.NewJobTitleFactory().Make()
	existing.ID = 1
	newLabel := "Updated Label"
	req := &dto.UpdateJobTitleRequest{Label: &newLabel}

	// KategoriID/RumpunProfesiID/Code tidak dikirim → tidak ada pengecekan
	// kategori/rumpun/duplikat kode yang dipanggil.
	s.repo.On("GetJobTitleByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdateJobTitle", mock.AnythingOfType("*models.JobTitle")).Return(nil)

	result, err := s.svc.UpdateJobTitle(context.Background(), 1, req, actor)

	s.NoError(err)
	s.Equal(newLabel, result.Label)
	s.repo.AssertNotCalled(s.T(), "CheckJobTitleKategori", mock.Anything)
	s.repo.AssertNotCalled(s.T(), "ExistsByCode", mock.Anything, mock.Anything)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_UpdateJobTitle_ChangeKategori_Success() {
	actor := superadminActor()
	existing := factories.NewJobTitleFactory().Make()
	existing.ID = 1
	newKategoriID := int64(9)
	req := &dto.UpdateJobTitleRequest{KategoriID: &newKategoriID}

	existing.Kategori = &models.JobTitleKategori{ID: newKategoriID, Code: "KAT-9", Label: "Kategori 9"}

	// UBAH BARIS INI: Hapus mock.Anything
	s.repo.On("GetJobTitleByID", int64(1)).Return(existing, nil).Twice()
	s.repo.On("CheckJobTitleKategori", mock.Anything, newKategoriID).Return(true, nil)
	s.repo.On("UpdateJobTitle", mock.MatchedBy(func(m *models.JobTitle) bool {
		return m.KategoriID == newKategoriID
	})).Return(nil)

	result, err := s.svc.UpdateJobTitle(context.Background(), 1, req, actor)

	s.NoError(err)
	s.Require().NotNil(result.Kategori)
	s.Equal(newKategoriID, result.Kategori.ID)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_UpdateJobTitle_ChangeKategori_Invalid() {
	actor := superadminActor()
	existing := factories.NewJobTitleFactory().Make()
	existing.ID = 1
	newKategoriID := int64(99)
	req := &dto.UpdateJobTitleRequest{KategoriID: &newKategoriID}

	s.repo.On("GetJobTitleByID", int64(1)).Return(existing, nil)
	s.repo.On("CheckJobTitleKategori", mock.Anything, newKategoriID).Return(false, nil)

	result, err := s.svc.UpdateJobTitle(context.Background(), 1, req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusUnprocessableEntity, appErr.Code)
	s.repo.AssertNotCalled(s.T(), "UpdateJobTitle", mock.Anything)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_UpdateJobTitle_ChangeRumpunProfesi_Invalid() {
	actor := superadminActor()
	existing := factories.NewJobTitleFactory().Make()
	existing.ID = 1
	newRumpunID := int64(99)
	req := &dto.UpdateJobTitleRequest{RumpunProfesiID: &newRumpunID}

	s.repo.On("GetJobTitleByID", int64(1)).Return(existing, nil)
	s.repo.On("CheckJobTitleRumpunProfesi", mock.Anything, newRumpunID).Return(false, nil)

	result, err := s.svc.UpdateJobTitle(context.Background(), 1, req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusUnprocessableEntity, appErr.Code)
	s.repo.AssertNotCalled(s.T(), "UpdateJobTitle", mock.Anything)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_UpdateJobTitle_DuplicateCode() {
	actor := superadminActor()
	existing := factories.NewJobTitleFactory().Make()
	existing.ID = 1
	newCode := "SUDAH-DIPAKAI"
	req := &dto.UpdateJobTitleRequest{Code: &newCode}

	s.repo.On("GetJobTitleByID", int64(1)).Return(existing, nil)
	s.repo.On("ExistsByCode", newCode, int64(1)).Return(true, nil)

	result, err := s.svc.UpdateJobTitle(context.Background(), 1, req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusUnprocessableEntity, appErr.Code)
	s.repo.AssertNotCalled(s.T(), "UpdateJobTitle", mock.Anything)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_UpdateJobTitle_PartialFields() {
	actor := superadminActor()
	existing := factories.NewJobTitleFactory().Make()
	existing.ID = 1
	origCode, origKategori := existing.Code, existing.KategoriID
	req := &dto.UpdateJobTitleRequest{
		MemerlukanSTR: ptrTo(true),
		JenjangMin:    ptrTo("S1"),
	}

	s.repo.On("GetJobTitleByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdateJobTitle", mock.MatchedBy(func(m *models.JobTitle) bool {
		return m.Code == origCode &&
			m.KategoriID == origKategori &&
			m.MemerlukanSTR &&
			m.JenjangMin != nil && *m.JenjangMin == "S1"
	})).Return(nil)

	result, err := s.svc.UpdateJobTitle(context.Background(), 1, req, actor)

	s.NoError(err)
	s.Equal(origCode, result.Code)
	s.True(result.MemerlukanSTR)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_UpdateJobTitle_CanSetFlagToFalse() {
	actor := superadminActor()
	existing := factories.NewJobTitleFactory().Make()
	existing.ID = 1
	existing.MemerlukanSIP = true
	req := &dto.UpdateJobTitleRequest{MemerlukanSIP: ptrTo(false)}

	s.repo.On("GetJobTitleByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdateJobTitle", mock.MatchedBy(func(m *models.JobTitle) bool {
		return !m.MemerlukanSIP
	})).Return(nil)

	result, err := s.svc.UpdateJobTitle(context.Background(), 1, req, actor)

	s.NoError(err)
	s.False(result.MemerlukanSIP)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_UpdateJobTitle_NotFound() {
	actor := superadminActor()
	req := &dto.UpdateJobTitleRequest{}

	s.repo.On("GetJobTitleByID", int64(999)).Return(nil, nil)

	result, err := s.svc.UpdateJobTitle(context.Background(), 999, req, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianJabatanServiceTestSuite) Test_UpdateJobTitle_Forbidden() {
	actor := regularActor()
	req := &dto.UpdateJobTitleRequest{}
	s.mockNoPermissions()

	result, err := s.svc.UpdateJobTitle(context.Background(), 1, req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_UpdateJobTitle_RepoError() {
	actor := superadminActor()
	existing := factories.NewJobTitleFactory().Make()
	existing.ID = 1
	req := &dto.UpdateJobTitleRequest{}

	s.repo.On("GetJobTitleByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdateJobTitle", mock.AnythingOfType("*models.JobTitle")).Return(fmt.Errorf("db error"))

	result, err := s.svc.UpdateJobTitle(context.Background(), 1, req, actor)

	s.Nil(result)
	s.Error(err)
}

// ── Delete ───────────────────────────────────────────────────────────────────

func (s *KepegawaianJabatanServiceTestSuite) Test_DeleteJobTitle_Success() {
	actor := superadminActor()
	existing := factories.NewJobTitleFactory().Make()
	existing.ID = 1

	s.repo.On("GetJobTitleByID", int64(1)).Return(existing, nil)
	s.repo.On("DeleteJobTitle", int64(1), actor.UserID).Return(nil)

	err := s.svc.DeleteJobTitle(context.Background(), 1, actor)

	s.NoError(err)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_DeleteJobTitle_NotFound() {
	actor := superadminActor()

	s.repo.On("GetJobTitleByID", int64(999)).Return(nil, nil)

	err := s.svc.DeleteJobTitle(context.Background(), 999, actor)

	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianJabatanServiceTestSuite) Test_DeleteJobTitle_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	err := s.svc.DeleteJobTitle(context.Background(), 1, actor)

	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}
