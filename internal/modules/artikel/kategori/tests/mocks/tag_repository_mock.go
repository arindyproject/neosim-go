package mocks

import (
	"neosim_go/internal/modules/artikel/kategori/dto"
	"neosim_go/internal/modules/artikel/kategori/models"
	"github.com/stretchr/testify/mock"
)

// TagRepositoryMock is a mock implementation of contracts.TagRepository
type TagRepositoryMock struct {
	mock.Mock
}

func (m *TagRepositoryMock) Create(item *models.Tag) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *TagRepositoryMock) GetByID(id int64) (*models.Tag, error) {
	args := m.Called(id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.Tag), args.Error(1)
}

func (m *TagRepositoryMock) List(page, pageSize int, filter *dto.FilterTagRequest) ([]models.Tag, int64, error) {
	args := m.Called(page, pageSize, filter)
	return args.Get(0).([]models.Tag), args.Get(1).(int64), args.Error(2)
}

func (m *TagRepositoryMock) Update(item *models.Tag) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *TagRepositoryMock) Delete(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}
