package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/minhtran/his/internal/domain"
	"github.com/minhtran/his/internal/dto"
	"github.com/minhtran/his/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrPrescriptionNotFound = errors.New("prescription not found")
)

// DispensingService handles dispensing business logic
type DispensingService struct {
	dispensingRepo   *repository.DispensingRepository
	inventoryRepo    *repository.InventoryRepository
	prescriptionRepo *repository.PrescriptionRepository
	db               *gorm.DB
}

// NewDispensingService creates a new dispensing service
func NewDispensingService(
	dispensingRepo *repository.DispensingRepository,
	inventoryRepo *repository.InventoryRepository,
	prescriptionRepo *repository.PrescriptionRepository,
	db *gorm.DB,
) *DispensingService {
	return &DispensingService{
		dispensingRepo:   dispensingRepo,
		inventoryRepo:    inventoryRepo,
		prescriptionRepo: prescriptionRepo,
		db:               db,
	}
}

// DispensePrescription dispenses prescription with stock deduction
func (s *DispensingService) DispensePrescription(req *dto.DispensePrescriptionRequest, pharmacistID uint) ([]*dto.DispensingResponse, error) {
	// Validate prescription exists
	prescription, err := s.prescriptionRepo.FindByID(req.PrescriptionID)
	if err != nil {
		return nil, fmt.Errorf("failed to find prescription: %w", err)
	}
	if prescription == nil {
		return nil, ErrPrescriptionNotFound
	}

	// Start transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var dispensings []*domain.Dispensing

	for _, item := range req.Items {
		// Get inventory
		inventory, err := s.inventoryRepo.FindByID(item.InventoryID)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to find inventory: %w", err)
		}
		if inventory == nil {
			tx.Rollback()
			return nil, fmt.Errorf("inventory not found")
		}

		// Check stock availability
		if inventory.Quantity < item.Quantity {
			tx.Rollback()
			return nil, ErrInsufficientStock
		}

		// Check expiry
		if inventory.ExpiryDate.Before(time.Now()) {
			tx.Rollback()
			return nil, ErrExpiredStock
		}

		// Generate dispensing code
		code, err := s.dispensingRepo.GenerateDispensingCode()
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to generate dispensing code: %w", err)
		}

		// Create dispensing record
		dispensing := &domain.Dispensing{
			DispensingCode:     code,
			PrescriptionID:     req.PrescriptionID,
			PrescriptionItemID: item.PrescriptionItemID,
			MedicationID:       inventory.MedicationID,
			InventoryID:        item.InventoryID,
			PatientID:          prescription.PatientID,
			PharmacistID:       pharmacistID,
			QuantityDispensed:  item.Quantity,
			BatchNumber:        inventory.BatchNumber,
			DispensedDate:      time.Now(),
		}

		if err := tx.Create(dispensing).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create dispensing: %w", err)
		}

		// Deduct from inventory
		newQuantity := inventory.Quantity - item.Quantity
		if err := tx.Model(&domain.Inventory{}).Where("id = ?", item.InventoryID).Update("quantity", newQuantity).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to update inventory: %w", err)
		}

		dispensings = append(dispensings, dispensing)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Reload to get relationships
	responses := make([]*dto.DispensingResponse, len(dispensings))
	for i, d := range dispensings {
		var dispensing domain.Dispensing
		err := s.db.Preload("Medication").Preload("Patient").Preload("Pharmacist").First(&dispensing, d.ID).Error
		if err == nil {
			responses[i] = s.toDispensingResponse(&dispensing)
		}
	}

	return responses, nil
}

// GetPrescriptionDispensingRecords gets prescription's dispensing records
func (s *DispensingService) GetPrescriptionDispensingRecords(prescriptionID uint) ([]*dto.DispensingListItem, error) {
	dispensings, err := s.dispensingRepo.FindByPrescriptionID(prescriptionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get dispensing records: %w", err)
	}

	items := make([]*dto.DispensingListItem, len(dispensings))
	for i, d := range dispensings {
		items[i] = s.toDispensingListItem(d)
	}
	return items, nil
}

// GetPatientDispensingHistory gets patient's dispensing history
func (s *DispensingService) GetPatientDispensingHistory(patientID uint) ([]*dto.DispensingListItem, error) {
	dispensings, err := s.dispensingRepo.FindByPatientID(patientID)
	if err != nil {
		return nil, fmt.Errorf("failed to get patient dispensing history: %w", err)
	}

	items := make([]*dto.DispensingListItem, len(dispensings))
	for i, d := range dispensings {
		items[i] = s.toDispensingListItem(d)
	}
	return items, nil
}

// Helper functions
func (s *DispensingService) toDispensingResponse(d *domain.Dispensing) *dto.DispensingResponse {
	resp := &dto.DispensingResponse{
		ID:                 d.ID,
		DispensingCode:     d.DispensingCode,
		PrescriptionID:     d.PrescriptionID,
		PrescriptionItemID: d.PrescriptionItemID,
		MedicationID:       d.MedicationID,
		PatientID:          d.PatientID,
		PharmacistID:       d.PharmacistID,
		QuantityDispensed:  d.QuantityDispensed,
		BatchNumber:        d.BatchNumber,
		DispensedDate:      d.DispensedDate,
		Notes:              d.Notes,
		CreatedAt:          d.CreatedAt,
	}
	if d.Medication != nil {
		resp.MedicationName = d.Medication.Name
	}
	if d.Patient != nil {
		resp.PatientName = d.Patient.FullName
	}
	if d.Pharmacist != nil {
		resp.PharmacistName = d.Pharmacist.FullName
	}
	return resp
}

func (s *DispensingService) toDispensingListItem(d *domain.Dispensing) *dto.DispensingListItem {
	item := &dto.DispensingListItem{
		ID:                d.ID,
		DispensingCode:    d.DispensingCode,
		QuantityDispensed: d.QuantityDispensed,
		DispensedDate:     d.DispensedDate,
	}
	if d.Medication != nil {
		item.MedicationName = d.Medication.Name
	}
	if d.Pharmacist != nil {
		item.PharmacistName = d.Pharmacist.FullName
	}
	return item
}
