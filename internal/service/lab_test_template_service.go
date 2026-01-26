package service

import (
	"fmt"

	"github.com/minhtran/his/internal/dto"
	"github.com/minhtran/his/internal/repository"
)

// LabTestTemplateService handles lab test template business logic
type LabTestTemplateService struct {
	templateRepo *repository.LabTestTemplateRepository
}

// NewLabTestTemplateService creates a new lab test template service
func NewLabTestTemplateService(templateRepo *repository.LabTestTemplateRepository) *LabTestTemplateService {
	return &LabTestTemplateService{templateRepo: templateRepo}
}

// SearchTemplates searches lab test templates
func (s *LabTestTemplateService) SearchTemplates(query string, limit int) ([]*dto.LabTestTemplateListItem, error) {
	templates, err := s.templateRepo.Search(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search templates: %w", err)
	}

	items := make([]*dto.LabTestTemplateListItem, len(templates))
	for i, t := range templates {
		items[i] = &dto.LabTestTemplateListItem{
			ID:         t.ID,
			Code:       t.Code,
			Name:       t.Name,
			Category:   string(t.Category),
			SampleType: string(t.SampleType),
			Price:      t.Price,
		}
	}
	return items, nil
}

// GetTemplateByCode gets template by code
func (s *LabTestTemplateService) GetTemplateByCode(code string) (*dto.LabTestTemplateResponse, error) {
	template, err := s.templateRepo.FindByCode(code)
	if err != nil {
		return nil, fmt.Errorf("failed to find template: %w", err)
	}
	if template == nil {
		return nil, ErrLabTestTemplateNotFound
	}

	resp := &dto.LabTestTemplateResponse{
		ID:                      template.ID,
		Code:                    template.Code,
		Name:                    template.Name,
		Category:                string(template.Category),
		Description:             template.Description,
		SampleType:              string(template.SampleType),
		PreparationInstructions: template.PreparationInstructions,
		TurnaroundTimeHours:     template.TurnaroundTimeHours,
		Price:                   template.Price,
		CreatedAt:               template.CreatedAt,
	}

	// Add parameters
	resp.Parameters = make([]*dto.TemplateParameterResponse, len(template.Parameters))
	for i, p := range template.Parameters {
		resp.Parameters[i] = &dto.TemplateParameterResponse{
			ID:              p.ID,
			ParameterName:   p.ParameterName,
			Unit:            p.Unit,
			NormalRangeMin:  p.NormalRangeMin,
			NormalRangeMax:  p.NormalRangeMax,
			NormalRangeText: p.NormalRangeText,
			DisplayOrder:    p.DisplayOrder,
		}
	}

	return resp, nil
}

// GetTemplatesByCategory gets templates by category
func (s *LabTestTemplateService) GetTemplatesByCategory(category string) ([]*dto.LabTestTemplateListItem, error) {
	templates, err := s.templateRepo.FindByCategory(category)
	if err != nil {
		return nil, fmt.Errorf("failed to get templates: %w", err)
	}

	items := make([]*dto.LabTestTemplateListItem, len(templates))
	for i, t := range templates {
		items[i] = &dto.LabTestTemplateListItem{
			ID:         t.ID,
			Code:       t.Code,
			Name:       t.Name,
			Category:   string(t.Category),
			SampleType: string(t.SampleType),
			Price:      t.Price,
		}
	}
	return items, nil
}
