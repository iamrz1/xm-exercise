package handlers_test

import (
	"github.com/stretchr/testify/mock"

	"xm-exercise/pkg/models"
)

// MockCompanyRepository is a mock implementation of db.CompanyRepository
type MockCompanyRepository struct {
	mock.Mock
}

func (m *MockCompanyRepository) Create(company *models.Company) error {
	args := m.Called(company)
	return args.Error(0)
}

func (m *MockCompanyRepository) GetByID(id string) (*models.Company, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Company), args.Error(1)
}

func (m *MockCompanyRepository) Update(company *models.Company) error {
	args := m.Called(company)
	return args.Error(0)
}

func (m *MockCompanyRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockCompanyRepository) ExistsByName(name string) (bool, error) {
	args := m.Called(name)
	return args.Bool(0), args.Error(1)
}

// MockKafkaProducer is a mock implementation of events.KafkaProducer
type MockKafkaProducer struct {
	mock.Mock
}

func (m *MockKafkaProducer) PublishCompanyUpdated(company *models.Company) error {
	args := m.Called(company)
	return args.Error(0)
}

func (m *MockKafkaProducer) PublishCompanyDeleted(companyID string) error {
	args := m.Called(companyID)
	return args.Error(0)
}

func (m *MockKafkaProducer) PublishCompanyCreated(company models.Company) error {
	args := m.Called(company)
	return args.Error(0)
}
