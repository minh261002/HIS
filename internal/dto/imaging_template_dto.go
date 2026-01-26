package dto

import "time"

// ImagingTemplateResponse represents imaging template details
type ImagingTemplateResponse struct {
	ID              uint      `json:"id"`
	Code            string    `json:"code"`
	Name            string    `json:"name"`
	Modality        string    `json:"modality"`
	BodyPart        string    `json:"body_part"`
	Description     string    `json:"description"`
	TemplateContent string    `json:"template_content"`
	Price           float64   `json:"price"`
	CreatedAt       time.Time `json:"created_at"`
}

// ImagingTemplateListItem represents simplified template for list view
type ImagingTemplateListItem struct {
	ID       uint    `json:"id"`
	Code     string  `json:"code"`
	Name     string  `json:"name"`
	Modality string  `json:"modality"`
	BodyPart string  `json:"body_part"`
	Price    float64 `json:"price"`
}
