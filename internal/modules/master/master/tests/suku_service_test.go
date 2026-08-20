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
// Suku
// =====================================================================

func (s *MasterServiceTestSuite) Test_CreateSuku_Superadmin_Success() {
	req := &dto.CreateMasterSukuRequest{Name: "Dokter"}
	actor := superadminActor()

	// Mock cek duplikasi nama (return nil agar dianggap belum ada)
	s.repo.On("GetByNameSuku", req.Name).Return(nil, nil)
	s.repo.On("CreateSuku", mock.AnythingOfType("*models.MasterSuku")).Return(nil)

	result, err := s.svc.CreateSuku(context.Background(), req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(req.Name, result.Name)
	s.repo.AssertExpectations(s.T())
}

func (s *MasterServiceTestSuite) Test_CreateSuku_WithPermission_Success() {
	req := &dto.CreateMasterSukuRequest{Name: "Dokter"}
	actor := regularActor()

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermMasterCreate, mock.Anything).Return(true, nil)
	s.repo.On("GetByNameSuku", req.Name).Return(nil, nil) // Added
	s.repo.On("CreateSuku", mock.AnythingOfType("*models.MasterSuku")).Return(nil)

	result, err := s.svc.CreateSuku(context.Background(), req, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *MasterServiceTestSuite) Test_CreateSuku_WithManagePermission_Success() {
	req := &dto.CreateMasterSukuRequest{Name: "Dokter"}
	actor := regularActor()

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermMasterCreate, mock.Anything).Return(false, nil)
	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermMasterManage, mock.Anything).Return(true, nil)
	s.repo.On("GetByNameSuku", req.Name).Return(nil, nil) // Added
	s.repo.On("CreateSuku", mock.AnythingOfType("*models.MasterSuku")).Return(nil)

	result, err := s.svc.CreateSuku(context.Background(), req, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *MasterServiceTestSuite) Test_CreateSuku_Forbidden() {
	req := &dto.CreateMasterSukuRequest{Name: "Dokter"}
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.CreateSuku(context.Background(), req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *MasterServiceTestSuite) Test_CreateSuku_RepoError() {
	req := &dto.CreateMasterSukuRequest{Name: "Dokter"}
	actor := superadminActor()

	s.repo.On("GetByNameSuku", req.Name).Return(nil, nil) // Added: harus lolos cek duplikat dulu
	s.repo.On("CreateSuku", mock.AnythingOfType("*models.MasterSuku")).Return(fmt.Errorf("db error"))

	result, err := s.svc.CreateSuku(context.Background(), req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *MasterServiceTestSuite) Test_GetByIDSuku_Success() {
	item := factories.NewSukuFactory().Make()
	item.ID = 1

	s.repo.On("GetByIDSuku", int64(1)).Return(item, nil)

	// FIXED: Typo GetByIDSuku -> GetByIDSuku
	result, err := s.svc.GetByIDSuku(context.Background(), 1)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(item.ID, result.ID)
	s.Equal(item.Name, result.Name)
}

func (s *MasterServiceTestSuite) Test_GetByIDSuku_NotFound() {
	s.repo.On("GetByIDSuku", int64(999)).Return(nil, nil)

	result, err := s.svc.GetByIDSuku(context.Background(), 999)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *MasterServiceTestSuite) Test_GetByIDSuku_RepoError() {
	s.repo.On("GetByIDSuku", int64(1)).Return(nil, fmt.Errorf("db error"))

	result, err := s.svc.GetByIDSuku(context.Background(), 1)

	s.Nil(result)
	s.Error(err)
}

func (s *MasterServiceTestSuite) Test_ListSuku_Success() {
	filter := &dto.FilterMasterSukuRequest{}
	items := []models.MasterSuku{
		*factories.NewSukuFactory().Make(),
		*factories.NewSukuFactory().Make(),
	}

	// FIXED: Typo ListMasterSuku -> ListSuku
	s.repo.On("ListSuku", 1, 10, filter).Return(items, int64(2), nil)

	result, total, err := s.svc.ListSuku(context.Background(), 1, 10, filter)

	s.NoError(err)
	s.Equal(int64(2), total)
	s.Len(result, 2)
}

func (s *MasterServiceTestSuite) Test_ListSuku_NotFound() {
	filter := &dto.FilterMasterSukuRequest{}

	// FIXED: Typo ListMasterSuku -> ListSuku
	s.repo.On("ListSuku", 1, 10, filter).Return([]models.MasterSuku{}, int64(0), nil)

	result, total, err := s.svc.ListSuku(context.Background(), 1, 10, filter)

	s.Nil(result)
	s.Equal(int64(0), total)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *MasterServiceTestSuite) Test_ListSuku_DefaultPagination() {
	filter := &dto.FilterMasterSukuRequest{}
	items := []models.MasterSuku{*factories.NewSukuFactory().Make()}

	// FIXED: Typo ListMasterSuku -> ListSuku
	s.repo.On("ListSuku", 1, 10, filter).Return(items, int64(1), nil)

	_, _, err := s.svc.ListSuku(context.Background(), 0, 0, filter)

	s.NoError(err)
	s.repo.AssertCalled(s.T(), "ListSuku", 1, 10, filter)
}

func (s *MasterServiceTestSuite) Test_ListSuku_PageSizeCapped() {
	filter := &dto.FilterMasterSukuRequest{}
	items := []models.MasterSuku{*factories.NewSukuFactory().Make()}

	// FIXED: Typo ListMasterSuku -> ListSuku
	s.repo.On("ListSuku", 1, 10, filter).Return(items, int64(1), nil)

	_, _, err := s.svc.ListSuku(context.Background(), 1, 999, filter)

	s.NoError(err)
	s.repo.AssertCalled(s.T(), "ListSuku", 1, 10, filter)
}

func (s *MasterServiceTestSuite) Test_UpdateSuku_Superadmin_Success() {
	actor := superadminActor()
	existing := factories.NewSukuFactory().Make()
	existing.ID = 1
	newName := "Updated Name"
	req := &dto.UpdateMasterSukuRequest{Name: &newName}

	s.repo.On("GetByIDSuku", int64(1)).Return(existing, nil)
	// Added: Mock cek duplikasi nama untuk update
	s.repo.On("GetByNameSuku", newName).Return(nil, nil)
	s.repo.On("UpdateSuku", mock.AnythingOfType("*models.MasterSuku")).Return(nil)

	result, err := s.svc.UpdateSuku(context.Background(), 1, req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(newName, result.Name)
}

func (s *MasterServiceTestSuite) Test_UpdateSuku_Forbidden() {
	actor := regularActor()
	req := &dto.UpdateMasterSukuRequest{}
	s.mockNoPermissions()

	result, err := s.svc.UpdateSuku(context.Background(), 1, req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *MasterServiceTestSuite) Test_UpdateSuku_NotFound() {
	actor := superadminActor()
	req := &dto.UpdateMasterSukuRequest{}

	s.repo.On("GetByIDSuku", int64(999)).Return(nil, nil)

	result, err := s.svc.UpdateSuku(context.Background(), 999, req, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *MasterServiceTestSuite) Test_UpdateSuku_PartialFields() {
	actor := superadminActor()
	existing := factories.NewSukuFactory().Make()
	existing.ID = 1
	newName := "New " + existing.Name
	req := &dto.UpdateMasterSukuRequest{Name: &newName}

	s.repo.On("GetByIDSuku", int64(1)).Return(existing, nil)
	// Added: Mock cek duplikasi nama untuk update
	s.repo.On("GetByNameSuku", newName).Return(nil, nil)
	s.repo.On("UpdateSuku", mock.MatchedBy(func(m *models.MasterSuku) bool {
		return m.Name == newName
	})).Return(nil)

	result, err := s.svc.UpdateSuku(context.Background(), 1, req, actor)

	s.NoError(err)
	s.Equal(newName, result.Name)
}

func (s *MasterServiceTestSuite) Test_UpdateSuku_RepoError() {
	actor := superadminActor()
	existing := factories.NewSukuFactory().Make()
	existing.ID = 1
	req := &dto.UpdateMasterSukuRequest{}

	s.repo.On("GetByIDSuku", int64(1)).Return(existing, nil)
	s.repo.On("UpdateSuku", mock.AnythingOfType("*models.MasterSuku")).Return(fmt.Errorf("db error"))

	result, err := s.svc.UpdateSuku(context.Background(), 1, req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *MasterServiceTestSuite) Test_DeleteSuku_Superadmin_Success() {
	actor := superadminActor()
	existing := factories.NewSukuFactory().Make()
	existing.ID = 1

	s.repo.On("GetByIDSuku", int64(1)).Return(existing, nil)
	s.repo.On("DeleteSuku", int64(1)).Return(nil)

	err := s.svc.DeleteSuku(context.Background(), 1, actor)

	s.NoError(err)
	s.repo.AssertExpectations(s.T())
}

func (s *MasterServiceTestSuite) Test_DeleteSuku_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	err := s.svc.DeleteSuku(context.Background(), 1, actor)

	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *MasterServiceTestSuite) Test_DeleteSuku_NotFound() {
	actor := superadminActor()

	s.repo.On("GetByIDSuku", int64(999)).Return(nil, nil)

	err := s.svc.DeleteSuku(context.Background(), 999, actor)

	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *MasterServiceTestSuite) Test_DeleteSuku_RepoError() {
	actor := superadminActor()
	existing := factories.NewSukuFactory().Make()
	existing.ID = 1

	s.repo.On("GetByIDSuku", int64(1)).Return(existing, nil)
	s.repo.On("DeleteSuku", int64(1)).Return(fmt.Errorf("db error"))

	err := s.svc.DeleteSuku(context.Background(), 1, actor)

	s.Error(err)
}
