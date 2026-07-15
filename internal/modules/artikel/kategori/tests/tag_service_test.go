package tests

import (
	"fmt"
	"net/http"

	"github.com/stretchr/testify/mock"

	"neosim_go/internal/modules/artikel/kategori/dto"
	"neosim_go/internal/modules/artikel/kategori/models"
	"neosim_go/internal/modules/artikel/kategori/tests/factories"

	appErrors "neosim_go/internal/shared/errors"
)

// Catatan: suite ArtikelKategoriServiceTestSuite, TestMain, dan helper
// (superadminActor, regularActor, mockNoPermissions, dst) sudah didefinisikan
// di kategori_service_test.go. File ini HANYA menambah skenario test untuk
// Tag, memakai s.svc / s.repo yang SAMA (satu service & repository
// untuk seluruh sub-module kategori).

func (s *ArtikelKategoriServiceTestSuite) Test_CreateTag_Superadmin_Success() {
	req := &dto.CreateTagRequest{Name: "Test Tag"}
	actor := superadminActor()

	s.repo.On("CreateTag", mock.AnythingOfType("*models.Tag")).Return(nil)

	result, err := s.svc.CreateTag(req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(req.Name, result.Name)
}

func (s *ArtikelKategoriServiceTestSuite) Test_CreateTag_Forbidden() {
	req := &dto.CreateTagRequest{Name: "Test"}
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.CreateTag(req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *ArtikelKategoriServiceTestSuite) Test_CreateTag_RepoError() {
	req := &dto.CreateTagRequest{Name: "Test"}
	actor := superadminActor()

	s.repo.On("CreateTag", mock.AnythingOfType("*models.Tag")).Return(fmt.Errorf("db error"))

	result, err := s.svc.CreateTag(req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *ArtikelKategoriServiceTestSuite) Test_GetTagByID_Success() {
	actor := superadminActor()
	item := factories.NewTagFactory().Make()
	item.ID = 1

	s.repo.On("GetTagByID", int64(1)).Return(item, nil)

	result, err := s.svc.GetTagByID(1, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(item.ID, result.ID)
}

func (s *ArtikelKategoriServiceTestSuite) Test_GetTagByID_NotFound() {
	actor := superadminActor()

	s.repo.On("GetTagByID", int64(999)).Return(nil, nil)

	result, err := s.svc.GetTagByID(999, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *ArtikelKategoriServiceTestSuite) Test_GetTagByID_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.GetTagByID(1, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *ArtikelKategoriServiceTestSuite) Test_ListTag_Success() {
	actor := superadminActor()
	filter := &dto.FilterTagRequest{}
	items := []models.Tag{
		*factories.NewTagFactory().Make(),
		*factories.NewTagFactory().Make(),
	}

	s.repo.On("ListTag", 1, 10, filter).Return(items, int64(2), nil)

	result, total, err := s.svc.ListTag(1, 10, filter, actor)

	s.NoError(err)
	s.Equal(int64(2), total)
	s.Len(result, 2)
}

func (s *ArtikelKategoriServiceTestSuite) Test_ListTag_Forbidden() {
	actor := regularActor()
	filter := &dto.FilterTagRequest{}
	s.mockNoPermissions()

	result, total, err := s.svc.ListTag(1, 10, filter, actor)

	s.Nil(result)
	s.Equal(int64(0), total)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *ArtikelKategoriServiceTestSuite) Test_UpdateTag_Success() {
	actor := superadminActor()
	existing := factories.NewTagFactory().Make()
	existing.ID = 1
	newName := "Updated Name"
	req := &dto.UpdateTagRequest{Name: &newName}

	s.repo.On("GetTagByID", int64(1)).Return(existing, nil)
	s.repo.On("UpdateTag", mock.AnythingOfType("*models.Tag")).Return(nil)

	result, err := s.svc.UpdateTag(1, req, actor)

	s.NoError(err)
	s.Equal(newName, result.Name)
}

func (s *ArtikelKategoriServiceTestSuite) Test_UpdateTag_NotFound() {
	actor := superadminActor()
	req := &dto.UpdateTagRequest{}

	s.repo.On("GetTagByID", int64(999)).Return(nil, nil)

	result, err := s.svc.UpdateTag(999, req, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *ArtikelKategoriServiceTestSuite) Test_UpdateTag_Forbidden() {
	actor := regularActor()
	req := &dto.UpdateTagRequest{}
	s.mockNoPermissions()

	result, err := s.svc.UpdateTag(1, req, actor)

	s.Nil(result)
	s.Error(err)
}

func (s *ArtikelKategoriServiceTestSuite) Test_DeleteTag_Success() {
	actor := superadminActor()
	existing := factories.NewTagFactory().Make()
	existing.ID = 1

	s.repo.On("GetTagByID", int64(1)).Return(existing, nil)
	s.repo.On("DeleteTag", int64(1)).Return(nil)

	err := s.svc.DeleteTag(1, actor)

	s.NoError(err)
}

func (s *ArtikelKategoriServiceTestSuite) Test_DeleteTag_NotFound() {
	actor := superadminActor()

	s.repo.On("GetTagByID", int64(999)).Return(nil, nil)

	err := s.svc.DeleteTag(999, actor)

	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *ArtikelKategoriServiceTestSuite) Test_DeleteTag_Forbidden() {
	actor := regularActor()
	s.mockNoPermissions()

	err := s.svc.DeleteTag(1, actor)

	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}
