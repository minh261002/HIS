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
	ErrPatientExists     = errors.New("patient already exists")
	ErrPatientNotFound   = errors.New("patient not found")
	ErrInvalidDateFormat = errors.New("invalid date format, use YYYY-MM-DD")
)

// PatientService handles patient business logic
type PatientService struct {
	patientRepo *repository.PatientRepository
}

// NewPatientService creates a new patient service
func NewPatientService(patientRepo *repository.PatientRepository) *PatientService {
	return &PatientService{
		patientRepo: patientRepo,
	}
}

// RegisterPatient registers a new patient
func (s *PatientService) RegisterPatient(req *dto.CreatePatientRequest, createdBy uint) (*dto.PatientResponse, error) {
	// Check if patient with national ID already exists
	if req.NationalID != "" {
		existing, err := s.patientRepo.FindByNationalID(req.NationalID)
		if err != nil {
			return nil, fmt.Errorf("failed to check existing patient: %w", err)
		}
		if existing != nil {
			return nil, ErrPatientExists
		}
	}

	// Parse date of birth
	dob, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil {
		return nil, ErrInvalidDateFormat
	}

	// Generate patient code
	patientCode, err := s.patientRepo.GeneratePatientCode()
	if err != nil {
		return nil, fmt.Errorf("failed to generate patient code: %w", err)
	}

	// Create patient
	patient := &domain.Patient{
		PatientCode:                  patientCode,
		FirstName:                    req.FirstName,
		LastName:                     req.LastName,
		DateOfBirth:                  dob,
		Gender:                       domain.Gender(req.Gender),
		BloodType:                    domain.BloodType(req.BloodType),
		PhoneNumber:                  req.PhoneNumber,
		Email:                        req.Email,
		Address:                      req.Address,
		City:                         req.City,
		State:                        req.State,
		PostalCode:                   req.PostalCode,
		Country:                      req.Country,
		NationalID:                   req.NationalID,
		InsuranceNumber:              req.InsuranceNumber,
		InsuranceProvider:            req.InsuranceProvider,
		EmergencyContactName:         req.EmergencyContactName,
		EmergencyContactPhone:        req.EmergencyContactPhone,
		EmergencyContactRelationship: req.EmergencyContactRelationship,
		Allergies:                    req.Allergies,
		ChronicConditions:            req.ChronicConditions,
		Notes:                        req.Notes,
		IsActive:                     true,
		CreatedBy:                    createdBy,
	}

	if patient.Country == "" {
		patient.Country = "Vietnam"
	}

	if err := s.patientRepo.Create(patient); err != nil {
		return nil, fmt.Errorf("failed to create patient: %w", err)
	}

	return s.toPatientResponse(patient), nil
}

// UpdatePatient updates patient information
func (s *PatientService) UpdatePatient(id uint, req *dto.UpdatePatientRequest, updatedBy uint) (*dto.PatientResponse, error) {
	patient, err := s.patientRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find patient: %w", err)
	}
	if patient == nil {
		return nil, ErrPatientNotFound
	}

	// Update fields
	if req.FirstName != "" {
		patient.FirstName = req.FirstName
	}
	if req.LastName != "" {
		patient.LastName = req.LastName
	}
	if req.DateOfBirth != "" {
		dob, err := time.Parse("2006-01-02", req.DateOfBirth)
		if err != nil {
			return nil, ErrInvalidDateFormat
		}
		patient.DateOfBirth = dob
	}
	if req.Gender != "" {
		patient.Gender = domain.Gender(req.Gender)
	}
	if req.BloodType != "" {
		patient.BloodType = domain.BloodType(req.BloodType)
	}
	if req.PhoneNumber != "" {
		patient.PhoneNumber = req.PhoneNumber
	}
	if req.Email != "" {
		patient.Email = req.Email
	}
	if req.Address != "" {
		patient.Address = req.Address
	}
	if req.City != "" {
		patient.City = req.City
	}
	if req.State != "" {
		patient.State = req.State
	}
	if req.PostalCode != "" {
		patient.PostalCode = req.PostalCode
	}
	if req.Country != "" {
		patient.Country = req.Country
	}
	if req.NationalID != "" {
		patient.NationalID = req.NationalID
	}
	if req.InsuranceNumber != "" {
		patient.InsuranceNumber = req.InsuranceNumber
	}
	if req.InsuranceProvider != "" {
		patient.InsuranceProvider = req.InsuranceProvider
	}
	if req.EmergencyContactName != "" {
		patient.EmergencyContactName = req.EmergencyContactName
	}
	if req.EmergencyContactPhone != "" {
		patient.EmergencyContactPhone = req.EmergencyContactPhone
	}
	if req.EmergencyContactRelationship != "" {
		patient.EmergencyContactRelationship = req.EmergencyContactRelationship
	}
	if req.Allergies != "" {
		patient.Allergies = req.Allergies
	}
	if req.ChronicConditions != "" {
		patient.ChronicConditions = req.ChronicConditions
	}
	if req.Notes != "" {
		patient.Notes = req.Notes
	}
	if req.IsActive != nil {
		patient.IsActive = *req.IsActive
	}

	patient.UpdatedBy = updatedBy

	if err := s.patientRepo.Update(patient); err != nil {
		return nil, fmt.Errorf("failed to update patient: %w", err)
	}

	return s.toPatientResponse(patient), nil
}

// GetPatientByID gets patient by ID
func (s *PatientService) GetPatientByID(id uint) (*dto.PatientResponse, error) {
	patient, err := s.patientRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find patient: %w", err)
	}
	if patient == nil {
		return nil, ErrPatientNotFound
	}
	return s.toPatientResponse(patient), nil
}

// GetPatientByCode gets patient by patient code
func (s *PatientService) GetPatientByCode(code string) (*dto.PatientResponse, error) {
	patient, err := s.patientRepo.FindByPatientCode(code)
	if err != nil {
		return nil, fmt.Errorf("failed to find patient: %w", err)
	}
	if patient == nil {
		return nil, ErrPatientNotFound
	}
	return s.toPatientResponse(patient), nil
}

// SearchPatients searches patients
func (s *PatientService) SearchPatients(query string, page, pageSize int) ([]*dto.PatientListItem, int64, error) {
	patients, total, err := s.patientRepo.Search(query, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search patients: %w", err)
	}

	items := make([]*dto.PatientListItem, len(patients))
	for i, patient := range patients {
		items[i] = s.toPatientListItem(patient)
	}

	return items, total, nil
}

// ListPatients lists patients with filters
func (s *PatientService) ListPatients(page, pageSize int, filters map[string]interface{}) ([]*dto.PatientListItem, int64, error) {
	patients, total, err := s.patientRepo.List(page, pageSize, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list patients: %w", err)
	}

	items := make([]*dto.PatientListItem, len(patients))
	for i, patient := range patients {
		items[i] = s.toPatientListItem(patient)
	}

	return items, total, nil
}

// DeletePatient soft deletes a patient
func (s *PatientService) DeletePatient(id uint) error {
	patient, err := s.patientRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("failed to find patient: %w", err)
	}
	if patient == nil {
		return ErrPatientNotFound
	}

	if err := s.patientRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete patient: %w", err)
	}

	return nil
}

// GetPatientStats gets patient statistics
func (s *PatientService) GetPatientStats() (map[string]interface{}, error) {
	stats, err := s.patientRepo.GetPatientStats()
	if err != nil {
		return nil, fmt.Errorf("failed to get patient stats: %w", err)
	}
	return stats, nil
}

// Helper to convert domain patient to DTO
func (s *PatientService) toPatientResponse(patient *domain.Patient) *dto.PatientResponse {
	return &dto.PatientResponse{
		ID:                           patient.ID,
		PatientCode:                  patient.PatientCode,
		FirstName:                    patient.FirstName,
		LastName:                     patient.LastName,
		FullName:                     patient.FullName,
		DateOfBirth:                  patient.DateOfBirth.Format("2006-01-02"),
		Age:                          patient.Age,
		Gender:                       string(patient.Gender),
		BloodType:                    string(patient.BloodType),
		PhoneNumber:                  patient.PhoneNumber,
		Email:                        patient.Email,
		Address:                      patient.Address,
		City:                         patient.City,
		State:                        patient.State,
		PostalCode:                   patient.PostalCode,
		Country:                      patient.Country,
		NationalID:                   patient.NationalID,
		InsuranceNumber:              patient.InsuranceNumber,
		InsuranceProvider:            patient.InsuranceProvider,
		EmergencyContactName:         patient.EmergencyContactName,
		EmergencyContactPhone:        patient.EmergencyContactPhone,
		EmergencyContactRelationship: patient.EmergencyContactRelationship,
		Allergies:                    patient.Allergies,
		ChronicConditions:            patient.ChronicConditions,
		Notes:                        patient.Notes,
		IsActive:                     patient.IsActive,
		CreatedAt:                    patient.CreatedAt,
		UpdatedAt:                    patient.UpdatedAt,
	}
}

// Helper to convert domain patient to list item
func (s *PatientService) toPatientListItem(patient *domain.Patient) *dto.PatientListItem {
	return &dto.PatientListItem{
		ID:          patient.ID,
		PatientCode: patient.PatientCode,
		FullName:    patient.FullName,
		Age:         patient.Age,
		Gender:      string(patient.Gender),
		PhoneNumber: patient.PhoneNumber,
		City:        patient.City,
		IsActive:    patient.IsActive,
		CreatedAt:   patient.CreatedAt.Format(time.RFC3339),
	}
}
