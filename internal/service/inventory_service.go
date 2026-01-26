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
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrExpiredStock      = errors.New("stock has expired")
)

// InventoryService handles inventory business logic
type InventoryService struct {
	inventoryRepo *repository.InventoryRepository
}

// NewInventoryService creates a new inventory service
func NewInventoryService(inventoryRepo *repository.InventoryRepository) *InventoryService {
	return &InventoryService{inventoryRepo: inventoryRepo}
}

// AddStock adds inventory
func (s *InventoryService) AddStock(req *dto.CreateInventoryRequest) (*dto.InventoryResponse, error) {
	expiryDate, err := time.Parse("2006-01-02", req.ExpiryDate)
	if err != nil {
		return nil, fmt.Errorf("invalid expiry date format: %w", err)
	}

	receivedDate, err := time.Parse("2006-01-02", req.ReceivedDate)
	if err != nil {
		return nil, fmt.Errorf("invalid received date format: %w", err)
	}

	// Validate expiry date is in the future
	if expiryDate.Before(time.Now()) {
		return nil, ErrExpiredStock
	}

	inventory := &domain.Inventory{
		MedicationID: req.MedicationID,
		BatchNumber:  req.BatchNumber,
		ExpiryDate:   expiryDate,
		Quantity:     req.Quantity,
		Unit:         domain.InventoryUnit(req.Unit),
		CostPrice:    req.CostPrice,
		Supplier:     req.Supplier,
		ReceivedDate: receivedDate,
	}

	if err := s.inventoryRepo.Create(inventory); err != nil {
		return nil, fmt.Errorf("failed to create inventory: %w", err)
	}

	// Reload to get medication
	inventory, _ = s.inventoryRepo.FindByID(inventory.ID)
	return s.toInventoryResponse(inventory), nil
}

// GetMedicationStock gets stock levels for a medication
func (s *InventoryService) GetMedicationStock(medicationID uint) ([]*dto.InventoryListItem, error) {
	inventories, err := s.inventoryRepo.FindByMedicationID(medicationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get medication stock: %w", err)
	}

	items := make([]*dto.InventoryListItem, len(inventories))
	for i, inv := range inventories {
		items[i] = &dto.InventoryListItem{
			ID:             inv.ID,
			MedicationName: inv.Medication.Name,
			BatchNumber:    inv.BatchNumber,
			ExpiryDate:     inv.ExpiryDate.Format("2006-01-02"),
			Quantity:       inv.Quantity,
			Unit:           string(inv.Unit),
		}
	}
	return items, nil
}

// GetLowStockAlerts gets low stock alerts
func (s *InventoryService) GetLowStockAlerts(threshold int) ([]*dto.InventoryListItem, error) {
	inventories, err := s.inventoryRepo.FindLowStock(threshold)
	if err != nil {
		return nil, fmt.Errorf("failed to get low stock alerts: %w", err)
	}

	items := make([]*dto.InventoryListItem, len(inventories))
	for i, inv := range inventories {
		items[i] = &dto.InventoryListItem{
			ID:             inv.ID,
			MedicationName: inv.Medication.Name,
			BatchNumber:    inv.BatchNumber,
			ExpiryDate:     inv.ExpiryDate.Format("2006-01-02"),
			Quantity:       inv.Quantity,
			Unit:           string(inv.Unit),
		}
	}
	return items, nil
}

// GetExpiringSoonAlerts gets expiring soon alerts
func (s *InventoryService) GetExpiringSoonAlerts(days int) ([]*dto.InventoryListItem, error) {
	inventories, err := s.inventoryRepo.FindExpiringSoon(days)
	if err != nil {
		return nil, fmt.Errorf("failed to get expiring soon alerts: %w", err)
	}

	items := make([]*dto.InventoryListItem, len(inventories))
	for i, inv := range inventories {
		items[i] = &dto.InventoryListItem{
			ID:             inv.ID,
			MedicationName: inv.Medication.Name,
			BatchNumber:    inv.BatchNumber,
			ExpiryDate:     inv.ExpiryDate.Format("2006-01-02"),
			Quantity:       inv.Quantity,
			Unit:           string(inv.Unit),
		}
	}
	return items, nil
}

// Helper functions
func (s *InventoryService) toInventoryResponse(inv *domain.Inventory) *dto.InventoryResponse {
	resp := &dto.InventoryResponse{
		ID:           inv.ID,
		MedicationID: inv.MedicationID,
		BatchNumber:  inv.BatchNumber,
		ExpiryDate:   inv.ExpiryDate.Format("2006-01-02"),
		Quantity:     inv.Quantity,
		Unit:         string(inv.Unit),
		CostPrice:    inv.CostPrice,
		Supplier:     inv.Supplier,
		ReceivedDate: inv.ReceivedDate.Format("2006-01-02"),
		CreatedAt:    inv.CreatedAt,
	}
	if inv.Medication != nil {
		resp.MedicationName = inv.Medication.Name
	}
	return resp
}
