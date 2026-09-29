package mocks

import (
	"context"

	"neosim_go/internal/modules/kepegawaian/jabatan/dto"
	"neosim_go/internal/modules/kepegawaian/jabatan/models"
)

// Method di bawah ini ditempelkan ke KepegawaianJabatanRepositoryMock yang
// sama dengan mock entitas utama (lihat tests/mocks/jabatan_repository_mock.go).

func (m *KepegawaianJabatanRepositoryMock) CreatePosition(ctx context.Context, item *models.Position) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *KepegawaianJabatanRepositoryMock) GetPositionByID(ctx context.Context, id int64) (*models.Position, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Position), args.Error(1)
}

func (m *KepegawaianJabatanRepositoryMock) ListSelectPosition(ctx context.Context, search string) ([]models.Position, error) {
	args := m.Called(ctx, search)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Position), args.Error(1)
}

func (m *KepegawaianJabatanRepositoryMock) ListPosition(ctx context.Context, page, pageSize int, filter *dto.FilterPositionRequest) ([]models.Position, int64, error) {
	args := m.Called(ctx, page, pageSize, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.Position), args.Get(1).(int64), args.Error(2)
}

func (m *KepegawaianJabatanRepositoryMock) UpdatePosition(ctx context.Context, item *models.Position) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *KepegawaianJabatanRepositoryMock) DeletePosition(ctx context.Context, id int64, deletedBy int64) error {
	args := m.Called(ctx, id, deletedBy)
	return args.Error(0)
}

func (m *KepegawaianJabatanRepositoryMock) ExistsByNameAndKategori(ctx context.Context, kategoriID int64, name string, excludeID int64) (bool, error) {
	args := m.Called(ctx, kategoriID, name, excludeID)
	return args.Bool(0), args.Error(1)
}

func (m *KepegawaianJabatanRepositoryMock) ExistsPositionByID(ctx context.Context, id int64) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *KepegawaianJabatanRepositoryMock) FindAllPositions(ctx context.Context, onlyAktif bool) ([]models.Position, error) {
	args := m.Called(ctx, onlyAktif)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Position), args.Error(1)
}

func (m *KepegawaianJabatanRepositoryMock) FindChildrenByParentID(ctx context.Context, parentID int64) ([]models.Position, error) {
	args := m.Called(ctx, parentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Position), args.Error(1)
}

func (m *KepegawaianJabatanRepositoryMock) FindRootPositions(ctx context.Context) ([]models.Position, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Position), args.Error(1)
}

func (m *KepegawaianJabatanRepositoryMock) GetPositionWithChildrenByID(ctx context.Context, id int64) (*models.Position, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Position), args.Error(1)
}

func (m *KepegawaianJabatanRepositoryMock) HasChildren(ctx context.Context, id int64) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}
