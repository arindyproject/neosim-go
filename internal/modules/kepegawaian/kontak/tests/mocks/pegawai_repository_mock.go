package mocks

import (
	"context"

	"neosim_go/internal/modules/kepegawaian/pegawai/dto"
	"neosim_go/internal/modules/kepegawaian/pegawai/models"

	"github.com/stretchr/testify/mock"
)

// KepegawaianPegawaiRepositoryMock adalah struct mock untuk repository Pegawai
type KepegawaianPegawaiRepositoryMock struct {
	mock.Mock
}

// ── Create ────────────────────────────────────────────────────────────────────
func (m *KepegawaianPegawaiRepositoryMock) CreatePegawai(ctx context.Context, item *models.KepegawaianPegawai) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

// ── GetByID ───────────────────────────────────────────────────────────────────
func (m *KepegawaianPegawaiRepositoryMock) GetPegawaiByID(ctx context.Context, id int64) (*models.KepegawaianPegawai, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.KepegawaianPegawai), args.Error(1)
}

// ── List ──────────────────────────────────────────────────────────────────────
func (m *KepegawaianPegawaiRepositoryMock) ListPegawai(ctx context.Context, page, pageSize int, filter *dto.FilterKepegawaianPegawaiRequest) ([]models.KepegawaianPegawai, int64, error) {
	args := m.Called(ctx, page, pageSize, filter)

	var items []models.KepegawaianPegawai
	if args.Get(0) != nil {
		items = args.Get(0).([]models.KepegawaianPegawai)
	}

	return items, args.Get(1).(int64), args.Error(2)
}

// ── Update ────────────────────────────────────────────────────────────────────
func (m *KepegawaianPegawaiRepositoryMock) UpdatePegawai(ctx context.Context, item *models.KepegawaianPegawai) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

// ── Delete ────────────────────────────────────────────────────────────────────
func (m *KepegawaianPegawaiRepositoryMock) DeletePegawai(ctx context.Context, id int64, deletedBy int64) error {
	args := m.Called(ctx, id, deletedBy)
	return args.Error(0)
}
