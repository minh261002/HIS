package service

import (
	"fmt"

	"github.com/minhtran/his/internal/dto"
	"github.com/minhtran/his/internal/repository"
)

// ICD10CodeService handles ICD-10 code business logic
type ICD10CodeService struct {
	icd10Repo *repository.ICD10CodeRepository
}

// NewICD10CodeService creates a new ICD-10 code service
func NewICD10CodeService(icd10Repo *repository.ICD10CodeRepository) *ICD10CodeService {
	return &ICD10CodeService{icd10Repo: icd10Repo}
}

// SearchICD10Codes searches ICD-10 codes
func (s *ICD10CodeService) SearchICD10Codes(query string, limit int) ([]*dto.ICD10CodeListItem, error) {
	codes, err := s.icd10Repo.Search(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search ICD-10 codes: %w", err)
	}

	items := make([]*dto.ICD10CodeListItem, len(codes))
	for i, code := range codes {
		items[i] = &dto.ICD10CodeListItem{
			ID:          code.ID,
			Code:        code.Code,
			Description: code.Description,
			Category:    code.Category,
		}
	}
	return items, nil
}

// GetICD10CodeByCode gets ICD-10 code by code
func (s *ICD10CodeService) GetICD10CodeByCode(code string) (*dto.ICD10CodeResponse, error) {
	icd10Code, err := s.icd10Repo.FindByCode(code)
	if err != nil {
		return nil, fmt.Errorf("failed to find ICD-10 code: %w", err)
	}
	if icd10Code == nil {
		return nil, ErrICD10CodeNotFound
	}

	return &dto.ICD10CodeResponse{
		ID:          icd10Code.ID,
		Code:        icd10Code.Code,
		Description: icd10Code.Description,
		Category:    icd10Code.Category,
	}, nil
}

// GetICD10CodesByCategory gets ICD-10 codes by category
func (s *ICD10CodeService) GetICD10CodesByCategory(category string) ([]*dto.ICD10CodeListItem, error) {
	codes, err := s.icd10Repo.FindByCategory(category)
	if err != nil {
		return nil, fmt.Errorf("failed to get ICD-10 codes by category: %w", err)
	}

	items := make([]*dto.ICD10CodeListItem, len(codes))
	for i, code := range codes {
		items[i] = &dto.ICD10CodeListItem{
			ID:          code.ID,
			Code:        code.Code,
			Description: code.Description,
			Category:    code.Category,
		}
	}
	return items, nil
}
