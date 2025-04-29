package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"xm-exercise/internal/api/handlers"
	"xm-exercise/internal/api/middleware"
	"xm-exercise/internal/logger"
	"xm-exercise/pkg/models"
)

// newTestCompanyHandler is helper function to create a new CompanyHandler with mocks
func newTestCompanyHandler() (*handlers.CompanyHandler, *MockCompanyRepository, *MockKafkaProducer) {
	mockRepo := new(MockCompanyRepository)
	mockProducer := new(MockKafkaProducer)
	handler := handlers.NewCompanyHandler(mockRepo, mockProducer)
	return handler, mockRepo, mockProducer
}

func TestCompanyHandler_Create(t *testing.T) {
	err := logger.Init(zap.WarnLevel.String(), false)
	assert.NoError(t, err)
	t.Run("Successful Creation", func(t *testing.T) {
		handler, mockRepo, mockProducer := newTestCompanyHandler()

		companyCreateReq := models.CompanyCreateRequest{
			Name:          "Test Company",
			Description:   aws.String("A company for testing"),
			EmployeeCount: 100,
			Registered:    aws.Bool(true),
			Type:          "Corporations",
		}
		jsonBody, _ := json.Marshal(companyCreateReq)

		mockRepo.On("ExistsByName", companyCreateReq.Name).Return(false, nil).Once()
		mockRepo.On("Create", mock.AnythingOfType("*models.Company")).Return(nil).Once()
		mockProducer.On("PublishCompanyCreated", mock.AnythingOfType("models.Company")).Return(nil).Once()

		req, _ := http.NewRequest("POST", "/companies", bytes.NewBuffer(jsonBody))
		req = req.WithContext(middleware.SetUserID(req.Context(), uuid.New().String()))
		rr := httptest.NewRecorder()
		handler.Create(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
		var companyRes models.CompanyResponse
		err := json.NewDecoder(rr.Body).Decode(&companyRes)
		assert.NoError(t, err)
		assert.Equal(t, companyCreateReq.Name, companyRes.Name)
		assert.Equal(t, companyCreateReq.Description, companyRes.Description)
		assert.Equal(t, companyCreateReq.EmployeeCount, companyRes.EmployeeCount)
		assert.Equal(t, companyCreateReq.Registered, companyRes.Registered)
		assert.Equal(t, companyCreateReq.Type, companyRes.Type)
		assert.NotEmpty(t, companyRes.ID)
		assert.NotZero(t, companyRes.CreatedAt)
		assert.NotZero(t, companyRes.UpdatedAt)

		mockRepo.AssertExpectations(t)
		mockProducer.AssertExpectations(t)
	})

	t.Run("Invalid Request Body", func(t *testing.T) {
		handler, _, _ := newTestCompanyHandler()

		// Send invalid JSON
		jsonBody := []byte(`{"name": "Test Company", "invalid_field": true`)

		req, _ := http.NewRequest("POST", "/companies", bytes.NewBuffer(jsonBody))
		// Correctly add a mock user ID (as uuid.UUID) to the context for authentication
		req = req.WithContext(middleware.SetUserID(req.Context(), uuid.New().String()))
		rr := httptest.NewRecorder()

		handler.Create(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "Invalid request body")
	})

	t.Run("Validation Error", func(t *testing.T) {
		handler, mockRepo, _ := newTestCompanyHandler()

		// Missing required field (Name)
		companyCreateReq := models.CompanyCreateRequest{
			Description:   aws.String("A company for testing"),
			EmployeeCount: 100,
			Registered:    aws.Bool(true),
			Type:          "Corporations",
		}
		jsonBody, _ := json.Marshal(companyCreateReq)

		// No repository or producer methods should be called due to validation error
		mockRepo.AssertNotCalled(t, "ExistsByName", mock.Anything)
		mockRepo.AssertNotCalled(t, "Create", mock.Anything)

		req, _ := http.NewRequest("POST", "/companies", bytes.NewBuffer(jsonBody))
		// Correctly add a mock user ID (as uuid.UUID) to the context for authentication
		req = req.WithContext(middleware.SetUserID(req.Context(), uuid.New().String()))
		rr := httptest.NewRecorder()

		handler.Create(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "name is required") // Assuming your Validate method returns this error
	})

	t.Run("Company Name Already Exists", func(t *testing.T) {
		handler, mockRepo, _ := newTestCompanyHandler()

		companyCreateReq := models.CompanyCreateRequest{
			Name:          "Existing",
			Description:   aws.String("A company for testing"),
			EmployeeCount: 100,
			Registered:    aws.Bool(true),
			Type:          "Corporations",
		}
		jsonBody, _ := json.Marshal(companyCreateReq)

		// Expectation: Check if company name exists (should return true)
		mockRepo.On("ExistsByName", companyCreateReq.Name).Return(true, nil).Once()
		// No create or publish methods should be called
		mockRepo.AssertNotCalled(t, "Create", mock.Anything)
		mockRepo.AssertNotCalled(t, "Update", mock.Anything)
		mockRepo.AssertNotCalled(t, "Delete", mock.Anything)

		req, _ := http.NewRequest("POST", "/companies", bytes.NewBuffer(jsonBody))
		// Correctly add a mock user ID (as uuid.UUID) to the context for authentication
		req = req.WithContext(middleware.SetUserID(req.Context(), uuid.New().String()))
		rr := httptest.NewRecorder()

		handler.Create(rr, req)

		assert.Equal(t, http.StatusConflict, rr.Code)
		assert.Contains(t, rr.Body.String(), "Company name already exists")

		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Error on ExistsByName", func(t *testing.T) {
		handler, mockRepo, _ := newTestCompanyHandler()

		companyCreateReq := models.CompanyCreateRequest{
			Name:          "Test Company",
			Description:   aws.String("A company for testing"),
			EmployeeCount: 100,
			Registered:    aws.Bool(true),
			Type:          "Corporations",
		}
		jsonBody, _ := json.Marshal(companyCreateReq)

		// Expectation: Error when checking if company name exists
		mockRepo.On("ExistsByName", companyCreateReq.Name).Return(false, errors.New("database error")).Once()
		// No create or publish methods should be called
		mockRepo.AssertNotCalled(t, "Create", mock.Anything)
		mockRepo.AssertNotCalled(t, "Update", mock.Anything)
		mockRepo.AssertNotCalled(t, "Delete", mock.Anything)

		req, _ := http.NewRequest("POST", "/companies", bytes.NewBuffer(jsonBody))
		// Correctly add a mock user ID (as uuid.UUID) to the context for authentication
		req = req.WithContext(middleware.SetUserID(req.Context(), uuid.New().String()))
		rr := httptest.NewRecorder()

		handler.Create(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), "Error checking name for uniqueness")

		mockRepo.AssertExpectations(t)
	})

	t.Run("Repository Error on Create", func(t *testing.T) {
		handler, mockRepo, mockProducer := newTestCompanyHandler()

		companyCreateReq := models.CompanyCreateRequest{
			Name:          "Test Company",
			Description:   aws.String("A company for testing"),
			EmployeeCount: 100,
			Registered:    aws.Bool(true),
			Type:          "Corporations",
		}
		jsonBody, _ := json.Marshal(companyCreateReq)

		// Expectation: Check if company name exists (should return false)
		mockRepo.On("ExistsByName", companyCreateReq.Name).Return(false, nil).Once()
		// Expectation: Error when creating the company
		mockRepo.On("Create", mock.AnythingOfType("*models.Company")).Return(errors.New("database error")).Once()
		// No publish method should be called
		mockProducer.AssertNotCalled(t, "PublishCompanyCreated", mock.Anything)

		req, _ := http.NewRequest("POST", "/companies", bytes.NewBuffer(jsonBody))
		req = req.WithContext(middleware.SetUserID(req.Context(), uuid.New().String()))
		rr := httptest.NewRecorder()

		handler.Create(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), "Error creating company")

		mockRepo.AssertExpectations(t)
	})

	t.Run("Unauthorized Access", func(t *testing.T) {
		handler, mockRepo, mockProducer := newTestCompanyHandler()

		companyCreateReq := models.CompanyCreateRequest{
			Name:          "Test Company",
			Description:   aws.String("A company for testing"),
			EmployeeCount: 100,
			Registered:    aws.Bool(true),
			Type:          "Corporations",
		}
		jsonBody, _ := json.Marshal(companyCreateReq)

		// No repository or producer methods should be called if unauthorized
		mockRepo.AssertNotCalled(t, "ExistsByName", mock.Anything)
		mockRepo.AssertNotCalled(t, "Create", mock.Anything)
		mockProducer.AssertNotCalled(t, "PublishCompanyCreated", mock.Anything)

		req, _ := http.NewRequest("POST", "/companies", bytes.NewBuffer(jsonBody))
		rr := httptest.NewRecorder()

		handler.Create(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Contains(t, rr.Body.String(), "Unauthorized")
	})
}

func TestCompanyHandler_Get(t *testing.T) {
	t.Run("Successful Get", func(t *testing.T) {
		handler, mockRepo, _ := newTestCompanyHandler()

		companyID := uuid.New().String()
		expectedCompany := &models.Company{
			ID:            companyID,
			Name:          "Test Company",
			Description:   aws.String("A company for testing"),
			EmployeeCount: 100,
			Registered:    aws.Bool(true),
			Type:          models.TypeCorporation,
			CreatedAt:     time.Now().UTC(),
			UpdatedAt:     time.Now().UTC(),
		}
		mockRepo.On("GetByID", companyID).Return(expectedCompany, nil).Once()

		req, _ := http.NewRequest("GET", "/companies/"+companyID, nil)
		rr := httptest.NewRecorder()
		handler.Get(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var companyRes models.CompanyResponse
		err := json.NewDecoder(rr.Body).Decode(&companyRes)
		assert.NoError(t, err)
		assert.Equal(t, expectedCompany.ID, companyRes.ID)
		assert.Equal(t, expectedCompany.Name, companyRes.Name)
		assert.Equal(t, expectedCompany.Description, companyRes.Description)
		assert.Equal(t, expectedCompany.EmployeeCount, companyRes.EmployeeCount)
		assert.Equal(t, expectedCompany.Registered, companyRes.Registered)
		assert.Equal(t, string(expectedCompany.Type), string(companyRes.Type))

		mockRepo.AssertExpectations(t)
	})

	t.Run("Company Not Found", func(t *testing.T) {
		handler, mockRepo, _ := newTestCompanyHandler()

		companyID := uuid.New().String()
		mockRepo.On("GetByID", companyID).Return(&models.Company{}, errors.New("company not found")).Once()

		req, _ := http.NewRequest("GET", "/companies/"+companyID, nil)
		rr := httptest.NewRecorder()

		handler.Get(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.Contains(t, rr.Body.String(), "Company not found")

		mockRepo.AssertExpectations(t)
	})
}

func TestCompanyHandler_Patch(t *testing.T) {
	err := logger.Init(zap.WarnLevel.String(), false)
	assert.NoError(t, err)
	t.Run("Successful Patch", func(t *testing.T) {
		handler, mockRepo, mockProducer := newTestCompanyHandler()

		companyID := uuid.New().String()
		existingCompany := &models.Company{
			ID:            companyID,
			Name:          "OName",
			Description:   aws.String("Original Description"),
			EmployeeCount: 50,
			Registered:    aws.Bool(false),
			Type:          models.TypeSoleProprietor,
			CreatedAt:     time.Now().UTC().Add(-time.Hour),
			UpdatedAt:     time.Now().UTC().Add(-time.Hour),
		}

		newName := "UName"
		updatedEmployeeCount := 200
		updates := models.CompanyUpdateRequest{
			Name:          &newName,
			EmployeeCount: &updatedEmployeeCount,
		}

		jsonBody, _ := json.Marshal(updates)
		mockRepo.On("GetByID", companyID).Return(existingCompany, nil).Once()
		mockRepo.On("ExistsByName", newName).Return(false, nil).Once()
		mockRepo.On("Update", mock.AnythingOfType("*models.Company")).Return(nil).Once()
		mockProducer.On("PublishCompanyUpdated", mock.AnythingOfType("*models.Company")).Return(nil).Once()

		req, _ := http.NewRequest("PATCH", "/companies/"+companyID, bytes.NewBuffer(jsonBody))
		req = req.WithContext(middleware.SetUserID(req.Context(), uuid.New().String()))
		rr := httptest.NewRecorder()

		handler.Patch(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var companyRes models.CompanyResponse
		err := json.NewDecoder(rr.Body).Decode(&companyRes)
		assert.NoError(t, err)
		assert.Equal(t, newName, companyRes.Name)
		assert.Equal(t, updatedEmployeeCount, companyRes.EmployeeCount)
		assert.Equal(t, existingCompany.Type, companyRes.Type)

		mockRepo.AssertExpectations(t)
		mockProducer.AssertExpectations(t)
	})

	t.Run("Patch with Invalid Request Body", func(t *testing.T) {
		handler, mockRepo, _ := newTestCompanyHandler()
		companyID := uuid.New().String()
		jsonBody := []byte(`{"name": "Test Company", "invalid_field": true`)

		mockRepo.On("GetByID", companyID).Return(&models.Company{}, nil).Once()

		req, _ := http.NewRequest("PATCH", "/companies/"+companyID, bytes.NewBuffer(jsonBody))
		req = req.WithContext(middleware.SetUserID(req.Context(), uuid.New().String()))
		rr := httptest.NewRecorder()

		handler.Patch(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "Invalid request body")
	})

	t.Run("Patch with Validation Error", func(t *testing.T) {
		handler, mockRepo, _ := newTestCompanyHandler()

		companyID := uuid.New().String()
		existingCompany := &models.Company{
			ID:            companyID,
			Name:          "Original Name",
			Description:   aws.String("Original Description"),
			EmployeeCount: 50,
			Registered:    aws.Bool(false),
			Type:          models.TypeSoleProprietor,
			CreatedAt:     time.Now().UTC().Add(-time.Hour),
			UpdatedAt:     time.Now().UTC().Add(-time.Hour),
		}
		longName := "This name is way too long for the 15 character limit"
		updates := models.CompanyUpdateRequest{
			Name: &longName,
		}
		jsonBody, _ := json.Marshal(updates)

		mockRepo.On("GetByID", companyID).Return(existingCompany, nil).Once()
		mockRepo.AssertNotCalled(t, "ExistsByName", mock.Anything)
		mockRepo.AssertNotCalled(t, "Update", mock.Anything)
		mockRepo.AssertNotCalled(t, "PublishCompanyUpdated", mock.Anything)

		req, _ := http.NewRequest("PATCH", "/companies/"+companyID, bytes.NewBuffer(jsonBody))
		req = req.WithContext(middleware.SetUserID(req.Context(), uuid.New().String()))
		rr := httptest.NewRecorder()

		handler.Patch(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "name must be 15 characters or less")

		mockRepo.AssertExpectations(t)
	})

	t.Run("Patch Company Not Found", func(t *testing.T) {
		handler, mockRepo, _ := newTestCompanyHandler()

		companyID := uuid.New().String()
		newName := "Updated Company Name"
		updates := models.CompanyUpdateRequest{
			Name: &newName,
		}
		jsonBody, _ := json.Marshal(updates)

		mockRepo.On("GetByID", companyID).Return(&models.Company{}, errors.New("company not found")).Once()
		mockRepo.AssertNotCalled(t, "ExistsByName", mock.Anything)
		mockRepo.AssertNotCalled(t, "Update", mock.Anything)
		mockRepo.AssertNotCalled(t, "PublishCompanyUpdated", mock.Anything)

		req, _ := http.NewRequest("PATCH", "/companies/"+companyID, bytes.NewBuffer(jsonBody))
		req = req.WithContext(middleware.SetUserID(req.Context(), uuid.New().String()))
		rr := httptest.NewRecorder()

		handler.Patch(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.Contains(t, rr.Body.String(), "Company not found")

		mockRepo.AssertExpectations(t)
	})

	t.Run("Patch Company Name Already Exists", func(t *testing.T) {
		handler, mockRepo, _ := newTestCompanyHandler()

		companyID := uuid.New().String()
		existingCompany := &models.Company{
			ID:            companyID,
			Name:          "OName",
			Description:   aws.String("Original Description"),
			EmployeeCount: 50,
			Registered:    aws.Bool(false),
			Type:          models.TypeSoleProprietor,
			CreatedAt:     time.Now().UTC().Add(-time.Hour),
			UpdatedAt:     time.Now().UTC().Add(-time.Hour),
		}
		duplicateName := "EName"
		updates := models.CompanyUpdateRequest{
			Name: &duplicateName,
		}
		jsonBody, _ := json.Marshal(updates)

		mockRepo.On("GetByID", companyID).Return(existingCompany, nil).Once()
		mockRepo.On("ExistsByName", duplicateName).Return(true, nil).Once()
		mockRepo.AssertNotCalled(t, "Update", mock.Anything)
		mockRepo.AssertNotCalled(t, "PublishCompanyUpdated", mock.Anything)

		req, _ := http.NewRequest("PATCH", "/companies/"+companyID, bytes.NewBuffer(jsonBody))
		req = req.WithContext(middleware.SetUserID(req.Context(), uuid.New().String()))
		rr := httptest.NewRecorder()

		handler.Patch(rr, req)

		assert.Equal(t, http.StatusConflict, rr.Code)
		assert.Contains(t, rr.Body.String(), "Company name already exists")

		mockRepo.AssertExpectations(t)
	})

	t.Run("Patch Repository Error on GetByID", func(t *testing.T) {
		handler, mockRepo, _ := newTestCompanyHandler()

		companyID := uuid.New().String()
		newName := "Updated Company Name"
		updates := models.CompanyUpdateRequest{
			Name: &newName,
		}
		jsonBody, _ := json.Marshal(updates)

		mockRepo.On("GetByID", companyID).Return(&models.Company{}, errors.New("database error")).Once()
		mockRepo.AssertNotCalled(t, "ExistsByName", mock.Anything)
		mockRepo.AssertNotCalled(t, "Update", mock.Anything)
		mockRepo.AssertNotCalled(t, "PublishCompanyUpdated", mock.Anything)

		req, _ := http.NewRequest("PATCH", "/companies/"+companyID, bytes.NewBuffer(jsonBody))
		req = req.WithContext(middleware.SetUserID(req.Context(), uuid.New().String()))
		rr := httptest.NewRecorder()

		handler.Patch(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.Contains(t, rr.Body.String(), "Company not found")

		mockRepo.AssertExpectations(t)
	})

	t.Run("Patch Repository Error on ExistsByName", func(t *testing.T) {
		handler, mockRepo, _ := newTestCompanyHandler()

		companyID := uuid.New().String()
		existingCompany := &models.Company{
			ID:            companyID,
			Name:          "OName",
			Description:   aws.String("Original Description"),
			EmployeeCount: 50,
			Registered:    aws.Bool(false),
			Type:          models.TypeSoleProprietor,
			CreatedAt:     time.Now().UTC().Add(-time.Hour),
			UpdatedAt:     time.Now().UTC().Add(-time.Hour),
		}
		newName := "UName"
		updates := models.CompanyUpdateRequest{
			Name: &newName,
		}
		jsonBody, _ := json.Marshal(updates)
		mockRepo.On("GetByID", companyID).Return(existingCompany, nil).Once()
		mockRepo.On("ExistsByName", newName).Return(false, errors.New("database error")).Once()
		mockRepo.AssertNotCalled(t, "Update", mock.Anything)
		mockRepo.AssertNotCalled(t, "PublishCompanyUpdated", mock.Anything)

		req, _ := http.NewRequest("PATCH", "/companies/"+companyID, bytes.NewBuffer(jsonBody))
		req = req.WithContext(middleware.SetUserID(req.Context(), uuid.New().String()))
		rr := httptest.NewRecorder()

		handler.Patch(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), "Error checking name uniqueness")

		mockRepo.AssertExpectations(t)
	})

	t.Run("Patch Unauthorized Access", func(t *testing.T) {
		handler, mockRepo, mockProducer := newTestCompanyHandler()

		companyID := uuid.New().String()
		newName := "Updated Company Name"
		updates := models.CompanyUpdateRequest{
			Name: &newName,
		}
		jsonBody, _ := json.Marshal(updates)

		mockRepo.AssertNotCalled(t, "GetByID", mock.Anything)
		mockRepo.AssertNotCalled(t, "ExistsByName", mock.Anything)
		mockRepo.AssertNotCalled(t, "Update", mock.Anything)
		mockProducer.AssertNotCalled(t, "PublishCompanyUpdated", mock.Anything)

		req, _ := http.NewRequest("PATCH", "/companies/"+companyID, bytes.NewBuffer(jsonBody))
		rr := httptest.NewRecorder()

		handler.Patch(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Contains(t, rr.Body.String(), "Unauthorized")
	})
}

func TestCompanyHandler_Delete(t *testing.T) {
	err := logger.Init(zap.WarnLevel.String(), false)
	assert.NoError(t, err)

	t.Run("Successful Delete", func(t *testing.T) {
		handler, mockRepo, mockProducer := newTestCompanyHandler()

		companyID := uuid.New().String()

		mockRepo.On("GetByID", companyID).Return(&models.Company{}, nil).Once()
		mockRepo.On("Delete", companyID).Return(nil).Once()
		mockProducer.On("PublishCompanyDeleted", companyID).Return(nil).Once()

		req, _ := http.NewRequest("DELETE", "/companies/"+companyID, nil)
		req = req.WithContext(middleware.SetUserID(req.Context(), uuid.New().String()))
		rr := httptest.NewRecorder()

		handler.Delete(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		mockRepo.AssertExpectations(t)
		mockProducer.AssertExpectations(t)
	})

	t.Run("Company Not Found", func(t *testing.T) {
		handler, mockRepo, mockProducer := newTestCompanyHandler()

		companyID := uuid.New().String()

		mockRepo.On("GetByID", companyID).Return(&models.Company{}, errors.New("company not found")).Once()
		req, _ := http.NewRequest("DELETE", "/companies/"+companyID, nil)
		req = req.WithContext(middleware.SetUserID(req.Context(), uuid.New().String()))
		rr := httptest.NewRecorder()

		handler.Delete(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.Contains(t, rr.Body.String(), "Company not found")

		mockRepo.AssertExpectations(t)
		mockProducer.AssertNotCalled(t, "PublishCompanyDeleted", mock.Anything)
	})

	t.Run("Delete Unauthorized Access", func(t *testing.T) {
		handler, mockRepo, mockProducer := newTestCompanyHandler()

		companyID := uuid.New().String()

		mockRepo.AssertNotCalled(t, "GetByID", mock.Anything)
		mockRepo.AssertNotCalled(t, "Delete", mock.Anything)
		mockProducer.AssertNotCalled(t, "PublishCompanyDeleted", mock.Anything)

		req, _ := http.NewRequest("DELETE", "/companies/"+companyID, nil)
		rr := httptest.NewRecorder()

		handler.Patch(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Contains(t, rr.Body.String(), "Unauthorized")
	})
}
