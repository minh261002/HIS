package service

import (
	"fmt"

	"github.com/minhtran/his/internal/dto"
	"github.com/minhtran/his/internal/repository"
)

// MedicationService handles medication business logic
type MedicationService struct {
	medicationRepo *repository.MedicationRepository
}

// NewMedicationService creates a new medication service
func NewMedicationService(medicationRepo *repository.MedicationRepository) *MedicationService {
	return &MedicationService{medicationRepo: medicationRepo}
}

// SearchMedications searches medications
func (s *MedicationService) SearchMedications(query string, limit int) ([]*dto.MedicationListItem, error) {
	medications, err := s.medicationRepo.Search(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search medications: %w", err)
	}

	items := make([]*dto.MedicationListItem, len(medications))
	for i, med := range medications {
		items[i] = &dto.MedicationListItem{
			ID:          med.ID,
			Name:        med.Name,
			GenericName: med.GenericName,
			DosageForm:  string(med.DosageForm),
			Strength:    med.Strength,
		}
	}
	return items, nil
}

// GetMedicationByID gets medication by ID
func (s *MedicationService) GetMedicationByID(id uint) (*dto.MedicationResponse, error) {
	medication, err := s.medicationRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find medication: %w", err)
	}
	if medication == nil {
		return nil, ErrMedicationNotFound
	}

	return &dto.MedicationResponse{
		ID:           medication.ID,
		Name:         medication.Name,
		GenericName:  medication.GenericName,
		DosageForm:   string(medication.DosageForm),
		Strength:     medication.Strength,
		Unit:         medication.Unit,
		Manufacturer: medication.Manufacturer,
	}, nil
}
