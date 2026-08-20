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
// StatusPernikahan
// =====================================================================

func (s *MasterServiceTestSuite) Test_CreateStatusPernikahan_Superadmin_Success() {
	req := &dto.CreateMasterStatusPernikahanRequest{Name: "Dokter"}
	actor := superadminActor()

	// Mock cek duplikasi nama (return nil agar dianggap belum ada)
	s.repo.On("GetByNameStatusPernikahan", req.Name).Return(nil, nil)
	s.repo.On("CreateStatusPernikahan", mock.AnythingOfType("*models.MasterStatusPernikahan")).Return(nil)

	result, err := s.svc.CreateStatusPernikahan(context.Background(), req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(req.Name, result.Name)
	s.repo.AssertExpectations(s.T())
}

func (s *MasterServiceTestSuite) Test_CreateStatusPernikahan_WithPermission_Success() {
	req := &dto.CreateMasterStatusPernikahanRequest{Name: "Dokter"}
	actor := regularActor()

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermMasterCreate, mock.Anything).Return(true, nil)
	s.repo.On("GetByNameStatusPernikahan", req.Name).Return(nil, nil) // Added
	s.repo.On("CreateStatusPernikahan", mock.AnythingOfType("*models.MasterStatusPernikahan")).Return(nil)

	result, err := s.svc.CreateStatusPernikahan(context.Background(), req, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *MasterServiceTestSuite) Test_CreateStatusPernikahan_WithManagePermission_Success() {
	req := &dto.CreateMasterStatusPernikahanRequest{Name: "Dokter"}
	actor := regularActor()

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermMasterCreate, mock.Anything).Return(false, nil)
	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermMasterManage, mock.Anything).Return(true, nil)
	s.repo.On("GetByNameStatusPernikahan", req.Name).Return(nil, nil) // Added
	s.repo.On("CreateStatusPernikahan", mock.AnythingOfType("*models.MasterStatusPernikahan")).Return(nil)

	result, err := s.svc.CreateStatusPernikahan(context.Background(), req, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *MasterServiceTestSuite) Test_CreateStatusPernikahan_Forbidden() {
	req := &dto.CreateMasterStatusPernikahanRequest{Name: "Dokter"}
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.CreateStatusPernikahan(context.Background(), req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *MasterServiceTestSuite) Test_CreateStatusPernikahan_RepoError() {
	req := &dto.CreateMasterStatusPernikahanRequest{Name: "Dokter"}
	actor := superadminActor()

	s.repo.On("GetByNameStatusPernikahan", req.Name).Return(nil, nil) // Added: harus lolos cek duplikat dulu
	s.repo.On("CreateStatusPernikahan", mock.AnythingOfType("*models.MasterStatusPernikahan")).Return(fmt.Errorf("db error"))

	result, err := s.svc.CreateStatusPernikahan(context.Background(), req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *MasterServiceTestSuite) Test_GetByIDStatusPernikahan_Success() {
	item := factories.NewStatusPernikahanFactory().Make()
	item.ID = 1

	s.repo.On("GetByIDStatusPernikahan", int64(1)).Return(item, nil)

	// FIXED: Typo GetByIDPekerjaan -> GetByIDStatusPernikahan
	result, err := s.svc.GetByIDStatusPernikahan(context.Background(), 1)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(item.ID, result.ID)
	s.Equal(item.Name, result.Name)
}

func (s *MasterServiceTestSuite) Test_GetByIDStatusPernikahan_NotFound() {
	s.repo.On("GetByIDStatusPernikahan", int64(999)).Return(nil, nil)

	result, err := s.svc.GetByIDStatusPernikahan(context.Background(), 999)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *MasterServiceTestSuite) Test_GetByIDStatusPernikahan_RepoError() {
	s.repo.On("GetByIDStatusPernikahan", int64(1)).Return(nil, fmt.Errorf("db error"))

	result, err := s.svc.GetByIDStatusPernikahan(context.Background(), 1)

	s.Nil(result)
	s.Error(err)
}

func (s *MasterServiceTestSuite) Test_ListStatusPernikahan_Success() {
	filter := &dto.FilterMasterStatusPernikahanRequest{}
	items := []models.MasterStatusPernikahan{
		*factories.NewStatusPernikahanFactory().Make(),
		*factories.NewStatusPernikahanFactory().Make(),
	}

	// FIXED: Typo ListMasterStatusPernikahan -> ListStatusPernikahan
	s.repo.On("ListStatusPernikahan", 1, 10, filter).Return(items, int64(2), nil)

	result, total, err := s.svc.ListStatusPernikahan(context.Background(), 1, 10, filter)

	s.NoError(err)
	s.Equal(int64(2), total)
	s.Len(result, 2)
}

func (s *MasterServiceTestSuite) Test_ListStatusPernikahan_NotFound() {
	filter := &dto.FilterMasterStatusPernikahanRequest{}

	// FIXED: Typo ListMasterStatusPernikahan -> ListStatusPernikahan
	s.repo.On("ListStatusPernikahan", 1, 10, filter).Return([]models.MasterStatusPernikahan{}, int64(0), nil)

	result, total, err := s.svc.ListStatusPernikahan(context.Background(), 1, 10, filter)

	s.Nil(result)
	s.Equal(int64(0), total)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *MasterServiceTestSuite) Test_ListStatusPernikahan_DefaultPagination() {
	filter := &dto.FilterMasterStatusPernikahanRequest{}
	items := []models.MasterStatusPernikahan{*factories.NewStatusPernikahanFactory().Make()}

	// FIXED: Typo ListMasterStatusPernikahan -> ListStatusPernikahan
	s.repo.On("ListStatusPernikahan", 1, 10, filter).Return(items, int64(1), nil)

	_, _, err := s.svc.ListStatusPernikahan(context.Background(), 0, 0, filter)

	s.NoError(err)
	s.repo.AssertCalled(s.T(), "ListStatusPernikahan", 1, 10, filter)
}

func (s *MasterServiceTestSuite) Test_ListStatusPernikahan_PageSizeCapped() {
	filter := &dto.FilterMasterStatusPernikahanRequest{}
	items := []models.MasterStatusPernikahan{*factories.NewStatusPernikahanFactory().Make()}

	// FIXED: Typo ListMasterStatusPernikahan -> ListStatusPernikahan
	s.repo.On("ListStatusPernikahan", 1, 10, filter).Return(items, int64(1), nil)

	_, _, err := s.svc.ListStatusPernikahan(context.Background(), 1, 999, filter)

	s.NoError(err)
	s.repo.AssertCalled(s.T(), "ListStatusPernikahan", 1, 10, filter)
}

func (s *MasterServiceTestSuite) Test_UpdateStatusPernikahan_Superadmin_Success() {
	actor := superadminActor()
	existing := factories.NewStatusPernikahanFactory().Make()
	existing.ID = 1
	newName := "Updated Name"
	req := &dto.UpdateMasterStatusPernikahanRequest{Name: &newName}

	s.repo.On("GetByIDStatusPernikahan", int64(1)).Return(existing, nil)
	// Added: Mock cek duplikasi nama untuk update
	s.repo.On("GetByNameStatusPernikahan", newName).Return(nil, nil)
	s.repo.On("UpdateStatusPernikahan", mock.AnythingOfType("*models.MasterStatusPernikahan")).Return(nil)

	result, err := s.svc.UpdateStatusPernikahan(context.Background(), 1, req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(newName, result.Name)
}

func (s *MasterServiceTestSuite) Test_UpdateStatusPernikahan_Forbidden() {
	actor := regularActor()
	req := &dto.UpdateMasterStatusPernikahanRequest{}
	s.mockNoPermissions()

	result, err := s.svc.UpdateStatusPernikahan(context.Background(), 1, req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *MasterServiceTestSuite) Test_UpdateStatusPernikahan_NotFound() {
	actor := superadminActor()
	req := &dto.UpdateMasterStatusPernikahanRequest{}

	s.repo.On("GetByIDStatusPernikahan", int64(999)).Return(nil, nil)

	result, err := s.svc.UpdateStatusPernikahan(context.Background(), 999, req, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *MasterServiceTestSuite) Test_UpdateStatusPernikahan_PartialFields() {
	actor := superadminActor()
	existing := factories.NewStatusPernikahanFactory().Make()
	existing.ID = 1
	newName := "New " + existing.Name
	req := &dto.UpdateMasterStatusPernikahanRequest{Name: &newName}

	s.repo.On("GetByIDStatusPernikahan", int64(1)).Return(existing, nil)
	// Added: Mock cek duplikasi nama untuk update
	s.repo.On("GetByNameStatusPernikahan", newName).Return(nil, nil)
	s.repo.On("UpdateStatusPernikahan", mock.MatchedBy(func(m *models.MasterStatusPernikahan) bool {
		return m.Name == newName
	})).Return(nil)

	result, err := s.svc.UpdateStatusPernikahan(context.Background(), 1, req, actor)

	s.NoError(err)
	s.Equal(newName, result.Name)
}

func (s *MasterServiceTestSuite) Test_UpdateStatusPernikahan_RepoError() {
	actor := superadminActor()
	existing := factories.NewStatusPernikahanFactory().Make()
	existing.ID = 1
	req := &dto.UpdateMasterStatusPernikahanRequest{}

	s.repo.On("GetByIDStatusPernikahan", int64(1)).Return(existing, nil)
	s.repo.On("UpdateStatusPernikahan", mock.AnythingOfType("*models.MasterStatusPernikahan")).Return(fmt.Errorf("db error"))

	result, err := s.svc.UpdateStatusPernikahan(context.Background(), 1, req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *MasterServiceTestSuite) Test_DeleteStatusPernikahan_Superadmin_Success() {
	actor := superadminActor()
	existing := factories.NewStatusPernikahanFactory().Make()
	existing.ID = 1

	s.repo.On("GetByIDStatusPernikahan", int64(1)).Return(existing, nil)
	s.repo.On("DeleteStatusPernikahan", int64(1)).Return(nil)

	err := s.svc.DeleteStatusPernikahan(context.Background(), 1, actor)

	s.NoError(err)
	s.repo.AssertExpectations(s.T())
}

func (s *MasterServiceTestSuite) Test_DeleteStatusPernikahan_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	err := s.svc.DeleteStatusPernikahan(context.Background(), 1, actor)

	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *MasterServiceTestSuite) Test_DeleteStatusPernikahan_NotFound() {
	actor := superadminActor()

	s.repo.On("GetByIDStatusPernikahan", int64(999)).Return(nil, nil)

	err := s.svc.DeleteStatusPernikahan(context.Background(), 999, actor)

	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *MasterServiceTestSuite) Test_DeleteStatusPernikahan_RepoError() {
	actor := superadminActor()
	existing := factories.NewStatusPernikahanFactory().Make()
	existing.ID = 1

	s.repo.On("GetByIDStatusPernikahan", int64(1)).Return(existing, nil)
	s.repo.On("DeleteStatusPernikahan", int64(1)).Return(fmt.Errorf("db error"))

	err := s.svc.DeleteStatusPernikahan(context.Background(), 1, actor)

	s.Error(err)
}
