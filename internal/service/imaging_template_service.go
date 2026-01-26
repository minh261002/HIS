package service

import (
	"fmt"

	"github.com/minhtran/his/internal/dto"
	"github.com/minhtran/his/internal/repository"
)

// ImagingTemplateService handles imaging template business logic
type ImagingTemplateService struct {
	templateRepo *repository.ImagingTemplateRepository
}

// NewImagingTemplateService creates a new imaging template service
func NewImagingTemplateService(templateRepo *repository.ImagingTemplateRepository) *ImagingTemplateService {
	return &ImagingTemplateService{templateRepo: templateRepo}
}

// SearchTemplates searches imaging templates
func (s *ImagingTemplateService) SearchTemplates(query string, limit int) ([]*dto.ImagingTemplateListItem, error) {
	templates, err := s.templateRepo.Search(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search templates: %w", err)
	}

	items := make([]*dto.ImagingTemplateListItem, len(templates))
	for i, t := range templates {
		items[i] = &dto.ImagingTemplateListItem{
			ID:       t.ID,
			Code:     t.Code,
			Name:     t.Name,
			Modality: string(t.Modality),
			BodyPart: string(t.BodyPart),
			Price:    t.Price,
		}
	}
	return items, nil
}

// GetTemplateByCode gets template by code
func (s *ImagingTemplateService) GetTemplateByCode(code string) (*dto.ImagingTemplateResponse, error) {
	template, err := s.templateRepo.FindByCode(code)
	if err != nil {
		return nil, fmt.Errorf("failed to find template: %w", err)
	}
	if template == nil {
		return nil, ErrImagingTemplateNotFound
	}

	return &dto.ImagingTemplateResponse{
		ID:              template.ID,
		Code:            template.Code,
		Name:            template.Name,
		Modality:        string(template.Modality),
		BodyPart:        string(template.BodyPart),
		Description:     template.Description,
		TemplateContent: template.TemplateContent,
		Price:           template.Price,
		CreatedAt:       template.CreatedAt,
	}, nil
}

// GetTemplatesByModality gets templates by modality
func (s *ImagingTemplateService) GetTemplatesByModality(modality string) ([]*dto.ImagingTemplateListItem, error) {
	templates, err := s.templateRepo.FindByModality(modality)
	if err != nil {
		return nil, fmt.Errorf("failed to get templates: %w", err)
	}

	items := make([]*dto.ImagingTemplateListItem, len(templates))
	for i, t := range templates {
		items[i] = &dto.ImagingTemplateListItem{
			ID:       t.ID,
			Code:     t.Code,
			Name:     t.Name,
			Modality: string(t.Modality),
			BodyPart: string(t.BodyPart),
			Price:    t.Price,
		}
	}
	return items, nil
}
