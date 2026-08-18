package mocks

import (
	"neosim_go/internal/modules/kepegawaian/kontak/dto"
	"neosim_go/internal/modules/kepegawaian/kontak/models"
)

// Method di bawah ini ditempelkan ke KepegawaianKontakRepositoryMock yang
// sama dengan mock entitas utama (lihat tests/mocks/kontak_repository_mock.go).

func (m *KepegawaianKontakRepositoryMock) CreateTipe(item *models.Tipe) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *KepegawaianKontakRepositoryMock) GetTipeByID(id int64) (*models.Tipe, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Tipe), args.Error(1)
}

func (m *KepegawaianKontakRepositoryMock) ListTipe(page, pageSize int, filter *dto.FilterTipeRequest) ([]models.Tipe, int64, error) {
	args := m.Called(page, pageSize, filter)
	return args.Get(0).([]models.Tipe), args.Get(1).(int64), args.Error(2)
}

func (m *KepegawaianKontakRepositoryMock) UpdateTipe(item *models.Tipe) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *KepegawaianKontakRepositoryMock) DeleteTipe(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}
