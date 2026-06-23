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
// JenisKelamin
// =====================================================================

func (s *MasterServiceTestSuite) Test_CreateJenisKelamin_Superadmin_Success() {
	req := &dto.CreateMasterJenisKelaminRequest{Name: "Dokter"}
	actor := superadminActor()

	// Mock cek duplikasi nama (return nil agar dianggap belum ada)
	s.repo.On("GetByNameJenisKelamin", req.Name).Return(nil, nil)
	s.repo.On("CreateJenisKelamin", mock.AnythingOfType("*models.MasterJenisKelamin")).Return(nil)

	result, err := s.svc.CreateJenisKelamin(req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(req.Name, result.Name)
	s.repo.AssertExpectations(s.T())
}

func (s *MasterServiceTestSuite) Test_CreateJenisKelamin_WithPermission_Success() {
	req := &dto.CreateMasterJenisKelaminRequest{Name: "Dokter"}
	actor := regularActor()

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermMasterCreate, mock.Anything).Return(true, nil)
	s.repo.On("GetByNameJenisKelamin", req.Name).Return(nil, nil) // Added
	s.repo.On("CreateJenisKelamin", mock.AnythingOfType("*models.MasterJenisKelamin")).Return(nil)

	result, err := s.svc.CreateJenisKelamin(req, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *MasterServiceTestSuite) Test_CreateJenisKelamin_WithManagePermission_Success() {
	req := &dto.CreateMasterJenisKelaminRequest{Name: "Dokter"}
	actor := regularActor()

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermMasterCreate, mock.Anything).Return(false, nil)
	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermMasterManage, mock.Anything).Return(true, nil)
	s.repo.On("GetByNameJenisKelamin", req.Name).Return(nil, nil) // Added
	s.repo.On("CreateJenisKelamin", mock.AnythingOfType("*models.MasterJenisKelamin")).Return(nil)

	result, err := s.svc.CreateJenisKelamin(req, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *MasterServiceTestSuite) Test_CreateJenisKelamin_Forbidden() {
	req := &dto.CreateMasterJenisKelaminRequest{Name: "Dokter"}
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.CreateJenisKelamin(req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *MasterServiceTestSuite) Test_CreateJenisKelamin_RepoError() {
	req := &dto.CreateMasterJenisKelaminRequest{Name: "Dokter"}
	actor := superadminActor()

	s.repo.On("GetByNameJenisKelamin", req.Name).Return(nil, nil) // Added: harus lolos cek duplikat dulu
	s.repo.On("CreateJenisKelamin", mock.AnythingOfType("*models.MasterJenisKelamin")).Return(fmt.Errorf("db error"))

	result, err := s.svc.CreateJenisKelamin(req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *MasterServiceTestSuite) Test_GetByIDJenisKelamin_Success() {
	item := factories.NewJenisKelaminFactory().Make()
	item.ID = 1

	s.repo.On("GetByIDJenisKelamin", int64(1)).Return(item, nil)

	// FIXED: Typo GetByIDJenisKelamin -> GetByIDJenisKelamin
	result, err := s.svc.GetByIDJenisKelamin(1)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(item.ID, result.ID)
	s.Equal(item.Name, result.Name)
}

func (s *MasterServiceTestSuite) Test_GetByIDJenisKelamin_NotFound() {
	s.repo.On("GetByIDJenisKelamin", int64(999)).Return(nil, nil)

	result, err := s.svc.GetByIDJenisKelamin(999)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *MasterServiceTestSuite) Test_GetByIDJenisKelamin_RepoError() {
	s.repo.On("GetByIDJenisKelamin", int64(1)).Return(nil, fmt.Errorf("db error"))

	result, err := s.svc.GetByIDJenisKelamin(1)

	s.Nil(result)
	s.Error(err)
}

func (s *MasterServiceTestSuite) Test_ListJenisKelamin_Success() {
	filter := &dto.FilterMasterJenisKelaminRequest{}
	items := []models.MasterJenisKelamin{
		*factories.NewJenisKelaminFactory().Make(),
		*factories.NewJenisKelaminFactory().Make(),
	}

	// FIXED: Typo ListMasterJenisKelamin -> ListJenisKelamin
	s.repo.On("ListJenisKelamin", 1, 10, filter).Return(items, int64(2), nil)

	result, total, err := s.svc.ListJenisKelamin(1, 10, filter)

	s.NoError(err)
	s.Equal(int64(2), total)
	s.Len(result, 2)
}

func (s *MasterServiceTestSuite) Test_ListJenisKelamin_NotFound() {
	filter := &dto.FilterMasterJenisKelaminRequest{}

	// FIXED: Typo ListMasterJenisKelamin -> ListJenisKelamin
	s.repo.On("ListJenisKelamin", 1, 10, filter).Return([]models.MasterJenisKelamin{}, int64(0), nil)

	result, total, err := s.svc.ListJenisKelamin(1, 10, filter)

	s.Nil(result)
	s.Equal(int64(0), total)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *MasterServiceTestSuite) Test_ListJenisKelamin_DefaultPagination() {
	filter := &dto.FilterMasterJenisKelaminRequest{}
	items := []models.MasterJenisKelamin{*factories.NewJenisKelaminFactory().Make()}

	// FIXED: Typo ListMasterJenisKelamin -> ListJenisKelamin
	s.repo.On("ListJenisKelamin", 1, 10, filter).Return(items, int64(1), nil)

	_, _, err := s.svc.ListJenisKelamin(0, 0, filter)

	s.NoError(err)
	s.repo.AssertCalled(s.T(), "ListJenisKelamin", 1, 10, filter)
}

func (s *MasterServiceTestSuite) Test_ListJenisKelamin_PageSizeCapped() {
	filter := &dto.FilterMasterJenisKelaminRequest{}
	items := []models.MasterJenisKelamin{*factories.NewJenisKelaminFactory().Make()}

	// FIXED: Typo ListMasterJenisKelamin -> ListJenisKelamin
	s.repo.On("ListJenisKelamin", 1, 10, filter).Return(items, int64(1), nil)

	_, _, err := s.svc.ListJenisKelamin(1, 999, filter)

	s.NoError(err)
	s.repo.AssertCalled(s.T(), "ListJenisKelamin", 1, 10, filter)
}

func (s *MasterServiceTestSuite) Test_UpdateJenisKelamin_Superadmin_Success() {
	actor := superadminActor()
	existing := factories.NewJenisKelaminFactory().Make()
	existing.ID = 1
	newName := "Updated Name"
	req := &dto.UpdateMasterJenisKelaminRequest{Name: &newName}

	s.repo.On("GetByIDJenisKelamin", int64(1)).Return(existing, nil)
	// Added: Mock cek duplikasi nama untuk update
	s.repo.On("GetByNameJenisKelamin", newName).Return(nil, nil)
	s.repo.On("UpdateJenisKelamin", mock.AnythingOfType("*models.MasterJenisKelamin")).Return(nil)

	result, err := s.svc.UpdateJenisKelamin(1, req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(newName, result.Name)
}

func (s *MasterServiceTestSuite) Test_UpdateJenisKelamin_Forbidden() {
	actor := regularActor()
	req := &dto.UpdateMasterJenisKelaminRequest{}
	s.mockNoPermissions()

	result, err := s.svc.UpdateJenisKelamin(1, req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *MasterServiceTestSuite) Test_UpdateJenisKelamin_NotFound() {
	actor := superadminActor()
	req := &dto.UpdateMasterJenisKelaminRequest{}

	s.repo.On("GetByIDJenisKelamin", int64(999)).Return(nil, nil)

	result, err := s.svc.UpdateJenisKelamin(999, req, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *MasterServiceTestSuite) Test_UpdateJenisKelamin_PartialFields() {
	actor := superadminActor()
	existing := factories.NewJenisKelaminFactory().Make()
	existing.ID = 1
	newName := "New " + existing.Name
	req := &dto.UpdateMasterJenisKelaminRequest{Name: &newName}

	s.repo.On("GetByIDJenisKelamin", int64(1)).Return(existing, nil)
	// Added: Mock cek duplikasi nama untuk update
	s.repo.On("GetByNameJenisKelamin", newName).Return(nil, nil)
	s.repo.On("UpdateJenisKelamin", mock.MatchedBy(func(m *models.MasterJenisKelamin) bool {
		return m.Name == newName
	})).Return(nil)

	result, err := s.svc.UpdateJenisKelamin(1, req, actor)

	s.NoError(err)
	s.Equal(newName, result.Name)
}

func (s *MasterServiceTestSuite) Test_UpdateJenisKelamin_RepoError() {
	actor := superadminActor()
	existing := factories.NewJenisKelaminFactory().Make()
	existing.ID = 1
	req := &dto.UpdateMasterJenisKelaminRequest{}

	s.repo.On("GetByIDJenisKelamin", int64(1)).Return(existing, nil)
	s.repo.On("UpdateJenisKelamin", mock.AnythingOfType("*models.MasterJenisKelamin")).Return(fmt.Errorf("db error"))

	result, err := s.svc.UpdateJenisKelamin(1, req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *MasterServiceTestSuite) Test_DeleteJenisKelamin_Superadmin_Success() {
	actor := superadminActor()
	existing := factories.NewJenisKelaminFactory().Make()
	existing.ID = 1

	s.repo.On("GetByIDJenisKelamin", int64(1)).Return(existing, nil)
	s.repo.On("DeleteJenisKelamin", int64(1)).Return(nil)

	err := s.svc.DeleteJenisKelamin(1, actor)

	s.NoError(err)
	s.repo.AssertExpectations(s.T())
}

func (s *MasterServiceTestSuite) Test_DeleteJenisKelamin_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	err := s.svc.DeleteJenisKelamin(1, actor)

	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *MasterServiceTestSuite) Test_DeleteJenisKelamin_NotFound() {
	actor := superadminActor()

	s.repo.On("GetByIDJenisKelamin", int64(999)).Return(nil, nil)

	err := s.svc.DeleteJenisKelamin(999, actor)

	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *MasterServiceTestSuite) Test_DeleteJenisKelamin_RepoError() {
	actor := superadminActor()
	existing := factories.NewJenisKelaminFactory().Make()
	existing.ID = 1

	s.repo.On("GetByIDJenisKelamin", int64(1)).Return(existing, nil)
	s.repo.On("DeleteJenisKelamin", int64(1)).Return(fmt.Errorf("db error"))

	err := s.svc.DeleteJenisKelamin(1, actor)

	s.Error(err)
}
