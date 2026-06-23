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
// GolonganDarah
// =====================================================================

func (s *MasterServiceTestSuite) Test_CreateGolonganDarah_Superadmin_Success() {
	req := &dto.CreateMasterGolonganDarahRequest{Name: "Dokter"}
	actor := superadminActor()

	// Mock cek duplikasi nama (return nil agar dianggap belum ada)
	s.repo.On("GetByNameGolonganDarah", req.Name).Return(nil, nil)
	s.repo.On("CreateGolonganDarah", mock.AnythingOfType("*models.MasterGolonganDarah")).Return(nil)

	result, err := s.svc.CreateGolonganDarah(req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(req.Name, result.Name)
	s.repo.AssertExpectations(s.T())
}

func (s *MasterServiceTestSuite) Test_CreateGolonganDarah_WithPermission_Success() {
	req := &dto.CreateMasterGolonganDarahRequest{Name: "Dokter"}
	actor := regularActor()

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermMasterCreate, mock.Anything).Return(true, nil)
	s.repo.On("GetByNameGolonganDarah", req.Name).Return(nil, nil) // Added
	s.repo.On("CreateGolonganDarah", mock.AnythingOfType("*models.MasterGolonganDarah")).Return(nil)

	result, err := s.svc.CreateGolonganDarah(req, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *MasterServiceTestSuite) Test_CreateGolonganDarah_WithManagePermission_Success() {
	req := &dto.CreateMasterGolonganDarahRequest{Name: "Dokter"}
	actor := regularActor()

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermMasterCreate, mock.Anything).Return(false, nil)
	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermMasterManage, mock.Anything).Return(true, nil)
	s.repo.On("GetByNameGolonganDarah", req.Name).Return(nil, nil) // Added
	s.repo.On("CreateGolonganDarah", mock.AnythingOfType("*models.MasterGolonganDarah")).Return(nil)

	result, err := s.svc.CreateGolonganDarah(req, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *MasterServiceTestSuite) Test_CreateGolonganDarah_Forbidden() {
	req := &dto.CreateMasterGolonganDarahRequest{Name: "Dokter"}
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.CreateGolonganDarah(req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *MasterServiceTestSuite) Test_CreateGolonganDarah_RepoError() {
	req := &dto.CreateMasterGolonganDarahRequest{Name: "Dokter"}
	actor := superadminActor()

	s.repo.On("GetByNameGolonganDarah", req.Name).Return(nil, nil) // Added: harus lolos cek duplikat dulu
	s.repo.On("CreateGolonganDarah", mock.AnythingOfType("*models.MasterGolonganDarah")).Return(fmt.Errorf("db error"))

	result, err := s.svc.CreateGolonganDarah(req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *MasterServiceTestSuite) Test_GetByIDGolonganDarah_Success() {
	item := factories.NewGolonganDarahFactory().Make()
	item.ID = 1

	s.repo.On("GetByIDGolonganDarah", int64(1)).Return(item, nil)

	// FIXED: Typo GetByIDGolonganDarah -> GetByIDGolonganDarah
	result, err := s.svc.GetByIDGolonganDarah(1)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(item.ID, result.ID)
	s.Equal(item.Name, result.Name)
}

func (s *MasterServiceTestSuite) Test_GetByIDGolonganDarah_NotFound() {
	s.repo.On("GetByIDGolonganDarah", int64(999)).Return(nil, nil)

	result, err := s.svc.GetByIDGolonganDarah(999)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *MasterServiceTestSuite) Test_GetByIDGolonganDarah_RepoError() {
	s.repo.On("GetByIDGolonganDarah", int64(1)).Return(nil, fmt.Errorf("db error"))

	result, err := s.svc.GetByIDGolonganDarah(1)

	s.Nil(result)
	s.Error(err)
}

func (s *MasterServiceTestSuite) Test_ListGolonganDarah_Success() {
	filter := &dto.FilterMasterGolonganDarahRequest{}
	items := []models.MasterGolonganDarah{
		*factories.NewGolonganDarahFactory().Make(),
		*factories.NewGolonganDarahFactory().Make(),
	}

	// FIXED: Typo ListMasterGolonganDarah -> ListGolonganDarah
	s.repo.On("ListGolonganDarah", 1, 10, filter).Return(items, int64(2), nil)

	result, total, err := s.svc.ListGolonganDarah(1, 10, filter)

	s.NoError(err)
	s.Equal(int64(2), total)
	s.Len(result, 2)
}

func (s *MasterServiceTestSuite) Test_ListGolonganDarah_NotFound() {
	filter := &dto.FilterMasterGolonganDarahRequest{}

	// FIXED: Typo ListMasterGolonganDarah -> ListGolonganDarah
	s.repo.On("ListGolonganDarah", 1, 10, filter).Return([]models.MasterGolonganDarah{}, int64(0), nil)

	result, total, err := s.svc.ListGolonganDarah(1, 10, filter)

	s.Nil(result)
	s.Equal(int64(0), total)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *MasterServiceTestSuite) Test_ListGolonganDarah_DefaultPagination() {
	filter := &dto.FilterMasterGolonganDarahRequest{}
	items := []models.MasterGolonganDarah{*factories.NewGolonganDarahFactory().Make()}

	// FIXED: Typo ListMasterGolonganDarah -> ListGolonganDarah
	s.repo.On("ListGolonganDarah", 1, 10, filter).Return(items, int64(1), nil)

	_, _, err := s.svc.ListGolonganDarah(0, 0, filter)

	s.NoError(err)
	s.repo.AssertCalled(s.T(), "ListGolonganDarah", 1, 10, filter)
}

func (s *MasterServiceTestSuite) Test_ListGolonganDarah_PageSizeCapped() {
	filter := &dto.FilterMasterGolonganDarahRequest{}
	items := []models.MasterGolonganDarah{*factories.NewGolonganDarahFactory().Make()}

	// FIXED: Typo ListMasterGolonganDarah -> ListGolonganDarah
	s.repo.On("ListGolonganDarah", 1, 10, filter).Return(items, int64(1), nil)

	_, _, err := s.svc.ListGolonganDarah(1, 999, filter)

	s.NoError(err)
	s.repo.AssertCalled(s.T(), "ListGolonganDarah", 1, 10, filter)
}

func (s *MasterServiceTestSuite) Test_UpdateGolonganDarah_Superadmin_Success() {
	actor := superadminActor()
	existing := factories.NewGolonganDarahFactory().Make()
	existing.ID = 1
	newName := "Updated Name"
	req := &dto.UpdateMasterGolonganDarahRequest{Name: &newName}

	s.repo.On("GetByIDGolonganDarah", int64(1)).Return(existing, nil)
	// Added: Mock cek duplikasi nama untuk update
	s.repo.On("GetByNameGolonganDarah", newName).Return(nil, nil)
	s.repo.On("UpdateGolonganDarah", mock.AnythingOfType("*models.MasterGolonganDarah")).Return(nil)

	result, err := s.svc.UpdateGolonganDarah(1, req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(newName, result.Name)
}

func (s *MasterServiceTestSuite) Test_UpdateGolonganDarah_Forbidden() {
	actor := regularActor()
	req := &dto.UpdateMasterGolonganDarahRequest{}
	s.mockNoPermissions()

	result, err := s.svc.UpdateGolonganDarah(1, req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *MasterServiceTestSuite) Test_UpdateGolonganDarah_NotFound() {
	actor := superadminActor()
	req := &dto.UpdateMasterGolonganDarahRequest{}

	s.repo.On("GetByIDGolonganDarah", int64(999)).Return(nil, nil)

	result, err := s.svc.UpdateGolonganDarah(999, req, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *MasterServiceTestSuite) Test_UpdateGolonganDarah_PartialFields() {
	actor := superadminActor()
	existing := factories.NewGolonganDarahFactory().Make()
	existing.ID = 1
	newName := "New " + existing.Name
	req := &dto.UpdateMasterGolonganDarahRequest{Name: &newName}

	s.repo.On("GetByIDGolonganDarah", int64(1)).Return(existing, nil)
	// Added: Mock cek duplikasi nama untuk update
	s.repo.On("GetByNameGolonganDarah", newName).Return(nil, nil)
	s.repo.On("UpdateGolonganDarah", mock.MatchedBy(func(m *models.MasterGolonganDarah) bool {
		return m.Name == newName
	})).Return(nil)

	result, err := s.svc.UpdateGolonganDarah(1, req, actor)

	s.NoError(err)
	s.Equal(newName, result.Name)
}

func (s *MasterServiceTestSuite) Test_UpdateGolonganDarah_RepoError() {
	actor := superadminActor()
	existing := factories.NewGolonganDarahFactory().Make()
	existing.ID = 1
	req := &dto.UpdateMasterGolonganDarahRequest{}

	s.repo.On("GetByIDGolonganDarah", int64(1)).Return(existing, nil)
	s.repo.On("UpdateGolonganDarah", mock.AnythingOfType("*models.MasterGolonganDarah")).Return(fmt.Errorf("db error"))

	result, err := s.svc.UpdateGolonganDarah(1, req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *MasterServiceTestSuite) Test_DeleteGolonganDarah_Superadmin_Success() {
	actor := superadminActor()
	existing := factories.NewGolonganDarahFactory().Make()
	existing.ID = 1

	s.repo.On("GetByIDGolonganDarah", int64(1)).Return(existing, nil)
	s.repo.On("DeleteGolonganDarah", int64(1)).Return(nil)

	err := s.svc.DeleteGolonganDarah(1, actor)

	s.NoError(err)
	s.repo.AssertExpectations(s.T())
}

func (s *MasterServiceTestSuite) Test_DeleteGolonganDarah_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	err := s.svc.DeleteGolonganDarah(1, actor)

	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *MasterServiceTestSuite) Test_DeleteGolonganDarah_NotFound() {
	actor := superadminActor()

	s.repo.On("GetByIDGolonganDarah", int64(999)).Return(nil, nil)

	err := s.svc.DeleteGolonganDarah(999, actor)

	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *MasterServiceTestSuite) Test_DeleteGolonganDarah_RepoError() {
	actor := superadminActor()
	existing := factories.NewGolonganDarahFactory().Make()
	existing.ID = 1

	s.repo.On("GetByIDGolonganDarah", int64(1)).Return(existing, nil)
	s.repo.On("DeleteGolonganDarah", int64(1)).Return(fmt.Errorf("db error"))

	err := s.svc.DeleteGolonganDarah(1, actor)

	s.Error(err)
}
