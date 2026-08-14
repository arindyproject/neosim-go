package mocks

import (
	"neosim_go/internal/modules/kepegawaian/kontak/dto"
	"neosim_go/internal/modules/kepegawaian/kontak/models"
	"github.com/stretchr/testify/mock"
)

// KepegawaianKontakRepositoryMock is a mock implementation of contracts.Repository.
// Ketika item ditambahkan (mode add-item), method mock untuk item tersebut
// ditempelkan ke struct INI JUGA (mis. tests/mocks/tag_repository_mock.go),
// bukan membuat mock struct baru.
type KepegawaianKontakRepositoryMock struct {
	mock.Mock
}

func (m *KepegawaianKontakRepositoryMock) CreateKontak(item *models.KepegawaianKontak) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *KepegawaianKontakRepositoryMock) GetKontakByID(id int64) (*models.KepegawaianKontak, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.KepegawaianKontak), args.Error(1)
}

func (m *KepegawaianKontakRepositoryMock) ListKontak(page, pageSize int, filter *dto.FilterKepegawaianKontakRequest) ([]models.KepegawaianKontak, int64, error) {
	args := m.Called(page, pageSize, filter)
	return args.Get(0).([]models.KepegawaianKontak), args.Get(1).(int64), args.Error(2)
}

func (m *KepegawaianKontakRepositoryMock) UpdateKontak(item *models.KepegawaianKontak) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *KepegawaianKontakRepositoryMock) DeleteKontak(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}
