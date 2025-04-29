package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"xm-exercise/internal/utils"

	"xm-exercise/internal/api/middleware"
	"xm-exercise/internal/db"
	"xm-exercise/internal/events"
	"xm-exercise/internal/logger"
	"xm-exercise/pkg/models"
)

// CompanyHandler handles company-related requests
type CompanyHandler struct {
	companyRepo db.CompanyRepositoryInterface
	producer    events.KafkaProducerInterface
}

// NewCompanyHandler creates a new company handler
func NewCompanyHandler(
	companyRepo db.CompanyRepositoryInterface,
	producer events.KafkaProducerInterface,
) *CompanyHandler {
	return &CompanyHandler{
		companyRepo: companyRepo,
		producer:    producer,
	}
}

// Create godoc
// @Summary Create a new company
// @Description Create a new company with the provided details.
// @Tags companies
// @Accept json
// @Produce json
// @Param company body models.CompanyCreateRequest true "Company details"
// @Param Authorization header string true "Bearer token" example:"Bearer {token}"
// @Success 201 {object} models.CompanyResponse "Company created successfully"
// @Failure 400 {string} string "Invalid request body or validation error"
// @Failure 401 {string} string "Unauthorized"
// @Failure 409 {string} string "Company name already exists"
// @Failure 500 {string} string "Internal server error"
// @Security Bearer
// @Router /companies [post]
func (h *CompanyHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.WithContext(ctx)

	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		log.Warn("Unauthorized company creation attempt")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var companyCreateReq models.CompanyCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&companyCreateReq); err != nil {
		log.Error("Failed to decode request body", zap.Error(err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := companyCreateReq.Validate(); err != nil {
		log.Warn("Company validation failed",
			zap.Error(err),
			zap.String("company_name", companyCreateReq.Name),
			zap.String("created_by", userID),
		)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	now := time.Now().UTC()
	company := models.Company{
		ID:            uuid.New().String(),
		Name:          companyCreateReq.Name,
		Description:   companyCreateReq.Description,
		EmployeeCount: companyCreateReq.EmployeeCount,
		Registered:    companyCreateReq.Registered,
		Type:          companyCreateReq.Type,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	exists, err := h.companyRepo.ExistsByName(company.Name)
	if err != nil {
		log.Error("Failed to check name for uniqueness", zap.Error(err), zap.String("name", company.Name))
		http.Error(w, "Error checking name for uniqueness", http.StatusInternalServerError)
		return
	}

	if exists {
		log.Warn("Company name already exists", zap.String("name", company.Name))
		http.Error(w, "Company name already exists", http.StatusConflict)
		return
	}

	if err := h.companyRepo.Create(&company); err != nil {
		log.Error("Failed to create company",
			zap.Error(err),
			zap.String("company_name", company.Name),
			zap.String("created_by", userID),
		)
		http.Error(w, "Error creating company", http.StatusInternalServerError)
		return
	}

	if err := h.producer.PublishCompanyCreated(company); err != nil {
		log.Error("Failed to publish company created event",
			zap.Error(err),
			zap.String("company_name", company.Name),
		)
	} else {
		log.Info("Company created event published",
			zap.String("company_id", company.ID),
			zap.String("company_name", company.Name),
		)
	}

	log.Info("Company created",
		zap.String("company_id", company.ID),
		zap.String("company_name", company.Name),
		zap.String("created_by", userID),
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err = json.NewEncoder(w).Encode(company.ToResponse()); err != nil {
		log.Error("Failed to encode response data",
			zap.Error(err),
		)
	}
}

// Get godoc
// @Summary Get a company by ID
// @Description Get detailed information about a company by its ID
// @Tags companies
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token" example:"Bearer {token}"
// @Param id path string true "Company ID" format(uuid)
// @Success 200 {object} models.CompanyResponse "Company found"
// @Failure 400 {string} string "Invalid company ID"
// @Failure 404 {string} string "Company not found"
// @Security Bearer
// @Router /companies/{id} [get]
func (h *CompanyHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.WithContext(ctx)
	id := utils.ExtractIDFromPath(r)
	company, err := h.companyRepo.GetByID(id)
	if err != nil {
		http.Error(w, "Company not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(company.ToResponse()); err != nil {
		log.Error("Failed to encode response data",
			zap.Error(err),
		)
	}
}

// Patch godoc
// @Summary Update a company
// @Description Update specific fields of a company
// @Tags companies
// @Accept json
// @Produce json
// @Param id path string true "Company ID" format(uuid)
// @Param company body object true "Fields to update"
// @Success 200 {object} models.CompanyUpdateRequest "Company updated successfully"
// @Failure 400 {string} string "Invalid request body or validation error"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Company not found"
// @Failure 409 {string} string "Company name already exists"
// @Failure 500 {string} string "Internal server error"
// @Security Bearer
// @Router /companies/{id} [patch]
func (h *CompanyHandler) Patch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.WithContext(ctx)

	_, ok := middleware.GetUserID(ctx)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id := utils.ExtractIDFromPath(r)
	existingCompany, err := h.companyRepo.GetByID(id)
	if err != nil {
		http.Error(w, "Company not found", http.StatusNotFound)
		return
	}

	var updates models.CompanyUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if err := updates.Validate(); err != nil {
		log.Warn("Company update validation failed",
			zap.Error(err),
		)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if updates.Name != nil && *updates.Name != existingCompany.Name {
		exists, err := h.companyRepo.ExistsByName(*updates.Name)
		if err != nil {
			http.Error(w, "Error checking name uniqueness", http.StatusInternalServerError)
			return
		}
		if exists {
			http.Error(w, "Company name already exists", http.StatusConflict)
			return
		}
	}

	if updates.Name != nil {
		existingCompany.Name = *updates.Name
	}
	if updates.Description != nil {
		existingCompany.Description = updates.Description
	}
	if updates.EmployeeCount != nil {
		existingCompany.EmployeeCount = *updates.EmployeeCount
	}
	if updates.Registered != nil {
		existingCompany.Registered = updates.Registered
	}
	if updates.Type != "" {
		existingCompany.Type = updates.Type
	}
	existingCompany.UpdatedAt = time.Now().UTC()

	if err := h.companyRepo.Update(existingCompany); err != nil {
		http.Error(w, "Error updating company", http.StatusInternalServerError)
		return
	}

	if err := h.producer.PublishCompanyUpdated(existingCompany); err != nil {
		log.Error("Failed to publish company updated event",
			zap.Error(err),
			zap.String("company_name", existingCompany.Name),
		)
	} else {
		log.Info("Company updated event published",
			zap.String("company_id", existingCompany.ID),
			zap.String("company_name", existingCompany.Name),
		)
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(existingCompany.ToResponse()); err != nil {
		log.Error("Failed to encode response data",
			zap.Error(err),
		)
	}
}

// Delete godoc
// @Summary Delete a company
// @Description Delete a company by its ID
// @Tags companies
// @Accept json
// @Produce json
// @Param id path string true "Company ID" format(uuid)
// @Success 204 {string} string "Company deleted successfully"
// @Failure 400 {string} string "Invalid company ID"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Company not found"
// @Failure 500 {string} string "Internal server error"
// @Security Bearer
// @Router /companies/{id} [delete]
func (h *CompanyHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.WithContext(ctx)
	_, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id := utils.ExtractIDFromPath(r)

	existingCompany, err := h.companyRepo.GetByID(id)
	if err != nil {
		http.Error(w, "Company not found", http.StatusNotFound)
		return
	}

	if err := h.companyRepo.Delete(id); err != nil {
		http.Error(w, "Error deleting company", http.StatusInternalServerError)
		return
	}

	if err := h.producer.PublishCompanyDeleted(id); err != nil {
		log.Error("Failed to publish company deleted event",
			zap.Error(err),
			zap.String("company_name", existingCompany.Name),
		)
	} else {
		log.Info("Company deleted event published",
			zap.String("company_id", existingCompany.ID),
			zap.String("company_name", existingCompany.Name),
		)
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(existingCompany.ToResponse()); err != nil {
		log.Error("Failed to encode response data",
			zap.Error(err),
		)
	}
}
