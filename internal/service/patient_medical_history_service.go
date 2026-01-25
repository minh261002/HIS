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
	ErrMedicalHistoryNotFound = errors.New("medical history not found")
)

// PatientMedicalHistoryService handles patient medical history business logic
type PatientMedicalHistoryService struct {
	historyRepo *repository.PatientMedicalHistoryRepository
	patientRepo *repository.PatientRepository
}

// NewPatientMedicalHistoryService creates a new patient medical history service
func NewPatientMedicalHistoryService(historyRepo *repository.PatientMedicalHistoryRepository, patientRepo *repository.PatientRepository) *PatientMedicalHistoryService {
	return &PatientMedicalHistoryService{
		historyRepo: historyRepo,
		patientRepo: patientRepo,
	}
}

// AddMedicalHistory adds a new medical history to a patient
func (s *PatientMedicalHistoryService) AddMedicalHistory(patientID uint, req *dto.CreateMedicalHistoryRequest, createdBy uint) (*dto.MedicalHistoryResponse, error) {
	// Verify patient exists
	patient, err := s.patientRepo.FindByID(patientID)
	if err != nil {
		return nil, fmt.Errorf("failed to find patient: %w", err)
	}
	if patient == nil {
		return nil, ErrPatientNotFound
	}

	// Parse diagnosis date if provided
	var diagnosisDate *time.Time
	if req.DiagnosisDate != "" {
		parsed, err := time.Parse("2006-01-02", req.DiagnosisDate)
		if err != nil {
			return nil, ErrInvalidDateFormat
		}
		diagnosisDate = &parsed
	}

	history := &domain.PatientMedicalHistory{
		PatientID:     patientID,
		ConditionName: req.ConditionName,
		ConditionType: domain.ConditionType(req.ConditionType),
		DiagnosisDate: diagnosisDate,
		Status:        domain.ConditionStatus(req.Status),
		Treatment:     req.Treatment,
		Notes:         req.Notes,
		IsActive:      true,
		CreatedBy:     createdBy,
	}

	if err := s.historyRepo.Create(history); err != nil {
		return nil, fmt.Errorf("failed to create medical history: %w", err)
	}

	return s.toMedicalHistoryResponse(history), nil
}

// UpdateMedicalHistory updates a medical history record
func (s *PatientMedicalHistoryService) UpdateMedicalHistory(id uint, req *dto.UpdateMedicalHistoryRequest, updatedBy uint) (*dto.MedicalHistoryResponse, error) {
	history, err := s.historyRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find medical history: %w", err)
	}
	if history == nil {
		return nil, ErrMedicalHistoryNotFound
	}

	// Update fields
	if req.ConditionName != "" {
		history.ConditionName = req.ConditionName
	}
	if req.ConditionType != "" {
		history.ConditionType = domain.ConditionType(req.ConditionType)
	}
	if req.DiagnosisDate != "" {
		parsed, err := time.Parse("2006-01-02", req.DiagnosisDate)
		if err != nil {
			return nil, ErrInvalidDateFormat
		}
		history.DiagnosisDate = &parsed
	}
	if req.Status != "" {
		history.Status = domain.ConditionStatus(req.Status)
	}
	if req.Treatment != "" {
		history.Treatment = req.Treatment
	}
	if req.Notes != "" {
		history.Notes = req.Notes
	}
	if req.IsActive != nil {
		history.IsActive = *req.IsActive
	}

	history.UpdatedBy = updatedBy

	if err := s.historyRepo.Update(history); err != nil {
		return nil, fmt.Errorf("failed to update medical history: %w", err)
	}

	return s.toMedicalHistoryResponse(history), nil
}

// GetPatientHistory gets all medical history for a patient
func (s *PatientMedicalHistoryService) GetPatientHistory(patientID uint) ([]*dto.MedicalHistoryListItem, error) {
	histories, err := s.historyRepo.FindByPatientID(patientID)
	if err != nil {
		return nil, fmt.Errorf("failed to get medical history: %w", err)
	}

	items := make([]*dto.MedicalHistoryListItem, len(histories))
	for i, history := range histories {
		items[i] = s.toMedicalHistoryListItem(history)
	}

	return items, nil
}

// GetActiveConditions gets active medical conditions for a patient
func (s *PatientMedicalHistoryService) GetActiveConditions(patientID uint) ([]*dto.MedicalHistoryListItem, error) {
	histories, err := s.historyRepo.FindActiveByPatientID(patientID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active conditions: %w", err)
	}

	items := make([]*dto.MedicalHistoryListItem, len(histories))
	for i, history := range histories {
		items[i] = s.toMedicalHistoryListItem(history)
	}

	return items, nil
}

// GetMedicalHistoryByID gets medical history details by ID
func (s *PatientMedicalHistoryService) GetMedicalHistoryByID(id uint) (*dto.MedicalHistoryResponse, error) {
	history, err := s.historyRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find medical history: %w", err)
	}
	if history == nil {
		return nil, ErrMedicalHistoryNotFound
	}

	return s.toMedicalHistoryResponse(history), nil
}

// DeleteMedicalHistory soft deletes a medical history record
func (s *PatientMedicalHistoryService) DeleteMedicalHistory(id uint) error {
	history, err := s.historyRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("failed to find medical history: %w", err)
	}
	if history == nil {
		return ErrMedicalHistoryNotFound
	}

	if err := s.historyRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete medical history: %w", err)
	}

	return nil
}

// Helper to convert domain to DTO
func (s *PatientMedicalHistoryService) toMedicalHistoryResponse(history *domain.PatientMedicalHistory) *dto.MedicalHistoryResponse {
	return &dto.MedicalHistoryResponse{
		ID:            history.ID,
		PatientID:     history.PatientID,
		ConditionName: history.ConditionName,
		ConditionType: string(history.ConditionType),
		DiagnosisDate: history.DiagnosisDate,
		Status:        string(history.Status),
		Treatment:     history.Treatment,
		Notes:         history.Notes,
		IsActive:      history.IsActive,
		CreatedAt:     history.CreatedAt,
		UpdatedAt:     history.UpdatedAt,
	}
}

// Helper to convert domain to list item
func (s *PatientMedicalHistoryService) toMedicalHistoryListItem(history *domain.PatientMedicalHistory) *dto.MedicalHistoryListItem {
	return &dto.MedicalHistoryListItem{
		ID:            history.ID,
		ConditionName: history.ConditionName,
		ConditionType: string(history.ConditionType),
		Status:        string(history.Status),
		IsActive:      history.IsActive,
	}
}
