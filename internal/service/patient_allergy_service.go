package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/minhtran/his/internal/domain"
	"github.com/minhtran/his/internal/dto"
	"github.com/minhtran/his/internal/repository"
)

var (
	ErrAllergyNotFound = errors.New("allergy not found")
)

// PatientAllergyService handles patient allergy business logic
type PatientAllergyService struct {
	allergyRepo *repository.PatientAllergyRepository
	patientRepo *repository.PatientRepository
}

// NewPatientAllergyService creates a new patient allergy service
func NewPatientAllergyService(allergyRepo *repository.PatientAllergyRepository, patientRepo *repository.PatientRepository) *PatientAllergyService {
	return &PatientAllergyService{
		allergyRepo: allergyRepo,
		patientRepo: patientRepo,
	}
}

// AddAllergy adds a new allergy to a patient
func (s *PatientAllergyService) AddAllergy(patientID uint, req *dto.CreateAllergyRequest, createdBy uint) (*dto.AllergyResponse, error) {
	// Verify patient exists
	patient, err := s.patientRepo.FindByID(patientID)
	if err != nil {
		return nil, fmt.Errorf("failed to find patient: %w", err)
	}
	if patient == nil {
		return nil, ErrPatientNotFound
	}

	// Parse diagnosed date if provided
	var diagnosedDate *time.Time
	if req.DiagnosedDate != "" {
		parsed, err := time.Parse("2006-01-02", req.DiagnosedDate)
		if err != nil {
			return nil, ErrInvalidDateFormat
		}
		diagnosedDate = &parsed
	}

	allergy := &domain.PatientAllergy{
		PatientID:     patientID,
		Allergen:      req.Allergen,
		AllergenType:  domain.AllergenType(req.AllergenType),
		Reaction:      req.Reaction,
		Severity:      domain.AllergySeverity(req.Severity),
		DiagnosedDate: diagnosedDate,
		Notes:         req.Notes,
		IsActive:      true,
		CreatedBy:     createdBy,
	}

	if err := s.allergyRepo.Create(allergy); err != nil {
		return nil, fmt.Errorf("failed to create allergy: %w", err)
	}

	return s.toAllergyResponse(allergy), nil
}

// UpdateAllergy updates an allergy record
func (s *PatientAllergyService) UpdateAllergy(id uint, req *dto.UpdateAllergyRequest, updatedBy uint) (*dto.AllergyResponse, error) {
	allergy, err := s.allergyRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find allergy: %w", err)
	}
	if allergy == nil {
		return nil, ErrAllergyNotFound
	}

	// Update fields
	if req.Allergen != "" {
		allergy.Allergen = req.Allergen
	}
	if req.AllergenType != "" {
		allergy.AllergenType = domain.AllergenType(req.AllergenType)
	}
	if req.Reaction != "" {
		allergy.Reaction = req.Reaction
	}
	if req.Severity != "" {
		allergy.Severity = domain.AllergySeverity(req.Severity)
	}
	if req.DiagnosedDate != "" {
		parsed, err := time.Parse("2006-01-02", req.DiagnosedDate)
		if err != nil {
			return nil, ErrInvalidDateFormat
		}
		allergy.DiagnosedDate = &parsed
	}
	if req.Notes != "" {
		allergy.Notes = req.Notes
	}
	if req.IsActive != nil {
		allergy.IsActive = *req.IsActive
	}

	allergy.UpdatedBy = updatedBy

	if err := s.allergyRepo.Update(allergy); err != nil {
		return nil, fmt.Errorf("failed to update allergy: %w", err)
	}

	return s.toAllergyResponse(allergy), nil
}

// GetPatientAllergies gets all allergies for a patient
func (s *PatientAllergyService) GetPatientAllergies(patientID uint) ([]*dto.AllergyListItem, error) {
	allergies, err := s.allergyRepo.FindByPatientID(patientID)
	if err != nil {
		return nil, fmt.Errorf("failed to get allergies: %w", err)
	}

	items := make([]*dto.AllergyListItem, len(allergies))
	for i, allergy := range allergies {
		items[i] = s.toAllergyListItem(allergy)
	}

	return items, nil
}

// GetActiveAllergies gets active allergies for a patient
func (s *PatientAllergyService) GetActiveAllergies(patientID uint) ([]*dto.AllergyListItem, error) {
	allergies, err := s.allergyRepo.FindActiveByPatientID(patientID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active allergies: %w", err)
	}

	items := make([]*dto.AllergyListItem, len(allergies))
	for i, allergy := range allergies {
		items[i] = s.toAllergyListItem(allergy)
	}

	return items, nil
}

// GetAllergyByID gets allergy details by ID
func (s *PatientAllergyService) GetAllergyByID(id uint) (*dto.AllergyResponse, error) {
	allergy, err := s.allergyRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find allergy: %w", err)
	}
	if allergy == nil {
		return nil, ErrAllergyNotFound
	}

	return s.toAllergyResponse(allergy), nil
}

// DeleteAllergy soft deletes an allergy
func (s *PatientAllergyService) DeleteAllergy(id uint) error {
	allergy, err := s.allergyRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("failed to find allergy: %w", err)
	}
	if allergy == nil {
		return ErrAllergyNotFound
	}

	if err := s.allergyRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete allergy: %w", err)
	}

	return nil
}

// Helper to convert domain to DTO
func (s *PatientAllergyService) toAllergyResponse(allergy *domain.PatientAllergy) *dto.AllergyResponse {
	return &dto.AllergyResponse{
		ID:            allergy.ID,
		PatientID:     allergy.PatientID,
		Allergen:      allergy.Allergen,
		AllergenType:  string(allergy.AllergenType),
		Reaction:      allergy.Reaction,
		Severity:      string(allergy.Severity),
		DiagnosedDate: allergy.DiagnosedDate,
		Notes:         allergy.Notes,
		IsActive:      allergy.IsActive,
		CreatedAt:     allergy.CreatedAt,
		UpdatedAt:     allergy.UpdatedAt,
	}
}

// Helper to convert domain to list item
func (s *PatientAllergyService) toAllergyListItem(allergy *domain.PatientAllergy) *dto.AllergyListItem {
	return &dto.AllergyListItem{
		ID:           allergy.ID,
		Allergen:     allergy.Allergen,
		AllergenType: string(allergy.AllergenType),
		Severity:     string(allergy.Severity),
		IsActive:     allergy.IsActive,
	}
}
