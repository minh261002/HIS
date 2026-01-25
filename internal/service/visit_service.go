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
	ErrVisitNotFound = errors.New("visit not found")
)

// VisitService handles visit business logic
type VisitService struct {
	visitRepo       *repository.VisitRepository
	patientRepo     *repository.PatientRepository
	userRepo        *repository.UserRepository
	appointmentRepo *repository.AppointmentRepository
}

// NewVisitService creates a new visit service
func NewVisitService(
	visitRepo *repository.VisitRepository,
	patientRepo *repository.PatientRepository,
	userRepo *repository.UserRepository,
	appointmentRepo *repository.AppointmentRepository,
) *VisitService {
	return &VisitService{
		visitRepo:       visitRepo,
		patientRepo:     patientRepo,
		userRepo:        userRepo,
		appointmentRepo: appointmentRepo,
	}
}

// CreateVisit creates a new visit
func (s *VisitService) CreateVisit(req *dto.CreateVisitRequest, createdBy uint) (*dto.VisitResponse, error) {
	// Validate patient exists
	patient, err := s.patientRepo.FindByID(req.PatientID)
	if err != nil {
		return nil, fmt.Errorf("failed to find patient: %w", err)
	}
	if patient == nil {
		return nil, ErrPatientNotFound
	}

	// Validate doctor exists
	doctor, err := s.userRepo.FindByID(req.DoctorID)
	if err != nil {
		return nil, fmt.Errorf("failed to find doctor: %w", err)
	}
	if doctor == nil {
		return nil, errors.New("doctor not found")
	}

	// Validate appointment if provided
	if req.AppointmentID != nil {
		appointment, err := s.appointmentRepo.FindByID(*req.AppointmentID)
		if err != nil {
			return nil, fmt.Errorf("failed to find appointment: %w", err)
		}
		if appointment == nil {
			return nil, errors.New("appointment not found")
		}
		if appointment.PatientID != req.PatientID {
			return nil, errors.New("appointment does not belong to this patient")
		}
	}

	// Generate visit code
	code, err := s.visitRepo.GenerateVisitCode()
	if err != nil {
		return nil, fmt.Errorf("failed to generate visit code: %w", err)
	}

	now := time.Now()
	visit := &domain.Visit{
		VisitCode:      code,
		AppointmentID:  req.AppointmentID,
		PatientID:      req.PatientID,
		DoctorID:       req.DoctorID,
		VisitDate:      now,
		VisitTime:      now,
		VisitType:      domain.VisitType(req.VisitType),
		Status:         domain.VisitStatusWaiting,
		ChiefComplaint: req.ChiefComplaint,
		CreatedBy:      createdBy,
	}

	// Set vital signs if provided
	if req.VitalSigns != nil {
		visit.Temperature = req.VitalSigns.Temperature
		visit.BloodPressureSystolic = req.VitalSigns.BloodPressureSystolic
		visit.BloodPressureDiastolic = req.VitalSigns.BloodPressureDiastolic
		visit.HeartRate = req.VitalSigns.HeartRate
		visit.RespiratoryRate = req.VitalSigns.RespiratoryRate
		visit.OxygenSaturation = req.VitalSigns.OxygenSaturation
		visit.Weight = req.VitalSigns.Weight
		visit.Height = req.VitalSigns.Height
	}

	if err := s.visitRepo.Create(visit); err != nil {
		return nil, fmt.Errorf("failed to create visit: %w", err)
	}

	// Update appointment status to IN_PROGRESS if linked
	if req.AppointmentID != nil {
		appointment, _ := s.appointmentRepo.FindByID(*req.AppointmentID)
		if appointment != nil {
			appointment.Status = domain.AppointmentStatusInProgress
			s.appointmentRepo.Update(appointment)
		}
	}

	// Reload to get relationships
	visit, _ = s.visitRepo.FindByID(visit.ID)
	return s.toVisitResponse(visit), nil
}

// UpdateVisit updates visit details
func (s *VisitService) UpdateVisit(id uint, req *dto.UpdateVisitRequest, updatedBy uint) (*dto.VisitResponse, error) {
	visit, err := s.visitRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find visit: %w", err)
	}
	if visit == nil {
		return nil, ErrVisitNotFound
	}

	// Update fields
	if req.Symptoms != "" {
		visit.Symptoms = req.Symptoms
	}
	if req.PhysicalExamination != "" {
		visit.PhysicalExamination = req.PhysicalExamination
	}
	if req.ClinicalNotes != "" {
		visit.ClinicalNotes = req.ClinicalNotes
	}
	if req.TreatmentPlan != "" {
		visit.TreatmentPlan = req.TreatmentPlan
	}
	if req.FollowUpInstructions != "" {
		visit.FollowUpInstructions = req.FollowUpInstructions
	}
	if req.NextVisitDate != "" {
		nextDate, err := time.Parse("2006-01-02", req.NextVisitDate)
		if err != nil {
			return nil, ErrInvalidDateFormat
		}
		visit.NextVisitDate = &nextDate
	}

	// Update vital signs if provided
	if req.VitalSigns != nil {
		if req.VitalSigns.Temperature > 0 {
			visit.Temperature = req.VitalSigns.Temperature
		}
		if req.VitalSigns.BloodPressureSystolic > 0 {
			visit.BloodPressureSystolic = req.VitalSigns.BloodPressureSystolic
		}
		if req.VitalSigns.BloodPressureDiastolic > 0 {
			visit.BloodPressureDiastolic = req.VitalSigns.BloodPressureDiastolic
		}
		if req.VitalSigns.HeartRate > 0 {
			visit.HeartRate = req.VitalSigns.HeartRate
		}
		if req.VitalSigns.RespiratoryRate > 0 {
			visit.RespiratoryRate = req.VitalSigns.RespiratoryRate
		}
		if req.VitalSigns.OxygenSaturation > 0 {
			visit.OxygenSaturation = req.VitalSigns.OxygenSaturation
		}
		if req.VitalSigns.Weight > 0 {
			visit.Weight = req.VitalSigns.Weight
		}
		if req.VitalSigns.Height > 0 {
			visit.Height = req.VitalSigns.Height
		}
	}

	visit.UpdatedBy = updatedBy

	if err := s.visitRepo.Update(visit); err != nil {
		return nil, fmt.Errorf("failed to update visit: %w", err)
	}

	visit, _ = s.visitRepo.FindByID(visit.ID)
	return s.toVisitResponse(visit), nil
}

// CompleteVisit marks visit as completed
func (s *VisitService) CompleteVisit(id uint, updatedBy uint) error {
	visit, err := s.visitRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("failed to find visit: %w", err)
	}
	if visit == nil {
		return ErrVisitNotFound
	}

	visit.Status = domain.VisitStatusCompleted
	visit.UpdatedBy = updatedBy

	if err := s.visitRepo.Update(visit); err != nil {
		return fmt.Errorf("failed to complete visit: %w", err)
	}

	// Update appointment status to COMPLETED if linked
	if visit.AppointmentID != nil {
		appointment, _ := s.appointmentRepo.FindByID(*visit.AppointmentID)
		if appointment != nil {
			appointment.Status = domain.AppointmentStatusCompleted
			s.appointmentRepo.Update(appointment)
		}
	}

	return nil
}

// CancelVisit cancels a visit
func (s *VisitService) CancelVisit(id uint, updatedBy uint) error {
	visit, err := s.visitRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("failed to find visit: %w", err)
	}
	if visit == nil {
		return ErrVisitNotFound
	}

	visit.Status = domain.VisitStatusCancelled
	visit.UpdatedBy = updatedBy

	return s.visitRepo.Update(visit)
}

// GetVisitByID gets visit by ID
func (s *VisitService) GetVisitByID(id uint) (*dto.VisitResponse, error) {
	visit, err := s.visitRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find visit: %w", err)
	}
	if visit == nil {
		return nil, ErrVisitNotFound
	}
	return s.toVisitResponse(visit), nil
}

// GetVisitByCode gets visit by code
func (s *VisitService) GetVisitByCode(code string) (*dto.VisitResponse, error) {
	visit, err := s.visitRepo.FindByCode(code)
	if err != nil {
		return nil, fmt.Errorf("failed to find visit: %w", err)
	}
	if visit == nil {
		return nil, ErrVisitNotFound
	}
	return s.toVisitResponse(visit), nil
}

// GetPatientVisits gets patient's visit history
func (s *VisitService) GetPatientVisits(patientID uint, filters map[string]interface{}) ([]*dto.VisitListItem, error) {
	visits, err := s.visitRepo.FindByPatientID(patientID, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get visits: %w", err)
	}

	items := make([]*dto.VisitListItem, len(visits))
	for i, visit := range visits {
		items[i] = s.toVisitListItem(visit)
	}
	return items, nil
}

// GetDoctorVisits gets doctor's visits for a date
func (s *VisitService) GetDoctorVisits(doctorID uint, dateStr string) ([]*dto.VisitListItem, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, ErrInvalidDateFormat
	}

	visits, err := s.visitRepo.FindByDoctorID(doctorID, date)
	if err != nil {
		return nil, fmt.Errorf("failed to get visits: %w", err)
	}

	items := make([]*dto.VisitListItem, len(visits))
	for i, visit := range visits {
		items[i] = s.toVisitListItem(visit)
	}
	return items, nil
}

// SearchVisits searches visits
func (s *VisitService) SearchVisits(filters map[string]interface{}, page, pageSize int) ([]*dto.VisitListItem, int64, error) {
	visits, total, err := s.visitRepo.Search(filters, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search visits: %w", err)
	}

	items := make([]*dto.VisitListItem, len(visits))
	for i, visit := range visits {
		items[i] = s.toVisitListItem(visit)
	}
	return items, total, nil
}

// Helper functions
func (s *VisitService) toVisitResponse(v *domain.Visit) *dto.VisitResponse {
	resp := &dto.VisitResponse{
		ID:                     v.ID,
		VisitCode:              v.VisitCode,
		AppointmentID:          v.AppointmentID,
		PatientID:              v.PatientID,
		DoctorID:               v.DoctorID,
		VisitDate:              v.VisitDate.Format("2006-01-02"),
		VisitTime:              v.VisitTime.Format("15:04"),
		VisitType:              string(v.VisitType),
		Status:                 string(v.Status),
		ChiefComplaint:         v.ChiefComplaint,
		Symptoms:               v.Symptoms,
		Temperature:            v.Temperature,
		BloodPressureSystolic:  v.BloodPressureSystolic,
		BloodPressureDiastolic: v.BloodPressureDiastolic,
		HeartRate:              v.HeartRate,
		RespiratoryRate:        v.RespiratoryRate,
		OxygenSaturation:       v.OxygenSaturation,
		Weight:                 v.Weight,
		Height:                 v.Height,
		BMI:                    v.BMI,
		PhysicalExamination:    v.PhysicalExamination,
		ClinicalNotes:          v.ClinicalNotes,
		TreatmentPlan:          v.TreatmentPlan,
		FollowUpInstructions:   v.FollowUpInstructions,
		NextVisitDate:          v.NextVisitDate,
		CreatedAt:              v.CreatedAt,
		UpdatedAt:              v.UpdatedAt,
	}

	if v.Patient != nil {
		resp.PatientName = v.Patient.FullName
	}
	if v.Doctor != nil {
		resp.DoctorName = v.Doctor.FullName
	}

	return resp
}

func (s *VisitService) toVisitListItem(v *domain.Visit) *dto.VisitListItem {
	item := &dto.VisitListItem{
		ID:             v.ID,
		VisitCode:      v.VisitCode,
		VisitDate:      v.VisitDate.Format("2006-01-02"),
		VisitTime:      v.VisitTime.Format("15:04"),
		VisitType:      string(v.VisitType),
		Status:         string(v.Status),
		ChiefComplaint: v.ChiefComplaint,
	}

	if v.Patient != nil {
		item.PatientName = v.Patient.FullName
	}
	if v.Doctor != nil {
		item.DoctorName = v.Doctor.FullName
	}

	return item
}
