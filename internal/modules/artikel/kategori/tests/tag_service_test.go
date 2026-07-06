package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"neosim_go/internal/modules/artikel/kategori/contracts"
	"neosim_go/internal/modules/artikel/kategori/dto"
	"neosim_go/internal/modules/artikel/kategori/models"
	"neosim_go/internal/modules/artikel/kategori/services"
	"neosim_go/internal/modules/artikel/kategori/tests/mocks"
	"neosim_go/config"
	rbacModels "neosim_go/internal/modules/rbac/models"
	appErrors "neosim_go/internal/shared/errors"
	he "neosim_go/internal/shared/httputil"
)

type TagServiceTestSuite struct {
	suite.Suite
	repo     *mocks.TagRepositoryMock
	rbacRepo *mocks.RBACRepositoryMock
	authRepo *mocks.AuthRepositoryMock
	svc      contracts.TagService
}

func (s *TagServiceTestSuite) SetupTest() {
	s.repo     = new(mocks.TagRepositoryMock)
	s.rbacRepo = new(mocks.RBACRepositoryMock)
	s.authRepo = new(mocks.AuthRepositoryMock)
	s.cfg      = &config.Config{DefaultPageSize: 10, DefaultPageSizeMax: 100}
	s.svc = services.NewTagService(s.repo, s.rbacRepo, s.authRepo, s.cfg)
}

func TestTagService(t *testing.T) {
	suite.Run(t, new(TagServiceTestSuite))
}

func (s *TagServiceTestSuite) superadminActor() he.AuthContext {
	return he.AuthContext{UserID: 1, IsSuperadmin: true}
}

func (s *TagServiceTestSuite) regularActor() he.AuthContext {
	return he.AuthContext{UserID: 2, IsSuperadmin: false}
}

func (s *TagServiceTestSuite) mockNoPermissions() {
	s.rbacRepo.On("HasPermission", int64(2), mock.Anything).Return(false, nil)
}

func (s *TagServiceTestSuite) Test_Create_Superadmin_Success() {
	actor := s.superadminActor()
	req := &dto.CreateTagRequest{Name: "Test Tag"}

	s.repo.On("Create", mock.AnythingOfType("*models.Tag")).Return(nil)

	result, err := s.svc.Create(req, &actor.UserID, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(req.Name, result.Name)
}

func (s *TagServiceTestSuite) Test_Create_WithPermission_Success() {
	actor := s.regularActor()
	req := &dto.CreateTagRequest{Name: "Test"}

	s.rbacRepo.On("HasPermission", actor.UserID, rbacModels.PermAnyCreate).Return(true, nil)
	s.repo.On("Create", mock.AnythingOfType("*models.Tag")).Return(nil)

	result, err := s.svc.Create(req, &actor.UserID, actor)

	s.NoError(err)
	s.NotNil(result)
}

func (s *TagServiceTestSuite) Test_Create_Forbidden() {
	actor := s.regularActor()
	req := &dto.CreateTagRequest{Name: "Test"}
	s.mockNoPermissions()

	result, err := s.svc.Create(req, &actor.UserID, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *TagServiceTestSuite) Test_Create_RepoError() {
	actor := s.superadminActor()
	req := &dto.CreateTagRequest{Name: "Test"}

	s.repo.On("Create", mock.AnythingOfType("*models.Tag")).Return(fmt.Errorf("db error"))

	result, err := s.svc.Create(req, &actor.UserID, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *TagServiceTestSuite) Test_GetByID_Superadmin_Success() {
	actor := s.superadminActor()
	item := &models.Tag{ID: 1, Name: "Test"}

	s.repo.On("GetByID", int64(1)).Return(item, nil)

	result, err := s.svc.GetByID(1, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(int64(1), result.ID)
}

func (s *TagServiceTestSuite) Test_GetByID_NotFound() {
	actor := s.superadminActor()

	s.repo.On("GetByID", int64(999)).Return(nil, nil)

	result, err := s.svc.GetByID(999, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *TagServiceTestSuite) Test_GetByID_Forbidden() {
	actor := s.regularActor()
	s.mockNoPermissions()

	result, err := s.svc.GetByID(1, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *TagServiceTestSuite) Test_List_Superadmin_Success() {
	actor := s.superadminActor()
	filter := &dto.FilterTagRequest{}
	itemA := models.Tag{ID: 1, Name: "A"}
	itemB := models.Tag{ID: 2, Name: "B"}
	items := []models.Tag{itemA, itemB}

	s.repo.On("List", 1, 10, filter).Return(items, int64(2), nil)

	result, total, err := s.svc.List(1, 10, filter, actor)

	s.NoError(err)
	s.Equal(int64(2), total)
	s.Len(result, 2)
}

func (s *TagServiceTestSuite) Test_List_DefaultPagination() {
	actor := s.superadminActor()
	filter := &dto.FilterTagRequest{}

	s.repo.On("List", 1, 10, filter).Return([]models.Tag{}, int64(0), nil)

	_, _, err := s.svc.List(0, 999, filter, actor)

	s.NoError(err)
	s.repo.AssertCalled(s.T(), "List", 1, 10, filter)
}

func (s *TagServiceTestSuite) Test_List_Forbidden() {
	actor := s.regularActor()
	filter := &dto.FilterTagRequest{}
	s.mockNoPermissions()

	result, total, err := s.svc.List(1, 10, filter, actor)

	s.Nil(result)
	s.Equal(int64(0), total)
	s.Error(err)
}

func (s *TagServiceTestSuite) Test_Update_Superadmin_Success() {
	actor := s.superadminActor()
	existing := &models.Tag{ID: 1, Name: "Old"}
	newName := "New Name"
	req := &dto.UpdateTagRequest{Name: &newName}

	s.repo.On("GetByID", int64(1)).Return(existing, nil)
	s.repo.On("Update", mock.AnythingOfType("*models.Tag")).Return(nil)

	result, err := s.svc.Update(1, req, &actor.UserID, actor)

	s.NoError(err)
	s.Equal(newName, result.Name)
}

func (s *TagServiceTestSuite) Test_Update_NotFound() {
	actor := s.superadminActor()
	req := &dto.UpdateTagRequest{}

	s.repo.On("GetByID", int64(999)).Return(nil, nil)

	result, err := s.svc.Update(999, req, &actor.UserID, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *TagServiceTestSuite) Test_Update_Forbidden() {
	actor := s.regularActor()
	req := &dto.UpdateTagRequest{}
	s.mockNoPermissions()

	result, err := s.svc.Update(1, req, &actor.UserID, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *TagServiceTestSuite) Test_Delete_Superadmin_Success() {
	actor := s.superadminActor()
	existing := &models.Tag{ID: 1, Name: "Test"}

	s.repo.On("GetByID", int64(1)).Return(existing, nil)
	s.repo.On("Delete", int64(1)).Return(nil)

	err := s.svc.Delete(1, actor)

	s.NoError(err)
}

func (s *TagServiceTestSuite) Test_Delete_NotFound() {
	actor := s.superadminActor()

	s.repo.On("GetByID", int64(999)).Return(nil, nil)

	err := s.svc.Delete(999, actor)

	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *TagServiceTestSuite) Test_Delete_Forbidden() {
	actor := s.regularActor()
	s.mockNoPermissions()

	err := s.svc.Delete(1, actor)

	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *TagServiceTestSuite) Test_Delete_RepoError() {
	actor := s.superadminActor()
	existing := &models.Tag{ID: 1, Name: "Test"}

	s.repo.On("GetByID", int64(1)).Return(existing, nil)
	s.repo.On("Delete", int64(1)).Return(fmt.Errorf("db error"))

	err := s.svc.Delete(1, actor)

	s.Error(err)
}
