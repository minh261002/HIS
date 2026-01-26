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
	ErrBedNotFound       = errors.New("bed not found")
	ErrAdmissionNotFound = errors.New("admission not found")
	ErrBedNotAvailable   = errors.New("bed is not available")
)

// AdmissionService handles admission business logic
type AdmissionService struct {
	admissionRepo   *repository.AdmissionRepository
	allocationRepo  *repository.BedAllocationRepository
	bedRepo         *repository.BedRepository
	visitRepo       *repository.VisitRepository
	nursingNoteRepo *repository.NursingNoteRepository
}

// NewAdmissionService creates a new admission service
func NewAdmissionService(
	admissionRepo *repository.AdmissionRepository,
	allocationRepo *repository.BedAllocationRepository,
	bedRepo *repository.BedRepository,
	visitRepo *repository.VisitRepository,
	nursingNoteRepo *repository.NursingNoteRepository,
) *AdmissionService {
	return &AdmissionService{
		admissionRepo:   admissionRepo,
		allocationRepo:  allocationRepo,
		bedRepo:         bedRepo,
		visitRepo:       visitRepo,
		nursingNoteRepo: nursingNoteRepo,
	}
}

// CreateAdmission creates an admission with optional bed allocation
func (s *AdmissionService) CreateAdmission(req *dto.CreateAdmissionRequest, createdBy uint) (*dto.AdmissionResponse, error) {
	// Validate visit exists
	visit, err := s.visitRepo.FindByID(req.VisitID)
	if err != nil {
		return nil, fmt.Errorf("failed to find visit: %w", err)
	}
	if visit == nil {
		return nil, ErrVisitNotFound
	}

	// Generate admission code
	code, err := s.admissionRepo.GenerateAdmissionCode()
	if err != nil {
		return nil, fmt.Errorf("failed to generate admission code: %w", err)
	}

	// Create admission
	admission := &domain.Admission{
		AdmissionCode:      code,
		VisitID:            req.VisitID,
		PatientID:          visit.PatientID,
		DoctorID:           visit.DoctorID,
		AdmissionDate:      time.Now(),
		AdmissionDiagnosis: req.AdmissionDiagnosis,
		Status:             domain.AdmissionStatusAdmitted,
		CreatedBy:          createdBy,
	}

	if err := s.admissionRepo.Create(admission); err != nil {
		return nil, fmt.Errorf("failed to create admission: %w", err)
	}

	// Allocate bed if provided
	if req.BedID != nil {
		if err := s.allocateBed(admission.ID, *req.BedID, "", createdBy); err != nil {
			return nil, fmt.Errorf("failed to allocate bed: %w", err)
		}
	}

	// Reload to get relationships
	admission, _ = s.admissionRepo.FindByID(admission.ID)
	return s.toAdmissionResponse(admission), nil
}

// DischargeAdmission discharges a patient
func (s *AdmissionService) DischargeAdmission(id uint, req *dto.DischargeAdmissionRequest, updatedBy uint) error {
	admission, err := s.admissionRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("failed to find admission: %w", err)
	}
	if admission == nil {
		return ErrAdmissionNotFound
	}

	// Release current bed if any
	currentAllocation, _ := s.allocationRepo.FindCurrentAllocation(id)
	if currentAllocation != nil {
		if err := s.allocationRepo.ReleaseAllocation(currentAllocation.ID); err != nil {
			return fmt.Errorf("failed to release bed: %w", err)
		}
		// Update bed status to available
		s.bedRepo.UpdateStatus(currentAllocation.BedID, domain.BedStatusAvailable)
	}

	// Update admission
	now := time.Now()
	admission.DischargeDate = &now
	admission.DischargeDiagnosis = req.DischargeDiagnosis
	admission.DischargeSummary = req.DischargeSummary
	admission.Status = domain.AdmissionStatusDischarged
	admission.UpdatedBy = updatedBy

	return s.admissionRepo.Update(admission)
}

// TransferBed transfers patient to a new bed
func (s *AdmissionService) TransferBed(admissionID uint, req *dto.TransferBedRequest, createdBy uint) error {
	// Release current bed
	currentAllocation, err := s.allocationRepo.FindCurrentAllocation(admissionID)
	if err != nil {
		return fmt.Errorf("failed to find current allocation: %w", err)
	}
	if currentAllocation != nil {
		if err := s.allocationRepo.ReleaseAllocation(currentAllocation.ID); err != nil {
			return fmt.Errorf("failed to release current bed: %w", err)
		}
		// Update old bed status to available
		s.bedRepo.UpdateStatus(currentAllocation.BedID, domain.BedStatusAvailable)
	}

	// Allocate new bed
	return s.allocateBed(admissionID, req.NewBedID, req.Notes, createdBy)
}

// allocateBed allocates a bed to an admission
func (s *AdmissionService) allocateBed(admissionID, bedID uint, notes string, createdBy uint) error {
	// Validate bed exists and is available
	bed, err := s.bedRepo.FindByID(bedID)
	if err != nil {
		return fmt.Errorf("failed to find bed: %w", err)
	}
	if bed == nil {
		return ErrBedNotFound
	}
	if bed.Status != domain.BedStatusAvailable {
		return ErrBedNotAvailable
	}

	// Create allocation
	allocation := &domain.BedAllocation{
		AdmissionID:   admissionID,
		BedID:         bedID,
		AllocatedDate: time.Now(),
		IsCurrent:     true,
		Notes:         notes,
		CreatedBy:     createdBy,
	}

	if err := s.allocationRepo.Create(allocation); err != nil {
		return fmt.Errorf("failed to create allocation: %w", err)
	}

	// Update bed status to occupied
	return s.bedRepo.UpdateStatus(bedID, domain.BedStatusOccupied)
}

// CreateNursingNote creates a nursing note
func (s *AdmissionService) CreateNursingNote(admissionID uint, req *dto.CreateNursingNoteRequest, nurseID uint) error {
	note := &domain.NursingNote{
		AdmissionID:   admissionID,
		NurseID:       nurseID,
		NoteDate:      time.Now(),
		VitalSigns:    req.VitalSigns,
		Observations:  req.Observations,
		Interventions: req.Interventions,
	}
	return s.nursingNoteRepo.Create(note)
}

// GetAdmissionByID gets admission by ID
func (s *AdmissionService) GetAdmissionByID(id uint) (*dto.AdmissionResponse, error) {
	admission, err := s.admissionRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find admission: %w", err)
	}
	if admission == nil {
		return nil, ErrAdmissionNotFound
	}
	return s.toAdmissionResponse(admission), nil
}

// GetAdmissionByCode gets admission by code
func (s *AdmissionService) GetAdmissionByCode(code string) (*dto.AdmissionResponse, error) {
	admission, err := s.admissionRepo.FindByCode(code)
	if err != nil {
		return nil, fmt.Errorf("failed to find admission: %w", err)
	}
	if admission == nil {
		return nil, ErrAdmissionNotFound
	}
	return s.toAdmissionResponse(admission), nil
}

// GetActiveAdmissions gets all active admissions
func (s *AdmissionService) GetActiveAdmissions() ([]*dto.AdmissionListItem, error) {
	admissions, err := s.admissionRepo.FindActiveAdmissions()
	if err != nil {
		return nil, fmt.Errorf("failed to get active admissions: %w", err)
	}

	items := make([]*dto.AdmissionListItem, len(admissions))
	for i, a := range admissions {
		items[i] = s.toAdmissionListItem(a)
	}
	return items, nil
}

// GetPatientAdmissions gets patient's admission history
func (s *AdmissionService) GetPatientAdmissions(patientID uint) ([]*dto.AdmissionListItem, error) {
	admissions, err := s.admissionRepo.FindByPatientID(patientID)
	if err != nil {
		return nil, fmt.Errorf("failed to get patient admissions: %w", err)
	}

	items := make([]*dto.AdmissionListItem, len(admissions))
	for i, a := range admissions {
		items[i] = s.toAdmissionListItem(a)
	}
	return items, nil
}

// GetAdmissionNursingNotes gets admission's nursing notes
func (s *AdmissionService) GetAdmissionNursingNotes(admissionID uint) ([]*dto.NursingNoteResponse, error) {
	notes, err := s.nursingNoteRepo.FindByAdmissionID(admissionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get nursing notes: %w", err)
	}

	responses := make([]*dto.NursingNoteResponse, len(notes))
	for i, n := range notes {
		responses[i] = &dto.NursingNoteResponse{
			ID:            n.ID,
			NurseID:       n.NurseID,
			NoteDate:      n.NoteDate,
			VitalSigns:    n.VitalSigns,
			Observations:  n.Observations,
			Interventions: n.Interventions,
			CreatedAt:     n.CreatedAt,
		}
		if n.Nurse != nil {
			responses[i].NurseName = n.Nurse.FullName
		}
	}
	return responses, nil
}

// Helper functions
func (s *AdmissionService) toAdmissionResponse(a *domain.Admission) *dto.AdmissionResponse {
	resp := &dto.AdmissionResponse{
		ID:                 a.ID,
		AdmissionCode:      a.AdmissionCode,
		VisitID:            a.VisitID,
		PatientID:          a.PatientID,
		DoctorID:           a.DoctorID,
		AdmissionDate:      a.AdmissionDate,
		DischargeDate:      a.DischargeDate,
		AdmissionDiagnosis: a.AdmissionDiagnosis,
		DischargeDiagnosis: a.DischargeDiagnosis,
		DischargeSummary:   a.DischargeSummary,
		Status:             string(a.Status),
		CreatedAt:          a.CreatedAt,
		UpdatedAt:          a.UpdatedAt,
	}

	if a.Patient != nil {
		resp.PatientName = a.Patient.FullName
	}
	if a.Doctor != nil {
		resp.DoctorName = a.Doctor.FullName
	}

	// Add bed allocations
	resp.BedAllocations = make([]*dto.BedAllocationResponse, len(a.BedAllocations))
	for i, alloc := range a.BedAllocations {
		resp.BedAllocations[i] = &dto.BedAllocationResponse{
			ID:            alloc.ID,
			AllocatedDate: alloc.AllocatedDate,
			ReleasedDate:  alloc.ReleasedDate,
			IsCurrent:     alloc.IsCurrent,
			Notes:         alloc.Notes,
		}
		if alloc.Bed != nil {
			resp.BedAllocations[i].BedNumber = alloc.Bed.BedNumber
			resp.BedAllocations[i].Ward = alloc.Bed.Ward
		}
	}

	return resp
}

func (s *AdmissionService) toAdmissionListItem(a *domain.Admission) *dto.AdmissionListItem {
	item := &dto.AdmissionListItem{
		ID:            a.ID,
		AdmissionCode: a.AdmissionCode,
		AdmissionDate: a.AdmissionDate,
		Status:        string(a.Status),
	}

	if a.Patient != nil {
		item.PatientName = a.Patient.FullName
	}

	// Get current bed
	for _, alloc := range a.BedAllocations {
		if alloc.IsCurrent && alloc.Bed != nil {
			item.CurrentBedNumber = alloc.Bed.BedNumber
			break
		}
	}

	return item
}
