package tests

import (
	"context"
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
// Pendidikan
// =====================================================================

func (s *MasterServiceTestSuite) Test_CreatePendidikan_Superadmin_Success() {
	req := &dto.CreateMasterPendidikanRequest{Name: "Dokter"}
	actor := superadminActor()

	// Mock cek duplikasi nama (return nil agar dianggap belum ada)
	s.repo.On("GetByNamePendidikan", req.Name).Return(nil, nil)
	s.repo.On("CreatePendidikan", mock.AnythingOfType("*models.MasterPendidikan")).Return(nil)

	result, err := s.svc.CreatePendidikan(context.Background(), req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(req.Name, result.Name)
	s.repo.AssertExpectations(s.T())
}

func (s *MasterServiceTestSuite) Test_CreatePendidikan_WithPermission_Success() {
	req := &dto.CreateMasterPendidikanRequest{Name: "Dokter"}
	actor := regularActor()

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermMasterCreate, mock.Anything).Return(true, nil)
	s.repo.On("GetByNamePendidikan", req.Name).Return(nil, nil) // Added
	s.repo.On("CreatePendidikan", mock.AnythingOfType("*models.MasterPendidikan")).Return(nil)

	result, err := s.svc.CreatePendidikan(context.Background(), req, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *MasterServiceTestSuite) Test_CreatePendidikan_WithManagePermission_Success() {
	req := &dto.CreateMasterPendidikanRequest{Name: "Dokter"}
	actor := regularActor()

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermMasterCreate, mock.Anything).Return(false, nil)
	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermMasterManage, mock.Anything).Return(true, nil)
	s.repo.On("GetByNamePendidikan", req.Name).Return(nil, nil) // Added
	s.repo.On("CreatePendidikan", mock.AnythingOfType("*models.MasterPendidikan")).Return(nil)

	result, err := s.svc.CreatePendidikan(context.Background(), req, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *MasterServiceTestSuite) Test_CreatePendidikan_Forbidden() {
	req := &dto.CreateMasterPendidikanRequest{Name: "Dokter"}
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.CreatePendidikan(context.Background(), req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *MasterServiceTestSuite) Test_CreatePendidikan_RepoError() {
	req := &dto.CreateMasterPendidikanRequest{Name: "Dokter"}
	actor := superadminActor()

	s.repo.On("GetByNamePendidikan", req.Name).Return(nil, nil) // Added: harus lolos cek duplikat dulu
	s.repo.On("CreatePendidikan", mock.AnythingOfType("*models.MasterPendidikan")).Return(fmt.Errorf("db error"))

	result, err := s.svc.CreatePendidikan(context.Background(), req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *MasterServiceTestSuite) Test_GetByIDPendidikan_Success() {
	item := factories.NewPendidikanFactory().Make()
	item.ID = 1

	s.repo.On("GetByIDPendidikan", int64(1)).Return(item, nil)

	// FIXED: Typo GetByIDPekerjaan -> GetByIDPendidikan
	result, err := s.svc.GetByIDPendidikan(context.Background(), 1)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(item.ID, result.ID)
	s.Equal(item.Name, result.Name)
}

func (s *MasterServiceTestSuite) Test_GetByIDPendidikan_NotFound() {
	s.repo.On("GetByIDPendidikan", int64(999)).Return(nil, nil)

	result, err := s.svc.GetByIDPendidikan(context.Background(), 999)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *MasterServiceTestSuite) Test_GetByIDPendidikan_RepoError() {
	s.repo.On("GetByIDPendidikan", int64(1)).Return(nil, fmt.Errorf("db error"))

	result, err := s.svc.GetByIDPendidikan(context.Background(), 1)

	s.Nil(result)
	s.Error(err)
}

func (s *MasterServiceTestSuite) Test_ListPendidikan_Success() {
	filter := &dto.FilterMasterPendidikanRequest{}
	items := []models.MasterPendidikan{
		*factories.NewPendidikanFactory().Make(),
		*factories.NewPendidikanFactory().Make(),
	}

	// FIXED: Typo ListMasterPendidikan -> ListPendidikan
	s.repo.On("ListPendidikan", 1, 10, filter).Return(items, int64(2), nil)

	result, total, err := s.svc.ListPendidikan(context.Background(), 1, 10, filter)

	s.NoError(err)
	s.Equal(int64(2), total)
	s.Len(result, 2)
}

func (s *MasterServiceTestSuite) Test_ListPendidikan_NotFound() {
	filter := &dto.FilterMasterPendidikanRequest{}

	// FIXED: Typo ListMasterPendidikan -> ListPendidikan
	s.repo.On("ListPendidikan", 1, 10, filter).Return([]models.MasterPendidikan{}, int64(0), nil)

	result, total, err := s.svc.ListPendidikan(context.Background(), 1, 10, filter)

	s.Nil(result)
	s.Equal(int64(0), total)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *MasterServiceTestSuite) Test_ListPendidikan_DefaultPagination() {
	filter := &dto.FilterMasterPendidikanRequest{}
	items := []models.MasterPendidikan{*factories.NewPendidikanFactory().Make()}

	// FIXED: Typo ListMasterPendidikan -> ListPendidikan
	s.repo.On("ListPendidikan", 1, 10, filter).Return(items, int64(1), nil)

	_, _, err := s.svc.ListPendidikan(context.Background(), 0, 0, filter)

	s.NoError(err)
	s.repo.AssertCalled(s.T(), "ListPendidikan", 1, 10, filter)
}

func (s *MasterServiceTestSuite) Test_ListPendidikan_PageSizeCapped() {
	filter := &dto.FilterMasterPendidikanRequest{}
	items := []models.MasterPendidikan{*factories.NewPendidikanFactory().Make()}

	// FIXED: Typo ListMasterPendidikan -> ListPendidikan
	s.repo.On("ListPendidikan", 1, 10, filter).Return(items, int64(1), nil)

	_, _, err := s.svc.ListPendidikan(context.Background(), 1, 999, filter)

	s.NoError(err)
	s.repo.AssertCalled(s.T(), "ListPendidikan", 1, 10, filter)
}

func (s *MasterServiceTestSuite) Test_UpdatePendidikan_Superadmin_Success() {
	actor := superadminActor()
	existing := factories.NewPendidikanFactory().Make()
	existing.ID = 1
	newName := "Updated Name"
	req := &dto.UpdateMasterPendidikanRequest{Name: &newName}

	s.repo.On("GetByIDPendidikan", int64(1)).Return(existing, nil)
	// Added: Mock cek duplikasi nama untuk update
	s.repo.On("GetByNamePendidikan", newName).Return(nil, nil)
	s.repo.On("UpdatePendidikan", mock.AnythingOfType("*models.MasterPendidikan")).Return(nil)

	result, err := s.svc.UpdatePendidikan(context.Background(), 1, req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(newName, result.Name)
}

func (s *MasterServiceTestSuite) Test_UpdatePendidikan_Forbidden() {
	actor := regularActor()
	req := &dto.UpdateMasterPendidikanRequest{}
	s.mockNoPermissions()

	result, err := s.svc.UpdatePendidikan(context.Background(), 1, req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *MasterServiceTestSuite) Test_UpdatePendidikan_NotFound() {
	actor := superadminActor()
	req := &dto.UpdateMasterPendidikanRequest{}

	s.repo.On("GetByIDPendidikan", int64(999)).Return(nil, nil)

	result, err := s.svc.UpdatePendidikan(context.Background(), 999, req, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *MasterServiceTestSuite) Test_UpdatePendidikan_PartialFields() {
	actor := superadminActor()
	existing := factories.NewPendidikanFactory().Make()
	existing.ID = 1
	newName := "New " + existing.Name
	req := &dto.UpdateMasterPendidikanRequest{Name: &newName}

	s.repo.On("GetByIDPendidikan", int64(1)).Return(existing, nil)
	// Added: Mock cek duplikasi nama untuk update
	s.repo.On("GetByNamePendidikan", newName).Return(nil, nil)
	s.repo.On("UpdatePendidikan", mock.MatchedBy(func(m *models.MasterPendidikan) bool {
		return m.Name == newName
	})).Return(nil)

	result, err := s.svc.UpdatePendidikan(context.Background(), 1, req, actor)

	s.NoError(err)
	s.Equal(newName, result.Name)
}

func (s *MasterServiceTestSuite) Test_UpdatePendidikan_RepoError() {
	actor := superadminActor()
	existing := factories.NewPendidikanFactory().Make()
	existing.ID = 1
	req := &dto.UpdateMasterPendidikanRequest{}

	s.repo.On("GetByIDPendidikan", int64(1)).Return(existing, nil)
	s.repo.On("UpdatePendidikan", mock.AnythingOfType("*models.MasterPendidikan")).Return(fmt.Errorf("db error"))

	result, err := s.svc.UpdatePendidikan(context.Background(), 1, req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *MasterServiceTestSuite) Test_DeletePendidikan_Superadmin_Success() {
	actor := superadminActor()
	existing := factories.NewPendidikanFactory().Make()
	existing.ID = 1

	s.repo.On("GetByIDPendidikan", int64(1)).Return(existing, nil)
	s.repo.On("DeletePendidikan", int64(1)).Return(nil)

	err := s.svc.DeletePendidikan(context.Background(), 1, actor)

	s.NoError(err)
	s.repo.AssertExpectations(s.T())
}

func (s *MasterServiceTestSuite) Test_DeletePendidikan_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	err := s.svc.DeletePendidikan(context.Background(), 1, actor)

	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *MasterServiceTestSuite) Test_DeletePendidikan_NotFound() {
	actor := superadminActor()

	s.repo.On("GetByIDPendidikan", int64(999)).Return(nil, nil)

	err := s.svc.DeletePendidikan(context.Background(), 999, actor)

	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *MasterServiceTestSuite) Test_DeletePendidikan_RepoError() {
	actor := superadminActor()
	existing := factories.NewPendidikanFactory().Make()
	existing.ID = 1

	s.repo.On("GetByIDPendidikan", int64(1)).Return(existing, nil)
	s.repo.On("DeletePendidikan", int64(1)).Return(fmt.Errorf("db error"))

	err := s.svc.DeletePendidikan(context.Background(), 1, actor)

	s.Error(err)
}
