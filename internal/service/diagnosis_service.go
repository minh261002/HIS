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
	ErrICD10CodeNotFound      = errors.New("ICD-10 code not found")
	ErrDiagnosisNotFound      = errors.New("diagnosis not found")
	ErrPrimaryDiagnosisExists = errors.New("primary diagnosis already exists for this visit")
)

// DiagnosisService handles diagnosis business logic
type DiagnosisService struct {
	diagnosisRepo *repository.DiagnosisRepository
	icd10Repo     *repository.ICD10CodeRepository
	visitRepo     *repository.VisitRepository
	patientRepo   *repository.PatientRepository
}

// NewDiagnosisService creates a new diagnosis service
func NewDiagnosisService(
	diagnosisRepo *repository.DiagnosisRepository,
	icd10Repo *repository.ICD10CodeRepository,
	visitRepo *repository.VisitRepository,
	patientRepo *repository.PatientRepository,
) *DiagnosisService {
	return &DiagnosisService{
		diagnosisRepo: diagnosisRepo,
		icd10Repo:     icd10Repo,
		visitRepo:     visitRepo,
		patientRepo:   patientRepo,
	}
}

// AddDiagnosis adds a diagnosis to a visit
func (s *DiagnosisService) AddDiagnosis(req *dto.CreateDiagnosisRequest, diagnosedBy uint) (*dto.DiagnosisResponse, error) {
	// Validate visit exists
	visit, err := s.visitRepo.FindByID(req.VisitID)
	if err != nil {
		return nil, fmt.Errorf("failed to find visit: %w", err)
	}
	if visit == nil {
		return nil, ErrVisitNotFound
	}

	// Validate ICD-10 code exists
	icd10Code, err := s.icd10Repo.FindByID(req.ICD10CodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to find ICD-10 code: %w", err)
	}
	if icd10Code == nil {
		return nil, ErrICD10CodeNotFound
	}

	// Check if primary diagnosis already exists
	if req.DiagnosisType == "PRIMARY" {
		existing, _ := s.diagnosisRepo.FindPrimaryDiagnosis(req.VisitID)
		if existing != nil {
			return nil, ErrPrimaryDiagnosisExists
		}
	}

	diagnosis := &domain.Diagnosis{
		VisitID:         req.VisitID,
		PatientID:       visit.PatientID,
		ICD10CodeID:     req.ICD10CodeID,
		DiagnosedBy:     diagnosedBy,
		DiagnosisType:   domain.DiagnosisType(req.DiagnosisType),
		DiagnosisStatus: domain.DiagnosisStatus(req.DiagnosisStatus),
		ClinicalNotes:   req.ClinicalNotes,
		DiagnosedAt:     time.Now(),
		CreatedBy:       diagnosedBy,
	}

	if err := s.diagnosisRepo.Create(diagnosis); err != nil {
		return nil, fmt.Errorf("failed to create diagnosis: %w", err)
	}

	// Reload to get relationships
	diagnosis, _ = s.diagnosisRepo.FindByID(diagnosis.ID)
	return s.toDiagnosisResponse(diagnosis), nil
}

// UpdateDiagnosis updates a diagnosis
func (s *DiagnosisService) UpdateDiagnosis(id uint, req *dto.UpdateDiagnosisRequest, updatedBy uint) (*dto.DiagnosisResponse, error) {
	diagnosis, err := s.diagnosisRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find diagnosis: %w", err)
	}
	if diagnosis == nil {
		return nil, ErrDiagnosisNotFound
	}

	if req.DiagnosisStatus != "" {
		diagnosis.DiagnosisStatus = domain.DiagnosisStatus(req.DiagnosisStatus)
	}
	if req.ClinicalNotes != "" {
		diagnosis.ClinicalNotes = req.ClinicalNotes
	}

	diagnosis.UpdatedBy = updatedBy

	if err := s.diagnosisRepo.Update(diagnosis); err != nil {
		return nil, fmt.Errorf("failed to update diagnosis: %w", err)
	}

	diagnosis, _ = s.diagnosisRepo.FindByID(diagnosis.ID)
	return s.toDiagnosisResponse(diagnosis), nil
}

// DeleteDiagnosis deletes a diagnosis
func (s *DiagnosisService) DeleteDiagnosis(id uint) error {
	diagnosis, err := s.diagnosisRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("failed to find diagnosis: %w", err)
	}
	if diagnosis == nil {
		return ErrDiagnosisNotFound
	}

	return s.diagnosisRepo.Delete(id)
}

// GetDiagnosisByID gets diagnosis by ID
func (s *DiagnosisService) GetDiagnosisByID(id uint) (*dto.DiagnosisResponse, error) {
	diagnosis, err := s.diagnosisRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find diagnosis: %w", err)
	}
	if diagnosis == nil {
		return nil, ErrDiagnosisNotFound
	}
	return s.toDiagnosisResponse(diagnosis), nil
}

// GetVisitDiagnoses gets all diagnoses for a visit
func (s *DiagnosisService) GetVisitDiagnoses(visitID uint) ([]*dto.DiagnosisListItem, error) {
	diagnoses, err := s.diagnosisRepo.FindByVisitID(visitID)
	if err != nil {
		return nil, fmt.Errorf("failed to get diagnoses: %w", err)
	}

	items := make([]*dto.DiagnosisListItem, len(diagnoses))
	for i, d := range diagnoses {
		items[i] = s.toDiagnosisListItem(d)
	}
	return items, nil
}

// GetPatientDiagnoses gets patient's diagnosis history
func (s *DiagnosisService) GetPatientDiagnoses(patientID uint, filters map[string]interface{}) ([]*dto.DiagnosisListItem, error) {
	diagnoses, err := s.diagnosisRepo.FindByPatientID(patientID, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get diagnoses: %w", err)
	}

	items := make([]*dto.DiagnosisListItem, len(diagnoses))
	for i, d := range diagnoses {
		items[i] = s.toDiagnosisListItem(d)
	}
	return items, nil
}

// Helper functions
func (s *DiagnosisService) toDiagnosisResponse(d *domain.Diagnosis) *dto.DiagnosisResponse {
	resp := &dto.DiagnosisResponse{
		ID:              d.ID,
		VisitID:         d.VisitID,
		PatientID:       d.PatientID,
		ICD10CodeID:     d.ICD10CodeID,
		DiagnosedBy:     d.DiagnosedBy,
		DiagnosisType:   string(d.DiagnosisType),
		DiagnosisStatus: string(d.DiagnosisStatus),
		ClinicalNotes:   d.ClinicalNotes,
		DiagnosedAt:     d.DiagnosedAt,
		CreatedAt:       d.CreatedAt,
		UpdatedAt:       d.UpdatedAt,
	}

	if d.ICD10Code != nil {
		resp.ICD10Code = d.ICD10Code.Code
		resp.ICD10Description = d.ICD10Code.Description
	}
	if d.Doctor != nil {
		resp.DoctorName = d.Doctor.FullName
	}
	if d.Patient != nil {
		resp.PatientName = d.Patient.FullName
	}

	return resp
}

func (s *DiagnosisService) toDiagnosisListItem(d *domain.Diagnosis) *dto.DiagnosisListItem {
	item := &dto.DiagnosisListItem{
		ID:              d.ID,
		DiagnosisType:   string(d.DiagnosisType),
		DiagnosisStatus: string(d.DiagnosisStatus),
		DiagnosedAt:     d.DiagnosedAt,
	}

	if d.ICD10Code != nil {
		item.ICD10Code = d.ICD10Code.Code
		item.ICD10Description = d.ICD10Code.Description
	}

	return item
}
