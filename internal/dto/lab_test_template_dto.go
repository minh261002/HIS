package dto

import "time"

// TemplateParameterResponse represents a template parameter
type TemplateParameterResponse struct {
	ID              uint    `json:"id"`
	ParameterName   string  `json:"parameter_name"`
	Unit            string  `json:"unit"`
	NormalRangeMin  float64 `json:"normal_range_min"`
	NormalRangeMax  float64 `json:"normal_range_max"`
	NormalRangeText string  `json:"normal_range_text"`
	DisplayOrder    int     `json:"display_order"`
}

// LabTestTemplateResponse represents lab test template details
type LabTestTemplateResponse struct {
	ID                      uint                         `json:"id"`
	Code                    string                       `json:"code"`
	Name                    string                       `json:"name"`
	Category                string                       `json:"category"`
	Description             string                       `json:"description"`
	SampleType              string                       `json:"sample_type"`
	PreparationInstructions string                       `json:"preparation_instructions"`
	TurnaroundTimeHours     int                          `json:"turnaround_time_hours"`
	Price                   float64                      `json:"price"`
	Parameters              []*TemplateParameterResponse `json:"parameters"`
	CreatedAt               time.Time                    `json:"created_at"`
}

// LabTestTemplateListItem represents simplified template for list view
type LabTestTemplateListItem struct {
	ID         uint    `json:"id"`
	Code       string  `json:"code"`
	Name       string  `json:"name"`
	Category   string  `json:"category"`
	SampleType string  `json:"sample_type"`
	Price      float64 `json:"price"`
}
