package service

import (
	"fmt"

	"github.com/minhtran/his/internal/dto"
	"github.com/minhtran/his/internal/repository"
)

// BedService handles bed business logic
type BedService struct {
	bedRepo *repository.BedRepository
}

// NewBedService creates a new bed service
func NewBedService(bedRepo *repository.BedRepository) *BedService {
	return &BedService{bedRepo: bedRepo}
}

// GetAvailableBeds gets available beds by department and bed type
func (s *BedService) GetAvailableBeds(department, bedType string) ([]*dto.BedListItem, error) {
	beds, err := s.bedRepo.FindAvailableBeds(department, bedType)
	if err != nil {
		return nil, fmt.Errorf("failed to get available beds: %w", err)
	}

	items := make([]*dto.BedListItem, len(beds))
	for i, b := range beds {
		items[i] = &dto.BedListItem{
			ID:         b.ID,
			BedNumber:  b.BedNumber,
			Department: string(b.Department),
			Ward:       b.Ward,
			BedType:    string(b.BedType),
			Status:     string(b.Status),
		}
	}
	return items, nil
}

// GetBedByNumber gets bed by bed number
func (s *BedService) GetBedByNumber(bedNumber string) (*dto.BedResponse, error) {
	bed, err := s.bedRepo.FindByBedNumber(bedNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to find bed: %w", err)
	}
	if bed == nil {
		return nil, ErrBedNotFound
	}

	return &dto.BedResponse{
		ID:         bed.ID,
		BedNumber:  bed.BedNumber,
		Department: string(bed.Department),
		Ward:       bed.Ward,
		BedType:    string(bed.BedType),
		Status:     string(bed.Status),
		CreatedAt:  bed.CreatedAt,
	}, nil
}
