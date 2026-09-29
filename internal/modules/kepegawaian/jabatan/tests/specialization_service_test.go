package tests

import (
	"context"
	"net/http"

	"github.com/stretchr/testify/mock"

	"neosim_go/internal/modules/kepegawaian/jabatan/dto"
	"neosim_go/internal/modules/kepegawaian/jabatan/models"
	"neosim_go/internal/modules/kepegawaian/jabatan/tests/factories"

	appErrors "neosim_go/internal/shared/errors"
)

func (s *KepegawaianJabatanServiceTestSuite) Test_CreateSpecialization_Superadmin_Success() {
	req := &dto.CreateSpecializationRequest{
		Code:       "SP-001",
		Label:      "Spesialis Test",
		JobTitleID: ptrTo(int64(1)),
		KategoriID: ptrTo(int64(1)),
	}
	actor := superadminActor()

	// 1. Mock pre-checks (1 argumen: code, sesuai m.Called(code))
	s.repo.On("GetSpecializationByCode", "SP-001").Return(nil, nil)

	// 2. Mock Create (1 argumen: item, sesuai m.Called(item))
	s.repo.On("CreateSpecialization", mock.AnythingOfType("*models.Specialization")).Return(nil).Run(func(args mock.Arguments) {
		// PERBAIKAN PENTING: Gunakan args.Get(0), bukan args.Get(1)
		// Karena di mock Anda: m.Called(item) hanya punya 1 argumen.
		args.Get(0).(*models.Specialization).ID = 1
	})

	// 3. Mock reload (1 argumen: id, sesuai m.Called(id))
	saved := factories.NewSpecializationFactory().
		With("code", "SP-001").
		With("label", "Spesialis Test").
		With("job_title_id", int64(1)).
		With("kategori_id", int64(1)).
		Make()
	saved.ID = 1
	s.repo.On("GetSpecializationByID", int64(1)).Return(saved, nil)

	result, err := s.svc.CreateSpecialization(context.Background(), req, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal("SP-001", result.Code)
	s.Equal("Spesialis Test", result.Label)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_CreateSpecialization_Forbidden() {
	req := &dto.CreateSpecializationRequest{Code: "SP-002", Label: "Test"}
	actor := regularActor()
	s.mockNoPermissions()

	result, err := s.svc.CreateSpecialization(context.Background(), req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusForbidden, appErr.Code)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_CreateSpecialization_CodeConflict() {
	req := &dto.CreateSpecializationRequest{Code: "SP-003", Label: "Test"}
	actor := superadminActor()

	existing := factories.NewSpecializationFactory().With("code", "SP-003").Make()
	existing.ID = 1

	// GetSpecializationByCode di mock ini dipanggil dengan 1 argumen (code saja, tanpa ctx)
	s.repo.On("GetSpecializationByCode", "SP-003").Return(existing, nil)

	result, err := s.svc.CreateSpecialization(context.Background(), req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusConflict, appErr.Code)

	s.repo.AssertNotCalled(s.T(), "CreateSpecialization", mock.Anything)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_GetSpecializationByID_Success() {
	actor := superadminActor()
	item := factories.NewSpecializationFactory().Make()
	item.ID = 1

	// 1 argumen: id
	s.repo.On("GetSpecializationByID", int64(1)).Return(item, nil)

	result, err := s.svc.GetSpecializationByID(context.Background(), 1, actor)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(item.ID, result.ID)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_GetSpecializationByID_NotFound() {
	actor := superadminActor()

	// 1 argumen: id
	s.repo.On("GetSpecializationByID", int64(999)).Return(nil, nil)

	result, err := s.svc.GetSpecializationByID(context.Background(), 999, actor)

	s.Nil(result)
	s.Error(err)
	s.Contains(err.Error(), "tidak ditemukan")
}

func (s *KepegawaianJabatanServiceTestSuite) Test_ListSpecialization_Success() {
	actor := superadminActor()
	filter := &dto.FilterSpecializationRequest{}
	items := []models.Specialization{
		*factories.NewSpecializationFactory().Make(),
		*factories.NewSpecializationFactory().Make(),
	}

	// 3 argumen: page, pageSize, filter (sesuai m.Called(page, pageSize, filter))
	s.repo.On("ListSpecialization", 1, 10, filter).Return(items, int64(2), nil)

	result, total, err := s.svc.ListSpecialization(context.Background(), 1, 10, filter, actor)

	s.NoError(err)
	s.Equal(int64(2), total)
	s.Len(result, 2)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_ListSelectSpecialization_Success() {
	actor := superadminActor()
	items := []models.Specialization{*factories.NewSpecializationFactory().Make()}

	// 1 argumen: search
	s.repo.On("ListSelectSpecialization", "").Return(items, nil)

	result, err := s.svc.ListSelectSpecialization(context.Background(), "", actor)

	s.NoError(err)
	s.Len(result, 1)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_UpdateSpecialization_Success() {
	actor := superadminActor()
	existing := factories.NewSpecializationFactory().Make()
	existing.ID = 1
	newLabel := "Updated Label"
	req := &dto.UpdateSpecializationRequest{Label: &newLabel}

	// 1 argumen: id (dipanggil 2x: validasi & reload)
	s.repo.On("GetSpecializationByID", int64(1)).Return(existing, nil).Twice()

	// 1 argumen: item
	s.repo.On("UpdateSpecialization", mock.AnythingOfType("*models.Specialization")).Return(nil)

	result, err := s.svc.UpdateSpecialization(context.Background(), 1, req, actor)

	s.NoError(err)
	s.Equal(newLabel, result.Label)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_UpdateSpecialization_CodeConflict() {
	actor := superadminActor()
	existing := factories.NewSpecializationFactory().Make()
	existing.ID = 1
	existing.Code = "OLD-CODE"
	newCode := "NEW-CODE"
	req := &dto.UpdateSpecializationRequest{Code: &newCode}

	conflictItem := factories.NewSpecializationFactory().Make()
	conflictItem.ID = 2
	conflictItem.Code = "NEW-CODE"

	s.repo.On("GetSpecializationByID", int64(1)).Return(existing, nil)
	s.repo.On("GetSpecializationByCode", "NEW-CODE").Return(conflictItem, nil)

	result, err := s.svc.UpdateSpecialization(context.Background(), 1, req, actor)

	s.Nil(result)
	s.Error(err)
	var appErr *appErrors.AppError
	s.ErrorAs(err, &appErr)
	s.Equal(http.StatusConflict, appErr.Code)
}

func (s *KepegawaianJabatanServiceTestSuite) Test_DeleteSpecialization_Success() {
	actor := superadminActor()
	existing := factories.NewSpecializationFactory().Make()
	existing.ID = 1

	// 1 argumen: id
	s.repo.On("GetSpecializationByID", int64(1)).Return(existing, nil)

	// 2 argumen: id, deletedBy (sesuai m.Called(id, deletedBy))
	s.repo.On("DeleteSpecialization", int64(1), actor.UserID).Return(nil)

	err := s.svc.DeleteSpecialization(context.Background(), 1, actor)

	s.NoError(err)
}
