package tests

import (
	"fmt"
	"neosim_go/internal/modules/master/master/dto"
	"neosim_go/internal/modules/master/master/models"
	"neosim_go/internal/modules/master/master/tests/factories"
	rbacModels "neosim_go/internal/modules/rbac/models"
	appErrors "neosim_go/internal/shared/errors"
	"net/http"

	"github.com/stretchr/testify/mock"
)

// =====================================================================
// Pekerjaan
// =====================================================================

func (s *MasterServiceTestSuite) Test_CreatePekerjaan_Superadmin_Success() {
	req := &dto.CreateMasterPekerjaanRequest{Name: "Dokter"}
	actor := superadminActor()

	// Mock cek duplikasi nama (return nil agar dianggap belum ada)
	s.repo.On("GetByNamePekerjaan", req.Name).Return(nil, nil)
	s.repo.On("CreatePekerjaan", mock.AnythingOfType("*models.MasterPekerjaan")).Return(nil)

	result, err := s.svc.CreatePekerjaan(req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(req.Name, result.Name)
	s.repo.AssertExpectations(s.T())
}

func (s *MasterServiceTestSuite) Test_CreatePekerjaan_WithPermission_Success() {
	req := &dto.CreateMasterPekerjaanRequest{Name: "Dokter"}
	actor := regularActor()

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermMasterCreate, mock.Anything).Return(true, nil)
	s.repo.On("GetByNamePekerjaan", req.Name).Return(nil, nil) // Added
	s.repo.On("CreatePekerjaan", mock.AnythingOfType("*models.MasterPekerjaan")).Return(nil)

	result, err := s.svc.CreatePekerjaan(req, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *MasterServiceTestSuite) Test_CreatePekerjaan_WithManagePermission_Success() {
	req := &dto.CreateMasterPekerjaanRequest{Name: "Dokter"}
	actor := regularActor()

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermMasterCreate, mock.Anything).Return(false, nil)
	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermMasterManage, mock.Anything).Return(true, nil)
	s.repo.On("GetByNamePekerjaan", req.Name).Return(nil, nil) // Added
	s.repo.On("CreatePekerjaan", mock.AnythingOfType("*models.MasterPekerjaan")).Return(nil)

	result, err := s.svc.CreatePekerjaan(req, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *MasterServiceTestSuite) Test_CreatePekerjaan_Forbidden() {
	req := &dto.CreateMasterPekerjaanRequest{Name: "Dokter"}
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.CreatePekerjaan(req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *MasterServiceTestSuite) Test_CreatePekerjaan_RepoError() {
	req := &dto.CreateMasterPekerjaanRequest{Name: "Dokter"}
	actor := superadminActor()

	s.repo.On("GetByNamePekerjaan", req.Name).Return(nil, nil) // Added: harus lolos cek duplikat dulu
	s.repo.On("CreatePekerjaan", mock.AnythingOfType("*models.MasterPekerjaan")).Return(fmt.Errorf("db error"))

	result, err := s.svc.CreatePekerjaan(req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *MasterServiceTestSuite) Test_GetByIDPekerjaan_Success() {
	item := factories.NewPekerjaanFactory().Make()
	item.ID = 1

	s.repo.On("GetByIDPekerjaan", int64(1)).Return(item, nil)

	// FIXED: Typo GetByIDPekerjaan -> GetByIDPekerjaan
	result, err := s.svc.GetByIDPekerjaan(1)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(item.ID, result.ID)
	s.Equal(item.Name, result.Name)
}

func (s *MasterServiceTestSuite) Test_GetByIDPekerjaan_NotFound() {
	s.repo.On("GetByIDPekerjaan", int64(999)).Return(nil, nil)

	result, err := s.svc.GetByIDPekerjaan(999)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *MasterServiceTestSuite) Test_GetByIDPekerjaan_RepoError() {
	s.repo.On("GetByIDPekerjaan", int64(1)).Return(nil, fmt.Errorf("db error"))

	result, err := s.svc.GetByIDPekerjaan(1)

	s.Nil(result)
	s.Error(err)
}

func (s *MasterServiceTestSuite) Test_ListPekerjaan_Success() {
	filter := &dto.FilterMasterPekerjaanRequest{}
	items := []models.MasterPekerjaan{
		*factories.NewPekerjaanFactory().Make(),
		*factories.NewPekerjaanFactory().Make(),
	}

	// FIXED: Typo ListMasterPekerjaan -> ListPekerjaan
	s.repo.On("ListPekerjaan", 1, 10, filter).Return(items, int64(2), nil)

	result, total, err := s.svc.ListPekerjaan(1, 10, filter)

	s.NoError(err)
	s.Equal(int64(2), total)
	s.Len(result, 2)
}

func (s *MasterServiceTestSuite) Test_ListPekerjaan_NotFound() {
	filter := &dto.FilterMasterPekerjaanRequest{}

	// FIXED: Typo ListMasterPekerjaan -> ListPekerjaan
	s.repo.On("ListPekerjaan", 1, 10, filter).Return([]models.MasterPekerjaan{}, int64(0), nil)

	result, total, err := s.svc.ListPekerjaan(1, 10, filter)

	s.Nil(result)
	s.Equal(int64(0), total)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *MasterServiceTestSuite) Test_ListPekerjaan_DefaultPagination() {
	filter := &dto.FilterMasterPekerjaanRequest{}
	items := []models.MasterPekerjaan{*factories.NewPekerjaanFactory().Make()}

	// FIXED: Typo ListMasterPekerjaan -> ListPekerjaan
	s.repo.On("ListPekerjaan", 1, 10, filter).Return(items, int64(1), nil)

	_, _, err := s.svc.ListPekerjaan(0, 0, filter)

	s.NoError(err)
	s.repo.AssertCalled(s.T(), "ListPekerjaan", 1, 10, filter)
}

func (s *MasterServiceTestSuite) Test_ListPekerjaan_PageSizeCapped() {
	filter := &dto.FilterMasterPekerjaanRequest{}
	items := []models.MasterPekerjaan{*factories.NewPekerjaanFactory().Make()}

	// FIXED: Typo ListMasterPekerjaan -> ListPekerjaan
	s.repo.On("ListPekerjaan", 1, 10, filter).Return(items, int64(1), nil)

	_, _, err := s.svc.ListPekerjaan(1, 999, filter)

	s.NoError(err)
	s.repo.AssertCalled(s.T(), "ListPekerjaan", 1, 10, filter)
}

func (s *MasterServiceTestSuite) Test_UpdatePekerjaan_Superadmin_Success() {
	actor := superadminActor()
	existing := factories.NewPekerjaanFactory().Make()
	existing.ID = 1
	newName := "Updated Name"
	req := &dto.UpdateMasterPekerjaanRequest{Name: &newName}

	s.repo.On("GetByIDPekerjaan", int64(1)).Return(existing, nil)
	// Added: Mock cek duplikasi nama untuk update
	s.repo.On("GetByNamePekerjaan", newName).Return(nil, nil)
	s.repo.On("UpdatePekerjaan", mock.AnythingOfType("*models.MasterPekerjaan")).Return(nil)

	result, err := s.svc.UpdatePekerjaan(1, req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(newName, result.Name)
}

func (s *MasterServiceTestSuite) Test_UpdatePekerjaan_Forbidden() {
	actor := regularActor()
	req := &dto.UpdateMasterPekerjaanRequest{}
	s.mockNoPermissions()

	result, err := s.svc.UpdatePekerjaan(1, req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *MasterServiceTestSuite) Test_UpdatePekerjaan_NotFound() {
	actor := superadminActor()
	req := &dto.UpdateMasterPekerjaanRequest{}

	s.repo.On("GetByIDPekerjaan", int64(999)).Return(nil, nil)

	result, err := s.svc.UpdatePekerjaan(999, req, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *MasterServiceTestSuite) Test_UpdatePekerjaan_PartialFields() {
	actor := superadminActor()
	existing := factories.NewPekerjaanFactory().Make()
	existing.ID = 1
	newName := "New " + existing.Name
	req := &dto.UpdateMasterPekerjaanRequest{Name: &newName}

	s.repo.On("GetByIDPekerjaan", int64(1)).Return(existing, nil)
	// Added: Mock cek duplikasi nama untuk update
	s.repo.On("GetByNamePekerjaan", newName).Return(nil, nil)
	s.repo.On("UpdatePekerjaan", mock.MatchedBy(func(m *models.MasterPekerjaan) bool {
		return m.Name == newName
	})).Return(nil)

	result, err := s.svc.UpdatePekerjaan(1, req, actor)

	s.NoError(err)
	s.Equal(newName, result.Name)
}

func (s *MasterServiceTestSuite) Test_UpdatePekerjaan_RepoError() {
	actor := superadminActor()
	existing := factories.NewPekerjaanFactory().Make()
	existing.ID = 1
	req := &dto.UpdateMasterPekerjaanRequest{}

	s.repo.On("GetByIDPekerjaan", int64(1)).Return(existing, nil)
	s.repo.On("UpdatePekerjaan", mock.AnythingOfType("*models.MasterPekerjaan")).Return(fmt.Errorf("db error"))

	result, err := s.svc.UpdatePekerjaan(1, req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *MasterServiceTestSuite) Test_DeletePekerjaan_Superadmin_Success() {
	actor := superadminActor()
	existing := factories.NewPekerjaanFactory().Make()
	existing.ID = 1

	s.repo.On("GetByIDPekerjaan", int64(1)).Return(existing, nil)
	s.repo.On("DeletePekerjaan", int64(1)).Return(nil)

	err := s.svc.DeletePekerjaan(1, actor)

	s.NoError(err)
	s.repo.AssertExpectations(s.T())
}

func (s *MasterServiceTestSuite) Test_DeletePekerjaan_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	err := s.svc.DeletePekerjaan(1, actor)

	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *MasterServiceTestSuite) Test_DeletePekerjaan_NotFound() {
	actor := superadminActor()

	s.repo.On("GetByIDPekerjaan", int64(999)).Return(nil, nil)

	err := s.svc.DeletePekerjaan(999, actor)

	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *MasterServiceTestSuite) Test_DeletePekerjaan_RepoError() {
	actor := superadminActor()
	existing := factories.NewPekerjaanFactory().Make()
	existing.ID = 1

	s.repo.On("GetByIDPekerjaan", int64(1)).Return(existing, nil)
	s.repo.On("DeletePekerjaan", int64(1)).Return(fmt.Errorf("db error"))

	err := s.svc.DeletePekerjaan(1, actor)

	s.Error(err)
}
