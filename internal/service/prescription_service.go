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
	ErrMedicationNotFound = errors.New("medication not found")
)

// PrescriptionService handles prescription business logic
type PrescriptionService struct {
	prescriptionRepo     *repository.PrescriptionRepository
	prescriptionItemRepo *repository.PrescriptionItemRepository
	medicationRepo       *repository.MedicationRepository
	visitRepo            *repository.VisitRepository
}

// NewPrescriptionService creates a new prescription service
func NewPrescriptionService(
	prescriptionRepo *repository.PrescriptionRepository,
	prescriptionItemRepo *repository.PrescriptionItemRepository,
	medicationRepo *repository.MedicationRepository,
	visitRepo *repository.VisitRepository,
) *PrescriptionService {
	return &PrescriptionService{
		prescriptionRepo:     prescriptionRepo,
		prescriptionItemRepo: prescriptionItemRepo,
		medicationRepo:       medicationRepo,
		visitRepo:            visitRepo,
	}
}

// CreatePrescription creates a prescription with items
func (s *PrescriptionService) CreatePrescription(req *dto.CreatePrescriptionRequest, prescribedBy uint) (*dto.PrescriptionResponse, error) {
	// Validate visit exists
	visit, err := s.visitRepo.FindByID(req.VisitID)
	if err != nil {
		return nil, fmt.Errorf("failed to find visit: %w", err)
	}
	if visit == nil {
		return nil, ErrVisitNotFound
	}

	// Validate medications exist
	for _, item := range req.Items {
		medication, err := s.medicationRepo.FindByID(item.MedicationID)
		if err != nil {
			return nil, fmt.Errorf("failed to find medication: %w", err)
		}
		if medication == nil {
			return nil, ErrMedicationNotFound
		}
	}

	// Generate prescription code
	code, err := s.prescriptionRepo.GeneratePrescriptionCode()
	if err != nil {
		return nil, fmt.Errorf("failed to generate prescription code: %w", err)
	}

	// Create prescription with items
	prescription := &domain.Prescription{
		PrescriptionCode: code,
		VisitID:          req.VisitID,
		PatientID:        visit.PatientID,
		DoctorID:         visit.DoctorID,
		DiagnosisID:      req.DiagnosisID,
		Status:           domain.PrescriptionStatusPending,
		PrescribedDate:   time.Now(),
		Notes:            req.Notes,
		CreatedBy:        prescribedBy,
	}

	// Add items
	prescription.Items = make([]*domain.PrescriptionItem, len(req.Items))
	for i, itemReq := range req.Items {
		prescription.Items[i] = &domain.PrescriptionItem{
			MedicationID: itemReq.MedicationID,
			Quantity:     itemReq.Quantity,
			Dosage:       itemReq.Dosage,
			Frequency:    itemReq.Frequency,
			DurationDays: itemReq.DurationDays,
			Instructions: itemReq.Instructions,
		}
	}

	if err := s.prescriptionRepo.Create(prescription); err != nil {
		return nil, fmt.Errorf("failed to create prescription: %w", err)
	}

	// Reload to get relationships
	prescription, _ = s.prescriptionRepo.FindByID(prescription.ID)
	return s.toPrescriptionResponse(prescription), nil
}

// UpdatePrescription updates prescription
func (s *PrescriptionService) UpdatePrescription(id uint, req *dto.UpdatePrescriptionRequest, updatedBy uint) (*dto.PrescriptionResponse, error) {
	prescription, err := s.prescriptionRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find prescription: %w", err)
	}
	if prescription == nil {
		return nil, ErrPrescriptionNotFound
	}

	if req.Notes != "" {
		prescription.Notes = req.Notes
	}
	prescription.UpdatedBy = updatedBy

	if err := s.prescriptionRepo.Update(prescription); err != nil {
		return nil, fmt.Errorf("failed to update prescription: %w", err)
	}

	prescription, _ = s.prescriptionRepo.FindByID(prescription.ID)
	return s.toPrescriptionResponse(prescription), nil
}

// DispensePrescription marks prescription as dispensed
func (s *PrescriptionService) DispensePrescription(id uint, updatedBy uint) error {
	prescription, err := s.prescriptionRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("failed to find prescription: %w", err)
	}
	if prescription == nil {
		return ErrPrescriptionNotFound
	}

	prescription.Status = domain.PrescriptionStatusDispensed
	prescription.UpdatedBy = updatedBy
	return s.prescriptionRepo.Update(prescription)
}

// CompletePrescription marks prescription as completed
func (s *PrescriptionService) CompletePrescription(id uint, updatedBy uint) error {
	prescription, err := s.prescriptionRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("failed to find prescription: %w", err)
	}
	if prescription == nil {
		return ErrPrescriptionNotFound
	}

	prescription.Status = domain.PrescriptionStatusCompleted
	prescription.UpdatedBy = updatedBy
	return s.prescriptionRepo.Update(prescription)
}

// CancelPrescription cancels prescription
func (s *PrescriptionService) CancelPrescription(id uint, updatedBy uint) error {
	prescription, err := s.prescriptionRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("failed to find prescription: %w", err)
	}
	if prescription == nil {
		return ErrPrescriptionNotFound
	}

	prescription.Status = domain.PrescriptionStatusCancelled
	prescription.UpdatedBy = updatedBy
	return s.prescriptionRepo.Update(prescription)
}

// GetPrescriptionByID gets prescription by ID
func (s *PrescriptionService) GetPrescriptionByID(id uint) (*dto.PrescriptionResponse, error) {
	prescription, err := s.prescriptionRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find prescription: %w", err)
	}
	if prescription == nil {
		return nil, ErrPrescriptionNotFound
	}
	return s.toPrescriptionResponse(prescription), nil
}

// GetPrescriptionByCode gets prescription by code
func (s *PrescriptionService) GetPrescriptionByCode(code string) (*dto.PrescriptionResponse, error) {
	prescription, err := s.prescriptionRepo.FindByCode(code)
	if err != nil {
		return nil, fmt.Errorf("failed to find prescription: %w", err)
	}
	if prescription == nil {
		return nil, ErrPrescriptionNotFound
	}
	return s.toPrescriptionResponse(prescription), nil
}

// GetVisitPrescriptions gets visit prescriptions
func (s *PrescriptionService) GetVisitPrescriptions(visitID uint) ([]*dto.PrescriptionListItem, error) {
	prescriptions, err := s.prescriptionRepo.FindByVisitID(visitID)
	if err != nil {
		return nil, fmt.Errorf("failed to get prescriptions: %w", err)
	}

	items := make([]*dto.PrescriptionListItem, len(prescriptions))
	for i, p := range prescriptions {
		items[i] = s.toPrescriptionListItem(p)
	}
	return items, nil
}

// GetPatientPrescriptions gets patient prescriptions
func (s *PrescriptionService) GetPatientPrescriptions(patientID uint, filters map[string]interface{}) ([]*dto.PrescriptionListItem, error) {
	prescriptions, err := s.prescriptionRepo.FindByPatientID(patientID, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get prescriptions: %w", err)
	}

	items := make([]*dto.PrescriptionListItem, len(prescriptions))
	for i, p := range prescriptions {
		items[i] = s.toPrescriptionListItem(p)
	}
	return items, nil
}

// Helper functions
func (s *PrescriptionService) toPrescriptionResponse(p *domain.Prescription) *dto.PrescriptionResponse {
	resp := &dto.PrescriptionResponse{
		ID:               p.ID,
		PrescriptionCode: p.PrescriptionCode,
		VisitID:          p.VisitID,
		PatientID:        p.PatientID,
		DoctorID:         p.DoctorID,
		DiagnosisID:      p.DiagnosisID,
		Status:           string(p.Status),
		PrescribedDate:   p.PrescribedDate.Format("2006-01-02"),
		Notes:            p.Notes,
		CreatedAt:        p.CreatedAt,
		UpdatedAt:        p.UpdatedAt,
	}

	if p.Patient != nil {
		resp.PatientName = p.Patient.FullName
	}
	if p.Doctor != nil {
		resp.DoctorName = p.Doctor.FullName
	}
	if p.Diagnosis != nil && p.Diagnosis.ICD10Code != nil {
		resp.DiagnosisCode = p.Diagnosis.ICD10Code.Code
		resp.DiagnosisDesc = p.Diagnosis.ICD10Code.Description
	}

	// Add items
	resp.Items = make([]*dto.PrescriptionItemResponse, len(p.Items))
	for i, item := range p.Items {
		resp.Items[i] = &dto.PrescriptionItemResponse{
			ID:           item.ID,
			MedicationID: item.MedicationID,
			Quantity:     item.Quantity,
			Dosage:       item.Dosage,
			Frequency:    item.Frequency,
			DurationDays: item.DurationDays,
			Instructions: item.Instructions,
		}
		if item.Medication != nil {
			resp.Items[i].MedicationName = item.Medication.Name
			resp.Items[i].DosageForm = string(item.Medication.DosageForm)
			resp.Items[i].Strength = item.Medication.Strength
		}
	}

	return resp
}

func (s *PrescriptionService) toPrescriptionListItem(p *domain.Prescription) *dto.PrescriptionListItem {
	item := &dto.PrescriptionListItem{
		ID:               p.ID,
		PrescriptionCode: p.PrescriptionCode,
		Status:           string(p.Status),
		PrescribedDate:   p.PrescribedDate.Format("2006-01-02"),
		ItemCount:        len(p.Items),
		CreatedAt:        p.CreatedAt,
	}

	if p.Patient != nil {
		item.PatientName = p.Patient.FullName
	}
	if p.Doctor != nil {
		item.DoctorName = p.Doctor.FullName
	}

	return item
}
