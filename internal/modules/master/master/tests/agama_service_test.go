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
// Agama
// =====================================================================

func (s *MasterServiceTestSuite) Test_CreateAgama_Superadmin_Success() {
	req := &dto.CreateMasterAgamaRequest{Name: "Dokter"}
	actor := superadminActor()

	// Mock cek duplikasi nama (return nil agar dianggap belum ada)
	s.repo.On("GetByNameAgama", req.Name).Return(nil, nil)
	s.repo.On("CreateAgama", mock.AnythingOfType("*models.MasterAgama")).Return(nil)

	result, err := s.svc.CreateAgama(req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(req.Name, result.Name)
	s.repo.AssertExpectations(s.T())
}

func (s *MasterServiceTestSuite) Test_CreateAgama_WithPermission_Success() {
	req := &dto.CreateMasterAgamaRequest{Name: "Dokter"}
	actor := regularActor()

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermMasterCreate, mock.Anything).Return(true, nil)
	s.repo.On("GetByNameAgama", req.Name).Return(nil, nil) // Added
	s.repo.On("CreateAgama", mock.AnythingOfType("*models.MasterAgama")).Return(nil)

	result, err := s.svc.CreateAgama(req, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *MasterServiceTestSuite) Test_CreateAgama_WithManagePermission_Success() {
	req := &dto.CreateMasterAgamaRequest{Name: "Dokter"}
	actor := regularActor()

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermMasterCreate, mock.Anything).Return(false, nil)
	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermMasterManage, mock.Anything).Return(true, nil)
	s.repo.On("GetByNameAgama", req.Name).Return(nil, nil) // Added
	s.repo.On("CreateAgama", mock.AnythingOfType("*models.MasterAgama")).Return(nil)

	result, err := s.svc.CreateAgama(req, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *MasterServiceTestSuite) Test_CreateAgama_Forbidden() {
	req := &dto.CreateMasterAgamaRequest{Name: "Dokter"}
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.CreateAgama(req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *MasterServiceTestSuite) Test_CreateAgama_RepoError() {
	req := &dto.CreateMasterAgamaRequest{Name: "Dokter"}
	actor := superadminActor()

	s.repo.On("GetByNameAgama", req.Name).Return(nil, nil) // Added: harus lolos cek duplikat dulu
	s.repo.On("CreateAgama", mock.AnythingOfType("*models.MasterAgama")).Return(fmt.Errorf("db error"))

	result, err := s.svc.CreateAgama(req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *MasterServiceTestSuite) Test_GetByIDAgama_Success() {
	item := factories.NewAgamaFactory().Make()
	item.ID = 1

	s.repo.On("GetByIDAgama", int64(1)).Return(item, nil)

	// FIXED: Typo GetByIDAgama -> GetByIDAgama
	result, err := s.svc.GetByIDAgama(1)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(item.ID, result.ID)
	s.Equal(item.Name, result.Name)
}

func (s *MasterServiceTestSuite) Test_GetByIDAgama_NotFound() {
	s.repo.On("GetByIDAgama", int64(999)).Return(nil, nil)

	result, err := s.svc.GetByIDAgama(999)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *MasterServiceTestSuite) Test_GetByIDAgama_RepoError() {
	s.repo.On("GetByIDAgama", int64(1)).Return(nil, fmt.Errorf("db error"))

	result, err := s.svc.GetByIDAgama(1)

	s.Nil(result)
	s.Error(err)
}

func (s *MasterServiceTestSuite) Test_ListAgama_Success() {
	filter := &dto.FilterMasterAgamaRequest{}
	items := []models.MasterAgama{
		*factories.NewAgamaFactory().Make(),
		*factories.NewAgamaFactory().Make(),
	}

	// FIXED: Typo ListMasterAgama -> ListAgama
	s.repo.On("ListAgama", 1, 10, filter).Return(items, int64(2), nil)

	result, total, err := s.svc.ListAgama(1, 10, filter)

	s.NoError(err)
	s.Equal(int64(2), total)
	s.Len(result, 2)
}

func (s *MasterServiceTestSuite) Test_ListAgama_NotFound() {
	filter := &dto.FilterMasterAgamaRequest{}

	// FIXED: Typo ListMasterAgama -> ListAgama
	s.repo.On("ListAgama", 1, 10, filter).Return([]models.MasterAgama{}, int64(0), nil)

	result, total, err := s.svc.ListAgama(1, 10, filter)

	s.Nil(result)
	s.Equal(int64(0), total)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *MasterServiceTestSuite) Test_ListAgama_DefaultPagination() {
	filter := &dto.FilterMasterAgamaRequest{}
	items := []models.MasterAgama{*factories.NewAgamaFactory().Make()}

	// FIXED: Typo ListMasterAgama -> ListAgama
	s.repo.On("ListAgama", 1, 10, filter).Return(items, int64(1), nil)

	_, _, err := s.svc.ListAgama(0, 0, filter)

	s.NoError(err)
	s.repo.AssertCalled(s.T(), "ListAgama", 1, 10, filter)
}

func (s *MasterServiceTestSuite) Test_ListAgama_PageSizeCapped() {
	filter := &dto.FilterMasterAgamaRequest{}
	items := []models.MasterAgama{*factories.NewAgamaFactory().Make()}

	// FIXED: Typo ListMasterAgama -> ListAgama
	s.repo.On("ListAgama", 1, 10, filter).Return(items, int64(1), nil)

	_, _, err := s.svc.ListAgama(1, 999, filter)

	s.NoError(err)
	s.repo.AssertCalled(s.T(), "ListAgama", 1, 10, filter)
}

func (s *MasterServiceTestSuite) Test_UpdateAgama_Superadmin_Success() {
	actor := superadminActor()
	existing := factories.NewAgamaFactory().Make()
	existing.ID = 1
	newName := "Updated Name"
	req := &dto.UpdateMasterAgamaRequest{Name: &newName}

	s.repo.On("GetByIDAgama", int64(1)).Return(existing, nil)
	// Added: Mock cek duplikasi nama untuk update
	s.repo.On("GetByNameAgama", newName).Return(nil, nil)
	s.repo.On("UpdateAgama", mock.AnythingOfType("*models.MasterAgama")).Return(nil)

	result, err := s.svc.UpdateAgama(1, req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(newName, result.Name)
}

func (s *MasterServiceTestSuite) Test_UpdateAgama_Forbidden() {
	actor := regularActor()
	req := &dto.UpdateMasterAgamaRequest{}
	s.mockNoPermissions()

	result, err := s.svc.UpdateAgama(1, req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *MasterServiceTestSuite) Test_UpdateAgama_NotFound() {
	actor := superadminActor()
	req := &dto.UpdateMasterAgamaRequest{}

	s.repo.On("GetByIDAgama", int64(999)).Return(nil, nil)

	result, err := s.svc.UpdateAgama(999, req, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *MasterServiceTestSuite) Test_UpdateAgama_PartialFields() {
	actor := superadminActor()
	existing := factories.NewAgamaFactory().Make()
	existing.ID = 1
	newName := "New " + existing.Name
	req := &dto.UpdateMasterAgamaRequest{Name: &newName}

	s.repo.On("GetByIDAgama", int64(1)).Return(existing, nil)
	// Added: Mock cek duplikasi nama untuk update
	s.repo.On("GetByNameAgama", newName).Return(nil, nil)
	s.repo.On("UpdateAgama", mock.MatchedBy(func(m *models.MasterAgama) bool {
		return m.Name == newName
	})).Return(nil)

	result, err := s.svc.UpdateAgama(1, req, actor)

	s.NoError(err)
	s.Equal(newName, result.Name)
}

func (s *MasterServiceTestSuite) Test_UpdateAgama_RepoError() {
	actor := superadminActor()
	existing := factories.NewAgamaFactory().Make()
	existing.ID = 1
	req := &dto.UpdateMasterAgamaRequest{}

	s.repo.On("GetByIDAgama", int64(1)).Return(existing, nil)
	s.repo.On("UpdateAgama", mock.AnythingOfType("*models.MasterAgama")).Return(fmt.Errorf("db error"))

	result, err := s.svc.UpdateAgama(1, req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *MasterServiceTestSuite) Test_DeleteAgama_Superadmin_Success() {
	actor := superadminActor()
	existing := factories.NewAgamaFactory().Make()
	existing.ID = 1

	s.repo.On("GetByIDAgama", int64(1)).Return(existing, nil)
	s.repo.On("DeleteAgama", int64(1)).Return(nil)

	err := s.svc.DeleteAgama(1, actor)

	s.NoError(err)
	s.repo.AssertExpectations(s.T())
}

func (s *MasterServiceTestSuite) Test_DeleteAgama_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	err := s.svc.DeleteAgama(1, actor)

	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *MasterServiceTestSuite) Test_DeleteAgama_NotFound() {
	actor := superadminActor()

	s.repo.On("GetByIDAgama", int64(999)).Return(nil, nil)

	err := s.svc.DeleteAgama(999, actor)

	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *MasterServiceTestSuite) Test_DeleteAgama_RepoError() {
	actor := superadminActor()
	existing := factories.NewAgamaFactory().Make()
	existing.ID = 1

	s.repo.On("GetByIDAgama", int64(1)).Return(existing, nil)
	s.repo.On("DeleteAgama", int64(1)).Return(fmt.Errorf("db error"))

	err := s.svc.DeleteAgama(1, actor)

	s.Error(err)
}
